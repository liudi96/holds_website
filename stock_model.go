package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

type Stock struct {
	Symbol              string                `json:"symbol"`
	Name                string                `json:"name"`
	Status              string                `json:"status,omitempty"`
	Action              string                `json:"action,omitempty"`
	CurrentPrice        float64               `json:"currentPrice,omitempty"`
	PreviousClose       float64               `json:"previousClose,omitempty"`
	TwentyDayClose      float64               `json:"twentyDayClose,omitempty"`
	TwentyDayCloseDate  string                `json:"twentyDayCloseDate,omitempty"`
	TwentyDayChange     *float64              `json:"twentyDayChange,omitempty"`
	MarketCap           *float64              `json:"marketCap,omitempty"`
	MarketCapCurrency   string                `json:"marketCapCurrency,omitempty"`
	CurrentPriceDate    string                `json:"currentPriceDate,omitempty"`
	PreviousCloseDate   string                `json:"previousCloseDate,omitempty"`
	MarginOfSafety      *float64              `json:"marginOfSafety,omitempty"`
	QualityScore        *float64              `json:"qualityScore,omitempty"`
	Risk                string                `json:"risk,omitempty"`
	Industry            string                `json:"industry,omitempty"`
	Category            string                `json:"category,omitempty"`
	Currency            string                `json:"currency,omitempty"`
	IntrinsicValue      *float64              `json:"intrinsicValue,omitempty"`
	FairValueRange      string                `json:"fairValueRange,omitempty"`
	TargetBuyPrice      *float64              `json:"targetBuyPrice,omitempty"`
	PriceLevels         *PriceLevels          `json:"priceLevels,omitempty"`
	ValuationConfidence string                `json:"valuationConfidence,omitempty"`
	BuyLogic           string                `json:"buyLogic,omitempty"`
	BusinessModel       *float64              `json:"businessModel,omitempty"`
	Moat                *float64              `json:"moat,omitempty"`
	Governance          *float64              `json:"governance,omitempty"`
	FinancialQuality    *float64              `json:"financialQuality,omitempty"`
	UpdatedAt           string                `json:"updatedAt,omitempty"`
	Notes               string                `json:"notes,omitempty"`
	KillCriteria        json.RawMessage       `json:"killCriteria,omitempty"`
	Reports             []Report              `json:"reports,omitempty"`
	Dividend            *Dividend             `json:"dividend,omitempty"`
	NetCash             *NetCashProfile       `json:"netCash,omitempty"`
	OwnerCashFlowAudit  *OwnerCashFlowAudit   `json:"ownerCashFlowAudit,omitempty"`
	ResearchUpdates     []ResearchUpdate      `json:"researchUpdates,omitempty"`
	Financials          *Financials           `json:"financials,omitempty"`
	Position            *StockPosition        `json:"position,omitempty"`
	Valuation           *ValuationAssumptions `json:"valuation,omitempty"`
	Screening           *ScreeningSummary     `json:"screening,omitempty"`
}

type StockPosition struct {
	Shares float64 `json:"shares"`
	Cost   float64 `json:"cost"`
}

type ScreeningWeights struct {
	Quality           int `json:"quality"`
	CashFlow          int `json:"cashFlow"`
	Valuation         int `json:"valuation"`
	ShareholderReturn int `json:"shareholderReturn"`
	Growth            int `json:"growth"`
}

type ScreeningSummary struct {
	HardBlocks []string             `json:"hardBlocks,omitempty"`
	Score      *float64             `json:"score,omitempty"`
	Subscores  map[string]float64   `json:"subscores,omitempty"`
	Reasons    []string             `json:"reasons,omitempty"`
	Sources    []ScreeningDataPoint `json:"sources,omitempty"`
}

type ScreeningDataPoint struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	Source     string `json:"source,omitempty"`
	AsOf       string `json:"asOf,omitempty"`
	Confidence string `json:"confidence,omitempty"`
}

type ValuationAssumptions struct {
	Currency     string              `json:"currency,omitempty"`
	CurrentPrice float64             `json:"currentPrice,omitempty"`
	RequiredMargin *float64          `json:"requiredMargin,omitempty"`
	UpdatedAt    string              `json:"updatedAt,omitempty"`
	Source       string              `json:"source,omitempty"`
	Scenarios    []ValuationScenario `json:"scenarios,omitempty"`
	Range        *ValuationRange     `json:"range,omitempty"`
}

