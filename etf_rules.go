package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ETFRuleStatus struct {
	Symbol        string          `json:"symbol"`
	Name          string          `json:"name"`
	Level         string          `json:"level,omitempty"`
	LevelLabel    string          `json:"levelLabel,omitempty"`
	MonthlyAmount float64         `json:"monthlyAmount,omitempty"`
	WeeklyAmount  float64         `json:"weeklyAmount,omitempty"`
	Complete      bool            `json:"complete"`
	Reason        string          `json:"reason,omitempty"`
	AsOf          string          `json:"asOf,omitempty"`
	UpdatedAt     string          `json:"updatedAt,omitempty"`
	Metrics       []ETFRuleMetric `json:"metrics,omitempty"`
	Sources       []ETFRuleSource `json:"sources,omitempty"`
}

type ETFRuleMetric struct {
	Key       string   `json:"key"`
	Label     string   `json:"label"`
	Value     *float64 `json:"value,omitempty"`
	Unit      string   `json:"unit,omitempty"`
	AsOf      string   `json:"asOf,omitempty"`
	Available bool     `json:"available"`
	Error     string   `json:"error,omitempty"`
}

type ETFRuleSource struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

const (
	etfRuleDailyMetricMaxAgeDays      = 7
	etfRuleMonthlyMetricMaxAgeDays    = 60
	etfRuleRuntimeTimestampDateLayout = "2006-01-02"
	historyOfMarketSP500PEURL         = "https://historyofmarket.com/api/sp500/pe.json"
	historyOfMarketNDXForwardPEURL    = "https://historyofmarket.com/api/ndx/forward-pe.json"
	multplShillerCAPEURL              = "https://www.multpl.com/shiller-pe/table/by-year"
	fundDBAPIHost                     = "https://api.jiucaishuo.com"
	fundDBIndexPageURL                = "https://funddb.cn/site/index"
	fundDBAPIVersion                  = "2.2.7"
	fundDBAPIReqKey                   = "EWf45rlv#kfsr@k#gfksgkr"
	primaryValuationMaxLagDays        = 3
)

type etfRuleLevel struct {
	Key   string
	Label string
}

type etfRuleConfig struct {
	Symbol              string
	Name                string
	PriceSymbol         string
	PriceSourceName     string
	PriceSourceURL      string
	ValuationMetricKey  string
	ValuationMetricName string
	ValuationSourceName string
	ValuationSourceURL  string
	Levels              map[string]etfRuleLevel
	Monthly             map[string]float64
	Weekly              map[string]float64
	Evaluate            func(etfRuleInputs) etfRuleEvaluation
}

type etfRuleInputs struct {
	Drawdown            *float64
	DrawdownAsOf        string
	ValuationPercentile *float64
	ValuationZScore     *float64
	ValuationAsOf       string
}

type etfRuleEvaluation struct {
	Level    string
	Complete bool
	Reason   string
}

var etfRuleLevels = map[string]etfRuleLevel{
	"quarter": {Key: "quarter", Label: "0.25倍"},
	"half":    {Key: "half", Label: "0.5倍"},
	"one":     {Key: "one", Label: "1倍"},
	"oneHalf": {Key: "oneHalf", Label: "1.5倍"},
	"two":     {Key: "two", Label: "2倍"},
}

var etfRuleConfigs = []etfRuleConfig{
	{
		Symbol:              "022434",
		Name:                "南方中证A500ETF联接A",
		PriceSymbol:         "000510.SH",
		PriceSourceName:     "东方财富沪深行情K线",
		PriceSourceURL:      "https://quote.eastmoney.com/",
		ValuationMetricKey:  "pePercentile",
		ValuationMetricName: "中证A500 PE分位",
		ValuationSourceName: "韭圈儿中证A500 PE分位（优先；乐咕乐股备援）",
		ValuationSourceURL:  fundDBIndexPageURL,
		Levels:              etfRuleLevels,
		Monthly:             map[string]float64{"quarter": 2800, "half": 5600, "one": 11200, "oneHalf": 16800, "two": 22400},
		Weekly:              map[string]float64{"quarter": 700, "half": 1400, "one": 2800, "oneHalf": 4200, "two": 5600},
		Evaluate:            evaluateA500Rule,
	},
	{
		Symbol:              "018738",
		Name:                "博时标普500ETF联接E(人民币)",
		PriceSymbol:         "^GSPC",
		PriceSourceName:     "东方财富/Yahoo/Nasdaq/Stooq标普500日线（自动选最新）",
		PriceSourceURL:      "https://finance.yahoo.com/quote/%5EGSPC/history/",
		ValuationMetricKey:  "pePercentile",
		ValuationMetricName: "标普500 PE分位",
		ValuationSourceName: "韭圈儿标普500 PE分位（优先；History of Market CAPE备援）",
		ValuationSourceURL:  fundDBIndexPageURL,
		Levels:              etfRuleLevels,
		Monthly:             map[string]float64{"quarter": 4200, "half": 8400, "one": 16800, "oneHalf": 25200, "two": 33600},
		Weekly:              map[string]float64{"quarter": 1050, "half": 2100, "one": 4200, "oneHalf": 6300, "two": 8400},
		Evaluate:            evaluateSP500Rule,
	},
	{
		Symbol:              "008163",
		Name:                "南方标普红利低波50ETF联接A",
		PriceSymbol:         "515450.SH",
		PriceSourceName:     "东方财富沪深行情K线",
		PriceSourceURL:      "https://quote.eastmoney.com/sh515450.html",
		ValuationMetricKey:  "dividendYield",
		ValuationMetricName: "515450股息率",
		ValuationSourceName: "天天基金/东方财富515450分红TTM备援口径",
		ValuationSourceURL:  "https://fundf10.eastmoney.com/fhsp_515450.html",
		Levels:              etfRuleLevels,
		Monthly:             map[string]float64{"quarter": 4200, "half": 8400, "one": 16800, "oneHalf": 25200, "two": 33600},
		Weekly:              map[string]float64{"quarter": 1050, "half": 2100, "one": 4200, "oneHalf": 6300, "two": 8400},
		Evaluate:            evaluateDividendLowVolRule,
	},
	{
		Symbol:              "021000",
		Name:                "南方纳斯达克100指数发起(QDII)I",
		PriceSymbol:         "^NDX",
		PriceSourceName:     "Nasdaq官方NDX历史/报价、Yahoo、东方财富、Stooq纳指100日线（自动选最新）",
		PriceSourceURL:      "https://finance.yahoo.com/quote/%5ENDX/history/",
		ValuationMetricKey:  "pePercentile",
		ValuationMetricName: "纳指100 PE分位",
		ValuationSourceName: "韭圈儿纳斯达克100 PE分位（优先；History of Market Forward PE备援）",
		ValuationSourceURL:  fundDBIndexPageURL,
		Levels:              etfRuleLevels,
		Monthly:             map[string]float64{"quarter": 2800, "half": 5600, "one": 11200, "oneHalf": 16800, "two": 22400},
		Weekly:              map[string]float64{"quarter": 700, "half": 1400, "one": 2800, "oneHalf": 4200, "two": 5600},
		Evaluate:            evaluateNasdaq100Rule,
	},
}

func updateETFRuleStatuses(client *http.Client, now time.Time) ([]ETFRuleStatus, []QuoteSkip) {
	statuses := make([]ETFRuleStatus, 0, len(etfRuleConfigs))
	skipped := []QuoteSkip{}
	for _, config := range etfRuleConfigs {
		status, err := fetchETFRuleStatus(client, config, now)
		if err != nil {
			skipped = append(skipped, QuoteSkip{Type: "etf-rule", Symbol: config.Symbol, Name: config.Name, Error: err.Error()})
		}
		statuses = append(statuses, status)
	}
	return statuses, skipped
}

