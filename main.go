package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const dataFile = "data/portfolio.json"
const decisionLogLimit = 500

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
	Symbol              string          `json:"symbol"`
	Name                string          `json:"name"`
	Shares              float64         `json:"shares"`
	Cost                float64         `json:"cost"`
	CurrentPrice        float64         `json:"currentPrice"`
	PreviousClose       float64         `json:"previousClose"`
	MarketCap           *float64        `json:"marketCap,omitempty"`
	MarketCapCurrency   string          `json:"marketCapCurrency,omitempty"`
	CurrentPriceDate    string          `json:"currentPriceDate"`
	PreviousCloseDate   string          `json:"previousCloseDate"`
	Action              string          `json:"action"`
	Status              string          `json:"status"`
	MarginOfSafety      *float64        `json:"marginOfSafety"`
	QualityScore        *float64        `json:"qualityScore"`
	Risk                string          `json:"risk"`
	Industry            string          `json:"industry"`
	Currency            string          `json:"currency"`
	IntrinsicValue      *float64        `json:"intrinsicValue"`
	FairValueRange      string          `json:"fairValueRange"`
	TargetBuyPrice      *float64        `json:"targetBuyPrice"`
	PriceLevels         *PriceLevels    `json:"priceLevels,omitempty"`
	ValuationConfidence string          `json:"valuationConfidence,omitempty"`
	BusinessModel       *float64        `json:"businessModel"`
	Moat                *float64        `json:"moat"`
	Governance          *float64        `json:"governance"`
	FinancialQuality    *float64        `json:"financialQuality"`
	UpdatedAt           string          `json:"updatedAt"`
	Notes               string          `json:"notes"`
	KillCriteria        json.RawMessage `json:"killCriteria,omitempty"`
	Reports             []Report        `json:"reports,omitempty"`
	Dividend            *Dividend       `json:"dividend,omitempty"`
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
	Symbol              string          `json:"symbol"`
	Name                string          `json:"name"`
	Status              string          `json:"status"`
	Action              string          `json:"action"`
	CurrentPrice        float64         `json:"currentPrice"`
	PreviousClose       float64         `json:"previousClose"`
	MarketCap           *float64        `json:"marketCap,omitempty"`
	MarketCapCurrency   string          `json:"marketCapCurrency,omitempty"`
	CurrentPriceDate    string          `json:"currentPriceDate"`
	PreviousCloseDate   string          `json:"previousCloseDate"`
	MarginOfSafety      *float64        `json:"marginOfSafety"`
	QualityScore        *float64        `json:"qualityScore"`
	Risk                string          `json:"risk"`
	Industry            string          `json:"industry"`
	Currency            string          `json:"currency"`
	IntrinsicValue      *float64        `json:"intrinsicValue"`
	FairValueRange      string          `json:"fairValueRange"`
	TargetBuyPrice      *float64        `json:"targetBuyPrice"`
	PriceLevels         *PriceLevels    `json:"priceLevels,omitempty"`
	ValuationConfidence string          `json:"valuationConfidence,omitempty"`
	BusinessModel       *float64        `json:"businessModel"`
	Moat                *float64        `json:"moat"`
	Governance          *float64        `json:"governance"`
	FinancialQuality    *float64        `json:"financialQuality"`
	UpdatedAt           string          `json:"updatedAt"`
	Notes               string          `json:"notes"`
	KillCriteria        json.RawMessage `json:"killCriteria,omitempty"`
	Reports             []Report        `json:"reports,omitempty"`
	Dividend            *Dividend       `json:"dividend,omitempty"`
}

type Rule struct {
	Dimension string  `json:"dimension"`
	Score     float64 `json:"score"`
	Standard  string  `json:"standard"`
}

type Server struct {
	mu    sync.Mutex
	state AppState
}

