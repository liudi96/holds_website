package main

import (
	"strings"
	"testing"
)

func TestAnalysisPageIsRemoved(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireContains(t, html, `data-page="overview"`)
	requireContains(t, html, `data-page="holdings"`)
	requireContains(t, html, `data-page="funds"`)
	requireContains(t, html, `data-page="trades"`)
	requireContains(t, html, `data-view="holdings" title="股票"`)
	requireContains(t, html, `data-view="funds" title="基金"`)
	requireNotContains(t, html, `data-page="screener"`)
	requireNotContains(t, html, `data-view="screener"`)
	requireNotContains(t, html, `id="sunny30Section"`)
	requireNotContains(t, html, `id="sunny30Body"`)
	requireNotContains(t, html, `id="sunny30CandidateDialog"`)
	requireNotContains(t, html, `id="openSunny30Candidate"`)
	requireNotContains(t, html, `data-page="portfolio"`)
	requireNotContains(t, html, `data-page="valuation"`)
	requireNotContains(t, html, `valuation-module-panel`)
	requireNotContains(t, html, `id="valuationModuleList"`)
	requireNotContains(t, html, `id="updateValuationHistory"`)

	holdingsStart := strings.Index(html, `data-page="holdings"`)
	fundsStart := strings.Index(html, `data-page="funds"`)
	tradesStart := strings.Index(html, `data-page="trades"`)
	overviewStart := strings.Index(html, `data-page="overview"`)
	allocationStart := strings.Index(html, `overview-allocation-panel`)
	mastersStart := strings.Index(html, `id="masters"`)
	fundPanelStart := strings.Index(html, `fund-holdings-panel`)
	etfTrackerStart := strings.Index(html, `id="etfRuleTracker"`)
	if overviewStart < 0 || holdingsStart < 0 || fundsStart < 0 || tradesStart < 0 || allocationStart < 0 || mastersStart < 0 || fundPanelStart < 0 || etfTrackerStart < 0 {
		t.Fatalf("missing required page anchors")
	}
	if !(overviewStart < allocationStart && allocationStart < etfTrackerStart && etfTrackerStart < holdingsStart) {
		t.Fatalf("expected ETF rule tracker under overview allocation panel")
	}
	if !(holdingsStart < mastersStart && mastersStart < fundsStart) {
		t.Fatalf("expected stock holdings table inside holdings page")
	}
	if !(fundsStart < fundPanelStart && fundPanelStart < tradesStart) {
		t.Fatalf("expected fund holdings panel inside funds page before logs page")
	}
	stockPageSection := extractBetween(t, html, `data-page="holdings"`, `data-page="funds"`)
	fundPageSection := extractBetween(t, html, `data-page="funds"`, `data-page="trades"`)
	requireContains(t, stockPageSection, `<h2>股票</h2>`)
	requireNotContains(t, stockPageSection, `id="fundsBody"`)
	requireContains(t, fundPageSection, `<h2>基金</h2>`)
	requireContains(t, fundPageSection, `id="fundsBody"`)
}

func TestRemovedAnalysisRoutesFallBackToOverview(t *testing.T) {
	js := readTextFile(t, "app.js")

	requireContains(t, js, `if (view === "screener" || view === "sunny30" || view === "candidates")`)
	requireContains(t, js, `return { view: "overview", page: "overview" }`)
	requireNotContains(t, js, `return { view: "screener", page: "screener" }`)
	requireNotContains(t, js, `if (route.page === "screener")`)
	overviewRender := extractBetween(t, js, `if (route.page === "overview") {`, `if (route.page === "holdings") {`)
	requireContains(t, overviewRender, `renderEtfRuleTracker();`)
	requireNotContains(t, js, `renderValuationModule`)
	requireNotContains(t, js, `valuationModuleList`)
	requireNotContains(t, js, `updateValuationHistory`)
	requireNotContains(t, js, `view === "valuation"`)
	requireNotContains(t, js, `route.page === "valuation"`)
	requireNotContains(t, js, `route.page === "portfolio"`)
}