func fetchETFRuleStatus(client *http.Client, config etfRuleConfig, now time.Time) (ETFRuleStatus, error) {
	inputs := etfRuleInputs{}
	metrics := []ETFRuleMetric{}
	sources := []ETFRuleSource{{Name: config.PriceSourceName, URL: config.PriceSourceURL}}
	statusErrs := []string{}

	drawdown, drawdownDate, err := fetchETFRuleDrawdown(client, config)
	if err != nil {
		statusErrs = append(statusErrs, "回撤："+err.Error())
		metrics = append(metrics, ETFRuleMetric{Key: "drawdown252", Label: "近252交易日回撤", Unit: "%", Available: false, Error: err.Error()})
	} else {
		inputs.Drawdown = &drawdown
		inputs.DrawdownAsOf = drawdownDate
		metrics = append(metrics, ETFRuleMetric{Key: "drawdown252", Label: "近252交易日回撤", Value: percentMetric(drawdown), Unit: "%", AsOf: drawdownDate, Available: true})
	}

	valuation, valuationErr := fetchETFRuleValuation(client, config)
	if valuationErr != nil {
		statusErrs = append(statusErrs, config.ValuationMetricName+"："+valuationErr.Error())
		metrics = append(metrics, ETFRuleMetric{Key: config.ValuationMetricKey, Label: config.ValuationMetricName, Unit: configValuationMetricUnit(config), Available: false, Error: valuationErr.Error()})
	} else {
		inputs.ValuationAsOf = valuation.Date
		metricValue := valuation.Value
		if valuation.Kind == "zScore" {
			inputs.ValuationZScore = &valuation.Value
		} else {
			inputs.ValuationPercentile = &valuation.Value
			metricValue = valuation.Value * 100
		}
		metrics = append(metrics, ETFRuleMetric{Key: config.ValuationMetricKey, Label: config.ValuationMetricName, Value: floatMetric(metricValue), Unit: valuation.Unit, AsOf: valuation.Date, Available: true})
	}
	if strings.TrimSpace(config.ValuationSourceName) != "" {
		sources = append(sources, ETFRuleSource{Name: config.ValuationSourceName, URL: config.ValuationSourceURL})
	}

	evaluation := config.Evaluate(inputs)
	level := config.Levels[evaluation.Level]
	status := ETFRuleStatus{
		Symbol:        config.Symbol,
		Name:          config.Name,
		Level:         evaluation.Level,
		LevelLabel:    level.Label,
		MonthlyAmount: config.Monthly[evaluation.Level],
		WeeklyAmount:  config.Weekly[evaluation.Level],
		Complete:      evaluation.Complete,
		Reason:        evaluation.Reason,
		AsOf:          firstNonEmpty(inputs.DrawdownAsOf, inputs.ValuationAsOf),
		UpdatedAt:     now.Format("2006-01-02 15:04:05"),
		Metrics:       metrics,
		Sources:       sources,
	}
	if status.Level == "" {
		status.LevelLabel = "待数据"
		status.Reason = firstNonEmpty(evaluation.Reason, strings.Join(statusErrs, "；"))
	}
	if len(statusErrs) > 0 && status.Complete {
		status.Reason = strings.TrimSpace(status.Reason + "；部分辅助指标未取到：" + strings.Join(statusErrs, "；"))
	}
	if status.AsOf == "" {
		status.AsOf = now.Format("2006-01-02")
	}
	status = enforceETFRuleStatusConfidence(status, config, now)
	if len(statusErrs) > 0 {
		return status, errors.New(strings.Join(statusErrs, "；"))
	}
	return status, nil
}

func fetchETFRuleDrawdown(client *http.Client, config etfRuleConfig) (float64, string, error) {
	closes, err := fetchRuleDailyCloses(client, config.PriceSymbol, 280)
	if err != nil {
		return 0, "", err
	}
	return drawdownFromRecentHigh(closes, 252)
}

func fetchRuleDailyCloses(client *http.Client, symbol string, limit int) ([]dailyClose, error) {
	normalized := normalizeSymbol(symbol)
	if strings.HasSuffix(normalized, ".SH") || strings.HasSuffix(normalized, ".SZ") || strings.HasSuffix(normalized, ".HK") {
		closes, err := fetchTencentDailyCloses(client, normalized, limit)
		if err == nil && len(closes) > 0 {
			return closes, nil
		}
		return fetchEastmoneyDailyCloses(client, normalized, limit)
	}
	if secID := eastmoneyGlobalIndexSecID(normalized); secID != "" {
		return fetchGlobalIndexDailyCloses(client, normalized, secID, limit)
	}
	closes, err := fetchYahooDailyCloses(client, normalized, "1y")
	if err == nil && len(closes) > 0 {
		return closes, nil
	}
	if strings.EqualFold(normalized, "^GSPC") || strings.EqualFold(normalized, "SPX") {
		nasdaqCloses, nasdaqErr := fetchNasdaqHistoricalCloses(client, "SPY", "etf", limit)
		if nasdaqErr == nil && len(nasdaqCloses) > 0 {
			return nasdaqCloses, nil
		}
	}
	stooqCloses, stooqErr := fetchStooqDailyCloses(client, normalized)
	if stooqErr == nil && len(stooqCloses) > 0 {
		if limit > 0 && len(stooqCloses) > limit {
			return stooqCloses[len(stooqCloses)-limit:], nil
		}
		return stooqCloses, nil
	}
	return nil, fmt.Errorf("yahoo: %v; stooq: %v", err, stooqErr)
}

type dailyCloseCandidate struct {
	Source string
	Closes []dailyClose
}

func fetchGlobalIndexDailyCloses(client *http.Client, symbol string, eastmoneySecID string, limit int) ([]dailyClose, error) {
	candidates := []dailyCloseCandidate{}
	errs := []string{}

	addCandidate := func(source string, closes []dailyClose, err error) {
		if err != nil {
			errs = append(errs, source+": "+err.Error())
			return
		}
		if len(closes) == 0 {
			errs = append(errs, source+": empty daily closes")
			return
		}
		candidates = append(candidates, dailyCloseCandidate{Source: source, Closes: closes})
	}

	if strings.EqualFold(symbol, "^NDX") || strings.EqualFold(symbol, "NDX") {
		closes, err := fetchNasdaqHistoricalCloses(client, "NDX", "index", limit)
		if err == nil && len(closes) > 0 {
			if latest, latestErr := fetchNasdaqLatestQuoteClose(client, "NDX", "index"); latestErr == nil {
				closes = appendOrReplaceLatestDailyClose(closes, latest)
			}
		}
		addCandidate("nasdaq", closes, err)
	}
	if strings.TrimSpace(eastmoneySecID) != "" {
		closes, err := fetchEastmoneyDailyClosesBySecID(client, eastmoneySecID, limit)
		addCandidate("eastmoney", closes, err)
	}
	closes, err := fetchYahooDailyCloses(client, symbol, "2y")
	addCandidate("yahoo", closes, err)
	if strings.EqualFold(symbol, "^GSPC") || strings.EqualFold(symbol, "SPX") {
		closes, err := fetchNasdaqHistoricalCloses(client, "SPY", "etf", limit)
		addCandidate("nasdaq-spy", closes, err)
	}
	closes, err = fetchStooqDailyCloses(client, symbol)
	addCandidate("stooq", closes, err)

	if closes, _, ok := selectLatestDailyCloseCandidate(candidates, limit); ok {
		return closes, nil
	}
	return nil, errors.New(strings.Join(errs, "; "))
}

func selectLatestDailyCloseCandidate(candidates []dailyCloseCandidate, limit int) ([]dailyClose, string, bool) {
	var best dailyCloseCandidate
	bestDate := ""
	for _, candidate := range candidates {
		date := latestDailyCloseDate(candidate.Closes)
		if date == "" {
			continue
		}
		if bestDate == "" || date > bestDate {
			best = candidate
			bestDate = date
		}
	}
	if bestDate == "" {
		return nil, "", false
	}
	closes := best.Closes
	if limit > 0 && len(closes) > limit {
		closes = closes[len(closes)-limit:]
	}
	return closes, best.Source, true
}

func latestDailyCloseDate(closes []dailyClose) string {
	if len(closes) == 0 {
		return ""
	}
	return strings.TrimSpace(closes[len(closes)-1].Date)
}

func appendOrReplaceLatestDailyClose(closes []dailyClose, latest dailyClose) []dailyClose {
	if latest.Date == "" || latest.Price <= 0 {
		return closes
	}
	if len(closes) == 0 {
		return []dailyClose{latest}
	}
	result := append([]dailyClose(nil), closes...)
	last := result[len(result)-1]
	if latest.Date < last.Date {
		return result
	}
	if latest.Date == last.Date {
		result[len(result)-1] = latest
		return result
	}
	result = append(result, latest)
	return result
}

func eastmoneyGlobalIndexSecID(symbol string) string {
	switch strings.ToUpper(strings.TrimSpace(symbol)) {
	case "^GSPC", "SPX":
		return "100.SPX"
	case "^NDX", "NDX":
		return "100.NDX"
	default:
		return ""
	}
}

