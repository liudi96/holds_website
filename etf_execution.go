package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	etfExecutionPlanActive    = "active"
	etfExecutionStagePending  = "pending"
	etfExecutionStagePartial  = "partial"
	etfExecutionStageComplete = "complete"
	etfExecutionStageCanceled = "canceled"
)

type ETFExecutionPlan struct {
	TrackerSymbol   string                `json:"trackerSymbol"`
	TacticalSymbol  string                `json:"tacticalSymbol"`
	TacticalName    string                `json:"tacticalName"`
	RoundID         string                `json:"roundId"`
	StartedAt       string                `json:"startedAt"`
	Status          string                `json:"status"`
	OpportunityPool float64               `json:"opportunityPool"`
	PeakValue       float64               `json:"peakValue,omitempty"`
	PeakAsOf        string                `json:"peakAsOf,omitempty"`
	Stages          []ETFExecutionStage   `json:"stages,omitempty"`
	ArchivedRounds  []ETFExecutionArchive `json:"archivedRounds,omitempty"`
	UpdatedAt       string                `json:"updatedAt,omitempty"`
}

type ETFExecutionArchive struct {
	RoundID         string              `json:"roundId"`
	StartedAt       string              `json:"startedAt"`
	ClosedAt        string              `json:"closedAt"`
	OpportunityPool float64             `json:"opportunityPool"`
	PeakValue       float64             `json:"peakValue,omitempty"`
	PeakAsOf        string              `json:"peakAsOf,omitempty"`
	Stages          []ETFExecutionStage `json:"stages,omitempty"`
}

type ETFExecutionStage struct {
	Key                string                    `json:"key"`
	Threshold          float64                   `json:"threshold"`
	PlannedAmount      float64                   `json:"plannedAmount,omitempty"`
	TriggeredAt        string                    `json:"triggeredAt,omitempty"`
	SignalAsOf         string                    `json:"signalAsOf,omitempty"`
	Status             string                    `json:"status"`
	InstallmentCount   int                       `json:"installmentCount,omitempty"`
	Installments       []ETFExecutionInstallment `json:"installments,omitempty"`
	CanceledAt         string                    `json:"canceledAt,omitempty"`
	CancellationReason string                    `json:"cancellationReason,omitempty"`
}

type ETFExecutionInstallment struct {
	Number         int     `json:"number"`
	PlannedAmount  float64 `json:"plannedAmount,omitempty"`
	ExecutedAmount float64 `json:"executedAmount,omitempty"`
	ExecutedAt     string  `json:"executedAt,omitempty"`
	TradeIDs       []int64 `json:"tradeIds,omitempty"`
	Override       bool    `json:"override,omitempty"`
	OverrideReason string  `json:"overrideReason,omitempty"`
}

type ETFExecutionTradeMeta struct {
	TrackerSymbol      string  `json:"trackerSymbol"`
	TacticalSymbol     string  `json:"tacticalSymbol"`
	RoundID            string  `json:"roundId"`
	StageKey           string  `json:"stageKey"`
	StageThreshold     float64 `json:"stageThreshold"`
	Installment        int     `json:"installment"`
	InstallmentCount   int     `json:"installmentCount"`
	SignalAsOf         string  `json:"signalAsOf"`
	RecommendedAmount  float64 `json:"recommendedAmount"`
	StagePlannedAmount float64 `json:"stagePlannedAmount"`
	Override           bool    `json:"override,omitempty"`
	OverrideReason     string  `json:"overrideReason,omitempty"`
}

type etfExecutionPlanConfig struct {
	TrackerSymbol   string
	TacticalSymbol  string
	TacticalName    string
	StartedAt       string
	OpportunityPool float64
	Thresholds      []float64
}

var etfExecutionPlanConfigs = []etfExecutionPlanConfig{
	{TrackerSymbol: "022434", TacticalSymbol: "159352", TacticalName: "中证A500ETF南方", StartedAt: "2026-07-14", OpportunityPool: 288795.36, Thresholds: []float64{7, 12, 18, 25, 35, 45}},
	{TrackerSymbol: "008163", TacticalSymbol: "515450", TacticalName: "红利低波50ETF南方", StartedAt: "2026-07-13", OpportunityPool: 252287.81, Thresholds: []float64{4, 6, 9, 12, 15, 20}},
	{TrackerSymbol: "018738", TacticalSymbol: "513650", TacticalName: "南方标普500ETF", StartedAt: "2026-07-14", OpportunityPool: 288993.85, Thresholds: []float64{8, 12, 18, 25, 35}},
	{TrackerSymbol: "021000", TacticalSymbol: "159659", TacticalName: "招商纳斯达克100ETF", StartedAt: "2026-07-14", OpportunityPool: 288996.89, Thresholds: []float64{10, 15, 20, 30, 40}},
}