func TestStockAndFundHoldingRoutesRenderSeparately(t *testing.T) {
	js := readTextFile(t, "app.js")

	requireContains(t, js, `holdings: "股票"`)
	requireContains(t, js, `funds: "基金"`)
	requireContains(t, js, `return { view: "funds", page: "funds" }`)
	requireContains(t, js, `if (route.page === "funds")`)
	stockRender := extractBetween(t, js, `if (route.page === "holdings") {`, `if (route.page === "funds") {`)
	fundRender := extractBetween(t, js, `if (route.page === "funds") {`, `if (route.page === "trades") {`)
	requireContains(t, stockRender, `renderPositions(positions);`)
	requireContains(t, stockRender, `renderMastersPage(positions);`)
	requireNotContains(t, stockRender, `renderFunds(`)
	requireContains(t, fundRender, `renderFunds(fundPositions);`)
	requireNotContains(t, fundRender, `renderPositions(`)
}

func TestEtfRuleTrackerRendersOnOverviewPage(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")
	css := readTextFile(t, "styles.css")

	overviewSection := extractBetween(t, html, `data-page="overview"`, `data-page="holdings"`)
	requireContains(t, overviewSection, `id="etfRuleTracker"`)
	requireContains(t, html, `id="etfRuleTracker"`)
	requireContains(t, html, `id="etfBuyDialog"`)
	requireContains(t, html, `id="etfBuyForm"`)
	requireContains(t, html, `id="etfBuyFundLabel"`)
	requireContains(t, html, `name="amount"`)
	requireContains(t, html, `aria-label="ETF追踪"`)
	requireContains(t, js, `<h2>ETF追踪</h2>`)
	requireNotContains(t, html, `<h2>股票追踪</h2>`)
	requireNotContains(t, js, `按每月4周执行；五档`)
	requireNotContains(t, js, `ETF_RULE_TRACKER_NOTES`)
	requireNotContains(t, js, `etf-rule-notes`)
	requireContains(t, js, `ETF_RULE_TRACKER_RULES`)
	requireContains(t, js, `etfRuleStatuses`)
	requireContains(t, js, `renderEtfRuleTracker();`)
	requireContains(t, js, `status?.complete`)
	requireContains(t, js, `renderEtfRuleMetric`)
	requireContains(t, js, `renderEtfRuleRulebook`)
	requireContains(t, js, `ETF_ALLOCATION_POOL_BASE`)
	requireContains(t, js, `ETF_ALLOCATION_MONTHLY_INFLOW`)
	requireContains(t, js, `ETF_RULE_TRADING_DAYS_PER_WEEK`)
	requireContains(t, js, `function etfRuleDailyAmount`)
	requireContains(t, js, `function etfAllocationSnapshot`)
	requireContains(t, js, `renderEtfAllocationBar`)
	requireContains(t, js, `renderEtfPoolProgress`)
	requireContains(t, js, `renderEtfRuleActionItem`)
	requireContains(t, js, `etf-rule-action-list`)
	requireContains(t, js, `etf-rule-detail-layer`)
	requireContains(t, js, `日计划`)
	requireContains(t, js, `今日计划`)
	requireContains(t, js, `data-etf-rule-buy`)
	requireContains(t, js, `function openEtfBuyDialog`)
	requireContains(t, js, `function tradeFromEtfRuleBuyForm`)
	requireContains(t, js, `async function saveEtfRuleBuy`)
	requireContains(t, js, `async function saveTradeRecord`)
	for _, symbol := range []string{"022434", "018738", "008163", "021000"} {
		requireContains(t, js, symbol)
	}
	for _, fundName := range []string{
		"南方中证A500ETF联接A",
		"博时标普500ETF联接E(人民币)",
		"南方标普红利低波50ETF联接A",
		"南方纳斯达克100指数发起(QDII)I",
	} {
		requireContains(t, js, fundName)
	}
	for _, level := range []string{`key: "quarter"`, `key: "half"`, `key: "one"`, `key: "oneHalf"`, `key: "two"`} {
		requireContains(t, js, level)
	}
	for _, condition := range []string{
		"PE分位>80%；若回撤<5%且高估则保持限速",
		"PE分位>95%；最低档不再下调",
		"PE分位20%—40%；回撤<15%则降为1倍",
		"PE分位>85%；或70%—85%且回撤<5%后限速",
		"PE分位20%—40%；或<20%且回撤<30%后限速",
		`quarter: 4200, half: 8400, one: 16800, oneHalf: 25200, two: 33600`,
		`quarter: 700, half: 1400, one: 2800, oneHalf: 4200, two: 5600`,
	} {
		requireContains(t, js, condition)
	}
	requireNotContains(t, js, `ETF_RULE_TRACKER_KEY`)
	requireNotContains(t, js, `data-etf-rule-level`)
	requireNotContains(t, js, `data-etf-rule-done`)
	requireNotContains(t, js, `updateEtfRuleTrackerEntry`)
	requireNotContains(t, js, `已手动覆盖`)
	requireNotContains(t, js, `周计划`)
	requireNotContains(t, js, `本周计划`)
	requireContains(t, js, `renderEtfRuleLiveStatus`)
	requireContains(t, css, `.etf-rule-panel`)
	requireContains(t, css, `.etf-rule-card`)
	requireContains(t, css, `.etf-rule-card-actions`)
	requireContains(t, css, `.etf-rule-buy-button`)
	requireContains(t, css, `.etf-buy-target`)
	requireContains(t, css, `.etf-pool-progress`)
	requireContains(t, css, `.etf-allocation-progress`)
	requireContains(t, css, `.etf-allocation-bar`)
	requireContains(t, css, `.etf-rule-action-list`)
	requireContains(t, css, `.etf-rule-action-item`)
	requireContains(t, css, `.etf-rule-detail-layer`)
	requireContains(t, css, `.etf-rule-live`)
	requireContains(t, css, `.etf-rule-metric`)
	requireContains(t, css, `.etf-rule-rulebook`)
	requireContains(t, css, `.etf-rule-condition.boost`)
	requireNotContains(t, js, `etf-rule-levels`)
	requireNotContains(t, js, `etf-rule-active-condition`)
	requireNotContains(t, js, `etf-rule-source`)
	requireNotContains(t, css, `.etf-rule-levels`)
	requireNotContains(t, css, `.etf-rule-active-condition`)
	requireNotContains(t, css, `.etf-rule-source`)
	requireNotContains(t, css, `etf-rule-confidence`)
}

