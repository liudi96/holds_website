package main

import (
	"math"
	"testing"
)

func TestFundTradeCreatesFundHolding(t *testing.T) {
	server := Server{state: AppState{
		Cash: 10000,
		FX:   map[string]float64{"CNY": 1},
	}}
	trade := Trade{
		AssetType: "fund",
		Name:      "中欧红利优享混合",
		Side:      "buy",
		Shares:    1000,
		Price:     1.25,
	}

	if err := server.resolveTradeInput(&trade); err != nil {
		t.Fatalf("resolveTradeInput() error = %v", err)
	}
	if err := validateTrade(trade); err != nil {
		t.Fatalf("validateTrade() error = %v", err)
	}
	trade.Date = "2026-05-15"
	server.applyTrade(trade)

	if len(server.state.Funds) != 1 {
		t.Fatalf("fund count = %d, want 1", len(server.state.Funds))
	}
	fund := server.state.Funds[0]
	if fund.Name != "中欧红利优享混合" || fund.Symbol != "中欧红利优享混合" {
		t.Fatalf("unexpected fund identity: %+v", fund)
	}
	if fund.Shares != 1000 || fund.Cost != 1.25 || fund.CurrentNAV != 1.25 {
		t.Fatalf("unexpected fund holding values: %+v", fund)
	}
	if server.state.Cash != 8750 {
		t.Fatalf("cash = %v, want 8750", server.state.Cash)
	}
	if len(server.state.Trades) != 1 || server.state.Trades[0].AssetType != "fund" {
		t.Fatalf("unexpected trades: %+v", server.state.Trades)
	}
}

func TestFundTradeClearsFundWithoutCandidate(t *testing.T) {
	server := Server{state: AppState{
		Cash: 0,
		FX:   map[string]float64{"CNY": 1},
		Funds: []Fund{{
			Symbol:     "稳健债基",
			Name:       "稳健债基",
			Shares:     100,
			Cost:       1,
			CurrentNAV: 1.1,
			Currency:   "CNY",
		}},
	}}
	trade := Trade{
		AssetType: "fund",
		Name:      "稳健债基",
		Side:      "sell",
		Shares:    100,
		Price:     1.1,
		Date:      "2026-05-15",
	}

	if err := server.resolveTradeInput(&trade); err != nil {
		t.Fatalf("resolveTradeInput() error = %v", err)
	}
	server.applyTrade(trade)

	if len(server.state.Funds) != 0 {
		t.Fatalf("fund count = %d, want 0", len(server.state.Funds))
	}
	if len(server.state.Candidates) != 0 {
		t.Fatalf("candidate count = %d, want 0", len(server.state.Candidates))
	}
	if math.Abs(server.state.Cash-110) > 0.000001 {
		t.Fatalf("cash = %v, want 110", server.state.Cash)
	}
}

func TestParseEastmoneyFundNAV(t *testing.T) {
	nav, err := parseEastmoneyFundNAV("004814.OF", `jsonpgz({"fundcode":"004814","name":"中欧红利优享混合","jzrq":"2026-05-14","dwjz":"2.0600","gsz":"2.0612","gztime":"2026-05-15 15:00"});`)
	if err != nil {
		t.Fatalf("parseEastmoneyFundNAV() error = %v", err)
	}
	if nav.Symbol != "004814.OF" || nav.SourceSymbol != "004814" {
		t.Fatalf("unexpected nav identity: %+v", nav)
	}
	if nav.Name != "中欧红利优享混合" || nav.Currency != "CNY" {
		t.Fatalf("unexpected nav metadata: %+v", nav)
	}
	if nav.CurrentNAV != 2.06 || nav.CurrentNAVDate != "2026-05-14" {
		t.Fatalf("unexpected nav values: %+v", nav)
	}
}

func TestApplyFundNAV(t *testing.T) {
	fund := Fund{
		Symbol:     "004814.OF",
		Name:       "中欧红利优享混合",
		Shares:     40000,
		Cost:       1.98,
		CurrentNAV: 2.01,
		Currency:   "CNY",
	}
	applyFundNAV(&fund, fundNAV{
		Symbol:         "004814.OF",
		CurrentNAV:     2.06,
		Currency:       "CNY",
		CurrentNAVDate: "2026-05-14",
		SourceSymbol:   "004814",
		SourceName:     "东方财富基金净值",
	}, "2026-05-15 15:30:00")

	if fund.CurrentNAV != 2.06 || fund.CurrentNAVDate != "2026-05-14" {
		t.Fatalf("unexpected fund nav: %+v", fund)
	}
	if fund.UpdatedAt != "2026-05-15 15:30:00；净值源 东方财富基金净值；代码 004814" {
		t.Fatalf("unexpected updatedAt: %q", fund.UpdatedAt)
	}
}
