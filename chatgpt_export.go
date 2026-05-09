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
	Symbol            string
	Name              string
	Status            string
	Currency          string
	Industry          string
	Shares            float64
	Cost              float64
	CurrentPrice      float64
	CurrentPriceDate  string
	PreviousClose     float64
	PreviousCloseDate string
	MarketCap         *float64
	MarketCapCurrency string
	IntrinsicValue    *float64
	FairValueRange    string
	TargetBuyPrice    *float64
	PriceLevels       *PriceLevels
	MarginOfSafety    *float64
	QualityScore      *float64
	BusinessModel     *float64
	Moat              *float64
	Governance        *float64
	FinancialQuality  *float64
	Action            string
	Risk              string
	UpdatedAt         string
	Notes             string
	Reports           []Report
	Dividend          *Dividend
	Plan              *PlanItem
	DecisionLogs      []DecisionLog
	MarketValueCNY    float64
	CostValueCNY      float64
	ProfitLossCNY     float64
	ProfitLossRate    *float64
	Weight            float64
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
		path := uniqueZipPath(usedPaths, "stocks/"+sanitizeZipFileName(record.Symbol+"_"+record.Name)+".md")
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
	return fmt.Sprintf("---\ngeneratedAt: %s\nsource: holds_website data/portfolio.json\ntimezone: %s\n---\n\n", generatedAt.Format(time.RFC3339), chatGPTExportTimezone)
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
		Symbol:            holding.Symbol,
		Name:              holding.Name,
		Status:            "持仓",
		Currency:          currency,
		Industry:          holding.Industry,
		Shares:            holding.Shares,
		Cost:              holding.Cost,
		CurrentPrice:      holding.CurrentPrice,
		CurrentPriceDate:  holding.CurrentPriceDate,
		PreviousClose:     holding.PreviousClose,
		PreviousCloseDate: holding.PreviousCloseDate,
		MarketCap:         holding.MarketCap,
		MarketCapCurrency: holding.MarketCapCurrency,
		IntrinsicValue:    holding.IntrinsicValue,
		FairValueRange:    holding.FairValueRange,
		TargetBuyPrice:    targetBuyPrice,
		PriceLevels:       priceLevelsFromTarget(targetBuyPrice),
		MarginOfSafety:    marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, holding.MarginOfSafety),
		QualityScore:      holding.QualityScore,
		BusinessModel:     holding.BusinessModel,
		Moat:              holding.Moat,
		Governance:        holding.Governance,
		FinancialQuality:  holding.FinancialQuality,
		Action:            holding.Action,
		Risk:              holding.Risk,
		UpdatedAt:         holding.UpdatedAt,
		Notes:             holding.Notes,
		Reports:           holding.Reports,
		Dividend:          holding.Dividend,
		Plan:              findPlanForDecisionLog(&state, holding.Symbol, holding.Name),
		DecisionLogs:      chatGPTLogsForStock(state.DecisionLogs, holding.Symbol),
		MarketValueCNY:    marketValue,
		CostValueCNY:      costValue,
		ProfitLossCNY:     marketValue - costValue,
		ProfitLossRate:    profitLossRate,
		Weight:            weight,
	}
}

