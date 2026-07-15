package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	nasdaqTacticalHistoryYears       = 10
	nasdaqTacticalMinimumWeeks       = 480
	nasdaqTacticalETFCode            = "159659"
	nasdaqVXNHistoryURL              = "https://cdn.cboe.com/api/global/us_indices/daily_prices/VXN_History.csv"
	nasdaqFXHistoryBaseURL           = "https://api.frankfurter.app"
	nasdaqUS10YTreasuryURL           = "https://home.treasury.gov/resource-center/data-chart-center/interest-rates/pages/xml"
	nasdaqSinaFuturesURL             = "https://hq.sinajs.cn/list=hf_NQ"
	nasdaqFuturesChartURL            = "https://query1.finance.yahoo.com/v8/finance/chart/NQ%3DF?interval=5m&range=1d"
	nasdaqTacticalETFQuoteURL        = "https://quote.eastmoney.com/sz159659.html"
	nasdaqTacticalETFNetValuePageURL = "https://fund.eastmoney.com/159659.html"
)

type nasdaqOpportunitySnapshot struct {
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
	VXN                     float64
	VXNDate                 string
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
}

type nasdaqOpportunityErrors struct {
	Drawdown    error
	Valuation   error
	VXN         error
	CNYDrawdown error
	Premium     error
}

type nasdaqValuationSnapshot struct {
	Date                string
	ForwardPE           float64
	ForwardPEPercentile float64
	US10YBondYield      float64
	Spread              float64
	SpreadPercentile    float64
	ObservationCount    int
}

type nasdaqFuturesSnapshot struct {
	Change float64
	Date   string
}

