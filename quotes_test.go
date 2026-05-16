package main

import (
	"encoding/json"
	"testing"
)

func ptrFloat(value float64) *float64 {
	return &value
}

func TestApplyDividendQuotePreservesAuditedDividend(t *testing.T) {
	current := &Dividend{
		FiscalYear:           "FY2025",
		DividendPerShare:     ptrFloat(4.3),
		DividendCurrency:     "CNY",
		CashDividendTotal:    ptrFloat(32360630265),
		CashDividendCurrency: "CNY",
	}
	trailingDividend := 4.0

	applyDividendQuote(&current, quote{
		DividendPerShare:   &trailingDividend,
		DividendCurrency:   "CNY",
		DividendFiscalYear: "TTM 2025-11-18",
		Currency:           "CNY",
	}, "CNY")

	if current.FiscalYear != "FY2025" {
		t.Fatalf("expected audited fiscal year to be preserved, got %q", current.FiscalYear)
	}
	if current.DividendPerShare == nil || *current.DividendPerShare != 4.3 {
		t.Fatalf("expected audited dividend per share to be preserved, got %v", current.DividendPerShare)
	}
}

func TestApplyDividendQuotePreservesAuditedDividendWithoutCashTotal(t *testing.T) {
	current := &Dividend{
		FiscalYear:       "FY2025",
		DividendPerShare: ptrFloat(1.38),
		DividendCurrency: "CNY",
	}
	trailingDividend := 1.7

	applyDividendQuote(&current, quote{
		DividendPerShare:   &trailingDividend,
		DividendCurrency:   "CNY",
		DividendFiscalYear: "TTM 2025-12-17",
		Currency:           "CNY",
	}, "CNY")

	if current.FiscalYear != "FY2025" {
		t.Fatalf("expected audited fiscal year to be preserved, got %q", current.FiscalYear)
	}
	if current.DividendPerShare == nil || *current.DividendPerShare != 1.38 {
		t.Fatalf("expected audited dividend per share to be preserved, got %v", current.DividendPerShare)
	}
}

func TestApplyDividendQuoteFillsMissingDividend(t *testing.T) {
	var current *Dividend
	trailingDividend := 1.2

	applyDividendQuote(&current, quote{
		DividendPerShare:   &trailingDividend,
		DividendCurrency:   "HKD",
		DividendFiscalYear: "TTM 2025-06-01",
		Currency:           "HKD",
	}, "HKD")

	if current == nil {
		t.Fatal("expected dividend to be created")
	}
	if current.FiscalYear != "TTM 2025-06-01" {
		t.Fatalf("expected quote fiscal year, got %q", current.FiscalYear)
	}
	if current.DividendPerShare == nil || *current.DividendPerShare != 1.2 {
		t.Fatalf("expected quote dividend per share, got %v", current.DividendPerShare)
	}
}

func TestPersistentStateDoesNotMutateRuntimeQuotes(t *testing.T) {
	state := AppState{
		Holdings: []Holding{{
			Symbol:            "0506.HK",
			CurrentPrice:      3.72,
			PreviousClose:     3.70,
			MarketCap:         ptrFloat(10290000000),
			MarketCapCurrency: "HKD",
			CurrentPriceDate:  "2026-05-11",
			PreviousCloseDate: "2026-05-08",
			UpdatedAt:         "行情源 Yahoo Finance",
		}},
		Candidates: []Candidate{{
			Symbol:            "600519.SH",
			CurrentPrice:      1372.99,
			PreviousClose:     1370,
			MarketCap:         ptrFloat(1720000000000),
			MarketCapCurrency: "CNY",
			CurrentPriceDate:  "2026-05-11",
			PreviousCloseDate: "2026-05-08",
			UpdatedAt:         "行情源 Yahoo Finance",
		}},
	}

	persisted := persistentState(state)

	if persisted.Holdings[0].CurrentPrice != 0 || persisted.Holdings[0].CurrentPriceDate != "" {
		t.Fatalf("expected persisted holding runtime quote to be cleared, got %+v", persisted.Holdings[0])
	}
	if persisted.Holdings[0].MarketCap != nil || persisted.Holdings[0].MarketCapCurrency != "" {
		t.Fatalf("expected persisted holding market cap to be cleared, got %+v", persisted.Holdings[0])
	}
	if persisted.Candidates[0].CurrentPrice != 0 || persisted.Candidates[0].CurrentPriceDate != "" {
		t.Fatalf("expected persisted candidate runtime quote to be cleared, got %+v", persisted.Candidates[0])
	}
	if persisted.Candidates[0].MarketCap != nil || persisted.Candidates[0].MarketCapCurrency != "" {
		t.Fatalf("expected persisted candidate market cap to be cleared, got %+v", persisted.Candidates[0])
	}
	if state.Holdings[0].CurrentPrice != 3.72 || state.Holdings[0].CurrentPriceDate != "2026-05-11" {
		t.Fatalf("persistentState mutated source holding, got %+v", state.Holdings[0])
	}
	if state.Holdings[0].MarketCap == nil || *state.Holdings[0].MarketCap != 10290000000 || state.Holdings[0].MarketCapCurrency != "HKD" {
		t.Fatalf("persistentState mutated source holding market cap, got %+v", state.Holdings[0])
	}
	if state.Candidates[0].CurrentPrice != 1372.99 || state.Candidates[0].CurrentPriceDate != "2026-05-11" {
		t.Fatalf("persistentState mutated source candidate, got %+v", state.Candidates[0])
	}
	if state.Candidates[0].MarketCap == nil || *state.Candidates[0].MarketCap != 1720000000000 || state.Candidates[0].MarketCapCurrency != "CNY" {
		t.Fatalf("persistentState mutated source candidate market cap, got %+v", state.Candidates[0])
	}
}

