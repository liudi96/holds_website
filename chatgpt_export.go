package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	chatGPTExportTimezone              = "Asia/Shanghai"
	chatGPTExportBuyProximity          = 0.05
	chatGPTExportAggressiveBuyDiscount = 0.10
)

type chatGPTStockRecord struct {
	Symbol             string
	Name               string
	Status             string
	Currency           string
	Industry           string
	Shares             float64
	Cost               float64
	CurrentPrice       float64
	CurrentPriceDate   string
	PreviousClose      float64
	PreviousCloseDate  string
	MarketCap          *float64
	MarketCapCurrency  string
	IntrinsicValue     *float64
	FairValueRange     string
	TargetBuyPrice     *float64
	PriceLevels        *PriceLevels
	MarginOfSafety     *float64
	QualityScore       *float64
	BusinessModel      *float64
	Moat               *float64
	Governance         *float64
	FinancialQuality   *float64
	Action             string
	Risk               string
	UpdatedAt          string
	Notes              string
	Reports            []Report
	Dividend           *Dividend
	NetCash            *NetCashProfile
	OwnerCashFlowAudit *OwnerCashFlowAudit
	ResearchUpdates    []ResearchUpdate
	Financials         *Financials
	Plan               *PlanItem
	DecisionLogs       []DecisionLog
	MarketValueCNY     float64
	CostValueCNY       float64
	ProfitLossCNY      float64
	ProfitLossRate     *float64
	Weight             float64
}

func (s *Server) handleExportChatGPTContext(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	state, err := loadState()
	s.mu.Unlock()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "读取组合数据失败")
		return
	}

	generatedAt := chatGPTExportNow()
	body, err := buildChatGPTContextZip(state, generatedAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "生成 ChatGPT 档案失败")
		return
	}

	filename := fmt.Sprintf("portfolio-context-%s.zip", generatedAt.Format("20060102-150405"))
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func chatGPTExportNow() time.Time {
	location, err := time.LoadLocation(chatGPTExportTimezone)
	if err != nil {
		return time.Now()
	}
	return time.Now().In(location)
}

