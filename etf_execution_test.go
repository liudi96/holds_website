package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestEnsureETFExecutionPlansCreatesFourFrozenRounds(t *testing.T) {
	state := AppState{}
	if !ensureETFExecutionPlans(&state, time.Date(2026, 7, 14, 10, 0, 0, 0, time.Local)) {
		t.Fatal("expected execution plans to be initialized")
	}
	if len(state.ETFExecutionPlans) != 4 {
		t.Fatalf("plan count = %d, want 4", len(state.ETFExecutionPlans))
	}
	index := findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")
	if index < 0 {
		t.Fatal("A500 execution plan missing")
	}
	plan := state.ETFExecutionPlans[index]
	if plan.OpportunityPool != 288795.36 || plan.RoundID != "022434-20260714" {
		t.Fatalf("unexpected A500 frozen round: %+v", plan)
	}
	if len(plan.Stages) != 6 || plan.Stages[0].Key != "dd-7" {
		t.Fatalf("unexpected A500 stages: %+v", plan.Stages)
	}
}

func TestETFExecutionTradePreventsDuplicateAndRebuildsAfterDelete(t *testing.T) {
	state := AppState{FX: map[string]float64{"CNY": 1}, Cash: 100000}
	ensureETFExecutionPlans(&state, time.Now())
	trade := Trade{
		ID: 1001, Date: "2026-07-14", AssetType: assetTypeFund, Symbol: "159352", Name: "中证A500ETF南方",
		Side: "buy", Shares: 7000, Price: 1.2893, CurrentPrice: 1.2893, Currency: "CNY", Reason: "第一批",
		ETFExecution: &ETFExecutionTradeMeta{
			TrackerSymbol: "022434", TacticalSymbol: "159352", RoundID: "022434-20260714", StageKey: "dd-7",
			StageThreshold: 7, Installment: 1, InstallmentCount: 2, SignalAsOf: "2026-07-13",
			RecommendedAmount: 9024.86, StagePlannedAmount: 18049.71,
		},
	}
	if err := validateETFExecutionTrade(&state, trade); err != nil {
		t.Fatalf("first installment rejected: %v", err)
	}
	applyTradeToState(&state, trade)
	recordETFExecutionTrade(&state, trade)
	plan := state.ETFExecutionPlans[findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")]
	stage := plan.Stages[findETFExecutionStageIndex(plan.Stages, "dd-7")]
	if stage.Status != etfExecutionStagePartial || len(stage.Installments) != 1 {
		t.Fatalf("first installment not recorded: %+v", stage)
	}
	if err := validateETFExecutionTrade(&state, trade); err == nil {
		t.Fatal("duplicate installment should be rejected")
	}

	state.Trades = nil
	reverseFundTradeFromState(&state, trade)
	rebuildETFExecutionStagesFromTrades(&state)
	plan = state.ETFExecutionPlans[findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")]
	stage = plan.Stages[findETFExecutionStageIndex(plan.Stages, "dd-7")]
	if stage.Status != etfExecutionStagePending || len(stage.Installments) != 0 {
		t.Fatalf("deleting trade should restore pending installment: %+v", stage)
	}
}

func TestETFExecutionTradeHandlersRecordAndRestoreInstallment(t *testing.T) {
	withTempPortfolioData(t)
	state := AppState{FX: map[string]float64{"CNY": 1}, Cash: 100000}
	ensureETFExecutionPlans(&state, time.Now())
	server := Server{state: state}
	plan := state.ETFExecutionPlans[findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")]
	trade := Trade{
		AssetType: assetTypeFund,
		Symbol:    "159352",
		Name:      "A500 tactical ETF",
		Side:      "buy",
		Shares:    7000,
		Price:     1.2893,
		Currency:  "CNY",
		Reason:    "A500 first tactical installment",
		ETFExecution: &ETFExecutionTradeMeta{
			TrackerSymbol:      "022434",
			TacticalSymbol:     "159352",
			RoundID:            plan.RoundID,
			StageKey:           "dd-7",
			StageThreshold:     7,
			Installment:        1,
			InstallmentCount:   2,
			SignalAsOf:         "2026-07-13",
			RecommendedAmount:  9024.86,
			StagePlannedAmount: 18049.71,
		},
	}
	payload, err := json.Marshal(trade)
	if err != nil {
		t.Fatalf("marshal trade: %v", err)
	}
	createReq := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(payload))
	createRec := httptest.NewRecorder()

	server.handleCreateTrade(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body=%s", createRec.Code, createRec.Body.String())
	}
	if len(server.state.Trades) != 1 {
		t.Fatalf("trade count = %d, want 1", len(server.state.Trades))
	}
	created := server.state.Trades[0]
	activePlan := server.state.ETFExecutionPlans[findETFExecutionPlanIndex(server.state.ETFExecutionPlans, "022434")]
	stage := activePlan.Stages[findETFExecutionStageIndex(activePlan.Stages, "dd-7")]
	if stage.Status != etfExecutionStagePartial || len(stage.Installments) != 1 || len(stage.Installments[0].TradeIDs) != 1 {
		t.Fatalf("first installment was not atomically recorded: %+v", stage)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/trades/"+strconv.FormatInt(created.ID, 10), nil)
	deleteRec := httptest.NewRecorder()
	server.handleDeleteTrade(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusOK {
		t.Fatalf("delete status = %d, body=%s", deleteRec.Code, deleteRec.Body.String())
	}
	if len(server.state.Trades) != 0 {
		t.Fatalf("trade should be deleted: %+v", server.state.Trades)
	}
	activePlan = server.state.ETFExecutionPlans[findETFExecutionPlanIndex(server.state.ETFExecutionPlans, "022434")]
	stage = activePlan.Stages[findETFExecutionStageIndex(activePlan.Stages, "dd-7")]
	if stage.Status != etfExecutionStagePending || len(stage.Installments) != 0 {
		t.Fatalf("deleting the trade should restore the first installment: %+v", stage)
	}
}

func TestCancelETFExecutionStagePersistsReason(t *testing.T) {
	withTempPortfolioData(t)
	state := AppState{FX: map[string]float64{"CNY": 1}}
	ensureETFExecutionPlans(&state, time.Now())
	planIndex := findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")
	if planIndex < 0 || state.ETFExecutionPlans[planIndex].RoundID != "022434-20260714" {
		t.Fatalf("unexpected seeded plan: %+v", state.ETFExecutionPlans)
	}
	cloned, err := cloneAppState(state)
	if err != nil || findETFExecutionPlanIndex(cloned.ETFExecutionPlans, "022434") < 0 {
		t.Fatalf("execution plans lost while cloning: err=%v plans=%+v", err, cloned.ETFExecutionPlans)
	}
	server := Server{state: state}
	payload := `{"trackerSymbol":"022434","roundId":"022434-20260714","stageKey":"dd-7","reason":"回撤收窄3个百分点"}`
	var parsed cancelETFExecutionStageRequest
	if err := json.Unmarshal([]byte(payload), &parsed); err != nil || parsed.RoundID != "022434-20260714" {
		t.Fatalf("invalid cancellation fixture: err=%v payload=%+v", err, parsed)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/etf/execution-stages/cancel", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	server.handleCancelETFExecutionStage(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	plan := server.state.ETFExecutionPlans[findETFExecutionPlanIndex(server.state.ETFExecutionPlans, "022434")]
	stage := plan.Stages[findETFExecutionStageIndex(plan.Stages, "dd-7")]
	if stage.Status != etfExecutionStageCanceled || stage.CancellationReason != "回撤收窄3个百分点" {
		t.Fatalf("stage cancellation not persisted: %+v", stage)
	}
}

func TestETFExecutionNewHighArchivesActiveRound(t *testing.T) {
	state := AppState{FX: map[string]float64{"CNY": 1}}
	ensureETFExecutionPlans(&state, time.Now())
	planIndex := findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")
	plan := &state.ETFExecutionPlans[planIndex]
	plan.Stages[0].Status = etfExecutionStagePartial
	plan.Stages[0].Installments = []ETFExecutionInstallment{{Number: 1, TradeIDs: []int64{1001}, ExecutedAmount: 9024.86}}
	drawdown := 0.0
	peak := 10000.0
	state.ETFRuleStatuses = []ETFRuleStatus{{
		Symbol: "022434",
		Metrics: []ETFRuleMetric{
			{Key: "drawdown3y", Value: &drawdown, AsOf: "2026-07-20", Available: true},
			{Key: "totalReturnPeak", Value: &peak, AsOf: "2026-07-20", Available: true},
		},
	}}

	syncETFExecutionPlans(&state, time.Date(2026, 7, 20, 16, 0, 0, 0, time.Local))

	plan = &state.ETFExecutionPlans[planIndex]
	if plan.RoundID != "022434-20260720" || plan.StartedAt != "2026-07-20" {
		t.Fatalf("new high did not create a new round: %+v", plan)
	}
	if len(plan.ArchivedRounds) != 1 || plan.ArchivedRounds[0].RoundID != "022434-20260714" {
		t.Fatalf("old round was not archived: %+v", plan.ArchivedRounds)
	}
	if plan.Stages[0].Status != etfExecutionStagePending || len(plan.Stages[0].Installments) != 0 {
		t.Fatalf("new round should start clean: %+v", plan.Stages[0])
	}
}

func TestETFExecutionOverrideRequiresAndStoresReason(t *testing.T) {
	state := AppState{FX: map[string]float64{"CNY": 1}}
	ensureETFExecutionPlans(&state, time.Now())
	plan := state.ETFExecutionPlans[findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")]
	trade := Trade{
		ID: 1002, Date: "2026-07-14", AssetType: assetTypeFund, Symbol: "159352", Side: "buy", Shares: 1000, Price: 1.3,
		ETFExecution: &ETFExecutionTradeMeta{
			TrackerSymbol: "022434", TacticalSymbol: "159352", RoundID: plan.RoundID, StageKey: "dd-7",
			StageThreshold: 7, Installment: 1, InstallmentCount: 2, RecommendedAmount: 9024.86, StagePlannedAmount: 18049.71,
			Override: true,
		},
	}
	if err := validateETFExecutionTrade(&state, trade); err == nil {
		t.Fatal("override without a reason should be rejected")
	}
	trade.ETFExecution.OverrideReason = "manual allocation correction"
	if err := validateETFExecutionTrade(&state, trade); err != nil {
		t.Fatalf("override with a reason rejected: %v", err)
	}
	recordETFExecutionTrade(&state, trade)
	plan = state.ETFExecutionPlans[findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")]
	stage := plan.Stages[findETFExecutionStageIndex(plan.Stages, "dd-7")]
	if len(stage.Installments) != 1 || !stage.Installments[0].Override || stage.Installments[0].OverrideReason == "" {
		t.Fatalf("override audit metadata was not stored: %+v", stage)
	}
}

func TestETFExecutionCreateRollbackWhenPersistenceUnavailable(t *testing.T) {
	withTempPortfolioData(t)
	state := AppState{FX: map[string]float64{"CNY": 1}, Cash: 100000}
	ensureETFExecutionPlans(&state, time.Now())
	server := Server{state: state}
	plan := state.ETFExecutionPlans[findETFExecutionPlanIndex(state.ETFExecutionPlans, "022434")]
	trade := Trade{
		AssetType: assetTypeFund, Symbol: "159352", Side: "buy", Shares: 7000, Price: 1.2893, CurrentPrice: 1.2893,
		Currency: "CNY", Reason: "persistence rollback test",
		ETFExecution: &ETFExecutionTradeMeta{
			TrackerSymbol: "022434", TacticalSymbol: "159352", RoundID: plan.RoundID, StageKey: "dd-7",
			StageThreshold: 7, Installment: 1, InstallmentCount: 2, RecommendedAmount: 9024.86, StagePlannedAmount: 18049.71,
		},
	}
	payload, err := json.Marshal(trade)
	if err != nil {
		t.Fatalf("marshal trade: %v", err)
	}
	dataFile = filepath.Join(t.TempDir(), "missing", "portfolio.json")
	req := httptest.NewRequest(http.MethodPost, "/api/trades", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	server.handleCreateTrade(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", rec.Code, rec.Body.String())
	}
	if len(server.state.Trades) != 0 || server.state.Cash != 100000 {
		t.Fatalf("server state changed after failed persistence: %+v", server.state)
	}
	activePlan := server.state.ETFExecutionPlans[findETFExecutionPlanIndex(server.state.ETFExecutionPlans, "022434")]
	stage := activePlan.Stages[findETFExecutionStageIndex(activePlan.Stages, "dd-7")]
	if stage.Status != etfExecutionStagePending || len(stage.Installments) != 0 {
		t.Fatalf("execution stage changed after failed persistence: %+v", stage)
	}
}
