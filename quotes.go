package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type QuoteUpdateResponse struct {
	Updated int         `json:"updated"`
	Skipped []QuoteSkip `json:"skipped"`
	State   AppState    `json:"state"`
}

type QuoteSkip struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Error  string `json:"error"`
}

type yahooChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency         string `json:"currency"`
				ExchangeTimezone string `json:"exchangeTimezoneName"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
			Events struct {
				Dividends map[string]struct {
					Amount float64 `json:"amount"`
					Date   int64   `json:"date"`
				} `json:"dividends"`
			} `json:"events"`
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}

type quote struct {
	Price              float64
	PreviousClose      float64
	PriceDate          string
	PreviousCloseDate  string
	Currency           string
	SourceSymbol       string
	SourceName         string
	DividendPerShare   *float64
	DividendCurrency   string
	DividendFiscalYear string
}

type dailyClose struct {
	Price float64
	Date  string
}

func (s *Server) handleUpdateQuotes(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := loadState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}

	now := time.Now()
	updated, skipped, quoteRecords := updateQuotes(&state, &http.Client{Timeout: 12 * time.Second}, now)
	if updated > 0 {
		if err := saveRuntimeQuoteRecords(quoteRecords, now.Format("2006-01-02 15:04:05")); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to save runtime quotes")
			return
		}
	}

	s.state = state
	writeJSON(w, http.StatusOK, QuoteUpdateResponse{
		Updated: updated,
		Skipped: skipped,
		State:   state,
	})
}

func updateQuotes(state *AppState, client *http.Client, now time.Time) (int, []QuoteSkip, []RuntimeQuote) {
	updated := 0
	skipped := []QuoteSkip{}
	quoteRecords := map[string]RuntimeQuote{}
	cache := make(map[string]quote)
	fallbackCache, fallbackErr := fetchFallbackQuotes(client, quoteSymbols(state))
	updateLabel := now.Format("2006-01-02 15:04:05")

	for i := range state.Holdings {
		holding := &state.Holdings[i]
		if strings.TrimSpace(holding.Symbol) == "" {
			continue
		}

		quote, err := fetchQuoteCached(client, cache, fallbackCache, fallbackErr, holding.Symbol)
		if err != nil {
			skipped = append(skipped, QuoteSkip{Type: "holding", Symbol: holding.Symbol, Name: holding.Name, Error: err.Error()})
			continue
		}

		applyHoldingQuote(holding, quote, updateLabel)
		quoteRecords[normalizeSymbol(holding.Symbol)] = runtimeQuoteFromQuote(holding.Symbol, quote, updateLabel)
		updated++
	}

	for i := range state.Candidates {
		candidate := &state.Candidates[i]
		if strings.TrimSpace(candidate.Symbol) == "" {
			continue
		}

		quote, err := fetchQuoteCached(client, cache, fallbackCache, fallbackErr, candidate.Symbol)
		if err != nil {
			skipped = append(skipped, QuoteSkip{Type: "candidate", Symbol: candidate.Symbol, Name: candidate.Name, Error: err.Error()})
			continue
		}

		applyCandidateQuote(candidate, quote, updateLabel)
		quoteRecords[normalizeSymbol(candidate.Symbol)] = runtimeQuoteFromQuote(candidate.Symbol, quote, updateLabel)
		updated++
	}

	return updated, skipped, runtimeQuoteList(quoteRecords)
}

func quoteSymbols(state *AppState) []string {
	symbols := []string{}
	seen := map[string]bool{}
	for _, holding := range state.Holdings {
		normalized := normalizeSymbol(holding.Symbol)
		if normalized != "" && !seen[normalized] {
			symbols = append(symbols, holding.Symbol)
			seen[normalized] = true
		}
	}
	for _, candidate := range state.Candidates {
		normalized := normalizeSymbol(candidate.Symbol)
		if normalized != "" && !seen[normalized] {
			symbols = append(symbols, candidate.Symbol)
			seen[normalized] = true
		}
	}
	return symbols
}

func fetchFallbackQuotes(client *http.Client, symbols []string) (map[string]quote, error) {
	quotes, err := fetchTencentQuotes(client, symbols)
	if err == nil {
		return quotes, nil
	}

	eastmoneyQuotes, eastmoneyErr := fetchEastmoneyQuotes(client, symbols)
	if eastmoneyErr == nil {
		return eastmoneyQuotes, nil
	}

	return nil, fmt.Errorf("tencent: %v; eastmoney: %v", err, eastmoneyErr)
}

func fetchQuoteCached(client *http.Client, cache map[string]quote, fallbackCache map[string]quote, fallbackErr error, symbol string) (quote, error) {
	normalized := normalizeSymbol(symbol)
	if cached, ok := cache[normalized]; ok {
		return cached, nil
	}

	quote, err := fetchQuote(client, normalized, fallbackCache, fallbackErr)
	if err != nil {
		return quote, err
	}
	cache[normalized] = quote
	return quote, nil
}

func applyHoldingQuote(holding *Holding, quote quote, updateLabel string) {
	holding.CurrentPrice = quote.Price
	holding.PreviousClose = quote.PreviousClose
	holding.CurrentPriceDate = quote.PriceDate
	holding.PreviousCloseDate = quote.PreviousCloseDate
	holding.MarginOfSafety = marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, holding.MarginOfSafety)
	applyDividendQuote(&holding.Dividend, quote, holding.Currency)
	if strings.TrimSpace(holding.Currency) == "" {
		holding.Currency = strings.ToUpper(strings.TrimSpace(quote.Currency))
	}
	holding.UpdatedAt = quoteUpdateLabel(updateLabel, quote)
}

func applyCandidateQuote(candidate *Candidate, quote quote, updateLabel string) {
	candidate.CurrentPrice = quote.Price
	candidate.PreviousClose = quote.PreviousClose
	candidate.CurrentPriceDate = quote.PriceDate
	candidate.PreviousCloseDate = quote.PreviousCloseDate
	candidate.MarginOfSafety = marginOfSafetyFromPrice(candidate.IntrinsicValue, candidate.CurrentPrice, candidate.MarginOfSafety)
	applyDividendQuote(&candidate.Dividend, quote, candidate.Currency)
	if strings.TrimSpace(candidate.Currency) == "" {
		candidate.Currency = strings.ToUpper(strings.TrimSpace(quote.Currency))
	}
	candidate.UpdatedAt = quoteUpdateLabel(updateLabel, quote)
}

func quoteUpdateLabel(updateLabel string, quote quote) string {
	sourceName := firstNonEmpty(quote.SourceName, "Yahoo Finance 日线收盘价")
	return fmt.Sprintf("%s；行情源 %s；代码 %s；币种 %s；日期 %s/%s", updateLabel, sourceName, quote.SourceSymbol, quote.Currency, quote.PreviousCloseDate, quote.PriceDate)
}

func quoteTriggerDecisionLog(state *AppState, symbol string, name string, currency string, beforePrice float64, currentPrice float64, intrinsicValue *float64, currentDate string, previousDate string, now time.Time) *DecisionLog {
	trigger := quoteTriggerText(beforePrice, currentPrice, intrinsicValue)
	if strings.TrimSpace(trigger) == "" {
		return nil
	}
	_, _, _, decision, discipline := decisionLogContext(state, symbol)
	return &DecisionLog{
		Date:       now.Format("2006-01-02 15:04:05"),
		Type:       "quote",
		Symbol:     symbol,
		Name:       name,
		Price:      pricePointer(currentPrice),
		Currency:   currency,
		Decision:   firstNonEmpty(decision, "行情触发"),
		Discipline: firstNonEmpty(discipline, "只在纪律区间变化时记录行情日志"),
		Detail:     fmt.Sprintf("%s；现价 %s %.4f；前值 %.4f；今收 %s；昨收 %s", trigger, strings.ToUpper(currency), currentPrice, beforePrice, firstNonEmpty(currentDate, "未知"), firstNonEmpty(previousDate, "未知")),
	}
}

func quoteTriggerText(beforePrice float64, currentPrice float64, intrinsicValue *float64) string {
	beforeZone := quotePriceZone(beforePrice, intrinsicValue)
	currentZone := quotePriceZone(currentPrice, intrinsicValue)
	if beforeZone != "" && currentZone != "" && beforeZone != currentZone {
		return fmt.Sprintf("进入/离开关键区间：%s -> %s", beforeZone, currentZone)
	}
	beforeMargin := quoteMarginZone(beforePrice, intrinsicValue)
	currentMargin := quoteMarginZone(currentPrice, intrinsicValue)
	if beforeMargin != "" && currentMargin != "" && beforeMargin != currentMargin && (currentMargin == "安全边际达标" || currentMargin == "高于内在价值") {
		return fmt.Sprintf("安全边际跨区：%s -> %s", beforeMargin, currentMargin)
	}
	return ""
}

func quotePriceZone(price float64, intrinsicValue *float64) string {
	if price <= 0 || intrinsicValue == nil || *intrinsicValue <= 0 {
		return ""
	}
	initialBuyPrice := *intrinsicValue * (1 - defaultSafetyMarginTarget)
	watchPrice := initialBuyPrice * 1.05
	aggressiveBuyPrice := initialBuyPrice * 0.9
	switch {
	case price <= aggressiveBuyPrice:
		return "重仓区"
	case price <= initialBuyPrice:
		return "首买区"
	case price <= watchPrice:
		return "观察区"
	default:
		return "等待区"
	}
}

func quoteMarginZone(price float64, intrinsicValue *float64) string {
	if price <= 0 || intrinsicValue == nil || *intrinsicValue <= 0 {
		return ""
	}
	margin := (*intrinsicValue - price) / *intrinsicValue
	switch {
	case margin >= defaultSafetyMarginTarget:
		return "安全边际达标"
	case margin < 0:
		return "高于内在价值"
	default:
		return "安全边际不足"
	}
}

func applyDividendQuote(current **Dividend, quote quote, fallbackCurrency string) {
	if quote.DividendPerShare == nil || *quote.DividendPerShare <= 0 {
		return
	}
	if *current == nil {
		*current = &Dividend{}
	}
	dividend := *current
	if shouldPreserveDividendQuote(dividend) {
		if strings.TrimSpace(dividend.DividendCurrency) == "" {
			dividend.DividendCurrency = strings.ToUpper(firstNonEmpty(quote.DividendCurrency, quote.Currency, fallbackCurrency))
		}
		dividend.DividendYield = nil
		dividend.EstimatedAnnualCash = nil
		return
	}
	dividend.FiscalYear = firstNonEmpty(quote.DividendFiscalYear, dividend.FiscalYear)
	dividend.DividendPerShare = quote.DividendPerShare
	dividend.DividendCurrency = strings.ToUpper(firstNonEmpty(quote.DividendCurrency, quote.Currency, fallbackCurrency))
	dividend.DividendYield = nil
	dividend.EstimatedAnnualCash = nil
}

func shouldPreserveDividendQuote(dividend *Dividend) bool {
	if dividend == nil {
		return false
	}
	if dividend.CashDividendTotal != nil && *dividend.CashDividendTotal > 0 {
		return true
	}
	fiscalYear := strings.ToUpper(strings.TrimSpace(dividend.FiscalYear))
	return dividend.DividendPerShare != nil &&
		*dividend.DividendPerShare > 0 &&
		fiscalYear != "" &&
		!strings.HasPrefix(fiscalYear, "TTM")
}

func fetchQuote(client *http.Client, symbol string, fallbackCache map[string]quote, fallbackErr error) (quote, error) {
	quote, yahooErr := fetchYahooQuote(client, symbol)
	if yahooErr == nil {
		return quote, nil
	}

	if quote, ok := fallbackCache[normalizeSymbol(symbol)]; ok {
		return quote, nil
	}
	if fallbackErr != nil {
		return quote, fmt.Errorf("yahoo: %v; fallback: %v", yahooErr, fallbackErr)
	}

	return quote, fmt.Errorf("yahoo: %v; fallback: no quote for %s", yahooErr, symbol)
}

func fetchYahooQuote(client *http.Client, symbol string) (quote, error) {
	sourceSymbol := yahooSymbol(symbol)
	endpoint := "https://query1.finance.yahoo.com/v8/finance/chart/" + url.PathEscape(sourceSymbol) + "?range=2y&interval=1d&events=div"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return quote{}, err
	}
	req.Header.Set("User-Agent", "holds-website quote updater")

	resp, err := client.Do(req)
	if err != nil {
		return quote{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return quote{}, fmt.Errorf("quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload yahooChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return quote{}, err
	}
	if payload.Chart.Error != nil || len(payload.Chart.Result) == 0 {
		return quote{}, errors.New("empty quote response")
	}

	result := payload.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return quote{}, errors.New("missing close series")
	}

	location := loadLocation(result.Meta.ExchangeTimezone)
	closes := result.Indicators.Quote[0].Close
	validCloses := make([]dailyClose, 0, len(closes))
	for i, closePrice := range closes {
		if closePrice > 0 {
			validCloses = append(validCloses, dailyClose{
				Price: closePrice,
				Date:  closeDate(result.Timestamp, i, location),
			})
		}
	}
	if len(validCloses) < 2 {
		return quote{}, errors.New("need at least two daily close prices")
	}

	priceClose := validCloses[len(validCloses)-1]
	previousClose := validCloses[len(validCloses)-2]
	dividendPerShare, dividendFiscalYear := trailingDividendFromEvents(result.Events.Dividends, priceClose.Date, location)
	return quote{
		Price:              priceClose.Price,
		PreviousClose:      previousClose.Price,
		PriceDate:          priceClose.Date,
		PreviousCloseDate:  previousClose.Date,
		Currency:           result.Meta.Currency,
		SourceSymbol:       sourceSymbol,
		SourceName:         "Yahoo Finance 日线收盘价",
		DividendPerShare:   dividendPerShare,
		DividendCurrency:   result.Meta.Currency,
		DividendFiscalYear: dividendFiscalYear,
	}, nil
}

func trailingDividendFromEvents(events map[string]struct {
	Amount float64 `json:"amount"`
	Date   int64   `json:"date"`
}, priceDate string, location *time.Location) (*float64, string) {
	if len(events) == 0 {
		return nil, ""
	}
	reference, err := time.ParseInLocation("2006-01-02", priceDate, location)
	if err != nil {
		reference = time.Now().In(location)
	}
	cutoff := reference.AddDate(-1, 0, 0)
	total := 0.0
	latest := time.Time{}
	for _, event := range events {
		if event.Amount <= 0 || event.Date <= 0 {
			continue
		}
		eventDate := time.Unix(event.Date, 0).In(location)
		if eventDate.Before(cutoff) || eventDate.After(reference.AddDate(0, 0, 1)) {
			continue
		}
		total += event.Amount
		if eventDate.After(latest) {
			latest = eventDate
		}
	}
	if total <= 0 {
		return nil, ""
	}
	labelDate := reference
	if !latest.IsZero() {
		labelDate = latest
	}
	return &total, "TTM " + labelDate.Format("2006-01-02")
}

func fetchEastmoneyQuote(client *http.Client, symbol string) (quote, error) {
	sourceSymbol, scale, err := eastmoneySymbol(symbol)
	if err != nil {
		return quote{}, err
	}

	endpoint := "https://push2.eastmoney.com/api/qt/stock/get?secid=" + url.QueryEscape(sourceSymbol) + "&fields=f43,f57,f58,f60,f86,f107"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return quote{}, err
	}
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Referer", "https://quote.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website quote updater)")

	resp, err := client.Do(req)
	if err != nil {
		return quote{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return quote{}, fmt.Errorf("quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		RC   int            `json:"rc"`
		Data map[string]any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return quote{}, err
	}
	if payload.RC != 0 || len(payload.Data) == 0 {
		return quote{}, errors.New("empty quote response")
	}

	rawPrice, err := numberField(payload.Data, "f43")
	if err != nil {
		return quote{}, err
	}
	rawPreviousClose, err := numberField(payload.Data, "f60")
	if err != nil {
		return quote{}, err
	}
	if rawPrice <= 0 || rawPreviousClose <= 0 {
		return quote{}, errors.New("missing price series")
	}

	priceDate := eastmoneyQuoteDate(payload.Data)
	return quote{
		Price:             rawPrice / scale,
		PreviousClose:     rawPreviousClose / scale,
		PriceDate:         priceDate,
		PreviousCloseDate: priceDate,
		Currency:          currencyForSymbol(symbol),
		SourceSymbol:      sourceSymbol,
		SourceName:        "东方财富实时行情",
	}, nil
}

func fetchTencentQuotes(client *http.Client, symbols []string) (map[string]quote, error) {
	querySymbols := []string{}
	normalizedByQuery := map[string]string{}
	for _, symbol := range symbols {
		querySymbol, normalized, err := tencentSymbol(symbol)
		if err != nil {
			continue
		}
		querySymbols = append(querySymbols, querySymbol)
		normalizedByQuery[querySymbol] = normalized
	}
	if len(querySymbols) == 0 {
		return map[string]quote{}, errors.New("no supported tencent symbols")
	}

	endpoint := "http://qt.gtimg.cn/q=" + strings.Join(querySymbols, ",")
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/plain,*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website quote updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	quotes := map[string]quote{}
	for _, line := range strings.Split(string(body), ";") {
		querySymbol, fields, ok := parseTencentLine(line)
		if !ok {
			continue
		}
		normalized := normalizedByQuery[querySymbol]
		if normalized == "" || len(fields) <= 30 {
			continue
		}
		price, err := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64)
		if err != nil || price <= 0 {
			continue
		}
		previousClose, err := strconv.ParseFloat(strings.TrimSpace(fields[4]), 64)
		if err != nil || previousClose <= 0 {
			continue
		}
		priceDate := tencentQuoteDate(fields[30])
		quotes[normalized] = quote{
			Price:             price,
			PreviousClose:     previousClose,
			PriceDate:         priceDate,
			PreviousCloseDate: priceDate,
			Currency:          currencyForSymbol(normalized),
			SourceSymbol:      querySymbol,
			SourceName:        "腾讯实时行情",
		}
	}
	if len(quotes) == 0 {
		return nil, errors.New("empty quote response")
	}
	return quotes, nil
}

func parseTencentLine(line string) (string, []string, bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "v_") {
		return "", nil, false
	}
	equalIndex := strings.Index(line, "=")
	if equalIndex <= 2 {
		return "", nil, false
	}
	querySymbol := strings.TrimSpace(strings.TrimPrefix(line[:equalIndex], "v_"))
	payload := strings.TrimSpace(line[equalIndex+1:])
	payload = strings.Trim(payload, "\"")
	if querySymbol == "" || payload == "" {
		return "", nil, false
	}
	return querySymbol, strings.Split(payload, "~"), true
}

func tencentQuoteDate(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= len("2006/01/02") && strings.Contains(value[:len("2006/01/02")], "/") {
		return strings.ReplaceAll(value[:len("2006/01/02")], "/", "-")
	}
	if len(value) >= len("20060102") {
		if parsed, err := time.Parse("20060102", value[:len("20060102")]); err == nil {
			return parsed.Format("2006-01-02")
		}
	}
	return time.Now().Format("2006-01-02")
}

func tencentSymbol(symbol string) (string, string, error) {
	symbol = normalizeSymbol(symbol)
	switch {
	case strings.HasSuffix(symbol, ".SH"):
		code := strings.TrimSuffix(symbol, ".SH")
		return "sh" + code, symbol, nil
	case strings.HasSuffix(symbol, ".SZ"):
		code := strings.TrimSuffix(symbol, ".SZ")
		return "sz" + code, symbol, nil
	case strings.HasSuffix(symbol, ".HK"):
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconv.Atoi(code); err == nil {
			code = fmt.Sprintf("%05d", value)
		}
		return "hk" + code, symbol, nil
	default:
		return "", "", fmt.Errorf("unsupported tencent symbol: %s", symbol)
	}
}

func fetchEastmoneyQuotes(client *http.Client, symbols []string) (map[string]quote, error) {
	secIDs := []string{}
	for _, symbol := range symbols {
		secID, _, err := eastmoneySymbol(symbol)
		if err == nil {
			secIDs = append(secIDs, secID)
		}
	}
	if len(secIDs) == 0 {
		return map[string]quote{}, errors.New("no supported eastmoney symbols")
	}

	endpoint := "https://push2.eastmoney.com/api/qt/ulist.np/get?secids=" + strings.Join(secIDs, ",") + "&fields=f2,f12,f13,f14,f18,f124"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Referer", "https://quote.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website quote updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		RC   int `json:"rc"`
		Data struct {
			Diff []map[string]any `json:"diff"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.RC != 0 || len(payload.Data.Diff) == 0 {
		return nil, errors.New("empty quote response")
	}

	quotes := make(map[string]quote, len(payload.Data.Diff))
	for _, item := range payload.Data.Diff {
		normalized, sourceSymbol, scale, err := eastmoneyDiffSymbol(item)
		if err != nil {
			continue
		}
		rawPrice, err := numberField(item, "f2")
		if err != nil {
			continue
		}
		rawPreviousClose, err := numberField(item, "f18")
		if err != nil {
			continue
		}
		if rawPrice <= 0 || rawPreviousClose <= 0 {
			continue
		}
		priceDate := eastmoneyQuoteDate(item)
		quotes[normalized] = quote{
			Price:             rawPrice / scale,
			PreviousClose:     rawPreviousClose / scale,
			PriceDate:         priceDate,
			PreviousCloseDate: priceDate,
			Currency:          currencyForSymbol(normalized),
			SourceSymbol:      sourceSymbol,
			SourceName:        "东方财富实时行情",
		}
	}
	if len(quotes) == 0 {
		return nil, errors.New("empty quote response")
	}
	return quotes, nil
}