func buildChatGPTContextZip(state AppState, generatedAt time.Time) ([]byte, error) {
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)

	if err := addZipDirectory(writer, "stocks/", generatedAt); err != nil {
		_ = writer.Close()
		return nil, err
	}

	meta := chatGPTExportMeta(generatedAt)
	records, totalPositions, totalAssets := buildChatGPTStockRecords(state)

	files := []struct {
		name    string
		content string
	}{
		{"00_project_instructions.md", renderProjectInstructions(meta)},
		{"00_reference_tables.md", renderReferenceTables(meta, state, records, totalPositions, totalAssets)},
		{"01_portfolio_snapshot.md", renderPortfolioSnapshot(meta, state, records, totalPositions, totalAssets)},
		{"02_decision_rules.md", renderDecisionRules(meta, state)},
		{"03_watchlist_and_triggers.md", renderWatchlistAndTriggers(meta, records)},
		{"04_recent_decision_logs.md", renderRecentDecisionLogs(meta, state.DecisionLogs, 100)},
		{"05_master_lens_tables.md", renderMasterLensTables(meta, records)},
		{"06_risk_committee_memo.md", renderRiskCommitteeMemo(meta, state, records, totalPositions, totalAssets)},
		{"research_loop_guide.md", renderResearchLoopGuide(meta)},
		{"import_schema.md", renderImportSchema(meta)},
	}

	for _, file := range files {
		if err := addZipFile(writer, file.name, file.content, generatedAt); err != nil {
			_ = writer.Close()
			return nil, err
		}
	}

	usedPaths := make(map[string]int)
	for _, record := range records {
		path := uniqueZipPath(usedPaths, stockZipPath(record))
		if err := addZipFile(writer, path, renderStockMarkdown(meta, record), generatedAt); err != nil {
			_ = writer.Close()
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func addZipDirectory(writer *zip.Writer, name string, modTime time.Time) error {
	header := &zip.FileHeader{Name: name, Method: zip.Store}
	header.SetModTime(modTime)
	_, err := writer.CreateHeader(header)
	return err
}

func addZipFile(writer *zip.Writer, name string, content string, modTime time.Time) error {
	header := &zip.FileHeader{Name: name, Method: zip.Deflate}
	header.SetModTime(modTime)
	fileWriter, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.WriteString(fileWriter, content)
	return err
}

func chatGPTExportMeta(generatedAt time.Time) string {
	return fmt.Sprintf("---\ngeneratedAt: %s\nsource: holds_website data/portfolio.json + data/runtime/quotes.json\ntimezone: %s\n---\n\n", generatedAt.Format(time.RFC3339), chatGPTExportTimezone)
}

func buildChatGPTStockRecords(state AppState) ([]chatGPTStockRecord, float64, float64) {
	totalPositions := 0.0
	for _, holding := range state.Holdings {
		currency := firstNonEmpty(holding.Currency, expectedCurrency(holding.Symbol))
		totalPositions += holding.Shares * holding.CurrentPrice * chatGPTCurrencyRate(state, currency)
	}
	totalAssets := totalPositions + state.Cash
	if totalAssets <= 0 {
		totalAssets = totalPositions
	}

	records := make([]chatGPTStockRecord, 0, len(state.Holdings)+len(state.Candidates))
	seen := make(map[string]bool)
	for _, holding := range state.Holdings {
		key := normalizeSymbol(holding.Symbol)
		seen[key] = true
		records = append(records, chatGPTHoldingRecord(state, holding, totalAssets))
	}
	for _, candidate := range state.Candidates {
		key := normalizeSymbol(candidate.Symbol)
		if seen[key] {
			continue
		}
		seen[key] = true
		records = append(records, chatGPTCandidateRecord(state, candidate, totalAssets))
	}
	return records, totalPositions, totalAssets
}

func chatGPTHoldingRecord(state AppState, holding Holding, totalAssets float64) chatGPTStockRecord {
	currency := firstNonEmpty(holding.Currency, expectedCurrency(holding.Symbol))
	rate := chatGPTCurrencyRate(state, currency)
	marketValue := holding.Shares * holding.CurrentPrice * rate
	costValue := holding.Shares * holding.Cost * rate
	var profitLossRate *float64
	if costValue > 0 {
		value := (marketValue - costValue) / costValue
		profitLossRate = &value
	}
	weight := 0.0
	if totalAssets > 0 {
		weight = marketValue / totalAssets
	}
	targetBuyPrice := targetBuyPriceFromIntrinsicValue(holding.IntrinsicValue)
	return chatGPTStockRecord{
		Symbol:             holding.Symbol,
		Name:               holding.Name,
		Status:             "持仓",
		Currency:           currency,
		Industry:           holding.Industry,
		Shares:             holding.Shares,
		Cost:               holding.Cost,
		CurrentPrice:       holding.CurrentPrice,
		CurrentPriceDate:   holding.CurrentPriceDate,
		PreviousClose:      holding.PreviousClose,
		PreviousCloseDate:  holding.PreviousCloseDate,
		MarketCap:          holding.MarketCap,
		MarketCapCurrency:  holding.MarketCapCurrency,
		IntrinsicValue:     holding.IntrinsicValue,
		FairValueRange:     holding.FairValueRange,
		TargetBuyPrice:     targetBuyPrice,
		PriceLevels:        priceLevelsFromTarget(targetBuyPrice),
		MarginOfSafety:     marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, holding.MarginOfSafety),
		QualityScore:       holding.QualityScore,
		BusinessModel:      holding.BusinessModel,
		Moat:               holding.Moat,
		Governance:         holding.Governance,
		FinancialQuality:   holding.FinancialQuality,
		Action:             holding.Action,
		Risk:               holding.Risk,
		UpdatedAt:          holding.UpdatedAt,
		Notes:              holding.Notes,
		Reports:            holding.Reports,
		Dividend:           holding.Dividend,
		NetCash:            holding.NetCash,
		OwnerCashFlowAudit: holding.OwnerCashFlowAudit,
		ResearchUpdates:    holding.ResearchUpdates,
		Financials:         holding.Financials,
		Plan:               findPlanForDecisionLog(&state, holding.Symbol, holding.Name),
		DecisionLogs:       chatGPTLogsForStock(state.DecisionLogs, holding.Symbol),
		MarketValueCNY:     marketValue,
		CostValueCNY:       costValue,
		ProfitLossCNY:      marketValue - costValue,
		ProfitLossRate:     profitLossRate,
		Weight:             weight,
	}
}

func chatGPTCandidateRecord(state AppState, candidate Candidate, totalAssets float64) chatGPTStockRecord {
	targetBuyPrice := targetBuyPriceFromIntrinsicValue(candidate.IntrinsicValue)
	return chatGPTStockRecord{
		Symbol:             candidate.Symbol,
		Name:               candidate.Name,
		Status:             "候选",
		Currency:           firstNonEmpty(candidate.Currency, expectedCurrency(candidate.Symbol)),
		Industry:           candidate.Industry,
		CurrentPrice:       candidate.CurrentPrice,
		CurrentPriceDate:   candidate.CurrentPriceDate,
		PreviousClose:      candidate.PreviousClose,
		PreviousCloseDate:  candidate.PreviousCloseDate,
		MarketCap:          candidate.MarketCap,
		MarketCapCurrency:  candidate.MarketCapCurrency,
		IntrinsicValue:     candidate.IntrinsicValue,
		FairValueRange:     candidate.FairValueRange,
		TargetBuyPrice:     targetBuyPrice,
		PriceLevels:        priceLevelsFromTarget(targetBuyPrice),
		MarginOfSafety:     marginOfSafetyFromPrice(candidate.IntrinsicValue, candidate.CurrentPrice, candidate.MarginOfSafety),
		QualityScore:       candidate.QualityScore,
		BusinessModel:      candidate.BusinessModel,
		Moat:               candidate.Moat,
		Governance:         candidate.Governance,
		FinancialQuality:   candidate.FinancialQuality,
		Action:             candidate.Action,
		Risk:               candidate.Risk,
		UpdatedAt:          candidate.UpdatedAt,
		Notes:              candidate.Notes,
		Reports:            candidate.Reports,
		Dividend:           candidate.Dividend,
		NetCash:            candidate.NetCash,
		OwnerCashFlowAudit: candidate.OwnerCashFlowAudit,
		ResearchUpdates:    candidate.ResearchUpdates,
		Financials:         candidate.Financials,
		Plan:               findPlanForDecisionLog(&state, candidate.Symbol, candidate.Name),
		DecisionLogs:       chatGPTLogsForStock(state.DecisionLogs, candidate.Symbol),
	}
}

func chatGPTCurrencyRate(state AppState, currency string) float64 {
	normalized := strings.ToUpper(strings.TrimSpace(currency))
	if normalized == "" || normalized == "CNY" {
		return 1
	}
	if rate, ok := state.FX[normalized]; ok && rate > 0 {
		return rate
	}
	return 1
}

func renderProjectInstructions(meta string) string {
	return meta + `# ChatGPT 项目使用说明

你正在读取 holds_website 导出的完整组合上下文包。请把这些 Markdown 文件作为长期项目档案使用，用于后续股票深度研究、组合复盘、风险检查和导入 JSON 生成。

## 使用原则

- 优先读取 ` + "`00_reference_tables.md`" + ` 建立横向比较视图，再读取 ` + "`01_portfolio_snapshot.md`" + ` 理解当前现金、仓位、汇率、持仓和候选池。
- 做个股分析时读取 ` + "`05_master_lens_tables.md`" + `，先按“双策略”判断主策略、辅策略、过渡观察或风险排除。
- 做仓位和风险判断时读取 ` + "`06_risk_committee_memo.md`" + `，优先检查 70% 回报蓝筹主策略、30% 净现金烟蒂辅策略的偏离度。
- 深度研究单只股票前，先读取 ` + "`stocks/`" + ` 下对应股票档案，避免重复询问已经存在的成本、目标价、风险和历史决策。
- 事件或财报更新请先读取 ` + "`research_loop_guide.md`" + `，判断应该输出完整重估 ` + "`fullReview`" + ` 还是增量更新 ` + "`eventUpdate`" + `。
- 所有建议必须同时考虑综合回报盾、DCF安全边际、长期股东现金流评分、净现金保护、自由现金流、仓位和既有投资纪律。
- 如果研究结论需要回写网站，请输出符合 ` + "`import_schema.md`" + ` 的 JSON；不要输出散乱字段。
- 遇到价格、财报或新闻这类会变化的信息时，先说明信息时点，再给出结论。
- 不把短期波动当成买卖理由，除非它改变了估值、安全边际或基本面判断。

## 输出偏好

- 单股研究：给出结论、估值区间、买入/持有/卖出纪律、关键风险、需要持续跟踪的触发条件。
- 组合复盘：先看仓位和风险集中度，再看候选池是否有更优风险收益比。
- 导入网站：只输出一个 JSON 对象，字段遵循 ` + "`import_schema.md`" + `。
`
}

func renderReferenceTables(meta string, state AppState, records []chatGPTStockRecord, totalPositions float64, totalAssets float64) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString("# 参考表格\n\n")
	builder.WriteString("这份文件是给 ChatGPT 快速建立组合全局认知的索引层。先用表格做横向比较，再进入 `stocks/` 里的单股档案查证据和细节。\n\n")

	builder.WriteString("## 组合总览表\n\n")
	builder.WriteString("| 项目 | 数值 |\n| --- | --- |\n")
	builder.WriteString(fmt.Sprintf("| 现金 | %s |\n", mdCell(formatCurrency(state.Cash, "CNY"))))
	builder.WriteString(fmt.Sprintf("| 持仓市值 | %s |\n", mdCell(formatCurrency(totalPositions, "CNY"))))
	builder.WriteString(fmt.Sprintf("| 总资产 | %s |\n", mdCell(formatCurrency(totalAssets, "CNY"))))
	builder.WriteString(fmt.Sprintf("| 现金占比 | %s |\n", formatPercent(ratioOrZero(state.Cash, totalAssets))))
	builder.WriteString(fmt.Sprintf("| 持仓占比 | %s |\n", formatPercent(ratioOrZero(totalPositions, totalAssets))))
	builder.WriteString(fmt.Sprintf("| 持仓数量 | %d |\n", len(state.Holdings)))
	builder.WriteString(fmt.Sprintf("| 候选数量 | %d |\n", len(state.Candidates)))
	builder.WriteString(fmt.Sprintf("| 决策日志数量 | %d |\n", len(state.DecisionLogs)))

	builder.WriteString("\n## 持仓参考表\n\n")
	builder.WriteString("| 档案 | 状态 | 货币 | 股数 | 成本 | 最新价 | 市值(CNY) | 盈亏(CNY) | 盈亏率 | 仓位 | 安全边际 | 质量分 | 当前动作 |\n")
	builder.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |\n")
	for _, record := range records {
		if record.Status != "持仓" {
			continue
		}
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(chatGPTStockLink(record)),
			mdCell(record.Status),
			mdCell(record.Currency),
			formatNumber(record.Shares),
			formatNumber(record.Cost),
			formatNumber(record.CurrentPrice),
			formatNumber(record.MarketValueCNY),
			formatNumber(record.ProfitLossCNY),
			formatPercentPtr(record.ProfitLossRate),
			formatPercent(record.Weight),
			formatPercentPtr(record.MarginOfSafety),
			formatScorePtr(record.QualityScore),
			mdCell(record.Action),
		))
	}

	builder.WriteString("\n## 候选池参考表\n\n")
	builder.WriteString("| 档案 | 行业 | 货币 | 最新价 | 内在价值 | 目标买入价 | 安全边际 | 距离目标买入价 | 质量分 | 当前动作 |\n")
	builder.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- |\n")
	for _, record := range records {
		if record.Status != "候选" {
			continue
		}
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(chatGPTStockLink(record)),
			mdCell(record.Industry),
			mdCell(record.Currency),
			formatNumber(record.CurrentPrice),
			formatFloatPtr(record.IntrinsicValue),
			formatFloatPtr(record.TargetBuyPrice),
			formatPercentPtr(record.MarginOfSafety),
			formatTargetDistance(record.CurrentPrice, record.TargetBuyPrice),
			formatScorePtr(record.QualityScore),
			mdCell(record.Action),
		))
	}

	builder.WriteString("\n## 估值与安全边际表\n\n")
	builder.WriteString("| 状态 | 档案 | 最新价 | 内在价值 | 公允区间 | 目标买入价 | 安全边际 | 距离目标买入价 | 质量分 |\n")
	builder.WriteString("| --- | --- | ---: | ---: | --- | ---: | ---: | ---: | ---: |\n")
	for _, record := range records {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Status),
			mdCell(chatGPTStockLink(record)),
			formatNumber(record.CurrentPrice),
			formatFloatPtr(record.IntrinsicValue),
			mdCell(record.FairValueRange),
			formatFloatPtr(record.TargetBuyPrice),
			formatPercentPtr(record.MarginOfSafety),
			formatTargetDistance(record.CurrentPrice, record.TargetBuyPrice),
			formatScorePtr(record.QualityScore),
		))
	}

	builder.WriteString("\n## 执行计划表\n\n")
	builder.WriteString("| 排名 | 标的 | 名称 | 状态 | 优先级 | 建议 | 纪律 |\n")
	builder.WriteString("| ---: | --- | --- | --- | --- | --- | --- |\n")
	plans := sortedPlanItems(state.Plan)
	for _, plan := range plans {
		record := recordForPlan(records, plan)
		status := "-"
		if record != nil {
			status = record.Status
		}
		builder.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s | %s |\n",
			plan.Rank,
			mdCell(plan.Symbol),
			mdCell(plan.Name),
			mdCell(status),
			mdCell(plan.Priority),
			mdCell(plan.Advice),
			mdCell(plan.Discipline),
		))
	}
	if len(plans) == 0 {
		builder.WriteString("| - | - | - | - | - | - | - |\n")
	}

	builder.WriteString("\n## 风险与否决项表\n\n")
	builder.WriteString("| 状态 | 档案 | 行业 | 质量分 | 安全边际 | 主要风险 |\n")
	builder.WriteString("| --- | --- | --- | ---: | ---: | --- |\n")
	for _, record := range records {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Status),
			mdCell(chatGPTStockLink(record)),
			mdCell(record.Industry),
			formatScorePtr(record.QualityScore),
			formatPercentPtr(record.MarginOfSafety),
			mdCell(record.Risk),
		))
	}

	builder.WriteString("\n## 最近决策日志索引\n\n")
	builder.WriteString("| 日期 | 类型 | 标的 | 名称 | 决策 | 纪律 |\n")
	builder.WriteString("| --- | --- | --- | --- | --- | --- |\n")
	logs := sortedDecisionLogs(state.DecisionLogs)
	if len(logs) > 30 {
		logs = logs[:30]
	}
	for _, logItem := range logs {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			mdCell(logItem.Date),
			mdCell(logItem.Type),
			mdCell(logItem.Symbol),
			mdCell(logItem.Name),
			mdCell(logItem.Decision),
			mdCell(logItem.Discipline),
		))
	}
	if len(logs) == 0 {
		builder.WriteString("| - | - | - | - | - | - |\n")
	}

	return builder.String()
}

