package main

import (
	"math"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func almostEqual(got float64, want float64, tolerance float64) bool {
	return math.Abs(got-want) <= tolerance
}

func TestDrawdownFromRecentHigh(t *testing.T) {
	drawdown, date, err := drawdownFromRecentHigh([]dailyClose{
		{Date: "2026-07-01", Price: 100},
		{Date: "2026-07-02", Price: 120},
		{Date: "2026-07-03", Price: 90},
	}, 252)
	if err != nil {
		t.Fatalf("drawdownFromRecentHigh returned error: %v", err)
	}
	if date != "2026-07-03" {
		t.Fatalf("date = %q, want 2026-07-03", date)
	}
	if !almostEqual(drawdown, 0.25, 0.000001) {
		t.Fatalf("drawdown = %.6f, want 0.25", drawdown)
	}
}

func TestParseMultplTable(t *testing.T) {
	rows, err := parseMultplTable([]byte(`
		<table>
			<tr><td>Jul 1, 2026</td><td>39.12</td></tr>
			<tr><td>Jun 1, 2026</td><td>&#x2002; 37.50</td></tr>
		</table>
	`))
	if err != nil {
		t.Fatalf("parseMultplTable returned error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows len = %d, want 2", len(rows))
	}
	if rows[0].Date != "2026-07-01" || rows[0].Price != 39.12 {
		t.Fatalf("first row = %+v", rows[0])
	}
}

func TestCAPEPercentileFromMonthlyValuesUsesTenYearWindow(t *testing.T) {
	rows := []dailyClose{
		{Date: "2014-06-01", Price: 10},
		{Date: "2016-07-01", Price: 20},
		{Date: "2017-07-01", Price: 30},
		{Date: "2018-07-01", Price: 40},
		{Date: "2019-07-01", Price: 50},
		{Date: "2020-07-01", Price: 60},
		{Date: "2021-07-01", Price: 70},
		{Date: "2022-07-01", Price: 80},
		{Date: "2023-07-01", Price: 90},
		{Date: "2024-07-01", Price: 100},
		{Date: "2025-07-01", Price: 110},
		{Date: "2026-07-01", Price: 120},
	}
	percentile, date, err := capePercentileFromMonthlyValues(rows, 10)
	if err != nil {
		t.Fatalf("capePercentileFromMonthlyValues returned error: %v", err)
	}
	if date != "2026-07-01" {
		t.Fatalf("date = %q, want 2026-07-01", date)
	}
	if !almostEqual(percentile, 1, 0.000001) {
		t.Fatalf("percentile = %.6f, want 1", percentile)
	}
}

func TestPercentileFromHistoryOfMarketPointsUsesTenYearWindow(t *testing.T) {
	percentile, date, err := percentileFromHistoryOfMarketPoints([]historyOfMarketPoint{
		{Date: "2015-05-18", Value: 5},
		{Date: "2016-05-18", Value: 20},
		{Date: "2020-05-18", Value: 10},
		{Date: "2024-05-18", Value: 30},
		{Date: "2026-05-18", Value: 25},
	}, 10, "forward PE")
	if err != nil {
		t.Fatalf("percentileFromHistoryOfMarketPoints returned error: %v", err)
	}
	if date != "2026-05-18" {
		t.Fatalf("date = %q, want 2026-05-18", date)
	}
	if !almostEqual(percentile, 0.75, 0.000001) {
		t.Fatalf("percentile = %.6f, want 0.75", percentile)
	}
}

func TestHistoryOfMarketPointsWithCurrentForwardUsesUpdatedDate(t *testing.T) {
	points := historyOfMarketPointsWithCurrentForward([]historyOfMarketPoint{
		{Date: "2016-07-07", Value: 10},
		{Date: "2018-07-07", Value: 12},
		{Date: "2020-07-07", Value: 20},
		{Date: "2022-07-07", Value: 25},
		{Date: "2026-05-18", Value: 28},
	}, "2026-07-07", historyOfMarketCurrentValuation{Forward: 30})

	percentile, date, err := percentileFromHistoryOfMarketPoints(points, 10, "Nasdaq 100 forward PE")
	if err != nil {
		t.Fatalf("percentileFromHistoryOfMarketPoints returned error: %v", err)
	}
	if date != "2026-07-07" {
		t.Fatalf("date = %q, want 2026-07-07", date)
	}
	if !almostEqual(percentile, 1, 0.000001) {
		t.Fatalf("percentile = %.6f, want 1", percentile)
	}
}

func TestValuationDateStale(t *testing.T) {
	now := time.Date(2026, 7, 8, 15, 30, 0, 0, time.UTC)
	tests := []struct {
		name string
		date string
		want bool
	}{
		{name: "old primary", date: "2026-07-02", want: true},
		{name: "recent primary", date: "2026-07-07", want: false},
		{name: "same day", date: "2026-07-08", want: false},
		{name: "future date", date: "2026-07-09", want: false},
		{name: "bad date", date: "bad-date", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valuationDateStale(tt.date, now, primaryValuationMaxLagDays); got != tt.want {
				t.Fatalf("valuationDateStale(%q) = %v, want %v", tt.date, got, tt.want)
			}
		})
	}
}

