package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const defaultSafetyMarginTarget = 0.25

type ResearchImport struct {
	UpdateType          string              `json:"updateType,omitempty"`
	Symbol              string              `json:"symbol"`
	Name                string              `json:"name"`
	AsOf                string              `json:"asOf"`
	Currency            string              `json:"currency"`
	Industry            string              `json:"industry"`
	Status              string              `json:"status"`
	Action              string              `json:"action"`
	Risk                string              `json:"risk"`
	Valuation           Valuation           `json:"valuation"`
	Quality             Quality             `json:"quality"`
	Plan                PlanInput           `json:"plan"`
	Dividend            *Dividend           `json:"dividend,omitempty"`
	NetCash             *NetCashProfile     `json:"netCash,omitempty"`
	OwnerCashFlowAudit  *OwnerCashFlowAudit `json:"ownerCashFlowAudit,omitempty"`
	ValuationConfidence string              `json:"valuationConfidence,omitempty"`
	KillCriteria        json.RawMessage     `json:"killCriteria,omitempty"`
	Notes               string              `json:"notes"`
	Event               *ResearchEvent      `json:"event,omitempty"`
	Impact              *ResearchImpact     `json:"impact,omitempty"`
	Updates             *ResearchPatch      `json:"updates,omitempty"`
}

type Valuation struct {
	IntrinsicValue *float64     `json:"intrinsicValue"`
	FairValueRange string       `json:"fairValueRange"`
	TargetBuyPrice *float64     `json:"targetBuyPrice"`
	MarginOfSafety *float64     `json:"marginOfSafety"`
	PriceLevels    *PriceLevels `json:"priceLevels,omitempty"`
}

type Quality struct {
	TotalScore       *float64 `json:"totalScore"`
	BusinessModel    *float64 `json:"businessModel"`
	Moat             *float64 `json:"moat"`
	Governance       *float64 `json:"governance"`
	FinancialQuality *float64 `json:"financialQuality"`
}

type PlanInput struct {
	Rank       int    `json:"rank"`
	Priority   string `json:"priority"`
	Advice     string `json:"advice"`
	Discipline string `json:"discipline"`
}

type ResearchPatch struct {
	Status              string              `json:"status,omitempty"`
	Action              string              `json:"action,omitempty"`
	Risk                string              `json:"risk,omitempty"`
	Notes               string              `json:"notes,omitempty"`
	NotesAppend         string              `json:"notesAppend,omitempty"`
	Valuation           *Valuation          `json:"valuation,omitempty"`
	Quality             *Quality            `json:"quality,omitempty"`
	Plan                *PlanInput          `json:"plan,omitempty"`
	Dividend            *Dividend           `json:"dividend,omitempty"`
	NetCash             *NetCashProfile     `json:"netCash,omitempty"`
	OwnerCashFlowAudit  *OwnerCashFlowAudit `json:"ownerCashFlowAudit,omitempty"`
	ValuationConfidence string              `json:"valuationConfidence,omitempty"`
	KillCriteria        json.RawMessage     `json:"killCriteria,omitempty"`
}

type ResearchResponse struct {
	Summary       string         `json:"summary"`
	TargetType    string         `json:"targetType"`
	Warnings      []string       `json:"warnings"`
	Research      ResearchImport `json:"research"`
	Plan          []PlanItem     `json:"plan"`
	ChangedFields []string       `json:"changedFields,omitempty"`
	BackupPath    string         `json:"backupPath,omitempty"`
	State         *AppState      `json:"state,omitempty"`
}

func (s *Server) handlePreviewResearch(w http.ResponseWriter, r *http.Request) {
	research, err := decodeResearch(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	warnings, err := validateResearch(research)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	state, err := loadState()
	s.mu.Unlock()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}

	summary, targetType, changedFields := applyResearch(&state, research)
	writeJSON(w, http.StatusOK, ResearchResponse{
		Summary:       summary,
		TargetType:    targetType,
		Warnings:      warnings,
		Research:      normalizeResearch(research),
		Plan:          state.Plan,
		ChangedFields: changedFields,
	})
}

func (s *Server) handleImportResearch(w http.ResponseWriter, r *http.Request) {
	research, err := decodeResearch(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	warnings, err := validateResearch(research)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := loadState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}

	summary, targetType, changedFields := applyResearch(&state, research)
	appendResearchDecisionLog(&state, research, summary, targetType)
	backupPath, err := saveStateWithBackup(state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save state")
		return
	}
	if err := hydrateState(&state); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}

	s.state = state
	writeJSON(w, http.StatusOK, ResearchResponse{
		Summary:       summary,
		TargetType:    targetType,
		Warnings:      warnings,
		Research:      normalizeResearch(research),
		Plan:          state.Plan,
		ChangedFields: changedFields,
		BackupPath:    backupPath,
		State:         &state,
	})
}

func decodeResearch(r *http.Request) (ResearchImport, error) {
	defer r.Body.Close()

	var research ResearchImport
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&research); err != nil {
		return ResearchImport{}, fmt.Errorf("invalid research JSON: %w", err)
	}
	return normalizeResearch(research), nil
}

func validateResearch(research ResearchImport) ([]string, error) {
	warnings := []string{}

	if strings.TrimSpace(research.Symbol) == "" {
		return nil, errors.New("symbol is required")
	}
	if strings.TrimSpace(research.Name) == "" {
		return nil, errors.New("name is required")
	}
	if strings.TrimSpace(research.AsOf) == "" {
		return nil, errors.New("asOf is required")
	}
	if research.UpdateType != "fullReview" && research.UpdateType != "eventUpdate" {
		return nil, errors.New("updateType must be fullReview or eventUpdate")
	}

	asOfDate, err := time.Parse("2006-01-02", research.AsOf)
	if err != nil {
		return nil, fmt.Errorf("asOf must be YYYY-MM-DD: %w", err)
	}
	today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if asOfDate.After(today) {
		return nil, fmt.Errorf("asOf cannot be later than today: %s", today.Format("2006-01-02"))
	}

	if research.UpdateType == "eventUpdate" {
		eventWarnings, err := validateEventUpdate(research)
		if err != nil {
			return nil, err
		}
		warnings = append(warnings, eventWarnings...)
	} else {
		if err := validateFullReview(research); err != nil {
			return nil, err
		}
	}

	if expected := expectedCurrency(research.Symbol); expected != "" && research.Currency != "" && research.Currency != expected {
		warnings = append(warnings, fmt.Sprintf("currency is %s, but %s usually uses %s", research.Currency, research.Symbol, expected))
	}
	qualityForCheck := research.Quality
	if research.UpdateType == "eventUpdate" && research.Updates != nil && research.Updates.Quality != nil {
		qualityForCheck = *research.Updates.Quality
	}
	if hasAllQualityScores(qualityForCheck) {
		sum := *qualityForCheck.BusinessModel + *qualityForCheck.Moat + *qualityForCheck.Governance + *qualityForCheck.FinancialQuality
		if math.Abs(*qualityForCheck.TotalScore-sum) > 0.01 {
			warnings = append(warnings, fmt.Sprintf("quality.totalScore is %.2f, but component scores add up to %.2f", *qualityForCheck.TotalScore, sum))
		}
	}
	if research.UpdateType == "eventUpdate" {
		return warnings, nil
	}
	if strings.TrimSpace(research.Status) == "" {
		warnings = append(warnings, "status is empty")
	}
	if strings.TrimSpace(research.Action) == "" {
		warnings = append(warnings, "action is empty")
	}
	if strings.TrimSpace(research.Risk) == "" {
		warnings = append(warnings, "risk is empty")
	}
	if strings.TrimSpace(research.Notes) == "" {
		warnings = append(warnings, "notes is empty")
	}
	if strings.TrimSpace(research.Plan.Priority) == "" &&
		strings.TrimSpace(research.Plan.Advice) == "" &&
		strings.TrimSpace(research.Plan.Discipline) == "" {
		warnings = append(warnings, "plan is empty; overview execution plan will not be updated")
	}

	return warnings, nil
}