func renderPortfolioSnapshot(meta string, state AppState, records []chatGPTStockRecord, totalPositions float64, totalAssets float64) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString("# 组合快照\n\n")
	builder.WriteString("## 总览\n\n")
	builder.WriteString("| 项目 | 数值 |\n| --- | --- |\n")
	builder.WriteString(fmt.Sprintf("| 现金 | %s |\n", mdCell(formatCurrency(state.Cash, "CNY"))))
	builder.WriteString(fmt.Sprintf("| 持仓市值 | %s |\n", mdCell(formatCurrency(totalPositions, "CNY"))))
	builder.WriteString(fmt.Sprintf("| 总资产 | %s |\n", mdCell(formatCurrency(totalAssets, "CNY"))))
	builder.WriteString(fmt.Sprintf("| 持仓数量 | %d |\n", len(state.Holdings)))
	builder.WriteString(fmt.Sprintf("| 候选数量 | %d |\n", len(state.Candidates)))
	builder.WriteString(fmt.Sprintf("| 决策日志数量 | %d |\n", len(state.DecisionLogs)))
	builder.WriteString("\n## 汇率\n\n")
	builder.WriteString("| 货币 | CNY 汇率 |\n| --- | ---: |\n")
	keys := make([]string, 0, len(state.FX))
	for key := range state.FX {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("| %s | %.6f |\n", mdCell(key), state.FX[key]))
	}
	if len(keys) == 0 {
		builder.WriteString("| - | - |\n")
	}

	builder.WriteString("\n## 持仓总览\n\n")
	builder.WriteString("| 标的 | 名称 | 货币 | 股数 | 成本 | 最新价 | 市值(CNY) | 盈亏(CNY) | 仓位 | 安全边际 | 质量分 | 动作 |\n")
	builder.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |\n")
	for _, record := range records {
		if record.Status != "持仓" {
			continue
		}
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Symbol),
			mdCell(record.Name),
			mdCell(record.Currency),
			formatNumber(record.Shares),
			formatNumber(record.Cost),
			formatNumber(record.CurrentPrice),
			formatNumber(record.MarketValueCNY),
			formatNumber(record.ProfitLossCNY),
			formatPercent(record.Weight),
			formatPercentPtr(record.MarginOfSafety),
			formatScorePtr(record.QualityScore),
			mdCell(record.Action),
		))
	}

	builder.WriteString("\n## 候选池总览\n\n")
	builder.WriteString("| 标的 | 名称 | 货币 | 最新价 | 内在价值 | 目标买入价 | 安全边际 | 质量分 | 动作 | 风险 |\n")
	builder.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | --- | --- |\n")
	for _, record := range records {
		if record.Status != "候选" {
			continue
		}
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Symbol),
			mdCell(record.Name),
			mdCell(record.Currency),
			formatNumber(record.CurrentPrice),
			formatFloatPtr(record.IntrinsicValue),
			formatFloatPtr(record.TargetBuyPrice),
			formatPercentPtr(record.MarginOfSafety),
			formatScorePtr(record.QualityScore),
			mdCell(record.Action),
			mdCell(record.Risk),
		))
	}

	return builder.String()
}

func renderDecisionRules(meta string, state AppState) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString("# 投资纪律与决策规则\n\n")
	builder.WriteString("## 网站内规则\n\n")
	builder.WriteString("| 维度 | 分数 | 标准 |\n| --- | ---: | --- |\n")
	for _, rule := range state.Rules {
		builder.WriteString(fmt.Sprintf("| %s | %.1f | %s |\n", mdCell(rule.Dimension), rule.Score, mdCell(rule.Standard)))
	}
	if len(state.Rules) == 0 {
		builder.WriteString("| - | - | - |\n")
	}

	builder.WriteString("\n## 质量评分\n\n")
	builder.WriteString("| 维度 | 解释 |\n| --- | --- |\n")
	builder.WriteString("| 总分 | 对商业质量、估值、风险和执行纪律的综合评分。 |\n")
	builder.WriteString("| 商业模式 | 关注需求稳定性、定价权、增长空间和现金流韧性。 |\n")
	builder.WriteString("| 护城河 | 关注品牌、网络效应、成本优势、渠道、监管壁垒和切换成本。 |\n")
	builder.WriteString("| 治理 | 关注管理层资本配置、股东回报、财务透明度和利益一致性。 |\n")
	builder.WriteString("| 财务质量 | 关注利润率、资产负债表、现金流、ROE/ROIC 和周期波动。 |\n")

	builder.WriteString("\n## 买入纪律\n\n")
	builder.WriteString("- 主策略固定为 70% 目标仓位：自选池大盘蓝筹，A股综合回报率≥6%或H股综合回报率≥8%，最近完整财年口径达标，DCF安全边际≥15%。\n")
	builder.WriteString("- 主策略买入前长期股东现金流评分需达到 75/100；review 项按部分分计入，不再要求七项全部通过。\n")
	builder.WriteString("- 辅策略固定为 30% 目标仓位：账上净现金保护、折扣后净现金可验证，A股 ex-cash PE≤10，H股 ex-cash PE≤8，并优先看 ex-cash P/FCF、FCF yield 和FCF连续性。\n")
	builder.WriteString("- 买入前必须同时检查策略归属、综合回报盾、DCF边际、长期需求、资产耐久、再投资需求、分红FCF支持、净现金折扣、自由现金流和现有仓位。\n")
	builder.WriteString("- 对已经持仓的股票，新增买入必须证明风险收益比优于候选池替代项。\n")
	builder.WriteString("- 旧持仓未达新阈值时默认过渡观察，不自动触发卖出；新资金只进入主策略达标或辅策略烟蒂达标标的。\n")
	builder.WriteString("- 如果最新价格、财报或重大新闻缺失，先补研究，不直接给交易建议。\n")

	builder.WriteString("\n## 风险否决原则\n\n")
	builder.WriteString("- 财务造假、治理失信、核心竞争力永久受损、杠杆失控属于优先否决项。\n")
	builder.WriteString("- 对高仓位股票，除非安全边际显著扩大，否则优先控制集中度。\n")
	builder.WriteString("- 不能用短期情绪覆盖估值纪律；不能用低估值掩盖质量恶化。\n")
	return builder.String()
}

func renderWatchlistAndTriggers(meta string, records []chatGPTStockRecord) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString("# 候选池与触发条件\n\n")
	builder.WriteString("## 候选池\n\n")
	builder.WriteString("| 标的 | 名称 | 行业 | 最新价 | 目标买入价 | 内在价值 | 安全边际 | 价格距离目标 | 行动建议 | 主要风险 |\n")
	builder.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | --- | --- |\n")
	for _, record := range records {
		if record.Status != "候选" {
			continue
		}
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Symbol),
			mdCell(record.Name),
			mdCell(record.Industry),
			formatNumber(record.CurrentPrice),
			formatFloatPtr(record.TargetBuyPrice),
			formatFloatPtr(record.IntrinsicValue),
			formatPercentPtr(record.MarginOfSafety),
			formatTargetDistance(record.CurrentPrice, record.TargetBuyPrice),
			mdCell(record.Action),
			mdCell(record.Risk),
		))
	}

	builder.WriteString("\n## 全部标的触发器\n\n")
	builder.WriteString("| 状态 | 标的 | 名称 | 最新价 | 目标买入价 | 安全边际 | 执行计划 | 达标状态 |\n")
	builder.WriteString("| --- | --- | --- | ---: | ---: | ---: | --- | --- |\n")
	for _, record := range records {
		planAdvice := "-"
		planPriority := "-"
		if record.Plan != nil {
			planAdvice = firstNonEmpty(record.Plan.Advice, "-")
			planPriority = firstNonEmpty(record.Plan.Priority, "-")
		}
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Status),
			mdCell(record.Symbol),
			mdCell(record.Name),
			formatNumber(record.CurrentPrice),
			formatFloatPtr(record.TargetBuyPrice),
			formatPercentPtr(record.MarginOfSafety),
			mdCell(planAdvice),
			mdCell(planPriority),
		))
	}
	return builder.String()
}

func renderMasterLensTables(meta string, records []chatGPTStockRecord) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString("# 双策略评分表\n\n")
	builder.WriteString("本文件用于让 ChatGPT 先按回报蓝筹主策略与净现金烟蒂辅策略审查标的，再回到单股档案补证据。\n\n")

	builder.WriteString("## 主策略：回报蓝筹\n\n")
	builder.WriteString("| 状态 | 档案 | 市场 | 最新价 | 综合回报率 | 股息率 | 预估股息率 | 回报门槛 | DCF边际 | 长期评分 | 当前动作 |\n")
	builder.WriteString("| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |\n")
	for _, record := range records {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Status),
			mdCell(chatGPTStockLink(record)),
			mdCell(recordMarketKind(record)),
			formatNumber(record.CurrentPrice),
			formatShareholderReturnYield(record),
			formatStockDividendYield(record),
			formatForecastDividendYield(record.Dividend),
			formatPercent(recordDividendTarget(record)),
			formatPercentPtr(record.MarginOfSafety),
			mdCell(formatOwnerAuditConclusion(record.OwnerCashFlowAudit)),
			mdCell(record.Action),
		))
	}

	builder.WriteString("\n## 辅策略：净现金烟蒂\n\n")
	builder.WriteString("| 状态 | 档案 | 调整后净现金 | ex-cash PE | ex-cash P/FCF | FCF yield | FCF为正年数 | 折扣说明 |\n")
	builder.WriteString("| --- | --- | ---: | ---: | ---: | ---: | ---: | --- |\n")
	for _, record := range records {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Status),
			mdCell(chatGPTStockLink(record)),
			formatNetCashAmount(record.NetCash, "adjusted"),
			formatFloatPtr(netCashFloat(record.NetCash, "pe")),
			formatFloatPtr(netCashFloat(record.NetCash, "pfcf")),
			formatPercentPtr(netCashFloat(record.NetCash, "fcfYield")),
			formatIntPtr(netCashInt(record.NetCash, "fcfYears")),
			mdCell(formatNetCashReason(record.NetCash)),
		))
	}

	builder.WriteString("\n## 过渡观察与风险排除\n\n")
	builder.WriteString("| 状态 | 档案 | 质量分 | 安全边际 | 主要风险 | 下一步 |\n")
	builder.WriteString("| --- | --- | ---: | ---: | --- | --- |\n")
	for _, record := range records {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			mdCell(record.Status),
			mdCell(chatGPTStockLink(record)),
			formatScorePtr(record.QualityScore),
			formatPercentPtr(record.MarginOfSafety),
			mdCell(record.Risk),
			mdCell(firstNonEmpty(record.Action, record.Notes)),
		))
	}

	return builder.String()
}

