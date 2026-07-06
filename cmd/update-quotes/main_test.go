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

func TestParseFundGZQuote(t *testing.T) {
	payload, err := parseFundGZQuote([]byte(`jsonpgz({"fundcode":"000001","name":"Test Fund","jzrq":"2026-07-03","dwjz":"1.2345","gsz":"1.2300","gztime":"2026-07-03 15:00"})`))
	if err != nil {
		t.Fatalf("parseFundGZQuote() error = %v", err)
	}
	if payload.FundCode != "000001" || payload.NAV != "1.2345" || payload.JZDate != "2026-07-03" {
		t.Fatalf("payload = %+v", payload)
	}
}

func TestQuoteSymbolsIncludesETFFundsOnly(t *testing.T) {
	state := AppState{
		Holdings: []Holding{{Symbol: "0700.HK"}},
		Funds: []Fund{
			{Symbol: "510300.SH", FundType: "etf"},
			{Symbol: "000001", FundType: "otc"},
			{Symbol: "004814.OF"},
		},
	}
	symbols := quoteSymbols(&state)
	if len(symbols) != 2 {
		t.Fatalf("symbols = %+v, want 2 entries", symbols)
	}
	if symbols[0] != "0700.HK" || symbols[1] != "510300.SH" {
		t.Fatalf("symbols = %+v, want stock plus ETF only", symbols)
	}
}
