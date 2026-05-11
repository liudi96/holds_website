package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	decisionLogLimit          = 500
	defaultSafetyMarginTarget = 0.25
	runtimeQuotesFile         = "data/runtime/quotes.json"
)

type AppState struct {
	TotalCapital float64            `json:"totalCapital"`
	Cash         float64            `json:"cash"`
	FX           map[string]float64 `json:"fx"`
	Trades       []Trade            `json:"trades"`
	DecisionLogs []DecisionLog      `json:"decisionLogs"`
	Holdings     []Holding          `json:"holdings"`
	Plan         []PlanItem         `json:"plan"`
	Candidates   []Candidate        `json:"candidates"`
	Rules        []Rule             `json:"rules"`
}

type Trade struct {
	ID           int64   `json:"id"`
	Date         string  `json:"date"`
	Symbol       string  `json:"symbol"`
	Name         string  `json:"name"`
	Side         string  `json:"side"`
	Shares       float64 `json:"shares"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	CurrentPrice float64 `json:"currentPrice"`
}

type DecisionLog struct {
	ID         int64    `json:"id"`
	Date       string   `json:"date"`
	Type       string   `json:"type"`
	Symbol     string   `json:"symbol,omitempty"`
	Name       string   `json:"name,omitempty"`
	Price      *float64 `json:"price,omitempty"`
	Currency   string   `json:"currency,omitempty"`
	Decision   string   `json:"decision"`
	Discipline string   `json:"discipline"`
	Detail     string   `json:"detail,omitempty"`
}

type Holding struct {
	Symbol              string              `json:"symbol"`
	Name                string              `json:"name"`
	Shares              float64             `json:"shares"`
	Cost                float64             `json:"cost"`
	CurrentPrice        float64             `json:"currentPrice"`
	PreviousClose       float64             `json:"previousClose"`
	MarketCap           *float64            `json:"marketCap,omitempty"`
	MarketCapCurrency   string              `json:"marketCapCurrency,omitempty"`
	CurrentPriceDate    string              `json:"currentPriceDate"`
	PreviousCloseDate   string              `json:"previousCloseDate"`
	Action              string              `json:"action"`
	Status              string              `json:"status"`
	MarginOfSafety      *float64            `json:"marginOfSafety"`
	QualityScore        *float64            `json:"qualityScore"`
	Risk                string              `json:"risk"`
	Industry            string              `json:"industry"`
	Currency            string              `json:"currency"`
	IntrinsicValue      *float64            `json:"intrinsicValue"`
	FairValueRange      string              `json:"fairValueRange"`
	TargetBuyPrice      *float64            `json:"targetBuyPrice"`
	PriceLevels         *PriceLevels        `json:"priceLevels,omitempty"`
	ValuationConfidence string              `json:"valuationConfidence,omitempty"`
	BusinessModel       *float64            `json:"businessModel"`
	Moat                *float64            `json:"moat"`
	Governance          *float64            `json:"governance"`
	FinancialQuality    *float64            `json:"financialQuality"`
	UpdatedAt           string              `json:"updatedAt"`
	Notes               string              `json:"notes"`
	KillCriteria        json.RawMessage     `json:"killCriteria,omitempty"`
	Reports             []Report            `json:"reports,omitempty"`
	Dividend            *Dividend           `json:"dividend,omitempty"`
	NetCash             *NetCashProfile     `json:"netCash,omitempty"`
	OwnerCashFlowAudit  *OwnerCashFlowAudit `json:"ownerCashFlowAudit,omitempty"`
	Financials          json.RawMessage     `json:"financials,omitempty"`
}

type PriceLevels struct {
	WatchPrice         *float64 `json:"watchPrice,omitempty"`
	InitialBuyPrice    *float64 `json:"initialBuyPrice,omitempty"`
	AggressiveBuyPrice *float64 `json:"aggressiveBuyPrice,omitempty"`
}

type Dividend struct {
	FiscalYear           string   `json:"fiscalYear,omitempty"`
	DividendPerShare     *float64 `json:"dividendPerShare,omitempty"`
	DividendCurrency     string   `json:"dividendCurrency,omitempty"`
	CashDividendTotal    *float64 `json:"cashDividendTotal,omitempty"`
	CashDividendCurrency string   `json:"cashDividendCurrency,omitempty"`
	BuybackAmount        *float64 `json:"buybackAmount,omitempty"`
	BuybackCurrency      string   `json:"buybackCurrency,omitempty"`
	DividendYield        *float64 `json:"dividendYield,omitempty"`
	PayoutRatio          *float64 `json:"payoutRatio,omitempty"`
	EstimatedAnnualCash  *float64 `json:"estimatedAnnualCash,omitempty"`
	Reliability          string   `json:"reliability,omitempty"`
	ForecastFiscalYear   string   `json:"forecastFiscalYear,omitempty"`
	ForecastPerShare     *float64 `json:"forecastPerShare,omitempty"`
	ForecastCurrency     string   `json:"forecastCurrency,omitempty"`
	ForecastYield        *float64 `json:"forecastYield,omitempty"`
}

type NetCashProfile struct {
	CashAndShortInvestments *float64 `json:"cashAndShortInvestments,omitempty"`
	InterestBearingDebt     *float64 `json:"interestBearingDebt,omitempty"`
	NetCash                 *float64 `json:"netCash,omitempty"`
	Currency                string   `json:"currency,omitempty"`
	Haircut                 *float64 `json:"haircut,omitempty"`
	HaircutReason           string   `json:"haircutReason,omitempty"`
	AdjustedNetCash         *float64 `json:"adjustedNetCash,omitempty"`
	ExCashPE                *float64 `json:"exCashPe,omitempty"`
	ExCashPFCF              *float64 `json:"exCashPfcf,omitempty"`
	FCFYield                *float64 `json:"fcfYield,omitempty"`
	ShareholderFCF          *float64 `json:"shareholderFcf,omitempty"`
	ShareholderFCFCurrency  string   `json:"shareholderFcfCurrency,omitempty"`
	ShareholderFCFBasis     string   `json:"shareholderFcfBasis,omitempty"`
	ConsolidatedFCF         *float64 `json:"consolidatedFcf,omitempty"`
	MinorityFCFAdjustment   *float64 `json:"minorityFcfAdjustment,omitempty"`
	FCFPositiveYears        *int     `json:"fcfPositiveYears,omitempty"`
	Note                    string   `json:"note,omitempty"`
}

type OwnerCashFlowAudit struct {
	TenYearDemand                  OwnerAuditItem `json:"tenYearDemand,omitempty"`
	AssetDurability                OwnerAuditItem `json:"assetDurability,omitempty"`
	MaintenanceCapexLight          OwnerAuditItem `json:"maintenanceCapexLight,omitempty"`
	DividendFCFSupport             OwnerAuditItem `json:"dividendFcfSupport,omitempty"`
	DividendReinvestmentEfficiency OwnerAuditItem `json:"dividendReinvestmentEfficiency,omitempty"`
	RoeRoicDurability              OwnerAuditItem `json:"roeRoicDurability,omitempty"`
	ValuationSystemRisk            OwnerAuditItem `json:"valuationSystemRisk,omitempty"`
}

type OwnerAuditItem struct {
	Status string `json:"status,omitempty"`
	Note   string `json:"note,omitempty"`
}

type Report struct {
	Period string `json:"period"`
	Kind   string `json:"kind"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Source string `json:"source"`
	URL    string `json:"url"`
}