func renderRiskCommitteeMemo(meta string, state AppState, records []chatGPTStockRecord, totalPositions float64, totalAssets float64) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString("# 风险委员会备忘录\n\n")
	builder.WriteString("霍华德马克斯视角不负责选出好股票，而是负责判断风险补偿、周期暴露和组合进攻/防守节奏。\n\n")

	cashRatio := ratioOrZero(state.Cash, totalAssets)
	positionsRatio := ratioOrZero(totalPositions, totalAssets)
	riskRecords := exportRiskReviewRecords(records)
	builder.WriteString("## 当前姿态\n\n")
	builder.WriteString("| 项目 | 数值 |\n| --- | --- |\n")
	builder.WriteString(fmt.Sprintf("| 建议姿态 | %s |\n", mdCell(exportPortfolioPosture(cashRatio, riskRecords))))
	builder.WriteString(fmt.Sprintf("| 现金比例 | %s |\n", formatPercent(cashRatio)))
	builder.WriteString(fmt.Sprintf("| 持仓比例 | %s |\n", formatPercent(positionsRatio)))
	builder.WriteString(fmt.Sprintf("| 风险复盘标的数 | %d |\n", len(riskRecords)))

	builder.WriteString("\n## 仓位集中度\n\n")
	builder.WriteString("| 档案 | 状态 | 仓位 | 市值(CNY) | 安全边际 | 风险主席判断 |\n")
	builder.WriteString("| --- | --- | ---: | ---: | ---: | --- |\n")
	for _, record := range exportTopWeightRecords(records, 8) {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			mdCell(chatGPTStockLink(record)),
			mdCell(record.Status),
			formatPercent(record.Weight),
			formatNumber(record.MarketValueCNY),
			formatPercentPtr(record.MarginOfSafety),
			mdCell(exportRiskStatus(record)),
		))
	}

	builder.WriteString("\n## 风险复盘清单\n\n")
	builder.WriteString("| 档案 | 状态 | 质量分 | 安全边际 | 触发原因 | 当前动作 |\n")
	builder.WriteString("| --- | --- | ---: | ---: | --- | --- |\n")
	for _, record := range riskRecords {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			mdCell(chatGPTStockLink(record)),
			mdCell(record.Status),
			formatScorePtr(record.QualityScore),
			formatPercentPtr(record.MarginOfSafety),
			mdCell(exportRiskReason(record)),
			mdCell(record.Action),
		))
	}
	if len(riskRecords) == 0 {
		builder.WriteString("| - | - | - | - | 暂无显著复盘项 | - |\n")
	}

	builder.WriteString("\n## 马克斯式问题清单\n\n")
	builder.WriteString("- 当前价格给出的风险补偿是否足够，还是仅仅因为资产质量好而想买？\n")
	builder.WriteString("- 组合是否过度暴露于同一周期、同一宏观变量或同一政策风险？\n")
	builder.WriteString("- 现金比例是否允许在更好赔率出现时行动？\n")
	builder.WriteString("- 高股息、低估值或预期差是否掩盖了治理、现金流或基本面恶化？\n")
	builder.WriteString("- 如果市场再下跌 20%，哪些持仓应该加仓，哪些只能被动承受？\n")
	return builder.String()
}

func renderRecentDecisionLogs(meta string, logs []DecisionLog, limit int) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString("# 最近决策日志\n\n")
	sortedLogs := sortedDecisionLogs(logs)
	if limit > 0 && len(sortedLogs) > limit {
		sortedLogs = sortedLogs[:limit]
	}
	for _, logItem := range sortedLogs {
		builder.WriteString(fmt.Sprintf("## %s %s %s\n\n", mdText(logItem.Date), mdText(logItem.Symbol), mdText(logItem.Name)))
		builder.WriteString(fmt.Sprintf("- 类型：%s\n", mdText(logItem.Type)))
		builder.WriteString(fmt.Sprintf("- 决策：%s\n", mdText(logItem.Decision)))
		builder.WriteString(fmt.Sprintf("- 纪律：%s\n", mdText(logItem.Discipline)))
		builder.WriteString(fmt.Sprintf("- 详情：%s\n", mdText(logItem.Detail)))
		builder.WriteString(fmt.Sprintf("- 价格：%s\n\n", formatDecisionLogPrice(logItem)))
	}
	if len(sortedLogs) == 0 {
		builder.WriteString("暂无决策日志。\n")
	}
	return builder.String()
}

func renderImportSchema(meta string) string {
	return meta + "# 网站导入 JSON Schema\n\n" + `ChatGPT 深度研究后，如需回写网站，请只输出一个 JSON 对象。不要 Markdown fences，不要解释文字；不要添加额外字段，未知数字字段使用 null。

## fullReview：完整重估

~~~json
{
  "updateType": "fullReview",
  "symbol": "0700.HK",
  "name": "腾讯控股",
  "asOf": "2026-05-09",
  "currency": "HKD",
  "industry": "互联网平台/游戏/广告/金融科技",
  "status": "未达标（安全边际<15%）",
  "action": "继续持有；新资金暂不追买，等待安全边际达标后再分批",
  "risk": "政策、地缘、AI投入回报周期和广告/游戏周期波动需折价",
  "valuationConfidence": "high",
  "valuation": {
    "intrinsicValue": 508,
    "fairValueRange": "HK$480-560",
    "targetBuyPrice": null,
    "marginOfSafety": 0.09
  },
  "quality": {
    "totalScore": 89,
    "businessModel": 28,
    "moat": 23,
    "governance": 17,
    "financialQuality": 21
  },
  "plan": {
    "rank": 1,
    "priority": "观察/低优先级",
    "advice": "等待安全边际达标后再分批，未达标不追买",
    "discipline": "优秀资产要求≥15%安全边际；未达标不追买"
  },
  "dividend": {
    "fiscalYear": "FY2025",
    "dividendPerShare": 4.5,
    "dividendCurrency": "HKD",
    "payoutRatio": 0.16,
    "reliability": "stable",
    "forecastFiscalYear": "FY2026E",
    "forecastPerShare": 5.2,
    "forecastCurrency": "HKD",
    "forecastYield": 0.083
  },
  "netCash": {
    "cashAndShortInvestments": 320000000000,
    "interestBearingDebt": 120000000000,
    "netCash": 200000000000,
    "currency": "HKD",
    "haircut": 0.7,
    "haircutReason": "平台现金流稳定但需保留监管和再投资折扣",
    "adjustedNetCash": 140000000000,
    "exCashPe": 13.5,
    "exCashPfcf": 14.2,
    "fcfYield": 0.065,
    "shareholderFcf": 9000000000,
    "shareholderFcfCurrency": "HKD",
    "shareholderFcfBasis": "普通股东 FCF：合并FCF扣除少数股东分流后口径",
    "consolidatedFcf": 12000000000,
    "minorityFcfAdjustment": 3000000000,
    "fcfPositiveYears": 5,
    "note": "净现金、FCF 和估值口径使用 FY2025 年报与当前市值。"
  },
  "ownerCashFlowAudit": {
    "tenYearDemand": { "status": "pass", "note": "核心产品/服务十年后仍有稳定需求。" },
    "assetDurability": { "status": "pass", "note": "品牌、资源或网络资产不易折旧。" },
    "maintenanceCapexLight": { "status": "review", "note": "需继续核实维持性资本开支。" },
    "dividendFcfSupport": { "status": "pass", "note": "分红由真实自由现金流覆盖。" },
    "dividendReinvestmentEfficiency": { "status": "review", "note": "当前估值对分红再投入效率一般。" },
    "roeRoicDurability": { "status": "pass", "note": "ROE/ROIC 有长期维持基础。" },
    "valuationSystemRisk": { "status": "pass", "note": "暂未发现行业估值体系永久改变。" }
  },
  "notes": "总结关键财务事实、估值假设、研究来源时点和需要跟踪的变化。"
}
~~~

## eventUpdate：事件/财报增量更新

~~~json
{
  "updateType": "eventUpdate",
  "symbol": "0700.HK",
  "name": "腾讯控股",
  "asOf": "2026-05-15",
  "event": {
    "type": "earnings",
    "title": "2026Q1 财报更新",
    "date": "2026-05-15",
    "source": "公司公告",
    "summary": "收入和自由现金流好于原假设，AI capex 继续上升。"
  },
  "impact": {
    "thesisChange": "minor",
    "valuationChange": "raise",
    "riskChange": "unchanged",
    "actionChange": "unchanged"
  },
  "updates": {
    "valuation": {
      "intrinsicValue": 550,
      "fairValueRange": "HK$500-610",
      "marginOfSafety": 0.15
    },
    "risk": "AI资本开支仍需跟踪，但短期现金流韧性增强",
    "notesAppend": "2026Q1 证实广告和游戏恢复，暂不改变买入纪律。"
  }
}
~~~

## 字段要求

- updateType 可为 fullReview 或 eventUpdate；旧 JSON 不写 updateType 时按 fullReview 处理。
- fullReview 用于年度/重大重估，会覆盖核心研究字段。
- eventUpdate 用于财报、公告、分红、回购、监管和重大新闻；只覆盖 updates 中明确给出的字段，未给字段保留网站原值。
- symbol 必填，使用网站现有代码格式，例如 0700.HK、000333.SZ。
- 金额字段使用该股票交易货币，不要换算成人民币，除非字段名明确写 CNY。
- valuation.intrinsicValue 是核心 DCF 估值输入；主策略要求显示 DCF 安全边际≥15%。
- marginOfSafety 使用小数，例如 0.25 表示 25%。
- 首买价默认按 intrinsicValue × 75% 计算，观察价为首买价 × 105%，重仓价为首买价 × 90%。
- 主策略综合回报盾：A股综合回报率≥6%，H股综合回报率≥8%；使用最近完整财年现金分红加回购相对总市值计算。
- ownerCashFlowAudit 会被网站折算成 100 分长期股东评分；主策略要求评分≥75。review 项给部分分，不再要求七项全部 pass。valuationSystemRisk=fail 仍会进入风险排除。
- 股息数据由“更新行情”尽量从行情源抓取；研究也可以提供 dividend.forecastFiscalYear、forecastPerShare、forecastCurrency、forecastYield 作为参考，但预估股息率不替代综合回报硬门槛。
- 辅策略烟蒂：提供 netCash 结构化字段，重点说明净现金折扣、折扣后净现金、ex-cash PE、ex-cash P/FCF、FCF yield 和 FCF 连续性。
- 净现金折扣约定：稳定分红100%，一般70%，弱/周期40%，重大风险0%。如果 haircut 为空，网站会按股息可靠性和风险文本自动分档。
- ownerCashFlowAudit 七项 status 只能是 pass、review、fail；评分权重为十年需求18、资产耐久14、轻再投资12、分红FCF18、再投资效率12、ROE/ROIC14、估值体系12。
- quality.totalScore 应等于 businessModel + moat + governance + financialQuality。
- asOf 必须为 YYYY-MM-DD。
- plan 不要写 symbol；网站会用顶层 symbol 关联执行计划。
- 行情字段 currentPrice、previousClose 和日期由网站 runtime quote 文件负责，不通过研究导入更新。
- 如果 symbol 匹配现有持仓，会更新持仓研究字段；匹配候选股会更新候选字段；新标的会加入候选池。
`
}