func fetchYahooDailyCloses(client *http.Client, symbol string, rangeParam string) ([]dailyClose, error) {
	sourceSymbol := yahooSymbol(symbol)
	if strings.TrimSpace(rangeParam) == "" {
		rangeParam = "1y"
	}
	endpoint := "https://query1.finance.yahoo.com/v8/finance/chart/" + url.PathEscape(sourceSymbol) + "?range=" + url.QueryEscape(rangeParam) + "&interval=1d"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "holds-website etf rule updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("daily close request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload yahooChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Chart.Error != nil || len(payload.Chart.Result) == 0 {
		return nil, errors.New("empty daily close response")
	}
	result := payload.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return nil, errors.New("missing close series")
	}
	location := loadLocation(result.Meta.ExchangeTimezone)
	closes := result.Indicators.Quote[0].Close
	validCloses := make([]dailyClose, 0, len(closes))
	for i, closePrice := range closes {
		if closePrice > 0 {
			validCloses = append(validCloses, dailyClose{Price: closePrice, Date: closeDate(result.Timestamp, i, location)})
		}
	}
	if len(validCloses) == 0 {
		return nil, errors.New("no valid close prices")
	}
	return validCloses, nil
}

func fetchNasdaqHistoricalCloses(client *http.Client, symbol string, assetClass string, limit int) ([]dailyClose, error) {
	if limit <= 0 {
		limit = 280
	}
	toDate := time.Now().Format("2006-01-02")
	fromDate := time.Now().AddDate(-2, 0, 0).Format("2006-01-02")
	values := url.Values{}
	values.Set("assetclass", assetClass)
	values.Set("fromdate", fromDate)
	values.Set("todate", toDate)
	values.Set("limit", strconv.Itoa(limit+80))
	endpoint := "https://api.nasdaq.com/api/quote/" + url.PathEscape(strings.ToUpper(strings.TrimSpace(symbol))) + "/historical?" + values.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.nasdaq.com")
	req.Header.Set("Referer", "https://www.nasdaq.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("nasdaq historical request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	closes, err := parseNasdaqHistoricalCloses(body)
	if err != nil {
		return nil, err
	}
	if len(closes) > limit {
		return closes[len(closes)-limit:], nil
	}
	return closes, nil
}

func fetchNasdaqLatestQuoteClose(client *http.Client, symbol string, assetClass string) (dailyClose, error) {
	values := url.Values{}
	values.Set("assetclass", assetClass)
	endpoint := "https://api.nasdaq.com/api/quote/" + url.PathEscape(strings.ToUpper(strings.TrimSpace(symbol))) + "/info?" + values.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return dailyClose{}, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.nasdaq.com")
	req.Header.Set("Referer", "https://www.nasdaq.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return dailyClose{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return dailyClose{}, fmt.Errorf("nasdaq quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		Data struct {
			PrimaryData struct {
				LastSalePrice      string `json:"lastSalePrice"`
				LastTradeTimestamp string `json:"lastTradeTimestamp"`
			} `json:"primaryData"`
		} `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil {
		return dailyClose{}, err
	}
	price, err := parseMarketNumber(payload.Data.PrimaryData.LastSalePrice)
	if err != nil || price <= 0 {
		return dailyClose{}, errors.New("missing nasdaq quote close price")
	}
	date := normalizeNasdaqQuoteDate(payload.Data.PrimaryData.LastTradeTimestamp)
	if date == "" {
		return dailyClose{}, errors.New("missing nasdaq quote close date")
	}
	return dailyClose{Date: date, Price: price}, nil
}

func parseNasdaqHistoricalCloses(body []byte) ([]dailyClose, error) {
	var payload struct {
		Data struct {
			TradesTable struct {
				Rows []struct {
					Date  string `json:"date"`
					Close string `json:"close"`
				} `json:"rows"`
			} `json:"tradesTable"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	rows := payload.Data.TradesTable.Rows
	if len(rows) == 0 {
		return nil, errors.New("missing nasdaq historical rows")
	}
	closes := make([]dailyClose, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		date := normalizeNasdaqHistoricalDate(rows[i].Date)
		if date == "" {
			continue
		}
		closePrice, err := parseMarketNumber(rows[i].Close)
		if err != nil || closePrice <= 0 {
			continue
		}
		closes = append(closes, dailyClose{Date: date, Price: closePrice})
	}
	if len(closes) == 0 {
		return nil, errors.New("missing nasdaq historical close prices")
	}
	return closes, nil
}

func normalizeNasdaqHistoricalDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"01/02/2006", "1/2/2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format("2006-01-02")
		}
	}
	return ""
}

func normalizeNasdaqQuoteDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"Jan 2, 2006", "Jan 02, 2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format("2006-01-02")
		}
	}
	return ""
}

func parseMarketNumber(value string) (float64, error) {
	cleaned := strings.NewReplacer(",", "", "$", "", " ", "").Replace(strings.TrimSpace(value))
	return strconv.ParseFloat(cleaned, 64)
}

func drawdownFromRecentHigh(closes []dailyClose, window int) (float64, string, error) {
	if len(closes) == 0 {
		return 0, "", errors.New("missing close prices")
	}
	if window <= 0 || window > len(closes) {
		window = len(closes)
	}
	recent := closes[len(closes)-window:]
	latest := recent[len(recent)-1]
	high := 0.0
	for _, close := range recent {
		if close.Price > high {
			high = close.Price
		}
	}
	if high <= 0 || latest.Price <= 0 {
		return 0, "", errors.New("invalid close prices")
	}
	drawdown := (high - latest.Price) / high
	return drawdown, latest.Date, nil
}

func fetchStooqDailyCloses(client *http.Client, symbol string) ([]dailyClose, error) {
	sourceSymbol, err := stooqSymbol(symbol)
	if err != nil {
		return nil, err
	}
	endpoint := "https://stooq.com/q/d/l/?s=" + url.QueryEscape(sourceSymbol) + "&i=d"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "holds-website etf rule updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("stooq request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	return parseStooqDailyCSV(body)
}

func stooqSymbol(symbol string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(symbol)) {
	case "^GSPC", "SPX":
		return "^spx", nil
	case "^NDX", "NDX":
		return "^ndx", nil
	default:
		return "", fmt.Errorf("unsupported stooq symbol: %s", symbol)
	}
}

func parseStooqDailyCSV(body []byte) ([]dailyClose, error) {
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(lines) < 2 {
		return nil, errors.New("empty stooq csv")
	}
	closes := make([]dailyClose, 0, len(lines)-1)
	for _, line := range lines[1:] {
		fields := strings.Split(strings.TrimSpace(line), ",")
		if len(fields) < 5 || strings.EqualFold(fields[4], "null") {
			continue
		}
		date := strings.TrimSpace(fields[0])
		if _, err := time.Parse("2006-01-02", date); err != nil {
			continue
		}
		price, err := strconv.ParseFloat(strings.TrimSpace(fields[4]), 64)
		if err != nil || price <= 0 {
			continue
		}
		closes = append(closes, dailyClose{Date: date, Price: price})
	}
	if len(closes) == 0 {
		return nil, errors.New("missing stooq close prices")
	}
	return closes, nil
}

type etfRuleValuation struct {
	Value float64
	Date  string
	Unit  string
	Kind  string
}

func fetchETFRuleValuation(client *http.Client, config etfRuleConfig) (etfRuleValuation, error) {
	switch config.Symbol {
	case "022434":
		value, date, err := fetchA500PEPercentile(client, time.Now())
		return etfRuleValuation{Value: value, Date: date, Unit: "%", Kind: "percentile"}, err
	case "018738":
		value, date, err := fetchSP500PEPercentile(client)
		return etfRuleValuation{Value: value, Date: date, Unit: "%", Kind: "percentile"}, err
	case "008163":
		value, date, err := fetchDividendLowVolYield(client)
		return etfRuleValuation{Value: value, Date: date, Unit: "%", Kind: "percentile"}, err
	case "021000":
		value, date, err := fetchNasdaq100PEPercentile(client)
		return etfRuleValuation{Value: value, Date: date, Unit: "%", Kind: "percentile"}, err
	default:
		return etfRuleValuation{}, errors.New("valuation source not configured")
	}
}

type fundDBIndexPEPoint struct {
	Percentile float64
	Date       string
}

func fetchFundDBIndexPEPercentile(client *http.Client, indexCode string, category string) (float64, string, error) {
	payload := map[string]any{
		"gu_code":     strings.TrimSpace(indexCode),
		"pe_category": "pe",
		"year":        10,
		"category":    strings.TrimSpace(category),
		"ver":         "new",
	}
	body, err := fundDBPost(client, "/v2/guzhi/newtubiaodata", payload)
	if err != nil {
		return 0, "", err
	}
	point, err := parseFundDBIndexPEPercentile(body)
	if err != nil {
		return 0, "", err
	}
	return point.Percentile, point.Date, nil
}

