package main

import (
	"encoding/json"
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
	a500TotalReturnIndexCode   = "000510CNY010"
	a500PriceIndexCode         = "000510"
	a500TacticalETFCode        = "159352"
	a500HistoryStartDate       = "2004-01-01"
	a500ValuationStartDate     = "2024-09-02"
	a500TotalReturnIndexURL    = "https://www.csindex.com.cn/csindex-home/perf/index-perf?indexCode=000510CNY010"
	a500PEHistoryURL           = "https://www.csindex.com.cn/csindex-home/perf/indexCsiDsPe?indexCode=000510"
	a500TacticalETFQuoteURL    = "https://quote.eastmoney.com/sz159352.html"
	a500TacticalETFNetValueURL = "https://fund.eastmoney.com/159352.html"
)

type a500PerformancePoint struct {
	Date         string
	Close        float64
	TradingValue float64
}

type a500PEPoint struct {
	Date string
	PE   float64
}

type a500TacticalMarketSnapshot struct {
	Price         float64
	PreviousClose float64
	Open          float64
	Bid           float64
	Ask           float64
	Date          string
}

type a500OpportunitySnapshot struct {
	Date                   string
	IndexClose             float64
	PeakClose              float64
	PeakDate               string
	Drawdown               float64
	PE                     float64
	PEPercentile           float64
	PEObservationCount     int
	China10YBondYield      float64
	BondDate               string
	EarningsYieldSpread    float64
	SpreadPercentile       float64
	SpreadObservationCount int
	RV20                   float64
	RV20Percentile         float64
	RV20ObservationCount   int
	FiveDayReturn          float64
	VolumeRatio            float64
	MarketPrice            float64
	MarketPriceDate        string
	Bid                    float64
	Ask                    float64
	BidAskSpread           float64
	OpeningGap             float64
	OfficialNAV            float64
	OfficialNAVDate        string
	EstimatedNAV           float64
	Premium                float64
}

type a500OpportunityErrors struct {
	Drawdown  error
	Valuation error
	Panic     error
	Trading   error
}

func fetchA500OpportunitySnapshot(client *http.Client, now time.Time) (a500OpportunitySnapshot, a500OpportunityErrors) {
	var (
		performance    []a500PerformancePoint
		peHistory      []a500PEPoint
		bondHistory    []datedRate
		market         a500TacticalMarketSnapshot
		nav            quote
		performanceErr error
		peErr          error
		bondErr        error
		marketErr      error
		navErr         error
	)
	start, _ := time.Parse(etfRuleRuntimeTimestampDateLayout, a500HistoryStartDate)
	valuationStart, _ := time.Parse(etfRuleRuntimeTimestampDateLayout, a500ValuationStartDate)
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		performance, performanceErr = fetchA500Performance(client, start, now)
	}()
	go func() {
		defer wg.Done()
		peHistory, peErr = fetchA500PEHistory(client, valuationStart, now)
	}()
	go func() {
		defer wg.Done()
		bondHistory, bondErr = fetchEastmoneyChina10YBondYieldHistory(client, a500ValuationStartDate)
	}()
	go func() {
		defer wg.Done()
		market, marketErr = fetchA500TacticalMarketQuote(client)
	}()
	go func() {
		defer wg.Done()
		nav, navErr = fetchOTCFundHistoryQuote(client, Fund{Symbol: a500TacticalETFCode})
	}()
	wg.Wait()

	snapshot := a500OpportunitySnapshot{}
	issues := a500OpportunityErrors{}
	if performanceErr != nil {
		issues.Drawdown = performanceErr
		issues.Panic = performanceErr
	} else {
		if err := applyA500MarketIndicators(&snapshot, performance); err != nil {
			issues.Drawdown = err
		}
		if err := applyA500PanicIndicators(&snapshot, performance); err != nil {
			issues.Panic = err
		}
	}
	if peErr != nil || bondErr != nil {
		issues.Valuation = errors.Join(peErr, bondErr)
	} else if err := applyA500Valuation(client, &snapshot, peHistory, bondHistory); err != nil {
		issues.Valuation = err
	}
	if marketErr != nil || navErr != nil || performanceErr != nil {
		issues.Trading = errors.Join(marketErr, navErr, performanceErr)
	} else if err := applyA500TradingIndicators(&snapshot, performance, market, nav); err != nil {
		issues.Trading = err
	}
	return snapshot, issues
}

