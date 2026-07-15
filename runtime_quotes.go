package main

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
	"time"
)

type RuntimeQuoteBook struct {
	UpdatedAt       string                   `json:"updatedAt,omitempty"`
	Quotes          map[string]RuntimeQuote  `json:"quotes"`
	ETFRuleStatuses map[string]ETFRuleStatus `json:"etfRuleStatuses,omitempty"`
}

type RuntimeQuote struct {
	Symbol             string   `json:"symbol"`
	CurrentPrice       float64  `json:"currentPrice,omitempty"`
	PreviousClose      float64  `json:"previousClose,omitempty"`
	TwentyDayClose     float64  `json:"twentyDayClose,omitempty"`
	TwentyDayCloseDate string   `json:"twentyDayCloseDate,omitempty"`
	TwentyDayChange    *float64 `json:"twentyDayChange,omitempty"`
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
	if book.ETFRuleStatuses == nil {
		book.ETFRuleStatuses = map[string]ETFRuleStatus{}
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
	normalizedETFStatuses := make(map[string]ETFRuleStatus, len(book.ETFRuleStatuses))
	for key, status := range book.ETFRuleStatuses {
		symbol := normalizeFundSymbol(firstNonEmpty(status.Symbol, key))
		if symbol == "" {
			continue
		}
		status.Symbol = symbol
		normalizedETFStatuses[symbol] = status
	}
	book.ETFRuleStatuses = normalizedETFStatuses
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
		preserveRuntimeTwentyDayQuote(&record, book.Quotes[symbol])
		book.Quotes[symbol] = record
	}
	return saveRuntimeQuoteBook(book)
}

func saveRuntimeMarketData(records []RuntimeQuote, statuses []ETFRuleStatus, updatedAt string) error {
	book, err := loadRuntimeQuoteBook()
	if err != nil {
		return err
	}
	if strings.TrimSpace(updatedAt) != "" {
		book.UpdatedAt = strings.TrimSpace(updatedAt)
	}
	if book.Quotes == nil {
		book.Quotes = map[string]RuntimeQuote{}
	}
	if book.ETFRuleStatuses == nil {
		book.ETFRuleStatuses = map[string]ETFRuleStatus{}
	}
	if statuses != nil {
		activeSymbols := make(map[string]struct{}, len(etfRuleConfigs))
		for _, config := range etfRuleConfigs {
			activeSymbols[normalizeFundSymbol(config.Symbol)] = struct{}{}
		}
		for symbol := range book.ETFRuleStatuses {
			if _, active := activeSymbols[normalizeFundSymbol(symbol)]; !active {
				delete(book.ETFRuleStatuses, symbol)
			}
		}
	}
	for _, record := range records {
		symbol := normalizeSymbol(record.Symbol)
		if symbol == "" {
			continue
		}
		record.Symbol = symbol
		preserveRuntimeTwentyDayQuote(&record, book.Quotes[symbol])
		book.Quotes[symbol] = record
	}
	for _, status := range statuses {
		symbol := normalizeFundSymbol(status.Symbol)
		if symbol == "" {
			continue
		}
		status.Symbol = symbol
		qualityTime := runtimeQuoteBookUpdateTime(updatedAt)
		status = mergeETFRuleStatusWithExisting(status, book.ETFRuleStatuses[symbol], qualityTime)
		applyETFStatusDataQuality(&status, qualityTime)
		book.ETFRuleStatuses[symbol] = status
	}
	return saveRuntimeQuoteBook(book)
}

func runtimeQuoteBookUpdateTime(updatedAt string) time.Time {
	if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(updatedAt), time.Local); err == nil {
		return parsed
	}
	return time.Now()
}

func preserveRuntimeTwentyDayQuote(record *RuntimeQuote, existing RuntimeQuote) {
	if record == nil || record.TwentyDayChange != nil || existing.TwentyDayChange == nil {
		return
	}
	record.TwentyDayClose = existing.TwentyDayClose
	record.TwentyDayCloseDate = existing.TwentyDayCloseDate
	record.TwentyDayChange = existing.TwentyDayChange
}