func main() {
	state, err := loadState()
	if err != nil {
		log.Fatal(err)
	}

	server := &Server{state: state}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/state", server.handleGetState)
	mux.HandleFunc("POST /api/reset", server.handleReset)
	mux.HandleFunc("POST /api/trades", server.handleCreateTrade)
	mux.HandleFunc("PUT /api/holdings/", server.handleUpdateHolding)
	mux.HandleFunc("POST /api/research/preview", server.handlePreviewResearch)
	mux.HandleFunc("POST /api/research/import", server.handleImportResearch)
	mux.HandleFunc("GET /api/chatgpt/export", server.handleExportChatGPTContext)
	mux.HandleFunc("POST /api/quotes/update", server.handleUpdateQuotes)
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("/", http.FileServer(http.Dir(".")))

	addr := "0.0.0.0:8080"
	log.Printf("portfolio desk listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := loadState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}
	s.state = state

	writeJSON(w, http.StatusOK, s.state)
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := loadState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}
	s.state = state
	writeJSON(w, http.StatusOK, s.state)
}

func preserveQuoteFields(nextHoldings, currentHoldings []Holding) {
	currentBySymbol := make(map[string]Holding, len(currentHoldings))
	for _, holding := range currentHoldings {
		currentBySymbol[strings.ToUpper(strings.TrimSpace(holding.Symbol))] = holding
	}

	for i := range nextHoldings {
		current, ok := currentBySymbol[strings.ToUpper(strings.TrimSpace(nextHoldings[i].Symbol))]
		if !ok {
			continue
		}
		if current.CurrentPrice > 0 {
			nextHoldings[i].CurrentPrice = current.CurrentPrice
		}
		if current.PreviousClose > 0 {
			nextHoldings[i].PreviousClose = current.PreviousClose
		}
		nextHoldings[i].CurrentPriceDate = current.CurrentPriceDate
		nextHoldings[i].PreviousCloseDate = current.PreviousCloseDate
		if strings.Contains(current.UpdatedAt, "行情源") {
			nextHoldings[i].UpdatedAt = current.UpdatedAt
		}
	}
}

func (s *Server) handleCreateTrade(w http.ResponseWriter, r *http.Request) {
	var trade Trade
	if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
		writeError(w, http.StatusBadRequest, "invalid trade payload")
		return
	}

	if err := validateTrade(trade); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	trade.ID = time.Now().UnixMilli()
	trade.Date = time.Now().Format("2006-01-02")
	trade.Side = strings.ToLower(strings.TrimSpace(trade.Side))
	trade.Symbol = strings.ToUpper(strings.TrimSpace(trade.Symbol))
	trade.Currency = strings.ToUpper(strings.TrimSpace(trade.Currency))

	s.applyTrade(trade)
	appendTradeDecisionLog(&s.state, trade)
	if err := saveState(s.state); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save state")
		return
	}

	writeJSON(w, http.StatusCreated, s.state)
}

func (s *Server) handleUpdateHolding(w http.ResponseWriter, r *http.Request) {
	symbol := strings.TrimPrefix(r.URL.Path, "/api/holdings/")
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		writeError(w, http.StatusBadRequest, "missing symbol")
		return
	}

	var patch Holding
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid holding payload")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.state.Holdings {
		if strings.EqualFold(s.state.Holdings[i].Symbol, symbol) {
			s.state.Holdings[i].Name = strings.TrimSpace(patch.Name)
			s.state.Holdings[i].Industry = strings.TrimSpace(patch.Industry)
			s.state.Holdings[i].Action = strings.TrimSpace(patch.Action)
			s.state.Holdings[i].Status = strings.TrimSpace(patch.Status)
			s.state.Holdings[i].MarginOfSafety = marginOfSafetyFromPrice(s.state.Holdings[i].IntrinsicValue, s.state.Holdings[i].CurrentPrice, patch.MarginOfSafety)
			s.state.Holdings[i].QualityScore = patch.QualityScore
			s.state.Holdings[i].Notes = strings.TrimSpace(patch.Notes)
			if err := saveState(s.state); err != nil {
				writeError(w, http.StatusInternalServerError, "failed to save state")
				return
			}
			writeJSON(w, http.StatusOK, s.state)
			return
		}
	}

	writeError(w, http.StatusNotFound, "holding not found")
}

