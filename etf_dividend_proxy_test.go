package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

type dividendProxyRoundTripFunc func(*http.Request) (*http.Response, error)

func (fn dividendProxyRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func jsonHTTPResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestFetchSouthernETFBasketParsesOfficialPCF(t *testing.T) {
	rows := make([]string, 0, 50)
	for i := 1; i <= 50; i++ {
		rows = append(rows, `{"stockCode":"`+leftPadNumber(i, 6)+`","stockName":"测试","stockQuality":"100"}`)
	}
	client := &http.Client{Transport: dividendProxyRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost || req.URL.String() != southernETFPCFAPIURL {
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL)
		}
		body, _ := io.ReadAll(req.Body)
		form, _ := url.ParseQuery(string(body))
		if form.Get("fundCode") != "515450" || form.Get("queryDate") != "20260714" {
			t.Fatalf("unexpected PCF form: %v", form)
		}
		return jsonHTTPResponse(`{"code":"ETS-5BP00000","message":"操作成功","data":{"TradingDay":"20260714","list":[` + strings.Join(rows, ",") + `]}}`), nil
	})}

	basket, err := fetchSouthernETFBasket(client, "20260714")
	if err != nil {
		t.Fatalf("fetchSouthernETFBasket returned error: %v", err)
	}
	if basket.Date != "2026-07-14" || len(basket.Components) != 50 {
		t.Fatalf("basket = %+v", basket)
	}
	if basket.Components[0].Code != "000001" || basket.Components[0].Quantity != 100 {
		t.Fatalf("first component = %+v", basket.Components[0])
	}
}

func leftPadNumber(value int, width int) string {
	text := "0000000000" + strconvItoa(value)
	return text[len(text)-width:]
}

func strconvItoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 10)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for left, right := 0, len(digits)-1; left < right; left, right = left+1, right-1 {
		digits[left], digits[right] = digits[right], digits[left]
	}
	return string(digits)
}

func TestFetchEastmoneyStockValuationsAndDividends(t *testing.T) {
	client := &http.Client{Transport: dividendProxyRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		switch req.URL.Query().Get("reportName") {
		case "RPT_VALUEANALYSIS_DET":
			if req.URL.Query().Get("filter") != "(TRADE_DATE='2026-07-14')" {
				t.Fatalf("unexpected valuation filter: %q", req.URL.Query().Get("filter"))
			}
			return jsonHTTPResponse(`{"success":true,"message":"ok","result":{"data":[{"SECURITY_CODE":"000651","CLOSE_PRICE":38.71,"PB_MRQ":1.43}]}}`), nil
		case "RPT_SHAREBONUS_DET":
			if req.URL.Query().Get("filter") != `(SECURITY_CODE in ("000651"))` {
				t.Fatalf("unexpected dividend filter: %q", req.URL.Query().Get("filter"))
			}
			return jsonHTTPResponse(`{"success":true,"message":"ok","result":{"data":[{"SECURITY_CODE":"000651","PRETAX_BONUS_RMB":20,"EX_DIVIDEND_DATE":"2025-08-29 00:00:00"},{"SECURITY_CODE":"000651","PRETAX_BONUS_RMB":10,"EX_DIVIDEND_DATE":"2026-01-23 00:00:00"}]}}`), nil
		default:
			t.Fatalf("unexpected report: %s", req.URL.Query().Get("reportName"))
			return nil, nil
		}
	})}

	valuations, err := fetchEastmoneyStockValuations(client, "2026-07-14")
	if err != nil {
		t.Fatalf("fetchEastmoneyStockValuations returned error: %v", err)
	}
	if got := valuations["000651"]; got.Close != 38.71 || got.PB != 1.43 {
		t.Fatalf("valuation = %+v", got)
	}
	events, err := fetchEastmoneyStockDividends(client, "000651")
	if err != nil {
		t.Fatalf("fetchEastmoneyStockDividends returned error: %v", err)
	}
	if len(events) != 2 || events[0].Amount != 2 || events[1].Amount != 1 {
		t.Fatalf("events = %+v", events)
	}
}

