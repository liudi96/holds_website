package main

import (
	"archive/zip"
	"bytes"
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

func TestTotalReturnClosesReinvestsCashDividend(t *testing.T) {
	got, err := totalReturnCloses([]dailyClose{
		{Date: "2026-03-19", Price: 100},
		{Date: "2026-03-20", Price: 99},
		{Date: "2026-03-23", Price: 101},
	}, []cashDividendEvent{{Date: "2026-03-20", Amount: 2}})
	if err != nil {
		t.Fatalf("totalReturnCloses returned error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len(got) = %d, want 3", len(got))
	}
	if !almostEqual(got[1].Price, 101, 0.000001) {
		t.Fatalf("ex-dividend total-return value = %.6f, want 101", got[1].Price)
	}
	if !almostEqual(got[2].Price, 103.040404, 0.000001) {
		t.Fatalf("following total-return value = %.6f, want 103.040404", got[2].Price)
	}
}

func TestParseEastmoneyFundNAVTrendUsesUnitNAVAndCashDividend(t *testing.T) {
	body := []byte(`var Data_netWorthTrend = [
		{"x":1579190400000,"y":1.0,"equityReturn":0,"unitMoney":""},
		{"x":1579708800000,"y":0.9,"equityReturn":-10,"unitMoney":"分红：每份派现金0.1元"}
	]; var Data_ACWorthTrend = [];`)
	closes, events, err := parseEastmoneyFundNAVTrend(body)
	if err != nil {
		t.Fatalf("parseEastmoneyFundNAVTrend returned error: %v", err)
	}
	if len(closes) != 2 || closes[0].Date != "2020-01-17" || closes[1].Date != "2020-01-23" {
		t.Fatalf("closes = %+v", closes)
	}
	if len(events) != 1 || events[0].Date != "2020-01-23" || !almostEqual(events[0].Amount, 0.1, 0.000001) {
		t.Fatalf("events = %+v", events)
	}
	totalReturns, err := totalReturnCloses(closes, events)
	if err != nil {
		t.Fatalf("totalReturnCloses returned error: %v", err)
	}
	if !almostEqual(totalReturns[1].Price, 1.0, 0.000001) {
		t.Fatalf("ex-dividend total return = %.6f, want 1", totalReturns[1].Price)
	}
}

func TestParseStockAnalysisDividends(t *testing.T) {
	body := []byte(`<html><body><h2>Dividend History</h2><table><tbody>
		<tr><td>Jun 18, 2026</td><td>$1.90352</td><td>Jun 18, 2026</td></tr>
		<tr><td>Mar 20, 2026</td><td>$1.797</td><td>Mar 20, 2026</td></tr>
	</tbody></table></body></html>`)
	events, err := parseStockAnalysisDividends(body)
	if err != nil {
		t.Fatalf("parseStockAnalysisDividends returned error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("len(events) = %d, want 2", len(events))
	}
	if events[0].Date != "2026-06-18" || !almostEqual(events[0].Amount, 1.90352, 0.000001) {
		t.Fatalf("first event = %+v", events[0])
	}
}

func TestParseStateStreetSPYDividends(t *testing.T) {
	var workbook bytes.Buffer
	archive := zip.NewWriter(&workbook)
	sharedStrings, err := archive.Create("xl/sharedStrings.xml")
	if err != nil {
		t.Fatalf("create shared strings: %v", err)
	}
	_, err = sharedStrings.Write([]byte(`<sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
		<si><t>SPY</t></si>
		<si><r><t>06/18/</t></r><r><t>2026</t></r></si>
		<si><t>1.903516</t></si>
		<si><t>03/20/2026</t></si>
		<si><t>1.796999</t></si>
		<si><t>QQQ</t></si>
	</sst>`))
	if err != nil {
		t.Fatalf("write shared strings: %v", err)
	}
	sheet, err := archive.Create("xl/worksheets/sheet1.xml")
	if err != nil {
		t.Fatalf("create dividend sheet: %v", err)
	}
	_, err = sheet.Write([]byte(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>
		<row r="2"><c r="B2" t="s"><v>0</v></c><c r="D2" t="s"><v>1</v></c><c r="G2" t="s"><v>2</v></c></row>
		<row r="3"><c r="B3" t="s"><v>0</v></c><c r="D3" t="s"><v>3</v></c><c r="G3" t="s"><v>4</v></c></row>
		<row r="4"><c r="B4" t="s"><v>5</v></c><c r="D4" t="s"><v>3</v></c><c r="G4" t="s"><v>4</v></c></row>
	</sheetData></worksheet>`))
	if err != nil {
		t.Fatalf("write dividend sheet: %v", err)
	}
	if err := archive.Close(); err != nil {
		t.Fatalf("close workbook: %v", err)
	}

	events, err := parseStateStreetSPYDividends(workbook.Bytes())
	if err != nil {
		t.Fatalf("parseStateStreetSPYDividends returned error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("len(events) = %d, want 2", len(events))
	}
	if events[0].Date != "2026-03-20" || !almostEqual(events[0].Amount, 1.796999, 0.000001) {
		t.Fatalf("first event = %+v, want 2026-03-20 / 1.796999", events[0])
	}
	if events[1].Date != "2026-06-18" || !almostEqual(events[1].Amount, 1.903516, 0.000001) {
		t.Fatalf("second event = %+v, want 2026-06-18 / 1.903516", events[1])
	}
}

func TestParseEastmoneyChina10YBondYields(t *testing.T) {
	body := []byte(`{"result":{"pages":3,"data":[
		{"SOLAR_DATE":"2026-07-10 00:00:00","EMM00166466":1.7398},
		{"SOLAR_DATE":"2026-07-09 00:00:00","EMM00166466":1.7376}
	]}}`)
	points, pages, err := parseEastmoneyChina10YBondYields(body)
	if err != nil {
		t.Fatalf("parseEastmoneyChina10YBondYields returned error: %v", err)
	}
	if pages != 3 || len(points) != 2 {
		t.Fatalf("pages=%d len(points)=%d, want 3 and 2", pages, len(points))
	}
	if points[0].Date != "2026-07-10" || !almostEqual(points[0].Value, 0.017398, 0.0000001) {
		t.Fatalf("first point = %+v", points[0])
	}
}

func TestParseChinaBondOfficial10YYield(t *testing.T) {
	body := []byte(`<table>
		<tr><td>Yield Curve Name</td><td>Date</td><td>3M</td><td>6M</td><td>1Y</td><td>3Y</td><td>5Y</td><td>7Y</td><td>10Y</td></tr>
		<tr><td>ChinaBond Government Bond Yield Curve</td><td>2026-07-09</td><td></td><td></td><td></td><td></td><td></td><td></td><td>1.7376</td></tr>
		<tr><td>ChinaBond Government Bond Yield Curve</td><td>2026-07-10</td><td></td><td></td><td></td><td></td><td></td><td></td><td>1.7398</td></tr>
	</table>`)
	point, err := parseChinaBondOfficial10YYield(body)
	if err != nil {
		t.Fatalf("parseChinaBondOfficial10YYield returned error: %v", err)
	}
	if point.Date != "2026-07-10" || !almostEqual(point.Value, 0.017398, 0.0000001) {
		t.Fatalf("point = %+v", point)
	}
}

func TestCalculateDividendSpreadUsesAlignedHistory(t *testing.T) {
	yieldHistory := make([]datedRate, 0, 12)
	bondHistory := make([]datedRate, 0, 12)
	start := mustParseDate(t, "2025-08-01")
	for i := 0; i < 12; i++ {
		date := start.AddDate(0, i, 0).Format("2006-01-02")
		yieldHistory = append(yieldHistory, datedRate{Date: date, Value: 0.03 + float64(i)*0.001})
		bondHistory = append(bondHistory, datedRate{Date: date, Value: 0.02})
	}
	current := datedRate{Date: "2026-07-01", Value: 0.045}
	official := datedRate{Date: "2026-07-01", Value: 0.02}
	spread, percentile, observations, err := calculateDividendSpread(current, yieldHistory, bondHistory, official)
	if err != nil {
		t.Fatalf("calculateDividendSpread returned error: %v", err)
	}
	if !almostEqual(spread, 0.025, 0.000001) || !almostEqual(percentile, 1, 0.000001) || observations != 12 {
		t.Fatalf("spread=%.6f percentile=%.6f observations=%d", spread, percentile, observations)
	}
}

func TestSelectLatestDailyCloseCandidatePrefersNewestDate(t *testing.T) {
	closes, source, ok := selectLatestDailyCloseCandidate([]dailyCloseCandidate{
		{
			Source: "nasdaq",
			Closes: []dailyClose{
				{Date: "2026-07-06", Price: 100},
				{Date: "2026-07-07", Price: 101},
			},
		},
		{
			Source: "yahoo",
			Closes: []dailyClose{
				{Date: "2026-07-07", Price: 101},
				{Date: "2026-07-08", Price: 102},
			},
		},
	}, 10)
	if !ok {
		t.Fatal("selectLatestDailyCloseCandidate returned no candidate")
	}
	if source != "yahoo" {
		t.Fatalf("source = %q, want yahoo", source)
	}
	if got := latestDailyCloseDate(closes); got != "2026-07-08" {
		t.Fatalf("latest date = %q, want 2026-07-08", got)
	}
}

func TestSelectLatestDailyCloseCandidateKeepsFirstSourceOnTie(t *testing.T) {
	_, source, ok := selectLatestDailyCloseCandidate([]dailyCloseCandidate{
		{Source: "nasdaq", Closes: []dailyClose{{Date: "2026-07-08", Price: 100}}},
		{Source: "yahoo", Closes: []dailyClose{{Date: "2026-07-08", Price: 101}}},
	}, 10)
	if !ok {
		t.Fatal("selectLatestDailyCloseCandidate returned no candidate")
	}
	if source != "nasdaq" {
		t.Fatalf("source = %q, want nasdaq", source)
	}
}

func TestAppendOrReplaceLatestDailyCloseAppendsNewerQuote(t *testing.T) {
	closes := appendOrReplaceLatestDailyClose([]dailyClose{
		{Date: "2026-07-06", Price: 100},
		{Date: "2026-07-07", Price: 101},
	}, dailyClose{Date: "2026-07-08", Price: 102})
	if len(closes) != 3 {
		t.Fatalf("len = %d, want 3", len(closes))
	}
	if latest := closes[len(closes)-1]; latest.Date != "2026-07-08" || latest.Price != 102 {
		t.Fatalf("latest = %+v", latest)
	}
}

func TestAppendOrReplaceLatestDailyCloseReplacesSameDateQuote(t *testing.T) {
	closes := appendOrReplaceLatestDailyClose([]dailyClose{
		{Date: "2026-07-06", Price: 100},
		{Date: "2026-07-07", Price: 101},
	}, dailyClose{Date: "2026-07-07", Price: 102})
	if len(closes) != 2 {
		t.Fatalf("len = %d, want 2", len(closes))
	}
	if latest := closes[len(closes)-1]; latest.Date != "2026-07-07" || latest.Price != 102 {
		t.Fatalf("latest = %+v", latest)
	}
}

func TestNormalizeNasdaqQuoteDate(t *testing.T) {
	if got := normalizeNasdaqQuoteDate("Jul 8, 2026"); got != "2026-07-08" {
		t.Fatalf("date = %q, want 2026-07-08", got)
	}
}

func TestETFRuleBaseAmountsKeepTargetAllocation(t *testing.T) {
	wantMonthly := map[string]float64{
		"022434": 5000,
		"008163": 5000,
		"018738": 5000,
		"021000": 5000,
	}
	totalMonthly := 0.0
	totalWeekly := 0.0
	for _, config := range etfRuleConfigs {
		monthly := config.Monthly["one"]
		weekly := config.Weekly["one"]
		if monthly != wantMonthly[config.Symbol] {
			t.Fatalf("%s monthly one = %.0f, want %.0f", config.Symbol, monthly, wantMonthly[config.Symbol])
		}
		if weekly*4 != monthly {
			t.Fatalf("%s weekly one %.0f does not reconcile with monthly %.0f", config.Symbol, weekly, monthly)
		}
		for level := range etfRuleLevels {
			if config.Monthly[level] != monthly || config.Weekly[level] != weekly {
				t.Fatalf("%s water level %s changed fixed plan", config.Symbol, level)
			}
		}
		totalMonthly += monthly
		totalWeekly += weekly
	}
	if totalMonthly != 20000 || totalWeekly != 5000 {
		t.Fatalf("totals monthly=%.0f weekly=%.0f, want 20000/5000", totalMonthly, totalWeekly)
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

func TestParseFundDBIndexPEPercentile(t *testing.T) {
	body := []byte(`{
		"code": 0,
		"message": "",
		"data": {
			"update_time": "2026-07-08",
			"top_data": [
				{"attribute":"close","name":"收盘价","new_percent_value":{"value":"93.7%"}},
				{"attribute":"pe","name":"市盈率","new_percent_value":{"value":"86.69%"}},
				{"attribute":"pb","name":"市净率","new_percent_value":{"value":"59.21%"}}
			]
		}
	}`)
	point, err := parseFundDBIndexPEPercentile(body)
	if err != nil {
		t.Fatalf("parseFundDBIndexPEPercentile returned error: %v", err)
	}
	if point.Date != "2026-07-08" {
		t.Fatalf("date = %q, want 2026-07-08", point.Date)
	}
	if !almostEqual(point.Percentile, 0.8669, 0.000001) {
		t.Fatalf("percentile = %.6f, want 0.8669", point.Percentile)
	}
	pbPoint, err := parseFundDBIndexPBPercentile(body)
	if err != nil {
		t.Fatalf("parseFundDBIndexPBPercentile returned error: %v", err)
	}
	if pbPoint.Date != "2026-07-08" || !almostEqual(pbPoint.Percentile, 0.5921, 0.000001) {
		t.Fatalf("PB point = %+v, want 59.21%% on 2026-07-08", pbPoint)
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

func TestEvaluateSP500RuleUsesForwardPEAndEarningsSpread(t *testing.T) {
	pePercentile := 0.82
	spreadPercentile := 0.19
	drawdown := 0.12
	got := evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &pePercentile, EarningsYieldSpreadPercentile: &spreadPercentile, Drawdown: &drawdown})
	if got.Level != "quarter" || !got.Complete {
		t.Fatalf("expensive PE and spread combination should choose quarter, got %+v", got)
	}

	spreadPercentile = 0.50
	drawdown = 0.03
	got = evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &pePercentile, EarningsYieldSpreadPercentile: &spreadPercentile, Drawdown: &drawdown})
	if got.Level != "one" || !got.Complete {
		t.Fatalf("high PE alone must not mark the combined valuation expensive, got %+v", got)
	}

	pePercentile = 0.30
	spreadPercentile = 0.61
	drawdown = 0.12
	got = evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &pePercentile, EarningsYieldSpreadPercentile: &spreadPercentile, Drawdown: &drawdown})
	if got.Level != "oneHalf" || !got.Complete {
		t.Fatalf("cheap PE and spread combination should choose oneHalf, got %+v", got)
	}

	drawdown = 0.30
	got = evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &pePercentile, EarningsYieldSpreadPercentile: &spreadPercentile, Drawdown: &drawdown})
	if got.Level != "oneHalf" || !got.Complete {
		t.Fatalf("drawdown must not change the combined valuation water level, got %+v", got)
	}

	got = evaluateSP500Rule(etfRuleInputs{ValuationPercentile: &pePercentile})
	if got.Complete || got.Level != "" {
		t.Fatalf("missing earnings spread percentile must fail closed, got %+v", got)
	}
}

func TestEvaluateNasdaq100RuleUsesForwardPEAndEarningsSpread(t *testing.T) {
	pePercentile := 0.81
	spreadPercentile := 0.19
	got := evaluateNasdaq100Rule(etfRuleInputs{
		ValuationPercentile:           &pePercentile,
		EarningsYieldSpreadPercentile: &spreadPercentile,
	})
	if got.Level != "quarter" || !got.Complete {
		t.Fatalf("expensive PE and spread combination should choose quarter, got %+v", got)
	}

	spreadPercentile = 0.50
	got = evaluateNasdaq100Rule(etfRuleInputs{
		ValuationPercentile:           &pePercentile,
		EarningsYieldSpreadPercentile: &spreadPercentile,
	})
	if got.Level != "one" || !got.Complete {
		t.Fatalf("high PE alone must not mark the combined valuation expensive, got %+v", got)
	}

	pePercentile = 0.29
	got = evaluateNasdaq100Rule(etfRuleInputs{
		ValuationPercentile:           &pePercentile,
		EarningsYieldSpreadPercentile: &spreadPercentile,
	})
	if got.Level != "oneHalf" || !got.Complete {
		t.Fatalf("low forward PE percentile should mark valuation cheap, got %+v", got)
	}

	got = evaluateNasdaq100Rule(etfRuleInputs{ValuationPercentile: &pePercentile})
	if got.Complete || got.Level != "" {
		t.Fatalf("missing earnings spread percentile must fail closed, got %+v", got)
	}
}

func TestEvaluateDividendLowVolDoesNotUseAbsoluteYieldFallback(t *testing.T) {
	yield := 0.059
	percentile := 0.90
	drawdown := 0.05
	got := evaluateDividendLowVolRule(etfRuleInputs{
		DividendYield:           &yield,
		DividendYieldPercentile: &percentile,
		Drawdown:                &drawdown,
	})
	if got.Complete || got.Level != "" {
		t.Fatalf("absolute yield must not produce an investable water level, got %+v", got)
	}
}

func TestEvaluateDividendLowVolPrefersSpreadPercentile(t *testing.T) {
	spreadPercentile := 0.90
	drawdown := 0.13
	got := evaluateDividendLowVolRule(etfRuleInputs{DividendSpreadPercentile: &spreadPercentile, Drawdown: &drawdown})
	if got.Level != "two" || !got.Complete {
		t.Fatalf("high spread percentile should display low water level, got %+v", got)
	}
}

func TestEvaluateDividendLowVolUsesSpreadAndPBScore(t *testing.T) {
	spreadPercentile := 0.80
	pbPercentile := 0.60
	score := dividendLowVolValuationScore(spreadPercentile, pbPercentile)
	if !almostEqual(score, 0.70, 0.000001) {
		t.Fatalf("valuation score = %.6f, want 0.70", score)
	}
	got := evaluateDividendLowVolRule(etfRuleInputs{
		DividendSpreadPercentile: &spreadPercentile,
		PBPercentile:             &pbPercentile,
		ValuationScore:           &score,
	})
	if got.Level != "oneHalf" || !got.Complete {
		t.Fatalf("70%% valuation score should display oneHalf water level, got %+v", got)
	}
}

func TestRuntimeMarketDataPersistsETFRuleStatuses(t *testing.T) {
	originalRuntimeQuotesFile := runtimeQuotesFile
	runtimeQuotesFile = filepath.Join(t.TempDir(), "runtime", "quotes.json")
	t.Cleanup(func() {
		runtimeQuotesFile = originalRuntimeQuotesFile
	})
	if err := saveRuntimeQuoteBook(RuntimeQuoteBook{ETFRuleStatuses: map[string]ETFRuleStatus{
		"022434": {Symbol: "022434", Name: "南方中证A500ETF联接A", Complete: true},
	}}); err != nil {
		t.Fatalf("seed retired ETF status: %v", err)
	}

	drawdown := 4.0
	pePercentile := 82.0
	err := saveRuntimeMarketData(nil, []ETFRuleStatus{{
		Symbol:        "018738",
		Name:          "博时标普500ETF联接E(人民币)",
		Level:         "half",
		LevelLabel:    "偏高",
		WeeklyAmount:  1250,
		MonthlyAmount: 5000,
		Complete:      true,
		Metrics: []ETFRuleMetric{
			{Key: "drawdown252", Label: "近252交易日回撤", Value: &drawdown, Unit: "%", AsOf: "2026-07-03", Available: true},
			{Key: "pePercentile", Label: "标普500 PE分位", Value: &pePercentile, Unit: "%", AsOf: "2026-07-02", Available: true},
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
	if status.Level != "half" || status.WeeklyAmount != 1250 {
		t.Fatalf("persisted status = %+v", status)
	}
	if _, exists := book.ETFRuleStatuses["022434"]; !exists {
		t.Fatal("active A500 status should be preserved")
	}
}

func TestMergeETFRuleStatusUsesFreshPreviousMetric(t *testing.T) {
	now := mustParseDate(t, "2026-07-06")
	existing := ETFRuleStatus{
		Symbol:  "018738",
		Metrics: sp500TestMetrics("2026-07-03", true),
	}
	next := ETFRuleStatus{
		Symbol:     "018738",
		LevelLabel: "待数据",
		Metrics:    sp500TestMetrics("", false),
		Sources:    []ETFRuleSource{{Name: "价格源"}, {Name: "估值源"}},
	}
	merged := mergeETFRuleStatusWithExisting(next, existing, now)
	if merged.Level != "one" {
		t.Fatalf("merged status = %+v", merged)
	}
	if strings.Contains(merged.Reason, "沿用上次成功值") {
		t.Fatalf("reason should not expose fallback state, got %q", merged.Reason)
	}
	if !merged.Complete || merged.WeeklyAmount != 1250 || merged.MonthlyAmount != 5000 {
		t.Fatalf("fresh fallback metrics should remain executable, got %+v", merged)
	}
}

func TestMergeETFRuleStatusDoesNotReuseRetiredDividendSpreadProxy(t *testing.T) {
	oldSpreadPercentile := 55.36
	drawdown := 6.14
	now := mustParseDate(t, "2026-07-13")
	existing := ETFRuleStatus{
		Symbol:   "008163",
		Complete: true,
		Level:    "one",
		Metrics: []ETFRuleMetric{
			{Key: "dividendSpreadPercentile", Label: "股息率利差历史分位", Value: &oldSpreadPercentile, Unit: "%", AsOf: "2026-07-10", Available: true},
		},
	}
	next := ETFRuleStatus{
		Symbol:     "008163",
		LevelLabel: "待数据",
		Complete:   false,
		Metrics: []ETFRuleMetric{
			{Key: "drawdown3y", Label: "成立以来总回报回撤", Value: &drawdown, Unit: "%", AsOf: "2026-07-13", Available: true},
			{Key: "dividendSpreadPercentile", Label: "标的指数股债利差历史分位", Unit: "%", Available: false, Error: "标的指数TTM股息率历史数据源暂不可用；基金现金分红不代替指数股息率"},
			{Key: "pbPercentile", Label: "标的指数PB历史分位", Unit: "%", Available: false, Error: "missing PB percentile"},
			{Key: "valuationScore", Label: "红利估值得分V", Unit: "%", Available: false, Error: "需要股债利差分位和PB历史分位"},
		},
	}
	merged := mergeETFRuleStatusWithExisting(next, existing, now)
	metric := ETFRuleMetric{}
	found := false
	for _, item := range merged.Metrics {
		if item.Key == "dividendSpreadPercentile" {
			metric = item
			found = true
			break
		}
	}
	if !found || metric.Available || metric.Value != nil {
		t.Fatalf("retired proxy metric must stay unavailable, got %+v", metric)
	}
	if merged.Complete || merged.Level != "" {
		t.Fatalf("retired proxy metric must not produce an executable status, got %+v", merged)
	}
}

func TestMergeETFRuleStatusAllowsStaleDrawdownWhenValuationIsFresh(t *testing.T) {
	now := mustParseDate(t, "2026-07-06")
	existingMetrics := sp500TestMetrics("2026-07-03", true)
	existingMetrics[0].AsOf = "2026-06-20"
	existing := ETFRuleStatus{
		Symbol:  "018738",
		Metrics: existingMetrics,
	}
	next := ETFRuleStatus{
		Symbol:     "018738",
		LevelLabel: "待数据",
		Metrics:    sp500TestMetrics("", false),
		Sources:    []ETFRuleSource{{Name: "价格源"}, {Name: "估值源"}},
	}
	merged := mergeETFRuleStatusWithExisting(next, existing, now)
	if !merged.Complete || merged.WeeklyAmount != 1250 || merged.MonthlyAmount != 5000 {
		t.Fatalf("stale drawdown should not block a non-accelerated valuation level, got %+v", merged)
	}
}

func TestETFRuleStatusConfidenceAcceptsFreshCompleteMetrics(t *testing.T) {
	config, ok := etfRuleConfigBySymbol("018738")
	if !ok {
		t.Fatal("missing S&P 500 config")
	}
	status := ETFRuleStatus{
		Symbol:        "018738",
		Level:         "half",
		LevelLabel:    "偏高",
		WeeklyAmount:  1250,
		MonthlyAmount: 5000,
		Complete:      true,
		Metrics:       sp500TestMetrics("2026-07-03", true),
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
	if !checked.Complete || checked.WeeklyAmount != 1250 {
		t.Fatalf("fresh complete metrics should remain executable, got %+v", checked)
	}
}

func TestETFRuleStatusConfidenceRejectsStaleValuation(t *testing.T) {
	config, ok := etfRuleConfigBySymbol("018738")
	if !ok {
		t.Fatal("missing S&P 500 config")
	}
	metrics := sp500TestMetrics("2026-07-03", true)
	for index := range metrics {
		if metrics[index].Key == "forwardPE" || metrics[index].Key == "forwardPEPercentile" || metrics[index].Key == "earningsYieldSpreadPercentile" {
			metrics[index].AsOf = "2026-04-01"
		}
	}
	status := ETFRuleStatus{
		Symbol:        "018738",
		Level:         "half",
		LevelLabel:    "偏高",
		WeeklyAmount:  1250,
		MonthlyAmount: 5000,
		Complete:      true,
		Metrics:       metrics,
		Sources: []ETFRuleSource{
			{Name: "价格源"},
			{Name: "估值源"},
		},
	}
	issues := etfRuleStatusConfidenceIssues(status, config, mustParseDate(t, "2026-07-06"))
	if len(issues) == 0 {
		t.Fatal("stale valuation should produce a confidence issue")
	}
	checked := enforceETFRuleStatusConfidence(status, config, mustParseDate(t, "2026-07-06"))
	if checked.Complete || checked.WeeklyAmount != 1250 || checked.MonthlyAmount != 5000 {
		t.Fatalf("stale water metrics should not remove the fixed plan, got %+v", checked)
	}
}

func TestStabilizeETFRuleLevelRequiresFiveDistinctTradingDates(t *testing.T) {
	config, ok := etfRuleConfigBySymbol("018738")
	if !ok {
		t.Fatal("missing S&P 500 config")
	}
	drawdown := 13.0
	existing := ETFRuleStatus{
		Symbol:            "018738",
		Level:             "one",
		LevelLabel:        "中性",
		Complete:          true,
		LevelUpdatedAt:    "2026-07-01",
		PendingLevel:      "oneHalf",
		PendingLevelLabel: "偏低",
		PendingSince:      "2026-07-06",
		PendingAsOf:       "2026-07-09",
		PendingDays:       4,
	}
	next := ETFRuleStatus{
		Symbol:        "018738",
		Level:         "oneHalf",
		LevelLabel:    "偏低",
		MonthlyAmount: config.Monthly["oneHalf"],
		WeeklyAmount:  config.Weekly["oneHalf"],
		Complete:      true,
		AsOf:          "2026-07-10",
		Metrics: []ETFRuleMetric{
			{Key: "drawdown3y", Label: "近3年总收益回撤", Value: &drawdown, Unit: "%", AsOf: "2026-07-10", Available: true},
		},
	}
	got := stabilizeETFRuleLevel(next, existing, config)
	if got.Level != "oneHalf" || got.PendingLevel != "" || got.LevelUpdatedAt != "2026-07-10" {
		t.Fatalf("fifth distinct trading date should switch level, got %+v", got)
	}
}

func TestStabilizeETFRuleLevelDoesNotCountSameTradingDateTwice(t *testing.T) {
	config, ok := etfRuleConfigBySymbol("018738")
	if !ok {
		t.Fatal("missing S&P 500 config")
	}
	drawdown := 13.0
	existing := ETFRuleStatus{
		Symbol:       "018738",
		Level:        "one",
		Complete:     true,
		PendingLevel: "oneHalf",
		PendingAsOf:  "2026-07-09",
		PendingDays:  3,
	}
	next := ETFRuleStatus{
		Symbol:        "018738",
		Level:         "oneHalf",
		LevelLabel:    "偏低",
		MonthlyAmount: config.Monthly["oneHalf"],
		WeeklyAmount:  config.Weekly["oneHalf"],
		Complete:      true,
		Metrics: []ETFRuleMetric{
			{Key: "drawdown3y", Label: "近3年总收益回撤", Value: &drawdown, Unit: "%", AsOf: "2026-07-09", Available: true},
		},
	}
	got := stabilizeETFRuleLevel(next, existing, config)
	if got.Level != "one" || got.PendingDays != 3 {
		t.Fatalf("same trading date must not advance confirmation, got %+v", got)
	}
}

func sp500TestMetrics(date string, available bool) []ETFRuleMetric {
	specs := []struct {
		key   string
		label string
		value float64
		unit  string
	}{
		{key: "drawdown3y", label: "SPTR全收益高点回撤", value: 4, unit: "%"},
		{key: "forwardPE", label: "标普500未来12个月PE", value: 20.5},
		{key: "forwardPEPercentile", label: "未来PE十年周度分位", value: 70, unit: "%"},
		{key: "us10YBondYield", label: "美国10年期国债收益率", value: 4.5, unit: "%"},
		{key: "earningsYieldSpreadPercentile", label: "盈利利差十年周度分位", value: 50, unit: "%"},
		{key: "vix", label: "VIX", value: 18},
		{key: "cnyTotalReturnDrawdown", label: "人民币口径全收益回撤", value: 3, unit: "%"},
		{key: "qdiiPremium", label: "513650估算溢价", value: 0.3, unit: "%"},
	}
	metrics := make([]ETFRuleMetric, 0, len(specs))
	for _, spec := range specs {
		metric := ETFRuleMetric{Key: spec.key, Label: spec.label, Unit: spec.unit, AsOf: date, Available: available}
		if available {
			value := spec.value
			metric.Value = &value
		} else {
			metric.Error = "temporary source error"
		}
		metrics = append(metrics, metric)
	}
	return metrics
}

func mustParseDate(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("parse date %q: %v", value, err)
	}
	return parsed
}