func fetchNasdaqOpportunitySnapshot(client *http.Client, now time.Time) (nasdaqOpportunitySnapshot, nasdaqOpportunityErrors) {
	snapshot := nasdaqOpportunitySnapshot{TacticalSymbol: nasdaqTacticalETFCode}
	issues := nasdaqOpportunityErrors{}
	if client == nil {
		err := errors.New("missing HTTP client")
		issues.Drawdown = err
		issues.Valuation = err
		issues.VXN = err
		issues.CNYDrawdown = err
		issues.Premium = err
		return snapshot, issues
	}

	start := now.AddDate(-nasdaqTacticalHistoryYears, 0, -14)
	var (
		xndxCloses    []dailyClose
		xndxErr       error
		fxCloses      []dailyClose
		fxErr         error
		valuation     nasdaqValuationSnapshot
		valuationErr  error
		vxn           dailyClose
		vxnErr        error
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
		xndxCloses, xndxErr = fetchNasdaqIndexHistoryChart(client, "XNDX", start, now)
	}()
	go func() {
		defer wait.Done()
		fxCloses, fxErr = fetchFrankfurterUSDToCNYHistory(client, start, now)
	}()
	go func() {
		defer wait.Done()
		valuation, valuationErr = fetchNasdaqTacticalValuation(client, now)
	}()
	go func() {
		defer wait.Done()
		vxn, vxnErr = fetchLatestVXN(client)
	}()
	go func() {
		defer wait.Done()
		market, marketErr = fetchNasdaqTacticalMarketQuote(client)
	}()
	go func() {
		defer wait.Done()
		nav, navErr = fetchOTCFundHistoryQuote(client, Fund{Symbol: nasdaqTacticalETFCode})
	}()
	go func() {
		defer wait.Done()
		futures, futuresErr = fetchNasdaqFuturesChange(client)
	}()
	go func() {
		defer wait.Done()
		fxIntraday, fxIntradayErr = fetchEastmoneyUSDCNHChange(client)
	}()
	wait.Wait()

	if xndxErr != nil {
		issues.Drawdown = fmt.Errorf("XNDX: %w", xndxErr)
	} else {
		drawdown, date, err := drawdownFromRecentHigh(xndxCloses, len(xndxCloses))
		if err != nil {
			issues.Drawdown = err
		} else {
			snapshot.Drawdown = drawdown
			snapshot.Date = date
		}
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

	if vxnErr != nil {
		issues.VXN = vxnErr
	} else {
		snapshot.VXN = vxn.Price
		snapshot.VXNDate = vxn.Date
	}

	if xndxErr != nil || fxErr != nil {
		issues.CNYDrawdown = combineNasdaqErrors("人民币全收益回撤", xndxErr, fxErr)
	} else {
		cnyCloses, err := calculateCNYTotalReturnCloses(xndxCloses, fxCloses)
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

	if marketErr != nil || navErr != nil || futuresErr != nil || xndxErr != nil || fxErr != nil {
		issues.Premium = combineNasdaqErrors("159659估算溢价", marketErr, navErr, futuresErr, xndxErr, fxErr)
	} else {
		estimatedNAV, err := estimateNasdaqQDIIRealtimeNAV(nav.Price, nav.PriceDate, xndxCloses, fxCloses, futures, fxIntraday, fxIntradayErr)
		if err != nil {
			issues.Premium = err
		} else if market.Price <= 0 || estimatedNAV <= 0 {
			issues.Premium = errors.New("invalid 159659 market price or estimated NAV")
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

func combineNasdaqErrors(label string, errs ...error) error {
	parts := make([]string, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			parts = append(parts, err.Error())
		}
	}
	if len(parts) == 0 {
		return nil
	}
	return fmt.Errorf("%s: %s", label, strings.Join(parts, "; "))
}

func fetchNasdaqTacticalValuation(client *http.Client, now time.Time) (nasdaqValuationSnapshot, error) {
	var payload struct {
		Updated string                          `json:"updated"`
		Current historyOfMarketCurrentValuation `json:"current"`
		Forward []historyOfMarketPoint          `json:"forward"`
	}
	if err := fetchHistoryOfMarketJSON(client, historyOfMarketNDXForwardPEURL, &payload); err != nil {
		return nasdaqValuationSnapshot{}, err
	}
	pePoints := historyOfMarketPointsWithCurrentForward(payload.Forward, payload.Updated, payload.Current)
	treasury, err := fetchUSTreasury10YHistory(client, now.AddDate(-nasdaqTacticalHistoryYears, 0, -14), now)
	if err != nil {
		return nasdaqValuationSnapshot{}, err
	}
	return calculateNasdaqTacticalValuation(pePoints, treasury, now)
}

func calculateNasdaqTacticalValuation(pePoints []historyOfMarketPoint, treasury []dailyClose, now time.Time) (nasdaqValuationSnapshot, error) {
	pePoints = normalizePositiveHistoryPoints(pePoints)
	treasury = normalizeDailyCloses(treasury)
	if len(pePoints) == 0 || len(treasury) == 0 {
		return nasdaqValuationSnapshot{}, errors.New("missing Nasdaq forward PE or Treasury history")
	}
	cutoff := now.AddDate(-nasdaqTacticalHistoryYears, 0, 0).Format(etfRuleRuntimeTimestampDateLayout)
	weeklyTreasury := latestClosePerISOWeek(treasury, cutoff)
	weeklyPE := make([]float64, 0, len(weeklyTreasury))
	weeklySpreads := make([]float64, 0, len(weeklyTreasury))
	peIndex := 0
	currentPE := 0.0
	for _, bond := range weeklyTreasury {
		for peIndex < len(pePoints) && pePoints[peIndex].Date <= bond.Date {
			currentPE = pePoints[peIndex].Value
			peIndex++
		}
		if currentPE <= 0 || bond.Price <= 0 {
			continue
		}
		weeklyPE = append(weeklyPE, currentPE)
		weeklySpreads = append(weeklySpreads, 1/currentPE-bond.Price)
	}
	if len(weeklySpreads) < nasdaqTacticalMinimumWeeks {
		return nasdaqValuationSnapshot{}, fmt.Errorf("insufficient 10-year weekly valuation history: %d", len(weeklySpreads))
	}
	latestPE := pePoints[len(pePoints)-1]
	latestBond := treasury[len(treasury)-1]
	if latestPE.Value <= 0 || latestBond.Price <= 0 {
		return nasdaqValuationSnapshot{}, errors.New("invalid latest Nasdaq valuation")
	}
	spread := 1/latestPE.Value - latestBond.Price
	return nasdaqValuationSnapshot{
		Date:                firstNonEmpty(latestBond.Date, latestPE.Date),
		ForwardPE:           latestPE.Value,
		ForwardPEPercentile: percentileRank(latestPE.Value, weeklyPE),
		US10YBondYield:      latestBond.Price,
		Spread:              spread,
		SpreadPercentile:    percentileRank(spread, weeklySpreads),
		ObservationCount:    len(weeklySpreads),
	}, nil
}

func normalizePositiveHistoryPoints(points []historyOfMarketPoint) []historyOfMarketPoint {
	byDate := map[string]historyOfMarketPoint{}
	for _, point := range points {
		if _, err := time.Parse(etfRuleRuntimeTimestampDateLayout, strings.TrimSpace(point.Date)); err != nil || point.Value <= 0 || math.IsNaN(point.Value) || math.IsInf(point.Value, 0) {
			continue
		}
		point.Date = strings.TrimSpace(point.Date)
		byDate[point.Date] = point
	}
	result := make([]historyOfMarketPoint, 0, len(byDate))
	for _, point := range byDate {
		result = append(result, point)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Date < result[j].Date })
	return result
}

func normalizeDailyCloses(closes []dailyClose) []dailyClose {
	byDate := map[string]dailyClose{}
	for _, close := range closes {
		if _, err := time.Parse(etfRuleRuntimeTimestampDateLayout, strings.TrimSpace(close.Date)); err != nil || close.Price <= 0 || math.IsNaN(close.Price) || math.IsInf(close.Price, 0) {
			continue
		}
		close.Date = strings.TrimSpace(close.Date)
		byDate[close.Date] = close
	}
	result := make([]dailyClose, 0, len(byDate))
	for _, close := range byDate {
		result = append(result, close)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Date < result[j].Date })
	return result
}

func latestClosePerISOWeek(closes []dailyClose, cutoff string) []dailyClose {
	byWeek := map[string]dailyClose{}
	for _, close := range closes {
		if close.Date < cutoff {
			continue
		}
		date, err := time.Parse(etfRuleRuntimeTimestampDateLayout, close.Date)
		if err != nil {
			continue
		}
		key := isoWeekKey(date)
		if previous, ok := byWeek[key]; !ok || close.Date > previous.Date {
			byWeek[key] = close
		}
	}
	result := make([]dailyClose, 0, len(byWeek))
	for _, close := range byWeek {
		result = append(result, close)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Date < result[j].Date })
	return result
}

func fetchUSTreasury10YHistory(client *http.Client, start time.Time, end time.Time) ([]dailyClose, error) {
	if client == nil {
		return nil, errors.New("missing HTTP client")
	}
	type result struct {
		closes []dailyClose
		err    error
	}
	years := end.Year() - start.Year() + 1
	results := make(chan result, years)
	sem := make(chan struct{}, 4)
	var wait sync.WaitGroup
	for year := start.Year(); year <= end.Year(); year++ {
		year := year
		wait.Add(1)
		go func() {
			defer wait.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			closes, err := fetchUSTreasury10YYear(client, year)
			results <- result{closes: closes, err: err}
		}()
	}
	go func() {
		wait.Wait()
		close(results)
	}()

	all := []dailyClose{}
	errs := []string{}
	for item := range results {
		if item.err != nil {
			errs = append(errs, item.err.Error())
			continue
		}
		all = append(all, item.closes...)
	}
	all = normalizeDailyCloses(all)
	filtered := all[:0]
	startKey := start.Format(etfRuleRuntimeTimestampDateLayout)
	endKey := end.Format(etfRuleRuntimeTimestampDateLayout)
	for _, close := range all {
		if close.Date >= startKey && close.Date <= endKey {
			filtered = append(filtered, close)
		}
	}
	if len(filtered) < 2000 {
		return nil, fmt.Errorf("insufficient US Treasury 10-year history (%d rows; %s)", len(filtered), strings.Join(errs, "; "))
	}
	return filtered, nil
}

func fetchUSTreasury10YYear(client *http.Client, year int) ([]dailyClose, error) {
	values := url.Values{}
	values.Set("data", "daily_treasury_yield_curve")
	values.Set("field_tdr_date_value", strconv.Itoa(year))
	req, err := http.NewRequest(http.MethodGet, nasdaqUS10YTreasuryURL+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/atom+xml,application/xml,text/xml")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("US Treasury %d request failed: %s %s", year, resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	closes, err := parseUSTreasury10YXML(body)
	if err != nil {
		return nil, fmt.Errorf("US Treasury %d: %w", year, err)
	}
	return closes, nil
}

func parseUSTreasury10YXML(body []byte) ([]dailyClose, error) {
	var feed struct {
		Entries []struct {
			Content struct {
				Properties struct {
					Date  string `xml:"NEW_DATE"`
					Yield string `xml:"BC_10YEAR"`
				} `xml:"properties"`
			} `xml:"content"`
		} `xml:"entry"`
	}
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}
	closes := make([]dailyClose, 0, len(feed.Entries))
	for _, entry := range feed.Entries {
		date := normalizeNasdaqTacticalDate(entry.Content.Properties.Date)
		yield, err := strconv.ParseFloat(strings.TrimSpace(entry.Content.Properties.Yield), 64)
		if date == "" || err != nil || yield <= 0 {
			continue
		}
		closes = append(closes, dailyClose{Date: date, Price: yield / 100})
	}
	if len(closes) == 0 {
		return nil, errors.New("missing US Treasury 10-year observations")
	}
	return normalizeDailyCloses(closes), nil
}

func normalizeNasdaqTacticalDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05", etfRuleRuntimeTimestampDateLayout, "01/02/2006"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format(etfRuleRuntimeTimestampDateLayout)
		}
	}
	return ""
}