func fetchA500Performance(client *http.Client, start time.Time, end time.Time) ([]a500PerformancePoint, error) {
	values := url.Values{}
	values.Set("indexCode", a500TotalReturnIndexCode)
	values.Set("startDate", start.Format("20060102"))
	values.Set("endDate", end.Format("20060102"))
	var payload struct {
		Code string `json:"code"`
		Data []struct {
			TradeDate    string  `json:"tradeDate"`
			Close        float64 `json:"close"`
			TradingValue float64 `json:"tradingValue"`
		} `json:"data"`
	}
	if err := fetchCSIJSON(client, "https://www.csindex.com.cn/csindex-home/perf/index-perf?"+values.Encode(), &payload); err != nil {
		return nil, err
	}
	if payload.Code != "200" {
		return nil, fmt.Errorf("CSI A500 total-return response code %s", payload.Code)
	}
	points := make([]a500PerformancePoint, 0, len(payload.Data))
	for _, row := range payload.Data {
		date, err := time.Parse("20060102", strings.TrimSpace(row.TradeDate))
		if err != nil || row.Close <= 0 {
			continue
		}
		points = append(points, a500PerformancePoint{Date: date.Format(etfRuleRuntimeTimestampDateLayout), Close: row.Close, TradingValue: row.TradingValue})
	}
	if len(points) < 30 {
		return nil, fmt.Errorf("insufficient CSI A500 total-return history: %d observations", len(points))
	}
	sort.Slice(points, func(i, j int) bool { return points[i].Date < points[j].Date })
	return points, nil
}

func fetchA500PEHistory(client *http.Client, start time.Time, end time.Time) ([]a500PEPoint, error) {
	values := url.Values{}
	values.Set("indexCode", a500PriceIndexCode)
	values.Set("startDate", start.Format("20060102"))
	values.Set("endDate", end.Format("20060102"))
	var payload struct {
		Code string `json:"code"`
		Data []struct {
			TradeDate string  `json:"tradeDate"`
			PE        float64 `json:"peg"`
		} `json:"data"`
	}
	if err := fetchCSIJSON(client, "https://www.csindex.com.cn/csindex-home/perf/indexCsiDsPe?"+values.Encode(), &payload); err != nil {
		return nil, err
	}
	if payload.Code != "200" {
		return nil, fmt.Errorf("CSI A500 PE response code %s", payload.Code)
	}
	points := make([]a500PEPoint, 0, len(payload.Data))
	for _, row := range payload.Data {
		date, err := time.Parse("20060102", strings.TrimSpace(row.TradeDate))
		if err != nil || row.PE <= 0 || row.PE > 100 {
			continue
		}
		points = append(points, a500PEPoint{Date: date.Format(etfRuleRuntimeTimestampDateLayout), PE: row.PE})
	}
	if len(points) < 60 {
		return nil, fmt.Errorf("insufficient CSI A500 PE history: %d observations", len(points))
	}
	sort.Slice(points, func(i, j int) bool { return points[i].Date < points[j].Date })
	return points, nil
}

