package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const dataFile = "data/portfolio.json"

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
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Handle("/", http.FileServer(http.Dir(".")))

	addr := ":8080"
	log.Printf("portfolio desk listening on http://127.0.0.1%s", addr)
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

	if state, err := loadState(); err == nil {
		s.state = state
	}

	nextState := defaultState()
	preserveQuoteFields(nextState.Holdings, s.state.Holdings)
	s.state = nextState
	if err := saveState(s.state); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save state")
		return
	}
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
			s.state.Holdings[i].MarginOfSafety = patch.MarginOfSafety
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
		TotalCapital: 1100000,
		Cash:         487794.75,
		FX:           map[string]float64{"CNY": 1, "HKD": 0.8716, "USD": 7.1},
		Trades:       []Trade{},
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
		},
		Plan: []PlanItem{
			{Rank: 1, Name: "腾讯控股", Priority: "观察/低优先级", Advice: "继续持有；新资金等待≤HK$432，HK$400-430可分批", Discipline: "优秀资产要求≥15%安全边际；当前约9%，未达标"},
			{Rank: 2, Name: "美的集团", Priority: "核心替补/中优先级", Advice: "A股等待≤¥76分批；H股≤HK$86-87优先；当前不追买", Discipline: "优秀资产要求≥20%安全边际；A股当前约15.3%，未达标"},
			{Rank: 3, Name: "海康威视", Priority: "重点预期差候选/中优先级", Advice: "不重仓；¥35-37仅适合小仓验证，¥30-32更从容；Q2验证后可升核心替补", Discipline: "质量分84，合格候选要求≥25%安全边际"},
			{Rank: 4, Name: "伊利股份", Priority: "核心替补/中低优先级", Advice: "暂不追买；¥25-26开始关注，≤¥24可考虑分批", Discipline: "质量分83，合格候选要求≥25%安全边际"},
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