type ValuationScenario struct {
	Name           string  `json:"name"`
	RevenueGrowth  float64 `json:"revenueGrowth,omitempty"`
	ProfitMargin   float64 `json:"profitMargin,omitempty"`
	FCF            float64 `json:"fcf,omitempty"`
	DiscountRate   float64 `json:"discountRate,omitempty"`
	ReasonablePE   float64 `json:"reasonablePe,omitempty"`
	ReasonablePFCF float64 `json:"reasonablePfcf,omitempty"`
	Shares         float64 `json:"shares,omitempty"`
	FairValue      float64 `json:"fairValue,omitempty"`
	Source         string  `json:"source,omitempty"`
	AsOf           string  `json:"asOf,omitempty"`
}

type ValuationRange struct {
	Low            float64  `json:"low"`
	Base           float64  `json:"base"`
	High           float64  `json:"high"`
	Currency       string   `json:"currency,omitempty"`
	MarginOfSafety *float64 `json:"marginOfSafety,omitempty"`
}

type ValuationHistoryPoint struct {
	Date  string   `json:"date"`
	Price *float64 `json:"price,omitempty"`
	PE    *float64 `json:"pe,omitempty"`
	PB    *float64 `json:"pb,omitempty"`
}

func DefaultScreeningWeights() ScreeningWeights {
	return ScreeningWeights{Quality: 30, CashFlow: 25, Valuation: 20, ShareholderReturn: 15, Growth: 10}
}

func (weights ScreeningWeights) Validate() error {
	total := weights.Quality + weights.CashFlow + weights.Valuation + weights.ShareholderReturn + weights.Growth
	if total != 100 {
		return fmt.Errorf("screening weights must sum to 100, got %d", total)
	}
	if weights.Quality < 0 || weights.CashFlow < 0 || weights.Valuation < 0 || weights.ShareholderReturn < 0 || weights.Growth < 0 {
		return errors.New("screening weights cannot be negative")
	}
	return nil
}

func normalizePortfolioState(state *AppState) {
	if state == nil {
		return
	}
	if state.FX == nil {
		state.FX = map[string]float64{"CNY": 1}
	}
	if len(state.Stocks) == 0 {
		state.Stocks = stocksFromLegacy(state.Holdings, state.Candidates)
	}
	rebuildLegacyBuckets(state)
	if state.ScreeningWeights == (ScreeningWeights{}) {
		state.ScreeningWeights = DefaultScreeningWeights()
	}
}

func stocksFromLegacy(holdings []Holding, candidates []Candidate) []Stock {
	bySymbol := map[string]*Stock{}
	order := []string{}
	for _, holding := range holdings {
		stock := stockFromHolding(holding)
		key := normalizeSymbol(stock.Symbol)
		if key == "" {
			continue
		}
		if _, ok := bySymbol[key]; !ok {
			order = append(order, key)
			bySymbol[key] = &stock
		} else {
			mergeStock(bySymbol[key], stock, false)
		}
	}
	for _, candidate := range candidates {
		stock := stockFromCandidate(candidate)
		key := normalizeSymbol(stock.Symbol)
		if key == "" {
			continue
		}
		if existing, ok := bySymbol[key]; ok {
			mergeStock(existing, stock, true)
			continue
		}
		order = append(order, key)
		bySymbol[key] = &stock
	}
	result := make([]Stock, 0, len(order))
	for _, key := range order {
		result = append(result, *bySymbol[key])
	}
	return result
}

func findStock(stocks []Stock, symbol string) *Stock {
	normalized := normalizeSymbol(symbol)
	for i := range stocks {
		if normalizeSymbol(stocks[i].Symbol) == normalized {
			return &stocks[i]
		}
	}
	return nil
}