func TestFetchStockDividendHistoriesUsesBatchAndKeepsNoDividendStocks(t *testing.T) {
	client := &http.Client{Transport: dividendProxyRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if filter := req.URL.Query().Get("filter"); filter != `(SECURITY_CODE in ("000651","600900"))` {
			t.Fatalf("unexpected batch filter: %q", filter)
		}
		return jsonHTTPResponse(`{"success":true,"message":"ok","result":{"data":[{"SECURITY_CODE":"600900","PRETAX_BONUS_RMB":8,"EX_DIVIDEND_DATE":"2026-05-20 00:00:00"}]}}`), nil
	})}

	histories, errs := fetchStockDividendHistories(client, map[string]struct{}{"600900": {}, "000651": {}})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if events, ok := histories["000651"]; !ok || len(events) != 0 {
		t.Fatalf("no-dividend stock should remain valid: %+v", histories)
	}
	if events := histories["600900"]; len(events) != 1 || events[0].Amount != 0.8 {
		t.Fatalf("batch dividend events = %+v", events)
	}
}

func TestCalculateDividendLowVolProxyPointUsesTTMAndHarmonicPB(t *testing.T) {
	basket := southernETFBasket{Date: "2026-07-14"}
	valuations := map[string]eastmoneyStockValuation{}
	dividends := map[string][]cashDividendEvent{}
	for i := 1; i <= 20; i++ {
		code := leftPadNumber(i, 6)
		quantity := 10.0
		if i == 20 {
			quantity = 0
		}
		basket.Components = append(basket.Components, southernETFBasketComponent{Code: code, Quantity: quantity})
		valuations[code] = eastmoneyStockValuation{Code: code, Close: 10, PB: 2}
		dividends[code] = []cashDividendEvent{
			{Date: "2025-07-14", Amount: 9},
			{Date: "2025-07-15", Amount: 0.5},
			{Date: "2026-07-14", Amount: 0.5},
		}
	}

	point, err := calculateDividendLowVolProxyPoint(basket, valuations, dividends)
	if err != nil {
		t.Fatalf("calculateDividendLowVolProxyPoint returned error: %v", err)
	}
	if !almostEqual(point.DividendYield, 0.10, 0.000001) || !almostEqual(point.PB, 2, 0.000001) {
		t.Fatalf("point = %+v", point)
	}
	if point.ValidComponentCount != 19 || !almostEqual(point.Coverage, 0.95, 0.000001) {
		t.Fatalf("coverage point = %+v", point)
	}

	basket.Components[18].Quantity = 0
	if _, err := calculateDividendLowVolProxyPoint(basket, valuations, dividends); err == nil || !strings.Contains(err.Error(), "coverage") {
		t.Fatalf("coverage below 95%% should fail, got %v", err)
	}
}

func TestNormalizeDividendLowVolProxyPointsKeepsLatestWeeklyPoint(t *testing.T) {
	points := normalizeDividendLowVolProxyPoints([]dividendLowVolProxyHistoryPoint{
		{Date: "2026-07-13", DividendYield: 0.05, PB: 0.80, Coverage: 0.98},
		{Date: "2026-07-17", DividendYield: 0.051, PB: 0.79, Coverage: 0.98},
		{Date: "2026-07-20", DividendYield: 0.052, PB: 0.78, Coverage: 0.94},
	})
	if len(points) != 1 || points[0].Date != "2026-07-17" {
		t.Fatalf("points = %+v", points)
	}
}

func TestDividendLowVolLegacyMetricsRequireBasketProxySource(t *testing.T) {
	legacy := ETFRuleStatus{Sources: []ETFRuleSource{{Name: "韭圈儿标的指数PB分位"}}}
	if dividendLowVolStatusUsesBasketProxy(legacy) {
		t.Fatal("legacy source must not be treated as basket proxy")
	}
	current := ETFRuleStatus{Sources: []ETFRuleSource{{Name: "南方基金515450申购赎回篮子 + 东方财富成分股PB与分红（场内代理）"}}}
	if !dividendLowVolStatusUsesBasketProxy(current) {
		t.Fatal("current basket source should be recognized")
	}
}

