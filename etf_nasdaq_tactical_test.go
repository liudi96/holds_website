package main

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestParseUSTreasury10YXML(t *testing.T) {
	body := []byte(`<?xml version="1.0"?><feed xmlns="http://www.w3.org/2005/Atom" xmlns:d="http://schemas.microsoft.com/ado/2007/08/dataservices" xmlns:m="http://schemas.microsoft.com/ado/2007/08/dataservices/metadata"><entry><content><m:properties><d:NEW_DATE>2026-07-10T00:00:00</d:NEW_DATE><d:BC_10YEAR>4.44</d:BC_10YEAR></m:properties></content></entry></feed>`)
	closes, err := parseUSTreasury10YXML(body)
	if err != nil {
		t.Fatalf("parseUSTreasury10YXML returned error: %v", err)
	}
	if len(closes) != 1 || closes[0].Date != "2026-07-10" || !almostEqual(closes[0].Price, 0.0444, 0.000001) {
		t.Fatalf("unexpected Treasury closes: %+v", closes)
	}
}

func TestParseVXNHistoryCSV(t *testing.T) {
	closes, err := parseVXNHistoryCSV([]byte("DATE,OPEN,HIGH,LOW,CLOSE\n07/10/2026,25,26,24,24.89\n07/13/2026,26,27,25,26.50\n"))
	if err != nil {
		t.Fatalf("parseVXNHistoryCSV returned error: %v", err)
	}
	if len(closes) != 2 || closes[1].Date != "2026-07-13" || closes[1].Price != 26.5 {
		t.Fatalf("unexpected VXN closes: %+v", closes)
	}
}

func TestParseSinaNasdaqFuturesQuote(t *testing.T) {
	snapshot, err := parseSinaNasdaqFuturesQuote([]byte(`var hq_str_hf_NQ="29430.435,,29430.500,29432.500,29479.500,29427.250,06:37:14,29475.750,29446.000,0,4,1,2026-07-14,纳斯达克指数期货,0";`))
	if err != nil {
		t.Fatalf("parseSinaNasdaqFuturesQuote returned error: %v", err)
	}
	if snapshot.Date != "2026-07-14" || !almostEqual(snapshot.Change, 29430.435/29475.750-1, 0.000001) {
		t.Fatalf("unexpected Sina Nasdaq futures snapshot: %+v", snapshot)
	}
}

func TestCalculateNasdaqTacticalValuationUsesTenYearWeeklySeries(t *testing.T) {
	now := time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC)
	start := now.AddDate(-10, 0, -7)
	pePoints := []historyOfMarketPoint{}
	for month := start; !month.After(now); month = month.AddDate(0, 1, 0) {
		pePoints = append(pePoints, historyOfMarketPoint{Date: month.Format("2006-01-02"), Value: 20 + float64(month.Month())/10})
	}
	pePoints = append(pePoints, historyOfMarketPoint{Date: "2026-07-13", Value: 25})
	treasury := []dailyClose{}
	for cursor := start; !cursor.After(now); cursor = cursor.AddDate(0, 0, 7) {
		treasury = append(treasury, dailyClose{Date: cursor.Format("2006-01-02"), Price: 0.03 + float64(cursor.Year()-2016)/1000})
	}
	snapshot, err := calculateNasdaqTacticalValuation(pePoints, treasury, now)
	if err != nil {
		t.Fatalf("calculateNasdaqTacticalValuation returned error: %v", err)
	}
	if snapshot.ObservationCount < nasdaqTacticalMinimumWeeks {
		t.Fatalf("observation count = %d", snapshot.ObservationCount)
	}
	if snapshot.ForwardPE != 25 || snapshot.ForwardPEPercentile <= 0 || snapshot.SpreadPercentile < 0 || snapshot.SpreadPercentile > 1 {
		t.Fatalf("unexpected valuation snapshot: %+v", snapshot)
	}
}

func TestCalculateCNYTotalReturnClosesIncludesCurrency(t *testing.T) {
	closes, err := calculateCNYTotalReturnCloses(
		buildDailyCloseSeries(2100, 100, 0.1),
		buildDailyCloseSeries(2100, 7, 0.001),
	)
	if err != nil {
		t.Fatalf("calculateCNYTotalReturnCloses returned error: %v", err)
	}
	if len(closes) != 2100 || !almostEqual(closes[0].Price, 700, 0.000001) {
		t.Fatalf("unexpected CNY total-return closes: first=%+v len=%d", closes[0], len(closes))
	}
}

func TestEstimateNasdaqQDIIRealtimeNAV(t *testing.T) {
	xndx := []dailyClose{{Date: "2026-07-10", Price: 100}, {Date: "2026-07-13", Price: 102}}
	fx := []dailyClose{{Date: "2026-07-10", Price: 7}, {Date: "2026-07-13", Price: 7.07}}
	estimated, err := estimateNasdaqQDIIRealtimeNAV(
		2,
		"2026-07-10",
		xndx,
		fx,
		nasdaqFuturesSnapshot{Change: 0.01, Date: "2026-07-14"},
		nasdaqFuturesSnapshot{Change: -0.002, Date: "2026-07-14"},
		nil,
	)
	if err != nil {
		t.Fatalf("estimateNasdaqQDIIRealtimeNAV returned error: %v", err)
	}
	want := 2 * 1.02 * 1.01 * 1.01 * 0.998
	if math.Abs(estimated-want) > 0.000001 {
		t.Fatalf("estimated NAV = %.8f, want %.8f", estimated, want)
	}
}

func TestNasdaqTacticalETFUses159659(t *testing.T) {
	if nasdaqTacticalETFCode != "159659" {
		t.Fatalf("nasdaq tactical ETF = %q", nasdaqTacticalETFCode)
	}
}

func TestNormalizeNasdaqTacticalEastmoneyQuoteUsesETFScale(t *testing.T) {
	got := normalizeNasdaqTacticalEastmoneyQuote(quote{Price: 22.87, PreviousClose: 23})
	if !almostEqual(got.Price, 2.287, 0.000001) || !almostEqual(got.PreviousClose, 2.3, 0.000001) {
		t.Fatalf("unexpected normalized ETF quote: %+v", got)
	}
}

func buildDailyCloseSeries(count int, start float64, step float64) []dailyClose {
	result := make([]dailyClose, 0, count)
	date := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	for index := 0; index < count; index++ {
		result = append(result, dailyClose{Date: date.AddDate(0, 0, index).Format("2006-01-02"), Price: start + float64(index)*step})
	}
	return result
}

func Example_nasdaqTacticalETFCode() {
	fmt.Println(nasdaqTacticalETFCode)
	// Output: 159659
}