func TestTradeAndFundTablesUseReadableNumbers(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")
	css := readTextFile(t, "styles.css")
	fundPageSection := extractBetween(t, html, `data-page="funds"`, `data-page="trades"`)

	requireNotContains(t, fundPageSection, `<th>类型</th>`)
	requireContains(t, fundPageSection, `<th>份额</th>`)
	requireContains(t, js, `function tradeQuantityText`)
	requireContains(t, js, `minimumFractionDigits: 4`)
	requireContains(t, js, `tradeQuantityText(trade.shares, isFundTrade)`)
	requireContains(t, js, `${escapeHTML(fund.symbol)} · ${escapeHTML(fundTypeLabel(fund))}`)
	requireContains(t, css, `.fund-table`)
	requireContains(t, css, `min-width: 880px`)
}

func TestTopbarDoesNotRenderRedundantDecisionShortcut(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireNotContains(t, html, `topbar-decision-link`)
	requireNotContains(t, html, `记录决策`)
}

func TestOverviewAutoRefreshesStockQuotesOnly(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")
	goMain := readTextFile(t, "main.go")

	requireNotContains(t, html, `id="updateQuotesButton"`)
	requireNotContains(t, html, `更新行情/净值`)
	requireNotContains(t, js, `updateQuotesButton`)
	requireContains(t, js, `async function updateQuotes()`)
	requireContains(t, js, `requestJSON("/api/quotes/update", { method: "POST" })`)
	requireContains(t, js, `async function autoRefreshOverviewStockQuotes()`)
	requireContains(t, js, `requestJSON("/api/quotes/stocks/update", { method: "POST", timeoutMs: 60000 })`)
	requireContains(t, js, `const STOCK_QUOTE_AUTO_REFRESH_MINUTES = 1`)
	requireContains(t, js, `if (route.page === "overview")`)
	requireContains(t, js, `void autoRefreshOverviewStockQuotes();`)
	requireContains(t, goMain, `mux.HandleFunc("POST /api/quotes/stocks/update", server.handleUpdateStockQuotes)`)
	requireContains(t, js, `if (window.location.hash.slice(1) === view)`)
	requireContains(t, js, `handleRoute(view);`)
}

