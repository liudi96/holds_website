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

func TestOverviewCockpitV2Markup(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireContains(t, html, `class="panel overview-cockpit"`)
	requireContains(t, html, `class="overview-conclusion`)
	requireContains(t, html, `class="metrics-grid overview-metrics-strip"`)
	requireContains(t, html, `class="overview-exposure"`)
	requireContains(t, html, `id="assetAllocationBar"`)
	requireContains(t, html, `class="decision-cockpit-title"`)
	requireContains(t, html, `class="decision-cockpit-grid"`)
	requireContains(t, html, `id="overviewBuyCandidates"`)
	requireContains(t, html, `id="overviewRiskReview"`)
	requireNotContains(t, html, `allocation-donut`)
}

func TestOverviewRuntimeStatusSinksStayHidden(t *testing.T) {
	html := readTextFile(t, "index.html")

	requireContains(t, html, `id="quoteUpdateStatus" hidden`)
	requireContains(t, html, `id="dataQualityMetric" hidden`)
	requireContains(t, html, `id="dataQualityDetail" hidden`)
}

func TestOverviewCockpitV2StyleHooks(t *testing.T) {
	css := readTextFile(t, "styles.css")

	requireContains(t, css, `.overview-cockpit`)
	requireContains(t, css, `.overview-conclusion`)
	requireContains(t, css, `.overview-metrics-strip`)
	requireContains(t, css, `.overview-exposure`)
	requireContains(t, css, `.exposure-bar`)
	requireContains(t, css, `.decision-cockpit-title`)
	requireContains(t, css, `.decision-cockpit-grid`)
	requireContains(t, css, `.cockpit-signal`)
}

func TestOverviewCockpitV2RenderingHooks(t *testing.T) {
	js := readTextFile(t, "app.js")

	requireContains(t, js, `assetAllocationBar: document.querySelector("#assetAllocationBar")`)
	requireContains(t, js, `overviewBuyCandidates: document.querySelector("#overviewBuyCandidates")`)
	requireContains(t, js, `overviewRiskReview: document.querySelector("#overviewRiskReview")`)
	requireContains(t, js, `function renderOverviewBuyCandidates`)
	requireContains(t, js, `function renderOverviewRiskReview`)
}

func TestOverviewDecisionQueueDoesNotRenderFauxActionButtons(t *testing.T) {
	js := readTextFile(t, "app.js")
	css := readTextFile(t, "styles.css")

	requireNotContains(t, js, `decision-action-controls`)
	requireNotContains(t, js, `decision-chip`)
	requireNotContains(t, css, `.decision-action-controls`)
	requireNotContains(t, css, `.decision-chip`)
}

func TestOverviewDecisionCockpitHasNoRedundantHeaderButtons(t *testing.T) {
	html := readTextFile(t, "index.html")
	css := readTextFile(t, "styles.css")

	requireNotContains(t, html, `筛选股票池`)
	requireNotContains(t, html, `批量复核估值`)
	requireNotContains(t, html, `decision-cockpit-actions`)
	requireNotContains(t, css, `.decision-cockpit-actions`)
}