func numberField(data map[string]any, key string) (float64, error) {
	value, ok := data[key]
	if !ok {
		return 0, fmt.Errorf("missing %s", key)
	}
	switch typed := value.(type) {
	case float64:
		return typed, nil
	case string:
		number, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid %s", key)
		}
		return number, nil
	default:
		return 0, fmt.Errorf("invalid %s", key)
	}
}

func eastmoneyQuoteDate(data map[string]any) string {
	for _, key := range []string{"f86", "f124"} {
		value, err := numberField(data, key)
		if err == nil && value > 0 {
			return time.Unix(int64(value), 0).In(loadLocation("Asia/Shanghai")).Format("2006-01-02")
		}
	}
	return time.Now().Format("2006-01-02")
}

func eastmoneyDiffSymbol(data map[string]any) (string, string, float64, error) {
	rawMarket, err := numberField(data, "f13")
	if err != nil {
		return "", "", 0, err
	}
	code, err := stringField(data, "f12")
	if err != nil {
		return "", "", 0, err
	}

	market := int(rawMarket)
	sourceSymbol := fmt.Sprintf("%d.%s", market, code)
	switch market {
	case 0:
		return normalizeSymbol(code + ".SZ"), sourceSymbol, 100, nil
	case 1:
		return normalizeSymbol(code + ".SH"), sourceSymbol, 100, nil
	case 116:
		if value, err := strconv.Atoi(code); err == nil {
			return normalizeSymbol(fmt.Sprintf("%d.HK", value)), sourceSymbol, 1000, nil
		}
		return normalizeSymbol(code + ".HK"), sourceSymbol, 1000, nil
	default:
		return "", "", 0, fmt.Errorf("unsupported eastmoney market: %d", market)
	}
}