func TestParseStooqDailyCSV(t *testing.T) {
	rows, err := parseStooqDailyCSV([]byte("Date,Open,High,Low,Close,Volume\n2026-07-02,1,2,1,100.5,0\n2026-07-03,1,2,1,101.25,0\n"))
	if err != nil {
		t.Fatalf("parseStooqDailyCSV returned error: %v", err)
	}
	if len(rows) != 2 || rows[1].Date != "2026-07-03" || rows[1].Price != 101.25 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestParseNasdaqHistoricalCloses(t *testing.T) {
	rows, err := parseNasdaqHistoricalCloses([]byte(`{
		"data": {
			"tradesTable": {
				"rows": [
					{"date": "07/03/2026", "close": "29,400.25"},
					{"date": "07/02/2026", "close": "29,329.21"}
				]
			}
		}
	}`))
	if err != nil {
		t.Fatalf("parseNasdaqHistoricalCloses returned error: %v", err)
	}
	if len(rows) != 2 || rows[0].Date != "2026-07-02" || rows[0].Price != 29329.21 || rows[1].Date != "2026-07-03" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestParseWorldPERatioNasdaq100(t *testing.T) {
	snapshot, err := parseWorldPERatioNasdaq100([]byte(`
		<p>The estimated <b>Price-to-Earnings (P/E) Ratio</b> for <b>Nasdaq 100 Index</b> is <b>32.74</b>, calculated on <b>02 July 2026</b>.</p>
		<table>
			<tr class="w3-center row1">
				<td><b>Last 1Y</b></td>
				<td>33.37</td>
				<td>0.63</td>
				<td>[32.11 · <font>32.74 , 34.00</font> · 34.63]</td>
				<td><div class="pe-range-container"></div></td>
				<td>-1.01 &sigma;</td>
				<td>Undervalued</td>
			</tr>
			<tr class="w3-center row10">
				<td><b>Last 10Y</b></td>
				<td>27.21</td>
				<td>4.06</td>
				<td>[19.08 · <font>23.14 , 31.27</font> · 35.33]</td>
				<td><div class="pe-range-container"></div></td>
				<td>+1.36 &sigma;</td>
				<td>Overvalued</td>
			</tr>
		</table>
	`))
	if err != nil {
		t.Fatalf("parseWorldPERatioNasdaq100 returned error: %v", err)
	}
	if snapshot.Date != "2026-07-02" || snapshot.CurrentPE != 32.74 {
		t.Fatalf("snapshot date/current = %+v", snapshot)
	}
	if snapshot.Average10Y != 27.21 || snapshot.StdDev10Y != 4.06 || snapshot.ZScore != 1.36 {
		t.Fatalf("snapshot Last 10Y metrics = %+v", snapshot)
	}
	if !almostEqual(snapshot.Percentile, 0.913085, 0.00001) {
		t.Fatalf("percentile = %.6f, want about 0.913085", snapshot.Percentile)
	}
}

func TestParseLeguleguIndexPERowsAndA500Percentile(t *testing.T) {
	rows, err := parseLeguleguIndexPERows([]byte(`{
		"data": [
			{"date":"2026-07-03","ttmPe":47.14,"ttmPeQuantile":0.96037},
			{"date":"2024-09-24","ttmPe":36.10,"ttmPeQuantile":0.24123},
			{"date":"2025-01-02","ttmPe":42.00,"ttmPeQuantile":0.61234},
			{"date":"2026-01-02","ttmPe":44.00,"ttmPeQuantile":0.72345}
		]
	}`))
	if err != nil {
		t.Fatalf("parseLeguleguIndexPERows returned error: %v", err)
	}
	if len(rows) != 4 || rows[0].Date != "2024-09-24" || rows[3].Date != "2026-07-03" {
		t.Fatalf("rows should be sorted by date, got %+v", rows)
	}
	percentile, date, err := a500PEPercentileFromRows(rows)
	if err != nil {
		t.Fatalf("a500PEPercentileFromRows returned error: %v", err)
	}
	if date != "2026-07-03" {
		t.Fatalf("date = %q, want 2026-07-03", date)
	}
	if !almostEqual(percentile, 0.96037, 0.000001) {
		t.Fatalf("percentile = %.6f, want 0.96037", percentile)
	}
}

func TestParseEastmoneyFundDividendsAndTrailingAmount(t *testing.T) {
	events, err := parseEastmoneyFundDividends([]byte(`
		<table>
			<tr><td>2025年</td><td>2025-12-12</td><td>2025-12-15</td><td>每份派现金0.0300元</td><td>2025-12-18</td></tr>
			<tr><td>2025年</td><td>2025-07-14</td><td>2025-07-15</td><td>每份派现金0.0450元</td><td>2025-07-18</td></tr>
			<tr><td>2024年</td><td>2024-12-16</td><td>2024-12-17</td><td>每份派现金0.0550元</td><td>2024-12-20</td></tr>
		</table>
	`))
	if err != nil {
		t.Fatalf("parseEastmoneyFundDividends returned error: %v", err)
	}
	if len(events) != 3 || events[0].Date != "2025-12-15" || events[0].Amount != 0.03 {
		t.Fatalf("events = %+v", events)
	}
	amount, err := trailingFundDividendAmount(events, "2026-07-03")
	if err != nil {
		t.Fatalf("trailingFundDividendAmount returned error: %v", err)
	}
	if !almostEqual(amount, 0.075, 0.000001) {
		t.Fatalf("amount = %.6f, want 0.075", amount)
	}
}

func TestEastmoneyGlobalIndexSecID(t *testing.T) {
	if got := eastmoneyGlobalIndexSecID("^GSPC"); got != "100.SPX" {
		t.Fatalf("S&P secid = %q", got)
	}
	if got := eastmoneyGlobalIndexSecID("^NDX"); got != "100.NDX" {
		t.Fatalf("Nasdaq 100 secid = %q", got)
	}
}

func TestEvaluateSP500RuleWithCAPEAndDrawdown(t *testing.T) {
	cape := 0.82
	drawdown := 0.12
	got := evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &cape, Drawdown: &drawdown})
	if got.Level != "half" || !got.Complete {
		t.Fatalf("high CAPE should choose half by base rule, got %+v", got)
	}

	cape = 0.50
	drawdown = 0.03
	got = evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &cape, Drawdown: &drawdown})
	if got.Level != "half" || !got.Complete {
		t.Fatalf("normal CAPE near high should downshift to half, got %+v", got)
	}

	cape = 0.50
	drawdown = 0.12
	got = evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &cape, Drawdown: &drawdown})
	if got.Level != "one" || !got.Complete {
		t.Fatalf("normal CAPE/drawdown should choose one, got %+v", got)
	}

	cape = 0.10
	drawdown = 0.30
	got = evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &cape, Drawdown: &drawdown})
	if got.Level != "two" || !got.Complete {
		t.Fatalf("low CAPE/deep drawdown should choose two, got %+v", got)
	}
}

