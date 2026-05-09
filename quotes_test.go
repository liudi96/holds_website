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