func parseFundDBIndexPEPercentile(body []byte) (fundDBIndexPEPoint, error) {
	var payload struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			UpdateTime string `json:"update_time"`
			TopData    []struct {
				Attribute       string `json:"attribute"`
				Name            string `json:"name"`
				NewPercentValue struct {
					Value string `json:"value"`
				} `json:"new_percent_value"`
			} `json:"top_data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return fundDBIndexPEPoint{}, err
	}
	if payload.Code != 0 {
		return fundDBIndexPEPoint{}, fmt.Errorf("funddb guzhi response code %d: %s", payload.Code, strings.TrimSpace(payload.Message))
	}
	date := strings.TrimSpace(payload.Data.UpdateTime)
	if _, err := time.Parse(etfRuleRuntimeTimestampDateLayout, date); err != nil {
		return fundDBIndexPEPoint{}, fmt.Errorf("funddb missing update date: %s", date)
	}
	for _, item := range payload.Data.TopData {
		if strings.TrimSpace(item.Attribute) != "pe" && !strings.Contains(item.Name, "市盈率") {
			continue
		}
		percentile, err := parseFundDBPercentile(item.NewPercentValue.Value)
		if err != nil {
			return fundDBIndexPEPoint{}, err
		}
		return fundDBIndexPEPoint{Percentile: percentile, Date: date}, nil
	}
	return fundDBIndexPEPoint{}, errors.New("funddb missing PE percentile")
}

func parseFundDBPercentile(value string) (float64, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "--" {
		return 0, errors.New("funddb percentile is empty")
	}
	number, err := firstTextNumber(trimmed)
	if err != nil {
		return 0, err
	}
	if strings.Contains(trimmed, "%") || number > 1 {
		number /= 100
	}
	if number < 0 || number > 1 || math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, fmt.Errorf("funddb percentile out of range: %s", value)
	}
	return number, nil
}

func fundDBPost(client *http.Client, path string, payload map[string]any) ([]byte, error) {
	signed := fundDBSignedPayload(payload, time.Now())
	body, err := json.Marshal(signed)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fundDBAPIHost+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://funddb.cn")
	req.Header.Set("Referer", fundDBIndexPageURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	requestClient := client
	if requestClient == nil {
		requestClient = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := requestClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("funddb request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return body, nil
}

func fundDBSignedPayload(payload map[string]any, now time.Time) map[string]any {
	signed := map[string]any{}
	for key, value := range payload {
		signed[key] = value
	}
	signed["type"] = "pc"
	signed["version"] = fundDBAPIVersion
	if _, ok := signed["authtoken"]; !ok {
		signed["authtoken"] = ""
	}
	signed["act_time"] = now.UnixMilli()

	keys := make([]string, 0, len(signed))
	for key := range signed {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	signatureText := strings.Builder{}
	for _, key := range keys {
		valueText, include := fundDBSignatureValue(signed[key])
		if include {
			signatureText.WriteString(valueText)
		}
	}
	signatureText.WriteString(fundDBAPIReqKey)
	sum := md5.Sum([]byte(signatureText.String()))
	fundDBApplySignatureFields(signed, fmt.Sprintf("%x", sum))
	return signed
}

func fundDBSignatureValue(value any) (string, bool) {
	switch typed := value.(type) {
	case nil:
		return "", false
	case string:
		return typed, strings.TrimSpace(typed) != ""
	case bool:
		if !typed {
			return "", false
		}
		return "true", true
	case int:
		return strconv.Itoa(typed), true
	case int64:
		return strconv.FormatInt(typed, 10), true
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64), true
	case map[string]any, []any:
		return "", false
	default:
		return fmt.Sprint(typed), true
	}
}

func fundDBApplySignatureFields(payload map[string]any, signature string) {
	if len(signature) < 32 {
		return
	}
	part := func(start int, end int) string { return signature[start:end] }
	payload["tirgkjfs"] = part(0, 2)
	payload["abiokytke"] = part(21, 23)
	payload["u54rg5d"] = part(2, 4)
	payload["kf54ge7"] = part(31, 32)
	payload["tiklsktr4"] = part(1, 2)
	payload["lksytkjh"] = part(17, 21)
	payload["sbnoywr"] = part(23, 25)
	payload["bgd7h8tyu54"] = part(6, 8)
	payload["y654b5fs3tr"] = part(11, 12)
	payload["bioduytlw"] = part(5, 6)
	payload["bd4uy742"] = part(26, 27)
	payload["h67456y"] = part(16, 19)
	payload["bvytikwqjk"] = part(6, 8)
	payload["ngd4uy551"] = part(17, 19)
	payload["bgiuytkw"] = part(9, 11)
	payload["nd354uy4752"] = part(30, 31)
	payload["ghtoiutkmlg"] = part(11, 14)
	payload["bd24y6421f"] = part(24, 26)
	payload["tbvdiuytk"] = part(16, 17)
	payload["ibvytiqjek"] = part(14, 16)
	payload["jnhf8u5231"] = part(9, 11)
	payload["fjlkatj"] = part(2, 5)
	payload["hy5641d321t"] = part(25, 27)
	payload["iogojti"] = part(25, 26)
	payload["ngd4yut78"] = part(12, 14)
	payload["nkjhrew"] = part(26, 27)
	payload["yt447e13f"] = part(8, 9)
	payload["n3bf4uj7y7"] = part(18, 19)
	payload["nbf4uj7y432"] = part(21, 23)
	payload["yi854tew"] = part(29, 31)
	payload["h13ey474"] = part(29, 32)
	payload["quikgdky"] = part(27, 29)
}

type leguleguIndexPERow struct {
	Date          string
	TtmPE         float64
	TtmPEQuantile *float64
}

func fetchA500PEPercentile(client *http.Client, now time.Time) (float64, string, error) {
	percentile, date, fundDBErr := fetchFundDBIndexPEPercentile(client, "000510.SH", "")
	if fundDBErr == nil && !valuationDateStale(date, now, primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	if fundDBErr == nil {
		fundDBErr = fmt.Errorf("funddb PE stale as of %s", date)
	}
	rows, err := fetchLeguleguA500PERows(client, now)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("funddb: %v; legulegu: %v", fundDBErr, err)
	}
	fallbackPercentile, fallbackDate, err := a500PEPercentileFromRows(rows)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("funddb: %v; legulegu: %v", fundDBErr, err)
	}
	if date != "" && !quoteDateAfter(fallbackDate, date) {
		return percentile, date, nil
	}
	return fallbackPercentile, fallbackDate, nil
}

func fetchSP500PEPercentile(client *http.Client) (float64, string, error) {
	now := time.Now()
	percentile, date, fundDBErr := fetchFundDBIndexPEPercentile(client, "SPX.GI", "5")
	if fundDBErr == nil && !valuationDateStale(date, now, primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	if fundDBErr == nil {
		fundDBErr = fmt.Errorf("funddb PE stale as of %s", date)
	}
	fallbackPercentile, fallbackDate, err := fetchSP500CAPEPercentile(client)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("funddb: %v; historyofmarket CAPE: %v", fundDBErr, err)
	}
	if date != "" && !quoteDateAfter(fallbackDate, date) {
		return percentile, date, nil
	}
	return fallbackPercentile, fallbackDate, nil
}

func fetchNasdaq100PEPercentile(client *http.Client) (float64, string, error) {
	now := time.Now()
	percentile, date, fundDBErr := fetchFundDBIndexPEPercentile(client, "NDX.GI", "5")
	if fundDBErr == nil && !valuationDateStale(date, now, primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	if fundDBErr == nil {
		fundDBErr = fmt.Errorf("funddb PE stale as of %s", date)
	}
	fallbackPercentile, fallbackDate, err := fetchNasdaq100ForwardPEPercentile(client)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("funddb: %v; historyofmarket forward PE: %v", fundDBErr, err)
	}
	if date != "" && !quoteDateAfter(fallbackDate, date) {
		return percentile, date, nil
	}
	return fallbackPercentile, fallbackDate, nil
}

func fetchLeguleguA500PEPercentile(client *http.Client, now time.Time) (float64, string, error) {
	rows, err := fetchLeguleguA500PERows(client, now)
	if err != nil {
		return 0, "", err
	}
	return a500PEPercentileFromRows(rows)
}

func fetchLeguleguA500PERows(client *http.Client, now time.Time) ([]leguleguIndexPERow, error) {
	pageURL := "https://legulegu.com/stockdata/index-ttm-lyr-pe?indexCode=000510.CSI"
	pageReq, err := http.NewRequest(http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, err
	}
	pageReq.Header.Set("Accept", "text/html,*/*")
	pageReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	pageReq.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	pageResp, err := client.Do(pageReq)
	if err != nil {
		return nil, err
	}
	cookieHeader := responseCookieHeader(pageResp)
	io.Copy(io.Discard, io.LimitReader(pageResp.Body, 1<<20))
	pageResp.Body.Close()
	if pageResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("legulegu page request failed: %s", pageResp.Status)
	}
	if strings.TrimSpace(cookieHeader) == "" {
		return nil, errors.New("missing legulegu session cookies")
	}

	var lastErr error
	for _, tokenDate := range []time.Time{now, now.AddDate(0, 0, -1), now.AddDate(0, 0, -2)} {
		rows, err := fetchLeguleguA500PERowsWithToken(client, cookieHeader, tokenDate)
		if err == nil && len(rows) > 0 {
			return rows, nil
		}
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("empty legulegu A500 PE response")
}

func fetchLeguleguA500PERowsWithToken(client *http.Client, cookieHeader string, tokenDate time.Time) ([]leguleguIndexPERow, error) {
	values := url.Values{}
	values.Set("indexCode", "000510.CSI")
	values.Set("token", leguleguToken(tokenDate))
	endpoint := "https://legulegu.com/api/stockdata/index-basic-pe?" + values.Encode()
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", cookieHeader)
	req.Header.Set("Referer", "https://legulegu.com/stockdata/index-ttm-lyr-pe?indexCode=000510.CSI")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("legulegu A500 PE request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return nil, errors.New("empty legulegu A500 PE response")
	}
	return parseLeguleguIndexPERows(body)
}

func responseCookieHeader(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	parts := []string{}
	for _, cookie := range resp.Cookies() {
		if strings.TrimSpace(cookie.Name) != "" {
			parts = append(parts, cookie.Name+"="+cookie.Value)
		}
	}
	return strings.Join(parts, "; ")
}

func leguleguToken(date time.Time) string {
	sum := md5.Sum([]byte(date.Format("2006-01-02")))
	return fmt.Sprintf("%x", sum)
}

func parseLeguleguIndexPERows(body []byte) ([]leguleguIndexPERow, error) {
	var payload struct {
		Data []struct {
			Date          string   `json:"date"`
			TtmPE         float64  `json:"ttmPe"`
			TtmPEQuantile *float64 `json:"ttmPeQuantile"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	rows := make([]leguleguIndexPERow, 0, len(payload.Data))
	for _, row := range payload.Data {
		if _, err := time.Parse("2006-01-02", row.Date); err != nil {
			continue
		}
		if row.TtmPE <= 0 && !validPercentilePointer(row.TtmPEQuantile) {
			continue
		}
		rows = append(rows, leguleguIndexPERow{Date: row.Date, TtmPE: row.TtmPE, TtmPEQuantile: row.TtmPEQuantile})
	}
	if len(rows) == 0 {
		return nil, errors.New("missing legulegu A500 PE rows")
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Date < rows[j].Date })
	return rows, nil
}