type PlanItem struct {
	Rank       int    `json:"rank"`
	Symbol     string `json:"symbol,omitempty"`
	Name       string `json:"name"`
	Priority   string `json:"priority"`
	Advice     string `json:"advice"`
	Discipline string `json:"discipline"`
}

type Candidate struct {
	Symbol              string              `json:"symbol"`
	Name                string              `json:"name"`
	Status              string              `json:"status"`
	Action              string              `json:"action"`
	CurrentPrice        float64             `json:"currentPrice"`
	PreviousClose       float64             `json:"previousClose"`
	MarketCap           *float64            `json:"marketCap,omitempty"`
	MarketCapCurrency   string              `json:"marketCapCurrency,omitempty"`
	CurrentPriceDate    string              `json:"currentPriceDate"`
	PreviousCloseDate   string              `json:"previousCloseDate"`
	MarginOfSafety      *float64            `json:"marginOfSafety"`
	QualityScore        *float64            `json:"qualityScore"`
	Risk                string              `json:"risk"`
	Industry            string              `json:"industry"`
	Currency            string              `json:"currency"`
	IntrinsicValue      *float64            `json:"intrinsicValue"`
	FairValueRange      string              `json:"fairValueRange"`
	TargetBuyPrice      *float64            `json:"targetBuyPrice"`
	PriceLevels         *PriceLevels        `json:"priceLevels,omitempty"`
	ValuationConfidence string              `json:"valuationConfidence,omitempty"`
	BusinessModel       *float64            `json:"businessModel"`
	Moat                *float64            `json:"moat"`
	Governance          *float64            `json:"governance"`
	FinancialQuality    *float64            `json:"financialQuality"`
	UpdatedAt           string              `json:"updatedAt"`
	Notes               string              `json:"notes"`
	KillCriteria        json.RawMessage     `json:"killCriteria,omitempty"`
	Reports             []Report            `json:"reports,omitempty"`
	Dividend            *Dividend           `json:"dividend,omitempty"`
	NetCash             *NetCashProfile     `json:"netCash,omitempty"`
	OwnerCashFlowAudit  *OwnerCashFlowAudit `json:"ownerCashFlowAudit,omitempty"`
	Financials          json.RawMessage     `json:"financials,omitempty"`
}

type Rule struct {
	Dimension string  `json:"dimension"`
	Score     float64 `json:"score"`
	Standard  string  `json:"standard"`
}

type yahooChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency         string `json:"currency"`
				ExchangeTimezone string `json:"exchangeTimezoneName"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
			Events struct {
				Dividends map[string]struct {
					Amount float64 `json:"amount"`
					Date   int64   `json:"date"`
				} `json:"dividends"`
			} `json:"events"`
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}

func main() {
	dataPath := flag.String("data", "data/portfolio.json", "portfolio JSON file to read")
	quotesPath := flag.String("quotes", runtimeQuotesFile, "runtime quote JSON file to update")
	dryRun := flag.Bool("dry-run", false, "print updates without writing the file")
	flag.Parse()

	state, err := loadState(*dataPath)
	if err != nil {
		fail(err)
	}
	if err := mergeRuntimeQuotes(&state, *quotesPath); err != nil {
		fail(err)
	}

	client := &http.Client{Timeout: 12 * time.Second}
	now := time.Now().Format("2006-01-02 15:04:05")
	updated := 0
	quoteRecords := map[string]RuntimeQuote{}
	cache := make(map[string]quote)
	fallbackCache, fallbackErr := fetchFallbackQuotes(client, quoteSymbols(&state))

	for i := range state.Holdings {
		holding := &state.Holdings[i]
		if strings.TrimSpace(holding.Symbol) == "" {
			continue
		}

		quote, err := fetchQuoteCached(client, cache, fallbackCache, fallbackErr, holding.Symbol)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", holding.Symbol, err)
			continue
		}

		fmt.Printf("%s %s: %.4f -> %.4f (%s), yesterday close %.4f (%s) [%s]\n", holding.Symbol, holding.Name, holding.CurrentPrice, quote.Price, quote.PriceDate, quote.PreviousClose, quote.PreviousCloseDate, quote.SourceSymbol)
		holding.CurrentPrice = quote.Price
		holding.PreviousClose = quote.PreviousClose
		holding.CurrentPriceDate = quote.PriceDate
		holding.PreviousCloseDate = quote.PreviousCloseDate
		holding.MarginOfSafety = marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, holding.MarginOfSafety)
		applyDividendQuote(&holding.Dividend, quote, holding.Currency)
		if strings.TrimSpace(holding.Currency) == "" {
			holding.Currency = strings.ToUpper(strings.TrimSpace(quote.Currency))
		}
		holding.UpdatedAt = quoteUpdateLabel(now, quote)
		quoteRecords[normalizeSymbol(holding.Symbol)] = runtimeQuoteFromQuote(holding.Symbol, quote, now)
		updated++
	}

	for i := range state.Candidates {
		candidate := &state.Candidates[i]
		if strings.TrimSpace(candidate.Symbol) == "" {
			continue
		}

		quote, err := fetchQuoteCached(client, cache, fallbackCache, fallbackErr, candidate.Symbol)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", candidate.Symbol, err)
			continue
		}

		fmt.Printf("%s %s: %.4f -> %.4f (%s), yesterday close %.4f (%s) [%s]\n", candidate.Symbol, candidate.Name, candidate.CurrentPrice, quote.Price, quote.PriceDate, quote.PreviousClose, quote.PreviousCloseDate, quote.SourceSymbol)
		candidate.CurrentPrice = quote.Price
		candidate.PreviousClose = quote.PreviousClose
		candidate.CurrentPriceDate = quote.PriceDate
		candidate.PreviousCloseDate = quote.PreviousCloseDate
		candidate.MarginOfSafety = marginOfSafetyFromPrice(candidate.IntrinsicValue, candidate.CurrentPrice, candidate.MarginOfSafety)
		applyDividendQuote(&candidate.Dividend, quote, candidate.Currency)
		if strings.TrimSpace(candidate.Currency) == "" {
			candidate.Currency = strings.ToUpper(strings.TrimSpace(quote.Currency))
		}
		candidate.UpdatedAt = quoteUpdateLabel(now, quote)
		quoteRecords[normalizeSymbol(candidate.Symbol)] = runtimeQuoteFromQuote(candidate.Symbol, quote, now)
		updated++
	}

	if updated == 0 {
		fail(errors.New("no quotes were updated"))
	}

	if *dryRun {
		fmt.Printf("dry run: %d quote records would be updated in %s\n", updated, *quotesPath)
		return
	}

	if err := saveRuntimeQuoteRecords(*quotesPath, runtimeQuoteList(quoteRecords), now); err != nil {
		fail(err)
	}
	fmt.Printf("updated %d quote records in %s\n", updated, *quotesPath)
}

type quote struct {
	Price              float64
	PreviousClose      float64
	PriceDate          string
	PreviousCloseDate  string
	Currency           string
	SourceSymbol       string
	SourceName         string
	DividendPerShare   *float64
	DividendCurrency   string
	DividendFiscalYear string
}

func fetchQuote(client *http.Client, symbol string, fallbackCache map[string]quote, fallbackErr error) (quote, error) {
	quote, yahooErr := fetchYahooQuote(client, symbol)
	if yahooErr == nil {
		return quote, nil
	}

	if quote, ok := fallbackCache[normalizeSymbol(symbol)]; ok {
		return quote, nil
	}
	if fallbackErr != nil {
		return quote, fmt.Errorf("yahoo: %v; fallback: %v", yahooErr, fallbackErr)
	}

	return quote, fmt.Errorf("yahoo: %v; fallback: no quote for %s", yahooErr, symbol)
}

func quoteUpdateLabel(updateLabel string, quote quote) string {
	sourceName := firstNonEmpty(quote.SourceName, "Yahoo Finance 日线收盘价")
	return fmt.Sprintf("%s；行情源 %s；代码 %s；币种 %s；日期 %s/%s", updateLabel, sourceName, quote.SourceSymbol, quote.Currency, quote.PreviousCloseDate, quote.PriceDate)
}

func applyDividendQuote(current **Dividend, quote quote, fallbackCurrency string) {
	if quote.DividendPerShare == nil || *quote.DividendPerShare <= 0 {
		return
	}
	if *current == nil {
		*current = &Dividend{}
	}
	dividend := *current
	if dividend.CashDividendTotal != nil && *dividend.CashDividendTotal > 0 {
		if strings.TrimSpace(dividend.DividendCurrency) == "" {
			dividend.DividendCurrency = strings.ToUpper(firstNonEmpty(quote.DividendCurrency, quote.Currency, fallbackCurrency))
		}
		dividend.DividendYield = nil
		dividend.EstimatedAnnualCash = nil
		return
	}
	dividend.FiscalYear = firstNonEmpty(quote.DividendFiscalYear, dividend.FiscalYear)
	dividend.DividendPerShare = quote.DividendPerShare
	dividend.DividendCurrency = strings.ToUpper(firstNonEmpty(quote.DividendCurrency, quote.Currency, fallbackCurrency))
	dividend.DividendYield = nil
	dividend.EstimatedAnnualCash = nil
}

func fetchYahooQuote(client *http.Client, symbol string) (quote, error) {
	sourceSymbol := yahooSymbol(symbol)
	endpoint := "https://query1.finance.yahoo.com/v8/finance/chart/" + url.PathEscape(sourceSymbol) + "?range=2y&interval=1d&events=div"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return quote{}, err
	}
	req.Header.Set("User-Agent", "holds-website quote updater")

	resp, err := client.Do(req)
	if err != nil {
		return quote{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return quote{}, fmt.Errorf("quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload yahooChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return quote{}, err
	}
	if payload.Chart.Error != nil || len(payload.Chart.Result) == 0 {
		return quote{}, errors.New("empty quote response")
	}

	result := payload.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return quote{}, errors.New("missing close series")
	}

	location := loadLocation(result.Meta.ExchangeTimezone)
	closes := result.Indicators.Quote[0].Close
	validCloses := make([]dailyClose, 0, len(closes))
	for i, closePrice := range closes {
		if closePrice > 0 {
			validCloses = append(validCloses, dailyClose{
				Price: closePrice,
				Date:  closeDate(result.Timestamp, i, location),
			})
		}
	}
	if len(validCloses) < 2 {
		return quote{}, errors.New("need at least two daily close prices")
	}

	priceClose := validCloses[len(validCloses)-1]
	previousClose := validCloses[len(validCloses)-2]
	dividendPerShare, dividendFiscalYear := trailingDividendFromEvents(result.Events.Dividends, priceClose.Date, location)
	return quote{
		Price:              priceClose.Price,
		PreviousClose:      previousClose.Price,
		PriceDate:          priceClose.Date,
		PreviousCloseDate:  previousClose.Date,
		Currency:           result.Meta.Currency,
		SourceSymbol:       sourceSymbol,
		SourceName:         "Yahoo Finance 日线收盘价",
		DividendPerShare:   dividendPerShare,
		DividendCurrency:   result.Meta.Currency,
		DividendFiscalYear: dividendFiscalYear,
	}, nil
}

func trailingDividendFromEvents(events map[string]struct {
	Amount float64 `json:"amount"`
	Date   int64   `json:"date"`
}, priceDate string, location *time.Location) (*float64, string) {
	if len(events) == 0 {
		return nil, ""
	}
	reference, err := time.ParseInLocation("2006-01-02", priceDate, location)
	if err != nil {
		reference = time.Now().In(location)
	}
	cutoff := reference.AddDate(-1, 0, 0)
	total := 0.0
	latest := time.Time{}
	for _, event := range events {
		if event.Amount <= 0 || event.Date <= 0 {
			continue
		}
		eventDate := time.Unix(event.Date, 0).In(location)
		if eventDate.Before(cutoff) || eventDate.After(reference.AddDate(0, 0, 1)) {
			continue
		}
		total += event.Amount
		if eventDate.After(latest) {
			latest = eventDate
		}
	}
	if total <= 0 {
		return nil, ""
	}
	labelDate := reference
	if !latest.IsZero() {
		labelDate = latest
	}
	return &total, "TTM " + labelDate.Format("2006-01-02")
}

func fetchEastmoneyQuote(client *http.Client, symbol string) (quote, error) {
	sourceSymbol, scale, err := eastmoneySymbol(symbol)
	if err != nil {
		return quote{}, err
	}

	endpoint := "https://push2.eastmoney.com/api/qt/stock/get?secid=" + url.QueryEscape(sourceSymbol) + "&fields=f43,f57,f58,f60,f86,f107"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return quote{}, err
	}
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Referer", "https://quote.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website quote updater)")

	resp, err := client.Do(req)
	if err != nil {
		return quote{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return quote{}, fmt.Errorf("quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		RC   int            `json:"rc"`
		Data map[string]any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return quote{}, err
	}
	if payload.RC != 0 || len(payload.Data) == 0 {
		return quote{}, errors.New("empty quote response")
	}

	rawPrice, err := numberField(payload.Data, "f43")
	if err != nil {
		return quote{}, err
	}
	rawPreviousClose, err := numberField(payload.Data, "f60")
	if err != nil {
		return quote{}, err
	}
	if rawPrice <= 0 || rawPreviousClose <= 0 {
		return quote{}, errors.New("missing price series")
	}

	priceDate := eastmoneyQuoteDate(payload.Data)
	return quote{
		Price:             rawPrice / scale,
		PreviousClose:     rawPreviousClose / scale,
		PriceDate:         priceDate,
		PreviousCloseDate: priceDate,
		Currency:          currencyForSymbol(symbol),
		SourceSymbol:      sourceSymbol,
		SourceName:        "东方财富实时行情",
	}, nil
}

func fetchFallbackQuotes(client *http.Client, symbols []string) (map[string]quote, error) {
	quotes, err := fetchTencentQuotes(client, symbols)
	if err == nil {
		return quotes, nil
	}

	eastmoneyQuotes, eastmoneyErr := fetchEastmoneyQuotes(client, symbols)
	if eastmoneyErr == nil {
		return eastmoneyQuotes, nil
	}

	return nil, fmt.Errorf("tencent: %v; eastmoney: %v", err, eastmoneyErr)
}

func fetchTencentQuotes(client *http.Client, symbols []string) (map[string]quote, error) {
	querySymbols := []string{}
	normalizedByQuery := map[string]string{}
	for _, symbol := range symbols {
		querySymbol, normalized, err := tencentSymbol(symbol)
		if err != nil {
			continue
		}
		querySymbols = append(querySymbols, querySymbol)
		normalizedByQuery[querySymbol] = normalized
	}
	if len(querySymbols) == 0 {
		return map[string]quote{}, errors.New("no supported tencent symbols")
	}

	endpoint := "http://qt.gtimg.cn/q=" + strings.Join(querySymbols, ",")
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/plain,*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website quote updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	quotes := map[string]quote{}
	for _, line := range strings.Split(string(body), ";") {
		querySymbol, fields, ok := parseTencentLine(line)
		if !ok {
			continue
		}
		normalized := normalizedByQuery[querySymbol]
		if normalized == "" || len(fields) <= 30 {
			continue
		}
		price, err := strconv.ParseFloat(strings.TrimSpace(fields[3]), 64)
		if err != nil || price <= 0 {
			continue
		}
		previousClose, err := strconv.ParseFloat(strings.TrimSpace(fields[4]), 64)
		if err != nil || previousClose <= 0 {
			continue
		}
		priceDate := tencentQuoteDate(fields[30])
		quotes[normalized] = quote{
			Price:             price,
			PreviousClose:     previousClose,
			PriceDate:         priceDate,
			PreviousCloseDate: priceDate,
			Currency:          currencyForSymbol(normalized),
			SourceSymbol:      querySymbol,
			SourceName:        "腾讯实时行情",
		}
	}
	if len(quotes) == 0 {
		return nil, errors.New("empty quote response")
	}
	return quotes, nil
}

func parseTencentLine(line string) (string, []string, bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "v_") {
		return "", nil, false
	}
	equalIndex := strings.Index(line, "=")
	if equalIndex <= 2 {
		return "", nil, false
	}
	querySymbol := strings.TrimSpace(strings.TrimPrefix(line[:equalIndex], "v_"))
	payload := strings.TrimSpace(line[equalIndex+1:])
	payload = strings.Trim(payload, "\"")
	if querySymbol == "" || payload == "" {
		return "", nil, false
	}
	return querySymbol, strings.Split(payload, "~"), true
}

func tencentQuoteDate(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= len("2006/01/02") && strings.Contains(value[:len("2006/01/02")], "/") {
		return strings.ReplaceAll(value[:len("2006/01/02")], "/", "-")
	}
	if len(value) >= len("20060102") {
		if parsed, err := time.Parse("20060102", value[:len("20060102")]); err == nil {
			return parsed.Format("2006-01-02")
		}
	}
	return time.Now().Format("2006-01-02")
}

func tencentSymbol(symbol string) (string, string, error) {
	symbol = normalizeSymbol(symbol)
	switch {
	case strings.HasSuffix(symbol, ".SH"):
		code := strings.TrimSuffix(symbol, ".SH")
		return "sh" + code, symbol, nil
	case strings.HasSuffix(symbol, ".SZ"):
		code := strings.TrimSuffix(symbol, ".SZ")
		return "sz" + code, symbol, nil
	case strings.HasSuffix(symbol, ".HK"):
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconv.Atoi(code); err == nil {
			code = fmt.Sprintf("%05d", value)
		}
		return "hk" + code, symbol, nil
	default:
		return "", "", fmt.Errorf("unsupported tencent symbol: %s", symbol)
	}
}

func fetchEastmoneyQuotes(client *http.Client, symbols []string) (map[string]quote, error) {
	secIDs := []string{}
	for _, symbol := range symbols {
		secID, _, err := eastmoneySymbol(symbol)
		if err == nil {
			secIDs = append(secIDs, secID)
		}
	}
	if len(secIDs) == 0 {
		return map[string]quote{}, errors.New("no supported eastmoney symbols")
	}

	endpoint := "https://push2.eastmoney.com/api/qt/ulist.np/get?secids=" + strings.Join(secIDs, ",") + "&fields=f2,f12,f13,f14,f18,f124"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Referer", "https://quote.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website quote updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		RC   int `json:"rc"`
		Data struct {
			Diff []map[string]any `json:"diff"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.RC != 0 || len(payload.Data.Diff) == 0 {
		return nil, errors.New("empty quote response")
	}

	quotes := make(map[string]quote, len(payload.Data.Diff))
	for _, item := range payload.Data.Diff {
		normalized, sourceSymbol, scale, err := eastmoneyDiffSymbol(item)
		if err != nil {
			continue
		}
		rawPrice, err := numberField(item, "f2")
		if err != nil {
			continue
		}
		rawPreviousClose, err := numberField(item, "f18")
		if err != nil {
			continue
		}
		if rawPrice <= 0 || rawPreviousClose <= 0 {
			continue
		}
		priceDate := eastmoneyQuoteDate(item)
		quotes[normalized] = quote{
			Price:             rawPrice / scale,
			PreviousClose:     rawPreviousClose / scale,
			PriceDate:         priceDate,
			PreviousCloseDate: priceDate,
			Currency:          currencyForSymbol(normalized),
			SourceSymbol:      sourceSymbol,
			SourceName:        "东方财富实时行情",
		}
	}
	if len(quotes) == 0 {
		return nil, errors.New("empty quote response")
	}
	return quotes, nil
}