func TestCloudSyncButtonAndEndpointAreWired(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")
	goMain := readTextFile(t, "main.go")
	css := readTextFile(t, "styles.css")
	topbarActions := extractBetween(t, html, `<div class="topbar-actions">`, `<section class="page active"`)

	requireContains(t, html, `styles.css?v=overview-layer-v1`)
	requireContains(t, html, `app.js?v=overview-layer-v1`)
	requireContains(t, html, `id="cloudSyncButton"`)
	requireContains(t, topbarActions, `id="cloudSyncButton"`)
	requireContains(t, topbarActions, `id="privacyToggle"`)
	requireContains(t, topbarActions, `overview-only-action`)
	requireContains(t, topbarActions, `portfolio.json" hidden>`)
	requireContains(t, topbarActions, `title="隐藏持仓数据" hidden>`)
	requireContains(t, css, `body:not([data-active-page="overview"]) .overview-only-action`)
	requireContains(t, css, `display: none !important;`)
	requireContains(t, js, `cloudSyncButton: document.querySelector("#cloudSyncButton")`)
	requireContains(t, js, `function setTopbarActionsVisibility(page)`)
	requireContains(t, js, `const showOnOverview = page === "overview"`)
	requireContains(t, js, `button.hidden = !showOnOverview`)
	requireContains(t, js, `setTopbarActionsVisibility(nextView)`)
	requireContains(t, js, `async function syncCloudPortfolio()`)
	requireContains(t, js, `requestJSON("/api/cloud/sync", { method: "POST", timeoutMs: 90000 })`)
	requireContains(t, js, `setLoadedState(result.state)`)
	requireContains(t, goMain, `mux.HandleFunc("POST /api/cloud/sync", server.handleSyncCloudPortfolio)`)
	requireContains(t, goMain, `func (s *Server) handleSyncCloudPortfolio`)
	requireContains(t, goMain, `backupPortfolioFile()`)
	requireContains(t, goMain, `saveState(remoteState)`)
}

func TestOverviewDailyPnlStartsFromRecordedHistory(t *testing.T) {
	js := readTextFile(t, "app.js")
	css := readTextFile(t, "styles.css")

	requireContains(t, js, `...(Array.isArray(state.pnlHistory) ? state.pnlHistory : [])`)
	requireContains(t, js, `function scrollOverviewPnlToLatest()`)
	requireContains(t, js, `scroller.scrollLeft = scroller.scrollWidth`)
	requireContains(t, js, `scrollOverviewPnlToLatest();`)
	requireContains(t, css, `grid-template-columns: repeat(var(--bar-count), minmax(56px, 1fr))`)
	requireContains(t, css, `font-variant-numeric: tabular-nums`)

	fallback := extractBetween(t, js, `function overviewFallbackPnlValue(periodKey, range, anchorDate, positions, fundPositions, stats) {`, `function overviewPnlSeries`)
	requireContains(t, fallback, `if (range === "day")`)
	requireContains(t, fallback, `return 0;`)
	requireNotContains(t, js, `function overviewDailyFallbackPnlValue`)
	requireNotContains(t, fallback, `currentPriceDate || item?.previousCloseDate`)
	requireNotContains(t, fallback, `periodKey === anchorDay`)
	requireNotContains(t, fallback, `overviewPnlValue(positions, "day") + overviewPnlValue(fundPositions, "day")`)
}

func TestRedundantMaintenanceEntrancesAreRemoved(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")
	css := readTextFile(t, "styles.css")

	requireNotContains(t, html, `id="openResearchPanel"`)
	requireNotContains(t, html, `id="exportChatGPTContext"`)
	requireNotContains(t, html, `id="clearDecisionLogs"`)
	requireNotContains(t, html, `清理日志`)
	requireNotContains(t, js, `mobileActionButton("stockTrade"`)
	requireNotContains(t, js, `action === "stockTrade"`)
	requireNotContains(t, js, `mobileContextActionbar`)
	requireNotContains(t, js, `mobilePortfolioControlsDialog`)
	requireNotContains(t, js, `POSITION_MOBILE_SORT_OPTIONS`)
	requireNotContains(t, js, `SUNNY30_MOBILE_SORT_OPTIONS`)
	requireNotContains(t, js, `POSITION_MOBILE_FILTERS`)
	requireNotContains(t, js, `FUND_MOBILE_SORT_OPTIONS`)
	requireNotContains(t, js, `data-mobile-position-filter`)
	requireNotContains(t, css, `.mobile-context-actionbar`)
	requireNotContains(t, css, `.mobile-control-option`)
}