func fetchCSIJSON(client *http.Client, endpoint string, target any) error {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.csindex.com.cn/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("CSI request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func applyA500MarketIndicators(snapshot *a500OpportunitySnapshot, points []a500PerformancePoint) error {
	if len(points) == 0 {
		return errors.New("missing CSI A500 total-return history")
	}
	peak := points[0]
	for _, point := range points[1:] {
		if point.Close > peak.Close {
			peak = point
		}
	}
	latest := points[len(points)-1]
	if peak.Close <= 0 || latest.Close <= 0 {
		return errors.New("invalid CSI A500 total-return close")
	}
	snapshot.Date = latest.Date
	snapshot.IndexClose = latest.Close
	snapshot.PeakClose = peak.Close
	snapshot.PeakDate = peak.Date
	snapshot.Drawdown = math.Max(0, 1-latest.Close/peak.Close)
	return nil
}

func applyA500Valuation(client *http.Client, snapshot *a500OpportunitySnapshot, peHistory []a500PEPoint, bondHistory []datedRate) error {
	if len(peHistory) == 0 {
		return errors.New("missing CSI A500 PE history")
	}
	latest := peHistory[len(peHistory)-1]
	latestDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, latest.Date)
	if err != nil {
		return err
	}
	officialBond, err := fetchChinaBondOfficial10YYield(client, latestDate)
	if err != nil {
		return err
	}
	peValues := make([]float64, 0, len(peHistory))
	spreads := make([]float64, 0, len(peHistory))
	for _, point := range peHistory {
		peValues = append(peValues, point.PE)
		bond, ok := datedRateOnOrBefore(bondHistory, point.Date, 10)
		if !ok || point.PE <= 0 {
			continue
		}
		spreads = append(spreads, 1/point.PE-bond.Value)
	}
	if len(spreads) < 60 {
		return fmt.Errorf("insufficient CSI A500 earnings-spread history: %d observations", len(spreads))
	}
	historyBond, ok := datedRateOnOrBefore(bondHistory, officialBond.Date, 0)
	if !ok {
		return fmt.Errorf("missing mirrored China 10Y yield for %s", officialBond.Date)
	}
	if math.Abs(historyBond.Value-officialBond.Value)/officialBond.Value > 0.01 {
		return fmt.Errorf("China 10Y yield sources differ: ChinaBond %.4f%%, Eastmoney %.4f%%", officialBond.Value*100, historyBond.Value*100)
	}
	currentSpread := 1/latest.PE - officialBond.Value
	snapshot.PE = latest.PE
	snapshot.PEPercentile = percentileRank(latest.PE, peValues)
	snapshot.PEObservationCount = len(peValues)
	snapshot.China10YBondYield = officialBond.Value
	snapshot.BondDate = officialBond.Date
	snapshot.EarningsYieldSpread = currentSpread
	snapshot.SpreadPercentile = percentileRank(currentSpread, spreads)
	snapshot.SpreadObservationCount = len(spreads)
	return nil
}

func applyA500PanicIndicators(snapshot *a500OpportunitySnapshot, points []a500PerformancePoint) error {
	if len(points) < 22 {
		return fmt.Errorf("insufficient A500 history for RV20: %d observations", len(points))
	}
	returns := make([]float64, 0, len(points)-1)
	for index := 1; index < len(points); index++ {
		if points[index-1].Close <= 0 || points[index].Close <= 0 {
			continue
		}
		returns = append(returns, math.Log(points[index].Close/points[index-1].Close))
	}
	if len(returns) < 20 {
		return errors.New("insufficient A500 return observations")
	}
	latestRV := annualizedRealizedVolatility(returns[len(returns)-20:])
	latestDate, _ := time.Parse(etfRuleRuntimeTimestampDateLayout, points[len(points)-1].Date)
	cutoff := latestDate.AddDate(-5, 0, 0).Format(etfRuleRuntimeTimestampDateLayout)
	history := make([]float64, 0, len(returns))
	for end := 20; end <= len(returns); end++ {
		pointDate := points[end].Date
		if pointDate < cutoff {
			continue
		}
		history = append(history, annualizedRealizedVolatility(returns[end-20:end]))
	}
	if len(history) < 60 {
		return fmt.Errorf("insufficient A500 five-year RV20 history: %d observations", len(history))
	}
	latest := points[len(points)-1]
	fiveDayBase := points[len(points)-6]
	snapshot.RV20 = latestRV
	snapshot.RV20Percentile = percentileRank(latestRV, history)
	snapshot.RV20ObservationCount = len(history)
	snapshot.FiveDayReturn = latest.Close/fiveDayBase.Close - 1
	if len(points) >= 22 {
		total := 0.0
		count := 0
		for _, point := range points[len(points)-21 : len(points)-1] {
			if point.TradingValue > 0 {
				total += point.TradingValue
				count++
			}
		}
		if count > 0 && latest.TradingValue > 0 {
			snapshot.VolumeRatio = latest.TradingValue / (total / float64(count))
		}
	}
	return nil
}

