package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	sp500TacticalETFCode            = "513650"
	sp500TotalReturnSymbol          = "^SP500TR"
	sp500TotalReturnDataURL         = "https://query2.finance.yahoo.com/v8/finance/chart/%5ESP500TR?range=max&interval=1d"
	sp500OfficialIndexURL           = "https://www.spglobal.com/spdji/en/indices/equity/sp-500/"
	sp500ForwardPEURL               = "https://historyofmarket.com/api/sp500/forward-pe.json"
	sp500VIXHistoryURL              = "https://cdn.cboe.com/api/global/us_indices/daily_prices/VIX_History.csv"
	sp500SinaFuturesURL             = "https://hq.sinajs.cn/list=hf_ES"
	sp500FuturesChartURL            = "https://query1.finance.yahoo.com/v8/finance/chart/ES%3DF?interval=5m&range=1d"
	sp500TacticalETFQuoteURL        = "https://quote.eastmoney.com/sh513650.html"
	sp500TacticalETFNetValuePageURL = "https://fund.eastmoney.com/513650.html"
	sp500SPYTotalReturnBackupURL    = "https://www.ssga.com/library-content/products/fund-data/etfs/us/spdr-etf-historical-distributions.xlsx"
)

type sp500OpportunitySnapshot struct {
	Date                    string
	Drawdown                float64
	CNYDrawdown             float64
	ForwardPE               float64
	ForwardPEPercentile     float64
	US10YBondYield          float64
	EarningsYieldSpread     float64
	SpreadPercentile        float64
	ValuationDate           string
	ValuationObservationCnt int
	VIX                     float64
	VIXDate                 string
	USDToCNY                float64
	FXDate                  string
	FuturesChange           float64
	FuturesDate             string
	TacticalSymbol          string
	MarketPrice             float64
	MarketPriceDate         string
	OfficialNAV             float64
	OfficialNAVDate         string
	EstimatedNAV            float64
	Premium                 float64
	IndexSource             string
	IndexSourceURL          string
	DirectSPTR              bool
}

type sp500OpportunityErrors struct {
	Drawdown         error
	Valuation        error
	EarningsRevision error
	VIX              error
	CNYDrawdown      error
	Premium          error
}