func a500PEPercentileFromRows(rows []leguleguIndexPERow) (float64, string, error) {
	if len(rows) == 0 {
		return 0, "", errors.New("missing A500 PE rows")
	}
	latest := rows[len(rows)-1]
	latestDate, err := time.Parse("2006-01-02", latest.Date)
	if err != nil {
		return 0, "", err
	}
	if validPercentilePointer(latest.TtmPEQuantile) {
		return *latest.TtmPEQuantile, latest.Date, nil
	}
	cutoff := latestDate.AddDate(-5, 0, 0)
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		rowDate, err := time.Parse("2006-01-02", row.Date)
		if err != nil || rowDate.Before(cutoff) || row.TtmPE <= 0 {
			continue
		}
		values = append(values, row.TtmPE)
	}
	if len(values) == 0 || latest.TtmPE <= 0 {
		return 0, "", errors.New("missing five-year A500 PE values")
	}
	return percentileRank(latest.TtmPE, values), latest.Date, nil
}

func validPercentilePointer(value *float64) bool {
	return value != nil && !math.IsNaN(*value) && !math.IsInf(*value, 0) && *value >= 0 && *value <= 1
}

func fetchSP500CAPEPercentile(client *http.Client) (float64, string, error) {
	percentile, date, err := fetchHistoryOfMarketSP500CAPEPercentile(client)
	if err == nil && !valuationDateStale(date, time.Now(), primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	historyErr := err
	if historyErr == nil {
		historyErr = fmt.Errorf("historyofmarket CAPE stale as of %s", date)
	}
	observations, err := fetchMultplMonthlyValues(client, multplShillerCAPEURL)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("historyofmarket: %v; multpl: %v", historyErr, err)
	}
	fallbackPercentile, fallbackDate, err := capePercentileFromMonthlyValues(observations, 10)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("historyofmarket: %v; multpl: %v", historyErr, err)
	}
	if date != "" && !quoteDateAfter(fallbackDate, date) {
		return percentile, date, nil
	}
	return fallbackPercentile, fallbackDate, nil
}

func fetchNasdaq100ForwardPEPercentile(client *http.Client) (float64, string, error) {
	percentile, date, err := fetchHistoryOfMarketNasdaq100ForwardPEPercentile(client)
	if err == nil && !valuationDateStale(date, time.Now(), primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	historyErr := err
	if historyErr == nil {
		historyErr = fmt.Errorf("historyofmarket Nasdaq 100 forward PE stale as of %s", date)
	}
	snapshot, err := fetchWorldPERatioNasdaq100(client, worldPERatioNasdaq100URL)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("historyofmarket: %v; worldperatio: %v", historyErr, err)
	}
	return snapshot.Percentile, snapshot.Date, nil
}

type historyOfMarketPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type historyOfMarketCurrentValuation struct {
	Forward float64 `json:"forward"`
}

func fetchHistoryOfMarketSP500CAPEPercentile(client *http.Client) (float64, string, error) {
	var payload struct {
		CAPE []historyOfMarketPoint `json:"cape"`
	}
	if err := fetchHistoryOfMarketJSON(client, historyOfMarketSP500PEURL, &payload); err != nil {
		return 0, "", err
	}
	return percentileFromHistoryOfMarketPoints(payload.CAPE, 10, "CAPE")
}

func fetchHistoryOfMarketNasdaq100ForwardPEPercentile(client *http.Client) (float64, string, error) {
	var payload struct {
		Updated string                          `json:"updated"`
		Current historyOfMarketCurrentValuation `json:"current"`
		Forward []historyOfMarketPoint          `json:"forward"`
	}
	if err := fetchHistoryOfMarketJSON(client, historyOfMarketNDXForwardPEURL, &payload); err != nil {
		return 0, "", err
	}
	points := historyOfMarketPointsWithCurrentForward(payload.Forward, payload.Updated, payload.Current)
	return percentileFromHistoryOfMarketPoints(points, 10, "Nasdaq 100 forward PE")
}

func historyOfMarketPointsWithCurrentForward(points []historyOfMarketPoint, updated string, current historyOfMarketCurrentValuation) []historyOfMarketPoint {
	combined := append([]historyOfMarketPoint(nil), points...)
	if current.Forward > 0 && strings.TrimSpace(updated) != "" {
		combined = append(combined, historyOfMarketPoint{
			Date:  strings.TrimSpace(updated),
			Value: current.Forward,
		})
	}
	return combined
}

func fetchHistoryOfMarketJSON(client *http.Client, endpoint string, target any) error {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "holds-website etf rule updater")
	requestClient := client
	if requestClient == nil {
		requestClient = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := requestClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("historyofmarket request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(target)
}

func percentileFromHistoryOfMarketPoints(points []historyOfMarketPoint, years int, label string) (float64, string, error) {
	observations := make([]dailyClose, 0, len(points))
	for _, point := range points {
		if _, err := time.Parse(etfRuleRuntimeTimestampDateLayout, point.Date); err != nil {
			continue
		}
		if point.Value <= 0 {
			continue
		}
		observations = append(observations, dailyClose{Date: point.Date, Price: point.Value})
	}
	return percentileFromDatedValues(observations, years, label)
}

func capePercentileFromMonthlyValues(observations []dailyClose, years int) (float64, string, error) {
	return percentileFromDatedValues(observations, years, "CAPE")
}

func valuationDateStale(date string, now time.Time, maxLagDays int) bool {
	parsed, err := time.Parse(etfRuleRuntimeTimestampDateLayout, strings.TrimSpace(date))
	if err != nil {
		return true
	}
	reference := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, parsed.Location())
	if reference.Before(parsed) {
		return false
	}
	return reference.Sub(parsed) > time.Duration(maxLagDays)*24*time.Hour
}