func validateFullReview(research ResearchImport) error {
	if err := validatePercent("valuation.marginOfSafety", research.Valuation.MarginOfSafety); err != nil {
		return err
	}
	if err := validateDividend(research.Dividend); err != nil {
		return err
	}
	if err := validateNetCash(research.NetCash); err != nil {
		return err
	}
	if err := validateOwnerCashFlowAudit(research.OwnerCashFlowAudit); err != nil {
		return err
	}
	return validateQuality(research.Quality)
}

func validateEventUpdate(research ResearchImport) ([]string, error) {
	warnings := []string{}
	if research.Event == nil {
		return nil, errors.New("event is required when updateType is eventUpdate")
	}
	if strings.TrimSpace(research.Event.Title) == "" {
		return nil, errors.New("event.title is required when updateType is eventUpdate")
	}
	if strings.TrimSpace(research.Event.Date) == "" {
		return nil, errors.New("event.date is required when updateType is eventUpdate")
	}
	if _, err := time.Parse("2006-01-02", research.Event.Date); err != nil {
		return nil, fmt.Errorf("event.date must be YYYY-MM-DD: %w", err)
	}
	if strings.TrimSpace(research.Event.Type) == "" {
		warnings = append(warnings, "event.type is empty")
	}
	if strings.TrimSpace(research.Event.Source) == "" {
		warnings = append(warnings, "event.source is empty")
	}
	if strings.TrimSpace(research.Event.Summary) == "" {
		warnings = append(warnings, "event.summary is empty")
	}
	if research.Updates == nil {
		warnings = append(warnings, "updates is empty; only research timeline will be appended")
		return warnings, nil
	}
	if research.Updates.Valuation != nil {
		if err := validatePercent("updates.valuation.marginOfSafety", research.Updates.Valuation.MarginOfSafety); err != nil {
			return nil, err
		}
	}
	if research.Updates.Dividend != nil {
		if err := validateDividend(research.Updates.Dividend); err != nil {
			return nil, err
		}
	}
	if research.Updates.NetCash != nil {
		if err := validateNetCash(research.Updates.NetCash); err != nil {
			return nil, err
		}
	}
	if research.Updates.OwnerCashFlowAudit != nil {
		if err := validateOwnerCashFlowAudit(research.Updates.OwnerCashFlowAudit); err != nil {
			return nil, err
		}
	}
	if research.Updates.Quality != nil {
		if err := validateQuality(*research.Updates.Quality); err != nil {
			return nil, err
		}
	}
	return warnings, nil
}

func validateQuality(quality Quality) error {
	if err := validateScore("quality.totalScore", quality.TotalScore, 100); err != nil {
		return err
	}
	if err := validateScore("quality.businessModel", quality.BusinessModel, 30); err != nil {
		return err
	}
	if err := validateScore("quality.moat", quality.Moat, 25); err != nil {
		return err
	}
	if err := validateScore("quality.governance", quality.Governance, 20); err != nil {
		return err
	}
	return validateScore("quality.financialQuality", quality.FinancialQuality, 25)
}

func applyResearch(state *AppState, research ResearchImport) (string, string, []string) {
	symbol := normalizeSymbol(research.Symbol)
	now := time.Now().Format("2006-01-02 15:04:05")
	updateLabel := fmt.Sprintf("%s；ChatGPT分析导入；分析日 %s", now, research.AsOf)

	for i := range state.Holdings {
		if normalizeSymbol(state.Holdings[i].Symbol) == symbol {
			if research.UpdateType == "eventUpdate" {
				changedFields := applyHoldingEventUpdate(&state.Holdings[i], research, now)
				upsertPlanFromPatch(state, research)
				return fmt.Sprintf("event update for holding %s (%s)", research.Symbol, research.Name), "holding", changedFields
			}
			applyHoldingResearch(&state.Holdings[i], research, updateLabel)
			upsertPlan(state, research)
			return fmt.Sprintf("updated holding %s (%s)", research.Symbol, research.Name), "holding", fullReviewChangedFields()
		}
	}

	for i := range state.Candidates {
		if normalizeSymbol(state.Candidates[i].Symbol) == symbol {
			if research.UpdateType == "eventUpdate" {
				changedFields := applyCandidateEventUpdate(&state.Candidates[i], research, now)
				upsertPlanFromPatch(state, research)
				return fmt.Sprintf("event update for candidate %s (%s)", research.Symbol, research.Name), "candidate", changedFields
			}
			applyCandidateResearch(&state.Candidates[i], research, updateLabel)
			upsertPlan(state, research)
			return fmt.Sprintf("updated candidate %s (%s)", research.Symbol, research.Name), "candidate", fullReviewChangedFields()
		}
	}

	state.Candidates = append(state.Candidates, Candidate{Symbol: normalizeDisplaySymbol(research.Symbol)})
	if research.UpdateType == "eventUpdate" {
		changedFields := applyCandidateEventUpdate(&state.Candidates[len(state.Candidates)-1], research, now)
		upsertPlanFromPatch(state, research)
		return fmt.Sprintf("added candidate event update %s (%s)", research.Symbol, research.Name), "newCandidate", changedFields
	}
	applyCandidateResearch(&state.Candidates[len(state.Candidates)-1], research, updateLabel)
	upsertPlan(state, research)
	return fmt.Sprintf("added candidate %s (%s)", research.Symbol, research.Name), "newCandidate", fullReviewChangedFields()
}