func fetchSP500OpportunitySnapshot(client *http.Client, now time.Time) (sp500OpportunitySnapshot, sp500OpportunityErrors) {
	snapshot := sp500OpportunitySnapshot{TacticalSymbol: sp500TacticalETFCode}
	issues := sp500OpportunityErrors{
		EarningsRevision: errors.New("free point-in-time three-month forward earnings revision series is not available"),
	}
	if client == nil {
		err := errors.New("missing HTTP client")
		issues.Drawdown = err
		issues.Valuation = err
		issues.VIX = err
		issues.CNYDrawdown = err
		issues.Premium = err
		return snapshot, issues
	}

	start := now.AddDate(-nasdaqTacticalHistoryYears, 0, -14)
	var (
		sptrCloses    []dailyClose
		sptrSource    string
		sptrSourceURL string
		sptrDirect    bool
		sptrErr       error
		fxCloses      []dailyClose
		fxErr         error
		valuation     nasdaqValuationSnapshot
		valuationErr  error
		vix           dailyClose
		vixErr        error
		market        quote
		marketErr     error
		nav           quote
		navErr        error
		futures       nasdaqFuturesSnapshot
		futuresErr    error
		fxIntraday    nasdaqFuturesSnapshot
		fxIntradayErr error
	)

	var wait sync.WaitGroup
	wait.Add(8)
	go func() {
		defer wait.Done()
		sptrCloses, sptrSource, sptrSourceURL, sptrDirect, sptrErr = fetchSP500TotalReturnClosesWithSource(client)
	}()
	go func() {
		defer wait.Done()
		fxCloses, fxErr = fetchFrankfurterUSDToCNYHistory(client, start, now)
	}()
	go func() {
		defer wait.Done()
		valuation, valuationErr = fetchSP500TacticalValuation(client, now)
	}()
	go func() {
		defer wait.Done()
		vix, vixErr = fetchLatestCboeVolatilityIndex(client, sp500VIXHistoryURL, "VIX")
	}()
	go func() {
		defer wait.Done()
		market, marketErr = fetchSP500TacticalMarketQuote(client)
	}()
	go func() {
		defer wait.Done()
		nav, navErr = fetchOTCFundHistoryQuote(client, Fund{Symbol: sp500TacticalETFCode})
	}()
	go func() {
		defer wait.Done()
		futures, futuresErr = fetchSP500FuturesChange(client)
	}()
	go func() {
		defer wait.Done()
		fxIntraday, fxIntradayErr = fetchEastmoneyUSDCNHChange(client)
	}()
	wait.Wait()
	snapshot.IndexSource = sptrSource
	snapshot.IndexSourceURL = sptrSourceURL
	snapshot.DirectSPTR = sptrDirect

	if sptrErr != nil {
		issues.Drawdown = fmt.Errorf("SPTR: %w", sptrErr)
	} else if drawdown, date, err := drawdownFromRecentHigh(sptrCloses, len(sptrCloses)); err != nil {
		issues.Drawdown = err
	} else {
		snapshot.Drawdown = drawdown
		snapshot.Date = date
	}

	if valuationErr != nil {
		issues.Valuation = valuationErr
	} else {
		snapshot.ForwardPE = valuation.ForwardPE
		snapshot.ForwardPEPercentile = valuation.ForwardPEPercentile
		snapshot.US10YBondYield = valuation.US10YBondYield
		snapshot.EarningsYieldSpread = valuation.Spread
		snapshot.SpreadPercentile = valuation.SpreadPercentile
		snapshot.ValuationDate = valuation.Date
		snapshot.ValuationObservationCnt = valuation.ObservationCount
	}

	if vixErr != nil {
		issues.VIX = vixErr
	} else {
		snapshot.VIX = vix.Price
		snapshot.VIXDate = vix.Date
	}

	if sptrErr != nil || fxErr != nil {
		issues.CNYDrawdown = combineNasdaqErrors("S&P 500 CNY total-return drawdown", sptrErr, fxErr)
	} else {
		cnyCloses, err := calculateCNYTotalReturnCloses(sptrCloses, fxCloses)
		if err != nil {
			issues.CNYDrawdown = err
		} else if drawdown, _, err := drawdownFromRecentHigh(cnyCloses, len(cnyCloses)); err != nil {
			issues.CNYDrawdown = err
		} else {
			snapshot.CNYDrawdown = drawdown
			latestFX := fxCloses[len(fxCloses)-1]
			snapshot.USDToCNY = latestFX.Price
			snapshot.FXDate = latestFX.Date
		}
	}

	if marketErr != nil || navErr != nil || futuresErr != nil || sptrErr != nil || fxErr != nil {
		issues.Premium = combineNasdaqErrors("513650 estimated premium", marketErr, navErr, futuresErr, sptrErr, fxErr)
	} else {
		estimatedNAV, err := estimateNasdaqQDIIRealtimeNAV(nav.Price, nav.PriceDate, sptrCloses, fxCloses, futures, fxIntraday, fxIntradayErr)
		if err != nil {
			issues.Premium = err
		} else if market.Price <= 0 || estimatedNAV <= 0 {
			issues.Premium = errors.New("invalid 513650 market price or estimated NAV")
		} else {
			snapshot.MarketPrice = market.Price
			snapshot.MarketPriceDate = market.PriceDate
			snapshot.OfficialNAV = nav.Price
			snapshot.OfficialNAVDate = nav.PriceDate
			snapshot.EstimatedNAV = estimatedNAV
			snapshot.Premium = market.Price/estimatedNAV - 1
			snapshot.FuturesChange = futures.Change
			snapshot.FuturesDate = futures.Date
		}
	}

	return snapshot, issues
}

func fetchSP500TotalReturnCloses(client *http.Client) ([]dailyClose, error) {
	closes, _, _, _, err := fetchSP500TotalReturnClosesWithSource(client)
	return closes, err
}

func fetchSP500TotalReturnClosesWithSource(client *http.Client) ([]dailyClose, string, string, bool, error) {
	errTexts := []string{}
	if closes, err := fetchYahooDailyCloses(client, sp500TotalReturnSymbol, "max"); err == nil {
		return closes, "Yahoo Finance SP500TR", sp500TotalReturnDataURL, true, nil
	} else {
		errTexts = append(errTexts, "Yahoo query2: "+err.Error())
	}
	query1URL := "https://query1.finance.yahoo.com/v8/finance/chart/%5ESP500TR?range=max&interval=1d"
	if closes, err := fetchYahooDailyClosesFromEndpoint(client, query1URL); err == nil {
		return closes, "Yahoo Finance SP500TR", query1URL, true, nil
	} else {
		errTexts = append(errTexts, "Yahoo query1: "+err.Error())
	}
	closes, err := fetchSPYTotalReturnCloses(client, nasdaqTacticalHistoryYears*252+120)
	if err != nil {
		return nil, "", "", false, fmt.Errorf("SPTR unavailable; SPY total-return backup failed: %s; %w", strings.Join(errTexts, "; "), err)
	}
	closes = normalizeDailyCloses(closes)
	if len(closes) < 2000 {
		return nil, "", "", false, fmt.Errorf("insufficient SPY total-return backup history: %d", len(closes))
	}
	return closes, "State Street SPY分红再投资总回报（SPTR备援）", sp500SPYTotalReturnBackupURL, false, nil
}