func renderStockMarkdown(meta string, record chatGPTStockRecord) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString(fmt.Sprintf("# %s %s\n\n", mdText(record.Symbol), mdText(record.Name)))
	builder.WriteString("## ChatGPT 分析入口\n\n")
	builder.WriteString(fmt.Sprintf("- 当前结论：%s\n", mdText(firstNonEmpty(record.Action, record.Status))))
	builder.WriteString(fmt.Sprintf("- 执行纪律：%s\n", mdText(planDiscipline(record.Plan))))
	builder.WriteString(fmt.Sprintf("- 本次分析必须先复核：最新价日期 %s、安全边际 %s、综合回报率 %s、长期股东评分 %s。\n",
		mdText(record.CurrentPriceDate),
		formatPercentPtr(record.MarginOfSafety),
		formatShareholderReturnYield(record),
		mdText(formatOwnerAuditConclusion(record.OwnerCashFlowAudit)),
	))
	builder.WriteString("- 若只是财报/公告/事件增量，请输出 `eventUpdate`；若需要系统性重估，请输出 `fullReview`。\n\n")
	builder.WriteString("## 当前状态\n\n")
	builder.WriteString("| 项目 | 数值 |\n| --- | --- |\n")
	builder.WriteString(fmt.Sprintf("| 状态 | %s |\n", mdCell(record.Status)))
	builder.WriteString(fmt.Sprintf("| 行业 | %s |\n", mdCell(record.Industry)))
	builder.WriteString(fmt.Sprintf("| 货币 | %s |\n", mdCell(record.Currency)))
	builder.WriteString(fmt.Sprintf("| 股数 | %s |\n", formatNumber(record.Shares)))
	builder.WriteString(fmt.Sprintf("| 成本 | %s |\n", formatNumber(record.Cost)))
	builder.WriteString(fmt.Sprintf("| 最新价 | %s |\n", formatNumber(record.CurrentPrice)))
	builder.WriteString(fmt.Sprintf("| 最新价日期 | %s |\n", mdCell(record.CurrentPriceDate)))
	builder.WriteString(fmt.Sprintf("| 昨收 | %s |\n", formatNumber(record.PreviousClose)))
	builder.WriteString(fmt.Sprintf("| 昨收日期 | %s |\n", mdCell(record.PreviousCloseDate)))
	builder.WriteString(fmt.Sprintf("| 市值(CNY) | %s |\n", formatNumber(record.MarketValueCNY)))
	builder.WriteString(fmt.Sprintf("| 成本市值(CNY) | %s |\n", formatNumber(record.CostValueCNY)))
	builder.WriteString(fmt.Sprintf("| 盈亏(CNY) | %s |\n", formatNumber(record.ProfitLossCNY)))
	builder.WriteString(fmt.Sprintf("| 盈亏率 | %s |\n", formatPercentPtr(record.ProfitLossRate)))
	builder.WriteString(fmt.Sprintf("| 仓位 | %s |\n", formatPercent(record.Weight)))

	builder.WriteString("\n## 估值\n\n")
	builder.WriteString("| 项目 | 数值 |\n| --- | --- |\n")
	builder.WriteString(fmt.Sprintf("| 内在价值 | %s |\n", formatFloatPtr(record.IntrinsicValue)))
	builder.WriteString(fmt.Sprintf("| 公允区间 | %s |\n", mdCell(record.FairValueRange)))
	builder.WriteString(fmt.Sprintf("| 目标买入价 | %s |\n", formatFloatPtr(record.TargetBuyPrice)))
	builder.WriteString(fmt.Sprintf("| 观察价 | %s |\n", formatPriceLevel(record.PriceLevels, "watch")))
	builder.WriteString(fmt.Sprintf("| 首买价 | %s |\n", formatPriceLevel(record.PriceLevels, "initial")))
	builder.WriteString(fmt.Sprintf("| 重仓价 | %s |\n", formatPriceLevel(record.PriceLevels, "aggressive")))
	builder.WriteString(fmt.Sprintf("| 安全边际 | %s |\n", formatPercentPtr(record.MarginOfSafety)))

	builder.WriteString("\n## 长期股东现金流评分\n\n")
	builder.WriteString(fmt.Sprintf("总评：%s\n\n", mdText(formatOwnerAuditConclusion(record.OwnerCashFlowAudit))))
	builder.WriteString("| 项目 | 状态 | 分数 | 说明 |\n| --- | --- | ---: | --- |\n")
	for _, item := range ownerAuditRows(record.OwnerCashFlowAudit) {
		builder.WriteString(fmt.Sprintf("| %s | %s | %d/%d | %s |\n", mdCell(item.Label), mdCell(item.Status), item.Points, item.MaxPoints, mdCell(item.Note)))
	}

	builder.WriteString("\n## 股息与现金流\n\n")
	builder.WriteString("| 项目 | 数值 |\n| --- | --- |\n")
	builder.WriteString(fmt.Sprintf("| 财年 | %s |\n", mdCell(dividendFiscalYear(record.Dividend))))
	builder.WriteString(fmt.Sprintf("| 每股分红 | %s |\n", formatDividendPerShare(record.Dividend)))
	builder.WriteString(fmt.Sprintf("| 股息率 | %s |\n", formatStockDividendYield(record)))
	builder.WriteString(fmt.Sprintf("| 回报门槛 | %s |\n", formatPercent(recordDividendTarget(record))))
	builder.WriteString(fmt.Sprintf("| 预估财年 | %s |\n", mdCell(dividendForecastFiscalYear(record.Dividend))))
	builder.WriteString(fmt.Sprintf("| 预估每股 | %s |\n", formatForecastDividendPerShare(record.Dividend)))
	builder.WriteString(fmt.Sprintf("| 预估股息率 | %s |\n", formatForecastDividendYield(record.Dividend)))
	builder.WriteString(fmt.Sprintf("| 综合回报率 | %s |\n", formatShareholderReturnYield(record)))
	builder.WriteString(fmt.Sprintf("| 预估年度现金 | %s |\n", formatStockEstimatedAnnualCash(record)))

	builder.WriteString("\n## 净现金烟蒂\n\n")
	builder.WriteString("| 项目 | 数值 |\n| --- | --- |\n")
	builder.WriteString(fmt.Sprintf("| 现金/短投 | %s |\n", formatNetCashAmount(record.NetCash, "cash")))
	builder.WriteString(fmt.Sprintf("| 有息债务 | %s |\n", formatNetCashAmount(record.NetCash, "debt")))
	builder.WriteString(fmt.Sprintf("| 净现金 | %s |\n", formatNetCashAmount(record.NetCash, "net")))
	builder.WriteString(fmt.Sprintf("| 净现金折扣 | %s |\n", formatNetCashHaircut(record.NetCash)))
	builder.WriteString(fmt.Sprintf("| 折扣说明 | %s |\n", mdCell(formatNetCashReason(record.NetCash))))
	builder.WriteString(fmt.Sprintf("| 调整后净现金 | %s |\n", formatNetCashAmount(record.NetCash, "adjusted")))
	builder.WriteString(fmt.Sprintf("| ex-cash PE | %s |\n", formatFloatPtr(netCashFloat(record.NetCash, "pe"))))
	builder.WriteString(fmt.Sprintf("| ex-cash P/FCF | %s |\n", formatFloatPtr(netCashFloat(record.NetCash, "pfcf"))))
	builder.WriteString(fmt.Sprintf("| 普通股东 FCF | %s |\n", formatNetCashShareholderFCF(record.NetCash)))
	builder.WriteString(fmt.Sprintf("| FCF yield | %s |\n", formatPercentPtr(netCashFloat(record.NetCash, "fcfYield"))))
	builder.WriteString(fmt.Sprintf("| FCF为正年数 | %s |\n", formatIntPtr(netCashInt(record.NetCash, "fcfYears"))))

	builder.WriteString("\n## 质量\n\n")
	builder.WriteString("| 维度 | 分数 |\n| --- | ---: |\n")
	builder.WriteString(fmt.Sprintf("| 总分 | %s |\n", formatScorePtr(record.QualityScore)))
	builder.WriteString(fmt.Sprintf("| 商业模式 | %s |\n", formatScorePtr(record.BusinessModel)))
	builder.WriteString(fmt.Sprintf("| 护城河 | %s |\n", formatScorePtr(record.Moat)))
	builder.WriteString(fmt.Sprintf("| 治理 | %s |\n", formatScorePtr(record.Governance)))
	builder.WriteString(fmt.Sprintf("| 财务质量 | %s |\n", formatScorePtr(record.FinancialQuality)))

	builder.WriteString(renderFinancialsMarkdown(record.Financials))

	builder.WriteString("\n## 决策\n\n")
	builder.WriteString(fmt.Sprintf("- 当前动作：%s\n", mdText(record.Action)))
	if record.Plan != nil {
		builder.WriteString(fmt.Sprintf("- 达标状态：%s\n", mdText(record.Plan.Priority)))
		builder.WriteString(fmt.Sprintf("- 执行计划：%s\n", mdText(record.Plan.Advice)))
		builder.WriteString(fmt.Sprintf("- 计划纪律：%s\n", mdText(record.Plan.Discipline)))
		builder.WriteString(fmt.Sprintf("- 执行排名：%d\n", record.Plan.Rank))
	} else {
		builder.WriteString("- 达标状态：-\n")
		builder.WriteString("- 执行计划：-\n")
	}
	builder.WriteString(fmt.Sprintf("- 主要风险：%s\n", mdText(record.Risk)))

	builder.WriteString(renderResearchUpdatesMarkdown(record.ResearchUpdates))

	builder.WriteString("\n## 研究材料\n\n")
	builder.WriteString("### Notes\n\n")
	builder.WriteString(mdText(record.Notes))
	builder.WriteString("\n\n### Reports\n\n")
	builder.WriteString("| 日期 | 周期 | 类型 | 标题 | 来源 | URL |\n")
	builder.WriteString("| --- | --- | --- | --- | --- | --- |\n")
	for _, report := range record.Reports {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			mdCell(report.Date),
			mdCell(report.Period),
			mdCell(report.Kind),
			mdCell(report.Title),
			mdCell(report.Source),
			mdCell(report.URL),
		))
	}
	if len(record.Reports) == 0 {
		builder.WriteString("| - | - | - | - | - | - |\n")
	}

	builder.WriteString("\n### 最近决策日志\n\n")
	logs := sortedDecisionLogs(record.DecisionLogs)
	for _, logItem := range logs {
		builder.WriteString(fmt.Sprintf("#### %s %s\n\n", mdText(logItem.Date), mdText(logItem.Type)))
		builder.WriteString(fmt.Sprintf("- 决策：%s\n", mdText(logItem.Decision)))
		builder.WriteString(fmt.Sprintf("- 纪律：%s\n", mdText(logItem.Discipline)))
		builder.WriteString(fmt.Sprintf("- 详情：%s\n", mdText(logItem.Detail)))
		builder.WriteString(fmt.Sprintf("- 价格：%s\n\n", formatDecisionLogPrice(logItem)))
	}
	if len(logs) == 0 {
		builder.WriteString("暂无该股票决策日志。\n")
	}
	return builder.String()
}

