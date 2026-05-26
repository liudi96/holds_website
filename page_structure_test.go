package main

import (
	"strings"
	"testing"
)

func TestHoldingsAndScreenerAreSeparatePages(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireContains(t, html, `data-page="holdings"`)
	requireContains(t, html, `data-page="screener"`)
	requireNotContains(t, html, `data-page="portfolio"`)
	requireNotContains(t, html, `data-page="valuation"`)

	holdingsStart := strings.Index(html, `data-page="holdings"`)
	screenerStart := strings.Index(html, `data-page="screener"`)
	mastersStart := strings.Index(html, `id="masters"`)
	sunny30Start := strings.Index(html, `id="sunny30Section"`)
	valuationStart := strings.Index(html, `id="valuationModuleList"`)
	if holdingsStart < 0 || screenerStart < 0 || mastersStart < 0 || sunny30Start < 0 || valuationStart < 0 {
		t.Fatalf("missing required page anchors")
	}
	if !(holdingsStart < mastersStart && mastersStart < screenerStart) {
		t.Fatalf("expected masters inside holdings page before screener page")
	}
	if !(screenerStart < sunny30Start && sunny30Start < valuationStart) {
		t.Fatalf("expected screener page to contain stock screener before valuation module")
	}
}

func TestScreenerOwnsValuationRoute(t *testing.T) {
	js := readTextFile(t, "app.js")

	requireContains(t, js, `return { view: "screener", page: "screener" }`)
	requireContains(t, js, `if (route.page === "screener")`)
	requireContains(t, js, `renderSunny30(positions);`)
	requireContains(t, js, `renderValuationModule(positions);`)
	requireNotContains(t, js, `route.page === "valuation"`)
	requireNotContains(t, js, `route.page === "portfolio"`)
}

func TestTopbarDoesNotRenderRedundantDecisionShortcut(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireNotContains(t, html, `topbar-decision-link`)
	requireNotContains(t, html, `记录决策`)
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

func TestScreenerRemovesLowSignalSortControls(t *testing.T) {
	html := readTextFile(t, "index.html")
	sunnySort := extractBetween(t, html, `<select id="sunny30MobileSort">`, `</select>`)

	requireNotContains(t, sunnySort, `value="type:asc"`)
	requireNotContains(t, sunnySort, `value="type:desc"`)
	requireNotContains(t, sunnySort, `value="return:desc"`)
	requireNotContains(t, sunnySort, `value="return:asc"`)
	requireNotContains(t, sunnySort, `value="moat:desc"`)
	requireNotContains(t, sunnySort, `value="moat:asc"`)
	requireNotContains(t, html, `data-sunny30-sort="reason"`)
}

func TestScreenerDoesNotRepeatSubscoresAsReasonColumn(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")

	requireNotContains(t, html, `排序原因`)
	requireNotContains(t, js, `data-label="排序原因"`)
	requireNotContains(t, js, `key === "reason"`)
	requireNotContains(t, js, `screening.reason`)
}

func TestScreenerMarksRejectedStocksInsteadOfHardRejectColumn(t *testing.T) {
	html := readTextFile(t, "index.html")
	js := readTextFile(t, "app.js")
	css := readTextFile(t, "styles.css")

	requireNotContains(t, html, `data-sunny30-sort="gate"`)
	requireNotContains(t, js, `data-label="硬否决"`)
	requireNotContains(t, js, `screeningGateCell`)
	requireContains(t, js, `screening-reject-name`)
	requireContains(t, css, `.screening-reject-name strong`)
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