func (s *Server) applyTrade(trade Trade) {
	idx := -1
	for i := range s.state.Holdings {
		if strings.EqualFold(s.state.Holdings[i].Symbol, trade.Symbol) {
			idx = i
			break
		}
	}

	if idx == -1 {
		s.state.Holdings = append(s.state.Holdings, Holding{
			Symbol:       trade.Symbol,
			Name:         trade.Name,
			Shares:       0,
			Cost:         trade.Price,
			CurrentPrice: trade.CurrentPrice,
			Currency:     trade.Currency,
		})
		idx = len(s.state.Holdings) - 1
	}

	holding := &s.state.Holdings[idx]
	if trade.Side == "buy" {
		totalCost := holding.Shares*holding.Cost + trade.Shares*trade.Price
		holding.Shares += trade.Shares
		if holding.Shares > 0 {
			holding.Cost = totalCost / holding.Shares
		}
	} else {
		holding.Shares -= trade.Shares
		if holding.Shares < 0 {
			holding.Shares = 0
		}
	}

	if strings.TrimSpace(trade.Name) != "" {
		holding.Name = strings.TrimSpace(trade.Name)
	}
	holding.Currency = trade.Currency
	holding.CurrentPrice = trade.CurrentPrice
	s.state.Trades = append(s.state.Trades, trade)

	multiplier := s.state.FX[trade.Currency]
	if multiplier == 0 {
		multiplier = 1
	}
	cashDelta := trade.Shares * trade.Price * multiplier
	if trade.Side == "buy" {
		s.state.Cash -= cashDelta
	} else {
		s.state.Cash += cashDelta
	}
}

func validateTrade(trade Trade) error {
	side := strings.ToLower(strings.TrimSpace(trade.Side))
	if side != "buy" && side != "sell" {
		return errors.New("side must be buy or sell")
	}
	if strings.TrimSpace(trade.Symbol) == "" {
		return errors.New("symbol is required")
	}
	if trade.Shares <= 0 {
		return errors.New("shares must be positive")
	}
	if trade.Price <= 0 {
		return errors.New("price must be positive")
	}
	if trade.CurrentPrice <= 0 {
		return errors.New("currentPrice must be positive")
	}
	return nil
}