func numberField(data map[string]any, key string) (float64, error) {
	value, ok := data[key]
	if !ok {
		return 0, fmt.Errorf("missing %s", key)
	}
	switch typed := value.(type) {
	case float64:
		return typed, nil
	case string:
		number, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid %s", key)
		}
		return number, nil
	default:
		return 0, fmt.Errorf("invalid %s", key)
	}
}

func eastmoneyQuoteDate(data map[string]any) string {
	for _, key := range []string{"f86", "f124"} {
		value, err := numberField(data, key)
		if err == nil && value > 0 {
			return time.Unix(int64(value), 0).In(loadLocation("Asia/Shanghai")).Format("2006-01-02")
		}
	}
	return time.Now().Format("2006-01-02")
}

func eastmoneyDiffSymbol(data map[string]any) (string, string, float64, error) {
	rawMarket, err := numberField(data, "f13")
	if err != nil {
		return "", "", 0, err
	}
	code, err := stringField(data, "f12")
	if err != nil {
		return "", "", 0, err
	}

	market := int(rawMarket)
	sourceSymbol := fmt.Sprintf("%d.%s", market, code)
	switch market {
	case 0:
		return normalizeSymbol(code + ".SZ"), sourceSymbol, 100, nil
	case 1:
		return normalizeSymbol(code + ".SH"), sourceSymbol, 100, nil
	case 116:
		if value, err := strconv.Atoi(code); err == nil {
			return normalizeSymbol(fmt.Sprintf("%d.HK", value)), sourceSymbol, 1000, nil
		}
		return normalizeSymbol(code + ".HK"), sourceSymbol, 1000, nil
	default:
		return "", "", 0, fmt.Errorf("unsupported eastmoney market: %d", market)
	}
}