func TestEvaluateNasdaq100RuleUsesPEPercentileThenDrawdownLimit(t *testing.T) {
	pePercentile := 0.80
	drawdown := 0.046
	got := evaluateNasdaq100Rule(etfRuleInputs{ValuationPercentile: &pePercentile, Drawdown: &drawdown})
	if got.Level != "quarter" || !got.Complete {
		t.Fatalf("high PE percentile with shallow drawdown should downshift to quarter, got %+v", got)
	}

	pePercentile = 0.80
	drawdown = 0.10
	got = evaluateNasdaq100Rule(etfRuleInputs{ValuationPercentile: &pePercentile, Drawdown: &drawdown})
	if got.Level != "half" || !got.Complete {
		t.Fatalf("high PE percentile without shallow drawdown should keep half, got %+v", got)
	}

	pePercentile = 0.10
	drawdown = 0.10
	got = evaluateNasdaq100Rule(etfRuleInputs{ValuationPercentile: &pePercentile, Drawdown: &drawdown})
	if got.Level != "oneHalf" || !got.Complete {
		t.Fatalf("very low PE percentile with insufficient drawdown should downshift to oneHalf, got %+v", got)
	}
}

func TestRuntimeMarketDataPersistsETFRuleStatuses(t *testing.T) {
	originalRuntimeQuotesFile := runtimeQuotesFile
	runtimeQuotesFile = filepath.Join(t.TempDir(), "runtime", "quotes.json")
	t.Cleanup(func() {
		runtimeQuotesFile = originalRuntimeQuotesFile
	})

	drawdown := 4.0
	cape := 82.0
	err := saveRuntimeMarketData(nil, []ETFRuleStatus{{
		Symbol:        "018738",
		Name:          "博时标普500ETF联接E(人民币)",
		Level:         "half",
		LevelLabel:    "0.5倍",
		WeeklyAmount:  2000,
		MonthlyAmount: 8000,
		Complete:      true,
		Metrics: []ETFRuleMetric{
			{Key: "drawdown252", Label: "近252交易日回撤", Value: &drawdown, Unit: "%", AsOf: "2026-07-03", Available: true},
			{Key: "capePercentile", Label: "S&P 500 Shiller CAPE近10年分位", Value: &cape, Unit: "%", AsOf: "2026-07-02", Available: true},
		},
		Sources: []ETFRuleSource{{Name: "价格源"}, {Name: "估值源"}},
	}}, "2026-07-06 08:00:00")
	if err != nil {
		t.Fatalf("saveRuntimeMarketData returned error: %v", err)
	}
	book, err := loadRuntimeQuoteBook()
	if err != nil {
		t.Fatalf("loadRuntimeQuoteBook returned error: %v", err)
	}
	status := book.ETFRuleStatuses["018738"]
	if status.Level != "half" || status.WeeklyAmount != 2000 {
		t.Fatalf("persisted status = %+v", status)
	}
}

