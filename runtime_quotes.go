package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type RuntimeQuoteBook struct {
	UpdatedAt string                  `json:"updatedAt,omitempty"`
	Quotes    map[string]RuntimeQuote `json:"quotes"`
}

type RuntimeQuote struct {
	Symbol             string   `json:"symbol"`
	CurrentPrice       float64  `json:"currentPrice,omitempty"`
	PreviousClose      float64  `json:"previousClose,omitempty"`
	MarketCap          *float64 `json:"marketCap,omitempty"`
	MarketCapCurrency  string   `json:"marketCapCurrency,omitempty"`
	CurrentPriceDate   string   `json:"currentPriceDate,omitempty"`
	PreviousCloseDate  string   `json:"previousCloseDate,omitempty"`
	Currency           string   `json:"currency,omitempty"`
	SourceSymbol       string   `json:"sourceSymbol,omitempty"`
	SourceName         string   `json:"sourceName,omitempty"`
	UpdatedAt          string   `json:"updatedAt,omitempty"`
	DividendPerShare   *float64 `json:"dividendPerShare,omitempty"`
	DividendCurrency   string   `json:"dividendCurrency,omitempty"`
	DividendFiscalYear string   `json:"dividendFiscalYear,omitempty"`
}

func loadRuntimeQuoteBook() (RuntimeQuoteBook, error) {
	book := RuntimeQuoteBook{Quotes: map[string]RuntimeQuote{}}
	body, err := os.ReadFile(runtimeQuotesFile)
	if errors.Is(err, os.ErrNotExist) {
		return book, nil
	}
	if err != nil {
		return book, err
	}
	if err := json.Unmarshal(body, &book); err != nil {
		return book, err
	}
	if book.Quotes == nil {
		book.Quotes = map[string]RuntimeQuote{}
	}
	normalized := make(map[string]RuntimeQuote, len(book.Quotes))
	for key, record := range book.Quotes {
		symbol := normalizeSymbol(firstNonEmpty(record.Symbol, key))
		if symbol == "" {
			continue
		}
		record.Symbol = symbol
		normalized[symbol] = record
	}
	book.Quotes = normalized
	return book, nil
}

func saveRuntimeQuoteRecords(records []RuntimeQuote, updatedAt string) error {
	book, err := loadRuntimeQuoteBook()
	if err != nil {
		return err
	}
	if strings.TrimSpace(updatedAt) != "" {
		book.UpdatedAt = strings.TrimSpace(updatedAt)
	}
	for _, record := range records {
		symbol := normalizeSymbol(record.Symbol)
		if symbol == "" {
			continue
		}
		record.Symbol = symbol
		book.Quotes[symbol] = record
	}
	return saveRuntimeQuoteBook(book)
}

func saveRuntimeQuoteBook(book RuntimeQuoteBook) error {
	if book.Quotes == nil {
		book.Quotes = map[string]RuntimeQuote{}
	}
	if err := os.MkdirAll(filepath.Dir(runtimeQuotesFile), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(runtimeQuotesFile, body, 0o644)
}

func mergeRuntimeQuotes(state *AppState) error {
	book, err := loadRuntimeQuoteBook()
	if err != nil {
		return err
	}
	for i := range state.Holdings {
		record, ok := book.Quotes[normalizeSymbol(state.Holdings[i].Symbol)]
		if ok {
			applyRuntimeQuoteToHolding(&state.Holdings[i], record)
		}
	}
	for i := range state.Candidates {
		record, ok := book.Quotes[normalizeSymbol(state.Candidates[i].Symbol)]
		if ok {
			applyRuntimeQuoteToCandidate(&state.Candidates[i], record)
		}
	}
	return nil
}

func applyRuntimeQuoteToHolding(holding *Holding, record RuntimeQuote) {
	if record.CurrentPrice > 0 {
		holding.CurrentPrice = record.CurrentPrice
		holding.MarginOfSafety = marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, holding.MarginOfSafety)
	}
	if record.PreviousClose > 0 {
		holding.PreviousClose = record.PreviousClose
	}
	if record.MarketCap != nil && *record.MarketCap > 0 {
		holding.MarketCap = record.MarketCap
		holding.MarketCapCurrency = strings.ToUpper(firstNonEmpty(record.MarketCapCurrency, record.Currency, holding.Currency))
	}
	holding.CurrentPriceDate = firstNonEmpty(record.CurrentPriceDate, holding.CurrentPriceDate)
	holding.PreviousCloseDate = firstNonEmpty(record.PreviousCloseDate, holding.PreviousCloseDate)
	if strings.TrimSpace(holding.Currency) == "" {
		holding.Currency = strings.ToUpper(strings.TrimSpace(record.Currency))
	}
	applyDividendQuote(&holding.Dividend, runtimeRecordAsQuote(record), holding.Currency)
	if strings.TrimSpace(record.UpdatedAt) != "" {
		holding.UpdatedAt = record.UpdatedAt
	}
}

