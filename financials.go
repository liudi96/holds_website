package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const financialYearsLimit = 10

type FinancialUpdateResponse struct {
	Symbol     string      `json:"symbol"`
	TargetType string      `json:"targetType"`
	Financials *Financials `json:"financials"`
	State      AppState    `json:"state"`
}

type Financials struct {
	UpdatedAt    string              `json:"updatedAt,omitempty"`
	Source       string              `json:"source,omitempty"`
	SourceSymbol string              `json:"sourceSymbol,omitempty"`
	Currency     string              `json:"currency,omitempty"`
	Annual       []FinancialAnnual   `json:"annual,omitempty"`
	Valuation    *FinancialValuation `json:"valuation,omitempty"`
}

type FinancialAnnual struct {
	FiscalYear                 string   `json:"fiscalYear,omitempty"`
	ReportDate                 string   `json:"reportDate,omitempty"`
	ReportType                 string   `json:"reportType,omitempty"`
	Currency                   string   `json:"currency,omitempty"`
	Revenue                    *float64 `json:"revenue,omitempty"`
	RevenueYoY                 *float64 `json:"revenueYoY,omitempty"`
	NetProfit                  *float64 `json:"netProfit,omitempty"`
	NetProfitYoY               *float64 `json:"netProfitYoY,omitempty"`
	OperatingProfit            *float64 `json:"operatingProfit,omitempty"`
	OperatingCashFlow          *float64 `json:"operatingCashFlow,omitempty"`
	OperatingCashFlowToRevenue *float64 `json:"operatingCashFlowToRevenue,omitempty"`
	CapitalExpenditure         *float64 `json:"capitalExpenditure,omitempty"`
	FreeCashFlow               *float64 `json:"freeCashFlow,omitempty"`
	EPS                        *float64 `json:"eps,omitempty"`
	BookValuePerShare          *float64 `json:"bookValuePerShare,omitempty"`
	TotalAssets                *float64 `json:"totalAssets,omitempty"`
	TotalLiabilities           *float64 `json:"totalLiabilities,omitempty"`
	TotalEquity                *float64 `json:"totalEquity,omitempty"`
	DebtRatio                  *float64 `json:"debtRatio,omitempty"`
	ROE                        *float64 `json:"roe,omitempty"`
	ROIC                       *float64 `json:"roic,omitempty"`
	GrossMargin                *float64 `json:"grossMargin,omitempty"`
	OperatingMargin            *float64 `json:"operatingMargin,omitempty"`
	NetMargin                  *float64 `json:"netMargin,omitempty"`
	Inventory                  *float64 `json:"inventory,omitempty"`
	AccountsReceivable         *float64 `json:"accountsReceivable,omitempty"`
	InventoryTurnoverDays      *float64 `json:"inventoryTurnoverDays,omitempty"`
	ReceivableTurnoverDays     *float64 `json:"receivableTurnoverDays,omitempty"`
	DividendPaid               *float64 `json:"dividendPaid,omitempty"`
	PEAtCurrentPrice           *float64 `json:"peAtCurrentPrice,omitempty"`
	PBAtCurrentPrice           *float64 `json:"pbAtCurrentPrice,omitempty"`
	PriceToOperatingCashFlow   *float64 `json:"priceToOperatingCashFlow,omitempty"`
	PriceToFreeCashFlow        *float64 `json:"priceToFreeCashFlow,omitempty"`
}

type FinancialValuation struct {
	Date          string       `json:"date,omitempty"`
	Price         *float64     `json:"price,omitempty"`
	Currency      string       `json:"currency,omitempty"`
	PE            *float64     `json:"pe,omitempty"`
	PB            *float64     `json:"pb,omitempty"`
	PEG           *float64     `json:"peg,omitempty"`
	DividendYield *float64     `json:"dividendYield,omitempty"`
	PERange       *MetricRange `json:"peRange,omitempty"`
	PBRange       *MetricRange `json:"pbRange,omitempty"`
	SourceNote    string       `json:"sourceNote,omitempty"`
}

type MetricRange struct {
	Min    *float64 `json:"min,omitempty"`
	Median *float64 `json:"median,omitempty"`
	Max    *float64 `json:"max,omitempty"`
}