func applyHoldingResearch(holding *Holding, research ResearchImport, updateLabel string) {
	holding.Symbol = normalizeDisplaySymbol(research.Symbol)
	holding.Name = strings.TrimSpace(research.Name)
	holding.Industry = strings.TrimSpace(research.Industry)
	holding.Status = strings.TrimSpace(research.Status)
	holding.Action = strings.TrimSpace(research.Action)
	holding.Risk = strings.TrimSpace(research.Risk)
	holding.Currency = prefer(holding.Currency, research.Currency)
	holding.QualityScore = research.Quality.TotalScore
	holding.IntrinsicValue = research.Valuation.IntrinsicValue
	holding.FairValueRange = strings.TrimSpace(research.Valuation.FairValueRange)
	holding.TargetBuyPrice = research.Valuation.TargetBuyPrice
	holding.MarginOfSafety = marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, research.Valuation.MarginOfSafety)
	holding.ValuationConfidence = strings.TrimSpace(research.ValuationConfidence)
	holding.BusinessModel = research.Quality.BusinessModel
	holding.Moat = research.Quality.Moat
	holding.Governance = research.Quality.Governance
	holding.FinancialQuality = research.Quality.FinancialQuality
	holding.UpdatedAt = updateLabel
	holding.Notes = strings.TrimSpace(research.Notes)
	holding.KillCriteria = cloneRawMessage(research.KillCriteria)
	if research.Dividend != nil {
		holding.Dividend = research.Dividend
	}
	if research.NetCash != nil {
		holding.NetCash = research.NetCash
	}
	if research.OwnerCashFlowAudit != nil {
		holding.OwnerCashFlowAudit = research.OwnerCashFlowAudit
	}
}

func applyCandidateResearch(candidate *Candidate, research ResearchImport, updateLabel string) {
	candidate.Symbol = normalizeDisplaySymbol(research.Symbol)
	candidate.Name = strings.TrimSpace(research.Name)
	candidate.Status = strings.TrimSpace(research.Status)
	candidate.Action = strings.TrimSpace(research.Action)
	candidate.Risk = strings.TrimSpace(research.Risk)
	candidate.Industry = strings.TrimSpace(research.Industry)
	candidate.Currency = prefer(candidate.Currency, research.Currency)
	candidate.MarginOfSafety = research.Valuation.MarginOfSafety
	candidate.QualityScore = research.Quality.TotalScore
	candidate.IntrinsicValue = research.Valuation.IntrinsicValue
	candidate.FairValueRange = strings.TrimSpace(research.Valuation.FairValueRange)
	candidate.TargetBuyPrice = research.Valuation.TargetBuyPrice
	candidate.ValuationConfidence = strings.TrimSpace(research.ValuationConfidence)
	candidate.BusinessModel = research.Quality.BusinessModel
	candidate.Moat = research.Quality.Moat
	candidate.Governance = research.Quality.Governance
	candidate.FinancialQuality = research.Quality.FinancialQuality
	candidate.UpdatedAt = updateLabel
	candidate.Notes = strings.TrimSpace(research.Notes)
	candidate.KillCriteria = cloneRawMessage(research.KillCriteria)
	if research.Dividend != nil {
		candidate.Dividend = research.Dividend
	}
	if research.NetCash != nil {
		candidate.NetCash = research.NetCash
	}
	if research.OwnerCashFlowAudit != nil {
		candidate.OwnerCashFlowAudit = research.OwnerCashFlowAudit
	}
}

func applyHoldingEventUpdate(holding *Holding, research ResearchImport, importedAt string) []string {
	holding.Symbol = normalizeDisplaySymbol(research.Symbol)
	if strings.TrimSpace(research.Name) != "" {
		holding.Name = strings.TrimSpace(research.Name)
	}
	changedFields := applyResearchPatchToHolding(holding, research)
	appendHoldingResearchUpdate(holding, research, importedAt, changedFields)
	return changedFields
}

func applyCandidateEventUpdate(candidate *Candidate, research ResearchImport, importedAt string) []string {
	candidate.Symbol = normalizeDisplaySymbol(research.Symbol)
	if strings.TrimSpace(research.Name) != "" {
		candidate.Name = strings.TrimSpace(research.Name)
	}
	changedFields := applyResearchPatchToCandidate(candidate, research)
	appendCandidateResearchUpdate(candidate, research, importedAt, changedFields)
	return changedFields
}

func applyResearchPatchToHolding(holding *Holding, research ResearchImport) []string {
	patch := research.Updates
	changed := []string{"researchUpdates"}
	if patch == nil {
		return changed
	}
	if text := strings.TrimSpace(patch.Status); text != "" {
		holding.Status = text
		changed = append(changed, "status")
	}
	if text := strings.TrimSpace(patch.Action); text != "" {
		holding.Action = text
		changed = append(changed, "action")
	}
	if text := strings.TrimSpace(patch.Risk); text != "" {
		holding.Risk = text
		changed = append(changed, "risk")
	}
	if text := strings.TrimSpace(patch.Notes); text != "" {
		holding.Notes = text
		changed = append(changed, "notes")
	}
	if text := strings.TrimSpace(patch.ValuationConfidence); text != "" {
		holding.ValuationConfidence = text
		changed = append(changed, "valuationConfidence")
	}
	if len(patch.KillCriteria) > 0 {
		holding.KillCriteria = cloneRawMessage(patch.KillCriteria)
		changed = append(changed, "killCriteria")
	}
	if patch.Valuation != nil {
		if applyValuationPatch(&holding.IntrinsicValue, &holding.FairValueRange, &holding.TargetBuyPrice, patch.Valuation) {
			holding.MarginOfSafety = marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, patch.Valuation.MarginOfSafety)
			changed = append(changed, "valuation")
		}
	}
	if patch.Quality != nil && applyQualityPatch(&holding.QualityScore, &holding.BusinessModel, &holding.Moat, &holding.Governance, &holding.FinancialQuality, patch.Quality) {
		changed = append(changed, "quality")
	}
	if mergeDividendPatch(&holding.Dividend, patch.Dividend, holding.Currency) {
		changed = append(changed, "dividend")
	}
	if mergeNetCashPatch(&holding.NetCash, patch.NetCash) {
		changed = append(changed, "netCash")
	}
	if mergeOwnerAuditPatch(&holding.OwnerCashFlowAudit, patch.OwnerCashFlowAudit) {
		changed = append(changed, "ownerCashFlowAudit")
	}
	return uniqueStrings(changed)
}