func chatGPTCandidateRecord(state AppState, candidate Candidate, totalAssets float64) chatGPTStockRecord {
	targetBuyPrice := targetBuyPriceFromIntrinsicValue(candidate.IntrinsicValue)
	return chatGPTStockRecord{
		Symbol:            candidate.Symbol,
		Name:              candidate.Name,
		Status:            "候选",
		Currency:          firstNonEmpty(candidate.Currency, expectedCurrency(candidate.Symbol)),
		Industry:          candidate.Industry,
		CurrentPrice:      candidate.CurrentPrice,
		CurrentPriceDate:  candidate.CurrentPriceDate,
		PreviousClose:     candidate.PreviousClose,
		PreviousCloseDate: candidate.PreviousCloseDate,
		MarketCap:         candidate.MarketCap,
		MarketCapCurrency: candidate.MarketCapCurrency,
		IntrinsicValue:    candidate.IntrinsicValue,
		FairValueRange:    candidate.FairValueRange,
		TargetBuyPrice:    targetBuyPrice,
		PriceLevels:       priceLevelsFromTarget(targetBuyPrice),
		MarginOfSafety:    marginOfSafetyFromPrice(candidate.IntrinsicValue, candidate.CurrentPrice, candidate.MarginOfSafety),
		QualityScore:      candidate.QualityScore,
		BusinessModel:     candidate.BusinessModel,
		Moat:              candidate.Moat,
		Governance:        candidate.Governance,
		FinancialQuality:  candidate.FinancialQuality,
		Action:            candidate.Action,
		Risk:              candidate.Risk,
		UpdatedAt:         candidate.UpdatedAt,
		Notes:             candidate.Notes,
		Reports:           candidate.Reports,
		Dividend:          candidate.Dividend,
		Plan:              findPlanForDecisionLog(&state, candidate.Symbol, candidate.Name),
		DecisionLogs:      chatGPTLogsForStock(state.DecisionLogs, candidate.Symbol),
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
- 深度研究单只股票前，先读取 ` + "`stocks/`" + ` 下对应股票档案，避免重复询问已经存在的成本、目标价、风险和历史决策。
- 所有建议必须同时考虑估值、安全边际、质量评分、仓位和既有投资纪律。
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
	builder.WriteString("- 买入前必须同时检查质量分、目标买入价、安全边际和现有仓位。\n")
	builder.WriteString("- 对已经持仓的股票，新增买入必须证明风险收益比优于候选池替代项。\n")
	builder.WriteString("- 候选股只有在价格进入目标区间、风险没有恶化、研究结论仍有效时才考虑执行。\n")
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
	return meta + "# 网站导入 JSON Schema\n\n" + `ChatGPT 深度研究后，如需回写网站，请只输出一个 JSON 对象。不要添加 Markdown 代码块以外的解释文字；不要添加额外字段，未知数字字段使用 null。

## 顶层结构（单只股票）

~~~json
{
  "symbol": "0700.HK",
  "name": "腾讯控股",
  "asOf": "2026-05-09",
  "currency": "HKD",
  "industry": "互联网平台/游戏/广告/金融科技",
  "status": "未达标（安全边际<15%）",
  "action": "继续持有；新资金暂不追买，等待安全边际达标后再分批",
  "risk": "政策、地缘、AI投入回报周期和广告/游戏周期波动需折价",
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
  "notes": "总结关键财务事实、估值假设、研究来源时点和需要跟踪的变化。"
}
~~~

## 字段要求

- symbol 必填，使用网站现有代码格式，例如 0700.HK、000333.SZ。
- 金额字段使用该股票交易货币，不要换算成人民币，除非字段名明确写 CNY。
- valuation.intrinsicValue 是核心估值输入；targetBuyPrice、priceLevels、dividend、dividendYield、estimatedAnnualCash 不需要提供，网站会自行计算或抓取。
- marginOfSafety 使用小数，例如 0.25 表示 25%。
- 首买价默认按 intrinsicValue × 75% 计算，观察价为首买价 × 105%，重仓价为首买价 × 90%。
- 股息数据由“更新行情”尽量从行情源抓取；股息率按最新完整财年现金分红总额 ÷ 公司总市值计算，综合回报率按现金分红总额 + 回购金额 ÷ 公司总市值计算。
- quality.totalScore 应等于 businessModel + moat + governance + financialQuality。
- asOf 必须为 YYYY-MM-DD。
- plan 不要写 symbol；网站会用顶层 symbol 关联执行计划。
- 行情字段 currentPrice、previousClose 和日期由网站“更新行情”负责，不通过研究导入更新。
- 如果 symbol 匹配现有持仓，会更新持仓研究字段；匹配候选股会更新候选字段；新标的会加入候选池。
`
}

func renderStockMarkdown(meta string, record chatGPTStockRecord) string {
	var builder strings.Builder
	builder.WriteString(meta)
	builder.WriteString(fmt.Sprintf("# %s %s\n\n", mdText(record.Symbol), mdText(record.Name)))
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

	builder.WriteString("\n## 股息与现金流\n\n")
	builder.WriteString("| 项目 | 数值 |\n| --- | --- |\n")
	builder.WriteString(fmt.Sprintf("| 财年 | %s |\n", mdCell(dividendFiscalYear(record.Dividend))))
	builder.WriteString(fmt.Sprintf("| 每股分红 | %s |\n", formatDividendPerShare(record.Dividend)))
	builder.WriteString(fmt.Sprintf("| 股息率 | %s |\n", formatStockDividendYield(record)))
	builder.WriteString(fmt.Sprintf("| 综合回报率 | %s |\n", formatShareholderReturnYield(record)))
	builder.WriteString(fmt.Sprintf("| 预估年度现金 | %s |\n", formatStockEstimatedAnnualCash(record)))

	builder.WriteString("\n## 质量\n\n")
	builder.WriteString("| 维度 | 分数 |\n| --- | ---: |\n")
	builder.WriteString(fmt.Sprintf("| 总分 | %s |\n", formatScorePtr(record.QualityScore)))
	builder.WriteString(fmt.Sprintf("| 商业模式 | %s |\n", formatScorePtr(record.BusinessModel)))
	builder.WriteString(fmt.Sprintf("| 护城河 | %s |\n", formatScorePtr(record.Moat)))
	builder.WriteString(fmt.Sprintf("| 治理 | %s |\n", formatScorePtr(record.Governance)))
	builder.WriteString(fmt.Sprintf("| 财务质量 | %s |\n", formatScorePtr(record.FinancialQuality)))

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
	path := "stocks/" + sanitizeZipFileName(record.Symbol+"_"+record.Name) + ".md"
	return fmt.Sprintf("[%s](%s)", label, path)
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