type financialTarget struct {
	Symbol       string
	Name         string
	Currency     string
	CurrentPrice float64
	MarketCap    *float64
	Dividend     *Dividend
}

func (s *Server) handleUpdateFinancials(w http.ResponseWriter, r *http.Request) {
	symbol := strings.TrimPrefix(r.URL.Path, "/api/financials/update/")
	symbol = normalizeSymbol(symbol)
	if symbol == "" {
		writeError(w, http.StatusBadRequest, "missing symbol")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := loadState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load state")
		return
	}

	target, targetType, apply := findFinancialTarget(&state, symbol)
	if apply == nil {
		writeError(w, http.StatusNotFound, "stock not found")
		return
	}

	financials, err := fetchFinancials(&http.Client{Timeout: 15 * time.Second}, target, time.Now())
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	apply(financials)
	appendDecisionLog(&state, DecisionLog{
		Type:       "financials",
		Symbol:     target.Symbol,
		Name:       target.Name,
		Currency:   financials.Currency,
		Decision:   "更新多年财务数据",
		Discipline: "三大师视角使用结构化财务数据复核质量、成长、杠杆和估值",
		Detail:     financialDecisionDetail(financials),
	})

	if err := saveState(state); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save state")
		return
	}

	s.state = state
	writeJSON(w, http.StatusOK, FinancialUpdateResponse{
		Symbol:     target.Symbol,
		TargetType: targetType,
		Financials: financials,
		State:      state,
	})
}

func findFinancialTarget(state *AppState, symbol string) (financialTarget, string, func(*Financials)) {
	normalized := normalizeSymbol(symbol)
	for i := range state.Holdings {
		if normalizeSymbol(state.Holdings[i].Symbol) != normalized {
			continue
		}
		target := financialTarget{
			Symbol:       state.Holdings[i].Symbol,
			Name:         state.Holdings[i].Name,
			Currency:     state.Holdings[i].Currency,
			CurrentPrice: state.Holdings[i].CurrentPrice,
			MarketCap:    state.Holdings[i].MarketCap,
			Dividend:     state.Holdings[i].Dividend,
		}
		return target, "holding", func(financials *Financials) {
			state.Holdings[i].Financials = financials
		}
	}

	for i := range state.Candidates {
		if normalizeSymbol(state.Candidates[i].Symbol) != normalized {
			continue
		}
		target := financialTarget{
			Symbol:       state.Candidates[i].Symbol,
			Name:         state.Candidates[i].Name,
			Currency:     state.Candidates[i].Currency,
			CurrentPrice: state.Candidates[i].CurrentPrice,
			MarketCap:    state.Candidates[i].MarketCap,
			Dividend:     state.Candidates[i].Dividend,
		}
		return target, "candidate", func(financials *Financials) {
			state.Candidates[i].Financials = financials
		}
	}

	return financialTarget{}, "", nil
}

func fetchFinancials(client *http.Client, target financialTarget, now time.Time) (*Financials, error) {
	symbol := normalizeSymbol(target.Symbol)
	switch {
	case strings.HasSuffix(symbol, ".HK"):
		return fetchHKFinancials(client, target, now)
	case strings.HasSuffix(symbol, ".SH"), strings.HasSuffix(symbol, ".SZ"):
		return fetchAshareFinancials(client, target, now)
	default:
		return nil, fmt.Errorf("unsupported financials symbol: %s", target.Symbol)
	}
}