func stringField(data map[string]any, key string) (string, error) {
	value, ok := data[key]
	if !ok {
		return "", fmt.Errorf("missing %s", key)
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed), nil
	case float64:
		return strconv.FormatInt(int64(typed), 10), nil
	default:
		return "", fmt.Errorf("invalid %s", key)
	}
}

func eastmoneySymbol(symbol string) (string, float64, error) {
	symbol = normalizeSymbol(symbol)
	switch {
	case strings.HasSuffix(symbol, ".SH"):
		return "1." + strings.TrimSuffix(symbol, ".SH"), 100, nil
	case strings.HasSuffix(symbol, ".SZ"):
		return "0." + strings.TrimSuffix(symbol, ".SZ"), 100, nil
	case strings.HasSuffix(symbol, ".HK"):
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconv.Atoi(code); err == nil {
			code = fmt.Sprintf("%05d", value)
		}
		return "116." + code, 1000, nil
	default:
		return "", 0, fmt.Errorf("unsupported eastmoney symbol: %s", symbol)
	}
}

func currencyForSymbol(symbol string) string {
	symbol = normalizeSymbol(symbol)
	switch {
	case strings.HasSuffix(symbol, ".HK"):
		return "HKD"
	case strings.HasSuffix(symbol, ".SH"), strings.HasSuffix(symbol, ".SZ"):
		return "CNY"
	default:
		return ""
	}
}