func TestValidateDividendLowVolProxyHistoryRequiresCoverageAndObservationCount(t *testing.T) {
	start := time.Date(2021, 7, 13, 0, 0, 0, 0, time.UTC)
	points := make([]dividendLowVolProxyHistoryPoint, 0, dividendLowVolMinimumObservations)
	for i := 0; i < dividendLowVolMinimumObservations; i++ {
		points = append(points, dividendLowVolProxyHistoryPoint{
			Date:          start.AddDate(0, 0, i*8).Format(etfRuleRuntimeTimestampDateLayout),
			DividendYield: 0.05,
			PB:            0.8,
			Coverage:      0.98,
		})
	}
	if err := validateDividendLowVolProxyHistory(points); err != nil {
		t.Fatalf("five-year history with %d observations should pass: %v", len(points), err)
	}
	if err := validateDividendLowVolProxyHistory(points[:len(points)-1]); err == nil || !strings.Contains(err.Error(), "weekly observations") {
		t.Fatalf("history below the observation floor should fail, got %v", err)
	}

	shortSpan := make([]dividendLowVolProxyHistoryPoint, dividendLowVolMinimumObservations)
	for i := range shortSpan {
		shortSpan[i] = dividendLowVolProxyHistoryPoint{
			Date:          start.AddDate(0, 0, i*7).Format(etfRuleRuntimeTimestampDateLayout),
			DividendYield: 0.05,
			PB:            0.8,
			Coverage:      0.98,
		}
	}
	if err := validateDividendLowVolProxyHistory(shortSpan); err == nil || !strings.Contains(err.Error(), "does not span five years") {
		t.Fatalf("short-span history should fail, got %v", err)
	}
}

func TestSummarizeDividendLowVolRefreshErrorsIsConcise(t *testing.T) {
	summary := summarizeDividendLowVolRefreshErrors([]string{strings.Repeat("x", 300), "second"})
	if len(summary) > 240 || !strings.Contains(summary, "2 fetch or calculation errors") || !strings.HasSuffix(summary, "...") {
		t.Fatalf("unexpected summary: %q", summary)
	}
}

func TestDividendLowVolHistoryTargetsSkipsPreviouslyScannedHistoricalGaps(t *testing.T) {
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	start := time.Date(2021, 7, 13, 0, 0, 0, 0, location)
	history := dividendLowVolProxyHistory{UpdatedAt: "2026-07-20 20:30:00"}
	for i := 0; i < dividendLowVolMinimumObservations; i++ {
		history.Points = append(history.Points, dividendLowVolProxyHistoryPoint{
			Date:          start.AddDate(0, 0, i*8).Format(etfRuleRuntimeTimestampDateLayout),
			DividendYield: 0.05,
			PB:            0.8,
			Coverage:      0.98,
		})
	}
	now := time.Date(2026, 8, 4, 0, 0, 0, 0, location)
	targets := dividendLowVolHistoryTargets(now, now.AddDate(-dividendLowVolHistoryYears, 0, 0), history)
	lastRefresh := parseDividendLowVolHistoryRefreshTime(history.UpdatedAt, location)
	if len(targets) == 0 {
		t.Fatal("weeks after the last refresh should be fetched")
	}
	for _, target := range targets {
		if !target.After(lastRefresh) {
			t.Fatalf("previously scanned historical gap was retried: %s", target.Format(etfRuleRuntimeTimestampDateLayout))
		}
	}

	history.UpdatedAt = ""
	uncachedTargets := dividendLowVolHistoryTargets(now, now.AddDate(-dividendLowVolHistoryYears, 0, 0), history)
	if len(uncachedTargets) <= len(targets) {
		t.Fatalf("history without a completed refresh should retry gaps: cached=%d uncached=%d", len(targets), len(uncachedTargets))
	}
}
