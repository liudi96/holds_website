package main

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const runtimeValuationHistoryFileName = "valuation_history.json"

var runtimeValuationHistoryFile = filepath.Join(dataDir, "runtime", runtimeValuationHistoryFileName)

type ValuationHistoryBook struct {
	UpdatedAt   string                                `json:"updatedAt"`
	History     map[string][]ValuationHistoryPoint    `json:"history"`
	Percentiles map[string]ValuationHistoryPercentile `json:"percentiles,omitempty"`
	Skipped     []ValuationHistorySkip                `json:"skipped,omitempty"`
}

type ValuationHistoryUpdateResponse struct {
	Updated int                  `json:"updated"`
	Book    ValuationHistoryBook `json:"book"`
	State   AppState             `json:"state"`
}

type ValuationHistoryPercentile struct {
	PE *float64 `json:"pe,omitempty"`
	PB *float64 `json:"pb,omitempty"`
}

type ValuationHistorySkip struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name,omitempty"`
	Error  string `json:"error"`
}

var fetchValuationMonthlyCloses = fetchYahooMonthlyCloses

func (s *Server) handleUpdateScreeningWeights(w http.ResponseWriter, r *http.Request) {
	var weights ScreeningWeights
	if err := json.NewDecoder(r.Body).Decode(&weights); err != nil {
		writeError(w, http.StatusBadRequest, "invalid screening weights payload")
		return
	}
	if err := weights.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	nextState, err := cloneAppState(s.state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	nextState.ScreeningWeights = weights
	if err := saveAndHydrateState(&nextState); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.state = nextState
	writeJSON(w, http.StatusOK, s.state)
}

func (s *Server) handleUpsertStock(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid stock payload")
		return
	}
	var stock Stock
	if err := json.Unmarshal(body, &stock); err != nil {
		writeError(w, http.StatusBadRequest, "invalid stock payload")
		return
	}
	var raw map[string]json.RawMessage
	_ = json.Unmarshal(body, &raw)
	stock.Symbol = normalizeSymbol(stock.Symbol)
	stock.Name = strings.TrimSpace(stock.Name)
	if stock.Symbol == "" {
		writeError(w, http.StatusBadRequest, "symbol is required")
		return
	}
	if stock.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if strings.TrimSpace(stock.Currency) == "" {
		stock.Currency = inferCurrencyFromSymbol(stock.Symbol)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	nextState, err := cloneAppState(s.state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	upsertStock(&nextState, stock)
	if err := saveAndHydrateState(&nextState); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	applyExplicitStockPayload(&nextState, stock.Symbol, raw)
	if err := saveAndHydrateState(&nextState); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.state = nextState
	writeJSON(w, http.StatusOK, s.state)
}

func applyExplicitStockPayload(state *AppState, symbol string, raw map[string]json.RawMessage) {
	if state == nil || len(raw) == 0 {
		return
	}
	stock := findStock(state.Stocks, symbol)
	if stock == nil {
		return
	}
	applyString := func(key string, dst *string) {
		body, ok := raw[key]
		if !ok {
			return
		}
		if string(body) == "null" {
			*dst = ""
			return
		}
		var value string
		if err := json.Unmarshal(body, &value); err == nil {
			*dst = strings.TrimSpace(value)
		}
	}
	applyFloatPtr := func(key string, dst **float64) {
		body, ok := raw[key]
		if !ok {
			return
		}
		if string(body) == "null" {
			*dst = nil
			return
		}
		var value float64
		if err := json.Unmarshal(body, &value); err == nil {
			*dst = &value
		}
	}

	applyString("action", &stock.Action)
	applyString("status", &stock.Status)
	applyString("risk", &stock.Risk)
	applyString("notes", &stock.Notes)
	applyString("buyLogic", &stock.BuyLogic)
	applyString("valuationConfidence", &stock.ValuationConfidence)
	applyString("fairValueRange", &stock.FairValueRange)
	applyFloatPtr("qualityScore", &stock.QualityScore)
	applyFloatPtr("businessModel", &stock.BusinessModel)
	applyFloatPtr("moat", &stock.Moat)
	applyFloatPtr("governance", &stock.Governance)
	applyFloatPtr("financialQuality", &stock.FinancialQuality)
	applyFloatPtr("intrinsicValue", &stock.IntrinsicValue)
	applyFloatPtr("targetBuyPrice", &stock.TargetBuyPrice)
	applyFloatPtr("marginOfSafety", &stock.MarginOfSafety)
	if body, ok := raw["killCriteria"]; ok {
		if string(body) == "null" {
			stock.KillCriteria = nil
		} else {
			stock.KillCriteria = append(stock.KillCriteria[:0], body...)
		}
	}
	if body, ok := raw["valuation"]; ok && string(body) == "null" {
		stock.Valuation = nil
	}
}

func (s *Server) handleDeleteStock(w http.ResponseWriter, r *http.Request) {
	symbol := normalizeSymbol(strings.TrimPrefix(r.URL.Path, "/api/stocks/"))
	if symbol == "" {
		writeError(w, http.StatusBadRequest, "missing symbol")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	nextState, err := cloneAppState(s.state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	stock := findStock(nextState.Stocks, symbol)
	if stock == nil {
		writeError(w, http.StatusNotFound, "stock not found")
		return
	}
	if stockHasOpenPosition(*stock) {
		writeError(w, http.StatusConflict, "cannot delete a stock with an open position")
		return
	}
	before := len(nextState.Stocks)
	nextState.Stocks = removeStock(nextState.Stocks, symbol)
	if len(nextState.Stocks) == before {
		writeError(w, http.StatusNotFound, "stock not found")
		return
	}
	normalizePortfolioState(&nextState)
	if err := saveAndHydrateState(&nextState); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.state = nextState
	writeJSON(w, http.StatusOK, s.state)
}

func (s *Server) handleCreateDecisionLog(w http.ResponseWriter, r *http.Request) {
	var entry DecisionLog
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		writeError(w, http.StatusBadRequest, "invalid decision log payload")
		return
	}
	if strings.TrimSpace(entry.Decision) == "" || strings.TrimSpace(entry.Detail) == "" {
		writeError(w, http.StatusBadRequest, "decision and detail are required")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	nextState, err := cloneAppState(s.state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	entry.Detail = decisionLogDetailWithSnapshot(nextState, entry)
	appendDecisionLog(&nextState, entry)
	if err := saveAndHydrateState(&nextState); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.state = nextState
	writeJSON(w, http.StatusCreated, s.state)
}

func (s *Server) handleUpdateValuationHistory(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	state, err := cloneAppState(s.state)
	s.mu.Unlock()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	normalizePortfolioState(&state)

	now := time.Now()
	book := ValuationHistoryBook{
		UpdatedAt:   now.Format("2006-01-02 15:04:05"),
		History:     map[string][]ValuationHistoryPoint{},
		Percentiles: map[string]ValuationHistoryPercentile{},
	}
	client := &http.Client{Timeout: 15 * time.Second}
	for i := range state.Stocks {
		stock := &state.Stocks[i]
		symbol := normalizeSymbol(stock.Symbol)
		if symbol == "" {
			continue
		}
		points, historyErr := valuationHistoryPointsForStock(client, *stock, now)
		if historyErr != nil {
			book.Skipped = append(book.Skipped, ValuationHistorySkip{Symbol: symbol, Name: stock.Name, Error: historyErr.Error()})
		}
		if len(points) == 0 {
			continue
		}
		book.History[symbol] = points

		valuation := stock.Financials.Valuation
		currentPE := pointerValue(valuation.PE)
		currentPB := pointerValue(valuation.PB)
		pePercentile, pbPercentile := ValuationPercentiles(points, currentPE, currentPB)
		if pePercentile != nil || pbPercentile != nil {
			book.Percentiles[symbol] = ValuationHistoryPercentile{PE: pePercentile, PB: pbPercentile}
		}
	}
	if len(book.Percentiles) == 0 {
		book.Percentiles = nil
	}
	if err := saveValuationHistoryBook(book); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save valuation history")
		return
	}

	s.mu.Lock()
	nextState, err := cloneAppState(s.state)
	if err != nil {
		s.mu.Unlock()
		writeError(w, http.StatusInternalServerError, "failed to prepare state")
		return
	}
	updated := applyValuationHistoryPercentiles(&nextState, book.Percentiles)
	if updated > 0 {
		if err := saveAndHydrateState(&nextState); err != nil {
			s.mu.Unlock()
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		s.state = nextState
	}
	response := ValuationHistoryUpdateResponse{Updated: updated, Book: book, State: s.state}
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, response)
}

func valuationHistoryPointsForStock(client *http.Client, stock Stock, now time.Time) ([]ValuationHistoryPoint, error) {
	if stock.Financials == nil || stock.Financials.Valuation == nil {
		return nil, nil
	}
	eps, bps := valuationPerShareBases(stock)
	if eps <= 0 && bps <= 0 {
		return nil, nil
	}

	closes, err := fetchValuationMonthlyCloses(client, stock.Symbol)
	if err != nil || len(closes) == 0 {
		point := valuationHistoryPoint(now.Format("2006-01-02"), valuationBasePrice(stock), eps, bps)
		if point.Price != nil || point.PE != nil || point.PB != nil {
			return []ValuationHistoryPoint{point}, err
		}
		return nil, err
	}

	points := make([]ValuationHistoryPoint, 0, len(closes))
	for _, close := range closes {
		point := valuationHistoryPoint(close.Date, close.Price, eps, bps)
		if point.Price != nil || point.PE != nil || point.PB != nil {
			points = append(points, point)
		}
	}
	return points, err
}

func valuationPerShareBases(stock Stock) (float64, float64) {
	valuation := stock.Financials.Valuation
	price := valuationBasePrice(stock)
	var eps float64
	var bps float64
	if price > 0 {
		if valuation.PE != nil && *valuation.PE > 0 {
			eps = price / *valuation.PE
		}
		if valuation.PB != nil && *valuation.PB > 0 {
			bps = price / *valuation.PB
		}
	}
	for _, annual := range stock.Financials.Annual {
		if eps <= 0 && annual.EPS != nil && *annual.EPS > 0 {
			eps = *annual.EPS
		}
		if bps <= 0 && annual.BookValuePerShare != nil && *annual.BookValuePerShare > 0 {
			bps = *annual.BookValuePerShare
		}
		if eps > 0 && bps > 0 {
			break
		}
	}
	return eps, bps
}

func valuationBasePrice(stock Stock) float64 {
	if stock.Financials != nil && stock.Financials.Valuation != nil && stock.Financials.Valuation.Price != nil && *stock.Financials.Valuation.Price > 0 {
		return *stock.Financials.Valuation.Price
	}
	return stock.CurrentPrice
}

func valuationHistoryPoint(date string, price float64, eps float64, bps float64) ValuationHistoryPoint {
	point := ValuationHistoryPoint{Date: date, Price: pricePointer(price)}
	if price > 0 && eps > 0 {
		point.PE = ptr(price / eps)
	}
	if price > 0 && bps > 0 {
		point.PB = ptr(price / bps)
	}
	return point
}

func pointerValue(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func stockHasOpenPosition(stock Stock) bool {
	return stock.Position != nil && stock.Position.Shares > 0
}

func applyValuationHistoryPercentiles(state *AppState, percentiles map[string]ValuationHistoryPercentile) int {
	if state == nil || len(percentiles) == 0 {
		return 0
	}
	normalizePortfolioState(state)
	updated := 0
	for i := range state.Stocks {
		percentile, ok := percentiles[normalizeSymbol(state.Stocks[i].Symbol)]
		if !ok {
			continue
		}
		if state.Stocks[i].Financials == nil {
			state.Stocks[i].Financials = &Financials{}
		}
		if state.Stocks[i].Financials.Valuation == nil {
			state.Stocks[i].Financials.Valuation = &FinancialValuation{}
		}
		state.Stocks[i].Financials.Valuation.PEPercentile = percentile.PE
		state.Stocks[i].Financials.Valuation.PBPercentile = percentile.PB
		updated++
	}
	normalizePortfolioState(state)
	return updated
}

func saveAndHydrateState(state *AppState) error {
	if _, err := saveStateWithBackup(*state); err != nil {
		return err
	}
	if err := hydrateState(state); err != nil {
		return err
	}
	return nil
}

func upsertStock(state *AppState, stock Stock) {
	normalizePortfolioState(state)
	for i := range state.Stocks {
		if normalizeSymbol(state.Stocks[i].Symbol) == normalizeSymbol(stock.Symbol) {
			mergeStock(&state.Stocks[i], stock, false)
			normalizePortfolioState(state)
			return
		}
	}
	state.Stocks = append(state.Stocks, stock)
	normalizePortfolioState(state)
}

func removeStock(stocks []Stock, symbol string) []Stock {
	normalized := normalizeSymbol(symbol)
	result := stocks[:0]
	for _, stock := range stocks {
		if normalizeSymbol(stock.Symbol) != normalized {
			result = append(result, stock)
		}
	}
	return result
}

func decisionLogDetailWithSnapshot(state AppState, entry DecisionLog) string {
	stock := findStock(state.Stocks, entry.Symbol)
	detail := strings.TrimSpace(entry.Detail)
	if stock == nil {
		return detail
	}
	parts := []string{detail}
	if stock.CurrentPrice > 0 {
		parts = append(parts, "价格 "+strings.ToUpper(stock.Currency)+" "+formatFixed(stock.CurrentPrice, 4))
	}
	if stock.Position != nil {
		parts = append(parts, "持仓 "+formatFixed(stock.Position.Shares, 2)+" 股；成本 "+strings.ToUpper(stock.Currency)+" "+formatFixed(stock.Position.Cost, 4))
	}
	if strings.TrimSpace(stock.FairValueRange) != "" {
		parts = append(parts, "估值区间 "+strings.TrimSpace(stock.FairValueRange))
	}
	if stock.MarginOfSafety != nil {
		parts = append(parts, "安全边际 "+formatFixed(*stock.MarginOfSafety*100, 2)+"%")
	}
	return strings.Join(parts, "；")
}

func formatFixed(value float64, digits int) string {
	return strconv.FormatFloat(value, 'f', digits, 64)
}

func saveValuationHistoryBook(book ValuationHistoryBook) error {
	body, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return writeFileAtomic(runtimeValuationHistoryFile, body, 0o644)
}