func quoteSymbols(state *AppState) []string {
	symbols := []string{}
	seen := map[string]bool{}
	for _, holding := range state.Holdings {
		normalized := normalizeSymbol(holding.Symbol)
		if normalized != "" && !seen[normalized] {
			symbols = append(symbols, holding.Symbol)
			seen[normalized] = true
		}
	}
	for _, candidate := range state.Candidates {
		normalized := normalizeSymbol(candidate.Symbol)
		if normalized != "" && !seen[normalized] {
			symbols = append(symbols, candidate.Symbol)
			seen[normalized] = true
		}
	}
	return symbols
}

func fetchQuoteCached(client *http.Client, cache map[string]quote, fallbackCache map[string]quote, fallbackErr error, symbol string) (quote, error) {
	normalized := strings.ToUpper(strings.TrimSpace(symbol))
	if cached, ok := cache[normalized]; ok {
		return cached, nil
	}
	quote, err := fetchQuote(client, normalized, fallbackCache, fallbackErr)
	if err != nil {
		return quote, err
	}
	cache[normalized] = quote
	return quote, nil
}

type dailyClose struct {
	Price float64
	Date  string
}

func loadLocation(name string) *time.Location {
	if strings.TrimSpace(name) == "" {
		return time.UTC
	}
	location, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return location
}

func closeDate(timestamps []int64, index int, location *time.Location) string {
	if index < 0 || index >= len(timestamps) || timestamps[index] <= 0 {
		return ""
	}
	return time.Unix(timestamps[index], 0).In(location).Format("2006-01-02")
}

func yahooSymbol(symbol string) string {
	symbol = normalizeSymbol(symbol)
	if strings.HasSuffix(symbol, ".SH") {
		return strings.TrimSuffix(symbol, ".SH") + ".SS"
	}
	return symbol
}

func marginOfSafetyFromPrice(intrinsicValue *float64, currentPrice float64, fallback *float64) *float64 {
	if intrinsicValue == nil || *intrinsicValue <= 0 || currentPrice <= 0 {
		return fallback
	}
	value := (*intrinsicValue - currentPrice) / *intrinsicValue
	return &value
}

func appendQuoteDecisionLog(state *AppState, symbol string, name string, currency string, beforePrice float64, currentPrice float64, intrinsicValue *float64, currentDate string, previousDate string, now string) {
	trigger := quoteTriggerText(beforePrice, currentPrice, intrinsicValue)
	if strings.TrimSpace(trigger) == "" {
		return
	}
	_, _, _, decision, discipline := decisionLogContext(state, symbol)
	appendDecisionLog(state, DecisionLog{
		Date:       now,
		Type:       "quote",
		Symbol:     symbol,
		Name:       name,
		Price:      pricePointer(currentPrice),
		Currency:   currency,
		Decision:   firstNonEmpty(decision, "行情触发"),
		Discipline: firstNonEmpty(discipline, "只在纪律区间变化时记录行情日志"),
		Detail:     fmt.Sprintf("%s；现价 %s %.4f；前值 %.4f；今收 %s；昨收 %s", trigger, strings.ToUpper(currency), currentPrice, beforePrice, firstNonEmpty(currentDate, "未知"), firstNonEmpty(previousDate, "未知")),
	})
}

func quoteTriggerText(beforePrice float64, currentPrice float64, intrinsicValue *float64) string {
	beforeZone := quotePriceZone(beforePrice, intrinsicValue)
	currentZone := quotePriceZone(currentPrice, intrinsicValue)
	if beforeZone != "" && currentZone != "" && beforeZone != currentZone {
		return fmt.Sprintf("进入/离开关键区间：%s -> %s", beforeZone, currentZone)
	}
	beforeMargin := quoteMarginZone(beforePrice, intrinsicValue)
	currentMargin := quoteMarginZone(currentPrice, intrinsicValue)
	if beforeMargin != "" && currentMargin != "" && beforeMargin != currentMargin && (currentMargin == "安全边际达标" || currentMargin == "高于内在价值") {
		return fmt.Sprintf("安全边际跨区：%s -> %s", beforeMargin, currentMargin)
	}
	return ""
}