func TestMergeETFRuleStatusUsesFreshPreviousMetric(t *testing.T) {
	drawdown := 4.0
	pePercentile := 70.0
	now := mustParseDate(t, "2026-07-06")
	existing := ETFRuleStatus{
		Symbol: "022434",
		Metrics: []ETFRuleMetric{
			{Key: "drawdown252", Label: "近252交易日回撤", Value: &drawdown, Unit: "%", AsOf: "2026-07-03", Available: true},
			{Key: "pePercentile", Label: "中证A500滚动PE分位", Value: &pePercentile, Unit: "%", AsOf: "2026-07-03", Available: true},
		},
	}
	next := ETFRuleStatus{
		Symbol:     "022434",
		LevelLabel: "待数据",
		Metrics: []ETFRuleMetric{
			{Key: "drawdown252", Label: "近252交易日回撤", Unit: "%", Available: false, Error: "temporary source error"},
			{Key: "pePercentile", Label: "中证A500滚动PE分位", Unit: "%", Available: false, Error: "source not configured"},
		},
		Sources: []ETFRuleSource{{Name: "价格源"}, {Name: "估值源"}},
	}
	merged := mergeETFRuleStatusWithExisting(next, existing, now)
	if merged.Level != "half" {
		t.Fatalf("merged status = %+v", merged)
	}
	if strings.Contains(merged.Reason, "沿用上次成功值") {
		t.Fatalf("reason should not expose fallback state, got %q", merged.Reason)
	}
	if !merged.Complete || merged.WeeklyAmount != 2500 || merged.MonthlyAmount != 10000 {
		t.Fatalf("fresh fallback metrics should remain executable, got %+v", merged)
	}
}