func fetchAshareFinancials(client *http.Client, target financialTarget, now time.Time) (*Financials, error) {
	secucode, err := eastmoneyFinancialSecucode(target.Symbol)
	if err != nil {
		return nil, err
	}

	mainRows, err := fetchEastmoneyDataRows(client, "https://datacenter.eastmoney.com/securities/api/data/get", url.Values{
		"source": []string{"HSF10"},
		"client": []string{"APP"},
		"type":   []string{"RPT_F10_FINANCE_MAINFINADATA"},
		"sty":    []string{"APP_F10_MAINFINADATA"},
		"filter": []string{fmt.Sprintf(`(SECUCODE="%s")`, secucode)},
		"ps":     []string{"80"},
		"sr":     []string{"-1"},
		"st":     []string{"REPORT_DATE"},
	})
	if err != nil {
		return nil, err
	}

	cashRows, err := fetchEastmoneyDataRows(client, "https://datacenter.eastmoney.com/securities/api/data/get", url.Values{
		"source": []string{"HSF10"},
		"client": []string{"APP"},
		"type":   []string{"RPT_F10_FINANCE_GCASHFLOW"},
		"sty":    []string{"APP_F10_GCASHFLOW"},
		"filter": []string{fmt.Sprintf(`(SECUCODE="%s")`, secucode)},
		"ps":     []string{"80"},
		"sr":     []string{"-1"},
		"st":     []string{"REPORT_DATE"},
	})
	if err != nil {
		cashRows = nil
	}

	cashByDate := rowsByReportDate(cashRows)
	annual := make([]FinancialAnnual, 0, financialYearsLimit)
	for _, row := range mainRows {
		if !isAnnualFinancialRow(row) {
			continue
		}
		reportDate := cleanReportDate(firstNonEmpty(mapString(row, "REPORT_DATE"), mapString(row, "STD_REPORT_DATE")))
		cashRow := cashByDate[reportDate]
		item := FinancialAnnual{
			FiscalYear:                 firstNonEmpty(mapString(row, "REPORT_YEAR"), financialYear(reportDate)),
			ReportDate:                 reportDate,
			ReportType:                 firstNonEmpty(mapString(row, "REPORT_DATE_NAME"), mapString(row, "REPORT_TYPE")),
			Currency:                   firstNonEmpty(mapString(row, "CURRENCY"), target.Currency, currencyForSymbol(target.Symbol)),
			Revenue:                    optionalNumber(row, "TOTALOPERATEREVE"),
			RevenueYoY:                 optionalPercent(row, "TOTALOPERATEREVETZ"),
			NetProfit:                  optionalNumber(row, "PARENTNETPROFIT"),
			NetProfitYoY:               optionalPercent(row, "PARENTNETPROFITTZ"),
			OperatingCashFlow:          optionalNumber(cashRow, "NETCASH_OPERATE"),
			CapitalExpenditure:         absNumber(optionalNumber(cashRow, "CONSTRUCT_LONG_ASSET")),
			EPS:                        optionalNumber(row, "EPSJB"),
			BookValuePerShare:          optionalNumber(row, "BPS"),
			TotalLiabilities:           optionalNumber(row, "LIABILITY"),
			DebtRatio:                  optionalPercent(row, "ZCFZL"),
			ROE:                        optionalPercent(row, "ROEJQ"),
			ROIC:                       optionalPercent(row, "ROIC"),
			GrossMargin:                optionalPercent(row, "XSMLL"),
			OperatingMargin:            optionalPercent(row, "PER_EBIT"),
			NetMargin:                  optionalPercent(row, "XSJLL"),
			InventoryTurnoverDays:      optionalNumber(row, "CHZZTS"),
			ReceivableTurnoverDays:     optionalNumber(row, "YSZKZZTS"),
			DividendPaid:               optionalNumber(cashRow, "ASSIGN_DIVIDEND_PORFIT"),
			OperatingCashFlowToRevenue: optionalNumber(row, "JYXJLYYSR"),
		}
		item.FreeCashFlow = freeCashFlow(item.OperatingCashFlow, item.CapitalExpenditure)
		annual = append(annual, item)
		if len(annual) >= financialYearsLimit {
			break
		}
	}

	if len(annual) == 0 {
		return nil, fmt.Errorf("empty financial data for %s", target.Symbol)
	}

	sortAnnualDesc(annual)
	enrichAnnualValuation(annual, target.CurrentPrice)
	fillDerivedGrowth(annual)
	return &Financials{
		UpdatedAt:    now.Format("2006-01-02 15:04:05"),
		Source:       "东方财富财务分析",
		SourceSymbol: secucode,
		Currency:     firstNonEmpty(annual[0].Currency, target.Currency, currencyForSymbol(target.Symbol)),
		Annual:       annual,
		Valuation:    buildFinancialValuation(annual, target),
	}, nil
}

