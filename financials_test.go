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

func TestApplyFinancialDividendUsesPerShareAndMarketCap(t *testing.T) {
	marketCap := 1000.0
	current := &Dividend{
		FiscalYear:       "TTM 2025-06-01",
		DividendPerShare: ptr(1.2),
		DividendCurrency: "HKD",
	}

	applyFinancialDividend(&current, &Financials{
		Currency: "HKD",
		Annual: []FinancialAnnual{{
			FiscalYear:       "2025",
			Currency:         "HKD",
			DividendPerShare: ptr(1.25),
			DividendCurrency: "HKD",
			BuybackAmount:    ptr(50),
			BuybackCurrency:  "HKD",
		}},
	}, "HKD", 10, &marketCap)

	if current.FiscalYear != "FY2025" {
		t.Fatalf("fiscal year = %q, want FY2025", current.FiscalYear)
	}
	if current.DividendPerShare == nil || *current.DividendPerShare != 1.25 {
		t.Fatalf("dividend per share = %v, want 1.25", current.DividendPerShare)
	}
	if current.CashDividendTotal == nil || *current.CashDividendTotal != 125 {
		t.Fatalf("cash dividend total = %v, want 125", current.CashDividendTotal)
	}
	if current.CashDividendCurrency != "HKD" {
		t.Fatalf("cash dividend currency = %q, want HKD", current.CashDividendCurrency)
	}
	if current.BuybackAmount == nil || *current.BuybackAmount != 50 {
		t.Fatalf("buyback amount = %v, want 50", current.BuybackAmount)
	}
	if current.BuybackCurrency != "HKD" {
		t.Fatalf("buyback currency = %q, want HKD", current.BuybackCurrency)
	}
}

func TestApplyFinancialDividendInfersCashTotalFromPerShare(t *testing.T) {
	marketCap := 1000.0
	current := &Dividend{
		FiscalYear:       "TTM 2026-01-16",
		DividendPerShare: ptr(3.0),
		DividendCurrency: "CNY",
		BuybackAmount:    ptr(20),
		BuybackCurrency:  "CNY",
	}

	applyFinancialDividend(&current, &Financials{
		Currency: "CNY",
		Annual: []FinancialAnnual{{
			FiscalYear: "2025",
			Currency:   "CNY",
		}},
	}, "CNY", 10, &marketCap)

	if current.FiscalYear != "FY2025" {
		t.Fatalf("fiscal year = %q, want FY2025", current.FiscalYear)
	}
	if current.CashDividendTotal == nil || *current.CashDividendTotal != 300 {
		t.Fatalf("cash dividend total = %v, want 300", current.CashDividendTotal)
	}
	if current.BuybackAmount != nil {
		t.Fatalf("expected stale buyback to be cleared, got %v", current.BuybackAmount)
	}
}

func TestApplyFinancialNetCashDerivesOrdinaryShareholderFCF(t *testing.T) {
	current := &NetCashProfile{
		Haircut:         ptrFloat(0.4),
		AdjustedNetCash: ptrFloat(10),
		ExCashPE:        ptrFloat(8),
		ExCashPFCF:      ptrFloat(9),
		FCFYield:        ptrFloat(0.1),
	}

	applyFinancialNetCash(&current, &Financials{
		Currency: "CNY",
		Annual: []FinancialAnnual{{
			FiscalYear:              "2025",
			Currency:                "CNY",
			CashAndShortInvestments: ptrFloat(100),
			InterestBearingDebt:     ptrFloat(20),
			FreeCashFlow:            ptrFloat(60),
			ParentEquity:            ptrFloat(70),
			TotalEquity:             ptrFloat(100),
		}},
	}, "CNY", "消费/饮料")

	if current.NetCash == nil || *current.NetCash != 80 {
		t.Fatalf("net cash = %v, want 80", current.NetCash)
	}
	if current.ShareholderFCF == nil || *current.ShareholderFCF != 42 {
		t.Fatalf("shareholder FCF = %v, want 42", current.ShareholderFCF)
	}
	if current.MinorityFCFAdjustment == nil || *current.MinorityFCFAdjustment != 18 {
		t.Fatalf("minority adjustment = %v, want 18", current.MinorityFCFAdjustment)
	}
	if current.AdjustedNetCash != nil || current.ExCashPE != nil || current.ExCashPFCF != nil || current.FCFYield != nil {
		t.Fatalf("expected derived valuation fields to be cleared, got %+v", current)
	}
}

func TestHKDividendAmount(t *testing.T) {
	amount, currency := hkDividendAmount("每股派港币0.1元")
	if amount != 0.1 || currency != "HKD" {
		t.Fatalf("hk dividend = %v %s, want 0.1 HKD", amount, currency)
	}

	amount, currency = hkDividendAmount("每股派人民币0.154元")
	if amount != 0.154 || currency != "CNY" {
		t.Fatalf("hk dividend = %v %s, want 0.154 CNY", amount, currency)
	}
}

func TestDividendDistributionsAddSumsSameFiscalYear(t *testing.T) {
	items := dividendDistributions{}
	items.add("2025", 0.1, "HKD")
	items.add("2025", 0.09, "HKD")
	items.add("2025", 0.01, "HKD")

	got := items["2025"]
	if got.PerShare == nil || *got.PerShare != 0.2 {
		t.Fatalf("per share = %v, want 0.2", got.PerShare)
	}
	if got.Currency != "HKD" {
		t.Fatalf("currency = %q, want HKD", got.Currency)
	}
}

func TestAshareBuybackFiscalYearPrefersFinancialReportDate(t *testing.T) {
	year := ashareBuybackFiscalYear(map[string]any{
		"REPORTDATE":     "2025-06-30 00:00:00",
		"FINISHDATE":     "2026-03-31 00:00:00",
		"REPURSTARTDATE": "2025-04-08 00:00:00",
	})
	if year != "2025" {
		t.Fatalf("fiscal year = %q, want 2025", year)
	}
}

func TestHKBuybackCashflowCode(t *testing.T) {
	rows := []map[string]any{
		{"REPORT_DATE": "2025-12-31 00:00:00", "STD_ITEM_CODE": "007008", "AMOUNT": float64(78181000000)},
	}
	byDate := hkCashflowByReportDate(rows)
	buyback := hkCashflowAmount(byDate["2025-12-31"], "007008")
	if buyback == nil || *buyback != 78181000000 {
		t.Fatalf("buyback = %v, want 78181000000", buyback)
	}
}
