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
	"strings"
	"time"
)

type AppState struct {
	TotalCapital float64            `json:"totalCapital"`
	Cash         float64            `json:"cash"`
	FX           map[string]float64 `json:"fx"`
	Trades       []Trade            `json:"trades"`
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
}

type PlanItem struct {
	Rank       int    `json:"rank"`
	Name       string `json:"name"`
	Priority   string `json:"priority"`
	Advice     string `json:"advice"`
	Discipline string `json:"discipline"`
}

type Candidate struct {
	Symbol         string   `json:"symbol"`
	Name           string   `json:"name"`
	Status         string   `json:"status"`
	Action         string   `json:"action"`
	MarginOfSafety *float64 `json:"marginOfSafety"`
	QualityScore   *float64 `json:"qualityScore"`
	Industry       string   `json:"industry"`
	Currency       string   `json:"currency"`
	IntrinsicValue *float64 `json:"intrinsicValue"`
	FairValueRange string   `json:"fairValueRange"`
	TargetBuyPrice *float64 `json:"targetBuyPrice"`
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

	for i := range state.Holdings {
		holding := &state.Holdings[i]
		if strings.TrimSpace(holding.Symbol) == "" || holding.CurrentPrice <= 0 {
			continue
		}

		quote, err := fetchQuote(client, holding.Symbol)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", holding.Symbol, err)
			continue
		}

		fmt.Printf("%s %s: %.4f -> %.4f (%s), yesterday close %.4f (%s) [%s]\n", holding.Symbol, holding.Name, holding.CurrentPrice, quote.Price, quote.PriceDate, quote.PreviousClose, quote.PreviousCloseDate, quote.SourceSymbol)
		holding.CurrentPrice = quote.Price
		holding.PreviousClose = quote.PreviousClose
		holding.CurrentPriceDate = quote.PriceDate
		holding.PreviousCloseDate = quote.PreviousCloseDate
		holding.UpdatedAt = fmt.Sprintf("%s；行情源 Yahoo Finance 日线收盘价；代码 %s；币种 %s；收盘日 %s/%s", now, quote.SourceSymbol, quote.Currency, quote.PreviousCloseDate, quote.PriceDate)
		updated++
	}

	if updated == 0 {
		fail(errors.New("no holdings were updated"))
	}

	if *dryRun {
		fmt.Printf("dry run: %d holdings would be updated\n", updated)
		return
	}

	if err := saveState(*dataPath, state); err != nil {
		fail(err)
	}
	fmt.Printf("updated %d holdings in %s\n", updated, *dataPath)
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
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if strings.HasSuffix(symbol, ".SH") {
		return strings.TrimSuffix(symbol, ".SH") + ".SS"
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
