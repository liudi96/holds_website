package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	decisionLogLimit          = 500
	defaultSafetyMarginTarget = 0.25
)

type AppState struct {
	TotalCapital float64            `json:"totalCapital"`
	Cash         float64            `json:"cash"`
	FX           map[string]float64 `json:"fx"`
	Trades       []Trade            `json:"trades"`
	DecisionLogs []DecisionLog      `json:"decisionLogs"`
	Holdings     []Holding          `json:"holdings"`
	Plan         []PlanItem         `json:"plan"`
	Candidates   []Candidate        `json:"candidates"`
	Rules        []Rule             `json:"rules"`
}

type Trade struct {
	ID           int64   `json:"id"`
	Date         string  `json:"date"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	Side         string  `json:"side"`
	Shares       float64 `json:"shares"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	CurrentPrice float64 `json:"currentPrice"`
}

type DecisionLog struct {
	ID         int64    `json:"id"`
	Date       string   `json:"date"`
	Type       string   `json:"type"`
	Symbol     string   `json:"symbol,omitempty"`
	Name       string   `json:"name,omitempty"`
	Price      *float64 `json:"price,omitempty"`
	Currency   string   `json:"currency,omitempty"`
	Decision   string   `json:"decision"`
	Discipline string   `json:"discipline"`
	Detail     string   `json:"detail,omitempty"`
}

type Holding struct {
	Symbol              string              `json:"symbol"`
	Name                string              `json:"name"`
	Shares              float64             `json:"shares"`
	Cost                float64             `json:"cost"`
	CurrentPrice        float64             `json:"currentPrice"`
	PreviousClose       float64             `json:"previousClose"`
	MarketCap           *float64            `json:"marketCap,omitempty"`
	MarketCapCurrency   string              `json:"marketCapCurrency,omitempty"`
	CurrentPriceDate    string              `json:"currentPriceDate"`
	PreviousCloseDate   string              `json:"previousCloseDate"`
	Action              string              `json:"action"`
	Status              string              `json:"status"`
	MarginOfSafety      *float64            `json:"marginOfSafety"`
	QualityScore        *float64            `json:"qualityScore"`
	Risk                string              `json:"risk"`
	Industry            string              `json:"industry"`
	Currency            string              `json:"currency"`
	IntrinsicValue      *float64            `json:"intrinsicValue"`
	FairValueRange      string              `json:"fairValueRange"`
	TargetBuyPrice      *float64            `json:"targetBuyPrice"`
	PriceLevels         *PriceLevels        `json:"priceLevels,omitempty"`
	ValuationConfidence string              `json:"valuationConfidence,omitempty"`
	BusinessModel       *float64            `json:"businessModel"`
	Moat                *float64            `json:"moat"`
	Governance          *float64            `json:"governance"`
	FinancialQuality    *float64            `json:"financialQuality"`
	UpdatedAt           string              `json:"updatedAt"`
	Notes               string              `json:"notes"`
	KillCriteria        json.RawMessage     `json:"killCriteria,omitempty"`
	Reports             []Report            `json:"reports,omitempty"`
	Dividend            *Dividend           `json:"dividend,omitempty"`
	NetCash             *NetCashProfile     `json:"netCash,omitempty"`
	OwnerCashFlowAudit  *OwnerCashFlowAudit `json:"ownerCashFlowAudit,omitempty"`
	Financials          json.RawMessage     `json:"financials,omitempty"`
}

type PriceLevels struct {
	WatchPrice         *float64 `json:"watchPrice,omitempty"`
	InitialBuyPrice    *float64 `json:"initialBuyPrice,omitempty"`
	AggressiveBuyPrice *float64 `json:"aggressiveBuyPrice,omitempty"`
}

type Dividend struct {
	FiscalYear           string   `json:"fiscalYear,omitempty"`
	DividendPerShare     *float64 `json:"dividendPerShare,omitempty"`
	DividendCurrency     string   `json:"dividendCurrency,omitempty"`
	CashDividendTotal    *float64 `json:"cashDividendTotal,omitempty"`
	CashDividendCurrency string   `json:"cashDividendCurrency,omitempty"`
	BuybackAmount        *float64 `json:"buybackAmount,omitempty"`
	BuybackCurrency      string   `json:"buybackCurrency,omitempty"`
	DividendYield        *float64 `json:"dividendYield,omitempty"`
	PayoutRatio          *float64 `json:"payoutRatio,omitempty"`
	EstimatedAnnualCash  *float64 `json:"estimatedAnnualCash,omitempty"`
	Reliability          string   `json:"reliability,omitempty"`
	ForecastFiscalYear   string   `json:"forecastFiscalYear,omitempty"`
	ForecastPerShare     *float64 `json:"forecastPerShare,omitempty"`
	ForecastCurrency     string   `json:"forecastCurrency,omitempty"`
	ForecastYield        *float64 `json:"forecastYield,omitempty"`
}