func stockFromHolding(holding Holding) Stock {
	return Stock{
		Symbol:              normalizeSymbol(holding.Symbol),
		Name:                holding.Name,
		Status:              holding.Status,
		Action:              holding.Action,
		CurrentPrice:        holding.CurrentPrice,
		PreviousClose:       holding.PreviousClose,
		TwentyDayClose:      holding.TwentyDayClose,
		TwentyDayCloseDate:  holding.TwentyDayCloseDate,
		TwentyDayChange:     holding.TwentyDayChange,
		MarketCap:           holding.MarketCap,
		MarketCapCurrency:   holding.MarketCapCurrency,
		CurrentPriceDate:    holding.CurrentPriceDate,
		PreviousCloseDate:   holding.PreviousCloseDate,
		MarginOfSafety:      holding.MarginOfSafety,
		QualityScore:        holding.QualityScore,
		Risk:                holding.Risk,
		Industry:            holding.Industry,
		Category:            holding.Category,
		Currency:            holding.Currency,
		IntrinsicValue:      holding.IntrinsicValue,
		FairValueRange:      holding.FairValueRange,
		TargetBuyPrice:      holding.TargetBuyPrice,
		PriceLevels:         holding.PriceLevels,
		ValuationConfidence: holding.ValuationConfidence,
		BusinessModel:       holding.BusinessModel,
		Moat:                holding.Moat,
		Governance:          holding.Governance,
		FinancialQuality:    holding.FinancialQuality,
		UpdatedAt:           holding.UpdatedAt,
		Notes:               holding.Notes,
		KillCriteria:        holding.KillCriteria,
		Reports:             holding.Reports,
		Dividend:            holding.Dividend,
		NetCash:             holding.NetCash,
		OwnerCashFlowAudit:  holding.OwnerCashFlowAudit,
		ResearchUpdates:     holding.ResearchUpdates,
		Financials:          holding.Financials,
		Position:            &StockPosition{Shares: holding.Shares, Cost: holding.Cost},
	}
}

func stockFromCandidate(candidate Candidate) Stock {
	return Stock{
		Symbol:              normalizeSymbol(candidate.Symbol),
		Name:                candidate.Name,
		Status:              candidate.Status,
		Action:              candidate.Action,
		CurrentPrice:        candidate.CurrentPrice,
		PreviousClose:       candidate.PreviousClose,
		TwentyDayClose:      candidate.TwentyDayClose,
		TwentyDayCloseDate:  candidate.TwentyDayCloseDate,
		TwentyDayChange:     candidate.TwentyDayChange,
		MarketCap:           candidate.MarketCap,
		MarketCapCurrency:   candidate.MarketCapCurrency,
		CurrentPriceDate:    candidate.CurrentPriceDate,
		PreviousCloseDate:   candidate.PreviousCloseDate,
		MarginOfSafety:      candidate.MarginOfSafety,
		QualityScore:        candidate.QualityScore,
		Risk:                candidate.Risk,
		Industry:            candidate.Industry,
		Category:            candidate.Category,
		Currency:            candidate.Currency,
		IntrinsicValue:      candidate.IntrinsicValue,
		FairValueRange:      candidate.FairValueRange,
		TargetBuyPrice:      candidate.TargetBuyPrice,
		PriceLevels:         candidate.PriceLevels,
		ValuationConfidence: candidate.ValuationConfidence,
		BusinessModel:       candidate.BusinessModel,
		Moat:                candidate.Moat,
		Governance:          candidate.Governance,
		FinancialQuality:    candidate.FinancialQuality,
		UpdatedAt:           candidate.UpdatedAt,
		Notes:               candidate.Notes,
		KillCriteria:        candidate.KillCriteria,
		Reports:             candidate.Reports,
		Dividend:            candidate.Dividend,
		NetCash:             candidate.NetCash,
		OwnerCashFlowAudit:  candidate.OwnerCashFlowAudit,
		ResearchUpdates:     candidate.ResearchUpdates,
		Financials:          candidate.Financials,
	}
}