func quotePriceZone(price float64, intrinsicValue *float64) string {
	if price <= 0 || intrinsicValue == nil || *intrinsicValue <= 0 {
		return ""
	}
	initialBuyPrice := *intrinsicValue * (1 - defaultSafetyMarginTarget)
	watchPrice := initialBuyPrice * 1.05
	aggressiveBuyPrice := initialBuyPrice * 0.9
	switch {
	case price <= aggressiveBuyPrice:
		return "重仓区"
	case price <= initialBuyPrice:
		return "首买区"
	case price <= watchPrice:
		return "观察区"
	default:
		return "等待区"
	}
}

func quoteMarginZone(price float64, intrinsicValue *float64) string {
	if price <= 0 || intrinsicValue == nil || *intrinsicValue <= 0 {
		return ""
	}
	margin := (*intrinsicValue - price) / *intrinsicValue
	switch {
	case margin >= defaultSafetyMarginTarget:
		return "安全边际达标"
	case margin < 0:
		return "高于内在价值"
	default:
		return "安全边际不足"
	}
}

func appendDecisionLog(state *AppState, entry DecisionLog) {
	entry.Type = strings.TrimSpace(entry.Type)
	if entry.Type == "" {
		entry.Type = "event"
	}
	entry.Symbol = normalizeSymbol(entry.Symbol)
	entry.Name = strings.TrimSpace(entry.Name)
	entry.Currency = normalizeSymbol(entry.Currency)
	entry.Decision = strings.TrimSpace(entry.Decision)
	entry.Discipline = strings.TrimSpace(entry.Discipline)
	entry.Detail = strings.TrimSpace(entry.Detail)
	if entry.ID == 0 {
		entry.ID = time.Now().UnixNano()
	}
	if strings.TrimSpace(entry.Date) == "" {
		entry.Date = time.Now().Format("2006-01-02 15:04:05")
	}

	state.DecisionLogs = append(state.DecisionLogs, entry)
	if len(state.DecisionLogs) > decisionLogLimit {
		state.DecisionLogs = state.DecisionLogs[len(state.DecisionLogs)-decisionLogLimit:]
	}
}

func decisionLogContext(state *AppState, symbol string) (string, *float64, string, string, string) {
	normalizedSymbol := normalizeSymbol(symbol)
	for i := range state.Holdings {
		holding := state.Holdings[i]
		if normalizeSymbol(holding.Symbol) != normalizedSymbol {
			continue
		}
		plan := findPlanForDecisionLog(state, holding.Symbol, holding.Name)
		return holding.Name, pricePointer(holding.CurrentPrice), holding.Currency, firstNonEmpty(holding.Action, holding.Status), firstNonEmpty(planDiscipline(plan), holding.Status)
	}

	for i := range state.Candidates {
		candidate := state.Candidates[i]
		if normalizeSymbol(candidate.Symbol) != normalizedSymbol {
			continue
		}
		plan := findPlanForDecisionLog(state, candidate.Symbol, candidate.Name)
		return candidate.Name, pricePointer(candidate.CurrentPrice), candidate.Currency, firstNonEmpty(candidate.Action, candidate.Status), firstNonEmpty(planDiscipline(plan), candidate.Status)
	}

	plan := findPlanForDecisionLog(state, symbol, "")
	return "", nil, "", "", planDiscipline(plan)
}

func findPlanForDecisionLog(state *AppState, symbol string, name string) *PlanItem {
	normalizedSymbol := normalizeSymbol(symbol)
	normalizedName := strings.TrimSpace(name)
	for i := range state.Plan {
		itemSymbol := normalizeSymbol(state.Plan[i].Symbol)
		if itemSymbol != "" && normalizedSymbol != "" && itemSymbol == normalizedSymbol {
			return &state.Plan[i]
		}
		itemName := strings.TrimSpace(state.Plan[i].Name)
		if itemName != "" && normalizedName != "" && (strings.EqualFold(itemName, normalizedName) || strings.Contains(normalizedName, itemName) || strings.Contains(itemName, normalizedName)) {
			return &state.Plan[i]
		}
	}
	return nil
}

func planDiscipline(plan *PlanItem) string {
	if plan == nil {
		return ""
	}
	return strings.TrimSpace(plan.Discipline)
}

func pricePointer(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if text := strings.TrimSpace(value); text != "" {
			return text
		}
	}
	return ""
}

func normalizeSymbol(symbol string) string {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if strings.HasSuffix(symbol, ".HK") {
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconv.Atoi(code); err == nil {
			return fmt.Sprintf("%04d.HK", value)
		}
	}
	return symbol
}

type RuntimeQuoteBook struct {
	UpdatedAt string                  `json:"updatedAt,omitempty"`
	Quotes    map[string]RuntimeQuote `json:"quotes"`
}

type RuntimeQuote struct {
	Symbol             string   `json:"symbol"`
	CurrentPrice       float64  `json:"currentPrice,omitempty"`
	PreviousClose      float64  `json:"previousClose,omitempty"`
	CurrentPriceDate   string   `json:"currentPriceDate,omitempty"`
	PreviousCloseDate  string   `json:"previousCloseDate,omitempty"`
	Currency           string   `json:"currency,omitempty"`
	SourceSymbol       string   `json:"sourceSymbol,omitempty"`
	SourceName         string   `json:"sourceName,omitempty"`
	UpdatedAt          string   `json:"updatedAt,omitempty"`
	DividendPerShare   *float64 `json:"dividendPerShare,omitempty"`
	DividendCurrency   string   `json:"dividendCurrency,omitempty"`
	DividendFiscalYear string   `json:"dividendFiscalYear,omitempty"`
}

func loadRuntimeQuoteBook(path string) (RuntimeQuoteBook, error) {
	book := RuntimeQuoteBook{Quotes: map[string]RuntimeQuote{}}
	body, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return book, nil
	}
	if err != nil {
		return book, err
	}
	if err := json.Unmarshal(body, &book); err != nil {
		return book, err
	}
	if book.Quotes == nil {
		book.Quotes = map[string]RuntimeQuote{}
	}
	normalized := make(map[string]RuntimeQuote, len(book.Quotes))
	for key, record := range book.Quotes {
		symbol := normalizeSymbol(firstNonEmpty(record.Symbol, key))
		if symbol == "" {
			continue
		}
		record.Symbol = symbol
		normalized[symbol] = record
	}
	book.Quotes = normalized
	return book, nil
}

