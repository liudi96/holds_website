package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNormalizePortfolioStateMergesLegacyHoldingsAndCandidatesIntoStocks(t *testing.T) {
	state := AppState{
		FX: map[string]float64{"CNY": 1, "HKD": 0.9},
		Holdings: []Holding{{
			Symbol:       "0700.HK",
			Name:         "腾讯控股",
			Shares:       200,
			Cost:         480,
			CurrentPrice: 460,
			Currency:     "HKD",
			Industry:     "互联网平台",
			Action:       "继续持有",
		}},
		Candidates: []Candidate{{
			Symbol:         "0700.HK",
			Name:           "腾讯控股",
			QualityScore:   ptrFloat(89),
			FairValueRange: "HK$480-590",
			Action:         "等待更好价格",
		}, {
			Symbol:       "600519.SH",
			Name:         "贵州茅台",
			CurrentPrice: 1500,
			Currency:     "CNY",
			Industry:     "白酒",
		}},
	}

	normalizePortfolioState(&state)

	if len(state.Stocks) != 2 {
		t.Fatalf("stocks = %d, want 2: %+v", len(state.Stocks), state.Stocks)
	}
	tencent := findStock(state.Stocks, "0700.HK")
	if tencent == nil {
		t.Fatal("expected Tencent stock")
	}
	if tencent.Position == nil {
		t.Fatalf("expected holding position, got %+v", tencent)
	}
	if tencent.Position.Shares != 200 || tencent.Position.Cost != 480 {
		t.Fatalf("unexpected position: %+v", tencent.Position)
	}
	if tencent.QualityScore == nil || *tencent.QualityScore != 89 {
		t.Fatalf("expected candidate quality to be merged, got %+v", tencent.QualityScore)
	}
	if tencent.Action != "继续持有" {
		t.Fatalf("holding action should win, got %q", tencent.Action)
	}
	if findStock(state.Stocks, "600519.SH") == nil {
		t.Fatal("expected candidate-only stock")
	}
}

func TestAppStateJSONExposesStocksAndHidesLegacyAssetBuckets(t *testing.T) {
	state := AppState{
		Stocks: []Stock{{
			Symbol: "0700.HK",
			Name:   "腾讯控股",
			Position: &StockPosition{
				Shares: 200,
				Cost:   480,
			},
		}},
		Holdings: []Holding{{Symbol: "0700.HK"}},
	}

	body, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("marshal AppState: %v", err)
	}
	text := string(body)
	if !strings.Contains(text, `"stocks"`) {
		t.Fatalf("expected stocks in json, got %s", text)
	}
	if strings.Contains(text, `"holdings"`) || strings.Contains(text, `"candidates"`) || strings.Contains(text, `"funds"`) {
		t.Fatalf("legacy buckets leaked into json: %s", text)
	}
}

func TestScreeningWeightsDefaultAndValidation(t *testing.T) {
	weights := DefaultScreeningWeights()
	if err := weights.Validate(); err != nil {
		t.Fatalf("default weights invalid: %v", err)
	}
	if weights.Quality != 30 || weights.CashFlow != 25 || weights.Valuation != 20 || weights.ShareholderReturn != 15 || weights.Growth != 10 {
		t.Fatalf("unexpected default weights: %+v", weights)
	}

	invalid := ScreeningWeights{Quality: 40, CashFlow: 25, Valuation: 20, ShareholderReturn: 15, Growth: 10}
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected invalid weight sum to be rejected")
	}
}

func TestTradeValidationAcceptsFundsAndRejectsMissingReason(t *testing.T) {
	fundTrade := Trade{AssetType: "fund", Symbol: "004814.OF", Name: "中欧红利优享混合", Side: "buy", Shares: 100, Price: 1.2, CurrentPrice: 1.2, Currency: "CNY", Reason: "配置基金"}
	if err := validateTrade(fundTrade); err != nil {
		t.Fatalf("fund trade should be accepted: %v", err)
	}

	stockTrade := Trade{AssetType: "stock", Symbol: "0700.HK", Name: "腾讯控股", Side: "buy", Shares: 100, Price: 420, CurrentPrice: 420, Currency: "HKD"}
	if err := validateTrade(stockTrade); err == nil {
		t.Fatal("expected missing reason to be rejected")
	}

	stockTrade.Reason = "安全边际达标，且 Q1 后 FCF 假设未恶化"
	if err := validateTrade(stockTrade); err != nil {
		t.Fatalf("stock trade with reason should pass: %v", err)
	}
}

func TestResolveStockTradeInputParsesCodeAndName(t *testing.T) {
	server := Server{state: AppState{FX: map[string]float64{"HKD": 0.9}}}
	trade := Trade{
		AssetType:    "stock",
		Name:         "9926.HK 康方生物",
		Side:         "buy",
		Shares:       100,
		Price:        45.2,
		CurrentPrice: 45.2,
		Reason:       "新开仓，先按小仓位跟踪执行",
	}

	if err := server.resolveTradeInput(&trade); err != nil {
		t.Fatalf("resolveTradeInput() error = %v", err)
	}
	if trade.Symbol != "9926.HK" || trade.Name != "康方生物" || trade.Currency != "HKD" {
		t.Fatalf("trade = %+v, want parsed stock code, name and HKD currency", trade)
	}
}