func appendDecisionLog(state *AppState, entry DecisionLog) {
	entry.Type = strings.TrimSpace(entry.Type)
	if entry.Type == "" {
		entry.Type = "event"
	}
	entry.Symbol = strings.ToUpper(strings.TrimSpace(entry.Symbol))
	entry.Name = strings.TrimSpace(entry.Name)
	entry.Currency = strings.ToUpper(strings.TrimSpace(entry.Currency))
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

func appendResearchDecisionLog(state *AppState, research ResearchImport, summary string, targetType string) {
	name, price, currency, decision, discipline := decisionLogContext(state, research.Symbol)
	name = firstNonEmpty(name, research.Name)
	currency = firstNonEmpty(currency, research.Currency)
	decision = firstNonEmpty(decision, research.Action, research.Status)
	discipline = firstNonEmpty(discipline, research.Plan.Discipline, research.Status)
	detail := strings.TrimSpace(summary)
	if strings.TrimSpace(targetType) != "" {
		detail = strings.TrimSpace(targetType + "；" + detail)
	}

	appendDecisionLog(state, DecisionLog{
		Type:       "research",
		Symbol:     research.Symbol,
		Name:       name,
		Price:      price,
		Currency:   currency,
		Decision:   decision,
		Discipline: discipline,
		Detail:     detail,
	})
}

func appendTradeDecisionLog(state *AppState, trade Trade) {
	name, _, currency, decision, discipline := decisionLogContext(state, trade.Symbol)
	name = firstNonEmpty(name, trade.Name)
	currency = firstNonEmpty(currency, trade.Currency)
	sideText := "买入"
	if trade.Side == "sell" {
		sideText = "卖出"
	}
	decision = firstNonEmpty(decision, fmt.Sprintf("%s %s", sideText, firstNonEmpty(name, trade.Symbol)))
	detail := fmt.Sprintf("%s %.2f 股；成交价 %s %.4f；录入最新价 %s %.4f", sideText, trade.Shares, strings.ToUpper(trade.Currency), trade.Price, strings.ToUpper(trade.Currency), trade.CurrentPrice)

	appendDecisionLog(state, DecisionLog{
		Type:       "trade",
		Symbol:     trade.Symbol,
		Name:       name,
		Price:      ptr(trade.Price),
		Currency:   currency,
		Decision:   decision,
		Discipline: discipline,
		Detail:     detail,
	})
}

func appendQuoteDecisionLogs(state *AppState, now time.Time) {
	updateLabel := now.Format("2006-01-02 15:04:05")
	for i := range state.Holdings {
		holding := state.Holdings[i]
		if !strings.HasPrefix(holding.UpdatedAt, updateLabel) {
			continue
		}
		_, _, _, decision, discipline := decisionLogContext(state, holding.Symbol)
		appendDecisionLog(state, DecisionLog{
			Date:       updateLabel,
			Type:       "quote",
			Symbol:     holding.Symbol,
			Name:       holding.Name,
			Price:      pricePointer(holding.CurrentPrice),
			Currency:   holding.Currency,
			Decision:   firstNonEmpty(decision, holding.Action, holding.Status),
			Discipline: firstNonEmpty(discipline, holding.Status),
			Detail:     quoteDecisionDetail(holding.CurrentPriceDate, holding.PreviousCloseDate),
		})
	}

	for i := range state.Candidates {
		candidate := state.Candidates[i]
		if !strings.HasPrefix(candidate.UpdatedAt, updateLabel) {
			continue
		}
		_, _, _, decision, discipline := decisionLogContext(state, candidate.Symbol)
		appendDecisionLog(state, DecisionLog{
			Date:       updateLabel,
			Type:       "quote",
			Symbol:     candidate.Symbol,
			Name:       candidate.Name,
			Price:      pricePointer(candidate.CurrentPrice),
			Currency:   candidate.Currency,
			Decision:   firstNonEmpty(decision, candidate.Action, candidate.Status),
			Discipline: firstNonEmpty(discipline, candidate.Status),
			Detail:     quoteDecisionDetail(candidate.CurrentPriceDate, candidate.PreviousCloseDate),
		})
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
		discipline := firstNonEmpty(planDiscipline(plan), holding.Status)
		return holding.Name, pricePointer(holding.CurrentPrice), holding.Currency, firstNonEmpty(holding.Action, holding.Status), discipline
	}

	for i := range state.Candidates {
		candidate := state.Candidates[i]
		if normalizeSymbol(candidate.Symbol) != normalizedSymbol {
			continue
		}
		plan := findPlanForDecisionLog(state, candidate.Symbol, candidate.Name)
		discipline := firstNonEmpty(planDiscipline(plan), candidate.Status)
		return candidate.Name, pricePointer(candidate.CurrentPrice), candidate.Currency, firstNonEmpty(candidate.Action, candidate.Status), discipline
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
	return ptr(value)
}

func quoteDecisionDetail(currentDate string, previousDate string) string {
	return fmt.Sprintf("今收 %s；昨收 %s", firstNonEmpty(currentDate, "未知"), firstNonEmpty(previousDate, "未知"))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if text := strings.TrimSpace(value); text != "" {
			return text
		}
	}
	return ""
}

func loadState() (AppState, error) {
	if _, err := os.Stat(dataFile); errors.Is(err, os.ErrNotExist) {
		state := defaultState()
		return state, saveState(state)
	}

	body, err := os.ReadFile(dataFile)
	if err != nil {
		return AppState{}, err
	}

	var state AppState
	if err := json.Unmarshal(body, &state); err != nil {
		return AppState{}, err
	}
	return state, nil
}

func saveState(state AppState) error {
	if err := os.MkdirAll(filepath.Dir(dataFile), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, body, 0o644)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func ptr(value float64) *float64 {
	return &value
}

func defaultState() AppState {
	return AppState{
		TotalCapital: 1150000,
		Cash:         477238.13560642325,
		FX:           map[string]float64{"CNY": 1, "HKD": 0.8716, "USD": 7.1},
		Trades:       []Trade{},
		DecisionLogs: []DecisionLog{},
		Holdings: []Holding{
			{Symbol: "0700.HK", Name: "腾讯控股", Shares: 200, Cost: 480.43, CurrentPrice: 463, PreviousClose: 463, Action: "继续持有；新资金暂不追买，放入核心替补", Status: "未达标（安全边际<15%）", MarginOfSafety: ptr(0.09), QualityScore: ptr(89), Risk: "无立即否决；政策/地缘/AI投入需折价", Industry: "互联网平台/游戏/广告/金融科技", Currency: "HKD", IntrinsicValue: ptr(508), FairValueRange: "HK$480-560", TargetBuyPrice: ptr(432), BusinessModel: ptr(28), Moat: ptr(23), Governance: ptr(17), FinancialQuality: ptr(21), UpdatedAt: "2026-05-06；最新价约HK$463.00；HKD/CNY约0.8716；FY2025", Notes: "FY2025：收入RMB7518亿、Non-IFRS净利RMB2596亿、FCF RMB1826亿、净现金RMB1071亿。"},
			{Symbol: "000333.SZ", Name: "美的集团", Shares: 600, Cost: 79.638, CurrentPrice: 80.44, PreviousClose: 80.44, Action: "放入核心替补；A股暂不追买，H股优先但等待≤HK$86-87", Status: "未达标（A股安全边际<20%；H股接近达标）", MarginOfSafety: ptr(0.153), QualityScore: ptr(88), Risk: "无立即否决；Q1扣非下滑、海外关税/汇率、价格战需跟踪", Industry: "家电/全球化制造/ToB楼宇科技/机器人自动化", Currency: "CNY", IntrinsicValue: ptr(95), FairValueRange: "¥90-100", TargetBuyPrice: ptr(76), BusinessModel: ptr(28), Moat: ptr(23), Governance: ptr(18), FinancialQuality: ptr(19), UpdatedAt: "2026-05-06；A股最新价约¥80.44；H股约HK$87.70；FY2025/2026Q1", Notes: "FY2025：营收RMB4585亿、归母净利RMB439.45亿、年度分红¥4.30/股。"},
			{Symbol: "002415.SZ", Name: "海康威视", Shares: 1200, Cost: 34.54, CurrentPrice: 36.29, PreviousClose: 36.29, Action: "重点预期差候选/核心替补边缘；可小仓验证，不宜重仓；Q2验证后再升级", Status: "未达标（安全边际约13.6%<25%；预期差仓可观察）", MarginOfSafety: ptr(0.136), QualityScore: ptr(84), Risk: "无一票否决；地缘/合规/实体清单、Q1经营现金流为负、AIoT重估需验证", Industry: "AIoT/安防/机器视觉/科技制造平台", Currency: "CNY", IntrinsicValue: ptr(42), FairValueRange: "¥34-48", TargetBuyPrice: ptr(31.5), BusinessModel: ptr(25), Moat: ptr(23), Governance: ptr(16), FinancialQuality: ptr(20), UpdatedAt: "2026-05-06；最新价约¥36.29；FY2025/2026Q1；董秘大额增持后修正", Notes: "FY2025：营收约RMB925.08亿、归母净利约RMB141.95亿；2026Q1归母净利同比+36.42%。"},
			{Symbol: "600887.SH", Name: "伊利股份", Shares: 1300, Cost: 26.469, CurrentPrice: 27.45, PreviousClose: 27.45, Action: "放入核心替补；暂不追买，等待¥24-26", Status: "未达标（安全边际约14.2%<25%）", MarginOfSafety: ptr(0.1421875), QualityScore: ptr(83), Risk: "无一票否决；需求弱复苏、原奶上涨传导不顺、液奶仍下滑、食品安全风险需跟踪", Industry: "乳制品/消费龙头/高股息/奶周期修复", Currency: "CNY", IntrinsicValue: ptr(32), FairValueRange: "¥28-36", TargetBuyPrice: ptr(24), BusinessModel: ptr(24), Moat: ptr(22), Governance: ptr(16), FinancialQuality: ptr(21), UpdatedAt: "2026-05-07；最新价约¥27.45；FY2025/2026Q1；奶周期底部右侧观察", Notes: "2025拟派息¥1.38/股，按¥27.45股息率约5.0%；达标买入价≤¥24。"},
			{Symbol: "600036.SH", Name: "招商银行", Shares: 500, Cost: 39.18, CurrentPrice: 39.18, PreviousClose: 39.18, Currency: "CNY"},
			{Symbol: "0696.HK", Name: "民航信", Shares: 11000, Cost: 10.648, CurrentPrice: 10.648, PreviousClose: 10.648, Currency: "HKD"},
			{Symbol: "0506.HK", Name: "中国食品", Shares: 22000, Cost: 4.041, CurrentPrice: 4.041, PreviousClose: 4.041, Currency: "HKD"},
			{Symbol: "2669.HK", Name: "中海物业", Shares: 20000, Cost: 4.468, CurrentPrice: 4.468, PreviousClose: 4.468, Currency: "HKD"},
			{Symbol: "6049.HK", Name: "保利物业", Shares: 2600, Cost: 32.663, CurrentPrice: 32.663, PreviousClose: 32.663, Currency: "HKD"},
			{Symbol: "0883.HK", Name: "中海油", Shares: 2000, Cost: 29.326, CurrentPrice: 29.326, PreviousClose: 29.326, Currency: "HKD"},
			{Symbol: "1448.HK", Name: "福寿园", Shares: 11000, Cost: 2.521, CurrentPrice: 2.64, PreviousClose: 2.64, CurrentPriceDate: "2026-05-07", PreviousCloseDate: "2026-05-07", Action: "暂不行动；不买入；不纳入核心替补，等待2025年报、审计意见、法证调查结论和复牌后再重估", Status: "未达标（停牌、年报延迟、治理与财务可靠性风险未解除）", MarginOfSafety: ptr(0), QualityScore: ptr(62), Risk: "已触发重大风险否决项：停牌、业绩延迟、现金及采购付款事项调查、管理层/内控可信度下降、墓穴ASP大幅下滑、资产和商誉减值风险", Industry: "殡葬服务/墓园运营/生命服务", Currency: "HKD", IntrinsicValue: ptr(2.65), FairValueRange: "HK$1.6-3.1", TargetBuyPrice: ptr(2), BusinessModel: ptr(22), Moat: ptr(16), Governance: ptr(5), FinancialQuality: ptr(19), UpdatedAt: "2026-05-07；停牌前最后价约HK$2.64；用户更新分析", Notes: "计划：剔除/仅风险观察。复牌前不行动；复牌后若审计无保留、调查无重大重述且价格≤HK$2.0-2.2，才重新评估普通候选价值。纪律：质量分低于75且有重大风险否决项；不因低估值或净现金买入，先等风险解除。最新市场状态：股份自2026-03-20起停牌，停牌前最后价约HK$2.64。最新可用财务口径：2024收入约RMB20.77亿，归母净利约RMB3.73亿，EPS约RMB0.164；2025H1收入约RMB6.11亿，归母亏损约RMB2.61亿，EPS约-RMB0.115。核心判断：福寿园当前不是单纯估值杀，而是业绩杀、治理杀和财报可信度风险叠加；内在价值区间仅为压力测试，不作为可执行买入依据。"},
			{Symbol: "07489.HK", Name: "岚图汽车", Shares: 2132, Cost: 0, CurrentPrice: 5.89, PreviousClose: 5.89, CurrentPriceDate: "2026-05-07", PreviousCloseDate: "2026-05-07", Action: "放入普通候选池观察；当前不买入，等待扣非利润和自由现金流验证", Status: "未达标（质量分<75且安全边际不足）", MarginOfSafety: ptr(0.16), QualityScore: ptr(72), Risk: "盈利质量受政府补助影响，梦想家单一车型依赖较高，新能源车价格战和智能化竞争可能压缩毛利率", Industry: "新能源乘用车/高端MPV/央企汽车", Currency: "HKD", IntrinsicValue: ptr(7), FairValueRange: "HK$4.5-8.5", TargetBuyPrice: ptr(4.8), BusinessModel: ptr(21), Moat: ptr(16), Governance: ptr(16), FinancialQuality: ptr(19), UpdatedAt: "2026-05-07；估值基于HK$5.89附近股价；用户更新分析", Notes: "2025年收入约人民币348.65亿元，毛利率约20.9%，净利润约人民币10.17亿元，首次年度盈利；2025年交付约150169辆，2026年1-4月交付约49038辆。估值基于HK$5.89附近股价、市值约HK$216.8亿、PE约16.9倍、PB约1.78倍。核心假设是2026年需验证扣非利润、经营现金流和自由现金流质量。"},
		},
		Plan: []PlanItem{
			{Rank: 1, Name: "腾讯控股", Priority: "观察/低优先级", Advice: "继续持有；新资金等待≤HK$432，HK$400-430可分批", Discipline: "优秀资产要求≥15%安全边际；当前约9%，未达标"},
			{Rank: 2, Name: "美的集团", Priority: "核心替补/中优先级", Advice: "A股等待≤¥76分批；H股≤HK$86-87优先；当前不追买", Discipline: "优秀资产要求≥20%安全边际；A股当前约15.3%，未达标"},
			{Rank: 3, Name: "海康威视", Priority: "重点预期差候选/中优先级", Advice: "不重仓；¥35-37仅适合小仓验证，¥30-32更从容；Q2验证后可升核心替补", Discipline: "质量分84，合格候选要求≥25%安全边际"},
			{Rank: 4, Name: "伊利股份", Priority: "核心替补/中低优先级", Advice: "暂不追买；¥25-26开始关注，≤¥24可考虑分批", Discipline: "质量分83，合格候选要求≥25%安全边际"},
			{Rank: 99, Name: "岚图汽车", Priority: "普通候选池/低优先级", Advice: "HK$4.2-4.8才接近可观察买入区；若2026H1扣非利润和自由现金流转正，可重新上修估值", Discipline: "质量分低于75原则上不进入核心资产池；安全边际不足时不试仓"},
		},
		Candidates: []Candidate{
			{Symbol: "600690.SH", Name: "海尔智家A", Status: "候选池", Action: "放入普通候选池观察；A股暂不追，H股赔率更优", MarginOfSafety: ptr(0.17), QualityScore: ptr(83), Industry: "家电/全球化白电/智慧家庭", Currency: "CNY", IntrinsicValue: ptr(26), FairValueRange: "¥24-28", TargetBuyPrice: ptr(19.5)},
		},
		Rules: []Rule{
			{Dimension: "商业模式", Score: 30, Standard: "需求刚性、收入可重复、定价权、资本开支、行业空间"},
			{Dimension: "护城河", Score: 25, Standard: "品牌/规模/网络效应/牌照/成本优势、份额稳定、利润率优于同行"},
			{Dimension: "管理层/企业文化/治理", Score: 20, Standard: "长期主义、资本配置、股东回报、披露透明、少画饼"},
			{Dimension: "财务质量", Score: 25, Standard: "ROE/ROIC、自由现金流、资产负债表、利润率、应收/存货/资本开支"},
		},
	}
}