func TestMergeETFRuleStatusRejectsStalePreviousMetric(t *testing.T) {
	drawdown := 4.0
	pePercentile := 70.0
	now := mustParseDate(t, "2026-07-06")
	existing := ETFRuleStatus{
		Symbol: "022434",
		Metrics: []ETFRuleMetric{
			{Key: "drawdown252", Label: "近252交易日回撤", Value: &drawdown, Unit: "%", AsOf: "2026-06-20", Available: true},
			{Key: "pePercentile", Label: "中证A500滚动PE分位", Value: &pePercentile, Unit: "%", AsOf: "2026-07-03", Available: true},
		},
	}
	next := ETFRuleStatus{
		Symbol:     "022434",
		LevelLabel: "待数据",
		Metrics: []ETFRuleMetric{
			{Key: "drawdown252", Label: "近252交易日回撤", Unit: "%", Available: false, Error: "temporary source error"},
			{Key: "pePercentile", Label: "中证A500滚动PE分位", Unit: "%", Available: false, Error: "source not configured"},
		},
		Sources: []ETFRuleSource{{Name: "价格源"}, {Name: "估值源"}},
	}
	merged := mergeETFRuleStatusWithExisting(next, existing, now)
	if merged.Complete || merged.WeeklyAmount != 0 || merged.MonthlyAmount != 0 {
		t.Fatalf("stale fallback metrics should not be executable, got %+v", merged)
	}
}

func TestETFRuleStatusConfidenceAcceptsFreshCompleteMetrics(t *testing.T) {
	drawdown := 4.0
	pePercentile := 70.0
	config, ok := etfRuleConfigBySymbol("022434")
	if !ok {
		t.Fatal("missing A500 config")
	}
	status := ETFRuleStatus{
		Symbol:        "022434",
		Level:         "half",
		LevelLabel:    "0.5倍",
		WeeklyAmount:  2500,
		MonthlyAmount: 10000,
		Complete:      true,
		Metrics: []ETFRuleMetric{
			{Key: "drawdown252", Label: "近252交易日回撤", Value: &drawdown, Unit: "%", AsOf: "2026-07-03", Available: true},
			{Key: "pePercentile", Label: "中证A500滚动PE分位", Value: &pePercentile, Unit: "%", AsOf: "2026-07-03", Available: true},
		},
		Sources: []ETFRuleSource{
			{Name: "价格源"},
			{Name: "估值源"},
		},
	}
	issues := etfRuleStatusConfidenceIssues(status, config, mustParseDate(t, "2026-07-06"))
	if len(issues) != 0 {
		t.Fatalf("fresh complete metrics should have no issues, got %+v", issues)
	}
	checked := enforceETFRuleStatusConfidence(status, config, mustParseDate(t, "2026-07-06"))
	if !checked.Complete || checked.WeeklyAmount != 2500 {
		t.Fatalf("fresh complete metrics should remain executable, got %+v", checked)
	}
}

func TestETFRuleStatusConfidenceRejectsStaleMetrics(t *testing.T) {
	drawdown := 4.0
	pePercentile := 70.0
	config, ok := etfRuleConfigBySymbol("022434")
	if !ok {
		t.Fatal("missing A500 config")
	}
	status := ETFRuleStatus{
		Symbol:        "022434",
		Level:         "half",
		LevelLabel:    "0.5倍",
		WeeklyAmount:  2500,
		MonthlyAmount: 10000,
		Complete:      true,
		Metrics: []ETFRuleMetric{
			{Key: "drawdown252", Label: "近252交易日回撤", Value: &drawdown, Unit: "%", AsOf: "2026-06-20", Available: true},
			{Key: "pePercentile", Label: "中证A500滚动PE分位", Value: &pePercentile, Unit: "%", AsOf: "2026-07-03", Available: true},
		},
		Sources: []ETFRuleSource{
			{Name: "价格源"},
			{Name: "估值源"},
		},
	}
	issues := etfRuleStatusConfidenceIssues(status, config, mustParseDate(t, "2026-07-06"))
	if len(issues) == 0 {
		t.Fatal("stale drawdown should produce a confidence issue")
	}
	checked := enforceETFRuleStatusConfidence(status, config, mustParseDate(t, "2026-07-06"))
	if checked.Complete || checked.WeeklyAmount != 0 || checked.MonthlyAmount != 0 {
		t.Fatalf("stale metrics should not be executable, got %+v", checked)
	}
}

func mustParseDate(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("parse date %q: %v", value, err)
	}
	return parsed
}