func percentileFromDatedValues(observations []dailyClose, years int, label string) (float64, string, error) {
	if len(observations) < 5 {
		return 0, "", fmt.Errorf("not enough %s observations", label)
	}
	ordered := append([]dailyClose(nil), observations...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Date < ordered[j].Date })
	latest := dailyClose{}
	for i := len(ordered) - 1; i >= 0; i-- {
		if ordered[i].Price > 0 && strings.TrimSpace(ordered[i].Date) != "" {
			latest = ordered[i]
			break
		}
	}
	if latest.Price <= 0 {
		return 0, "", fmt.Errorf("missing %s values", label)
	}
	latestDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, latest.Date)
	if err != nil {
		return 0, "", err
	}
	cutoff := latestDate.AddDate(-years, 0, 0)
	values := make([]float64, 0, len(ordered))
	for _, observation := range ordered {
		observationDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, observation.Date)
		if err != nil || observationDate.Before(cutoff) || observation.Price <= 0 {
			continue
		}
		values = append(values, observation.Price)
	}
	if len(values) == 0 {
		return 0, "", fmt.Errorf("missing ten-year %s values", label)
	}
	return percentileRank(latest.Price, values), latest.Date, nil
}

func fetchMultplMonthlyValues(client *http.Client, endpoint string) ([]dailyClose, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "holds-website etf rule updater")
	req.Header.Set("Referer", "https://www.multpl.com/")

	requestClient := client
	if requestClient == nil || (requestClient.Timeout > 0 && requestClient.Timeout < 30*time.Second) {
		requestClient = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := requestClient.Do(req)
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		req, retryReqErr := http.NewRequest(http.MethodGet, endpoint, nil)
		if retryReqErr != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "holds-website etf rule updater")
		req.Header.Set("Referer", "https://www.multpl.com/")
		resp, err = requestClient.Do(req)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("valuation request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	return parseMultplTable(body)
}

func parseMultplTable(body []byte) ([]dailyClose, error) {
	text := string(body)
	rowPattern := regexp.MustCompile(`(?is)<tr[^>]*>.*?</tr>`)
	cellPattern := regexp.MustCompile("(?is)<td[^>]*>(.*?)</td>")
	rows := rowPattern.FindAllString(text, -1)
	values := make([]dailyClose, 0, len(rows))
	for _, row := range rows {
		cells := cellPattern.FindAllStringSubmatch(row, -1)
		if len(cells) < 2 {
			continue
		}
		dateText := htmlPlainText(cells[0][1])
		valueText := htmlPlainText(cells[1][1])
		value, err := firstTextNumber(valueText)
		if err != nil || value <= 0 {
			continue
		}
		date := normalizeMultplDate(dateText)
		if date == "" {
			continue
		}
		values = append(values, dailyClose{Date: date, Price: value})
	}
	if len(values) == 0 {
		return nil, errors.New("missing valuation rows")
	}
	return values, nil
}

type cashDividendEvent struct {
	Date   string
	Amount float64
}

func fetchDividendLowVolYield(client *http.Client) (float64, string, error) {
	closes, err := fetchRuleDailyCloses(client, "515450.SH", 30)
	if err != nil {
		return 0, "", err
	}
	if len(closes) == 0 || closes[len(closes)-1].Price <= 0 {
		return 0, "", errors.New("missing 515450 close price")
	}
	latest := closes[len(closes)-1]
	events, err := fetchEastmoneyFundDividends(client, "515450")
	if err != nil {
		return 0, "", err
	}
	trailingAmount, err := trailingFundDividendAmount(events, latest.Date)
	if err != nil {
		return 0, "", err
	}
	return trailingAmount / latest.Price, latest.Date, nil
}

