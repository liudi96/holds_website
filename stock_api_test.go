package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleUpdateScreeningWeightsRejectsInvalidSum(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{ScreeningWeights: DefaultScreeningWeights()}}
	req := httptest.NewRequest(http.MethodPut, "/api/screening-weights", strings.NewReader(`{"quality":40,"cashFlow":25,"valuation":20,"shareholderReturn":15,"growth":10}`))
	rec := httptest.NewRecorder()

	server.handleUpdateScreeningWeights(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandleUpdateScreeningWeightsSavesValidWeights(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{ScreeningWeights: DefaultScreeningWeights()}}
	req := httptest.NewRequest(http.MethodPut, "/api/screening-weights", strings.NewReader(`{"quality":25,"cashFlow":25,"valuation":25,"shareholderReturn":15,"growth":10}`))
	rec := httptest.NewRecorder()

	server.handleUpdateScreeningWeights(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if server.state.ScreeningWeights.Quality != 25 || server.state.ScreeningWeights.Valuation != 25 {
		t.Fatalf("weights not saved: %+v", server.state.ScreeningWeights)
	}
}

func TestHandleUpsertStockCreatesUnifiedStock(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{FX: map[string]float64{"CNY": 1}}}
	req := httptest.NewRequest(http.MethodPost, "/api/stocks", strings.NewReader(`{"symbol":"600519.SH","name":"贵州茅台","industry":"白酒","currency":"CNY","qualityScore":92}`))
	rec := httptest.NewRecorder()

	server.handleUpsertStock(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	stock := findStock(server.state.Stocks, "600519.SH")
	if stock == nil || stock.Name != "贵州茅台" || stock.QualityScore == nil || *stock.QualityScore != 92 {
		t.Fatalf("stock not upserted: %+v", server.state.Stocks)
	}
}

func TestHandleUpsertStockUpdatesExistingStockFields(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{
		FX: map[string]float64{"CNY": 1},
		Stocks: []Stock{{
			Symbol:       "600519.SH",
			Name:         "贵州茅台",
			Industry:     "白酒",
			Notes:        "旧备注",
			QualityScore: ptrFloat(80),
		}},
	}}
	req := httptest.NewRequest(http.MethodPut, "/api/stocks/600519.SH", strings.NewReader(`{"symbol":"600519.SH","name":"贵州茅台","industry":"高端白酒","notes":"新备注","qualityScore":92}`))
	rec := httptest.NewRecorder()

	server.handleUpsertStock(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	stock := findStock(server.state.Stocks, "600519.SH")
	if stock == nil {
		t.Fatal("missing stock after update")
	}
	if stock.Industry != "高端白酒" || stock.Notes != "新备注" || stock.QualityScore == nil || *stock.QualityScore != 92 {
		t.Fatalf("stock fields not updated: %+v", stock)
	}
}

func TestHandleUpsertStockPersistsEditableHumanInputs(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{
		FX: map[string]float64{"HKD": 0.9},
		Stocks: []Stock{{
			Symbol:       "0700.HK",
			Name:         "腾讯控股",
			BuyLogic:     "旧买入逻辑",
			Notes:        "旧备注",
			QualityScore: ptrFloat(88),
		}},
	}}
	body := `{
		"symbol":"0700.HK",
		"name":"腾讯控股",
		"buyLogic":"自由现金流稳定，广告和视频号仍有经营杠杆",
		"notes":"",
		"qualityScore":null,
		"valuation":{
			"currency":"HKD",
			"currentPrice":460,
			"requiredMargin":0.18,
			"scenarios":[
				{"name":"保守","fcf":1500,"reasonablePfcf":12,"shares":10},
				{"name":"基准","fcf":1800,"reasonablePfcf":15,"shares":10},
				{"name":"乐观","fcf":2100,"reasonablePfcf":18,"shares":10}
			],
			"range":{"low":1800,"base":2700,"high":3780,"currency":"HKD","marginOfSafety":0.83}
		}
	}`
	req := httptest.NewRequest(http.MethodPut, "/api/stocks/0700.HK", strings.NewReader(body))
	rec := httptest.NewRecorder()

	server.handleUpsertStock(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	stock := findStock(server.state.Stocks, "0700.HK")
	if stock == nil {
		t.Fatal("missing stock after update")
	}
	if stock.BuyLogic != "自由现金流稳定，广告和视频号仍有经营杠杆" {
		t.Fatalf("buy logic not saved: %+v", stock)
	}
	if stock.Notes != "" {
		t.Fatalf("notes should be cleared, got %q", stock.Notes)
	}
	if stock.QualityScore != nil {
		t.Fatalf("quality score should be cleared, got %v", *stock.QualityScore)
	}
	if stock.Valuation == nil || stock.Valuation.RequiredMargin == nil || *stock.Valuation.RequiredMargin != 0.18 || stock.Valuation.Range == nil || stock.Valuation.Range.Base != 2700 || len(stock.Valuation.Scenarios) != 3 {
		t.Fatalf("valuation assumptions not saved: %+v", stock.Valuation)
	}
}

func TestHandleDeleteStockRejectsPositionedStock(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{
		FX: map[string]float64{"HKD": 0.9},
		Stocks: []Stock{{
			Symbol:   "0700.HK",
			Name:     "腾讯控股",
			Currency: "HKD",
			Position: &StockPosition{Shares: 100, Cost: 420},
		}},
	}}
	req := httptest.NewRequest(http.MethodDelete, "/api/stocks/0700.HK", nil)
	rec := httptest.NewRecorder()

	server.handleDeleteStock(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	stock := findStock(server.state.Stocks, "0700.HK")
	if stock == nil || stock.Position == nil || stock.Position.Shares != 100 {
		t.Fatalf("positioned stock should remain: %+v", server.state.Stocks)
	}
}

func TestHandleCreateDecisionLogAddsEvidenceSnapshot(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{
		FX: map[string]float64{"HKD": 0.9},
		Stocks: []Stock{{
			Symbol:         "0700.HK",
			Name:           "腾讯控股",
			CurrentPrice:   460,
			Currency:       "HKD",
			FairValueRange: "HK$480-590",
			MarginOfSafety: ptrFloat(0.12),
			Position:       &StockPosition{Shares: 200, Cost: 480},
		}},
	}}
	normalizePortfolioState(&server.state)
	req := httptest.NewRequest(http.MethodPost, "/api/decision-logs", strings.NewReader(`{"type":"hold","symbol":"0700.HK","decision":"继续持有","detail":"FCF 没有恶化，安全边际不足所以不加仓"}`))
	rec := httptest.NewRecorder()

	server.handleCreateDecisionLog(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", rec.Code, rec.Body.String())
	}
	if len(server.state.DecisionLogs) != 1 {
		t.Fatalf("decision log count = %d, want 1", len(server.state.DecisionLogs))
	}
	detail := server.state.DecisionLogs[0].Detail
	if !strings.Contains(detail, "FCF 没有恶化") || !strings.Contains(detail, "价格 HKD 460.0000") || !strings.Contains(detail, "估值区间 HK$480-590") {
		t.Fatalf("missing evidence snapshot: %q", detail)
	}
}