func applyResearchPatchToCandidate(candidate *Candidate, research ResearchImport) []string {
	patch := research.Updates
	changed := []string{"researchUpdates"}
	if patch == nil {
		return changed
	}
	if text := strings.TrimSpace(patch.Status); text != "" {
		candidate.Status = text
		changed = append(changed, "status")
	}
	if text := strings.TrimSpace(patch.Action); text != "" {
		candidate.Action = text
		changed = append(changed, "action")
	}
	if text := strings.TrimSpace(patch.Risk); text != "" {
		candidate.Risk = text
		changed = append(changed, "risk")
	}
	if text := strings.TrimSpace(patch.Notes); text != "" {
		candidate.Notes = text
		changed = append(changed, "notes")
	}
	if text := strings.TrimSpace(patch.ValuationConfidence); text != "" {
		candidate.ValuationConfidence = text
		changed = append(changed, "valuationConfidence")
	}
	if len(patch.KillCriteria) > 0 {
		candidate.KillCriteria = cloneRawMessage(patch.KillCriteria)
		changed = append(changed, "killCriteria")
	}
	if patch.Valuation != nil {
		if applyValuationPatch(&candidate.IntrinsicValue, &candidate.FairValueRange, &candidate.TargetBuyPrice, patch.Valuation) {
			candidate.MarginOfSafety = marginOfSafetyFromPrice(candidate.IntrinsicValue, candidate.CurrentPrice, patch.Valuation.MarginOfSafety)
			changed = append(changed, "valuation")
		}
	}
	if patch.Quality != nil && applyQualityPatch(&candidate.QualityScore, &candidate.BusinessModel, &candidate.Moat, &candidate.Governance, &candidate.FinancialQuality, patch.Quality) {
		changed = append(changed, "quality")
	}
	if mergeDividendPatch(&candidate.Dividend, patch.Dividend, candidate.Currency) {
		changed = append(changed, "dividend")
	}
	if mergeNetCashPatch(&candidate.NetCash, patch.NetCash) {
		changed = append(changed, "netCash")
	}
	if mergeOwnerAuditPatch(&candidate.OwnerCashFlowAudit, patch.OwnerCashFlowAudit) {
		changed = append(changed, "ownerCashFlowAudit")
	}
	return uniqueStrings(changed)
}

func applyValuationPatch(intrinsicValue **float64, fairValueRange *string, targetBuyPrice **float64, valuation *Valuation) bool {
	changed := false
	if valuation.IntrinsicValue != nil {
		*intrinsicValue = valuation.IntrinsicValue
		changed = true
	}
	if text := strings.TrimSpace(valuation.FairValueRange); text != "" {
		*fairValueRange = text
		changed = true
	}
	if valuation.TargetBuyPrice != nil {
		*targetBuyPrice = valuation.TargetBuyPrice
		changed = true
	}
	if valuation.MarginOfSafety != nil {
		changed = true
	}
	return changed
}

func applyQualityPatch(totalScore **float64, businessModel **float64, moat **float64, governance **float64, financialQuality **float64, quality *Quality) bool {
	changed := false
	if quality.TotalScore != nil {
		*totalScore = quality.TotalScore
		changed = true
	}
	if quality.BusinessModel != nil {
		*businessModel = quality.BusinessModel
		changed = true
	}
	if quality.Moat != nil {
		*moat = quality.Moat
		changed = true
	}
	if quality.Governance != nil {
		*governance = quality.Governance
		changed = true
	}
	if quality.FinancialQuality != nil {
		*financialQuality = quality.FinancialQuality
		changed = true
	}
	return changed
}

func mergeDividendPatch(current **Dividend, patch *Dividend, fallbackCurrency string) bool {
	if patch == nil {
		return false
	}
	if *current == nil {
		*current = &Dividend{}
	}
	changed := false
	dividend := *current
	if text := strings.TrimSpace(patch.FiscalYear); text != "" {
		dividend.FiscalYear = text
		changed = true
	}
	if patch.DividendPerShare != nil {
		dividend.DividendPerShare = patch.DividendPerShare
		changed = true
	}
	if text := strings.TrimSpace(patch.DividendCurrency); text != "" {
		dividend.DividendCurrency = strings.ToUpper(text)
		changed = true
	} else if dividend.DividendCurrency == "" && fallbackCurrency != "" {
		dividend.DividendCurrency = strings.ToUpper(fallbackCurrency)
	}
	if patch.CashDividendTotal != nil {
		dividend.CashDividendTotal = patch.CashDividendTotal
		changed = true
	}
	if text := strings.TrimSpace(patch.CashDividendCurrency); text != "" {
		dividend.CashDividendCurrency = strings.ToUpper(text)
		changed = true
	}
	if patch.BuybackAmount != nil {
		dividend.BuybackAmount = patch.BuybackAmount
		changed = true
	}
	if text := strings.TrimSpace(patch.BuybackCurrency); text != "" {
		dividend.BuybackCurrency = strings.ToUpper(text)
		changed = true
	}
	if patch.DividendYield != nil {
		dividend.DividendYield = patch.DividendYield
		changed = true
	}
	if patch.PayoutRatio != nil {
		dividend.PayoutRatio = patch.PayoutRatio
		changed = true
	}
	if patch.EstimatedAnnualCash != nil {
		dividend.EstimatedAnnualCash = patch.EstimatedAnnualCash
		changed = true
	}
	if text := strings.TrimSpace(patch.Reliability); text != "" {
		dividend.Reliability = text
		changed = true
	}
	if text := strings.TrimSpace(patch.ForecastFiscalYear); text != "" {
		dividend.ForecastFiscalYear = text
		changed = true
	}
	if patch.ForecastPerShare != nil {
		dividend.ForecastPerShare = patch.ForecastPerShare
		changed = true
	}
	if text := strings.TrimSpace(patch.ForecastCurrency); text != "" {
		dividend.ForecastCurrency = strings.ToUpper(text)
		changed = true
	}
	if patch.ForecastYield != nil {
		dividend.ForecastYield = patch.ForecastYield
		changed = true
	}
	return changed
}

