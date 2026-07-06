package main

import (
	"os"
	"strings"
	"testing"
)

func readTextFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func requireContains(t *testing.T, content, needle string) {
	t.Helper()
	if !strings.Contains(content, needle) {
		t.Fatalf("expected content to contain %q", needle)
	}
}

func requireNotContains(t *testing.T, content, needle string) {
	t.Helper()
	if strings.Contains(content, needle) {
		t.Fatalf("expected content not to contain %q", needle)
	}
}

func TestOverviewAssetDashboardMarkup(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireContains(t, html, `class="overview-dashboard"`)
	requireContains(t, html, `class="panel overview-assets-panel"`)
	requireContains(t, html, `class="panel overview-pnl-panel"`)
	requireContains(t, html, `class="panel overview-allocation-panel"`)
	requireContains(t, html, `权益总资产`)
	requireContains(t, html, `股票 + 基金`)
	requireContains(t, html, `id="stockMarketValue"`)
	requireContains(t, html, `id="fundMarketValue"`)
	requireContains(t, html, `id="overviewPnlBars"`)
	requireContains(t, html, `data-overview-pnl-range="day"`)
	requireContains(t, html, `data-overview-pnl-range="month"`)
	requireContains(t, html, `data-overview-pnl-range="year"`)
	requireContains(t, html, `id="assetAllocationPie"`)
	requireContains(t, html, `id="assetAllocationLegend"`)
	requireNotContains(t, html, `id="cashValueMetric"`)
	requireNotContains(t, html, `股票 + 基金 + 现金`)
	requireNotContains(t, html, `股票 / 基金 / 现金`)
	requireNotContains(t, html, `class="decision-cockpit-title"`)
	requireNotContains(t, html, `class="decision-cockpit-grid"`)
	requireNotContains(t, html, `id="committeeConsensus"`)
	requireNotContains(t, html, `id="overviewBuyCandidates"`)
	requireNotContains(t, html, `id="overviewRiskReview"`)
	requireNotContains(t, html, `id="actionConclusion"`)
}

func TestOverviewRuntimeStatusSinksStayHidden(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireContains(t, html, `id="quoteUpdateStatus" hidden`)
	requireContains(t, html, `id="dataQualityMetric" hidden`)
	requireContains(t, html, `id="dataQualityDetail" hidden`)
}

func TestOverviewAssetDashboardStyleHooks(t *testing.T) {
	css := readTextFile(t, "styles.css")

	requireContains(t, css, `.overview-dashboard`)
	requireContains(t, css, `.overview-asset-cards`)
	requireContains(t, css, `.overview-pnl-bars`)
	requireContains(t, css, `.overview-pnl-chart`)
	requireContains(t, css, `.overview-pnl-column`)
	requireContains(t, css, `.overview-range-tabs`)
	requireContains(t, css, `.overview-allocation-pie`)
	requireContains(t, css, `.overview-allocation-body`)
}

func TestOverviewAssetDashboardRenderingHooks(t *testing.T) {
	js := readTextFile(t, "app.js")

	requireContains(t, js, `assetAllocationPie: document.querySelector("#assetAllocationPie")`)
	requireContains(t, js, `overviewPnlBars: document.querySelector("#overviewPnlBars")`)
	requireContains(t, js, `overviewPnlRangeTabs: document.querySelector("#overviewPnlRangeTabs")`)
	requireContains(t, js, `stockMarketValue: document.querySelector("#stockMarketValue")`)
	requireContains(t, js, `fundMarketValue: document.querySelector("#fundMarketValue")`)
	requireNotContains(t, js, `cashValueMetric: document.querySelector("#cashValueMetric")`)
	requireContains(t, js, `let overviewPnlRange = "day"`)
	requireContains(t, js, `function overviewPnlSeries`)
	requireContains(t, js, `function renderOverviewPnlChart`)
	requireContains(t, js, `renderOverviewPnlChart(positions, fundPositions);`)
	requireNotContains(t, js, `renderDecisionArea(positions);`)
	requireNotContains(t, js, `renderCommitteeOverview(positions);`)
}

func TestOverviewDecisionQueueDoesNotRenderFauxActionButtons(t *testing.T) {
	js := readTextFile(t, "app.js")
	css := readTextFile(t, "styles.css")

	requireNotContains(t, js, `decision-action-controls`)
	requireNotContains(t, js, `decision-chip`)
	requireNotContains(t, css, `.decision-action-controls`)
	requireNotContains(t, css, `.decision-chip`)
}

func TestOverviewDashboardHasNoDecisionCockpitHeader(t *testing.T) {
	html := readTextFile(t, "index.html")
	css := readTextFile(t, "styles.css")

	requireNotContains(t, html, `筛选股票池`)
	requireNotContains(t, html, `批量复核估值`)
	requireNotContains(t, html, `decision-cockpit-actions`)
	requireNotContains(t, css, `.decision-cockpit-actions`)
}