func applyRuntimeQuoteToCandidate(candidate *Candidate, record RuntimeQuote) {
	if record.CurrentPrice > 0 {
		candidate.CurrentPrice = record.CurrentPrice
		candidate.MarginOfSafety = marginOfSafetyFromPrice(candidate.IntrinsicValue, candidate.CurrentPrice, candidate.MarginOfSafety)
	}
	if record.PreviousClose > 0 {
		candidate.PreviousClose = record.PreviousClose
	}
	if record.MarketCap != nil && *record.MarketCap > 0 {
		candidate.MarketCap = record.MarketCap
		candidate.MarketCapCurrency = strings.ToUpper(firstNonEmpty(record.MarketCapCurrency, record.Currency, candidate.Currency))
	}
	candidate.CurrentPriceDate = firstNonEmpty(record.CurrentPriceDate, candidate.CurrentPriceDate)
	candidate.PreviousCloseDate = firstNonEmpty(record.PreviousCloseDate, candidate.PreviousCloseDate)
	if strings.TrimSpace(candidate.Currency) == "" {
		candidate.Currency = strings.ToUpper(strings.TrimSpace(record.Currency))
	}
	applyDividendQuote(&candidate.Dividend, runtimeRecordAsQuote(record), candidate.Currency)
	if strings.TrimSpace(record.UpdatedAt) != "" {
		candidate.UpdatedAt = record.UpdatedAt
	}
}

func runtimeQuoteFromQuote(symbol string, quote quote, updateLabel string) RuntimeQuote {
	return RuntimeQuote{
		Symbol:             normalizeSymbol(symbol),
		CurrentPrice:       quote.Price,
		PreviousClose:      quote.PreviousClose,
		MarketCap:          quote.MarketCap,
		MarketCapCurrency:  strings.ToUpper(strings.TrimSpace(quote.MarketCapCurrency)),
		CurrentPriceDate:   quote.PriceDate,
		PreviousCloseDate:  quote.PreviousCloseDate,
		Currency:           strings.ToUpper(strings.TrimSpace(quote.Currency)),
		SourceSymbol:       quote.SourceSymbol,
		SourceName:         quote.SourceName,
		UpdatedAt:          quoteUpdateLabel(updateLabel, quote),
		DividendPerShare:   quote.DividendPerShare,
		DividendCurrency:   strings.ToUpper(strings.TrimSpace(quote.DividendCurrency)),
		DividendFiscalYear: quote.DividendFiscalYear,
	}
}

func runtimeRecordAsQuote(record RuntimeQuote) quote {
	return quote{
		Price:              record.CurrentPrice,
		PreviousClose:      record.PreviousClose,
		MarketCap:          record.MarketCap,
		MarketCapCurrency:  record.MarketCapCurrency,
		PriceDate:          record.CurrentPriceDate,
		PreviousCloseDate:  record.PreviousCloseDate,
		Currency:           record.Currency,
		SourceSymbol:       record.SourceSymbol,
		SourceName:         record.SourceName,
		DividendPerShare:   record.DividendPerShare,
		DividendCurrency:   record.DividendCurrency,
		DividendFiscalYear: record.DividendFiscalYear,
	}
}

func runtimeQuoteList(records map[string]RuntimeQuote) []RuntimeQuote {
	keys := make([]string, 0, len(records))
	for key := range records {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	list := make([]RuntimeQuote, 0, len(keys))
	for _, key := range keys {
		list = append(list, records[key])
	}
	return list
}

func persistentState(state AppState) AppState {
	state.Holdings = append([]Holding(nil), state.Holdings...)
	state.Candidates = append([]Candidate(nil), state.Candidates...)
	state.Funds = append([]Fund(nil), state.Funds...)
	state.Industries = nil
	for i := range state.Holdings {
		clearHoldingRuntimeQuote(&state.Holdings[i])
	}
	for i := range state.Candidates {
		clearCandidateRuntimeQuote(&state.Candidates[i])
	}
	return state
}

func clearHoldingRuntimeQuote(holding *Holding) {
	holding.CurrentPrice = 0
	holding.PreviousClose = 0
	holding.MarketCap = nil
	holding.MarketCapCurrency = ""
	holding.CurrentPriceDate = ""
	holding.PreviousCloseDate = ""
	if strings.Contains(holding.UpdatedAt, "行情源") {
		holding.UpdatedAt = ""
	}
}

func clearCandidateRuntimeQuote(candidate *Candidate) {
	candidate.CurrentPrice = 0
	candidate.PreviousClose = 0
	candidate.MarketCap = nil
	candidate.MarketCapCurrency = ""
	candidate.CurrentPriceDate = ""
	candidate.PreviousCloseDate = ""
	if strings.Contains(candidate.UpdatedAt, "行情源") {
		candidate.UpdatedAt = ""
	}
}
