package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type AppState struct {
	TotalCapital float64            `json:"totalCapital"`
	Cash         float64            `json:"cash"`
	FX           map[string]float64 `json:"fx"`
	Trades       []Trade            `json:"trades"`
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

type Holding struct {
	Symbol            string   `json:"symbol"`
	Name              string   `json:"name"`
	Shares            float64  `json:"shares"`
	Cost              float64  `json:"cost"`
	CurrentPrice      float64  `json:"currentPrice"`
	PreviousClose     float64  `json:"previousClose"`
	CurrentPriceDate  string   `json:"currentPriceDate"`
	PreviousCloseDate string   `json:"previousCloseDate"`
	Action            string   `json:"action"`
	Status            string   `json:"status"`
	MarginOfSafety    *float64 `json:"marginOfSafety"`
	QualityScore      *float64 `json:"qualityScore"`
	Risk              string   `json:"risk"`
	Industry          string   `json:"industry"`
	Currency          string   `json:"currency"`
	IntrinsicValue    *float64 `json:"intrinsicValue"`
	FairValueRange    string   `json:"fairValueRange"`
	TargetBuyPrice    *float64 `json:"targetBuyPrice"`
	BusinessModel     *float64 `json:"businessModel"`
	Moat              *float64 `json:"moat"`
	Governance        *float64 `json:"governance"`
	FinancialQuality  *float64 `json:"financialQuality"`
	UpdatedAt         string   `json:"updatedAt"`
	Notes             string   `json:"notes"`
}

type PlanItem struct {
	Rank       int    `json:"rank"`
	Name       string `json:"name"`
	Priority   string `json:"priority"`
	Advice     string `json:"advice"`
	Discipline string `json:"discipline"`
}

type Candidate struct {
	Symbol         string   `json:"symbol"`
	Name           string   `json:"name"`
	Status         string   `json:"status"`
	Action         string   `json:"action"`
	MarginOfSafety *float64 `json:"marginOfSafety"`
	QualityScore   *float64 `json:"qualityScore"`
	Industry       string   `json:"industry"`
	Currency       string   `json:"currency"`
	IntrinsicValue *float64 `json:"intrinsicValue"`
	FairValueRange string   `json:"fairValueRange"`
	TargetBuyPrice *float64 `json:"targetBuyPrice"`
}

type Rule struct {
	Dimension string  `json:"dimension"`
	Score     float64 `json:"score"`
	Standard  string  `json:"standard"`
}

type ResearchImport struct {
	Symbol    string    `json:"symbol"`
	Name      string    `json:"name"`
	AsOf      string    `json:"asOf"`
	Currency  string    `json:"currency"`
	Industry  string    `json:"industry"`
	Status    string    `json:"status"`
	Action    string    `json:"action"`
	Risk      string    `json:"risk"`
	Valuation Valuation `json:"valuation"`
	Quality   Quality   `json:"quality"`
	Plan      PlanInput `json:"plan"`
	Notes     string    `json:"notes"`
}

type Valuation struct {
	IntrinsicValue *float64 `json:"intrinsicValue"`
	FairValueRange string   `json:"fairValueRange"`
	TargetBuyPrice *float64 `json:"targetBuyPrice"`
	MarginOfSafety *float64 `json:"marginOfSafety"`
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
			applyCandidateResearch(&state.Candidates[i], research)
			upsertPlan(state, research)
			return fmt.Sprintf("updated candidate %s (%s)", research.Symbol, research.Name)
		}
	}

	state.Candidates = append(state.Candidates, Candidate{Symbol: normalizeDisplaySymbol(research.Symbol)})
	applyCandidateResearch(&state.Candidates[len(state.Candidates)-1], research)
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
	holding.MarginOfSafety = research.Valuation.MarginOfSafety
	holding.QualityScore = research.Quality.TotalScore
	holding.IntrinsicValue = research.Valuation.IntrinsicValue
	holding.FairValueRange = strings.TrimSpace(research.Valuation.FairValueRange)
	holding.TargetBuyPrice = research.Valuation.TargetBuyPrice
	holding.BusinessModel = research.Quality.BusinessModel
	holding.Moat = research.Quality.Moat
	holding.Governance = research.Quality.Governance
	holding.FinancialQuality = research.Quality.FinancialQuality
	holding.UpdatedAt = updateLabel
	holding.Notes = strings.TrimSpace(research.Notes)
}

func applyCandidateResearch(candidate *Candidate, research ResearchImport) {
	candidate.Symbol = normalizeDisplaySymbol(research.Symbol)
	candidate.Name = strings.TrimSpace(research.Name)
	candidate.Status = strings.TrimSpace(research.Status)
	candidate.Action = strings.TrimSpace(research.Action)
	candidate.Industry = strings.TrimSpace(research.Industry)
	candidate.Currency = prefer(candidate.Currency, research.Currency)
	candidate.MarginOfSafety = research.Valuation.MarginOfSafety
	candidate.QualityScore = research.Quality.TotalScore
	candidate.IntrinsicValue = research.Valuation.IntrinsicValue
	candidate.FairValueRange = strings.TrimSpace(research.Valuation.FairValueRange)
	candidate.TargetBuyPrice = research.Valuation.TargetBuyPrice
}

func upsertPlan(state *AppState, research ResearchImport) {
	if strings.TrimSpace(research.Plan.Priority) == "" &&
		strings.TrimSpace(research.Plan.Advice) == "" &&
		strings.TrimSpace(research.Plan.Discipline) == "" {
		return
	}

	next := PlanItem{
		Rank:       research.Plan.Rank,
		Name:       strings.TrimSpace(research.Name),
		Priority:   strings.TrimSpace(research.Plan.Priority),
		Advice:     strings.TrimSpace(research.Plan.Advice),
		Discipline: strings.TrimSpace(research.Plan.Discipline),
	}
	if next.Rank <= 0 {
		next.Rank = nextPlanRank(state.Plan)
	}

	for i := range state.Plan {
		if strings.EqualFold(state.Plan[i].Name, next.Name) {
			state.Plan[i] = next
			sortPlan(state.Plan)
			return
		}
	}

	state.Plan = append(state.Plan, next)
	sortPlan(state.Plan)
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

func sortPlan(plan []PlanItem) {
	sort.SliceStable(plan, func(i, j int) bool {
		return plan[i].Rank < plan[j].Rank
	})
}

func prefer(current, next string) string {
	if strings.TrimSpace(next) != "" {
		return strings.ToUpper(strings.TrimSpace(next))
	}
	return strings.ToUpper(strings.TrimSpace(current))
}

func normalizeSymbol(symbol string) string {
	return strings.ToUpper(strings.TrimSpace(symbol))
}

func normalizeDisplaySymbol(symbol string) string {
	return normalizeSymbol(symbol)
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