func fetchLatestVXN(client *http.Client) (dailyClose, error) {
	return fetchLatestCboeVolatilityIndex(client, nasdaqVXNHistoryURL, "VXN")
}

func parseVXNHistoryCSV(body []byte) ([]dailyClose, error) {
	reader := csv.NewReader(strings.NewReader(strings.TrimSpace(string(body))))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	closes := []dailyClose{}
	for index, record := range records {
		if index == 0 || len(record) < 5 {
			continue
		}
		date := normalizeNasdaqTacticalDate(record[0])
		value, err := strconv.ParseFloat(strings.TrimSpace(record[4]), 64)
		if date == "" || err != nil || value <= 0 {
			continue
		}
		closes = append(closes, dailyClose{Date: date, Price: value})
	}
	closes = normalizeDailyCloses(closes)
	if len(closes) == 0 {
		return nil, errors.New("missing VXN history")
	}
	return closes, nil
}

func fetchFrankfurterUSDToCNYHistory(client *http.Client, start time.Time, end time.Time) ([]dailyClose, error) {
	endpoint := fmt.Sprintf("%s/%s..%s?from=USD&to=CNY", nasdaqFXHistoryBaseURL, start.Format(etfRuleRuntimeTimestampDateLayout), end.Format(etfRuleRuntimeTimestampDateLayout))
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("Frankfurter USD/CNY request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var payload struct {
		Rates map[string]map[string]float64 `json:"rates"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&payload); err != nil {
		return nil, err
	}
	closes := make([]dailyClose, 0, len(payload.Rates))
	for date, rates := range payload.Rates {
		if value := rates["CNY"]; value > 0 {
			closes = append(closes, dailyClose{Date: date, Price: value})
		}
	}
	closes = normalizeDailyCloses(closes)
	if len(closes) < 2000 {
		return nil, fmt.Errorf("insufficient USD/CNY history: %d", len(closes))
	}
	return closes, nil
}

func calculateCNYTotalReturnCloses(xndx []dailyClose, fx []dailyClose) ([]dailyClose, error) {
	xndx = normalizeDailyCloses(xndx)
	fx = normalizeDailyCloses(fx)
	if len(xndx) == 0 || len(fx) == 0 {
		return nil, errors.New("missing XNDX or USD/CNY history")
	}
	result := make([]dailyClose, 0, len(xndx))
	fxIndex := 0
	latestFX := 0.0
	for _, indexClose := range xndx {
		for fxIndex < len(fx) && fx[fxIndex].Date <= indexClose.Date {
			latestFX = fx[fxIndex].Price
			fxIndex++
		}
		if latestFX <= 0 {
			continue
		}
		result = append(result, dailyClose{Date: indexClose.Date, Price: indexClose.Price * latestFX})
	}
	if len(result) < 2000 {
		return nil, fmt.Errorf("insufficient CNY total-return history: %d", len(result))
	}
	return result, nil
}

func fetchNasdaqFuturesChange(client *http.Client) (nasdaqFuturesSnapshot, error) {
	if snapshot, err := fetchSinaNasdaqFuturesChange(client); err == nil {
		return snapshot, nil
	}
	endpoints := []string{
		nasdaqFuturesChartURL,
		"https://query2.finance.yahoo.com/v8/finance/chart/NQ%3DF?interval=5m&range=1d",
	}
	errs := []string{}
	for _, endpoint := range endpoints {
		snapshot, err := fetchNasdaqFuturesChangeFromEndpoint(client, endpoint)
		if err == nil {
			return snapshot, nil
		}
		errs = append(errs, err.Error())
	}
	return nasdaqFuturesSnapshot{}, errors.New(strings.Join(errs, "; "))
}

func fetchSinaNasdaqFuturesChange(client *http.Client) (nasdaqFuturesSnapshot, error) {
	req, err := http.NewRequest(http.MethodGet, nasdaqSinaFuturesURL, nil)
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
		return nasdaqFuturesSnapshot{}, fmt.Errorf("Sina Nasdaq futures request failed: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	return parseSinaNasdaqFuturesQuote(body)
}

func parseSinaNasdaqFuturesQuote(body []byte) (nasdaqFuturesSnapshot, error) {
	text := strings.TrimSpace(string(body))
	start := strings.Index(text, "\"")
	end := strings.LastIndex(text, "\"")
	if start < 0 || end <= start {
		return nasdaqFuturesSnapshot{}, errors.New("invalid Sina Nasdaq futures payload")
	}
	fields := strings.Split(text[start+1:end], ",")
	if len(fields) < 13 {
		return nasdaqFuturesSnapshot{}, errors.New("incomplete Sina Nasdaq futures payload")
	}
	price, priceErr := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
	previous, previousErr := strconv.ParseFloat(strings.TrimSpace(fields[7]), 64)
	date := normalizeNasdaqTacticalDate(fields[12])
	if priceErr != nil || previousErr != nil || price <= 0 || previous <= 0 || date == "" {
		return nasdaqFuturesSnapshot{}, errors.New("invalid Sina Nasdaq futures quote")
	}
	return nasdaqFuturesSnapshot{Change: price/previous - 1, Date: date}, nil
}

func fetchNasdaqFuturesChangeFromEndpoint(client *http.Client, endpoint string) (nasdaqFuturesSnapshot, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://finance.yahoo.com/quote/NQ=F/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nasdaqFuturesSnapshot{}, fmt.Errorf("Nasdaq futures request failed: %s", resp.Status)
	}
	var payload struct {
		Chart struct {
			Result []struct {
				Meta struct {
					Price         float64 `json:"regularMarketPrice"`
					PreviousClose float64 `json:"previousClose"`
					Timestamp     int64   `json:"regularMarketTime"`
				} `json:"meta"`
			} `json:"result"`
		} `json:"chart"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	if len(payload.Chart.Result) == 0 {
		return nasdaqFuturesSnapshot{}, errors.New("missing Nasdaq futures result")
	}
	meta := payload.Chart.Result[0].Meta
	if meta.Price <= 0 || meta.PreviousClose <= 0 || meta.Timestamp <= 0 {
		return nasdaqFuturesSnapshot{}, errors.New("invalid Nasdaq futures quote")
	}
	return nasdaqFuturesSnapshot{
		Change: meta.Price/meta.PreviousClose - 1,
		Date:   time.Unix(meta.Timestamp, 0).In(loadLocation("America/New_York")).Format(etfRuleRuntimeTimestampDateLayout),
	}, nil
}

func fetchNasdaqTacticalMarketQuote(client *http.Client) (quote, error) {
	errTexts := []string{}
	if quotes, err := fetchTencentQuotes(client, []string{nasdaqTacticalETFCode + ".SZ"}); err == nil {
		if result, ok := quotes[normalizeSymbol(nasdaqTacticalETFCode+".SZ")]; ok {
			return result, nil
		}
	} else {
		errTexts = append(errTexts, "Tencent: "+err.Error())
	}
	for attempt := 0; attempt < 3; attempt++ {
		result, err := fetchEastmoneyQuote(client, nasdaqTacticalETFCode+".SZ")
		if err == nil {
			return normalizeNasdaqTacticalEastmoneyQuote(result), nil
		}
		errTexts = append(errTexts, err.Error())
		if attempt < 2 {
			time.Sleep(time.Duration(attempt+1) * 150 * time.Millisecond)
		}
	}
	return quote{}, fmt.Errorf("%s quote failed after retries: %s", nasdaqTacticalETFCode, strings.Join(errTexts, "; "))
}

func normalizeNasdaqTacticalEastmoneyQuote(result quote) quote {
	// Eastmoney f43/f60 use three decimals for exchange-traded funds, while the shared stock parser uses two.
	result.Price /= 10
	result.PreviousClose /= 10
	return result
}

func fetchEastmoneyUSDCNHChange(client *http.Client) (nasdaqFuturesSnapshot, error) {
	endpoint := "https://push2.eastmoney.com/api/qt/stock/get?secid=133.USDCNH&fields=f43,f60,f86"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Referer", "https://quote.eastmoney.com/center/gridlist.html#forex")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	defer resp.Body.Close()
	var payload struct {
		RC   int            `json:"rc"`
		Data map[string]any `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil {
		return nasdaqFuturesSnapshot{}, err
	}
	if payload.RC != 0 || len(payload.Data) == 0 {
		return nasdaqFuturesSnapshot{}, errors.New("missing USD/CNH quote")
	}
	price, priceErr := numberField(payload.Data, "f43")
	previous, previousErr := numberField(payload.Data, "f60")
	if priceErr != nil || previousErr != nil || price <= 0 || previous <= 0 {
		return nasdaqFuturesSnapshot{}, errors.New("invalid USD/CNH quote")
	}
	return nasdaqFuturesSnapshot{Change: price/previous - 1, Date: eastmoneyQuoteDate(payload.Data)}, nil
}

func estimateNasdaqQDIIRealtimeNAV(officialNAV float64, navDate string, xndx []dailyClose, fx []dailyClose, futures nasdaqFuturesSnapshot, fxIntraday nasdaqFuturesSnapshot, fxIntradayErr error) (float64, error) {
	if officialNAV <= 0 || navDate == "" {
		return 0, errors.New("invalid official QDII NAV")
	}
	baseIndex, ok := dailyCloseOnOrBefore(xndx, navDate)
	if !ok {
		return 0, fmt.Errorf("missing XNDX close on or before %s", navDate)
	}
	baseFX, ok := dailyCloseOnOrBefore(fx, navDate)
	if !ok {
		return 0, fmt.Errorf("missing USD/CNY close on or before %s", navDate)
	}
	latestIndex := normalizeDailyCloses(xndx)
	latestFX := normalizeDailyCloses(fx)
	if len(latestIndex) == 0 || len(latestFX) == 0 {
		return 0, errors.New("missing current XNDX or USD/CNY close")
	}
	estimated := officialNAV * latestIndex[len(latestIndex)-1].Price / baseIndex.Price * latestFX[len(latestFX)-1].Price / baseFX.Price
	if futures.Date > latestIndex[len(latestIndex)-1].Date {
		estimated *= 1 + futures.Change
	}
	if fxIntradayErr == nil && fxIntraday.Date > latestFX[len(latestFX)-1].Date {
		estimated *= 1 + fxIntraday.Change
	}
	if estimated <= 0 || math.IsNaN(estimated) || math.IsInf(estimated, 0) {
		return 0, errors.New("invalid estimated QDII NAV")
	}
	return estimated, nil
}

func dailyCloseOnOrBefore(closes []dailyClose, target string) (dailyClose, bool) {
	closes = normalizeDailyCloses(closes)
	index := sort.Search(len(closes), func(index int) bool { return closes[index].Date > target })
	if index == 0 {
		return dailyClose{}, false
	}
	return closes[index-1], true
}