func TestResolveStockTradeInputRejectsNameOnlyNewBuy(t *testing.T) {
	server := Server{state: AppState{FX: map[string]float64{"HKD": 0.9}}}
	trade := Trade{
		AssetType:    "stock",
		Name:         "康方生物",
		Side:         "buy",
		Shares:       100,
		Price:        45.2,
		CurrentPrice: 45.2,
		Reason:       "新开仓，先按小仓位跟踪执行",
	}

	if err := server.resolveTradeInput(&trade); err == nil || !strings.Contains(err.Error(), "请输入股票代码") {
		t.Fatalf("resolveTradeInput() error = %v, want stock-code requirement", err)
	}
}

func TestApplyTradeToStateUpdatesUnifiedStocks(t *testing.T) {
	state := AppState{
		FX:   map[string]float64{"CNY": 1},
		Cash: 100000,
		Stocks: []Stock{{
			Symbol:       "600519.SH",
			Name:         "贵州茅台",
			Currency:     "CNY",
			CurrentPrice: 1500,
			Industry:     "白酒",
		}},
	}
	normalizePortfolioState(&state)
	trade := Trade{
		AssetType:    "stock",
		Symbol:       "600519.SH",
		Name:         "贵州茅台",
		Side:         "buy",
		Shares:       10,
		Price:        1400,
		CurrentPrice: 1500,
		Currency:     "CNY",
		Reason:       "安全边际达标，买入后仓位仍受控",
	}

	applyTradeToState(&state, trade)

	stock := findStock(state.Stocks, "600519.SH")
	if stock == nil {
		t.Fatal("missing traded stock")
	}
	if stock.Position == nil || stock.Position.Shares != 10 || stock.Position.Cost != 1400 {
		t.Fatalf("trade did not update unified position: %+v", stock.Position)
	}
	if len(state.Trades) != 1 {
		t.Fatalf("trade count = %d, want 1", len(state.Trades))
	}
}

func TestApplyStockSellRemovesHoldingWithoutCandidate(t *testing.T) {
	state := AppState{
		FX:   map[string]float64{"HKD": 0.9},
		Cash: 1000,
		Holdings: []Holding{{
			Symbol:       "9926.HK",
			Name:         "康方生物",
			Shares:       100,
			Cost:         45,
			CurrentPrice: 48,
			Currency:     "HKD",
		}},
		Candidates: []Candidate{{Symbol: "0700.HK", Name: "腾讯控股"}},
	}

	applyTradeToState(&state, Trade{
		AssetType:    "stock",
		Symbol:       "9926.HK",
		Name:         "康方生物",
		Side:         "sell",
		Shares:       100,
		Price:        48,
		CurrentPrice: 48,
		Currency:     "HKD",
		Reason:       "卖出后不再进入跟踪池",
	})

	if findHoldingIndex(state.Holdings, "9926.HK") != -1 {
		t.Fatalf("sold-out holding should be removed: %+v", state.Holdings)
	}
	if len(state.Candidates) != 1 || state.Candidates[0].Symbol != "0700.HK" {
		t.Fatalf("sell should not add candidate tracking entries: %+v", state.Candidates)
	}
}

func TestValuationRangeUsesScenarioAssumptions(t *testing.T) {
	assumptions := ValuationAssumptions{
		Currency:     "CNY",
		CurrentPrice: 90,
		Scenarios: []ValuationScenario{
			{Name: "bear", RevenueGrowth: 0.02, ProfitMargin: 0.12, FCF: 800, ReasonablePE: 10, ReasonablePFCF: 9, Shares: 100},
			{Name: "base", RevenueGrowth: 0.05, ProfitMargin: 0.15, FCF: 1000, ReasonablePE: 12, ReasonablePFCF: 11, Shares: 100},
			{Name: "bull", RevenueGrowth: 0.08, ProfitMargin: 0.18, FCF: 1200, ReasonablePE: 15, ReasonablePFCF: 14, Shares: 100},
		},
	}

	result, err := CalculateValuationRange(assumptions)
	if err != nil {
		t.Fatalf("CalculateValuationRange error: %v", err)
	}
	if result.Low <= 0 || result.Base <= result.Low || result.High <= result.Base {
		t.Fatalf("unexpected valuation range: %+v", result)
	}
	if result.MarginOfSafety == nil || *result.MarginOfSafety <= 0 {
		t.Fatalf("expected positive margin of safety, got %+v", result.MarginOfSafety)
	}
}

func TestValuationPercentileUsesHistoricalPoints(t *testing.T) {
	points := []ValuationHistoryPoint{
		{Date: "2024-01-31", PE: ptrFloat(10), PB: ptrFloat(1)},
		{Date: "2024-02-29", PE: ptrFloat(15), PB: ptrFloat(1.5)},
		{Date: "2024-03-31", PE: ptrFloat(20), PB: ptrFloat(2)},
		{Date: "2024-04-30", PE: ptrFloat(25), PB: ptrFloat(2.5)},
	}

	pe, pb := ValuationPercentiles(points, 20, 1.5)
	if pe == nil || *pe != 0.75 {
		t.Fatalf("PE percentile = %v, want 0.75", pe)
	}
	if pb == nil || *pb != 0.5 {
		t.Fatalf("PB percentile = %v, want 0.5", pb)
	}
}
