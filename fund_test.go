package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApplyFundBuyAndSellTrade(t *testing.T) {
	state := AppState{Cash: 10000, FX: map[string]float64{"CNY": 1}}

	applyTradeToState(&state, Trade{
		AssetType:    "fund",
		Symbol:       "000001",
		Name:         "Test Fund",
		Side:         "buy",
		Shares:       1000,
		Price:        1.2,
		Currency:     "CNY",
		CurrentPrice: 1.25,
	})

	if len(state.Funds) != 1 {
		t.Fatalf("funds len = %d, want 1", len(state.Funds))
	}
	if state.Funds[0].Shares != 1000 || state.Funds[0].Cost != 1.2 {
		t.Fatalf("fund position after buy = %+v", state.Funds[0])
	}
	if state.Cash != 8800 {
		t.Fatalf("cash after buy = %.2f, want 8800", state.Cash)
	}

	applyTradeToState(&state, Trade{
		AssetType:    "fund",
		Symbol:       "000001",
		Name:         "Test Fund",
		Side:         "sell",
		Shares:       200,
		Price:        1.3,
		Currency:     "CNY",
		CurrentPrice: 1.32,
	})

	if state.Funds[0].Shares != 800 || state.Funds[0].Cost != 1.2 {
		t.Fatalf("fund position after sell = %+v", state.Funds[0])
	}
	if state.Cash != 9060 {
		t.Fatalf("cash after sell = %.2f, want 9060", state.Cash)
	}
}

func TestReverseFundTradeRestoresCashAndShares(t *testing.T) {
	state := AppState{
		Cash: 8800,
		FX:   map[string]float64{"CNY": 1},
		Funds: []Fund{{
			Symbol:       "000001",
			Name:         "Test Fund",
			FundType:     "otc",
			Shares:       1000,
			Cost:         1.2,
			CurrentPrice: 1.25,
			Currency:     "CNY",
		}},
	}
	reverseTradeFromState(&state, Trade{
		AssetType:    "fund",
		Symbol:       "000001",
		Name:         "Test Fund",
		Side:         "buy",
		Shares:       200,
		Price:        1.2,
		Currency:     "CNY",
		CurrentPrice: 1.25,
	})

	if len(state.Funds) != 1 || state.Funds[0].Shares != 800 {
		t.Fatalf("fund position after reverse = %+v", state.Funds)
	}
	if state.Cash != 9040 {
		t.Fatalf("cash after reverse = %.2f, want 9040", state.Cash)
	}
}

func TestParseFundGZQuote(t *testing.T) {
	payload, err := parseFundGZQuote([]byte(`jsonpgz({"fundcode":"000001","name":"Test Fund","jzrq":"2026-07-03","dwjz":"1.2345","gsz":"1.2300","gztime":"2026-07-03 15:00"})`))
	if err != nil {
		t.Fatalf("parseFundGZQuote() error = %v", err)
	}
	if payload.FundCode != "000001" || payload.NAV != "1.2345" || payload.JZDate != "2026-07-03" {
		t.Fatalf("payload = %+v", payload)
	}
}

func TestNormalizeFundSymbolStripsOF(t *testing.T) {
	if got := normalizeFundSymbol("004814.OF"); got != "004814" {
		t.Fatalf("normalizeFundSymbol() = %q, want 004814", got)
	}
}

func TestResolveFundTradeInputParsesCodeAndName(t *testing.T) {
	server := Server{state: AppState{FX: map[string]float64{"CNY": 1}}}
	trade := Trade{
		AssetType:    "fund",
		Name:         "000001 Test Fund",
		Side:         "buy",
		Shares:       100,
		Price:        1.2,
		CurrentPrice: 1.2,
		Currency:     "CNY",
		Reason:       "allocation",
	}
	if err := server.resolveTradeInput(&trade); err != nil {
		t.Fatalf("resolveTradeInput() error = %v", err)
	}
	if trade.Symbol != "000001" || trade.Name != "Test Fund" {
		t.Fatalf("trade = %+v, want parsed fund code and name", trade)
	}
}

func TestHandleUpsertFundCreatesFund(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{FX: map[string]float64{"CNY": 1}}}
	req := httptest.NewRequest(http.MethodPost, "/api/funds", strings.NewReader(`{"symbol":"000001","name":"Test Fund","fundType":"otc","shares":100,"cost":1.2,"currency":"CNY"}`))
	rec := httptest.NewRecorder()

	server.handleUpsertFund(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	idx := findFundIndex(server.state.Funds, "000001")
	if idx == -1 || server.state.Funds[idx].Name != "Test Fund" || server.state.Funds[idx].Shares != 100 {
		t.Fatalf("fund not created: %+v", server.state.Funds)
	}
}

func TestHandleDeleteFundRejectsPositionedFund(t *testing.T) {
	withTempPortfolioData(t)
	server := Server{state: AppState{
		FX:    map[string]float64{"CNY": 1},
		Funds: []Fund{{Symbol: "000001", Name: "Test Fund", Shares: 100, Cost: 1.2, Currency: "CNY"}},
	}}
	req := httptest.NewRequest(http.MethodDelete, "/api/funds/000001", nil)
	rec := httptest.NewRecorder()

	server.handleDeleteFund(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409; body=%s", rec.Code, rec.Body.String())
	}
	if len(server.state.Funds) != 1 {
		t.Fatalf("positioned fund should remain: %+v", server.state.Funds)
	}
}

func TestMarshalJSONIncludesFunds(t *testing.T) {
	body, err := json.Marshal(AppState{Funds: []Fund{{Symbol: "000001", Name: "Test Fund"}}})
	if err != nil {
		t.Fatalf("marshal AppState: %v", err)
	}
	if !json.Valid(body) || !strings.Contains(string(body), `"funds"`) {
		t.Fatalf("expected funds in JSON, got %s", string(body))
	}
}