func fetchYahooDailyClosesFromEndpoint(client *http.Client, endpoint string) ([]dailyClose, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://finance.yahoo.com/quote/%5ESP500TR/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/138 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("SP500TR daily close request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var payload yahooChartResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 8<<20)).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Chart.Error != nil || len(payload.Chart.Result) == 0 || len(payload.Chart.Result[0].Indicators.Quote) == 0 {
		return nil, errors.New("empty SP500TR daily close response")
	}
	result := payload.Chart.Result[0]
	location := loadLocation(result.Meta.ExchangeTimezone)
	closes := make([]dailyClose, 0, len(result.Timestamp))
	for index, price := range result.Indicators.Quote[0].Close {
		if price > 0 {
			closes = append(closes, dailyClose{Date: closeDate(result.Timestamp, index, location), Price: price})
		}
	}
	closes = normalizeDailyCloses(closes)
	if len(closes) < 2000 {
		return nil, fmt.Errorf("insufficient SP500TR history: %d", len(closes))
	}
	return closes, nil
}

func fetchSP500TacticalValuation(client *http.Client, now time.Time) (nasdaqValuationSnapshot, error) {
	var payload struct {
		Updated string                          `json:"updated"`
		Current historyOfMarketCurrentValuation `json:"current"`
		Forward []historyOfMarketPoint          `json:"forward"`
	}
	if err := fetchHistoryOfMarketJSON(client, sp500ForwardPEURL, &payload); err != nil {
		return nasdaqValuationSnapshot{}, err
	}
	pePoints := historyOfMarketPointsWithCurrentForward(payload.Forward, payload.Updated, payload.Current)
	treasury, err := fetchUSTreasury10YHistory(client, now.AddDate(-nasdaqTacticalHistoryYears, 0, -14), now)
	if err != nil {
		return nasdaqValuationSnapshot{}, err
	}
	return calculateNasdaqTacticalValuation(pePoints, treasury, now)
}

func fetchLatestCboeVolatilityIndex(client *http.Client, endpoint string, name string) (dailyClose, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return dailyClose{}, err
	}
	req.Header.Set("Accept", "text/csv,*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return dailyClose{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return dailyClose{}, fmt.Errorf("Cboe %s request failed: %s %s", name, resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return dailyClose{}, err
	}
	closes, err := parseVXNHistoryCSV(body)
	if err != nil {
		return dailyClose{}, fmt.Errorf("%s: %w", name, err)
	}
	return closes[len(closes)-1], nil
}

func fetchSP500FuturesChange(client *http.Client) (nasdaqFuturesSnapshot, error) {
	if snapshot, err := fetchSinaSP500FuturesChange(client); err == nil {
		return snapshot, nil
	}
	endpoints := []string{
		sp500FuturesChartURL,
		"https://query2.finance.yahoo.com/v8/finance/chart/ES%3DF?interval=5m&range=1d",
	}
	errTexts := []string{}
	for _, endpoint := range endpoints {
		snapshot, err := fetchNasdaqFuturesChangeFromEndpoint(client, endpoint)
		if err == nil {
			return snapshot, nil
		}
		errTexts = append(errTexts, err.Error())
	}
	return nasdaqFuturesSnapshot{}, errors.New(strings.Join(errTexts, "; "))
}

func fetchSinaSP500FuturesChange(client *http.Client) (nasdaqFuturesSnapshot, error) {
	req, err := http.NewRequest(http.MethodGet, sp500SinaFuturesURL, nil)
	if err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	req.Header.Set("Accept", "text/plain,*/*")
	req.Header.Set("Referer", "https://finance.sina.com.cn/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/138.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nasdaqFuturesSnapshot{}, fmt.Errorf("Sina S&P futures request failed: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	return parseSinaNasdaqFuturesQuote(body)
}

func fetchSP500TacticalMarketQuote(client *http.Client) (quote, error) {
	errTexts := []string{}
	if quotes, err := fetchTencentQuotes(client, []string{sp500TacticalETFCode + ".SH"}); err == nil {
		if result, ok := quotes[normalizeSymbol(sp500TacticalETFCode+".SH")]; ok {
			return result, nil
		}
	} else {
		errTexts = append(errTexts, "Tencent: "+err.Error())
	}
	for attempt := 0; attempt < 3; attempt++ {
		result, err := fetchEastmoneyQuote(client, sp500TacticalETFCode+".SH")
		if err == nil {
			return normalizeNasdaqTacticalEastmoneyQuote(result), nil
		}
		errTexts = append(errTexts, err.Error())
		if attempt < 2 {
			time.Sleep(time.Duration(attempt+1) * 150 * time.Millisecond)
		}
	}
	return quote{}, fmt.Errorf("%s quote failed after retries: %s", sp500TacticalETFCode, strings.Join(errTexts, "; "))
}