func TestAnalysisStockTrackerMarkupIsRemoved(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireNotContains(t, html, `<select id="sunny30MobileSort">`)
	requireNotContains(t, html, `data-sunny30-sort=`)
	requireNotContains(t, html, `data-sunny30-sort="reason"`)
	requireNotContains(t, html, `排序原因`)
	requireNotContains(t, html, `data-sunny30-sort="gate"`)
}

func TestHoldingsDecisionTableSeparatesTodayAndTotalPnl(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")
	holdingsTable := extractBetween(t, html, `<table class="decision-table">`, `</table>`)
	renderPositions := extractBetween(t, js, `function renderPositions(positions) {`, `function fundTypeLabel(fund) {`)

	requireNotContains(t, html, `id="positionCategorySummary"`)
	requireContains(t, holdingsTable, `data-position-sort="shares"`)
	requireContains(t, holdingsTable, `<span>股数</span>`)
	requireContains(t, holdingsTable, `data-position-sort="marketValue"`)
	requireContains(t, holdingsTable, `<span>市值</span>`)
	requireContains(t, holdingsTable, `data-position-sort="currentPrice"`)
	requireContains(t, holdingsTable, `<span>现价</span>`)
	requireContains(t, holdingsTable, `data-position-sort="dayChange"`)
	requireContains(t, holdingsTable, `<span>今日盈亏</span>`)
	requireContains(t, holdingsTable, `data-position-sort="twentyDayChange"`)
	requireContains(t, holdingsTable, `<span>20日涨跌幅</span>`)
	requireContains(t, holdingsTable, `<span>总盈亏</span>`)
	requireContains(t, renderPositions, `data-label="股数"`)
	requireContains(t, renderPositions, `data-label="市值"`)
	requireContains(t, renderPositions, `data-label="现价"`)
	requireContains(t, renderPositions, `data-label="今日盈亏"`)
	requireContains(t, renderPositions, `data-label="20日涨跌幅"`)
	requireContains(t, renderPositions, `data-label="总盈亏"`)
	requireContains(t, js, `renderMobileStat("总盈亏"`)
	requireContains(t, js, `renderMobileStat("20日涨跌幅"`)
	requireNotContains(t, holdingsTable, `<span>分类</span>`)
	requireNotContains(t, holdingsTable, `<span>安全边际</span>`)
	requireNotContains(t, holdingsTable, `<span>长期评分</span>`)
	requireNotContains(t, holdingsTable, `<th>健康状态</th>`)
	requireNotContains(t, renderPositions, `data-label="分类"`)
	requireNotContains(t, renderPositions, `data-label="安全边际"`)
	requireNotContains(t, renderPositions, `data-label="长期评分"`)
	requireNotContains(t, renderPositions, `data-label="健康状态"`)
	requireNotContains(t, renderPositions, `data-label="市值/现价"`)
	requireNotContains(t, renderPositions, `data-label="盈亏"`)
	requireNotContains(t, js, `renderMobileStat("累计盈亏"`)
}

func TestStockDetailOwnsEditableHumanInputs(t *testing.T) {
	js := readTextFile(t, "app.js")

	requireContains(t, js, `data-stock-human-input-form`)
	requireContains(t, js, `name="buyLogic"`)
	requireContains(t, js, `name="valuationConfidence"`)
	requireContains(t, js, `name="killCriteria"`)
	requireContains(t, js, `data-stock-valuation-form`)
	requireContains(t, js, `{ key: "base", label: "基准" }`)
	requireContains(t, js, `.revenueGrowth`)
	requireContains(t, js, `saveStockHumanInputs`)
	requireContains(t, js, `saveStockValuationInputs`)
}

