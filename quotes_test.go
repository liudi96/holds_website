package main

import "testing"

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
			CurrentPriceDate:  "2026-05-11",
			PreviousCloseDate: "2026-05-08",
			UpdatedAt:         "行情源 Yahoo Finance",
		}},
		Candidates: []Candidate{{
			Symbol:            "600519.SH",
			CurrentPrice:      1372.99,
			PreviousClose:     1370,
			CurrentPriceDate:  "2026-05-11",
			PreviousCloseDate: "2026-05-08",
			UpdatedAt:         "行情源 Yahoo Finance",
		}},
	}

	persisted := persistentState(state)

	if persisted.Holdings[0].CurrentPrice != 0 || persisted.Holdings[0].CurrentPriceDate != "" {
		t.Fatalf("expected persisted holding runtime quote to be cleared, got %+v", persisted.Holdings[0])
	}
	if persisted.Candidates[0].CurrentPrice != 0 || persisted.Candidates[0].CurrentPriceDate != "" {
		t.Fatalf("expected persisted candidate runtime quote to be cleared, got %+v", persisted.Candidates[0])
	}
	if state.Holdings[0].CurrentPrice != 3.72 || state.Holdings[0].CurrentPriceDate != "2026-05-11" {
		t.Fatalf("persistentState mutated source holding, got %+v", state.Holdings[0])
	}
	if state.Candidates[0].CurrentPrice != 1372.99 || state.Candidates[0].CurrentPriceDate != "2026-05-11" {
		t.Fatalf("persistentState mutated source candidate, got %+v", state.Candidates[0])
	}
}