func mergeNetCashPatch(current **NetCashProfile, patch *NetCashProfile) bool {
	if patch == nil {
		return false
	}
	if *current == nil {
		*current = &NetCashProfile{}
	}
	changed := false
	netCash := *current
	if patch.CashAndShortInvestments != nil {
		netCash.CashAndShortInvestments = patch.CashAndShortInvestments
		changed = true
	}
	if patch.InterestBearingDebt != nil {
		netCash.InterestBearingDebt = patch.InterestBearingDebt
		changed = true
	}
	if patch.NetCash != nil {
		netCash.NetCash = patch.NetCash
		changed = true
	}
	if text := strings.TrimSpace(patch.Currency); text != "" {
		netCash.Currency = strings.ToUpper(text)
		changed = true
	}
	if patch.Haircut != nil {
		netCash.Haircut = patch.Haircut
		changed = true
	}
	if text := strings.TrimSpace(patch.HaircutReason); text != "" {
		netCash.HaircutReason = text
		changed = true
	}
	if patch.AdjustedNetCash != nil {
		netCash.AdjustedNetCash = patch.AdjustedNetCash
		changed = true
	}
	if patch.ExCashPE != nil {
		netCash.ExCashPE = patch.ExCashPE
		changed = true
	}
	if patch.ExCashPFCF != nil {
		netCash.ExCashPFCF = patch.ExCashPFCF
		changed = true
	}
	if patch.FCFYield != nil {
		netCash.FCFYield = patch.FCFYield
		changed = true
	}
	if patch.ShareholderFCF != nil {
		netCash.ShareholderFCF = patch.ShareholderFCF
		changed = true
	}
	if text := strings.TrimSpace(patch.ShareholderFCFCurrency); text != "" {
		netCash.ShareholderFCFCurrency = strings.ToUpper(text)
		changed = true
	}
	if text := strings.TrimSpace(patch.ShareholderFCFBasis); text != "" {
		netCash.ShareholderFCFBasis = text
		changed = true
	}
	if patch.ConsolidatedFCF != nil {
		netCash.ConsolidatedFCF = patch.ConsolidatedFCF
		changed = true
	}
	if patch.MinorityFCFAdjustment != nil {
		netCash.MinorityFCFAdjustment = patch.MinorityFCFAdjustment
		changed = true
	}
	if patch.FCFPositiveYears != nil {
		netCash.FCFPositiveYears = patch.FCFPositiveYears
		changed = true
	}
	if text := strings.TrimSpace(patch.Note); text != "" {
		netCash.Note = text
		changed = true
	}
	return changed
}

func mergeOwnerAuditPatch(current **OwnerCashFlowAudit, patch *OwnerCashFlowAudit) bool {
	if patch == nil {
		return false
	}
	if *current == nil {
		*current = &OwnerCashFlowAudit{}
	}
	changed := false
	audit := *current
	changed = mergeOwnerAuditItem(&audit.TenYearDemand, patch.TenYearDemand) || changed
	changed = mergeOwnerAuditItem(&audit.AssetDurability, patch.AssetDurability) || changed
	changed = mergeOwnerAuditItem(&audit.MaintenanceCapexLight, patch.MaintenanceCapexLight) || changed
	changed = mergeOwnerAuditItem(&audit.DividendFCFSupport, patch.DividendFCFSupport) || changed
	changed = mergeOwnerAuditItem(&audit.DividendReinvestmentEfficiency, patch.DividendReinvestmentEfficiency) || changed
	changed = mergeOwnerAuditItem(&audit.RoeRoicDurability, patch.RoeRoicDurability) || changed
	changed = mergeOwnerAuditItem(&audit.ValuationSystemRisk, patch.ValuationSystemRisk) || changed
	return changed
}

func mergeOwnerAuditItem(current *OwnerAuditItem, patch OwnerAuditItem) bool {
	changed := false
	if text := strings.TrimSpace(patch.Status); text != "" {
		current.Status = text
		changed = true
	}
	if text := strings.TrimSpace(patch.Note); text != "" {
		current.Note = text
		changed = true
	}
	return changed
}

func appendHoldingResearchUpdate(holding *Holding, research ResearchImport, importedAt string, changedFields []string) {
	holding.ResearchUpdates = appendResearchUpdate(holding.ResearchUpdates, research, importedAt, changedFields)
}

func appendCandidateResearchUpdate(candidate *Candidate, research ResearchImport, importedAt string, changedFields []string) {
	candidate.ResearchUpdates = appendResearchUpdate(candidate.ResearchUpdates, research, importedAt, changedFields)
}

func appendResearchUpdate(updates []ResearchUpdate, research ResearchImport, importedAt string, changedFields []string) []ResearchUpdate {
	event := ResearchEvent{}
	if research.Event != nil {
		event = *research.Event
	}
	impact := ResearchImpact{}
	if research.Impact != nil {
		impact = *research.Impact
	}
	notesAppend := ""
	if research.Updates != nil {
		notesAppend = strings.TrimSpace(research.Updates.NotesAppend)
	}
	entry := ResearchUpdate{
		ID:            time.Now().UnixNano(),
		ImportedAt:    importedAt,
		AsOf:          research.AsOf,
		UpdateType:    research.UpdateType,
		Event:         event,
		Impact:        impact,
		Summary:       firstNonEmpty(event.Summary, notesAppend),
		ChangedFields: uniqueStrings(changedFields),
		NotesAppend:   notesAppend,
	}
	updates = append(updates, entry)
	if len(updates) > 50 {
		updates = updates[len(updates)-50:]
	}
	return updates
}

func fullReviewChangedFields() []string {
	return []string{"valuation", "quality", "status", "action", "risk", "plan", "dividend", "netCash", "ownerCashFlowAudit", "notes"}
}

func upsertPlan(state *AppState, research ResearchImport) {
	if strings.TrimSpace(research.Plan.Priority) == "" &&
		strings.TrimSpace(research.Plan.Advice) == "" &&
		strings.TrimSpace(research.Plan.Discipline) == "" {
		return
	}

	next := PlanItem{
		Rank:       research.Plan.Rank,
		Symbol:     normalizeDisplaySymbol(research.Symbol),
		Name:       strings.TrimSpace(research.Name),
		Priority:   strings.TrimSpace(research.Plan.Priority),
		Advice:     strings.TrimSpace(research.Plan.Advice),
		Discipline: strings.TrimSpace(research.Plan.Discipline),
	}
	if next.Rank <= 0 {
		next.Rank = nextPlanRank(state.Plan)
	}

	for i := range state.Plan {
		if samePlanItem(state.Plan[i], next) {
			state.Plan[i] = next
			normalizePlanRanks(state.Plan)
			return
		}
	}

	state.Plan = append(state.Plan, next)
	normalizePlanRanks(state.Plan)
}

func upsertPlanFromPatch(state *AppState, research ResearchImport) {
	if research.Updates == nil || research.Updates.Plan == nil {
		return
	}
	plan := research.Updates.Plan
	if strings.TrimSpace(plan.Priority) == "" &&
		strings.TrimSpace(plan.Advice) == "" &&
		strings.TrimSpace(plan.Discipline) == "" {
		return
	}
	upsertPlan(state, ResearchImport{
		Symbol: normalizeDisplaySymbol(research.Symbol),
		Name:   strings.TrimSpace(research.Name),
		Plan:   *plan,
	})
}

func samePlanItem(current, next PlanItem) bool {
	if normalizeSymbol(current.Symbol) != "" && normalizeSymbol(current.Symbol) == normalizeSymbol(next.Symbol) {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(current.Name), strings.TrimSpace(next.Name))
}

func nextPlanRank(plan []PlanItem) int {
	next := 1
	for _, item := range plan {
		if item.Rank >= next {
			next = item.Rank + 1
		}
	}
	return next
}