func fetchHKFinancials(client *http.Client, target financialTarget, now time.Time) (*Financials, error) {
	secucode, err := eastmoneyFinancialSecucode(target.Symbol)
	if err != nil {
		return nil, err
	}

	rows, err := fetchEastmoneyDataRows(client, "https://datacenter.eastmoney.com/securities/api/data/v1/get", url.Values{
		"source":       []string{"F10"},
		"client":       []string{"PC"},
		"reportName":   []string{"RPT_HKF10_FN_MAININDICATOR"},
		"columns":      []string{"ALL"},
		"quoteColumns": []string{""},
		"filter":       []string{fmt.Sprintf(`(SECUCODE="%s")(DATE_TYPE_CODE="001")`, secucode)},
		"pageNumber":   []string{"1"},
		"pageSize":     []string{"10"},
		"sortTypes":    []string{"-1"},
		"sortColumns":  []string{"STD_REPORT_DATE"},
	})
	if err != nil {
		return nil, err
	}

	reportDates := make([]string, 0, financialYearsLimit)
	for _, row := range rows {
		reportDate := cleanReportDate(firstNonEmpty(mapString(row, "STD_REPORT_DATE"), mapString(row, "REPORT_DATE")))
		if reportDate != "" {
			reportDates = append(reportDates, reportDate)
		}
		if len(reportDates) >= financialYearsLimit {
			break
		}
	}

	cashRows, err := fetchHKCashflowRows(client, secucode, reportDates)
	if err != nil {
		cashRows = nil
	}
	cashByDate := hkCashflowByReportDate(cashRows)

	annual := make([]FinancialAnnual, 0, financialYearsLimit)
	for _, row := range rows {
		reportDate := cleanReportDate(firstNonEmpty(mapString(row, "STD_REPORT_DATE"), mapString(row, "REPORT_DATE")))
		cashItems := cashByDate[reportDate]
		item := FinancialAnnual{
			FiscalYear:        financialYear(reportDate),
			ReportDate:        reportDate,
			ReportType:        mapString(row, "REPORT_TYPE"),
			Currency:          firstNonEmpty(mapString(row, "CURRENCY"), target.Currency, currencyForSymbol(target.Symbol)),
			Revenue:           optionalNumber(row, "OPERATE_INCOME"),
			RevenueYoY:        optionalPercent(row, "OPERATE_INCOME_YOY"),
			NetProfit:         optionalNumber(row, "HOLDER_PROFIT"),
			NetProfitYoY:      optionalPercent(row, "HOLDER_PROFIT_YOY"),
			OperatingProfit:   optionalNumber(row, "OPERATE_PROFIT"),
			OperatingCashFlow: firstNumber(optionalNumber(row, "NETCASH_OPERATE"), hkCashflowAmount(cashItems, "003999")),
			CapitalExpenditure: absNumber(firstNumber(
				hkCashflowAmount(cashItems, "005005"),
				optionalNumber(row, "CAPITAL_EXPENDITURE", "CAPEX"),
			)),
			EPS:               optionalNumber(row, "BASIC_EPS"),
			BookValuePerShare: optionalNumber(row, "BPS"),
			TotalAssets:       optionalNumber(row, "TOTAL_ASSETS"),
			TotalLiabilities:  optionalNumber(row, "TOTAL_LIABILITIES"),
			TotalEquity:       optionalNumber(row, "TOTAL_PARENT_EQUITY"),
			ROE:               optionalPercent(row, "ROE_AVG"),
			ROIC:              optionalPercent(row, "ROIC_YEARLY"),
			GrossMargin:       optionalPercent(row, "GROSS_PROFIT_RATIO"),
			NetMargin:         optionalPercent(row, "NET_PROFIT_RATIO"),
			DividendPaid:      hkCashflowAmount(cashItems, "007004"),
		}
		item.OperatingMargin = ratioFromValues(item.OperatingProfit, item.Revenue)
		item.OperatingCashFlowToRevenue = ratioFromValues(item.OperatingCashFlow, item.Revenue)
		item.DebtRatio = ratioFromValues(item.TotalLiabilities, item.TotalAssets)
		item.FreeCashFlow = freeCashFlow(item.OperatingCashFlow, item.CapitalExpenditure)
		annual = append(annual, item)
		if len(annual) >= financialYearsLimit {
			break
		}
	}

	if len(annual) == 0 {
		return nil, fmt.Errorf("empty financial data for %s", target.Symbol)
	}

	sortAnnualDesc(annual)
	enrichAnnualValuation(annual, target.CurrentPrice)
	fillDerivedGrowth(annual)
	return &Financials{
		UpdatedAt:    now.Format("2006-01-02 15:04:05"),
		Source:       "东方财富港股财务分析",
		SourceSymbol: secucode,
		Currency:     firstNonEmpty(annual[0].Currency, target.Currency, currencyForSymbol(target.Symbol)),
		Annual:       annual,
		Valuation:    buildFinancialValuation(annual, target),
	}, nil
}

