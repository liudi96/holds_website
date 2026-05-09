package main

import "testing"

func TestEastmoneyFinancialSecucode(t *testing.T) {
	cases := map[string]string{
		"600519.SH": "600519.SH",
		"000333.SZ": "000333.SZ",
		"700.HK":    "00700.HK",
		"0700.HK":   "00700.HK",
	}
	for symbol, want := range cases {
		got, err := eastmoneyFinancialSecucode(symbol)
		if err != nil {
			t.Fatalf("eastmoneyFinancialSecucode(%q) returned error: %v", symbol, err)
		}
		if got != want {
			t.Fatalf("eastmoneyFinancialSecucode(%q) = %q, want %q", symbol, got, want)
		}
	}
}

func TestFillDerivedGrowth(t *testing.T) {
	annual := []FinancialAnnual{
		{FiscalYear: "2025", Revenue: ptr(120), NetProfit: ptr(24), OperatingCashFlow: ptr(30)},
		{FiscalYear: "2024", Revenue: ptr(100), NetProfit: ptr(20), OperatingCashFlow: ptr(18)},
	}

	fillDerivedGrowth(annual)

	if annual[0].RevenueYoY == nil || *annual[0].RevenueYoY != 0.2 {
		t.Fatalf("expected revenue growth 0.2, got %v", annual[0].RevenueYoY)
	}
	if annual[0].NetProfitYoY == nil || *annual[0].NetProfitYoY != 0.2 {
		t.Fatalf("expected net profit growth 0.2, got %v", annual[0].NetProfitYoY)
	}
	if annual[0].OperatingCashFlowToRevenue == nil || *annual[0].OperatingCashFlowToRevenue != 0.25 {
		t.Fatalf("expected OCF/revenue 0.25, got %v", annual[0].OperatingCashFlowToRevenue)
	}
}

func TestHKCashflowByReportDate(t *testing.T) {
	rows := []map[string]any{
		{"REPORT_DATE": "2025-12-31 00:00:00", "STD_ITEM_CODE": "003999", "AMOUNT": float64(120)},
		{"REPORT_DATE": "2025-12-31 00:00:00", "STD_ITEM_CODE": "005005", "AMOUNT": float64(30)},
	}

	byDate := hkCashflowByReportDate(rows)
	operatingCashFlow := hkCashflowAmount(byDate["2025-12-31"], "003999")
	capex := hkCashflowAmount(byDate["2025-12-31"], "005005")
	fcf := freeCashFlow(operatingCashFlow, capex)

	if operatingCashFlow == nil || *operatingCashFlow != 120 {
		t.Fatalf("expected operating cash flow 120, got %v", operatingCashFlow)
	}
	if capex == nil || *capex != 30 {
		t.Fatalf("expected capex 30, got %v", capex)
	}
	if fcf == nil || *fcf != 90 {
		t.Fatalf("expected FCF 90, got %v", fcf)
	}
}

func TestQuotedList(t *testing.T) {
	got := quotedList([]string{"2025-12-31", "", "2024-12-31"})
	want := "'2025-12-31','2024-12-31'"
	if got != want {
		t.Fatalf("quotedList = %q, want %q", got, want)
	}
}