func normalizePlanRanks(plan []PlanItem) {
	sort.SliceStable(plan, func(i, j int) bool {
		if plan[i].Rank == plan[j].Rank {
			return false
		}
		if plan[i].Rank <= 0 {
			return false
		}
		if plan[j].Rank <= 0 {
			return true
		}
		return plan[i].Rank < plan[j].Rank
	})
	for i := range plan {
		plan[i].Rank = i + 1
	}
}

func backupPortfolioFile() (string, error) {
	body, err := os.ReadFile(dataFile)
	if err != nil {
		return "", err
	}

	backupDir := filepath.Join(filepath.Dir(dataFile), "backups")
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return "", err
	}

	now := time.Now()
	backupPath := filepath.Join(backupDir, fmt.Sprintf("portfolio-%s-%03d.json", now.Format("20060102-150405"), now.UnixMilli()%1000))
	if err := writeFileAtomic(backupPath, body, 0o644); err != nil {
		return "", err
	}
	return backupPath, nil
}

func normalizeResearch(research ResearchImport) ResearchImport {
	research.UpdateType = strings.TrimSpace(research.UpdateType)
	if research.UpdateType == "" {
		research.UpdateType = "fullReview"
	}
	research.Symbol = normalizeDisplaySymbol(research.Symbol)
	research.Name = strings.TrimSpace(research.Name)
	research.AsOf = strings.TrimSpace(research.AsOf)
	research.Currency = strings.ToUpper(strings.TrimSpace(research.Currency))
	research.Industry = strings.TrimSpace(research.Industry)
	research.Status = strings.TrimSpace(research.Status)
	research.Action = strings.TrimSpace(research.Action)
	research.Risk = strings.TrimSpace(research.Risk)
	research.Valuation.FairValueRange = strings.TrimSpace(research.Valuation.FairValueRange)
	research.Valuation.TargetBuyPrice = targetBuyPriceFromIntrinsicValue(research.Valuation.IntrinsicValue)
	research.Valuation.PriceLevels = nil
	research.Plan.Priority = strings.TrimSpace(research.Plan.Priority)
	research.Plan.Advice = strings.TrimSpace(research.Plan.Advice)
	research.Plan.Discipline = strings.TrimSpace(research.Plan.Discipline)
	research.ValuationConfidence = strings.TrimSpace(research.ValuationConfidence)
	research.KillCriteria = normalizeRawMessage(research.KillCriteria)
	research.Dividend = normalizeDividend(research.Dividend, research.Currency)
	research.NetCash = normalizeNetCash(research.NetCash, research.Currency)
	research.OwnerCashFlowAudit = normalizeOwnerCashFlowAudit(research.OwnerCashFlowAudit)
	research.Notes = strings.TrimSpace(research.Notes)
	if research.Event != nil {
		research.Event.Type = strings.TrimSpace(research.Event.Type)
		research.Event.Title = strings.TrimSpace(research.Event.Title)
		research.Event.Date = strings.TrimSpace(research.Event.Date)
		research.Event.Source = strings.TrimSpace(research.Event.Source)
		research.Event.Summary = strings.TrimSpace(research.Event.Summary)
	}
	if research.Impact != nil {
		research.Impact.ThesisChange = strings.TrimSpace(research.Impact.ThesisChange)
		research.Impact.ValuationChange = strings.TrimSpace(research.Impact.ValuationChange)
		research.Impact.RiskChange = strings.TrimSpace(research.Impact.RiskChange)
		research.Impact.ActionChange = strings.TrimSpace(research.Impact.ActionChange)
	}
	research.Updates = normalizeResearchPatch(research.Updates, research.Currency)
	return research
}

func normalizeResearchPatch(patch *ResearchPatch, currency string) *ResearchPatch {
	if patch == nil {
		return nil
	}
	patch.Status = strings.TrimSpace(patch.Status)
	patch.Action = strings.TrimSpace(patch.Action)
	patch.Risk = strings.TrimSpace(patch.Risk)
	patch.Notes = strings.TrimSpace(patch.Notes)
	patch.NotesAppend = strings.TrimSpace(patch.NotesAppend)
	patch.ValuationConfidence = strings.TrimSpace(patch.ValuationConfidence)
	patch.KillCriteria = normalizeRawMessage(patch.KillCriteria)
	if patch.Valuation != nil {
		patch.Valuation.FairValueRange = strings.TrimSpace(patch.Valuation.FairValueRange)
		if patch.Valuation.TargetBuyPrice == nil && patch.Valuation.IntrinsicValue != nil {
			patch.Valuation.TargetBuyPrice = targetBuyPriceFromIntrinsicValue(patch.Valuation.IntrinsicValue)
		}
		patch.Valuation.PriceLevels = nil
	}
	if patch.Plan != nil {
		patch.Plan.Priority = strings.TrimSpace(patch.Plan.Priority)
		patch.Plan.Advice = strings.TrimSpace(patch.Plan.Advice)
		patch.Plan.Discipline = strings.TrimSpace(patch.Plan.Discipline)
	}
	patch.Dividend = normalizeDividend(patch.Dividend, currency)
	patch.NetCash = normalizeNetCash(patch.NetCash, currency)
	patch.OwnerCashFlowAudit = normalizeOwnerCashFlowAudit(patch.OwnerCashFlowAudit)
	return patch
}

func normalizeRawMessage(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 || strings.TrimSpace(string(raw)) == "" || string(raw) == "null" {
		return nil
	}
	return cloneRawMessage(raw)
}

func cloneRawMessage(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return nil
	}
	next := make(json.RawMessage, len(raw))
	copy(next, raw)
	return next
}

func hasAllQualityScores(quality Quality) bool {
	return quality.TotalScore != nil &&
		quality.BusinessModel != nil &&
		quality.Moat != nil &&
		quality.Governance != nil &&
		quality.FinancialQuality != nil
}

func validatePercent(field string, value *float64) error {
	if value == nil {
		return nil
	}
	if *value < -1 || *value > 1 {
		return fmt.Errorf("%s must be decimal ratio, for example 0.09 for 9%%", field)
	}
	return nil
}

func validatePriceLevels(levels *PriceLevels) error {
	if levels == nil {
		return nil
	}
	if err := validatePositiveAmount("valuation.priceLevels.watchPrice", levels.WatchPrice); err != nil {
		return err
	}
	if err := validatePositiveAmount("valuation.priceLevels.initialBuyPrice", levels.InitialBuyPrice); err != nil {
		return err
	}
	return validatePositiveAmount("valuation.priceLevels.aggressiveBuyPrice", levels.AggressiveBuyPrice)
}