func TestHandleDeleteTradeRemovesTradeAndReversesBuy(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{
		FX:   map[string]float64{"HKD": 0.9},
		Cash: 96220,
		Trades: []Trade{{
			ID:           123,
			Date:         "2026-05-26",
			AssetType:    "stock",
			Symbol:       "0700.HK",
			Name:         "腾讯控股",
			Side:         "buy",
			Shares:       10,
			Price:        420,
			CurrentPrice: 430,
			Currency:     "HKD",
			Reason:       "误新增",
		}},
		DecisionLogs: []DecisionLog{{
			ID:       456,
			Date:     "2026-05-26",
			Type:     "trade",
			Symbol:   "0700.HK",
			Name:     "腾讯控股",
			Decision: "买入 腾讯控股",
			Detail:   "买入 10 股；成交价 HKD 420.0000；录入最新价 HKD 430.0000；理由：误新增",
		}},
		Stocks: []Stock{{
			Symbol:       "0700.HK",
			Name:         "腾讯控股",
			CurrentPrice: 430,
			Currency:     "HKD",
			Position:     &StockPosition{Shares: 10, Cost: 420},
		}},
	}}
	normalizePortfolioState(&server.state)
	req := httptest.NewRequest(http.MethodDelete, "/api/trades/123", nil)
	rec := httptest.NewRecorder()

	server.handleDeleteTrade(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if len(server.state.Trades) != 0 {
		t.Fatalf("trade should be deleted: %+v", server.state.Trades)
	}
	if len(server.state.DecisionLogs) != 0 {
		t.Fatalf("generated trade log should be deleted: %+v", server.state.DecisionLogs)
	}
	if server.state.Cash != 100000 {
		t.Fatalf("cash = %.2f, want 100000", server.state.Cash)
	}
	stock := findStock(server.state.Stocks, "0700.HK")
	if stock == nil {
		t.Fatal("stock should remain in tracking pool")
	}
	if stock.Position != nil && stock.Position.Shares > 0 {
		t.Fatalf("position should be reversed: %+v", stock.Position)
	}
}