func mergeStock(dst *Stock, src Stock, preserveExistingText bool) {
	if dst == nil {
		return
	}
	if strings.TrimSpace(src.Symbol) != "" && (!preserveExistingText || strings.TrimSpace(dst.Symbol) == "") {
		dst.Symbol = src.Symbol
	}
	if strings.TrimSpace(src.Name) != "" && (!preserveExistingText || strings.TrimSpace(dst.Name) == "") {
		dst.Name = src.Name
	}
	mergeStockString(&dst.Status, src.Status, preserveExistingText)
	mergeStockString(&dst.Action, src.Action, preserveExistingText)
	mergeStockString(&dst.Risk, src.Risk, preserveExistingText)
	mergeStockString(&dst.Industry, src.Industry, preserveExistingText)
	mergeStockString(&dst.Category, src.Category, preserveExistingText)
	mergeStockString(&dst.Currency, src.Currency, preserveExistingText)
	mergeStockString(&dst.FairValueRange, src.FairValueRange, preserveExistingText)
	mergeStockString(&dst.ValuationConfidence, src.ValuationConfidence, preserveExistingText)
	mergeStockString(&dst.BuyLogic, src.BuyLogic, preserveExistingText)
	mergeStockString(&dst.UpdatedAt, src.UpdatedAt, preserveExistingText)
	mergeStockString(&dst.Notes, src.Notes, preserveExistingText)
	mergeStockString(&dst.MarketCapCurrency, src.MarketCapCurrency, preserveExistingText)
	mergeStockString(&dst.CurrentPriceDate, src.CurrentPriceDate, preserveExistingText)
	mergeStockString(&dst.PreviousCloseDate, src.PreviousCloseDate, preserveExistingText)
	mergeStockString(&dst.TwentyDayCloseDate, src.TwentyDayCloseDate, preserveExistingText)
	mergeStockFloat(&dst.CurrentPrice, src.CurrentPrice, preserveExistingText)
	mergeStockFloat(&dst.PreviousClose, src.PreviousClose, preserveExistingText)
	mergeStockFloat(&dst.TwentyDayClose, src.TwentyDayClose, preserveExistingText)
	mergeStockPtr(&dst.TwentyDayChange, src.TwentyDayChange, preserveExistingText)
	mergeStockPtr(&dst.MarketCap, src.MarketCap, preserveExistingText)
	mergeStockPtr(&dst.MarginOfSafety, src.MarginOfSafety, preserveExistingText)
	mergeStockPtr(&dst.QualityScore, src.QualityScore, preserveExistingText)
	mergeStockPtr(&dst.IntrinsicValue, src.IntrinsicValue, preserveExistingText)
	mergeStockPtr(&dst.TargetBuyPrice, src.TargetBuyPrice, preserveExistingText)
	mergeStockPtr(&dst.BusinessModel, src.BusinessModel, preserveExistingText)
	mergeStockPtr(&dst.Moat, src.Moat, preserveExistingText)
	mergeStockPtr(&dst.Governance, src.Governance, preserveExistingText)
	mergeStockPtr(&dst.FinancialQuality, src.FinancialQuality, preserveExistingText)
	if src.PriceLevels != nil && (!preserveExistingText || dst.PriceLevels == nil) {
		dst.PriceLevels = src.PriceLevels
	}
	if len(src.KillCriteria) > 0 && (!preserveExistingText || len(dst.KillCriteria) == 0) {
		dst.KillCriteria = src.KillCriteria
	}
	if len(src.Reports) > 0 && (!preserveExistingText || len(dst.Reports) == 0) {
		dst.Reports = src.Reports
	}
	if src.Dividend != nil && (!preserveExistingText || dst.Dividend == nil) {
		dst.Dividend = src.Dividend
	}
	if src.NetCash != nil && (!preserveExistingText || dst.NetCash == nil) {
		dst.NetCash = src.NetCash
	}
	if src.OwnerCashFlowAudit != nil && (!preserveExistingText || dst.OwnerCashFlowAudit == nil) {
		dst.OwnerCashFlowAudit = src.OwnerCashFlowAudit
	}
	if len(src.ResearchUpdates) > 0 && (!preserveExistingText || len(dst.ResearchUpdates) == 0) {
		dst.ResearchUpdates = src.ResearchUpdates
	}
	if src.Financials != nil && (!preserveExistingText || dst.Financials == nil) {
		dst.Financials = src.Financials
	}
	if src.Position != nil && (!preserveExistingText || dst.Position == nil) {
		dst.Position = src.Position
	}
	if src.Valuation != nil && (!preserveExistingText || dst.Valuation == nil) {
		dst.Valuation = src.Valuation
	}
	if src.Screening != nil && (!preserveExistingText || dst.Screening == nil) {
		dst.Screening = src.Screening
	}
}

func mergeStockString(dst *string, value string, preserveExisting bool) {
	if strings.TrimSpace(value) == "" {
		return
	}
	if !preserveExisting || strings.TrimSpace(*dst) == "" {
		*dst = value
	}
}

func mergeStockFloat(dst *float64, value float64, preserveExisting bool) {
	if value <= 0 {
		return
	}
	if !preserveExisting || *dst <= 0 {
		*dst = value
	}
}

func mergeStockPtr(dst **float64, value *float64, preserveExisting bool) {
	if value == nil {
		return
	}
	if !preserveExisting || *dst == nil {
		*dst = value
	}
}