func saveRuntimeQuoteBook(book RuntimeQuoteBook) error {
	if book.Quotes == nil {
		book.Quotes = map[string]RuntimeQuote{}
	}
	body, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return writeFileAtomic(runtimeQuotesFile, body, 0o644)
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
	for i := range state.Funds {
		record, ok := book.Quotes[normalizeFundSymbol(state.Funds[i].Symbol)]
		if ok {
			applyRuntimeQuoteToFund(&state.Funds[i], record)
		}
	}
	state.ETFRuleStatuses = runtimeETFRuleStatusList(book.ETFRuleStatuses)
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
	if record.TwentyDayClose > 0 {
		holding.TwentyDayClose = record.TwentyDayClose
	}
	if strings.TrimSpace(record.TwentyDayCloseDate) != "" {
		holding.TwentyDayCloseDate = record.TwentyDayCloseDate
	}
	if record.TwentyDayChange != nil {
		holding.TwentyDayChange = record.TwentyDayChange
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
	if record.TwentyDayClose > 0 {
		candidate.TwentyDayClose = record.TwentyDayClose
	}
	if strings.TrimSpace(record.TwentyDayCloseDate) != "" {
		candidate.TwentyDayCloseDate = record.TwentyDayCloseDate
	}
	if record.TwentyDayChange != nil {
		candidate.TwentyDayChange = record.TwentyDayChange
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
		TwentyDayClose:     quote.TwentyDayClose,
		TwentyDayCloseDate: quote.TwentyDayCloseDate,
		TwentyDayChange:    quote.TwentyDayChange,
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
		TwentyDayClose:     record.TwentyDayClose,
		TwentyDayCloseDate: record.TwentyDayCloseDate,
		TwentyDayChange:    record.TwentyDayChange,
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
	normalizePortfolioState(&state)
	state.Holdings = append([]Holding(nil), state.Holdings...)
	state.Candidates = append([]Candidate(nil), state.Candidates...)
	state.Stocks = append([]Stock(nil), state.Stocks...)
	state.Funds = append([]Fund(nil), state.Funds...)
	state.DataStatus = nil
	state.ETFRuleStatuses = nil
	for i := range state.Holdings {
		clearHoldingRuntimeQuote(&state.Holdings[i])
	}
	for i := range state.Candidates {
		clearCandidateRuntimeQuote(&state.Candidates[i])
	}
	for i := range state.Stocks {
		clearStockRuntimeQuote(&state.Stocks[i])
	}
	for i := range state.Funds {
		clearFundRuntimeQuote(&state.Funds[i])
	}
	return state
}

func applyRuntimeQuoteToFund(fund *Fund, record RuntimeQuote) {
	if record.CurrentPrice > 0 {
		fund.CurrentPrice = record.CurrentPrice
	}
	if record.PreviousClose > 0 {
		fund.PreviousClose = record.PreviousClose
	}
	fund.CurrentPriceDate = firstNonEmpty(record.CurrentPriceDate, fund.CurrentPriceDate)
	fund.PreviousCloseDate = firstNonEmpty(record.PreviousCloseDate, fund.PreviousCloseDate)
	if strings.TrimSpace(fund.Currency) == "" {
		fund.Currency = strings.ToUpper(strings.TrimSpace(record.Currency))
	}
	if strings.TrimSpace(record.UpdatedAt) != "" {
		fund.UpdatedAt = record.UpdatedAt
	}
	*fund = normalizeFund(*fund)
}

func clearFundRuntimeQuote(fund *Fund) {
	fund.CurrentPrice = 0
	fund.PreviousClose = 0
	fund.CurrentPriceDate = ""
	fund.PreviousCloseDate = ""
	if strings.Contains(fund.UpdatedAt, "行情") || strings.Contains(fund.UpdatedAt, "quote") || strings.Contains(fund.UpdatedAt, "NAV") {
		fund.UpdatedAt = ""
	}
}

func clearStockRuntimeQuote(stock *Stock) {
	stock.CurrentPrice = 0
	stock.PreviousClose = 0
	stock.TwentyDayClose = 0
	stock.TwentyDayCloseDate = ""
	stock.TwentyDayChange = nil
	stock.MarketCap = nil
	stock.MarketCapCurrency = ""
	stock.CurrentPriceDate = ""
	stock.PreviousCloseDate = ""
	if strings.Contains(stock.UpdatedAt, "行情源") {
		stock.UpdatedAt = ""
	}
}

func clearHoldingRuntimeQuote(holding *Holding) {
	holding.CurrentPrice = 0
	holding.PreviousClose = 0
	holding.TwentyDayClose = 0
	holding.TwentyDayCloseDate = ""
	holding.TwentyDayChange = nil
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
	candidate.TwentyDayClose = 0
	candidate.TwentyDayCloseDate = ""
	candidate.TwentyDayChange = nil
	candidate.MarketCap = nil
	candidate.MarketCapCurrency = ""
	candidate.CurrentPriceDate = ""
	candidate.PreviousCloseDate = ""
	if strings.Contains(candidate.UpdatedAt, "行情源") {
		candidate.UpdatedAt = ""
	}
}