func saveRuntimeQuoteRecords(path string, records []RuntimeQuote, updatedAt string) error {
	book, err := loadRuntimeQuoteBook(path)
	if err != nil {
		return err
	}
	if strings.TrimSpace(updatedAt) != "" {
		book.UpdatedAt = strings.TrimSpace(updatedAt)
	}
	for _, record := range records {
		symbol := normalizeSymbol(record.Symbol)
		if symbol == "" {
			continue
		}
		record.Symbol = symbol
		book.Quotes[symbol] = record
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(path, body, 0o644)
}

func mergeRuntimeQuotes(state *AppState, path string) error {
	book, err := loadRuntimeQuoteBook(path)
	if err != nil {
		return err
	}
	for i := range state.Holdings {
		record, ok := book.Quotes[normalizeSymbol(state.Holdings[i].Symbol)]
		if ok {
			applyRuntimeQuoteToHolding(&state.Holdings[i], record)
		}
	}
	for i := range state.Candidates {
		record, ok := book.Quotes[normalizeSymbol(state.Candidates[i].Symbol)]
		if ok {
			applyRuntimeQuoteToCandidate(&state.Candidates[i], record)
		}
	}
	return nil
}

func applyRuntimeQuoteToHolding(holding *Holding, record RuntimeQuote) {
	if record.CurrentPrice > 0 {
		holding.CurrentPrice = record.CurrentPrice
		holding.MarginOfSafety = marginOfSafetyFromPrice(holding.IntrinsicValue, holding.CurrentPrice, holding.MarginOfSafety)
	}
	if record.PreviousClose > 0 {
		holding.PreviousClose = record.PreviousClose
	}
	holding.CurrentPriceDate = firstNonEmpty(record.CurrentPriceDate, holding.CurrentPriceDate)
	holding.PreviousCloseDate = firstNonEmpty(record.PreviousCloseDate, holding.PreviousCloseDate)
	if strings.TrimSpace(holding.Currency) == "" {
		holding.Currency = strings.ToUpper(strings.TrimSpace(record.Currency))
	}
	applyDividendQuote(&holding.Dividend, runtimeRecordAsQuote(record), holding.Currency)
	if strings.TrimSpace(record.UpdatedAt) != "" {
		holding.UpdatedAt = record.UpdatedAt
	}
}

func applyRuntimeQuoteToCandidate(candidate *Candidate, record RuntimeQuote) {
	if record.CurrentPrice > 0 {
		candidate.CurrentPrice = record.CurrentPrice
		candidate.MarginOfSafety = marginOfSafetyFromPrice(candidate.IntrinsicValue, candidate.CurrentPrice, candidate.MarginOfSafety)
	}
	if record.PreviousClose > 0 {
		candidate.PreviousClose = record.PreviousClose
	}
	candidate.CurrentPriceDate = firstNonEmpty(record.CurrentPriceDate, candidate.CurrentPriceDate)
	candidate.PreviousCloseDate = firstNonEmpty(record.PreviousCloseDate, candidate.PreviousCloseDate)
	if strings.TrimSpace(candidate.Currency) == "" {
		candidate.Currency = strings.ToUpper(strings.TrimSpace(record.Currency))
	}
	applyDividendQuote(&candidate.Dividend, runtimeRecordAsQuote(record), candidate.Currency)
	if strings.TrimSpace(record.UpdatedAt) != "" {
		candidate.UpdatedAt = record.UpdatedAt
	}
}

func runtimeQuoteFromQuote(symbol string, quote quote, updateLabel string) RuntimeQuote {
	return RuntimeQuote{
		Symbol:             normalizeSymbol(symbol),
		CurrentPrice:       quote.Price,
		PreviousClose:      quote.PreviousClose,
		CurrentPriceDate:   quote.PriceDate,
		PreviousCloseDate:  quote.PreviousCloseDate,
		Currency:           strings.ToUpper(strings.TrimSpace(quote.Currency)),
		SourceSymbol:       quote.SourceSymbol,
		SourceName:         quote.SourceName,
		UpdatedAt:          quoteUpdateLabel(updateLabel, quote),
		DividendPerShare:   quote.DividendPerShare,
		DividendCurrency:   strings.ToUpper(strings.TrimSpace(quote.DividendCurrency)),
		DividendFiscalYear: quote.DividendFiscalYear,
	}
}

func runtimeRecordAsQuote(record RuntimeQuote) quote {
	return quote{
		Price:              record.CurrentPrice,
		PreviousClose:      record.PreviousClose,
		PriceDate:          record.CurrentPriceDate,
		PreviousCloseDate:  record.PreviousCloseDate,
		Currency:           record.Currency,
		SourceSymbol:       record.SourceSymbol,
		SourceName:         record.SourceName,
		DividendPerShare:   record.DividendPerShare,
		DividendCurrency:   record.DividendCurrency,
		DividendFiscalYear: record.DividendFiscalYear,
	}
}

func runtimeQuoteList(records map[string]RuntimeQuote) []RuntimeQuote {
	keys := make([]string, 0, len(records))
	for key := range records {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	list := make([]RuntimeQuote, 0, len(keys))
	for _, key := range keys {
		list = append(list, records[key])
	}
	return list
}

func loadState(path string) (AppState, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return AppState{}, err
	}

	var state AppState
	if err := json.Unmarshal(body, &state); err != nil {
		return AppState{}, err
	}
	return state, nil
}

func saveState(path string, state AppState) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	body = append(body, '\n')
	return os.WriteFile(path, body, 0o644)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