func renderFinancialsMarkdown(financials *Financials) string {
	if financials == nil || len(financials.Annual) == 0 {
		return "\n## 多年财务数据\n\n暂无结构化多年财务数据。\n"
	}

	var builder strings.Builder
	builder.WriteString("\n## 多年财务数据\n\n")
	builder.WriteString(fmt.Sprintf("- 来源：%s\n", mdText(financials.Source)))
	builder.WriteString(fmt.Sprintf("- 更新时间：%s\n", mdText(financials.UpdatedAt)))
	if financials.Valuation != nil {
		builder.WriteString(fmt.Sprintf("- 当前 PE/PB/PEG：%s / %s / %s\n",
			formatFloatPtr(financials.Valuation.PE),
			formatFloatPtr(financials.Valuation.PB),
			formatFloatPtr(financials.Valuation.PEG),
		))
	}
	builder.WriteString("\n| 年度 | 收入 | 收入同比 | 归母利润 | 利润同比 | 经营现金流 | FCF | ROE | ROIC | 负债率 | PE | PB |\n")
	builder.WriteString("| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |\n")
	for _, item := range financials.Annual {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(firstNonEmpty(item.FiscalYear, item.ReportDate)),
			formatFloatPtr(item.Revenue),
			formatPercentPtr(item.RevenueYoY),
			formatFloatPtr(item.NetProfit),
			formatPercentPtr(item.NetProfitYoY),
			formatFloatPtr(item.OperatingCashFlow),
			formatFloatPtr(item.FreeCashFlow),
			formatPercentPtr(item.ROE),
			formatPercentPtr(item.ROIC),
			formatPercentPtr(item.DebtRatio),
			formatFloatPtr(item.PEAtCurrentPrice),
			formatFloatPtr(item.PBAtCurrentPrice),
		))
	}
	return builder.String()
}

func renderResearchUpdatesMarkdown(updates []ResearchUpdate) string {
	var builder strings.Builder
	builder.WriteString("\n## 最近研究更新\n\n")
	if len(updates) == 0 {
		builder.WriteString("暂无事件/财报增量更新记录。\n")
		return builder.String()
	}
	sorted := append([]ResearchUpdate(nil), updates...)
	sort.SliceStable(sorted, func(i, j int) bool {
		left := firstNonEmpty(sorted[i].ImportedAt, sorted[i].AsOf, sorted[i].Event.Date)
		right := firstNonEmpty(sorted[j].ImportedAt, sorted[j].AsOf, sorted[j].Event.Date)
		return left > right
	})
	if len(sorted) > 12 {
		sorted = sorted[:12]
	}
	builder.WriteString("| 导入时间 | 事件日期 | 类型 | 标题 | 影响 | 更新字段 | 摘要 |\n")
	builder.WriteString("| --- | --- | --- | --- | --- | --- | --- |\n")
	for _, item := range sorted {
		impact := strings.Join(nonEmptyStrings(
			impactText("thesis", item.Impact.ThesisChange),
			impactText("valuation", item.Impact.ValuationChange),
			impactText("risk", item.Impact.RiskChange),
			impactText("action", item.Impact.ActionChange),
		), "；")
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
			mdCell(item.ImportedAt),
			mdCell(firstNonEmpty(item.Event.Date, item.AsOf)),
			mdCell(firstNonEmpty(item.Event.Type, item.UpdateType)),
			mdCell(item.Event.Title),
			mdCell(firstNonEmpty(impact, "-")),
			mdCell(strings.Join(item.ChangedFields, "/")),
			mdCell(firstNonEmpty(item.Event.Summary, item.Summary, item.NotesAppend)),
		))
	}
	return builder.String()
}

func impactText(label string, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return label + "=" + value
}

func nonEmptyStrings(values ...string) []string {
	result := []string{}
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			result = append(result, strings.TrimSpace(value))
		}
	}
	return result
}

func chatGPTLogsForStock(logs []DecisionLog, symbol string) []DecisionLog {
	normalized := normalizeSymbol(symbol)
	filtered := make([]DecisionLog, 0)
	for _, logItem := range logs {
		if normalizeSymbol(logItem.Symbol) == normalized {
			filtered = append(filtered, logItem)
		}
	}
	sorted := sortedDecisionLogs(filtered)
	if len(sorted) > 20 {
		return sorted[:20]
	}
	return sorted
}

func sortedDecisionLogs(logs []DecisionLog) []DecisionLog {
	sorted := append([]DecisionLog(nil), logs...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Date == sorted[j].Date {
			return sorted[i].ID > sorted[j].ID
		}
		return sorted[i].Date > sorted[j].Date
	})
	return sorted
}

func sortedPlanItems(plan []PlanItem) []PlanItem {
	sorted := append([]PlanItem(nil), plan...)
	sort.SliceStable(sorted, func(i, j int) bool {
		leftRank := sorted[i].Rank
		rightRank := sorted[j].Rank
		if leftRank == 0 {
			leftRank = 999
		}
		if rightRank == 0 {
			rightRank = 999
		}
		if leftRank == rightRank {
			return sorted[i].Name < sorted[j].Name
		}
		return leftRank < rightRank
	})
	return sorted
}

func renderResearchLoopGuide(meta string) string {
	return meta + `# ChatGPT 研究闭环指南

本网站使用双模式导入，让研究可以持续迭代，而不是每次从零开始。

## 先读什么

1. 先读 ` + "`00_project_instructions.md`" + ` 和 ` + "`00_reference_tables.md`" + ` 建立组合上下文。
2. 再读目标股票的 ` + "`stocks/{symbol}.md`" + `，确认已有结论、估值、执行纪律、最近研究更新和财报材料。
3. 如果要回写网站，最后读 ` + "`import_schema.md`" + `，只输出一个 JSON 对象。

## 什么时候用 fullReview

用于年度报告、重大策略重估、首次纳入候选池，或原有内在价值/质量分/风险判断需要系统性重写。` + "`fullReview`" + ` 会覆盖股票当前核心研究字段。

## 什么时候用 eventUpdate

用于季度财报、分红/回购公告、监管或经营事件、业绩预告、重要管理层变化等增量信息。` + "`eventUpdate`" + ` 只写明本次需要改变的字段；未写字段会保留网站原值，避免误清空历史分析。

## 输出原则

- 不回写 currentPrice、previousClose 或收盘日期；行情由网站 runtime quote 文件维护。
- 每次事件更新必须说明事件日期、来源、摘要和影响判断。
- 若事件不改变估值或动作，也可以只追加 researchUpdates，用于保留复盘脉络。
- 不要孤立分析；结论必须解释相对既有内在价值、执行纪律和历史判断的变化。
`
}

func exportGrahamStatus(record chatGPTStockRecord) string {
	margin := exportPtrFloat(record.MarginOfSafety)
	financial := exportPtrFloat(record.FinancialQuality)
	if exportHasMajorRisk(record) {
		return "风险否决/先复盘"
	}
	if margin >= 0.25 && financial >= 20 {
		return "防御合格"
	}
	if margin >= 0.15 {
		return "勉强可观察"
	}
	return "不够便宜"
}

func exportBuffettStatus(record chatGPTStockRecord) string {
	quality := exportPtrFloat(record.QualityScore)
	moat := exportPtrFloat(record.Moat)
	financial := exportPtrFloat(record.FinancialQuality)
	margin := exportPtrFloat(record.MarginOfSafety)
	if exportHasMajorRisk(record) || (quality > 0 && quality < 75) {
		return "能力圈外/特殊观察"
	}
	if quality >= 85 && moat >= 22 && financial >= 20 && margin >= 0.1 {
		return "长期核心候选"
	}
	if quality >= 78 {
		return "好生意等价格"
	}
	return "普通机会"
}

func exportLynchStatus(record chatGPTStockRecord) string {
	if exportHasMajorRisk(record) {
		return "等待验证"
	}
	quality := exportPtrFloat(record.QualityScore)
	margin := exportPtrFloat(record.MarginOfSafety)
	hasGrowthCue := exportHasGrowthCue(record)
	if quality >= 82 && margin >= 0.1 && hasGrowthCue {
		return "成长故事清晰"
	}
	if (quality >= 75 && hasGrowthCue) || strings.Contains(exportLynchCategory(record), "预期差") {
		return "预期差可验证"
	}
	if !hasGrowthCue {
		return "故事不足"
	}
	return "继续跟踪"
}

func exportRiskStatus(record chatGPTStockRecord) string {
	margin := exportPtrFloat(record.MarginOfSafety)
	if exportHasMajorRisk(record) {
		return "风险复盘"
	}
	if margin >= 0.2 {
		return "补偿足够"
	}
	if margin >= 0.15 {
		return "仅观察"
	}
	return "等待补偿"
}

func exportLynchCategory(record chatGPTStockRecord) string {
	text := exportRecordText(record)
	if exportHasMajorRisk(record) {
		return "困境反转/问题股"
	}
	if strings.Contains(text, "预期差") || strings.Contains(text, "反转") || strings.Contains(text, "复苏") || strings.Contains(text, "修复") || strings.Contains(text, "验证") || strings.Contains(text, "重估") {
		return "困境反转/预期差"
	}
	if containsAnyText(text, []string{"新能源", "AI", "机器人", "自动化", "科技", "互联网", "智能"}) {
		return "快速增长/高波动"
	}
	if containsAnyText(text, []string{"油气", "银行", "地产", "物业", "航空", "周期", "白酒", "乳制品", "家电"}) {
		return "稳定增长/周期敏感"
	}
	return "稳定增长/普通成长"
}

