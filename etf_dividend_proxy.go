package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	dividendLowVolFundCode             = "515450"
	dividendLowVolHistoryVersion       = 1
	dividendLowVolHistoryYears         = 5
	dividendLowVolMinimumObservations  = 230
	dividendLowVolMinimumCoverage      = 0.95
	dividendLowVolFetchWorkers         = 8
	dividendLowVolBasketLookbackDays   = 14
	southernETFPCFAPIURL               = "https://www.southernfund.com/nfwebApi/trade/subAndRedempList"
	southernETFPCFPageURL              = "https://www.southernfund.com/new/transaction-guide/detailed-list.html"
	eastmoneyValueAnalysisAPIURL       = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	eastmoneyValueAnalysisPageURL      = "https://data.eastmoney.com/gzfx/"
	eastmoneyShareBonusPageURL         = "https://data.eastmoney.com/yjfp/"
	dividendLowVolHistoryCacheFileName = "dividend_low_vol_history.json"
)

type dividendLowVolIndexValuation struct {
	Date                    string
	DividendYield           float64
	DividendYieldPercentile float64
	PB                      float64
	PBPercentile            float64
	BondYield               float64
	BondDate                string
	Spread                  float64
	SpreadPercentile        float64
	ValuationScore          float64
	Coverage                float64
	ComponentCount          int
	ValidComponentCount     int
	ObservationCount        int
}

type dividendLowVolProxyHistory struct {
	Version   int                               `json:"version"`
	UpdatedAt string                            `json:"updatedAt,omitempty"`
	Points    []dividendLowVolProxyHistoryPoint `json:"points"`
}

type dividendLowVolProxyHistoryPoint struct {
	Date                string  `json:"date"`
	DividendYield       float64 `json:"dividendYield"`
	PB                  float64 `json:"pb"`
	Coverage            float64 `json:"coverage"`
	ComponentCount      int     `json:"componentCount"`
	ValidComponentCount int     `json:"validComponentCount"`
}

type southernETFBasket struct {
	Date       string
	Components []southernETFBasketComponent
}

type southernETFBasketComponent struct {
	Code     string
	Name     string
	Quantity float64
}

type eastmoneyStockValuation struct {
	Code  string
	Close float64
	PB    float64
}

type stockDividendHistory struct {
	Code   string
	Events []cashDividendEvent
}

func fetchDividendLowVolIndexValuation(client *http.Client, now time.Time) (dividendLowVolIndexValuation, error) {
	history, err := loadDividendLowVolProxyHistory()
	if err != nil {
		return dividendLowVolIndexValuation{}, err
	}
	history, changed, refreshErr := refreshDividendLowVolProxyHistory(client, now, history)
	if changed {
		if err := saveDividendLowVolProxyHistory(history); err != nil {
			return dividendLowVolIndexValuation{}, err
		}
	}
	snapshot, err := calculateDividendLowVolIndexValuation(client, history.Points)
	if err != nil {
		if refreshErr != nil {
			return dividendLowVolIndexValuation{}, fmt.Errorf("%v; refresh: %w", err, refreshErr)
		}
		return dividendLowVolIndexValuation{}, err
	}
	return snapshot, nil
}

func dividendLowVolProxyHistoryPath() string {
	return filepath.Join(filepath.Dir(runtimeQuotesFile), dividendLowVolHistoryCacheFileName)
}

func loadDividendLowVolProxyHistory() (dividendLowVolProxyHistory, error) {
	history := dividendLowVolProxyHistory{Version: dividendLowVolHistoryVersion, Points: []dividendLowVolProxyHistoryPoint{}}
	body, err := os.ReadFile(dividendLowVolProxyHistoryPath())
	if errors.Is(err, os.ErrNotExist) {
		return history, nil
	}
	if err != nil {
		return history, err
	}
	if err := json.Unmarshal(body, &history); err != nil {
		return dividendLowVolProxyHistory{}, fmt.Errorf("parse dividend low-vol history: %w", err)
	}
	if history.Version != dividendLowVolHistoryVersion {
		return dividendLowVolProxyHistory{Version: dividendLowVolHistoryVersion, Points: []dividendLowVolProxyHistoryPoint{}}, nil
	}
	history.Points = normalizeDividendLowVolProxyPoints(history.Points)
	return history, nil
}

