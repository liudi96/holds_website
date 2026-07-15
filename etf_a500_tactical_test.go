package main

import (
	"math"
	"testing"
)

func TestA500TrackerUsesOfficialTotalReturnAnd159352(t *testing.T) {
	config, ok := etfRuleConfigBySymbol("022434")
	if !ok {
		t.Fatal("missing A500 tracker config")
	}
	if config.PriceSymbol != a500TotalReturnIndexCode || a500TotalReturnIndexCode != "000510CNY010" {
		t.Fatalf("A500 total-return source = %q", config.PriceSymbol)
	}
	if a500TacticalETFCode != "159352" {
		t.Fatalf("A500 tactical ETF = %q", a500TacticalETFCode)
	}
	if config.Monthly["one"] != 5000 || config.Weekly["one"] != 1250 {
		t.Fatalf("A500 fixed plan = %.0f/%.0f", config.Monthly["one"], config.Weekly["one"])
	}
}

func TestApplyA500MarketIndicatorsUsesAllTimeTotalReturnPeak(t *testing.T) {
	points := []a500PerformancePoint{
		{Date: "2026-06-30", Close: 9234.39},
		{Date: "2026-07-10", Close: 8731.70},
		{Date: "2026-07-13", Close: 8509.36},
	}
	var snapshot a500OpportunitySnapshot
	if err := applyA500MarketIndicators(&snapshot, points); err != nil {
		t.Fatalf("applyA500MarketIndicators returned error: %v", err)
	}
	want := 1 - 8509.36/9234.39
	if math.Abs(snapshot.Drawdown-want) > 0.000001 || snapshot.PeakDate != "2026-06-30" {
		t.Fatalf("snapshot = %+v, want drawdown %.8f", snapshot, want)
	}
}

func TestParseA500TacticalMarketQuoteIncludesBidAskAndOpen(t *testing.T) {
	body := []byte(`v_sz159352="51~A500ETF南方~159352~1.297~1.335~1.321~27709491~14263311~13446180~1.297~34602~1.296~2509~1.295~2674~1.294~15226~1.293~3823~1.298~314~1.300~1929~1.301~411~1.302~831~1.303~1014~~20260713161418";`)
	market, err := parseA500TacticalMarketQuote(body)
	if err != nil {
		t.Fatalf("parseA500TacticalMarketQuote returned error: %v", err)
	}
	if market.Price != 1.297 || market.PreviousClose != 1.335 || market.Open != 1.321 || market.Bid != 1.297 || market.Ask != 1.298 || market.Date != "2026-07-13" {
		t.Fatalf("market = %+v", market)
	}
}

func TestApplyA500TradingIndicatorsUsesIndexAdjustedNAV(t *testing.T) {
	points := []a500PerformancePoint{
		{Date: "2026-07-10", Close: 8731.70},
		{Date: "2026-07-13", Close: 8509.36},
	}
	market := a500TacticalMarketSnapshot{Price: 1.297, PreviousClose: 1.335, Open: 1.321, Bid: 1.297, Ask: 1.298, Date: "2026-07-13"}
	nav := quote{Price: 1.3314, PriceDate: "2026-07-10"}
	var snapshot a500OpportunitySnapshot
	if err := applyA500TradingIndicators(&snapshot, points, market, nav); err != nil {
		t.Fatalf("applyA500TradingIndicators returned error: %v", err)
	}
	wantNAV := 1.3314 * 8509.36 / 8731.70
	if math.Abs(snapshot.EstimatedNAV-wantNAV) > 0.000001 {
		t.Fatalf("estimated NAV = %.8f, want %.8f", snapshot.EstimatedNAV, wantNAV)
	}
	if snapshot.BidAskSpread <= 0 || snapshot.Premium > 0.003 {
		t.Fatalf("trading snapshot = %+v", snapshot)
	}
}

func TestEvaluateA500RuleUsesNeutralFallbackAndCombinedValuation(t *testing.T) {
	neutral := evaluateA500Rule(etfRuleInputs{})
	if neutral.Level != "one" || !neutral.Complete {
		t.Fatalf("missing valuation should be neutral, got %+v", neutral)
	}
	cheapPE, cheapSpread := 0.30, 0.70
	cheap := evaluateA500Rule(etfRuleInputs{ValuationPercentile: &cheapPE, EarningsYieldSpreadPercentile: &cheapSpread})
	if cheap.Level != "oneHalf" || !cheap.Complete {
		t.Fatalf("cheap valuation = %+v", cheap)
	}
	expensivePE, expensiveSpread := 0.75, 0.30
	expensive := evaluateA500Rule(etfRuleInputs{ValuationPercentile: &expensivePE, EarningsYieldSpreadPercentile: &expensiveSpread})
	if expensive.Level != "quarter" || !expensive.Complete {
		t.Fatalf("expensive valuation = %+v", expensive)
	}
}