func exportGrowthCue(record chatGPTStockRecord) string {
	text := exportRecordText(record)
	cues := []string{"增长", "复苏", "修复", "预期差", "验证", "扩张", "现金流", "回购", "分红", "海外", "AI", "机器人", "新能源", "重估"}
	found := make([]string, 0, 3)
	for _, cue := range cues {
		if strings.Contains(text, cue) {
			found = append(found, cue)
		}
		if len(found) >= 3 {
			break
		}
	}
	if len(found) == 0 {
		return "成长线索待补充"
	}
	return strings.Join(found, "、")
}

func exportHasGrowthCue(record chatGPTStockRecord) bool {
	return exportGrowthCue(record) != "成长线索待补充"
}

func exportRiskReviewRecords(records []chatGPTStockRecord) []chatGPTStockRecord {
	riskRecords := make([]chatGPTStockRecord, 0)
	for _, record := range records {
		margin := exportPtrFloat(record.MarginOfSafety)
		quality := exportPtrFloat(record.QualityScore)
		if exportRiskStatus(record) == "风险复盘" || (quality > 0 && quality < 75) || (margin > 0 && margin < 0.1) {
			riskRecords = append(riskRecords, record)
		}
	}
	sort.SliceStable(riskRecords, func(i, j int) bool {
		return exportRiskSortScore(riskRecords[i]) > exportRiskSortScore(riskRecords[j])
	})
	return riskRecords
}

func exportTopWeightRecords(records []chatGPTStockRecord, limit int) []chatGPTStockRecord {
	held := make([]chatGPTStockRecord, 0)
	for _, record := range records {
		if record.Status == "持仓" {
			held = append(held, record)
		}
	}
	sort.SliceStable(held, func(i, j int) bool {
		return held[i].Weight > held[j].Weight
	})
	if limit > 0 && len(held) > limit {
		return held[:limit]
	}
	return held
}

func exportPortfolioPosture(cashRatio float64, riskRecords []chatGPTStockRecord) string {
	if len(riskRecords) >= 3 {
		return "先复盘风险，暂缓进攻"
	}
	if cashRatio >= 0.35 {
		return "现金充足，等待高赔率机会"
	}
	if cashRatio < 0.12 {
		return "现金偏低，控制新增仓位"
	}
	return "防守等待，小额验证"
}

func exportRiskReason(record chatGPTStockRecord) string {
	reasons := make([]string, 0, 3)
	margin := exportPtrFloat(record.MarginOfSafety)
	quality := exportPtrFloat(record.QualityScore)
	if exportHasMajorRisk(record) {
		reasons = append(reasons, "重大风险/治理或数据可信度")
	}
	if quality > 0 && quality < 75 {
		reasons = append(reasons, "质量分低于 75")
	}
	if margin > 0 && margin < 0.1 {
		reasons = append(reasons, "安全边际低于 10%")
	}
	if len(reasons) == 0 {
		reasons = append(reasons, exportRiskStatus(record))
	}
	return strings.Join(reasons, "；")
}

func exportRiskSortScore(record chatGPTStockRecord) float64 {
	score := record.Weight * 100
	if exportHasMajorRisk(record) {
		score += 10
	}
	quality := exportPtrFloat(record.QualityScore)
	if quality > 0 && quality < 75 {
		score += 5
	}
	margin := exportPtrFloat(record.MarginOfSafety)
	if margin > 0 && margin < 0.1 {
		score += 3
	}
	return score
}

func exportHasMajorRisk(record chatGPTStockRecord) bool {
	text := exportRecordText(record)
	text = strings.ReplaceAll(text, "无一票否决", "")
	text = strings.ReplaceAll(text, "无立即否决", "")
	text = strings.ReplaceAll(text, "没有一票否决", "")
	return containsAnyText(text, []string{"否决", "停牌", "造假", "调查", "重大风险", "低可信", "内控", "退市", "财报可信", "治理风险", "治理与财务可靠性", "质量分<75", "低于75"})
}

func exportRecordText(record chatGPTStockRecord) string {
	return strings.Join([]string{record.Industry, record.Action, record.Risk, record.Notes, record.FairValueRange}, " ")
}

func containsAnyText(text string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(text, needle) {
			return true
		}
	}
	return false
}