func saveDividendLowVolProxyHistory(history dividendLowVolProxyHistory) error {
	history.Version = dividendLowVolHistoryVersion
	history.Points = normalizeDividendLowVolProxyPoints(history.Points)
	body, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return writeFileAtomic(dividendLowVolProxyHistoryPath(), body, 0o644)
}

func normalizeDividendLowVolProxyPoints(points []dividendLowVolProxyHistoryPoint) []dividendLowVolProxyHistoryPoint {
	byWeek := map[string]dividendLowVolProxyHistoryPoint{}
	for _, point := range points {
		date, err := time.Parse(etfRuleRuntimeTimestampDateLayout, strings.TrimSpace(point.Date))
		if err != nil || point.DividendYield <= 0 || point.PB <= 0 || point.Coverage < dividendLowVolMinimumCoverage {
			continue
		}
		key := isoWeekKey(date)
		if previous, ok := byWeek[key]; !ok || point.Date > previous.Date {
			byWeek[key] = point
		}
	}
	result := make([]dividendLowVolProxyHistoryPoint, 0, len(byWeek))
	for _, point := range byWeek {
		result = append(result, point)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Date < result[j].Date })
	return result
}

func isoWeekKey(date time.Time) string {
	year, week := date.ISOWeek()
	return fmt.Sprintf("%04d-%02d", year, week)
}

func refreshDividendLowVolProxyHistory(client *http.Client, now time.Time, history dividendLowVolProxyHistory) (dividendLowVolProxyHistory, bool, error) {
	if client == nil {
		return history, false, errors.New("missing HTTP client")
	}
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate := nowDate.AddDate(-dividendLowVolHistoryYears, 0, 0)
	targets := dividendLowVolHistoryTargets(nowDate, startDate, history)
	if len(targets) == 0 {
		return history, false, nil
	}

	baskets, basketErrors := fetchSouthernETFBaskets(client, targets)
	if len(baskets) == 0 {
		return history, false, fmt.Errorf("missing 515450 baskets: %s", summarizeDividendLowVolRefreshErrors(basketErrors))
	}
	valuationDates := make([]string, 0, len(baskets))
	seenDates := map[string]struct{}{}
	stockCodes := map[string]struct{}{}
	for _, basket := range baskets {
		if _, seen := seenDates[basket.Date]; !seen {
			seenDates[basket.Date] = struct{}{}
			valuationDates = append(valuationDates, basket.Date)
		}
		for _, component := range basket.Components {
			stockCodes[component.Code] = struct{}{}
		}
	}
	valuations, valuationErrors := fetchEastmoneyStockValuationBatches(client, valuationDates)
	dividends, dividendErrors := fetchStockDividendHistories(client, stockCodes)

	pointsByWeek := map[string]dividendLowVolProxyHistoryPoint{}
	for _, point := range history.Points {
		date, err := time.Parse(etfRuleRuntimeTimestampDateLayout, point.Date)
		if err == nil {
			pointsByWeek[isoWeekKey(date)] = point
		}
	}
	calculationErrors := []string{}
	added := 0
	for _, basket := range baskets {
		point, err := calculateDividendLowVolProxyPoint(basket, valuations[basket.Date], dividends)
		if err != nil {
			calculationErrors = append(calculationErrors, basket.Date+": "+err.Error())
			continue
		}
		date, _ := time.Parse(etfRuleRuntimeTimestampDateLayout, point.Date)
		key := isoWeekKey(date)
		if previous, ok := pointsByWeek[key]; !ok || point.Date >= previous.Date {
			pointsByWeek[key] = point
			added++
		}
	}
	if added == 0 {
		errs := append(append(append(basketErrors, valuationErrors...), dividendErrors...), calculationErrors...)
		return history, false, fmt.Errorf("no dividend low-vol history points calculated: %s", summarizeDividendLowVolRefreshErrors(errs))
	}
	history.Points = history.Points[:0]
	for _, point := range pointsByWeek {
		history.Points = append(history.Points, point)
	}
	history.Points = normalizeDividendLowVolProxyPoints(history.Points)
	history.UpdatedAt = now.Format("2006-01-02 15:04:05")

	allErrors := append(append(append(basketErrors, valuationErrors...), dividendErrors...), calculationErrors...)
	if len(allErrors) > 0 && len(history.Points) < dividendLowVolMinimumObservations {
		return history, true, fmt.Errorf("partial dividend low-vol history (%d points): %s", len(history.Points), summarizeDividendLowVolRefreshErrors(allErrors))
	}
	return history, true, nil
}