func rebuildLegacyBuckets(state *AppState) {
	holdings := []Holding{}
	candidates := []Candidate{}
	for _, stock := range state.Stocks {
		if stock.Position != nil {
			holdings = append(holdings, holdingFromStock(stock))
		}
		candidates = append(candidates, candidateFromStock(stock))
	}
	state.Holdings = holdings
	state.Candidates = candidates
}

func syncStocksFromLegacy(state *AppState) {
	if state == nil {
		return
	}
	existing := map[string]Stock{}
	for _, stock := range state.Stocks {
		existing[normalizeSymbol(stock.Symbol)] = stock
	}
	state.Stocks = stocksFromLegacy(state.Holdings, state.Candidates)
	for i := range state.Stocks {
		prior := existing[normalizeSymbol(state.Stocks[i].Symbol)]
		if state.Stocks[i].Valuation == nil {
			state.Stocks[i].Valuation = prior.Valuation
		}
		if state.Stocks[i].Screening == nil {
			state.Stocks[i].Screening = prior.Screening
		}
		state.Stocks[i].BuyLogic = prior.BuyLogic
	}
}

func holdingFromStock(stock Stock) Holding {
	holding := Holding{
		Symbol:              stock.Symbol,
		Name:                stock.Name,
		CurrentPrice:        stock.CurrentPrice,
		PreviousClose:       stock.PreviousClose,
		TwentyDayClose:      stock.TwentyDayClose,
		TwentyDayCloseDate:  stock.TwentyDayCloseDate,
		TwentyDayChange:     stock.TwentyDayChange,
		MarketCap:           stock.MarketCap,
		MarketCapCurrency:   stock.MarketCapCurrency,
		CurrentPriceDate:    stock.CurrentPriceDate,
		PreviousCloseDate:   stock.PreviousCloseDate,
		Action:              stock.Action,
		Status:              stock.Status,
		MarginOfSafety:      stock.MarginOfSafety,
		QualityScore:        stock.QualityScore,
		Risk:                stock.Risk,
		Industry:            stock.Industry,
		Category:            stock.Category,
		Currency:            stock.Currency,
		IntrinsicValue:      stock.IntrinsicValue,
		FairValueRange:      stock.FairValueRange,
		TargetBuyPrice:      stock.TargetBuyPrice,
		PriceLevels:         stock.PriceLevels,
		ValuationConfidence: stock.ValuationConfidence,
		BusinessModel:       stock.BusinessModel,
		Moat:                stock.Moat,
		Governance:          stock.Governance,
		FinancialQuality:    stock.FinancialQuality,
		UpdatedAt:           stock.UpdatedAt,
		Notes:               stock.Notes,
		KillCriteria:        stock.KillCriteria,
		Reports:             stock.Reports,
		Dividend:            stock.Dividend,
		NetCash:             stock.NetCash,
		OwnerCashFlowAudit:  stock.OwnerCashFlowAudit,
		ResearchUpdates:     stock.ResearchUpdates,
		Financials:          stock.Financials,
	}
	if stock.Position != nil {
		holding.Shares = stock.Position.Shares
		holding.Cost = stock.Position.Cost
	}
	return holding
}

func candidateFromStock(stock Stock) Candidate {
	return Candidate{
		Symbol:              stock.Symbol,
		Name:                stock.Name,
		Status:              stock.Status,
		Action:              stock.Action,
		CurrentPrice:        stock.CurrentPrice,
		PreviousClose:       stock.PreviousClose,
		TwentyDayClose:      stock.TwentyDayClose,
		TwentyDayCloseDate:  stock.TwentyDayCloseDate,
		TwentyDayChange:     stock.TwentyDayChange,
		MarketCap:           stock.MarketCap,
		MarketCapCurrency:   stock.MarketCapCurrency,
		CurrentPriceDate:    stock.CurrentPriceDate,
		PreviousCloseDate:   stock.PreviousCloseDate,
		MarginOfSafety:      stock.MarginOfSafety,
		QualityScore:        stock.QualityScore,
		Risk:                stock.Risk,
		Industry:            stock.Industry,
		Category:            stock.Category,
		Currency:            stock.Currency,
		IntrinsicValue:      stock.IntrinsicValue,
		FairValueRange:      stock.FairValueRange,
		TargetBuyPrice:      stock.TargetBuyPrice,
		PriceLevels:         stock.PriceLevels,
		ValuationConfidence: stock.ValuationConfidence,
		BusinessModel:       stock.BusinessModel,
		Moat:                stock.Moat,
		Governance:          stock.Governance,
		FinancialQuality:    stock.FinancialQuality,
		UpdatedAt:           stock.UpdatedAt,
		Notes:               stock.Notes,
		KillCriteria:        stock.KillCriteria,
		Reports:             stock.Reports,
		Dividend:            stock.Dividend,
		NetCash:             stock.NetCash,
		OwnerCashFlowAudit:  stock.OwnerCashFlowAudit,
		ResearchUpdates:     stock.ResearchUpdates,
		Financials:          stock.Financials,
	}
}