func exportPtrFloat(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func recordForPlan(records []chatGPTStockRecord, plan PlanItem) *chatGPTStockRecord {
	normalizedSymbol := normalizeSymbol(plan.Symbol)
	for i := range records {
		if normalizedSymbol != "" && normalizeSymbol(records[i].Symbol) == normalizedSymbol {
			return &records[i]
		}
	}

	planName := strings.TrimSpace(plan.Name)
	for i := range records {
		recordName := strings.TrimSpace(records[i].Name)
		if planName != "" && recordName != "" && (strings.EqualFold(planName, recordName) || strings.Contains(planName, recordName) || strings.Contains(recordName, planName)) {
			return &records[i]
		}
	}
	return nil
}

func chatGPTStockLink(record chatGPTStockRecord) string {
	label := strings.TrimSpace(record.Symbol + " " + record.Name)
	if label == "" {
		label = "股票档案"
	}
	return fmt.Sprintf("[%s](%s)", label, stockZipPath(record))
}

func stockZipPath(record chatGPTStockRecord) string {
	symbol := strings.ToUpper(strings.TrimSpace(record.Symbol))
	if symbol == "" {
		symbol = "stock"
	}
	return "stocks/" + sanitizeZipFileName(symbol) + ".md"
}

func uniqueZipPath(used map[string]int, path string) string {
	if _, ok := used[path]; !ok {
		used[path] = 1
		return path
	}
	used[path]++
	extension := ".md"
	base := strings.TrimSuffix(path, extension)
	return fmt.Sprintf("%s_%d%s", base, used[path], extension)
}

func sanitizeZipFileName(name string) string {
	name = strings.TrimSpace(name)
	invalid := `/\:*?"<>|`
	var builder strings.Builder
	lastWasSpace := false
	for _, char := range name {
		if strings.ContainsRune(invalid, char) || char < 32 {
			if !lastWasSpace {
				builder.WriteByte('_')
				lastWasSpace = true
			}
			continue
		}
		if char == '\t' || char == '\n' || char == '\r' {
			if !lastWasSpace {
				builder.WriteByte('_')
				lastWasSpace = true
			}
			continue
		}
		builder.WriteRune(char)
		lastWasSpace = false
	}
	clean := strings.Trim(builder.String(), " ._")
	if clean == "" {
		return "stock"
	}
	return clean
}

func formatCurrency(value float64, currency string) string {
	return fmt.Sprintf("%s %.2f", currency, value)
}

func ratioOrZero(value float64, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return value / total
}

func formatNumber(value float64) string {
	if value == 0 {
		return "-"
	}
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func formatFloatPtr(value *float64) string {
	if value == nil || *value == 0 {
		return "-"
	}
	return strconv.FormatFloat(*value, 'f', 2, 64)
}

func priceLevelsFromTarget(targetBuyPrice *float64) *PriceLevels {
	if targetBuyPrice == nil || *targetBuyPrice <= 0 {
		return nil
	}
	watchPrice := *targetBuyPrice * (1 + chatGPTExportBuyProximity)
	initialBuyPrice := *targetBuyPrice
	aggressiveBuyPrice := *targetBuyPrice * (1 - chatGPTExportAggressiveBuyDiscount)
	return &PriceLevels{
		WatchPrice:         &watchPrice,
		InitialBuyPrice:    &initialBuyPrice,
		AggressiveBuyPrice: &aggressiveBuyPrice,
	}
}

func formatPriceLevel(levels *PriceLevels, level string) string {
	if levels == nil {
		return "-"
	}
	switch level {
	case "watch":
		return formatFloatPtr(levels.WatchPrice)
	case "initial":
		return formatFloatPtr(levels.InitialBuyPrice)
	case "aggressive":
		return formatFloatPtr(levels.AggressiveBuyPrice)
	default:
		return "-"
	}
}

func dividendFiscalYear(dividend *Dividend) string {
	if dividend == nil || strings.TrimSpace(dividend.FiscalYear) == "" {
		return "-"
	}
	return strings.TrimSpace(dividend.FiscalYear)
}

func dividendCurrencyCode(dividend *Dividend) string {
	if dividend == nil || strings.TrimSpace(dividend.DividendCurrency) == "" {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(dividend.DividendCurrency))
}

func cashDividendCurrencyCode(dividend *Dividend) string {
	if dividend == nil || strings.TrimSpace(dividend.CashDividendCurrency) == "" {
		return dividendCurrencyCode(dividend)
	}
	return strings.ToUpper(strings.TrimSpace(dividend.CashDividendCurrency))
}

func formatDividendPerShare(dividend *Dividend) string {
	if dividend == nil || dividend.DividendPerShare == nil {
		return "-"
	}
	currency := dividendCurrencyCode(dividend)
	if currency == "" {
		return formatFloatPtr(dividend.DividendPerShare)
	}
	return formatCurrency(*dividend.DividendPerShare, currency)
}

func recordMarketKind(record chatGPTStockRecord) string {
	symbol := strings.ToUpper(strings.TrimSpace(record.Symbol))
	if strings.HasSuffix(symbol, ".HK") || strings.EqualFold(record.Currency, "HKD") {
		return "HK"
	}
	return "A"
}

func recordDividendTarget(record chatGPTStockRecord) float64 {
	if recordMarketKind(record) == "HK" {
		return 0.08
	}
	return 0.06
}

func dividendForecastFiscalYear(dividend *Dividend) string {
	if dividend == nil || strings.TrimSpace(dividend.ForecastFiscalYear) == "" {
		return "-"
	}
	return strings.TrimSpace(dividend.ForecastFiscalYear)
}

func forecastDividendCurrencyCode(dividend *Dividend) string {
	if dividend == nil || strings.TrimSpace(dividend.ForecastCurrency) == "" {
		return dividendCurrencyCode(dividend)
	}
	return strings.ToUpper(strings.TrimSpace(dividend.ForecastCurrency))
}

func formatForecastDividendPerShare(dividend *Dividend) string {
	if dividend == nil || dividend.ForecastPerShare == nil {
		return "-"
	}
	currency := forecastDividendCurrencyCode(dividend)
	if currency == "" {
		return formatFloatPtr(dividend.ForecastPerShare)
	}
	return formatCurrency(*dividend.ForecastPerShare, currency)
}

func formatForecastDividendYield(dividend *Dividend) string {
	if dividend == nil {
		return "-"
	}
	return formatPercentPtr(dividend.ForecastYield)
}

func formatStockDividendYield(record chatGPTStockRecord) string {
	if record.Dividend == nil {
		return "-"
	}
	if record.Dividend.CashDividendTotal != nil && record.MarketCap != nil && *record.Dividend.CashDividendTotal > 0 && *record.MarketCap > 0 {
		dividendCurrency := cashDividendCurrencyCode(record.Dividend)
		marketCapCurrency := firstNonEmpty(record.MarketCapCurrency, record.Currency)
		if dividendCurrency != "" && marketCapCurrency != "" && !strings.EqualFold(dividendCurrency, marketCapCurrency) {
			return "-"
		}
		value := *record.Dividend.CashDividendTotal / *record.MarketCap
		return formatPercentPtr(&value)
	}
	if record.Dividend.DividendPerShare == nil || record.CurrentPrice <= 0 {
		return "-"
	}
	dividendCurrency := dividendCurrencyCode(record.Dividend)
	if dividendCurrency != "" && record.Currency != "" && !strings.EqualFold(dividendCurrency, record.Currency) {
		return "-"
	}
	value := *record.Dividend.DividendPerShare / record.CurrentPrice
	return formatPercentPtr(&value)
}

func formatShareholderReturnYield(record chatGPTStockRecord) string {
	if record.Dividend == nil {
		return "-"
	}
	if record.Dividend.CashDividendTotal != nil && record.MarketCap != nil && *record.Dividend.CashDividendTotal > 0 && *record.MarketCap > 0 {
		dividendCurrency := cashDividendCurrencyCode(record.Dividend)
		marketCapCurrency := firstNonEmpty(record.MarketCapCurrency, record.Currency)
		if dividendCurrency != "" && marketCapCurrency != "" && !strings.EqualFold(dividendCurrency, marketCapCurrency) {
			return "-"
		}
		buyback := 0.0
		if record.Dividend.BuybackAmount != nil && *record.Dividend.BuybackAmount > 0 {
			buybackCurrency := firstNonEmpty(record.Dividend.BuybackCurrency, record.Currency)
			if buybackCurrency != "" && marketCapCurrency != "" && !strings.EqualFold(buybackCurrency, marketCapCurrency) {
				return "-"
			}
			buyback = *record.Dividend.BuybackAmount
		}
		value := (*record.Dividend.CashDividendTotal + buyback) / *record.MarketCap
		return formatPercentPtr(&value)
	}
	return formatStockDividendYield(record)
}

func formatStockEstimatedAnnualCash(record chatGPTStockRecord) string {
	if record.Dividend == nil {
		return "-"
	}
	currency := dividendCurrencyCode(record.Dividend)
	if record.Dividend.DividendPerShare != nil && record.Shares > 0 {
		value := *record.Dividend.DividendPerShare * record.Shares
		if currency == "" {
			return formatNumber(value)
		}
		return formatCurrency(value, currency)
	}
	if record.Dividend.EstimatedAnnualCash == nil {
		return "-"
	}
	if currency == "" {
		return formatFloatPtr(record.Dividend.EstimatedAnnualCash)
	}
	return formatCurrency(*record.Dividend.EstimatedAnnualCash, currency)
}

func netCashCurrencyCode(netCash *NetCashProfile) string {
	if netCash == nil || strings.TrimSpace(netCash.Currency) == "" {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(netCash.Currency))
}

func netCashFloat(netCash *NetCashProfile, field string) *float64 {
	if netCash == nil {
		return nil
	}
	switch field {
	case "cash":
		return netCash.CashAndShortInvestments
	case "debt":
		return netCash.InterestBearingDebt
	case "net":
		return netCash.NetCash
	case "haircut":
		return netCash.Haircut
	case "adjusted":
		return netCash.AdjustedNetCash
	case "pe":
		return netCash.ExCashPE
	case "pfcf":
		return netCash.ExCashPFCF
	case "fcfYield":
		return netCash.FCFYield
	case "shareholderFcf":
		return netCash.ShareholderFCF
	default:
		return nil
	}
}

func netCashInt(netCash *NetCashProfile, field string) *int {
	if netCash == nil {
		return nil
	}
	if field == "fcfYears" {
		return netCash.FCFPositiveYears
	}
	return nil
}

func formatNetCashAmount(netCash *NetCashProfile, field string) string {
	value := netCashFloat(netCash, field)
	if value == nil || *value == 0 {
		return "-"
	}
	currency := netCashCurrencyCode(netCash)
	if currency == "" {
		return formatFloatPtr(value)
	}
	return formatCurrency(*value, currency)
}

func formatNetCashHaircut(netCash *NetCashProfile) string {
	return formatPercentPtr(netCashFloat(netCash, "haircut"))
}

func formatNetCashShareholderFCF(netCash *NetCashProfile) string {
	if netCash == nil || netCash.ShareholderFCF == nil {
		return "-"
	}
	currency := strings.ToUpper(strings.TrimSpace(netCash.ShareholderFCFCurrency))
	if currency == "" {
		currency = netCashCurrencyCode(netCash)
	}
	return fmt.Sprintf("%s %s", currency, formatFloatPtr(netCash.ShareholderFCF))
}

func formatNetCashReason(netCash *NetCashProfile) string {
	if netCash == nil {
		return "-"
	}
	return firstNonEmpty(strings.TrimSpace(netCash.HaircutReason), strings.TrimSpace(netCash.Note), "-")
}

type ownerAuditRow struct {
	Label     string
	Status    string
	Points    int
	MaxPoints int
	Note      string
}

type ownerAuditField struct {
	Key    string
	Label  string
	Weight int
}

func ownerAuditFields() []ownerAuditField {
	return []ownerAuditField{
		{Key: "tenYearDemand", Label: "十年需求", Weight: 18},
		{Key: "assetDurability", Label: "资产耐久", Weight: 14},
		{Key: "maintenanceCapexLight", Label: "轻再投资", Weight: 12},
		{Key: "dividendFcfSupport", Label: "分红FCF", Weight: 18},
		{Key: "dividendReinvestmentEfficiency", Label: "再投资效率", Weight: 12},
		{Key: "roeRoicDurability", Label: "ROE/ROIC", Weight: 14},
		{Key: "valuationSystemRisk", Label: "估值体系", Weight: 12},
	}
}

func ownerAuditRows(audit *OwnerCashFlowAudit) []ownerAuditRow {
	fields := ownerAuditFields()
	result := make([]ownerAuditRow, 0, len(fields))
	_, hasAudit := ownerAuditScore(audit)
	for _, field := range fields {
		item := ownerAuditItem(audit, field.Key)
		result = append(result, ownerAuditRow{
			Label:     field.Label,
			Status:    formatOwnerAuditStatus(item.Status),
			Points:    ownerAuditItemPoints(item.Status, field.Weight, hasAudit),
			MaxPoints: field.Weight,
			Note:      firstNonEmpty(strings.TrimSpace(item.Note), "待补充"),
		})
	}
	return result
}

func ownerAuditItem(audit *OwnerCashFlowAudit, key string) OwnerAuditItem {
	if audit == nil {
		return OwnerAuditItem{Status: "review"}
	}
	switch key {
	case "tenYearDemand":
		return audit.TenYearDemand
	case "assetDurability":
		return audit.AssetDurability
	case "maintenanceCapexLight":
		return audit.MaintenanceCapexLight
	case "dividendFcfSupport":
		return audit.DividendFCFSupport
	case "dividendReinvestmentEfficiency":
		return audit.DividendReinvestmentEfficiency
	case "roeRoicDurability":
		return audit.RoeRoicDurability
	case "valuationSystemRisk":
		return audit.ValuationSystemRisk
	default:
		return OwnerAuditItem{Status: "review"}
	}
}

func formatOwnerAuditConclusion(audit *OwnerCashFlowAudit) string {
	score, hasAudit := ownerAuditScore(audit)
	if !hasAudit {
		return "待评分"
	}
	if score >= 85 {
		return fmt.Sprintf("%d/100 长期股东强", score)
	}
	if score >= 75 {
		return fmt.Sprintf("%d/100 长期股东达标", score)
	}
	if score >= 60 {
		return fmt.Sprintf("%d/100 长期股东观察", score)
	}
	return fmt.Sprintf("%d/100 长期股东偏弱", score)
}

func ownerAuditScore(audit *OwnerCashFlowAudit) (int, bool) {
	hasAudit := false
	points := 0
	maxPoints := 0
	for _, field := range ownerAuditFields() {
		item := ownerAuditItem(audit, field.Key)
		if strings.TrimSpace(item.Status) != "" || strings.TrimSpace(item.Note) != "" {
			hasAudit = true
		}
		points += ownerAuditItemPoints(item.Status, field.Weight, audit != nil)
		maxPoints += field.Weight
	}
	if !hasAudit || maxPoints <= 0 {
		return 0, false
	}
	return (points*100 + maxPoints/2) / maxPoints, true
}

func ownerAuditItemPoints(status string, weight int, hasAudit bool) int {
	if !hasAudit {
		return 0
	}
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pass":
		return weight
	case "fail":
		return 0
	default:
		return (weight*60 + 50) / 100
	}
}

func formatOwnerAuditStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pass":
		return "通过"
	case "fail":
		return "失败"
	default:
		return "复核"
	}
}

func formatScorePtr(value *float64) string {
	if value == nil || *value == 0 {
		return "-"
	}
	return strconv.FormatFloat(*value, 'f', 1, 64)
}

func formatPercent(value float64) string {
	if value == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", value*100)
}

func formatPercentPtr(value *float64) string {
	if value == nil || *value == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", *value*100)
}

func formatIntPtr(value *int) string {
	if value == nil {
		return "-"
	}
	return strconv.Itoa(*value)
}

func formatTargetDistance(currentPrice float64, targetPrice *float64) string {
	if currentPrice <= 0 || targetPrice == nil || *targetPrice <= 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", (currentPrice-*targetPrice)/(*targetPrice)*100)
}

func formatDecisionLogPrice(logItem DecisionLog) string {
	if logItem.Price == nil || *logItem.Price <= 0 {
		return "-"
	}
	return strings.TrimSpace(fmt.Sprintf("%s %s", logItem.Currency, formatFloatPtr(logItem.Price)))
}

func mdCell(value string) string {
	value = mdText(value)
	value = strings.ReplaceAll(value, "|", `\|`)
	value = strings.ReplaceAll(value, "\r\n", "<br>")
	value = strings.ReplaceAll(value, "\n", "<br>")
	return value
}

func mdText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	return value
}