func stringField(data map[string]any, key string) (string, error) {
	value, ok := data[key]
	if !ok {
		return "", fmt.Errorf("missing %s", key)
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed), nil
	case float64:
		return strconv.FormatInt(int64(typed), 10), nil
	default:
		return "", fmt.Errorf("invalid %s", key)
	}
}

func eastmoneySymbol(symbol string) (string, float64, error) {
	symbol = normalizeSymbol(symbol)
	switch {
	case strings.HasSuffix(symbol, ".SH"):
		return "1." + strings.TrimSuffix(symbol, ".SH"), 100, nil
	case strings.HasSuffix(symbol, ".SZ"):
		return "0." + strings.TrimSuffix(symbol, ".SZ"), 100, nil
	case strings.HasSuffix(symbol, ".HK"):
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconv.Atoi(code); err == nil {
			code = fmt.Sprintf("%05d", value)
		}
		return "116." + code, 1000, nil
	default:
		return "", 0, fmt.Errorf("unsupported eastmoney symbol: %s", symbol)
	}
}

func currencyForSymbol(symbol string) string {
	symbol = normalizeSymbol(symbol)
	switch {
	case strings.HasSuffix(symbol, ".HK"):
		return "HKD"
	case strings.HasSuffix(symbol, ".SH"), strings.HasSuffix(symbol, ".SZ"):
		return "CNY"
	default:
		return ""
	}
}

func loadLocation(name string) *time.Location {
	if strings.TrimSpace(name) == "" {
		return time.UTC
	}
	location, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return location
}

func closeDate(timestamps []int64, index int, location *time.Location) string {
	if index < 0 || index >= len(timestamps) || timestamps[index] <= 0 {
		return ""
	}
	return time.Unix(timestamps[index], 0).In(location).Format("2006-01-02")
}

func yahooSymbol(symbol string) string {
	symbol = normalizeSymbol(symbol)
	if strings.HasSuffix(symbol, ".SH") {
		return strings.TrimSuffix(symbol, ".SH") + ".SS"
	}
	return symbol
}