func (state AppState) MarshalJSON() ([]byte, error) {
	normalizePortfolioState(&state)
	type appStateJSON struct {
		TotalCapital     float64            `json:"totalCapital"`
		Cash             float64            `json:"cash"`
		FX               map[string]float64 `json:"fx"`
		Trades           []Trade            `json:"trades"`
		DecisionLogs     []DecisionLog      `json:"decisionLogs"`
		Stocks           []Stock            `json:"stocks"`
		ScreeningWeights ScreeningWeights   `json:"screeningWeights"`
		Plan             []PlanItem         `json:"plan"`
		Industries       []IndustryResearch `json:"industries,omitempty"`
		Rules            []Rule             `json:"rules"`
		DataStatus       *DataStatus        `json:"dataStatus,omitempty"`
	}
	return json.Marshal(appStateJSON{
		TotalCapital:     state.TotalCapital,
		Cash:             state.Cash,
		FX:               state.FX,
		Trades:           state.Trades,
		DecisionLogs:     state.DecisionLogs,
		Stocks:           state.Stocks,
		ScreeningWeights: state.ScreeningWeights,
		Plan:             state.Plan,
		Industries:       state.Industries,
		Rules:            state.Rules,
		DataStatus:       state.DataStatus,
	})
}

func CalculateValuationRange(assumptions ValuationAssumptions) (ValuationRange, error) {
	if len(assumptions.Scenarios) == 0 {
		return ValuationRange{}, errors.New("valuation scenarios are required")
	}
	values := make([]float64, 0, len(assumptions.Scenarios))
	for _, scenario := range assumptions.Scenarios {
		value := scenarioFairValue(scenario)
		if value <= 0 || math.IsNaN(value) || math.IsInf(value, 0) {
			return ValuationRange{}, fmt.Errorf("invalid valuation scenario %q", scenario.Name)
		}
		values = append(values, value)
	}
	sort.Float64s(values)
	result := ValuationRange{
		Low:      values[0],
		Base:     values[len(values)/2],
		High:     values[len(values)-1],
		Currency: assumptions.Currency,
	}
	if assumptions.CurrentPrice > 0 && result.Base > 0 {
		margin := (result.Base - assumptions.CurrentPrice) / result.Base
		result.MarginOfSafety = &margin
	}
	return result, nil
}

func scenarioFairValue(scenario ValuationScenario) float64 {
	if scenario.FairValue > 0 {
		return scenario.FairValue
	}
	if scenario.Shares <= 0 {
		return 0
	}
	values := []float64{}
	if scenario.FCF > 0 && scenario.ReasonablePFCF > 0 {
		values = append(values, scenario.FCF*scenario.ReasonablePFCF/scenario.Shares)
	}
	if scenario.FCF > 0 && scenario.ReasonablePE > 0 {
		values = append(values, scenario.FCF*scenario.ReasonablePE/scenario.Shares)
	}
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func ValuationPercentiles(points []ValuationHistoryPoint, currentPE float64, currentPB float64) (*float64, *float64) {
	return metricPercentile(points, currentPE, func(point ValuationHistoryPoint) *float64 { return point.PE }),
		metricPercentile(points, currentPB, func(point ValuationHistoryPoint) *float64 { return point.PB })
}

func metricPercentile(points []ValuationHistoryPoint, current float64, pick func(ValuationHistoryPoint) *float64) *float64 {
	if current <= 0 {
		return nil
	}
	count := 0
	lessOrEqual := 0
	for _, point := range points {
		value := pick(point)
		if value == nil || *value <= 0 {
			continue
		}
		count++
		if *value <= current {
			lessOrEqual++
		}
	}
	if count == 0 {
		return nil
	}
	percentile := float64(lessOrEqual) / float64(count)
	return &percentile
}