func dividendLowVolHistoryTargets(nowDate time.Time, startDate time.Time, history dividendLowVolProxyHistory) []time.Time {
	existingWeeks := map[string]dividendLowVolProxyHistoryPoint{}
	for _, point := range history.Points {
		date, err := time.Parse(etfRuleRuntimeTimestampDateLayout, point.Date)
		if err == nil {
			existingWeeks[isoWeekKey(date)] = point
		}
	}
	historyReady := validateDividendLowVolProxyHistory(normalizeDividendLowVolProxyPoints(history.Points)) == nil
	lastRefresh := parseDividendLowVolHistoryRefreshTime(history.UpdatedAt, nowDate.Location())
	targets := []time.Time{}
	for cursor := nowDate; !cursor.Before(startDate); cursor = cursor.AddDate(0, 0, -7) {
		key := isoWeekKey(cursor)
		point, exists := existingWeeks[key]
		isCurrentWeek := key == isoWeekKey(nowDate)
		alreadyScannedHistoricalGap := !exists && !isCurrentWeek && historyReady && !lastRefresh.IsZero() && !cursor.After(lastRefresh)
		if alreadyScannedHistoricalGap {
			continue
		}
		if !exists || (isCurrentWeek && point.Date < nowDate.Format(etfRuleRuntimeTimestampDateLayout)) {
			targets = append(targets, cursor)
		}
	}
	return targets
}

func parseDividendLowVolHistoryRefreshTime(value string, location *time.Location) time.Time {
	for _, layout := range []string{"2006-01-02 15:04:05", etfRuleRuntimeTimestampDateLayout} {
		if parsed, err := time.ParseInLocation(layout, strings.TrimSpace(value), location); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

func summarizeDividendLowVolRefreshErrors(errs []string) string {
	if len(errs) == 0 {
		return "no error details"
	}
	first := strings.TrimSpace(errs[0])
	const maxLength = 180
	if len(first) > maxLength {
		first = first[:maxLength] + "..."
	}
	return fmt.Sprintf("%d fetch or calculation errors; first: %s", len(errs), first)
}

func fetchSouthernETFBaskets(client *http.Client, targets []time.Time) ([]southernETFBasket, []string) {
	type result struct {
		Basket southernETFBasket
		Err    error
		Target string
	}
	jobs := make(chan time.Time)
	results := make(chan result, len(targets))
	workers := minInt(dividendLowVolFetchWorkers, len(targets))
	var wait sync.WaitGroup
	for i := 0; i < workers; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for target := range jobs {
				basket, err := fetchSouthernETFBasketOnOrBefore(client, target)
				results <- result{Basket: basket, Err: err, Target: target.Format(etfRuleRuntimeTimestampDateLayout)}
			}
		}()
	}
	go func() {
		for _, target := range targets {
			jobs <- target
		}
		close(jobs)
		wait.Wait()
		close(results)
	}()

	basketsByDate := map[string]southernETFBasket{}
	errs := []string{}
	for item := range results {
		if item.Err != nil {
			errs = append(errs, item.Target+": "+item.Err.Error())
			continue
		}
		basketsByDate[item.Basket.Date] = item.Basket
	}
	baskets := make([]southernETFBasket, 0, len(basketsByDate))
	for _, basket := range basketsByDate {
		baskets = append(baskets, basket)
	}
	sort.Slice(baskets, func(i, j int) bool { return baskets[i].Date < baskets[j].Date })
	return baskets, errs
}

