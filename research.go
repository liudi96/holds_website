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

type ResearchResponse struct {
	Summary    string         `json:"summary"`
	TargetType string         `json:"targetType"`
	Warnings   []string       `json:"warnings"`
	Research   ResearchImport `json:"research"`
	Plan       []PlanItem     `json:"plan"`
	BackupPath string         `json:"backupPath,omitempty"`
	State      *AppState      `json:"state,omitempty"`
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

	summary, targetType := applyResearch(&state, research)
	writeJSON(w, http.StatusOK, ResearchResponse{
		Summary:    summary,
		TargetType: targetType,
		Warnings:   warnings,
		Research:   normalizeResearch(research),
		Plan:       state.Plan,
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

	summary, targetType := applyResearch(&state, research)
	appendResearchDecisionLog(&state, research, summary, targetType)
	backupPath, err := backupPortfolioFile()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to backup portfolio data")
		return
	}
	if err := saveState(state); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save state")
		return
	}

	s.state = state
	writeJSON(w, http.StatusOK, ResearchResponse{
		Summary:    summary,
		TargetType: targetType,
		Warnings:   warnings,
		Research:   normalizeResearch(research),
		Plan:       state.Plan,
		BackupPath: backupPath,
		State:      &state,
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

	asOfDate, err := time.Parse("2006-01-02", research.AsOf)
	if err != nil {
		return nil, fmt.Errorf("asOf must be YYYY-MM-DD: %w", err)
	}
	today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if asOfDate.After(today) {
		return nil, fmt.Errorf("asOf cannot be later than today: %s", today.Format("2006-01-02"))
	}

	if err := validatePercent("valuation.marginOfSafety", research.Valuation.MarginOfSafety); err != nil {
		return nil, err
	}
	if err := validateDividend(research.Dividend); err != nil {
		return nil, err
	}
	if err := validateNetCash(research.NetCash); err != nil {
		return nil, err
	}
	if err := validateOwnerCashFlowAudit(research.OwnerCashFlowAudit); err != nil {
		return nil, err
	}
	if err := validateScore("quality.totalScore", research.Quality.TotalScore, 100); err != nil {
		return nil, err
	}
	if err := validateScore("quality.businessModel", research.Quality.BusinessModel, 30); err != nil {
		return nil, err
	}
	if err := validateScore("quality.moat", research.Quality.Moat, 25); err != nil {
		return nil, err
	}
	if err := validateScore("quality.governance", research.Quality.Governance, 20); err != nil {
		return nil, err
	}
	if err := validateScore("quality.financialQuality", research.Quality.FinancialQuality, 25); err != nil {
		return nil, err
	}

	if expected := expectedCurrency(research.Symbol); expected != "" && research.Currency != "" && research.Currency != expected {
		warnings = append(warnings, fmt.Sprintf("currency is %s, but %s usually uses %s", research.Currency, research.Symbol, expected))
	}
	if hasAllQualityScores(research.Quality) {
		sum := *research.Quality.BusinessModel + *research.Quality.Moat + *research.Quality.Governance + *research.Quality.FinancialQuality
		if math.Abs(*research.Quality.TotalScore-sum) > 0.01 {
			warnings = append(warnings, fmt.Sprintf("quality.totalScore is %.2f, but component scores add up to %.2f", *research.Quality.TotalScore, sum))
		}
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

func applyResearch(state *AppState, research ResearchImport) (string, string) {
	symbol := normalizeSymbol(research.Symbol)
	now := time.Now().Format("2006-01-02 15:04:05")
	updateLabel := fmt.Sprintf("%s；ChatGPT分析导入；分析日 %s", now, research.AsOf)

	for i := range state.Holdings {
		if normalizeSymbol(state.Holdings[i].Symbol) == symbol {
			applyHoldingResearch(&state.Holdings[i], research, updateLabel)
			upsertPlan(state, research)
			return fmt.Sprintf("updated holding %s (%s)", research.Symbol, research.Name), "holding"
		}
	}

	for i := range state.Candidates {
		if normalizeSymbol(state.Candidates[i].Symbol) == symbol {
			applyCandidateResearch(&state.Candidates[i], research, updateLabel)
			upsertPlan(state, research)
			return fmt.Sprintf("updated candidate %s (%s)", research.Symbol, research.Name), "candidate"
		}
	}

	state.Candidates = append(state.Candidates, Candidate{Symbol: normalizeDisplaySymbol(research.Symbol)})
	applyCandidateResearch(&state.Candidates[len(state.Candidates)-1], research, updateLabel)
	upsertPlan(state, research)
	return fmt.Sprintf("added candidate %s (%s)", research.Symbol, research.Name), "newCandidate"
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
	if err := os.WriteFile(backupPath, body, 0o644); err != nil {
		return "", err
	}
	return backupPath, nil
}

func normalizeResearch(research ResearchImport) ResearchImport {
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
	return research
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
	if err := validatePositiveAmount("dividend.forecastPerShare", dividend.ForecastPerShare); err != nil {
		return err
	}
	if err := validateRatio("dividend.forecastYield", dividend.ForecastYield, 5); err != nil {
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
	next := &Dividend{
		FiscalYear:           strings.TrimSpace(dividend.FiscalYear),
		DividendPerShare:     cloneFloat(dividend.DividendPerShare),
		DividendCurrency:     strings.ToUpper(strings.TrimSpace(dividend.DividendCurrency)),
		CashDividendTotal:    nil,
		CashDividendCurrency: "",
		BuybackAmount:        nil,
		BuybackCurrency:      "",
		DividendYield:        nil,
		PayoutRatio:          cloneFloat(dividend.PayoutRatio),
		EstimatedAnnualCash:  nil,
		Reliability:          strings.TrimSpace(dividend.Reliability),
		ForecastFiscalYear:   strings.TrimSpace(dividend.ForecastFiscalYear),
		ForecastPerShare:     cloneFloat(dividend.ForecastPerShare),
		ForecastCurrency:     strings.ToUpper(strings.TrimSpace(dividend.ForecastCurrency)),
		ForecastYield:        cloneFloat(dividend.ForecastYield),
	}
	if next.DividendCurrency == "" {
		next.DividendCurrency = strings.ToUpper(strings.TrimSpace(fallbackCurrency))
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
		FCFPositiveYears:        cloneInt(netCash.FCFPositiveYears),
		Note:                    strings.TrimSpace(netCash.Note),
	}
	if next.Currency == "" {
		next.Currency = strings.ToUpper(strings.TrimSpace(fallbackCurrency))
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
	if strings.HasSuffix(symbol, ".HK") {
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconv.Atoi(code); err == nil {
			return fmt.Sprintf("%04d.HK", value)
		}
	}
	return symbol
}