func TestHandleUpdateValuationHistoryWritesRuntimeBook(t *testing.T) {
	withTempPortfolioData(t)
	oldFetcher := fetchValuationMonthlyCloses
	fetchValuationMonthlyCloses = func(_ *http.Client, _ string) ([]dailyClose, error) {
		return []dailyClose{
			{Date: "2016-01-31", Price: 100},
			{Date: "2016-02-29", Price: 120},
		}, nil
	}
	defer func() {
		fetchValuationMonthlyCloses = oldFetcher
	}()

	server := Server{state: AppState{Stocks: []Stock{{
		Symbol:       "0700.HK",
		Name:         "腾讯控股",
		CurrentPrice: 100,
		Financials: &Financials{Valuation: &FinancialValuation{
			Price: ptrFloat(100),
			PE:    ptrFloat(20),
			PB:    ptrFloat(4),
		}},
	}}}}
	req := httptest.NewRequest(http.MethodPost, "/api/valuation-history/update", nil)
	rec := httptest.NewRecorder()

	server.handleUpdateValuationHistory(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var response ValuationHistoryUpdateResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Updated != 1 {
		t.Fatalf("updated = %d, want 1", response.Updated)
	}
	stock := findStock(response.State.Stocks, "0700.HK")
	if stock == nil || stock.Financials == nil || stock.Financials.Valuation == nil || stock.Financials.Valuation.PEPercentile == nil {
		t.Fatalf("state missing valuation percentiles: %+v", response.State.Stocks)
	}
	if *stock.Financials.Valuation.PEPercentile != 0.5 {
		t.Fatalf("PE percentile = %v, want 0.5", *stock.Financials.Valuation.PEPercentile)
	}
	var book ValuationHistoryBook
	if err := json.Unmarshal(readTestFile(t, runtimeValuationHistoryFile), &book); err != nil {
		t.Fatalf("decode valuation history: %v", err)
	}
	points := book.History["0700.HK"]
	if len(points) != 2 {
		t.Fatalf("history = %+v, want monthly Tencent points", book.History)
	}
	if points[0].Price == nil || *points[0].Price != 100 || points[1].PE == nil || *points[1].PE != 24 {
		t.Fatalf("unexpected valuation points: %+v", points)
	}
	if book.Percentiles["0700.HK"].PE == nil || *book.Percentiles["0700.HK"].PE != 0.5 {
		t.Fatalf("missing PE percentile: %+v", book.Percentiles)
	}
}

func readTestFile(t *testing.T, path string) []byte {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return body
}

func withTempPortfolioData(t *testing.T) {
	t.Helper()
	oldDataDir := dataDir
	oldDataFile := dataFile
	oldRuntimeQuotesFile := runtimeQuotesFile
	oldRuntimeIndustryMetricsFile := runtimeIndustryMetricsFile
	oldRuntimeValuationHistoryFile := runtimeValuationHistoryFile
	dir := t.TempDir()
	dataDir = dir
	dataFile = filepath.Join(dir, "portfolio.json")
	runtimeQuotesFile = filepath.Join(dir, "runtime", "quotes.json")
	runtimeIndustryMetricsFile = filepath.Join(dir, "runtime", "industry_metrics.json")
	runtimeValuationHistoryFile = filepath.Join(dir, "runtime", "valuation_history.json")
	if err := os.WriteFile(dataFile, []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("seed temp portfolio: %v", err)
	}
	t.Cleanup(func() {
		dataDir = oldDataDir
		dataFile = oldDataFile
		runtimeQuotesFile = oldRuntimeQuotesFile
		runtimeIndustryMetricsFile = oldRuntimeIndustryMetricsFile
		runtimeValuationHistoryFile = oldRuntimeValuationHistoryFile
	})
}
