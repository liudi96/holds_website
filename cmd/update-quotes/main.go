package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const decisionLogLimit = 500

type AppState struct {
	TotalCapital float64            `json:"totalCapital"`
	Cash         float64            `json:"cash"`
	FX           map[string]float64 `json:"fx"`
	Trades       []Trade            `json:"trades"`
	DecisionLogs []DecisionLog      `json:"decisionLogs"`
	Holdings     []Holding          `json:"holdings"`
	Plan         []PlanItem         `json:"plan"`
	Candidates   []Candidate        `json:"candidates"`
	Rules        []Rule             `json:"rules"`
}

type Trade struct {
	ID           int64   `json:"id"`
	Date         string  `json:"date"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	Side         string  `json:"side"`
	Shares       float64 `json:"shares"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	CurrentPrice float64 `json:"currentPrice"`
}

type DecisionLog struct {
	ID         int64    `json:"id"`
	Date       string   `json:"date"`
	Type       string   `json:"type"`
	Symbol     string   `json:"symbol,omitempty"`
	Name       string   `json:"name,omitempty"`
	Price      *float64 `json:"price,omitempty"`
	Currency   string   `json:"currency,omitempty"`
	Decision   string   `json:"decision"`
	Discipline string   `json:"discipline"`
	Detail     string   `json:"detail,omitempty"`
}

type Holding struct {
	Symbol            string   `json:"symbol"`
	Name              string   `json:"name"`
	Shares            float64  `json:"shares"`
	Cost              float64  `json:"cost"`
	CurrentPrice      float64  `json:"currentPrice"`
	PreviousClose     float64  `json:"previousClose"`
	CurrentPriceDate  string   `json:"currentPriceDate"`
	PreviousCloseDate string   `json:"previousCloseDate"`
	Action            string   `json:"action"`
	Status            string   `json:"status"`
	MarginOfSafety    *float64 `json:"marginOfSafety"`
	QualityScore      *float64 `json:"qualityScore"`
	Risk              string   `json:"risk"`
	Industry          string   `json:"industry"`
	Currency          string   `json:"currency"`
	IntrinsicValue    *float64 `json:"intrinsicValue"`
	FairValueRange    string   `json:"fairValueRange"`
	TargetBuyPrice    *float64 `json:"targetBuyPrice"`
	BusinessModel     *float64 `json:"businessModel"`
	Moat              *float64 `json:"moat"`
	Governance        *float64 `json:"governance"`
	FinancialQuality  *float64 `json:"financialQuality"`
	UpdatedAt         string   `json:"updatedAt"`
	Notes             string   `json:"notes"`
	Reports           []Report `json:"reports,omitempty"`
}