func validateDividend(dividend *Dividend) error {
	if dividend == nil {
		return nil
	}
	if err := validatePositiveAmount("dividend.dividendPerShare", dividend.DividendPerShare); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("dividend.cashDividendTotal", dividend.CashDividendTotal); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("dividend.buybackAmount", dividend.BuybackAmount); err != nil {
		return err
	}
	if err := validateRatio("dividend.dividendYield", dividend.DividendYield, 5); err != nil {
		return err
	}
	if err := validatePositiveAmount("dividend.forecastPerShare", dividend.ForecastPerShare); err != nil {
		return err
	}
	if err := validateRatio("dividend.forecastYield", dividend.ForecastYield, 5); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("dividend.estimatedAnnualCash", dividend.EstimatedAnnualCash); err != nil {
		return err
	}
	return validateRatio("dividend.payoutRatio", dividend.PayoutRatio, 5)
}

func validateNetCash(netCash *NetCashProfile) error {
	if netCash == nil {
		return nil
	}
	if err := validateNonNegativeAmount("netCash.cashAndShortInvestments", netCash.CashAndShortInvestments); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("netCash.interestBearingDebt", netCash.InterestBearingDebt); err != nil {
		return err
	}
	if err := validateRatio("netCash.haircut", netCash.Haircut, 1); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("netCash.exCashPe", netCash.ExCashPE); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("netCash.exCashPfcf", netCash.ExCashPFCF); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("netCash.shareholderFcf", netCash.ShareholderFCF); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("netCash.consolidatedFcf", netCash.ConsolidatedFCF); err != nil {
		return err
	}
	if err := validateNonNegativeAmount("netCash.minorityFcfAdjustment", netCash.MinorityFCFAdjustment); err != nil {
		return err
	}
	return validateRatio("netCash.fcfYield", netCash.FCFYield, 5)
}

func validateOwnerCashFlowAudit(audit *OwnerCashFlowAudit) error {
	if audit == nil {
		return nil
	}
	for _, item := range ownerAuditItems(audit) {
		status := strings.ToLower(strings.TrimSpace(item.value.Status))
		if status == "" {
			continue
		}
		if status != "pass" && status != "review" && status != "fail" {
			return fmt.Errorf("%s.status must be pass, review, or fail", item.field)
		}
	}
	return nil
}

func validatePositiveAmount(field string, value *float64) error {
	if value == nil {
		return nil
	}
	if *value <= 0 {
		return fmt.Errorf("%s must be positive", field)
	}
	return nil
}

func validateNonNegativeAmount(field string, value *float64) error {
	if value == nil {
		return nil
	}
	if *value < 0 {
		return fmt.Errorf("%s must be non-negative", field)
	}
	return nil
}

func validateRatio(field string, value *float64, max float64) error {
	if value == nil {
		return nil
	}
	if *value < 0 || *value > max {
		return fmt.Errorf("%s must be a decimal ratio between 0 and %.0f", field, max)
	}
	return nil
}

func validateScore(field string, value *float64, max float64) error {
	if value == nil {
		return nil
	}
	if *value < 0 || *value > max {
		return fmt.Errorf("%s must be between 0 and %.0f", field, max)
	}
	return nil
}

func targetBuyPriceFromIntrinsicValue(intrinsicValue *float64) *float64 {
	if intrinsicValue == nil || *intrinsicValue <= 0 {
		return nil
	}
	value := *intrinsicValue * (1 - defaultSafetyMarginTarget)
	return &value
}

func normalizeDividend(dividend *Dividend, fallbackCurrency string) *Dividend {
	if dividend == nil {
		return nil
	}
	hasContent := strings.TrimSpace(dividend.FiscalYear) != "" ||
		dividend.DividendPerShare != nil ||
		strings.TrimSpace(dividend.DividendCurrency) != "" ||
		dividend.CashDividendTotal != nil ||
		strings.TrimSpace(dividend.CashDividendCurrency) != "" ||
		dividend.BuybackAmount != nil ||
		strings.TrimSpace(dividend.BuybackCurrency) != "" ||
		dividend.DividendYield != nil ||
		dividend.PayoutRatio != nil ||
		dividend.EstimatedAnnualCash != nil ||
		strings.TrimSpace(dividend.Reliability) != "" ||
		strings.TrimSpace(dividend.ForecastFiscalYear) != "" ||
		dividend.ForecastPerShare != nil ||
		strings.TrimSpace(dividend.ForecastCurrency) != "" ||
		dividend.ForecastYield != nil
	if !hasContent {
		return nil
	}
	next := &Dividend{
		FiscalYear:           strings.TrimSpace(dividend.FiscalYear),
		DividendPerShare:     cloneFloat(dividend.DividendPerShare),
		DividendCurrency:     strings.ToUpper(strings.TrimSpace(dividend.DividendCurrency)),
		CashDividendTotal:    cloneFloat(dividend.CashDividendTotal),
		CashDividendCurrency: strings.ToUpper(strings.TrimSpace(dividend.CashDividendCurrency)),
		BuybackAmount:        cloneFloat(dividend.BuybackAmount),
		BuybackCurrency:      strings.ToUpper(strings.TrimSpace(dividend.BuybackCurrency)),
		DividendYield:        cloneFloat(dividend.DividendYield),
		PayoutRatio:          cloneFloat(dividend.PayoutRatio),
		EstimatedAnnualCash:  cloneFloat(dividend.EstimatedAnnualCash),
		Reliability:          strings.TrimSpace(dividend.Reliability),
		ForecastFiscalYear:   strings.TrimSpace(dividend.ForecastFiscalYear),
		ForecastPerShare:     cloneFloat(dividend.ForecastPerShare),
		ForecastCurrency:     strings.ToUpper(strings.TrimSpace(dividend.ForecastCurrency)),
		ForecastYield:        cloneFloat(dividend.ForecastYield),
	}
	if next.DividendCurrency == "" {
		next.DividendCurrency = strings.ToUpper(strings.TrimSpace(fallbackCurrency))
	}
	if next.CashDividendCurrency == "" {
		next.CashDividendCurrency = next.DividendCurrency
	}
	if next.BuybackCurrency == "" {
		next.BuybackCurrency = next.CashDividendCurrency
	}
	if next.ForecastCurrency == "" {
		next.ForecastCurrency = strings.ToUpper(strings.TrimSpace(fallbackCurrency))
	}
	if next.FiscalYear == "" &&
		next.DividendPerShare == nil &&
		next.DividendCurrency == "" &&
		next.CashDividendTotal == nil &&
		next.CashDividendCurrency == "" &&
		next.BuybackAmount == nil &&
		next.BuybackCurrency == "" &&
		next.DividendYield == nil &&
		next.PayoutRatio == nil &&
		next.EstimatedAnnualCash == nil &&
		next.Reliability == "" &&
		next.ForecastFiscalYear == "" &&
		next.ForecastPerShare == nil &&
		next.ForecastCurrency == "" &&
		next.ForecastYield == nil {
		return nil
	}
	return next
}