func fetchEastmoneyFundDividends(client *http.Client, code string) ([]cashDividendEvent, error) {
	endpoint := "https://fundf10.eastmoney.com/fhsp_" + url.PathEscape(normalizeFundSymbol(code)) + ".html"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,*/*")
	req.Header.Set("Referer", "https://fundf10.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("fund dividend request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	return parseEastmoneyFundDividends(body)
}

func parseEastmoneyFundDividends(body []byte) ([]cashDividendEvent, error) {
	rowPattern := regexp.MustCompile(`(?is)<tr[^>]*>.*?</tr>`)
	cellPattern := regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)
	rows := rowPattern.FindAllString(string(body), -1)
	events := make([]cashDividendEvent, 0, len(rows))
	for _, row := range rows {
		cells := cellPattern.FindAllStringSubmatch(row, -1)
		if len(cells) < 4 {
			continue
		}
		exDate := htmlPlainText(cells[2][1])
		if _, err := time.Parse("2006-01-02", exDate); err != nil {
			continue
		}
		amountText := htmlPlainText(cells[3][1])
		amount, err := fundDividendAmount(amountText)
		if err != nil || amount <= 0 {
			continue
		}
		events = append(events, cashDividendEvent{Date: exDate, Amount: amount})
	}
	if len(events) == 0 {
		return nil, errors.New("missing fund dividend rows")
	}
	return events, nil
}

func fundDividendAmount(value string) (float64, error) {
	amount, err := firstTextNumber(value)
	if err != nil {
		return 0, errors.New("dividend amount not found")
	}
	return amount, nil
}

func trailingFundDividendAmount(events []cashDividendEvent, referenceDate string) (float64, error) {
	reference, err := time.Parse("2006-01-02", referenceDate)
	if err != nil {
		return 0, err
	}
	cutoff := reference.AddDate(-1, 0, 0)
	total := 0.0
	for _, event := range events {
		eventDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil || event.Amount <= 0 {
			continue
		}
		if eventDate.Before(cutoff) || eventDate.After(reference) {
			continue
		}
		total += event.Amount
	}
	if total <= 0 {
		return 0, errors.New("missing trailing fund dividend")
	}
	return total, nil
}

func normalizeMultplDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"Jan 2, 2006", "January 2, 2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format("2006-01-02")
		}
	}
	return ""
}

func percentileRank(value float64, values []float64) float64 {
	if len(values) == 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	sortedValues := append([]float64(nil), values...)
	sort.Float64s(sortedValues)
	count := 0
	for _, item := range sortedValues {
		if item <= value {
			count++
		}
	}
	return float64(count) / float64(len(sortedValues))
}

const worldPERatioNasdaq100URL = "https://worldperatio.com/index/nasdaq-100/"

type worldPERatioSnapshot struct {
	CurrentPE  float64
	Average10Y float64
	StdDev10Y  float64
	ZScore     float64
	Percentile float64
	Date       string
}

func fetchWorldPERatioNasdaq100(client *http.Client, endpoint string) (worldPERatioSnapshot, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return worldPERatioSnapshot{}, err
	}
	req.Header.Set("User-Agent", "holds-website etf rule updater")
	req.Header.Set("Referer", "https://worldperatio.com/")

	resp, err := client.Do(req)
	if err != nil {
		return worldPERatioSnapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return worldPERatioSnapshot{}, fmt.Errorf("worldperatio request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return worldPERatioSnapshot{}, err
	}
	return parseWorldPERatioNasdaq100(body)
}

func parseWorldPERatioNasdaq100(body []byte) (worldPERatioSnapshot, error) {
	pageText := htmlPlainText(string(body))
	currentPE, date, err := parseWorldPERatioCurrentPE(pageText)
	if err != nil {
		return worldPERatioSnapshot{}, err
	}

	row := worldPERatioPeriodRow(string(body), "Last 10Y")
	if strings.TrimSpace(row) == "" {
		return worldPERatioSnapshot{}, errors.New("missing Last 10Y row")
	}
	cellPattern := regexp.MustCompile("(?is)<td[^>]*>(.*?)</td>")
	cellMatches := cellPattern.FindAllStringSubmatch(row, -1)
	if len(cellMatches) < 6 {
		return worldPERatioSnapshot{}, errors.New("incomplete Last 10Y row")
	}
	cells := make([]string, 0, len(cellMatches))
	for _, match := range cellMatches {
		cells = append(cells, htmlPlainText(match[1]))
	}

	average, err := firstTextNumber(cells[1])
	if err != nil {
		return worldPERatioSnapshot{}, fmt.Errorf("missing Last 10Y average: %w", err)
	}
	stdDev, err := firstTextNumber(cells[2])
	if err != nil {
		return worldPERatioSnapshot{}, fmt.Errorf("missing Last 10Y standard deviation: %w", err)
	}
	if stdDev <= 0 {
		return worldPERatioSnapshot{}, errors.New("invalid Last 10Y standard deviation")
	}
	zScore, err := firstSigmaValue(cells[5])
	if err != nil {
		zScore = (currentPE - average) / stdDev
	}
	return worldPERatioSnapshot{
		CurrentPE:  currentPE,
		Average10Y: average,
		StdDev10Y:  stdDev,
		ZScore:     zScore,
		Percentile: normalPercentileFromZ(zScore),
		Date:       date,
	}, nil
}

func worldPERatioPeriodRow(pageHTML string, period string) string {
	rowPattern := regexp.MustCompile("(?is)<tr[^>]*>.*?</tr>")
	periodPattern := regexp.MustCompile("(?i)\\b" + regexp.QuoteMeta(period) + "\\b")
	for _, row := range rowPattern.FindAllString(pageHTML, -1) {
		if periodPattern.MatchString(htmlPlainText(row)) {
			return row
		}
	}
	return ""
}

func parseWorldPERatioCurrentPE(pageText string) (float64, string, error) {
	pattern := regexp.MustCompile("(?i)Price-to-Earnings\\s*\\(P/E\\)\\s+Ratio\\s+for\\s+Nasdaq\\s+100\\s+Index\\s+is\\s+([0-9]+(?:\\.[0-9]+)?)\\s*,\\s+calculated\\s+on\\s+([0-9]{1,2}\\s+[A-Za-z]+\\s+[0-9]{4})")
	match := pattern.FindStringSubmatch(pageText)
	if len(match) < 3 {
		return 0, "", errors.New("missing current P/E paragraph")
	}
	currentPE, err := strconv.ParseFloat(match[1], 64)
	if err != nil || currentPE <= 0 {
		return 0, "", errors.New("invalid current P/E")
	}
	date := normalizeWorldPERatioDate(match[2])
	if date == "" {
		return 0, "", fmt.Errorf("unsupported current P/E date: %s", match[2])
	}
	return currentPE, date, nil
}

func normalizeWorldPERatioDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"02 January 2006", "2 January 2006", "02 Jan 2006", "2 Jan 2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format("2006-01-02")
		}
	}
	return ""
}

func htmlPlainText(value string) string {
	withoutScripts := regexp.MustCompile("(?is)<script[^>]*>.*?</script>|<style[^>]*>.*?</style>").ReplaceAllString(value, " ")
	withoutTags := regexp.MustCompile("(?is)<[^>]+>").ReplaceAllString(withoutScripts, " ")
	decoded := html.UnescapeString(withoutTags)
	return strings.Join(strings.Fields(decoded), " ")
}

func firstTextNumber(value string) (float64, error) {
	pattern := regexp.MustCompile("[-+]?[0-9]+(?:\\.[0-9]+)?")
	match := pattern.FindString(value)
	if match == "" {
		return 0, errors.New("number not found")
	}
	return strconv.ParseFloat(match, 64)
}

func firstSigmaValue(value string) (float64, error) {
	pattern := regexp.MustCompile("([-+]?[0-9]+(?:\\.[0-9]+)?)\\s*(?:σ|sigma)")
	match := pattern.FindStringSubmatch(value)
	if len(match) < 2 {
		return 0, errors.New("sigma value not found")
	}
	return strconv.ParseFloat(match[1], 64)
}

func normalPercentileFromZ(zScore float64) float64 {
	if math.IsNaN(zScore) || math.IsInf(zScore, 0) {
		return 0
	}
	percentile := 0.5 * (1 + math.Erf(zScore/math.Sqrt2))
	if percentile < 0 {
		return 0
	}
	if percentile > 1 {
		return 1
	}
	return percentile
}

func evaluateA500Rule(inputs etfRuleInputs) etfRuleEvaluation {
	drawdown := valueOrNaN(inputs.Drawdown)
	valuation := valueOrNaN(inputs.ValuationPercentile)
	if !known(valuation) {
		return pendingRule("需要中证A500 PE分位决定基础倍数")
	}
	base := percentileBaseLevel(valuation, 0.80, 0.60, 0.40, 0.20)
	if !known(drawdown) {
		return partialRule(base, "已按PE分位得到基础倍数；回撤数据缺失，暂未做限速调整")
	}
	switch {
	case valuation > 0.80 && drawdown < 0.05:
		return completeRule(downshiftLevel(base), "PE分位>80%且回撤<5%，高位限速")
	case valuation >= 0.20 && valuation < 0.40 && drawdown < 0.12:
		return completeRule("one", "PE分位20%—40%但回撤<12%，低估确认不足")
	case valuation < 0.20 && drawdown < 0.18:
		return completeRule("oneHalf", "PE分位<20%但回撤<18%，极低确认不足")
	default:
		return completeRule(base, "按PE分位基础倍数执行，回撤未触发限速")
	}
}

func evaluateSP500Rule(inputs etfRuleInputs) etfRuleEvaluation {
	drawdown := valueOrNaN(inputs.Drawdown)
	pePercentile := valueOrNaN(inputs.ValuationPercentile)
	if !known(pePercentile) {
		return pendingRule("需要标普500 PE分位决定基础倍数")
	}
	base := percentileBaseLevel(pePercentile, 0.95, 0.80, 0.40, 0.20)
	if !known(drawdown) {
		return partialRule(base, "已按PE分位得到基础倍数；回撤数据缺失，暂未做限速调整")
	}
	switch {
	case pePercentile > 0.80 && drawdown < 0.05:
		return completeRule(downshiftLevel(base), "PE分位>80%且回撤<5%，高估限速")
	case pePercentile >= 0.40 && pePercentile <= 0.80 && drawdown < 0.05:
		return completeRule("half", "PE分位40%—80%但回撤<5%，正常估值高位限速")
	case pePercentile >= 0.20 && pePercentile < 0.40 && drawdown < 0.15:
		return completeRule("one", "PE分位20%—40%但回撤<15%，低估确认不足")
	case pePercentile < 0.20 && drawdown < 0.20:
		return completeRule("oneHalf", "PE分位<20%但回撤<20%，极低确认不足")
	default:
		return completeRule(base, "按PE分位基础倍数执行，回撤未触发限速")
	}
}

func evaluateDividendLowVolRule(inputs etfRuleInputs) etfRuleEvaluation {
	drawdown := valueOrNaN(inputs.Drawdown)
	yield := valueOrNaN(inputs.ValuationPercentile)
	if !known(yield) {
		return pendingRule("需要515450股息率决定基础倍数")
	}
	base := dividendYieldBaseLevel(yield)
	if !known(drawdown) {
		return partialRule(base, "已按股息率得到基础倍数；回撤数据缺失，暂未做限速调整")
	}
	switch {
	case yield < 0.050 && drawdown < 0.05:
		return completeRule(downshiftLevel(base), "股息率<5.0%且回撤<5%，股息率偏低且价格高位")
	case yield >= 0.058 && yield <= 0.062 && drawdown < 0.08:
		return completeRule("one", "股息率5.8%—6.2%但回撤<8%，低位确认不足")
	case yield > 0.062 && drawdown < 0.12:
		return completeRule("oneHalf", "股息率>6.2%但回撤<12%，极低确认不足")
	default:
		return completeRule(base, "按股息率基础倍数执行，回撤未触发限速")
	}
}

func evaluateNasdaq100Rule(inputs etfRuleInputs) etfRuleEvaluation {
	drawdown := valueOrNaN(inputs.Drawdown)
	pePercentile := valueOrNaN(inputs.ValuationPercentile)
	if !known(pePercentile) {
		return pendingRule("需要纳指100 PE分位决定基础倍数")
	}
	base := percentileBaseLevel(pePercentile, 0.85, 0.70, 0.40, 0.20)
	if !known(drawdown) {
		return partialRule(base, "已按PE分位得到基础倍数；回撤数据缺失，暂未做限速调整")
	}
	switch {
	case pePercentile > 0.70 && drawdown < 0.05:
		return completeRule(downshiftLevel(base), "PE分位>70%且回撤<5%，高估限速")
	case pePercentile >= 0.40 && pePercentile <= 0.70 && drawdown < 0.05:
		return completeRule("half", "PE分位40%—70%但回撤<5%，正常估值高位限速")
	case pePercentile >= 0.20 && pePercentile < 0.40 && drawdown < 0.20:
		return completeRule("one", "PE分位20%—40%但回撤<20%，低估确认不足")
	case pePercentile < 0.20 && drawdown < 0.30:
		return completeRule("oneHalf", "PE分位<20%但回撤<30%，极低确认不足")
	default:
		return completeRule(base, "按PE分位基础倍数执行，回撤未触发限速")
	}
}

func percentileBaseLevel(value float64, quarterThreshold float64, halfThreshold float64, oneThreshold float64, oneHalfThreshold float64) string {
	switch {
	case value > quarterThreshold:
		return "quarter"
	case value >= halfThreshold:
		return "half"
	case value >= oneThreshold:
		return "one"
	case value >= oneHalfThreshold:
		return "oneHalf"
	default:
		return "two"
	}
}

func dividendYieldBaseLevel(yield float64) string {
	switch {
	case yield < 0.047:
		return "quarter"
	case yield <= 0.050:
		return "half"
	case yield <= 0.058:
		return "one"
	case yield <= 0.062:
		return "oneHalf"
	default:
		return "two"
	}
}

func zScoreBaseLevel(zScore float64) string {
	switch {
	case zScore > 2:
		return "quarter"
	case zScore >= 1:
		return "half"
	case zScore >= -1:
		return "one"
	case zScore >= -2:
		return "oneHalf"
	default:
		return "two"
	}
}

func downshiftLevel(level string) string {
	switch level {
	case "two":
		return "oneHalf"
	case "oneHalf":
		return "one"
	case "one":
		return "half"
	case "half":
		return "quarter"
	default:
		return "quarter"
	}
}

func completeRule(level string, reason string) etfRuleEvaluation {
	return etfRuleEvaluation{Level: level, Complete: true, Reason: reason}
}

func partialRule(level string, reason string) etfRuleEvaluation {
	return etfRuleEvaluation{Level: level, Complete: false, Reason: reason}
}

func pendingRule(reason string) etfRuleEvaluation {
	return etfRuleEvaluation{Complete: false, Reason: reason}
}

func valueOrNaN(value *float64) float64 {
	if value == nil {
		return math.NaN()
	}
	return *value
}

func known(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func percentMetric(value float64) *float64 {
	percent := value * 100
	return &percent
}

func floatMetric(value float64) *float64 {
	return &value
}

func configValuationMetricUnit(config etfRuleConfig) string {
	if config.ValuationMetricKey == "peZScore" {
		return "σ"
	}
	return "%"
}

func runtimeETFRuleStatusList(records map[string]ETFRuleStatus) []ETFRuleStatus {
	keys := make([]string, 0, len(records))
	for key := range records {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	list := make([]ETFRuleStatus, 0, len(keys))
	for _, key := range keys {
		list = append(list, records[key])
	}
	return list
}

func mergeETFRuleStatusWithExisting(next ETFRuleStatus, existing ETFRuleStatus, now time.Time) ETFRuleStatus {
	if strings.TrimSpace(existing.Symbol) == "" {
		if config, ok := etfRuleConfigBySymbol(next.Symbol); ok {
			return enforceETFRuleStatusConfidence(next, config, now)
		}
		return next
	}
	existingMetrics := map[string]ETFRuleMetric{}
	for _, metric := range existing.Metrics {
		if strings.TrimSpace(metric.Key) != "" {
			existingMetrics[metric.Key] = metric
		}
	}
	usedFallback := false
	for i := range next.Metrics {
		if next.Metrics[i].Available {
			continue
		}
		previous, ok := existingMetrics[next.Metrics[i].Key]
		if !ok || !previous.Available {
			continue
		}
		previous.Error = ""
		next.Metrics[i] = previous
		usedFallback = true
	}
	if usedFallback {
		next = refreshETFRuleStatusFromMetrics(next)
	}
	if config, ok := etfRuleConfigBySymbol(next.Symbol); ok {
		next = enforceETFRuleStatusConfidence(next, config, now)
	}
	return next
}

func refreshETFRuleStatusFromMetrics(status ETFRuleStatus) ETFRuleStatus {
	config, ok := etfRuleConfigBySymbol(status.Symbol)
	if !ok {
		return status
	}
	inputs := etfRuleInputs{}
	for _, metric := range status.Metrics {
		if !metric.Available || metric.Value == nil {
			continue
		}
		switch metric.Key {
		case "drawdown252":
			value := *metric.Value / 100
			inputs.Drawdown = &value
			inputs.DrawdownAsOf = metric.AsOf
		case config.ValuationMetricKey:
			if metric.Key == "peZScore" || strings.TrimSpace(metric.Unit) == "σ" {
				value := *metric.Value
				inputs.ValuationZScore = &value
			} else {
				value := *metric.Value / 100
				inputs.ValuationPercentile = &value
			}
			inputs.ValuationAsOf = metric.AsOf
		}
	}
	evaluation := config.Evaluate(inputs)
	level := config.Levels[evaluation.Level]
	status.Level = evaluation.Level
	status.LevelLabel = level.Label
	status.MonthlyAmount = config.Monthly[evaluation.Level]
	status.WeeklyAmount = config.Weekly[evaluation.Level]
	status.Complete = evaluation.Complete
	status.Reason = evaluation.Reason
	status.AsOf = firstNonEmpty(inputs.DrawdownAsOf, inputs.ValuationAsOf, status.AsOf)
	if status.Level == "" {
		status.LevelLabel = "待数据"
	}
	return status
}

func etfRuleConfigBySymbol(symbol string) (etfRuleConfig, bool) {
	normalized := normalizeFundSymbol(symbol)
	for _, config := range etfRuleConfigs {
		if normalizeFundSymbol(config.Symbol) == normalized {
			return config, true
		}
	}
	return etfRuleConfig{}, false
}

func enforceETFRuleStatusConfidence(status ETFRuleStatus, config etfRuleConfig, now time.Time) ETFRuleStatus {
	if !status.Complete {
		return zeroETFRuleExecutableAmount(status)
	}
	if len(etfRuleStatusConfidenceIssues(status, config, now)) == 0 {
		return status
	}
	status = zeroETFRuleExecutableAmount(status)
	if strings.TrimSpace(status.Reason) == "" {
		status.Reason = "等待指标刷新"
	}
	return status
}

func zeroETFRuleExecutableAmount(status ETFRuleStatus) ETFRuleStatus {
	status.Complete = false
	status.MonthlyAmount = 0
	status.WeeklyAmount = 0
	return status
}

func etfRuleStatusConfidenceIssues(status ETFRuleStatus, config etfRuleConfig, now time.Time) []string {
	issues := []string{}
	metricsByKey := map[string]ETFRuleMetric{}
	for _, metric := range status.Metrics {
		if strings.TrimSpace(metric.Key) == "" {
			continue
		}
		metricsByKey[metric.Key] = metric
	}
	for _, key := range []string{"drawdown252", config.ValuationMetricKey} {
		metric, ok := metricsByKey[key]
		if !ok {
			issues = append(issues, key+"缺失")
			continue
		}
		issues = append(issues, etfRuleMetricConfidenceIssues(metric, config, now)...)
	}
	if len(status.Sources) < 2 {
		issues = append(issues, "数据源不足")
	}
	return issues
}

func etfRuleMetricConfidenceIssues(metric ETFRuleMetric, config etfRuleConfig, now time.Time) []string {
	issues := []string{}
	if !metric.Available || metric.Value == nil {
		return append(issues, firstNonEmpty(metric.Label, metric.Key)+"不可用")
	}
	value := *metric.Value
	if math.IsNaN(value) || math.IsInf(value, 0) || !etfRuleMetricValueInExpectedRange(metric, config, value) {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"数值异常")
	}
	if strings.TrimSpace(metric.Error) != "" {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"沿用旧值")
	}
	metricDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, strings.TrimSpace(metric.AsOf))
	if err != nil {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"日期缺失")
		return issues
	}
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, metricDate.Location())
	if metricDate.After(nowDate.Add(24 * time.Hour)) {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"日期异常")
		return issues
	}
	maxAgeDays := etfRuleMetricMaxAgeDays(metric, config)
	if nowDate.Sub(metricDate).Hours()/24 > float64(maxAgeDays) {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"过期")
	}
	return issues
}

func etfRuleMetricValueInExpectedRange(metric ETFRuleMetric, config etfRuleConfig, value float64) bool {
	switch metric.Key {
	case "drawdown252":
		return value >= 0 && value <= 100
	case config.ValuationMetricKey:
		if metric.Key == "peZScore" || strings.TrimSpace(metric.Unit) == "σ" {
			return value >= -6 && value <= 6
		}
		if metric.Key == "dividendYield" {
			return value > 0 && value <= 20
		}
		return value >= 0 && value <= 100
	default:
		return true
	}
}

func etfRuleMetricMaxAgeDays(metric ETFRuleMetric, config etfRuleConfig) int {
	if metric.Key == config.ValuationMetricKey && (config.Symbol == "018738" || config.Symbol == "021000") {
		return etfRuleMonthlyMetricMaxAgeDays
	}
	return etfRuleDailyMetricMaxAgeDays
}