type cancelETFExecutionStageRequest struct {
	TrackerSymbol string `json:"trackerSymbol"`
	RoundID       string `json:"roundId"`
	StageKey      string `json:"stageKey"`
	Reason        string `json:"reason"`
}

func etfExecutionConfig(trackerSymbol string) (etfExecutionPlanConfig, bool) {
	normalized := normalizeFundSymbol(trackerSymbol)
	for _, config := range etfExecutionPlanConfigs {
		if normalizeFundSymbol(config.TrackerSymbol) == normalized {
			return config, true
		}
	}
	return etfExecutionPlanConfig{}, false
}

func ensureETFExecutionPlans(state *AppState, now time.Time) bool {
	if state == nil {
		return false
	}
	changed := false
	for _, config := range etfExecutionPlanConfigs {
		index := findETFExecutionPlanIndex(state.ETFExecutionPlans, config.TrackerSymbol)
		if index < 0 {
			state.ETFExecutionPlans = append(state.ETFExecutionPlans, newETFExecutionPlan(config, now))
			changed = true
			continue
		}
		plan := &state.ETFExecutionPlans[index]
		if strings.TrimSpace(plan.TacticalSymbol) == "" {
			plan.TacticalSymbol = config.TacticalSymbol
			changed = true
		}
		if strings.TrimSpace(plan.TacticalName) == "" {
			plan.TacticalName = config.TacticalName
			changed = true
		}
		if plan.OpportunityPool <= 0 {
			plan.OpportunityPool = config.OpportunityPool
			changed = true
		}
		if strings.TrimSpace(plan.RoundID) == "" {
			plan.RoundID = etfExecutionRoundID(config.TrackerSymbol, firstNonEmpty(plan.StartedAt, config.StartedAt))
			changed = true
		}
		if strings.TrimSpace(plan.StartedAt) == "" {
			plan.StartedAt = config.StartedAt
			changed = true
		}
		if strings.TrimSpace(plan.Status) == "" {
			plan.Status = etfExecutionPlanActive
			changed = true
		}
		if ensureETFExecutionStageSkeletons(plan, config) {
			changed = true
		}
	}
	return changed
}

func newETFExecutionPlan(config etfExecutionPlanConfig, now time.Time) ETFExecutionPlan {
	plan := ETFExecutionPlan{
		TrackerSymbol:   config.TrackerSymbol,
		TacticalSymbol:  config.TacticalSymbol,
		TacticalName:    config.TacticalName,
		RoundID:         etfExecutionRoundID(config.TrackerSymbol, config.StartedAt),
		StartedAt:       config.StartedAt,
		Status:          etfExecutionPlanActive,
		OpportunityPool: config.OpportunityPool,
		UpdatedAt:       now.Format("2006-01-02 15:04:05"),
	}
	ensureETFExecutionStageSkeletons(&plan, config)
	return plan
}

func etfExecutionRoundID(symbol string, startedAt string) string {
	date := strings.ReplaceAll(firstNonEmpty(startedAt, time.Now().Format("2006-01-02")), "-", "")
	return normalizeFundSymbol(symbol) + "-" + date
}

func ensureETFExecutionStageSkeletons(plan *ETFExecutionPlan, config etfExecutionPlanConfig) bool {
	changed := false
	for _, threshold := range config.Thresholds {
		key := etfExecutionStageKey(threshold)
		if findETFExecutionStageIndex(plan.Stages, key) >= 0 {
			continue
		}
		plan.Stages = append(plan.Stages, ETFExecutionStage{Key: key, Threshold: threshold, Status: etfExecutionStagePending})
		changed = true
	}
	return changed
}

func etfExecutionStageKey(threshold float64) string {
	return fmt.Sprintf("dd-%g", threshold)
}

func findETFExecutionPlanIndex(plans []ETFExecutionPlan, trackerSymbol string) int {
	for i := range plans {
		if normalizeFundSymbol(plans[i].TrackerSymbol) == normalizeFundSymbol(trackerSymbol) {
			return i
		}
	}
	return -1
}

func findETFExecutionStageIndex(stages []ETFExecutionStage, stageKey string) int {
	for i := range stages {
		if strings.EqualFold(strings.TrimSpace(stages[i].Key), strings.TrimSpace(stageKey)) {
			return i
		}
	}
	return -1
}