type Report struct {
	Period string `json:"period"`
	Kind   string `json:"kind"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Source string `json:"source"`
	URL    string `json:"url"`
}

type PlanItem struct {
	Rank       int    `json:"rank"`
	Symbol     string `json:"symbol,omitempty"`
	Name       string `json:"name"`
	Priority   string `json:"priority"`
	Advice     string `json:"advice"`
	Discipline string `json:"discipline"`
}

type Candidate struct {
	Symbol            string   `json:"symbol"`
	Name              string   `json:"name"`
	Status            string   `json:"status"`
	Action            string   `json:"action"`
	CurrentPrice      float64  `json:"currentPrice"`
	PreviousClose     float64  `json:"previousClose"`
	CurrentPriceDate  string   `json:"currentPriceDate"`
	PreviousCloseDate string   `json:"previousCloseDate"`
	MarginOfSafety    *float64 `json:"marginOfSafety"`
	QualityScore      *float64 `json:"qualityScore"`
	Risk              string   `json:"risk"`
	Industry          string   `json:"industry"`
	Currency          string   `json:"currency"`
	IntrinsicValue    *float64 `json:"intrinsicValue"`
	FairValueRange    string   `json:"fairValueRange"`
	TargetBuyPrice    *float64 `json:"targetBuyPrice"`
	BusinessModel     *float64 `json:"businessModel"`
	Moat              *float64 `json:"moat"`
	Governance        *float64 `json:"governance"`
	FinancialQuality  *float64 `json:"financialQuality"`
	UpdatedAt         string   `json:"updatedAt"`
	Notes             string   `json:"notes"`
	Reports           []Report `json:"reports,omitempty"`
}

type Rule struct {
	Dimension string  `json:"dimension"`
	Score     float64 `json:"score"`
	Standard  string  `json:"standard"`
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
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}

func main() {
	dataPath := flag.String("data", "data/portfolio.json", "portfolio JSON file to update")
	dryRun := flag.Bool("dry-run", false, "print updates without writing the file")
	flag.Parse()

	state, err := loadState(*dataPath)
	if err != nil {
		fail(err)
	}

	client := &http.Client{Timeout: 12 * time.Second}
	now := time.Now().Format("2006-01-02 15:04:05")
	updated := 0
	cache := make(map[string]quote)

	for i := range state.Holdings {
		holding := &state.Holdings[i]
		if strings.TrimSpace(holding.Symbol) == "" {
			continue
		}

		quote, err := fetchQuoteCached(client, cache, holding.Symbol)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", holding.Symbol, err)
			continue
		}

		fmt.Printf("%s %s: %.4f -> %.4f (%s), yesterday close %.4f (%s) [%s]\n", holding.Symbol, holding.Name, holding.CurrentPrice, quote.Price, quote.PriceDate, quote.PreviousClose, quote.PreviousCloseDate, quote.SourceSymbol)
		holding.CurrentPrice = quote.Price
		holding.PreviousClose = quote.PreviousClose
		holding.CurrentPriceDate = quote.PriceDate
		holding.PreviousCloseDate = quote.PreviousCloseDate
		holding.MarginOfSafety = marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, holding.MarginOfSafety)
		holding.UpdatedAt = fmt.Sprintf("%s；行情源 Yahoo Finance 日线收盘价；代码 %s；币种 %s；收盘日 %s/%s", now, quote.SourceSymbol, quote.Currency, quote.PreviousCloseDate, quote.PriceDate)
		appendQuoteDecisionLog(&state, holding.Symbol, holding.Name, holding.Currency, holding.CurrentPrice, holding.CurrentPriceDate, holding.PreviousCloseDate, now)
		updated++
	}

	for i := range state.Candidates {
		candidate := &state.Candidates[i]
		if strings.TrimSpace(candidate.Symbol) == "" {
			continue
		}

		quote, err := fetchQuoteCached(client, cache, candidate.Symbol)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", candidate.Symbol, err)
			continue
		}

		fmt.Printf("%s %s: %.4f -> %.4f (%s), yesterday close %.4f (%s) [%s]\n", candidate.Symbol, candidate.Name, candidate.CurrentPrice, quote.Price, quote.PriceDate, quote.PreviousClose, quote.PreviousCloseDate, quote.SourceSymbol)
		candidate.CurrentPrice = quote.Price
		candidate.PreviousClose = quote.PreviousClose
		candidate.CurrentPriceDate = quote.PriceDate
		candidate.PreviousCloseDate = quote.PreviousCloseDate
		candidate.MarginOfSafety = marginOfSafetyFromPrice(candidate.IntrinsicValue, candidate.CurrentPrice, candidate.MarginOfSafety)
		if strings.TrimSpace(candidate.Currency) == "" {
			candidate.Currency = strings.ToUpper(strings.TrimSpace(quote.Currency))
		}
		candidate.UpdatedAt = fmt.Sprintf("%s；行情源 Yahoo Finance 日线收盘价；代码 %s；币种 %s；收盘日 %s/%s", now, quote.SourceSymbol, quote.Currency, quote.PreviousCloseDate, quote.PriceDate)
		appendQuoteDecisionLog(&state, candidate.Symbol, candidate.Name, candidate.Currency, candidate.CurrentPrice, candidate.CurrentPriceDate, candidate.PreviousCloseDate, now)
		updated++
	}

	if updated == 0 {
		fail(errors.New("no quotes were updated"))
	}

	if *dryRun {
		fmt.Printf("dry run: %d quote records would be updated\n", updated)
		return
	}

	if err := saveState(*dataPath, state); err != nil {
		fail(err)
	}
	fmt.Printf("updated %d quote records in %s\n", updated, *dataPath)
}

type quote struct {
	Price             float64
	PreviousClose     float64
	PriceDate         string
	PreviousCloseDate string
	Currency          string
	SourceSymbol      string
}

func fetchQuote(client *http.Client, symbol string) (quote, error) {
	sourceSymbol := yahooSymbol(symbol)
	endpoint := "https://query1.finance.yahoo.com/v8/finance/chart/" + url.PathEscape(sourceSymbol) + "?range=5d&interval=1d"

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
	return quote{
		Price:             priceClose.Price,
		PreviousClose:     previousClose.Price,
		PriceDate:         priceClose.Date,
		PreviousCloseDate: previousClose.Date,
		Currency:          result.Meta.Currency,
		SourceSymbol:      sourceSymbol,
	}, nil
}

func fetchQuoteCached(client *http.Client, cache map[string]quote, symbol string) (quote, error) {
	normalized := strings.ToUpper(strings.TrimSpace(symbol))
	if cached, ok := cache[normalized]; ok {
		return cached, nil
	}
	quote, err := fetchQuote(client, normalized)
	if err != nil {
		return quote, err
	}
	cache[normalized] = quote
	return quote, nil
}

type dailyClose struct {
	Price float64
	Date  string
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

func marginOfSafetyFromPrice(intrinsicValue *float64, currentPrice float64, fallback *float64) *float64 {
	if intrinsicValue == nil || *intrinsicValue <= 0 || currentPrice <= 0 {
		return fallback
	}
	value := (*intrinsicValue - currentPrice) / *intrinsicValue
	return &value
}

func appendQuoteDecisionLog(state *AppState, symbol string, name string, currency string, currentPrice float64, currentDate string, previousDate string, now string) {
	_, _, _, decision, discipline := decisionLogContext(state, symbol)
	appendDecisionLog(state, DecisionLog{
		Date:       now,
		Type:       "quote",
		Symbol:     symbol,
		Name:       name,
		Price:      pricePointer(currentPrice),
		Currency:   currency,
		Decision:   decision,
		Discipline: discipline,
		Detail:     fmt.Sprintf("今收 %s；昨收 %s", firstNonEmpty(currentDate, "未知"), firstNonEmpty(previousDate, "未知")),
	})
}

func appendDecisionLog(state *AppState, entry DecisionLog) {
	entry.Type = strings.TrimSpace(entry.Type)
	if entry.Type == "" {
		entry.Type = "event"
	}
	entry.Symbol = normalizeSymbol(entry.Symbol)
	entry.Name = strings.TrimSpace(entry.Name)
	entry.Currency = normalizeSymbol(entry.Currency)
	entry.Decision = strings.TrimSpace(entry.Decision)
	entry.Discipline = strings.TrimSpace(entry.Discipline)
	entry.Detail = strings.TrimSpace(entry.Detail)
	if entry.ID == 0 {
		entry.ID = time.Now().UnixNano()
	}
	if strings.TrimSpace(entry.Date) == "" {
		entry.Date = time.Now().Format("2006-01-02 15:04:05")
	}

	state.DecisionLogs = append(state.DecisionLogs, entry)
	if len(state.DecisionLogs) > decisionLogLimit {
		state.DecisionLogs = state.DecisionLogs[len(state.DecisionLogs)-decisionLogLimit:]
	}
}

func decisionLogContext(state *AppState, symbol string) (string, *float64, string, string, string) {
	normalizedSymbol := normalizeSymbol(symbol)
	for i := range state.Holdings {
		holding := state.Holdings[i]
		if normalizeSymbol(holding.Symbol) != normalizedSymbol {
			continue
		}
		plan := findPlanForDecisionLog(state, holding.Symbol, holding.Name)
		return holding.Name, pricePointer(holding.CurrentPrice), holding.Currency, firstNonEmpty(holding.Action, holding.Status), firstNonEmpty(planDiscipline(plan), holding.Status)
	}

	for i := range state.Candidates {
		candidate := state.Candidates[i]
		if normalizeSymbol(candidate.Symbol) != normalizedSymbol {
			continue
		}
		plan := findPlanForDecisionLog(state, candidate.Symbol, candidate.Name)
		return candidate.Name, pricePointer(candidate.CurrentPrice), candidate.Currency, firstNonEmpty(candidate.Action, candidate.Status), firstNonEmpty(planDiscipline(plan), candidate.Status)
	}

	plan := findPlanForDecisionLog(state, symbol, "")
	return "", nil, "", "", planDiscipline(plan)
}

func findPlanForDecisionLog(state *AppState, symbol string, name string) *PlanItem {
	normalizedSymbol := normalizeSymbol(symbol)
	normalizedName := strings.TrimSpace(name)
	for i := range state.Plan {
		itemSymbol := normalizeSymbol(state.Plan[i].Symbol)
		if itemSymbol != "" && normalizedSymbol != "" && itemSymbol == normalizedSymbol {
			return &state.Plan[i]
		}
		itemName := strings.TrimSpace(state.Plan[i].Name)
		if itemName != "" && normalizedName != "" && (strings.EqualFold(itemName, normalizedName) || strings.Contains(normalizedName, itemName) || strings.Contains(itemName, normalizedName)) {
			return &state.Plan[i]
		}
	}
	return nil
}

func planDiscipline(plan *PlanItem) string {
	if plan == nil {
		return ""
	}
	return strings.TrimSpace(plan.Discipline)
}

func pricePointer(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if text := strings.TrimSpace(value); text != "" {
			return text
		}
	}
	return ""
}

func normalizeSymbol(symbol string) string {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if strings.HasSuffix(symbol, ".HK") {
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconv.Atoi(code); err == nil {
			return fmt.Sprintf("%04d.HK", value)
		}
	}
	return symbol
}

func loadState(path string) (AppState, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return AppState{}, err
	}

	var state AppState
	if err := json.Unmarshal(body, &state); err != nil {
		return AppState{}, err
	}
	return state, nil
}

func saveState(path string, state AppState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(path, body, 0o644)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
