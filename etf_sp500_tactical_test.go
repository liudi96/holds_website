package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSP500TacticalETFUses513650(t *testing.T) {
	if sp500TacticalETFCode != "513650" {
		t.Fatalf("S&P 500 tactical ETF = %q", sp500TacticalETFCode)
	}
	config, ok := etfRuleConfigBySymbol("018738")
	if !ok {
		t.Fatal("missing S&P 500 rule config")
	}
	if config.PriceSymbol != sp500TotalReturnSymbol || config.ValuationMetricKey != "forwardPEPercentile" {
		t.Fatalf("unexpected S&P 500 config: %+v", config)
	}
}

func TestFetchLatestCboeVolatilityIndexParsesVIX(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/csv")
		_, _ = writer.Write([]byte("DATE,OPEN,HIGH,LOW,CLOSE\n07/10/2026,15,16,14,15.03\n07/13/2026,16,17,15,16.25\n"))
	}))
	defer server.Close()

	latest, err := fetchLatestCboeVolatilityIndex(server.Client(), server.URL, "VIX")
	if err != nil {
		t.Fatalf("fetchLatestCboeVolatilityIndex returned error: %v", err)
	}
	if latest.Date != "2026-07-13" || latest.Price != 16.25 {
		t.Fatalf("unexpected VIX close: %+v", latest)
	}
}

func TestSP500ValuationUsesTenYearForwardPEAndTreasurySeries(t *testing.T) {
	now := time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC)
	start := now.AddDate(-10, 0, -7)
	pePoints := []historyOfMarketPoint{}
	for quarter := start; !quarter.After(now); quarter = quarter.AddDate(0, 3, 0) {
		pePoints = append(pePoints, historyOfMarketPoint{Date: quarter.Format("2006-01-02"), Value: 16 + float64(quarter.Month())/10})
	}
	pePoints = append(pePoints, historyOfMarketPoint{Date: "2026-07-11", Value: 20.92})
	treasury := []dailyClose{}
	for cursor := start; !cursor.After(now); cursor = cursor.AddDate(0, 0, 7) {
		treasury = append(treasury, dailyClose{Date: cursor.Format("2006-01-02"), Price: 0.03 + float64(cursor.Year()-2016)/1000})
	}

	snapshot, err := calculateNasdaqTacticalValuation(pePoints, treasury, now)
	if err != nil {
		t.Fatalf("calculate tactical valuation returned error: %v", err)
	}
	if snapshot.ForwardPE != 20.92 || snapshot.ObservationCount < nasdaqTacticalMinimumWeeks {
		t.Fatalf("unexpected S&P 500 valuation snapshot: %+v", snapshot)
	}
	if snapshot.ForwardPEPercentile < 0 || snapshot.ForwardPEPercentile > 1 || snapshot.SpreadPercentile < 0 || snapshot.SpreadPercentile > 1 {
		t.Fatalf("valuation percentiles out of range: %+v", snapshot)
	}
}

func TestSP500OpportunitySnapshotDisablesUnverifiedEarningsAcceleration(t *testing.T) {
	_, issues := fetchSP500OpportunitySnapshot(nil, time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC))
	if issues.EarningsRevision == nil {
		t.Fatal("unverified three-month earnings revision must disable early-stage acceleration")
	}
}