func fetchHKCashflowRows(client *http.Client, secucode string, reportDates []string) ([]map[string]any, error) {
	reportDates = uniqueStrings(reportDates)
	if len(reportDates) == 0 {
		return nil, nil
	}
	return fetchEastmoneyDataRows(client, "https://datacenter.eastmoney.com/securities/api/data/v1/get", url.Values{
		"source":       []string{"F10"},
		"client":       []string{"PC"},
		"reportName":   []string{"RPT_HKF10_FN_CASHFLOW_PC"},
		"columns":      []string{"SECUCODE,SECURITY_CODE,SECURITY_NAME_ABBR,ORG_CODE,REPORT_DATE,DATE_TYPE_CODE,FISCAL_YEAR,START_DATE,STD_ITEM_CODE,STD_ITEM_NAME,AMOUNT"},
		"quoteColumns": []string{""},
		"filter":       []string{fmt.Sprintf(`(SECUCODE="%s")(REPORT_DATE in (%s))`, secucode, quotedList(reportDates))},
		"pageNumber":   []string{"1"},
		"pageSize":     []string{"300"},
		"sortTypes":    []string{"-1,1"},
		"sortColumns":  []string{"REPORT_DATE,STD_ITEM_CODE"},
	})
}

func hkCashflowByReportDate(rows []map[string]any) map[string]map[string]*float64 {
	byDate := make(map[string]map[string]*float64)
	for _, row := range rows {
		reportDate := cleanReportDate(mapString(row, "REPORT_DATE"))
		itemCode := mapString(row, "STD_ITEM_CODE")
		amount := optionalNumber(row, "AMOUNT")
		if reportDate == "" || itemCode == "" || amount == nil {
			continue
		}
		if byDate[reportDate] == nil {
			byDate[reportDate] = make(map[string]*float64)
		}
		byDate[reportDate][itemCode] = amount
	}
	return byDate
}

func hkCashflowAmount(items map[string]*float64, itemCodes ...string) *float64 {
	for _, itemCode := range itemCodes {
		if value := items[itemCode]; value != nil {
			return value
		}
	}
	return nil
}