type NetCashProfile struct {
	CashAndShortInvestments *float64 `json:"cashAndShortInvestments,omitempty"`
	InterestBearingDebt     *float64 `json:"interestBearingDebt,omitempty"`
	NetCash                 *float64 `json:"netCash,omitempty"`
	Currency                string   `json:"currency,omitempty"`
	Haircut                 *float64 `json:"haircut,omitempty"`
	HaircutReason           string   `json:"haircutReason,omitempty"`
	AdjustedNetCash         *float64 `json:"adjustedNetCash,omitempty"`
	ExCashPE                *float64 `json:"exCashPe,omitempty"`
	ExCashPFCF              *float64 `json:"exCashPfcf,omitempty"`
	FCFYield                *float64 `json:"fcfYield,omitempty"`
	ShareholderFCF          *float64 `json:"shareholderFcf,omitempty"`
	ShareholderFCFCurrency  string   `json:"shareholderFcfCurrency,omitempty"`
	ShareholderFCFBasis     string   `json:"shareholderFcfBasis,omitempty"`
	ConsolidatedFCF         *float64 `json:"consolidatedFcf,omitempty"`
	MinorityFCFAdjustment   *float64 `json:"minorityFcfAdjustment,omitempty"`
	FCFPositiveYears        *int     `json:"fcfPositiveYears,omitempty"`
	Note                    string   `json:"note,omitempty"`
}

type OwnerCashFlowAudit struct {
	TenYearDemand                  OwnerAuditItem `json:"tenYearDemand,omitempty"`
	AssetDurability                OwnerAuditItem `json:"assetDurability,omitempty"`
	MaintenanceCapexLight          OwnerAuditItem `json:"maintenanceCapexLight,omitempty"`
	DividendFCFSupport             OwnerAuditItem `json:"dividendFcfSupport,omitempty"`
	DividendReinvestmentEfficiency OwnerAuditItem `json:"dividendReinvestmentEfficiency,omitempty"`
	RoeRoicDurability              OwnerAuditItem `json:"roeRoicDurability,omitempty"`
	ValuationSystemRisk            OwnerAuditItem `json:"valuationSystemRisk,omitempty"`
}

type OwnerAuditItem struct {
	Status string `json:"status,omitempty"`
	Note   string `json:"note,omitempty"`
}