func syncETFExecutionPlans(state *AppState, now time.Time) {
	if state == nil {
		return
	}
	ensureETFExecutionPlans(state, now)
	for i := range state.ETFExecutionPlans {
		plan := &state.ETFExecutionPlans[i]
		status := findETFRuleStatus(state.ETFRuleStatuses, plan.TrackerSymbol)
		if status == nil {
			continue
		}
		peakMetric := findETFRuleMetric(status.Metrics, "totalReturnPeak")
		if peakMetric != nil && peakMetric.Available && peakMetric.Value != nil && *peakMetric.Value > 0 {
			plan.PeakValue = *peakMetric.Value
			plan.PeakAsOf = peakMetric.AsOf
		}
		drawdown := findETFRuleMetric(status.Metrics, "drawdown3y")
		if drawdown == nil || !drawdown.Available || drawdown.Value == nil || *drawdown.Value > 0.0001 || !etfExecutionPlanHasActivity(*plan) {
			continue
		}
		signalDate := firstNonEmpty(drawdown.AsOf, status.AsOf)
		if signalDate == "" || signalDate <= plan.StartedAt {
			continue
		}
		archive := ETFExecutionArchive{
			RoundID:         plan.RoundID,
			StartedAt:       plan.StartedAt,
			ClosedAt:        signalDate,
			OpportunityPool: plan.OpportunityPool,
			PeakValue:       plan.PeakValue,
			PeakAsOf:        plan.PeakAsOf,
			Stages:          append([]ETFExecutionStage(nil), plan.Stages...),
		}
		config, ok := etfExecutionConfig(plan.TrackerSymbol)
		if !ok {
			continue
		}
		archives := append(plan.ArchivedRounds, archive)
		*plan = newETFExecutionPlan(config, now)
		plan.StartedAt = signalDate
		plan.RoundID = etfExecutionRoundID(config.TrackerSymbol, signalDate)
		plan.ArchivedRounds = archives
	}
	rebuildETFExecutionStagesFromTrades(state)
}

func etfExecutionPlanHasActivity(plan ETFExecutionPlan) bool {
	for _, stage := range plan.Stages {
		if stage.Status != etfExecutionStagePending || len(stage.Installments) > 0 {
			return true
		}
	}
	return false
}

func findETFRuleStatus(statuses []ETFRuleStatus, symbol string) *ETFRuleStatus {
	for i := range statuses {
		if normalizeFundSymbol(statuses[i].Symbol) == normalizeFundSymbol(symbol) {
			return &statuses[i]
		}
	}
	return nil
}

func findETFRuleMetric(metrics []ETFRuleMetric, key string) *ETFRuleMetric {
	for i := range metrics {
		if metrics[i].Key == key {
			return &metrics[i]
		}
	}
	return nil
}

func validateETFExecutionTrade(state *AppState, trade Trade) error {
	meta := trade.ETFExecution
	if meta == nil {
		return nil
	}
	if normalizeAssetType(trade.AssetType) != assetTypeFund || !strings.EqualFold(trade.Side, "buy") {
		return errors.New("ETF execution record must be a fund buy")
	}
	config, ok := etfExecutionConfig(meta.TrackerSymbol)
	if !ok {
		return errors.New("unknown ETF execution tracker")
	}
	if normalizeFundSymbol(meta.TacticalSymbol) != normalizeFundSymbol(config.TacticalSymbol) || normalizeFundSymbol(trade.Symbol) != normalizeFundSymbol(config.TacticalSymbol) {
		return errors.New("ETF execution symbol does not match tracker")
	}
	ensureETFExecutionPlans(state, time.Now())
	planIndex := findETFExecutionPlanIndex(state.ETFExecutionPlans, config.TrackerSymbol)
	if planIndex < 0 {
		return errors.New("ETF execution plan not found")
	}
	plan := &state.ETFExecutionPlans[planIndex]
	if strings.TrimSpace(meta.RoundID) != strings.TrimSpace(plan.RoundID) {
		return errors.New("ETF execution round is no longer active")
	}
	stageIndex := findETFExecutionStageIndex(plan.Stages, meta.StageKey)
	if stageIndex < 0 {
		return errors.New("ETF execution stage not found")
	}
	stage := &plan.Stages[stageIndex]
	if meta.Installment < 1 || meta.Installment > 3 {
		return errors.New("ETF execution installment is invalid")
	}
	if meta.InstallmentCount < 1 || meta.InstallmentCount > 3 || meta.Installment > meta.InstallmentCount {
		return errors.New("ETF execution installment count is invalid")
	}
	if meta.RecommendedAmount <= 0 || meta.StagePlannedAmount <= 0 {
		return errors.New("ETF execution recommended amount is required")
	}
	if meta.Override && strings.TrimSpace(meta.OverrideReason) == "" {
		return errors.New("override reason is required")
	}
	if stage.Status == etfExecutionStageCanceled && !meta.Override {
		return errors.New("ETF execution stage was canceled")
	}
	for _, installment := range stage.Installments {
		if installment.Number == meta.Installment && len(installment.TradeIDs) > 0 && !meta.Override {
			return errors.New("ETF execution installment already recorded")
		}
	}
	return nil
}