func fetchEastmoneyDataRows(client *http.Client, endpoint string, params url.Values) ([]map[string]any, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Referer", "https://emweb.securities.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website financials updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("financial request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}

	var payload struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Code    int    `json:"code"`
		Result  struct {
			Data []map[string]any `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if !payload.Success && len(payload.Result.Data) == 0 {
		if strings.TrimSpace(payload.Message) != "" {
			return nil, errors.New(payload.Message)
		}
		return nil, fmt.Errorf("financial request failed with code %d", payload.Code)
	}
	if len(payload.Result.Data) == 0 {
		return nil, errors.New("empty financial response")
	}
	return payload.Result.Data, nil
}

func eastmoneyFinancialSecucode(symbol string) (string, error) {
	symbol = normalizeSymbol(symbol)
	switch {
	case strings.HasSuffix(symbol, ".SH"), strings.HasSuffix(symbol, ".SZ"):
		return symbol, nil
	case strings.HasSuffix(symbol, ".HK"):
		code := strings.TrimSuffix(symbol, ".HK")
		if value, err := strconvAtoi(code); err == nil {
			code = fmt.Sprintf("%05d", value)
		}
		return code + ".HK", nil
	default:
		return "", fmt.Errorf("unsupported eastmoney financial symbol: %s", symbol)
	}
}

func rowsByReportDate(rows []map[string]any) map[string]map[string]any {
	byDate := make(map[string]map[string]any, len(rows))
	for _, row := range rows {
		date := cleanReportDate(firstNonEmpty(mapString(row, "REPORT_DATE"), mapString(row, "STD_REPORT_DATE")))
		if date != "" {
			byDate[date] = row
		}
	}
	return byDate
}

func isAnnualFinancialRow(row map[string]any) bool {
	if strings.EqualFold(mapString(row, "DATE_TYPE_CODE"), "001") {
		return true
	}
	reportType := mapString(row, "REPORT_TYPE") + " " + mapString(row, "REPORT_DATE_NAME")
	if strings.Contains(reportType, "年报") {
		return true
	}
	return strings.HasSuffix(cleanReportDate(mapString(row, "REPORT_DATE")), "-12-31") ||
		strings.HasSuffix(cleanReportDate(mapString(row, "STD_REPORT_DATE")), "-12-31")
}

func mapString(data map[string]any, key string) string {
	if data == nil {
		return ""
	}
	value, ok := data[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		if math.Trunc(typed) == typed {
			return fmt.Sprintf("%.0f", typed)
		}
		return fmt.Sprintf("%f", typed)
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func optionalNumber(data map[string]any, keys ...string) *float64 {
	for _, key := range keys {
		value, err := numberField(data, key)
		if err == nil && !math.IsNaN(value) && !math.IsInf(value, 0) {
			return ptr(value)
		}
	}
	return nil
}

func firstNumber(values ...*float64) *float64 {
	for _, value := range values {
		if value != nil && !math.IsNaN(*value) && !math.IsInf(*value, 0) {
			return value
		}
	}
	return nil
}

func optionalPercent(data map[string]any, keys ...string) *float64 {
	value := optionalNumber(data, keys...)
	if value == nil {
		return nil
	}
	return ptr(*value / 100)
}

func absNumber(value *float64) *float64 {
	if value == nil {
		return nil
	}
	return ptr(math.Abs(*value))
}

func ratioFromValues(numerator *float64, denominator *float64) *float64 {
	if numerator == nil || denominator == nil || *denominator == 0 {
		return nil
	}
	return ptr(*numerator / *denominator)
}

func freeCashFlow(operatingCashFlow *float64, capitalExpenditure *float64) *float64 {
	if operatingCashFlow == nil {
		return nil
	}
	if capitalExpenditure == nil {
		return ptr(*operatingCashFlow)
	}
	return ptr(*operatingCashFlow - math.Abs(*capitalExpenditure))
}

func cleanReportDate(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= len("2006-01-02") {
		return value[:len("2006-01-02")]
	}
	return value
}

func financialYear(reportDate string) string {
	reportDate = cleanReportDate(reportDate)
	if len(reportDate) >= 4 {
		return reportDate[:4]
	}
	return ""
}

func sortAnnualDesc(annual []FinancialAnnual) {
	sort.SliceStable(annual, func(i, j int) bool {
		return annual[i].ReportDate > annual[j].ReportDate
	})
}

func enrichAnnualValuation(annual []FinancialAnnual, currentPrice float64) {
	if currentPrice <= 0 {
		return
	}
	for i := range annual {
		if annual[i].EPS != nil && *annual[i].EPS > 0 {
			annual[i].PEAtCurrentPrice = ptr(currentPrice / *annual[i].EPS)
		}
		if annual[i].BookValuePerShare != nil && *annual[i].BookValuePerShare > 0 {
			annual[i].PBAtCurrentPrice = ptr(currentPrice / *annual[i].BookValuePerShare)
		}
	}
}

func fillDerivedGrowth(annual []FinancialAnnual) {
	for i := range annual {
		if i+1 >= len(annual) {
			continue
		}
		if annual[i].RevenueYoY == nil {
			annual[i].RevenueYoY = growthFromValues(annual[i].Revenue, annual[i+1].Revenue)
		}
		if annual[i].NetProfitYoY == nil {
			annual[i].NetProfitYoY = growthFromValues(annual[i].NetProfit, annual[i+1].NetProfit)
		}
		if annual[i].OperatingCashFlowToRevenue == nil {
			annual[i].OperatingCashFlowToRevenue = ratioFromValues(annual[i].OperatingCashFlow, annual[i].Revenue)
		}
	}
}

func growthFromValues(current *float64, previous *float64) *float64 {
	if current == nil || previous == nil || *previous == 0 {
		return nil
	}
	return ptr((*current - *previous) / math.Abs(*previous))
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		unique = append(unique, value)
	}
	return unique
}

func quotedList(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.ReplaceAll(strings.TrimSpace(value), `'`, `''`)
		if value != "" {
			quoted = append(quoted, fmt.Sprintf("'%s'", value))
		}
	}
	return strings.Join(quoted, ",")
}