func normalizeNetCash(netCash *NetCashProfile, fallbackCurrency string) *NetCashProfile {
	if netCash == nil {
		return nil
	}
	next := &NetCashProfile{
		CashAndShortInvestments: cloneFloat(netCash.CashAndShortInvestments),
		InterestBearingDebt:     cloneFloat(netCash.InterestBearingDebt),
		NetCash:                 cloneFloat(netCash.NetCash),
		Currency:                strings.ToUpper(strings.TrimSpace(netCash.Currency)),
		Haircut:                 cloneFloat(netCash.Haircut),
		HaircutReason:           strings.TrimSpace(netCash.HaircutReason),
		AdjustedNetCash:         cloneFloat(netCash.AdjustedNetCash),
		ExCashPE:                cloneFloat(netCash.ExCashPE),
		ExCashPFCF:              cloneFloat(netCash.ExCashPFCF),
		FCFYield:                cloneFloat(netCash.FCFYield),
		ShareholderFCF:          cloneFloat(netCash.ShareholderFCF),
		ShareholderFCFCurrency:  strings.ToUpper(strings.TrimSpace(netCash.ShareholderFCFCurrency)),
		ShareholderFCFBasis:     strings.TrimSpace(netCash.ShareholderFCFBasis),
		ConsolidatedFCF:         cloneFloat(netCash.ConsolidatedFCF),
		MinorityFCFAdjustment:   cloneFloat(netCash.MinorityFCFAdjustment),
		FCFPositiveYears:        cloneInt(netCash.FCFPositiveYears),
		Note:                    strings.TrimSpace(netCash.Note),
	}
	if next.Currency == "" {
		next.Currency = strings.ToUpper(strings.TrimSpace(fallbackCurrency))
	}
	if next.ShareholderFCFCurrency == "" {
		next.ShareholderFCFCurrency = next.Currency
	}
	if next.NetCash == nil && next.CashAndShortInvestments != nil && next.InterestBearingDebt != nil {
		value := *next.CashAndShortInvestments - *next.InterestBearingDebt
		next.NetCash = &value
	}
	if next.AdjustedNetCash == nil && next.NetCash != nil && next.Haircut != nil {
		value := *next.NetCash * *next.Haircut
		next.AdjustedNetCash = &value
	}
	if next.CashAndShortInvestments == nil &&
		next.InterestBearingDebt == nil &&
		next.NetCash == nil &&
		next.Haircut == nil &&
		next.HaircutReason == "" &&
		next.AdjustedNetCash == nil &&
		next.ExCashPE == nil &&
		next.ExCashPFCF == nil &&
		next.FCFYield == nil &&
		next.ShareholderFCF == nil &&
		next.ShareholderFCFCurrency == "" &&
		next.ShareholderFCFBasis == "" &&
		next.ConsolidatedFCF == nil &&
		next.MinorityFCFAdjustment == nil &&
		next.FCFPositiveYears == nil &&
		next.Note == "" {
		return nil
	}
	return next
}

func normalizeOwnerCashFlowAudit(audit *OwnerCashFlowAudit) *OwnerCashFlowAudit {
	if audit == nil {
		return nil
	}
	next := &OwnerCashFlowAudit{
		TenYearDemand:                  normalizeOwnerAuditItem(audit.TenYearDemand),
		AssetDurability:                normalizeOwnerAuditItem(audit.AssetDurability),
		MaintenanceCapexLight:          normalizeOwnerAuditItem(audit.MaintenanceCapexLight),
		DividendFCFSupport:             normalizeOwnerAuditItem(audit.DividendFCFSupport),
		DividendReinvestmentEfficiency: normalizeOwnerAuditItem(audit.DividendReinvestmentEfficiency),
		RoeRoicDurability:              normalizeOwnerAuditItem(audit.RoeRoicDurability),
		ValuationSystemRisk:            normalizeOwnerAuditItem(audit.ValuationSystemRisk),
	}
	for _, item := range ownerAuditItems(next) {
		if item.value.Status != "" || item.value.Note != "" {
			return next
		}
	}
	return nil
}

func normalizeOwnerAuditItem(item OwnerAuditItem) OwnerAuditItem {
	return OwnerAuditItem{
		Status: strings.ToLower(strings.TrimSpace(item.Status)),
		Note:   strings.TrimSpace(item.Note),
	}
}

func ownerAuditItems(audit *OwnerCashFlowAudit) []struct {
	field string
	value OwnerAuditItem
} {
	if audit == nil {
		return nil
	}
	return []struct {
		field string
		value OwnerAuditItem
	}{
		{"ownerCashFlowAudit.tenYearDemand", audit.TenYearDemand},
		{"ownerCashFlowAudit.assetDurability", audit.AssetDurability},
		{"ownerCashFlowAudit.maintenanceCapexLight", audit.MaintenanceCapexLight},
		{"ownerCashFlowAudit.dividendFcfSupport", audit.DividendFCFSupport},
		{"ownerCashFlowAudit.dividendReinvestmentEfficiency", audit.DividendReinvestmentEfficiency},
		{"ownerCashFlowAudit.roeRoicDurability", audit.RoeRoicDurability},
		{"ownerCashFlowAudit.valuationSystemRisk", audit.ValuationSystemRisk},
	}
}

func cloneFloat(value *float64) *float64 {
	if value == nil {
		return nil
	}
	next := *value
	return &next
}

func cloneInt(value *int) *int {
	if value == nil {
		return nil
	}
	next := *value
	return &next
}

func expectedCurrency(symbol string) string {
	symbol = normalizeSymbol(symbol)
	switch {
	case strings.HasSuffix(symbol, ".HK"):
		return "HKD"
	case strings.HasSuffix(symbol, ".SH"), strings.HasSuffix(symbol, ".SZ"), strings.HasSuffix(symbol, ".SS"):
		return "CNY"
	default:
		return ""
	}
}

func marginOfSafetyFromPrice(intrinsicValue *float64, currentPrice float64, fallback *float64) *float64 {
	if intrinsicValue == nil || *intrinsicValue <= 0 || currentPrice <= 0 {
		return fallback
	}
	value := (*intrinsicValue - currentPrice) / *intrinsicValue
	return &value
}

func prefer(current, next string) string {
	if strings.TrimSpace(next) != "" {
		return strings.ToUpper(strings.TrimSpace(next))
	}
	return strings.ToUpper(strings.TrimSpace(current))
}

func normalizeSymbol(symbol string) string {
	return normalizeDisplaySymbol(symbol)
}

func normalizeDisplaySymbol(symbol string) string {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if strings.HasPrefix(symbol, "HK") {
		code := strings.TrimPrefix(symbol, "HK")
		if value, err := strconv.Atoi(code); err == nil {
			return fmt.Sprintf("%04d.HK", value)
		}
	}
	if strings.HasSuffix(symbol, ".HK") {
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconv.Atoi(code); err == nil {
			return fmt.Sprintf("%04d.HK", value)
		}
	}
	return symbol
}
