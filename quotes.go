package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}

type quote struct {
	Price             float64
	PreviousClose     float64
	PriceDate         string
	PreviousCloseDate string
	Currency          string
	SourceSymbol      string
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
	updated, skipped := updateQuotes(&state, &http.Client{Timeout: 12 * time.Second}, now)
	if updated > 0 {
		appendQuoteDecisionLogs(&state, now)
		if err := saveState(state); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to save state")
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

func updateQuotes(state *AppState, client *http.Client, now time.Time) (int, []QuoteSkip) {
	updated := 0
	skipped := []QuoteSkip{}
	cache := make(map[string]quote)
	updateLabel := now.Format("2006-01-02 15:04:05")

	for i := range state.Holdings {
		holding := &state.Holdings[i]
		if strings.TrimSpace(holding.Symbol) == "" {
			continue
		}

		quote, err := fetchQuoteCached(client, cache, holding.Symbol)
		if err != nil {
			skipped = append(skipped, QuoteSkip{Type: "holding", Symbol: holding.Symbol, Name: holding.Name, Error: err.Error()})
			continue
		}

		applyHoldingQuote(holding, quote, updateLabel)
		updated++
	}

	for i := range state.Candidates {
		candidate := &state.Candidates[i]
		if strings.TrimSpace(candidate.Symbol) == "" {
			continue
		}

		quote, err := fetchQuoteCached(client, cache, candidate.Symbol)
		if err != nil {
			skipped = append(skipped, QuoteSkip{Type: "candidate", Symbol: candidate.Symbol, Name: candidate.Name, Error: err.Error()})
			continue
		}

		applyCandidateQuote(candidate, quote, updateLabel)
		updated++
	}

	return updated, skipped
}

func fetchQuoteCached(client *http.Client, cache map[string]quote, symbol string) (quote, error) {
	normalized := normalizeSymbol(symbol)
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

func applyHoldingQuote(holding *Holding, quote quote, updateLabel string) {
	holding.CurrentPrice = quote.Price
	holding.PreviousClose = quote.PreviousClose
	holding.CurrentPriceDate = quote.PriceDate
	holding.PreviousCloseDate = quote.PreviousCloseDate
	holding.MarginOfSafety = marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, holding.MarginOfSafety)
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
	if strings.TrimSpace(candidate.Currency) == "" {
		candidate.Currency = strings.ToUpper(strings.TrimSpace(quote.Currency))
	}
	candidate.UpdatedAt = quoteUpdateLabel(updateLabel, quote)
}

func quoteUpdateLabel(updateLabel string, quote quote) string {
	return fmt.Sprintf("%s；行情源 Yahoo Finance 日线收盘价；代码 %s；币种 %s；收盘日 %s/%s", updateLabel, quote.SourceSymbol, quote.Currency, quote.PreviousCloseDate, quote.PriceDate)
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