func recordETFExecutionTrade(state *AppState, trade Trade) {
	meta := trade.ETFExecution
	if state == nil || meta == nil {
		return
	}
	planIndex := findETFExecutionPlanIndex(state.ETFExecutionPlans, meta.TrackerSymbol)
	if planIndex < 0 {
		return
	}
	plan := &state.ETFExecutionPlans[planIndex]
	if strings.TrimSpace(plan.RoundID) != strings.TrimSpace(meta.RoundID) {
		return
	}
	stageIndex := findETFExecutionStageIndex(plan.Stages, meta.StageKey)
	if stageIndex < 0 {
		return
	}
	stage := &plan.Stages[stageIndex]
	if stage.PlannedAmount <= 0 {
		stage.PlannedAmount = meta.StagePlannedAmount
	}
	stage.Threshold = meta.StageThreshold
	stage.TriggeredAt = firstNonEmpty(stage.TriggeredAt, trade.Date)
	stage.SignalAsOf = firstNonEmpty(meta.SignalAsOf, stage.SignalAsOf)
	stage.InstallmentCount = meta.InstallmentCount
	installmentIndex := -1
	for i := range stage.Installments {
		if stage.Installments[i].Number == meta.Installment {
			installmentIndex = i
			break
		}
	}
	if installmentIndex < 0 {
		stage.Installments = append(stage.Installments, ETFExecutionInstallment{Number: meta.Installment})
		installmentIndex = len(stage.Installments) - 1
	}
	installment := &stage.Installments[installmentIndex]
	installment.PlannedAmount = meta.RecommendedAmount
	installment.ExecutedAmount += trade.Shares * trade.Price
	installment.ExecutedAt = trade.Date
	installment.TradeIDs = append(installment.TradeIDs, trade.ID)
	installment.Override = meta.Override
	installment.OverrideReason = meta.OverrideReason
	refreshETFExecutionStageStatus(stage)
	plan.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
}

func refreshETFExecutionStageStatus(stage *ETFExecutionStage) {
	if stage == nil {
		return
	}
	if strings.TrimSpace(stage.CanceledAt) != "" {
		stage.Status = etfExecutionStageCanceled
		return
	}
	completed := 0
	for _, installment := range stage.Installments {
		if len(installment.TradeIDs) > 0 {
			completed++
		}
	}
	if completed == 0 {
		stage.Status = etfExecutionStagePending
	} else if stage.InstallmentCount > 0 && completed >= stage.InstallmentCount {
		stage.Status = etfExecutionStageComplete
	} else {
		stage.Status = etfExecutionStagePartial
	}
}

func rebuildETFExecutionStagesFromTrades(state *AppState) {
	if state == nil {
		return
	}
	for i := range state.ETFExecutionPlans {
		for j := range state.ETFExecutionPlans[i].Stages {
			stage := &state.ETFExecutionPlans[i].Stages[j]
			stage.Installments = nil
			if strings.TrimSpace(stage.CanceledAt) == "" {
				stage.Status = etfExecutionStagePending
			}
		}
	}
	for _, trade := range state.Trades {
		if trade.ETFExecution == nil {
			continue
		}
		recordETFExecutionTrade(state, trade)
	}
}

func (s *Server) handleCancelETFExecutionStage(w http.ResponseWriter, r *http.Request) {
	var request cancelETFExecutionStageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid ETF stage cancellation payload")
		return
	}
	request.Reason = strings.TrimSpace(request.Reason)
	if request.Reason == "" {
		writeError(w, http.StatusBadRequest, "cancellation reason is required")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	nextState, err := cloneAppState(s.state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	planIndex := findETFExecutionPlanIndex(nextState.ETFExecutionPlans, request.TrackerSymbol)
	if planIndex < 0 || nextState.ETFExecutionPlans[planIndex].RoundID != request.RoundID {
		writeError(w, http.StatusNotFound, "ETF execution round not found")
		return
	}
	stageIndex := findETFExecutionStageIndex(nextState.ETFExecutionPlans[planIndex].Stages, request.StageKey)
	if stageIndex < 0 {
		writeError(w, http.StatusNotFound, "ETF execution stage not found")
		return
	}
	stage := &nextState.ETFExecutionPlans[planIndex].Stages[stageIndex]
	stage.CanceledAt = time.Now().Format("2006-01-02")
	stage.CancellationReason = request.Reason
	stage.Status = etfExecutionStageCanceled
	nextState.ETFExecutionPlans[planIndex].UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	if err := saveAndHydrateState(&nextState); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.state = nextState
	writeJSON(w, http.StatusOK, s.state)
}

func normalizeETFExecutionAmount(value float64) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	return math.Round(value*100) / 100
}