func TestMergeQuoteSupplementFillsMarketCap(t *testing.T) {
	item := quote{Price: 10, Currency: "HKD"}
	mergeQuoteSupplement(&item, quote{
		MarketCap:         ptrFloat(1000),
		MarketCapCurrency: "HKD",
	})

	if item.MarketCap == nil || *item.MarketCap != 1000 {
		t.Fatalf("market cap = %v, want 1000", item.MarketCap)
	}
	if item.MarketCapCurrency != "HKD" {
		t.Fatalf("market cap currency = %q, want HKD", item.MarketCapCurrency)
	}
}

func TestTencentMarketCapUsesHundredMillionUnit(t *testing.T) {
	fields := make([]string, 46)
	fields[45] = "104.0567"

	marketCap := tencentMarketCap(fields)
	if marketCap == nil || *marketCap != 10405670000 {
		t.Fatalf("market cap = %v, want 10405670000", marketCap)
	}
}

func TestParseEastmoneyDailyClose(t *testing.T) {
	close, err := parseEastmoneyDailyClose("2026-05-15,82.00,82.57,83.10,81.90,100,200")
	if err != nil {
		t.Fatalf("parseEastmoneyDailyClose() error = %v", err)
	}
	if close.Date != "2026-05-15" || close.Price != 82.57 {
		t.Fatalf("daily close = %+v, want 2026-05-15 82.57", close)
	}
}

func TestParseTencentKlineRows(t *testing.T) {
	rows, err := parseTencentKlineRows(json.RawMessage(`[
		["2026-05-14","81.00","82.57","83.10","80.90","100"],
		["2026-05-15","82.00","83.20","83.50","81.90","120"]
	]`))
	if err != nil {
		t.Fatalf("parseTencentKlineRows() error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2", len(rows))
	}
	if rows[1].Date != "2026-05-15" || rows[1].Price != 83.20 {
		t.Fatalf("second row = %+v, want 2026-05-15 83.20", rows[1])
	}
}

func TestApplyTwentyDayChangeFromClosesUsesCurrentPrice(t *testing.T) {
	closes := make([]dailyClose, 21)
	for i := range closes {
		closes[i] = dailyClose{Date: "2026-05-01", Price: 100 + float64(i)}
	}
	closes[0] = dailyClose{Date: "2026-04-15", Price: 80}
	item := quote{Price: 100}

	applyTwentyDayChangeFromCloses(&item, closes)

	if item.TwentyDayClose != 80 || item.TwentyDayCloseDate != "2026-04-15" {
		t.Fatalf("twenty day base = %.2f %q, want 80 2026-04-15", item.TwentyDayClose, item.TwentyDayCloseDate)
	}
	if item.TwentyDayChange == nil || *item.TwentyDayChange != 0.25 {
		t.Fatalf("twenty day change = %v, want 0.25", item.TwentyDayChange)
	}
}

func TestPreserveRuntimeTwentyDayQuote(t *testing.T) {
	existingChange := 0.12
	record := RuntimeQuote{Symbol: "0700.HK", CurrentPrice: 456.4}

	preserveRuntimeTwentyDayQuote(&record, RuntimeQuote{
		Symbol:             "0700.HK",
		TwentyDayClose:     407.5,
		TwentyDayCloseDate: "2026-04-16",
		TwentyDayChange:    &existingChange,
	})

	if record.TwentyDayClose != 407.5 || record.TwentyDayCloseDate != "2026-04-16" {
		t.Fatalf("twenty day quote = %.2f %q, want 407.5 2026-04-16", record.TwentyDayClose, record.TwentyDayCloseDate)
	}
	if record.TwentyDayChange == nil || *record.TwentyDayChange != existingChange {
		t.Fatalf("twenty day change = %v, want %v", record.TwentyDayChange, existingChange)
	}
}