type Report struct {
	Period string `json:"period"`
	Kind   string `json:"kind"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Source string `json:"source"`
	URL    string `json:"url"`
}

type PlanItem struct {
	Rank       int    `json:"rank"`
	Symbol     string `json:"symbol,omitempty"`
	Name       string `json:"name"`
	Priority   string `json:"priority"`
	Advice     string `json:"advice"`
	Discipline string `json:"discipline"`
}

type Candidate struct {
	Symbol              string              `json:"symbol"`
	Name                string              `json:"name"`
	Status              string              `json:"status"`
	Action              string              `json:"action"`
	CurrentPrice        float64             `json:"currentPrice"`
	PreviousClose       float64             `json:"previousClose"`
	MarketCap           *float64            `json:"marketCap,omitempty"`
	MarketCapCurrency   string              `json:"marketCapCurrency,omitempty"`
	CurrentPriceDate    string              `json:"currentPriceDate"`
	PreviousCloseDate   string              `json:"previousCloseDate"`
	MarginOfSafety      *float64            `json:"marginOfSafety"`
	QualityScore        *float64            `json:"qualityScore"`
	Risk                string              `json:"risk"`
	Industry            string              `json:"industry"`
	Currency            string              `json:"currency"`
	IntrinsicValue      *float64            `json:"intrinsicValue"`
	FairValueRange      string              `json:"fairValueRange"`
	TargetBuyPrice      *float64            `json:"targetBuyPrice"`
	PriceLevels         *PriceLevels        `json:"priceLevels,omitempty"`
	ValuationConfidence string              `json:"valuationConfidence,omitempty"`
	BusinessModel       *float64            `json:"businessModel"`
	Moat                *float64            `json:"moat"`
	Governance          *float64            `json:"governance"`
	FinancialQuality    *float64            `json:"financialQuality"`
	UpdatedAt           string              `json:"updatedAt"`
	Notes               string              `json:"notes"`
	KillCriteria        json.RawMessage     `json:"killCriteria,omitempty"`
	Reports             []Report            `json:"reports,omitempty"`
	Dividend            *Dividend           `json:"dividend,omitempty"`
	NetCash             *NetCashProfile     `json:"netCash,omitempty"`
	OwnerCashFlowAudit  *OwnerCashFlowAudit `json:"ownerCashFlowAudit,omitempty"`
	Financials          json.RawMessage     `json:"financials,omitempty"`
}

type Rule struct {
	Dimension string  `json:"dimension"`
	Score     float64 `json:"score"`
	Standard  string  `json:"standard"`
}

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

func main() {
	dataPath := flag.String("data", "data/portfolio.json", "portfolio JSON file to update")
	dryRun := flag.Bool("dry-run", false, "validate and print changes without writing")
	flag.Parse()

	if flag.NArg() != 1 {
		fail(errors.New("usage: go run ./cmd/import-research [-data data/portfolio.json] path/to/research.json"))
	}

	research, err := loadResearch(flag.Arg(0))
	if err != nil {
		fail(err)
	}
	if err := validateResearch(research); err != nil {
		fail(err)
	}

	state, err := loadState(*dataPath)
	if err != nil {
		fail(err)
	}

	summary := applyResearch(&state, research)
	appendResearchDecisionLog(&state, research, summary)
	if *dryRun {
		fmt.Println(summary)
		fmt.Println("dry run: portfolio data was not changed")
		return
	}

	if err := saveState(*dataPath, state); err != nil {
		fail(err)
	}
	fmt.Println(summary)
}

func loadResearch(path string) (ResearchImport, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return ResearchImport{}, err
	}

	var research ResearchImport
	if err := json.Unmarshal(body, &research); err != nil {
		return ResearchImport{}, err
	}
	return research, nil
}

func validateResearch(research ResearchImport) error {
	if strings.TrimSpace(research.Symbol) == "" {
		return errors.New("symbol is required")
	}
	if strings.TrimSpace(research.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(research.AsOf) == "" {
		return errors.New("asOf is required")
	}
	if _, err := time.Parse("2006-01-02", research.AsOf); err != nil {
		return fmt.Errorf("asOf must be YYYY-MM-DD: %w", err)
	}
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
	if err := validateScore("quality.totalScore", research.Quality.TotalScore, 100); err != nil {
		return err
	}
	if err := validateScore("quality.businessModel", research.Quality.BusinessModel, 30); err != nil {
		return err
	}
	if err := validateScore("quality.moat", research.Quality.Moat, 25); err != nil {
		return err
	}
	if err := validateScore("quality.governance", research.Quality.Governance, 20); err != nil {
		return err
	}
	return validateScore("quality.financialQuality", research.Quality.FinancialQuality, 25)
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

func applyResearch(state *AppState, research ResearchImport) string {
	research = normalizeResearch(research)
	symbol := normalizeSymbol(research.Symbol)
	now := time.Now().Format("2006-01-02 15:04:05")
	updateLabel := fmt.Sprintf("%s；ChatGPT分析导入；分析日 %s", now, research.AsOf)

	for i := range state.Holdings {
		if normalizeSymbol(state.Holdings[i].Symbol) == symbol {
			applyHoldingResearch(&state.Holdings[i], research, updateLabel)
			upsertPlan(state, research)
			return fmt.Sprintf("updated holding %s (%s)", research.Symbol, research.Name)
		}
	}

	for i := range state.Candidates {
		if normalizeSymbol(state.Candidates[i].Symbol) == symbol {
			applyCandidateResearch(&state.Candidates[i], research, updateLabel)
			upsertPlan(state, research)
			return fmt.Sprintf("updated candidate %s (%s)", research.Symbol, research.Name)
		}
	}

	state.Candidates = append(state.Candidates, Candidate{Symbol: normalizeDisplaySymbol(research.Symbol)})
	applyCandidateResearch(&state.Candidates[len(state.Candidates)-1], research, updateLabel)
	upsertPlan(state, research)
	return fmt.Sprintf("added candidate %s (%s)", research.Symbol, research.Name)
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

func prefer(current, next string) string {
	if strings.TrimSpace(next) != "" {
		return strings.ToUpper(strings.TrimSpace(next))
	}
	return strings.ToUpper(strings.TrimSpace(current))
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

func marginOfSafetyFromPrice(intrinsicValue *float64, currentPrice float64, fallback *float64) *float64 {
	if intrinsicValue == nil || *intrinsicValue <= 0 || currentPrice <= 0 {
		return fallback
	}
	value := (*intrinsicValue - currentPrice) / *intrinsicValue
	return &value
}

func appendResearchDecisionLog(state *AppState, research ResearchImport, summary string) {
	name, price, currency, decision, discipline := decisionLogContext(state, research.Symbol)
	appendDecisionLog(state, DecisionLog{
		Type:       "research",
		Symbol:     research.Symbol,
		Name:       firstNonEmpty(name, research.Name),
		Price:      price,
		Currency:   firstNonEmpty(currency, research.Currency),
		Decision:   firstNonEmpty(decision, research.Action, research.Status),
		Discipline: firstNonEmpty(discipline, research.Plan.Discipline, research.Status),
		Detail:     summary,
	})
}

func appendDecisionLog(state *AppState, entry DecisionLog) {
	entry.Type = strings.TrimSpace(entry.Type)
	if entry.Type == "" {
		entry.Type = "event"
	}
	entry.Symbol = normalizeSymbol(entry.Symbol)
	entry.Name = strings.TrimSpace(entry.Name)
	entry.Currency = normalizeSymbol(entry.Currency)
	entry.Decision = strings.TrimSpace(entry.Decision)
	entry.Discipline = strings.TrimSpace(entry.Discipline)
	entry.Detail = strings.TrimSpace(entry.Detail)
	if entry.ID == 0 {
		entry.ID = time.Now().UnixNano()
	}
	if strings.TrimSpace(entry.Date) == "" {
		entry.Date = time.Now().Format("2006-01-02 15:04:05")
	}

	state.DecisionLogs = append(state.DecisionLogs, entry)
	if len(state.DecisionLogs) > decisionLogLimit {
		state.DecisionLogs = state.DecisionLogs[len(state.DecisionLogs)-decisionLogLimit:]
	}
}

func decisionLogContext(state *AppState, symbol string) (string, *float64, string, string, string) {
	normalizedSymbol := normalizeSymbol(symbol)
	for i := range state.Holdings {
		holding := state.Holdings[i]
		if normalizeSymbol(holding.Symbol) != normalizedSymbol {
			continue
		}
		plan := findPlanForDecisionLog(state, holding.Symbol, holding.Name)
		return holding.Name, pricePointer(holding.CurrentPrice), holding.Currency, firstNonEmpty(holding.Action, holding.Status), firstNonEmpty(planDiscipline(plan), holding.Status)
	}

	for i := range state.Candidates {
		candidate := state.Candidates[i]
		if normalizeSymbol(candidate.Symbol) != normalizedSymbol {
			continue
		}
		plan := findPlanForDecisionLog(state, candidate.Symbol, candidate.Name)
		return candidate.Name, pricePointer(candidate.CurrentPrice), candidate.Currency, firstNonEmpty(candidate.Action, candidate.Status), firstNonEmpty(planDiscipline(plan), candidate.Status)
	}

	plan := findPlanForDecisionLog(state, symbol, "")
	return "", nil, "", "", planDiscipline(plan)
}

func findPlanForDecisionLog(state *AppState, symbol string, name string) *PlanItem {
	normalizedSymbol := normalizeSymbol(symbol)
	normalizedName := strings.TrimSpace(name)
	for i := range state.Plan {
		itemSymbol := normalizeSymbol(state.Plan[i].Symbol)
		if itemSymbol != "" && normalizedSymbol != "" && itemSymbol == normalizedSymbol {
			return &state.Plan[i]
		}
		itemName := strings.TrimSpace(state.Plan[i].Name)
		if itemName != "" && normalizedName != "" && (strings.EqualFold(itemName, normalizedName) || strings.Contains(normalizedName, itemName) || strings.Contains(itemName, normalizedName)) {
			return &state.Plan[i]
		}
	}
	return nil
}

func planDiscipline(plan *PlanItem) string {
	if plan == nil {
		return ""
	}
	return strings.TrimSpace(plan.Discipline)
}

func pricePointer(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if text := strings.TrimSpace(value); text != "" {
			return text
		}
	}
	return ""
}

func loadState(path string) (AppState, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return AppState{}, err
	}

	var state AppState
	if err := json.Unmarshal(body, &state); err != nil {
		return AppState{}, err
	}
	return state, nil
}

func saveState(path string, state AppState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(path, body, 0o644)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