func annualizedRealizedVolatility(returns []float64) float64 {
	if len(returns) < 2 {
		return 0
	}
	mean := 0.0
	for _, value := range returns {
		mean += value
	}
	mean /= float64(len(returns))
	variance := 0.0
	for _, value := range returns {
		delta := value - mean
		variance += delta * delta
	}
	variance /= float64(len(returns) - 1)
	return math.Sqrt(variance) * math.Sqrt(252)
}

func fetchA500TacticalMarketQuote(client *http.Client) (a500TacticalMarketSnapshot, error) {
	endpoint := "http://qt.gtimg.cn/q=sz" + a500TacticalETFCode
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return a500TacticalMarketSnapshot{}, err
	}
	req.Header.Set("Accept", "text/plain,*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return a500TacticalMarketSnapshot{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return a500TacticalMarketSnapshot{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return a500TacticalMarketSnapshot{}, fmt.Errorf("Tencent A500 ETF quote failed: %s", resp.Status)
	}
	return parseA500TacticalMarketQuote(body)
}

func parseA500TacticalMarketQuote(body []byte) (a500TacticalMarketSnapshot, error) {
	for _, line := range strings.Split(string(body), ";") {
		_, fields, ok := parseTencentLine(line)
		if !ok || len(fields) <= 30 {
			continue
		}
		parse := func(index int) float64 {
			if index >= len(fields) {
				return 0
			}
			value, _ := strconv.ParseFloat(strings.TrimSpace(fields[index]), 64)
			return value
		}
		market := a500TacticalMarketSnapshot{
			Price: parse(3), PreviousClose: parse(4), Open: parse(5), Bid: parse(9), Ask: parse(19), Date: tencentQuoteDate(fields[30]),
		}
		if market.Price > 0 && market.PreviousClose > 0 && market.Bid > 0 && market.Ask >= market.Bid {
			return market, nil
		}
	}
	return a500TacticalMarketSnapshot{}, errors.New("missing Tencent 159352 market depth")
}

func applyA500TradingIndicators(snapshot *a500OpportunitySnapshot, points []a500PerformancePoint, market a500TacticalMarketSnapshot, nav quote) error {
	if market.Price <= 0 || market.Bid <= 0 || market.Ask < market.Bid || nav.Price <= 0 {
		return errors.New("invalid 159352 quote or NAV")
	}
	navIndex, ok := a500PerformanceOnOrBefore(points, nav.PriceDate)
	if !ok {
		return fmt.Errorf("missing A500 total-return close for NAV date %s", nav.PriceDate)
	}
	marketIndex, ok := a500PerformanceOnOrBefore(points, market.Date)
	if !ok {
		return fmt.Errorf("missing A500 total-return close for market date %s", market.Date)
	}
	estimatedNAV := nav.Price * marketIndex.Close / navIndex.Close
	if estimatedNAV <= 0 {
		return errors.New("invalid 159352 estimated NAV")
	}
	mid := (market.Bid + market.Ask) / 2
	snapshot.MarketPrice = market.Price
	snapshot.MarketPriceDate = market.Date
	snapshot.Bid = market.Bid
	snapshot.Ask = market.Ask
	snapshot.BidAskSpread = (market.Ask - market.Bid) / mid
	if market.Open > 0 && market.PreviousClose > 0 {
		snapshot.OpeningGap = market.Open/market.PreviousClose - 1
	}
	snapshot.OfficialNAV = nav.Price
	snapshot.OfficialNAVDate = nav.PriceDate
	snapshot.EstimatedNAV = estimatedNAV
	snapshot.Premium = market.Price/estimatedNAV - 1
	return nil
}

func a500PerformanceOnOrBefore(points []a500PerformancePoint, target string) (a500PerformancePoint, bool) {
	for index := len(points) - 1; index >= 0; index-- {
		if points[index].Date <= target {
			return points[index], true
		}
	}
	return a500PerformancePoint{}, false
}