func TestTradeListHasTrashDeleteAction(t *testing.T) {
	js := readTextFile(t, "app.js")
	css := readTextFile(t, "styles.css")

	requireContains(t, js, `data-delete-trade`)
	requireContains(t, js, `aria-label="删除交易记录`)
	requireContains(t, js, `deleteTradeRecord`)
	requireContains(t, css, `.trade-delete-button`)
}

func TestFundNavUsesFourDecimals(t *testing.T) {
	js := readTextFile(t, "app.js")

	requireContains(t, js, `function fundNav`)
	requireContains(t, js, `function isExchangeFundCode`)
	requireContains(t, js, `isExchangeFundCode(normalized) ? "etf" : "otc"`)
	requireContains(t, js, `minimumFractionDigits: 4`)
	requireContains(t, js, `maximumFractionDigits: 4`)

	renderFunds := extractBetween(t, js, `function renderFunds(fundPositions = computeFundPositions()) {`, `function renderAllocation(positions) {`)
	requireContains(t, renderFunds, `data-label="成本净值">${escapeHTML(privateText(fundNav(fund.cost, fund.currency)))}`)
	requireContains(t, renderFunds, `data-label="最新净值">${escapeHTML(privateText(fundNav(fund.currentPrice, fund.currency)))}`)
	requireContains(t, renderFunds, `renderMobileStat("成本净值", privateText(fundNav(fund.cost, fund.currency)))`)
	requireContains(t, renderFunds, `renderMobileStat("最新净值", privateText(fundNav(fund.currentPrice, fund.currency)))`)
	requireNotContains(t, renderFunds, `currency(fund.cost, fund.currency)`)
	requireNotContains(t, renderFunds, `currency(fund.currentPrice, fund.currency)`)

	renderTrades := extractBetween(t, js, `function renderTrades() {`, `function parseValuationRangeText(text) {`)
	requireContains(t, renderTrades, `const isFundTrade = normalizeAssetType(trade.assetType) === "fund";`)
	requireContains(t, renderTrades, `const tradePriceText = isFundTrade ? fundNav(trade.price, trade.currency) : currency(trade.price, trade.currency);`)
	requireContains(t, renderTrades, `const currentPriceText = isFundTrade ? fundNav(trade.currentPrice, trade.currency) : currency(trade.currentPrice, trade.currency);`)
	requireContains(t, renderTrades, `const currentLabel = isFundTrade ? "最新净值" : "最新价";`)
}

func TestStockTradeCreatesDirectHoldingWithoutTrackingPool(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")

	requireContains(t, html, `股票代码/名称`)
	requireContains(t, html, `placeholder="9926.HK 康方生物"`)
	requireContains(t, js, `function parseStockTradeInput`)
	requireContains(t, js, `请输入股票代码，名称可以跟在代码后面，例如 9926.HK 康方生物`)
	requireNotContains(t, js, `请先把它加入持仓或晴仓30`)

	tradeNames := extractBetween(t, js, `function renderTradeStockNames() {`, `function routeInfo`)
	requireContains(t, tradeNames, `optionType: "持仓"`)
	requireNotContains(t, tradeNames, `state.candidates`)
	requireNotContains(t, tradeNames, `optionType: "跟踪"`)

	tradeForm := extractBetween(t, js, `function tradeFromSimpleForm(formData) {`, `async function addTrade(formData)`)
	requireContains(t, tradeForm, `parseStockTradeInput(nameInput)`)
	requireNotContains(t, tradeForm, `state.candidates`)
	requireNotContains(t, tradeForm, `晴仓30`)

	addTrade := extractBetween(t, js, `async function addTrade(formData) {`, `async function deleteTradeRecord`)
	requireNotContains(t, addTrade, `removeCandidate(symbol)`)
	requireNotContains(t, addTrade, `upsertCandidate(clearedCandidateFromHolding`)
	requireNotContains(t, addTrade, `holdingFromCandidate`)
}

func extractBetween(t *testing.T, content, start, end string) string {
	t.Helper()
	startIndex := strings.Index(content, start)
	if startIndex < 0 {
		t.Fatalf("missing start marker %q", start)
	}
	rest := content[startIndex+len(start):]
	endIndex := strings.Index(rest, end)
	if endIndex < 0 {
		t.Fatalf("missing end marker %q", end)
	}
	return rest[:endIndex]
}