func buildFinancialValuation(annual []FinancialAnnual, target financialTarget) *FinancialValuation {
	if len(annual) == 0 {
		return nil
	}
	current := annual[0]
	valuation := &FinancialValuation{
		Date:          current.ReportDate,
		Currency:      firstNonEmpty(target.Currency, current.Currency, currencyForSymbol(target.Symbol)),
		PE:            current.PEAtCurrentPrice,
		PB:            current.PBAtCurrentPrice,
		DividendYield: financialDividendYield(target),
		PERange:       metricRangeFromAnnual(annual, func(item FinancialAnnual) *float64 { return item.PEAtCurrentPrice }),
		PBRange:       metricRangeFromAnnual(annual, func(item FinancialAnnual) *float64 { return item.PBAtCurrentPrice }),
		SourceNote:    "PE/PB 区间为按当前价格除以历史 EPS/BPS 的回看口径，不等同历史市场估值分位",
	}
	if target.CurrentPrice > 0 {
		valuation.Price = ptr(target.CurrentPrice)
	}
	if valuation.PE != nil && current.NetProfitYoY != nil && *current.NetProfitYoY > 0 {
		valuation.PEG = ptr(*valuation.PE / (*current.NetProfitYoY * 100))
	}
	return valuation
}

func metricRangeFromAnnual(annual []FinancialAnnual, valueFor func(FinancialAnnual) *float64) *MetricRange {
	values := []float64{}
	for _, item := range annual {
		value := valueFor(item)
		if value != nil && *value > 0 && !math.IsNaN(*value) && !math.IsInf(*value, 0) {
			values = append(values, *value)
		}
	}
	if len(values) == 0 {
		return nil
	}
	sort.Float64s(values)
	return &MetricRange{
		Min:    ptr(values[0]),
		Median: ptr(median(values)),
		Max:    ptr(values[len(values)-1]),
	}
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle]
	}
	return (values[middle-1] + values[middle]) / 2
}

func financialDividendYield(target financialTarget) *float64 {
	if target.Dividend == nil {
		return nil
	}
	dividend := target.Dividend
	if dividend.DividendYield != nil && *dividend.DividendYield > 0 {
		return dividend.DividendYield
	}
	if dividend.CashDividendTotal != nil && *dividend.CashDividendTotal > 0 && target.MarketCap != nil && *target.MarketCap > 0 {
		return ptr(*dividend.CashDividendTotal / *target.MarketCap)
	}
	if dividend.DividendPerShare != nil && *dividend.DividendPerShare > 0 && target.CurrentPrice > 0 {
		return ptr(*dividend.DividendPerShare / target.CurrentPrice)
	}
	return nil
}

func financialDecisionDetail(financials *Financials) string {
	if financials == nil || len(financials.Annual) == 0 {
		return "未取得可展示年度数据"
	}
	latest := financials.Annual[0]
	parts := []string{
		fmt.Sprintf("来源 %s", firstNonEmpty(financials.Source, "未知")),
		fmt.Sprintf("年度 %d 年", len(financials.Annual)),
		fmt.Sprintf("最新期 %s", firstNonEmpty(latest.FiscalYear, latest.ReportDate, "未知")),
	}
	if latest.Revenue != nil {
		parts = append(parts, fmt.Sprintf("收入 %.0f", *latest.Revenue))
	}
	if latest.NetProfit != nil {
		parts = append(parts, fmt.Sprintf("利润 %.0f", *latest.NetProfit))
	}
	if latest.ROE != nil {
		parts = append(parts, fmt.Sprintf("ROE %.2f%%", *latest.ROE*100))
	}
	if latest.ROIC != nil {
		parts = append(parts, fmt.Sprintf("ROIC %.2f%%", *latest.ROIC*100))
	}
	if latest.FreeCashFlow != nil {
		parts = append(parts, fmt.Sprintf("FCF %.0f", *latest.FreeCashFlow))
	}
	return strings.Join(parts, "；")
}

func strconvAtoi(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, errors.New("empty integer")
	}
	number := 0
	for _, char := range value {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("invalid integer: %s", value)
		}
		number = number*10 + int(char-'0')
	}
	return number, nil
}