func fetchSouthernETFBasketOnOrBefore(client *http.Client, target time.Time) (southernETFBasket, error) {
	errs := []string{}
	for lag := 0; lag <= dividendLowVolBasketLookbackDays; lag++ {
		queryDate := target.AddDate(0, 0, -lag)
		basket, err := fetchSouthernETFBasket(client, queryDate.Format("20060102"))
		if err == nil {
			return basket, nil
		}
		errs = append(errs, err.Error())
	}
	return southernETFBasket{}, errors.New(strings.Join(errs, "; "))
}

func fetchSouthernETFBasket(client *http.Client, queryDate string) (southernETFBasket, error) {
	values := url.Values{}
	values.Set("fundCode", dividendLowVolFundCode)
	values.Set("queryDate", strings.TrimSpace(queryDate))
	req, err := http.NewRequest(http.MethodPost, southernETFPCFAPIURL, strings.NewReader(values.Encode()))
	if err != nil {
		return southernETFBasket{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", southernETFPCFPageURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return southernETFBasket{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return southernETFBasket{}, fmt.Errorf("Southern Fund PCF request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var payload struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			TradingDay string `json:"TradingDay"`
			List       []struct {
				StockCode    string `json:"stockCode"`
				StockName    string `json:"stockName"`
				StockQuality string `json:"stockQuality"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&payload); err != nil {
		return southernETFBasket{}, err
	}
	if payload.Code != "ETS-5BP00000" {
		return southernETFBasket{}, fmt.Errorf("Southern Fund PCF error: %s %s", payload.Code, payload.Message)
	}
	date := compactDate(payload.Data.TradingDay)
	if date == "" || len(payload.Data.List) < 48 {
		return southernETFBasket{}, fmt.Errorf("incomplete 515450 basket for %s", queryDate)
	}
	basket := southernETFBasket{Date: date, Components: make([]southernETFBasketComponent, 0, len(payload.Data.List))}
	for _, row := range payload.Data.List {
		code := strings.TrimSpace(row.StockCode)
		quantity, err := strconv.ParseFloat(strings.TrimSpace(row.StockQuality), 64)
		if code == "" || err != nil || quantity < 0 {
			continue
		}
		basket.Components = append(basket.Components, southernETFBasketComponent{Code: code, Name: strings.TrimSpace(row.StockName), Quantity: quantity})
	}
	if len(basket.Components) < 48 {
		return southernETFBasket{}, fmt.Errorf("insufficient 515450 basket components for %s", date)
	}
	return basket, nil
}

func compactDate(value string) string {
	date, err := time.Parse("20060102", strings.TrimSpace(value))
	if err != nil {
		return ""
	}
	return date.Format(etfRuleRuntimeTimestampDateLayout)
}

func fetchEastmoneyStockValuationBatches(client *http.Client, dates []string) (map[string]map[string]eastmoneyStockValuation, []string) {
	type result struct {
		Date string
		Rows map[string]eastmoneyStockValuation
		Err  error
	}
	jobs := make(chan string)
	results := make(chan result, len(dates))
	workers := minInt(dividendLowVolFetchWorkers, len(dates))
	var wait sync.WaitGroup
	for i := 0; i < workers; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			for date := range jobs {
				rows, err := fetchEastmoneyStockValuations(client, date)
				results <- result{Date: date, Rows: rows, Err: err}
			}
		}()
	}
	go func() {
		for _, date := range dates {
			jobs <- date
		}
		close(jobs)
		wait.Wait()
		close(results)
	}()
	all := map[string]map[string]eastmoneyStockValuation{}
	errs := []string{}
	for item := range results {
		if item.Err != nil {
			errs = append(errs, item.Date+": "+item.Err.Error())
			continue
		}
		all[item.Date] = item.Rows
	}
	return all, errs
}

func fetchEastmoneyStockValuations(client *http.Client, date string) (map[string]eastmoneyStockValuation, error) {
	if _, err := time.Parse(etfRuleRuntimeTimestampDateLayout, strings.TrimSpace(date)); err != nil {
		return nil, fmt.Errorf("invalid valuation date: %s", date)
	}
	values := url.Values{}
	values.Set("reportName", "RPT_VALUEANALYSIS_DET")
	values.Set("columns", "SECURITY_CODE,CLOSE_PRICE,PB_MRQ")
	values.Set("pageNumber", "1")
	values.Set("pageSize", "10000")
	values.Set("sortColumns", "SECURITY_CODE")
	values.Set("sortTypes", "1")
	values.Set("source", "WEB")
	values.Set("client", "WEB")
	values.Set("filter", "(TRADE_DATE='"+date+"')")
	req, err := http.NewRequest(http.MethodGet, eastmoneyValueAnalysisAPIURL+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", eastmoneyValueAnalysisPageURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("Eastmoney valuation request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var payload struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Result  struct {
			Data []struct {
				Code  string  `json:"SECURITY_CODE"`
				Close float64 `json:"CLOSE_PRICE"`
				PB    float64 `json:"PB_MRQ"`
			} `json:"data"`
		} `json:"result"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 12<<20)).Decode(&payload); err != nil {
		return nil, err
	}
	if !payload.Success || len(payload.Result.Data) == 0 {
		return nil, fmt.Errorf("Eastmoney valuation error: %s", payload.Message)
	}
	rows := make(map[string]eastmoneyStockValuation, len(payload.Result.Data))
	for _, row := range payload.Result.Data {
		code := strings.TrimSpace(row.Code)
		if code == "" || row.Close <= 0 {
			continue
		}
		rows[code] = eastmoneyStockValuation{Code: code, Close: row.Close, PB: row.PB}
	}
	return rows, nil
}

func fetchStockDividendHistories(client *http.Client, codes map[string]struct{}) (map[string][]cashDividendEvent, []string) {
	codeList := make([]string, 0, len(codes))
	for code := range codes {
		codeList = append(codeList, code)
	}
	sort.Strings(codeList)
	histories := make(map[string][]cashDividendEvent, len(codeList))
	for _, code := range codeList {
		histories[code] = []cashDividendEvent{}
	}
	errs := []string{}
	const batchSize = 100
	for start := 0; start < len(codeList); start += batchSize {
		end := minInt(start+batchSize, len(codeList))
		batch, err := fetchEastmoneyStockDividendBatch(client, codeList[start:end])
		if err != nil {
			errs = append(errs, fmt.Sprintf("codes %d-%d: %v", start+1, end, err))
			continue
		}
		for code, events := range batch {
			histories[code] = events
		}
	}
	return histories, errs
}

func fetchEastmoneyStockDividends(client *http.Client, code string) ([]cashDividendEvent, error) {
	code = strings.TrimSpace(code)
	batch, err := fetchEastmoneyStockDividendBatch(client, []string{code})
	if err != nil {
		return nil, err
	}
	return batch[code], nil
}

func fetchEastmoneyStockDividendBatch(client *http.Client, codes []string) (map[string][]cashDividendEvent, error) {
	quotedCodes := make([]string, 0, len(codes))
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if len(code) != 6 {
			continue
		}
		quotedCodes = append(quotedCodes, `"`+code+`"`)
	}
	if len(quotedCodes) == 0 {
		return nil, errors.New("missing stock codes for Eastmoney dividend request")
	}
	values := url.Values{}
	values.Set("reportName", "RPT_SHAREBONUS_DET")
	values.Set("columns", "SECURITY_CODE,PRETAX_BONUS_RMB,EX_DIVIDEND_DATE")
	values.Set("pageNumber", "1")
	values.Set("pageSize", "10000")
	values.Set("sortTypes", "1,-1")
	values.Set("sortColumns", "SECURITY_CODE,REPORT_DATE")
	values.Set("source", "WEB")
	values.Set("client", "WEB")
	values.Set("filter", `(SECURITY_CODE in (`+strings.Join(quotedCodes, ",")+`))`)
	req, err := http.NewRequest(http.MethodGet, eastmoneyValueAnalysisAPIURL+"?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", eastmoneyShareBonusPageURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("Eastmoney dividend request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var payload struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Result  struct {
			Data []struct {
				Code   string  `json:"SECURITY_CODE"`
				Amount float64 `json:"PRETAX_BONUS_RMB"`
				Date   string  `json:"EX_DIVIDEND_DATE"`
			} `json:"data"`
		} `json:"result"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&payload); err != nil {
		return nil, err
	}
	if !payload.Success {
		return nil, fmt.Errorf("Eastmoney dividend error: %s", payload.Message)
	}
	histories := make(map[string][]cashDividendEvent, len(codes))
	for _, code := range codes {
		histories[strings.TrimSpace(code)] = []cashDividendEvent{}
	}
	for _, row := range payload.Result.Data {
		code := strings.TrimSpace(row.Code)
		date := normalizeEastmoneyDate(row.Date)
		if code == "" || date == "" || row.Amount <= 0 {
			continue
		}
		histories[code] = append(histories[code], cashDividendEvent{Date: date, Amount: row.Amount / 10})
	}
	for code := range histories {
		events := histories[code]
		sort.Slice(events, func(i, j int) bool { return events[i].Date < events[j].Date })
		histories[code] = events
	}
	return histories, nil
}

func normalizeEastmoneyDate(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 10 {
		value = value[:10]
	}
	date, err := time.Parse(etfRuleRuntimeTimestampDateLayout, value)
	if err != nil {
		return ""
	}
	return date.Format(etfRuleRuntimeTimestampDateLayout)
}

func calculateDividendLowVolProxyPoint(basket southernETFBasket, valuations map[string]eastmoneyStockValuation, dividends map[string][]cashDividendEvent) (dividendLowVolProxyHistoryPoint, error) {
	date, err := time.Parse(etfRuleRuntimeTimestampDateLayout, basket.Date)
	if err != nil {
		return dividendLowVolProxyHistoryPoint{}, err
	}
	if len(basket.Components) == 0 || len(valuations) == 0 {
		return dividendLowVolProxyHistoryPoint{}, errors.New("missing basket or valuations")
	}
	cutoff := date.AddDate(-1, 0, 0).Format(etfRuleRuntimeTimestampDateLayout)
	eligibleValue := 0.0
	validValue := 0.0
	dividendNumerator := 0.0
	pbDenominator := 0.0
	validCount := 0
	for _, component := range basket.Components {
		valuation, ok := valuations[component.Code]
		if !ok || component.Quantity <= 0 || valuation.Close <= 0 {
			continue
		}
		marketValue := component.Quantity * valuation.Close
		eligibleValue += marketValue
		events, dividendOK := dividends[component.Code]
		if !dividendOK || valuation.PB <= 0 {
			continue
		}
		trailingDividend := 0.0
		for _, event := range events {
			if event.Date > cutoff && event.Date <= basket.Date {
				trailingDividend += event.Amount
			}
		}
		validValue += marketValue
		dividendNumerator += marketValue * trailingDividend / valuation.Close
		pbDenominator += marketValue / valuation.PB
		validCount++
	}
	componentCoverage := float64(validCount) / float64(len(basket.Components))
	weightCoverage := 0.0
	if eligibleValue > 0 {
		weightCoverage = validValue / eligibleValue
	}
	coverage := math.Min(componentCoverage, weightCoverage)
	if coverage < dividendLowVolMinimumCoverage {
		return dividendLowVolProxyHistoryPoint{}, fmt.Errorf("basket coverage %.2f%% below %.0f%%", coverage*100, dividendLowVolMinimumCoverage*100)
	}
	if validValue <= 0 || pbDenominator <= 0 {
		return dividendLowVolProxyHistoryPoint{}, errors.New("invalid basket valuation totals")
	}
	dividendYield := dividendNumerator / validValue
	pb := validValue / pbDenominator
	if dividendYield <= 0 || dividendYield > 0.20 || pb <= 0.05 || pb > 20 {
		return dividendLowVolProxyHistoryPoint{}, fmt.Errorf("basket valuation out of range: yield %.4f PB %.4f", dividendYield, pb)
	}
	return dividendLowVolProxyHistoryPoint{
		Date:                basket.Date,
		DividendYield:       dividendYield,
		PB:                  pb,
		Coverage:            coverage,
		ComponentCount:      len(basket.Components),
		ValidComponentCount: validCount,
	}, nil
}

func calculateDividendLowVolIndexValuation(client *http.Client, points []dividendLowVolProxyHistoryPoint) (dividendLowVolIndexValuation, error) {
	points = normalizeDividendLowVolProxyPoints(points)
	if err := validateDividendLowVolProxyHistory(points); err != nil {
		return dividendLowVolIndexValuation{}, err
	}
	latest := points[len(points)-1]
	latestDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, latest.Date)
	if err != nil {
		return dividendLowVolIndexValuation{}, err
	}

	bondHistory, err := fetchEastmoneyChina10YBondYieldHistory(client, points[0].Date)
	if err != nil {
		return dividendLowVolIndexValuation{}, err
	}
	officialBond, err := fetchChinaBondOfficial10YYield(client, latestDate)
	if err != nil {
		return dividendLowVolIndexValuation{}, err
	}
	yieldHistory := make([]datedRate, 0, len(points))
	dividendValues := make([]float64, 0, len(points))
	pbValues := make([]float64, 0, len(points))
	for _, point := range points {
		yieldHistory = append(yieldHistory, datedRate{Date: point.Date, Value: point.DividendYield})
		dividendValues = append(dividendValues, point.DividendYield)
		pbValues = append(pbValues, point.PB)
	}
	spread, spreadPercentile, observations, err := calculateDividendSpread(
		datedRate{Date: latest.Date, Value: latest.DividendYield},
		yieldHistory,
		bondHistory,
		officialBond,
	)
	if err != nil {
		return dividendLowVolIndexValuation{}, err
	}
	pbPercentile := percentileRank(latest.PB, pbValues)
	spreadPercentile = clampUnitValue(spreadPercentile)
	pbPercentile = clampUnitValue(pbPercentile)
	return dividendLowVolIndexValuation{
		Date:                    latest.Date,
		DividendYield:           latest.DividendYield,
		DividendYieldPercentile: clampUnitValue(percentileRank(latest.DividendYield, dividendValues)),
		PB:                      latest.PB,
		PBPercentile:            pbPercentile,
		BondYield:               officialBond.Value,
		BondDate:                officialBond.Date,
		Spread:                  spread,
		SpreadPercentile:        spreadPercentile,
		ValuationScore:          dividendLowVolValuationScore(spreadPercentile, pbPercentile),
		Coverage:                latest.Coverage,
		ComponentCount:          latest.ComponentCount,
		ValidComponentCount:     latest.ValidComponentCount,
		ObservationCount:        observations,
	}, nil
}

func validateDividendLowVolProxyHistory(points []dividendLowVolProxyHistoryPoint) error {
	if len(points) < dividendLowVolMinimumObservations {
		return fmt.Errorf("insufficient 515450 basket history: %d weekly observations", len(points))
	}
	oldestDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, points[0].Date)
	if err != nil {
		return err
	}
	latest := points[len(points)-1]
	latestDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, latest.Date)
	if err != nil {
		return err
	}
	if oldestDate.After(latestDate.AddDate(-dividendLowVolHistoryYears, 0, 7)) {
		return fmt.Errorf("515450 basket history does not span five years: %s to %s", points[0].Date, latest.Date)
	}
	return nil
}

func clampUnitValue(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func minInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
}
