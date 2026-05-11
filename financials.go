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
	"strconv"
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
	DividendPerShare           *float64 `json:"dividendPerShare,omitempty"`
	DividendCurrency           string   `json:"dividendCurrency,omitempty"`
	BuybackAmount              *float64 `json:"buybackAmount,omitempty"`
	BuybackCurrency            string   `json:"buybackCurrency,omitempty"`
	CashAndShortInvestments    *float64 `json:"cashAndShortInvestments,omitempty"`
	InterestBearingDebt        *float64 `json:"interestBearingDebt,omitempty"`
	NetCash                    *float64 `json:"netCash,omitempty"`
	TotalAssets                *float64 `json:"totalAssets,omitempty"`
	TotalLiabilities           *float64 `json:"totalLiabilities,omitempty"`
	TotalEquity                *float64 `json:"totalEquity,omitempty"`
	ParentEquity               *float64 `json:"parentEquity,omitempty"`
	MinorityEquity             *float64 `json:"minorityEquity,omitempty"`
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
	Industry     string
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
			Industry:     state.Holdings[i].Industry,
			Currency:     state.Holdings[i].Currency,
			CurrentPrice: state.Holdings[i].CurrentPrice,
			MarketCap:    state.Holdings[i].MarketCap,
			Dividend:     state.Holdings[i].Dividend,
		}
		return target, "holding", func(financials *Financials) {
			applyFinancialFactsToHolding(&state.Holdings[i], financials)
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
			Industry:     state.Candidates[i].Industry,
			Currency:     state.Candidates[i].Currency,
			CurrentPrice: state.Candidates[i].CurrentPrice,
			MarketCap:    state.Candidates[i].MarketCap,
			Dividend:     state.Candidates[i].Dividend,
		}
		return target, "candidate", func(financials *Financials) {
			applyFinancialFactsToCandidate(&state.Candidates[i], financials)
			state.Candidates[i].Financials = financials
		}
	}

	return financialTarget{}, "", nil
}

func applyFinancialFactsToHolding(holding *Holding, financials *Financials) {
	if holding == nil || financials == nil || len(financials.Annual) == 0 {
		return
	}
	applyFinancialDividend(&holding.Dividend, financials, holding.Currency, holding.CurrentPrice, holding.MarketCap)
	applyFinancialNetCash(&holding.NetCash, financials, holding.Currency, holding.Industry)
	if financials.Valuation != nil {
		financials.Valuation.DividendYield = financialDividendYield(financialTarget{
			Symbol:       holding.Symbol,
			Currency:     holding.Currency,
			CurrentPrice: holding.CurrentPrice,
			MarketCap:    holding.MarketCap,
			Dividend:     holding.Dividend,
		})
	}
}

func applyFinancialFactsToCandidate(candidate *Candidate, financials *Financials) {
	if candidate == nil || financials == nil || len(financials.Annual) == 0 {
		return
	}
	applyFinancialDividend(&candidate.Dividend, financials, candidate.Currency, candidate.CurrentPrice, candidate.MarketCap)
	applyFinancialNetCash(&candidate.NetCash, financials, candidate.Currency, candidate.Industry)
	if financials.Valuation != nil {
		financials.Valuation.DividendYield = financialDividendYield(financialTarget{
			Symbol:       candidate.Symbol,
			Currency:     candidate.Currency,
			CurrentPrice: candidate.CurrentPrice,
			MarketCap:    candidate.MarketCap,
			Dividend:     candidate.Dividend,
		})
	}
}

func applyFinancialDividend(current **Dividend, financials *Financials, fallbackCurrency string, currentPrice float64, marketCap *float64) {
	if financials == nil || len(financials.Annual) == 0 {
		return
	}
	latest := financials.Annual[0]
	currency := strings.ToUpper(firstNonEmpty(latest.Currency, financials.Currency, fallbackCurrency))
	if *current == nil {
		*current = &Dividend{}
	}
	dividend := *current
	if strings.TrimSpace(latest.FiscalYear) != "" {
		dividend.FiscalYear = "FY" + strings.TrimPrefix(strings.TrimSpace(latest.FiscalYear), "FY")
	}
	if latest.DividendPerShare != nil && *latest.DividendPerShare > 0 {
		dividend.DividendPerShare = cloneFloat(latest.DividendPerShare)
		dividend.DividendCurrency = strings.ToUpper(firstNonEmpty(latest.DividendCurrency, currency))
	}
	if latest.BuybackAmount != nil && *latest.BuybackAmount > 0 {
		dividend.BuybackAmount = cloneFloat(latest.BuybackAmount)
		dividend.BuybackCurrency = strings.ToUpper(firstNonEmpty(latest.BuybackCurrency, currency))
	} else {
		dividend.BuybackAmount = nil
	}
	if dividend.DividendPerShare != nil && *dividend.DividendPerShare > 0 && marketCap != nil && *marketCap > 0 && currentPrice > 0 {
		dividend.CashDividendTotal = ptr((*dividend.DividendPerShare * *marketCap) / currentPrice)
		dividend.CashDividendCurrency = strings.ToUpper(firstNonEmpty(dividend.DividendCurrency, currency))
	}
	if strings.TrimSpace(dividend.DividendCurrency) == "" {
		dividend.DividendCurrency = currency
	}
	if strings.TrimSpace(dividend.CashDividendCurrency) == "" {
		dividend.CashDividendCurrency = dividend.DividendCurrency
	}
	if strings.TrimSpace(dividend.BuybackCurrency) == "" {
		dividend.BuybackCurrency = dividend.CashDividendCurrency
	}
	dividend.DividendYield = nil
	dividend.EstimatedAnnualCash = nil
}

func applyFinancialNetCash(current **NetCashProfile, financials *Financials, fallbackCurrency string, industry string) {
	if financials == nil || len(financials.Annual) == 0 || !netCashApplicableForIndustry(industry) {
		return
	}
	latest := financials.Annual[0]
	if latest.CashAndShortInvestments == nil && latest.InterestBearingDebt == nil && latest.FreeCashFlow == nil {
		return
	}
	if *current == nil {
		*current = &NetCashProfile{}
	}
	profile := *current
	currency := strings.ToUpper(firstNonEmpty(latest.Currency, financials.Currency, fallbackCurrency))
	profile.Currency = currency
	profile.CashAndShortInvestments = cloneFloat(latest.CashAndShortInvestments)
	profile.InterestBearingDebt = cloneFloat(latest.InterestBearingDebt)
	profile.NetCash = firstNumber(cloneFloat(latest.NetCash), netCashAmount(latest.CashAndShortInvestments, latest.InterestBearingDebt))
	profile.ConsolidatedFCF = cloneFloat(latest.FreeCashFlow)
	profile.FCFPositiveYears = positiveFreeCashFlowYears(financials.Annual)
	profile.ShareholderFCF, profile.ShareholderFCFBasis, profile.MinorityFCFAdjustment = ordinaryShareholderFreeCashFlow(latest)
	if profile.ShareholderFCF != nil {
		profile.ShareholderFCFCurrency = currency
	}

	// These are formulas that depend on runtime market cap and current strategy haircut.
	// Clear imported/stale values so the frontend recomputes them from button-owned data.
	profile.AdjustedNetCash = nil
	profile.ExCashPE = nil
	profile.ExCashPFCF = nil
	profile.FCFYield = nil
	profile.Note = "由更新财务根据最新年度资产负债表、现金流量表和少数股东权益自动刷新。"
}

func netCashApplicableForIndustry(industry string) bool {
	for _, part := range strings.Split(industry, "/") {
		part = strings.TrimSpace(part)
		switch part {
		case "银行", "保险", "券商", "证券", "信托", "财富管理":
			return false
		}
	}
	return true
}

func ordinaryShareholderFreeCashFlow(latest FinancialAnnual) (*float64, string, *float64) {
	if latest.FreeCashFlow == nil {
		return nil, "", nil
	}
	totalEquity := cloneFloat(latest.TotalEquity)
	if totalEquity == nil && latest.ParentEquity != nil && latest.MinorityEquity != nil {
		totalEquity = ptr(*latest.ParentEquity + *latest.MinorityEquity)
	}
	if latest.ParentEquity != nil && totalEquity != nil && *latest.ParentEquity > 0 && *totalEquity > 0 && *latest.ParentEquity < *totalEquity*0.95 {
		value := *latest.FreeCashFlow * (*latest.ParentEquity / *totalEquity)
		adjustment := *latest.FreeCashFlow - value
		return &value, "普通股东 FCF：合并 FCF 按归属普通股东权益比例扣除少数股东分流", &adjustment
	}
	value := *latest.FreeCashFlow
	return &value, "普通股东 FCF：未发现重大少数股东权益分流，使用合并 FCF", nil
}

func positiveFreeCashFlowYears(annual []FinancialAnnual) *int {
	if len(annual) == 0 {
		return nil
	}
	count := 0
	for _, item := range annual {
		if item.FreeCashFlow != nil && *item.FreeCashFlow > 0 {
			count++
		}
	}
	return &count
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

	balanceRows, err := fetchEastmoneyDataRows(client, "https://datacenter.eastmoney.com/securities/api/data/get", url.Values{
		"source": []string{"HSF10"},
		"client": []string{"PC"},
		"type":   []string{"RPT_F10_FINANCE_GBALANCE"},
		"sty":    []string{"F10_FINANCE_GBALANCE"},
		"filter": []string{fmt.Sprintf(`(SECUCODE="%s")`, secucode)},
		"ps":     []string{"80"},
		"sr":     []string{"-1"},
		"st":     []string{"REPORT_DATE"},
	})
	if err != nil {
		balanceRows = nil
	}

	dividends, err := fetchAshareDividendDistributions(client, secucode)
	if err != nil {
		dividends = dividendDistributions{}
	}
	buybacks, err := fetchAshareBuybacks(client, secucode)
	if err != nil {
		buybacks = map[string]*float64{}
	}
	cashByDate := rowsByReportDate(cashRows)
	balanceByDate := rowsByReportDate(balanceRows)
	annual := make([]FinancialAnnual, 0, financialYearsLimit)
	for _, row := range mainRows {
		if !isAnnualFinancialRow(row) {
			continue
		}
		reportDate := cleanReportDate(firstNonEmpty(mapString(row, "REPORT_DATE"), mapString(row, "STD_REPORT_DATE")))
		fiscalYear := firstNonEmpty(mapString(row, "REPORT_YEAR"), financialYear(reportDate))
		dividend := dividends[fiscalYear]
		cashRow := cashByDate[reportDate]
		balanceRow := balanceByDate[reportDate]
		item := FinancialAnnual{
			FiscalYear:                 fiscalYear,
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
			DividendPerShare:           cloneFloat(dividend.PerShare),
			DividendCurrency:           dividend.Currency,
			BuybackAmount:              cloneFloat(buybacks[fiscalYear]),
			BuybackCurrency:            "CNY",
			CashAndShortInvestments:    ashareCashAndShortInvestments(balanceRow),
			InterestBearingDebt:        ashareInterestBearingDebt(balanceRow),
			TotalAssets:                firstNumber(optionalNumber(balanceRow, "TOTAL_ASSETS", "ASSET_BALANCE"), optionalNumber(row, "TOTAL_ASSETS")),
			TotalLiabilities:           firstNumber(optionalNumber(balanceRow, "TOTAL_LIABILITIES", "LIAB_BALANCE"), optionalNumber(row, "LIABILITY")),
			TotalEquity:                optionalNumber(balanceRow, "TOTAL_EQUITY", "EQUITY_BALANCE"),
			ParentEquity:               optionalNumber(balanceRow, "TOTAL_PARENT_EQUITY", "PARENT_EQUITY_BALANCE"),
			MinorityEquity:             optionalNumber(balanceRow, "MINORITY_EQUITY"),
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
		item.NetCash = netCashAmount(item.CashAndShortInvestments, item.InterestBearingDebt)
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

	balanceRows, err := fetchHKBalanceRows(client, secucode, reportDates)
	if err != nil {
		balanceRows = nil
	}
	balanceByDate := hkRowsByReportDate(balanceRows)

	dividends, err := fetchHKDividendDistributions(client, secucode)
	if err != nil {
		dividends = dividendDistributions{}
	}
	annual := make([]FinancialAnnual, 0, financialYearsLimit)
	for _, row := range rows {
		reportDate := cleanReportDate(firstNonEmpty(mapString(row, "STD_REPORT_DATE"), mapString(row, "REPORT_DATE")))
		fiscalYear := financialYear(reportDate)
		dividend := dividends[fiscalYear]
		cashItems := cashByDate[reportDate]
		balanceItems := balanceByDate[reportDate]
		item := FinancialAnnual{
			FiscalYear:        fiscalYear,
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
			EPS:                     optionalNumber(row, "BASIC_EPS"),
			BookValuePerShare:       optionalNumber(row, "BPS"),
			DividendPerShare:        cloneFloat(dividend.PerShare),
			DividendCurrency:        dividend.Currency,
			BuybackAmount:           hkCashflowAmount(cashItems, "007008"),
			BuybackCurrency:         "HKD",
			CashAndShortInvestments: hkBalanceCashAndShortInvestments(balanceItems),
			InterestBearingDebt:     hkBalanceInterestBearingDebt(balanceItems),
			TotalAssets:             firstNumber(hkBalanceAmount(balanceItems, "004009999", "004039999"), optionalNumber(row, "TOTAL_ASSETS")),
			TotalLiabilities:        firstNumber(hkBalanceAmount(balanceItems, "004025999"), optionalNumber(row, "TOTAL_LIABILITIES")),
			TotalEquity:             firstNumber(hkBalanceAmount(balanceItems, "004036999"), optionalNumber(row, "TOTAL_EQUITY")),
			ParentEquity:            firstNumber(hkBalanceAmount(balanceItems, "004030999"), optionalNumber(row, "TOTAL_PARENT_EQUITY")),
			MinorityEquity:          hkBalanceAmount(balanceItems, "004027999"),
			ROE:                     optionalPercent(row, "ROE_AVG"),
			ROIC:                    optionalPercent(row, "ROIC_YEARLY"),
			GrossMargin:             optionalPercent(row, "GROSS_PROFIT_RATIO"),
			NetMargin:               optionalPercent(row, "NET_PROFIT_RATIO"),
			DividendPaid:            hkCashflowAmount(cashItems, "007004"),
		}
		item.OperatingMargin = ratioFromValues(item.OperatingProfit, item.Revenue)
		item.OperatingCashFlowToRevenue = ratioFromValues(item.OperatingCashFlow, item.Revenue)
		item.DebtRatio = ratioFromValues(item.TotalLiabilities, item.TotalAssets)
		item.NetCash = netCashAmount(item.CashAndShortInvestments, item.InterestBearingDebt)
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

func fetchHKBalanceRows(client *http.Client, secucode string, reportDates []string) ([]map[string]any, error) {
	reportDates = uniqueStrings(reportDates)
	if len(reportDates) == 0 {
		return nil, nil
	}
	return fetchEastmoneyDataRows(client, "https://datacenter.eastmoney.com/securities/api/data/v1/get", url.Values{
		"source":       []string{"F10"},
		"client":       []string{"PC"},
		"reportName":   []string{"RPT_HKF10_FN_BALANCE_PC"},
		"columns":      []string{"SECUCODE,SECURITY_CODE,SECURITY_NAME_ABBR,ORG_CODE,REPORT_DATE,DATE_TYPE_CODE,FISCAL_YEAR,STD_ITEM_CODE,STD_ITEM_NAME,AMOUNT,STD_REPORT_DATE"},
		"quoteColumns": []string{""},
		"filter":       []string{fmt.Sprintf(`(SECUCODE="%s")(REPORT_DATE in (%s))`, secucode, quotedList(reportDates))},
		"pageNumber":   []string{"1"},
		"pageSize":     []string{"500"},
		"sortTypes":    []string{"-1,1"},
		"sortColumns":  []string{"REPORT_DATE,STD_ITEM_CODE"},
	})
}

type dividendDistribution struct {
	PerShare *float64
	Currency string
}

type dividendDistributions map[string]dividendDistribution

func (items dividendDistributions) add(fiscalYear string, amount float64, currency string) {
	fiscalYear = strings.TrimSpace(fiscalYear)
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if fiscalYear == "" || amount <= 0 {
		return
	}
	item := items[fiscalYear]
	if item.PerShare == nil {
		item.PerShare = ptr(0)
	}
	*item.PerShare += amount
	if item.Currency == "" {
		item.Currency = currency
	}
	items[fiscalYear] = item
}

func fetchAshareDividendDistributions(client *http.Client, secucode string) (dividendDistributions, error) {
	code := strings.TrimSuffix(normalizeSymbol(secucode), ".SH")
	code = strings.TrimSuffix(code, ".SZ")
	rows, err := fetchEastmoneyDataRows(client, "https://datacenter-web.eastmoney.com/api/data/v1/get", url.Values{
		"source":       []string{"WEB"},
		"client":       []string{"WEB"},
		"reportName":   []string{"RPT_SHAREBONUS_DET"},
		"columns":      []string{"ALL"},
		"quoteColumns": []string{""},
		"filter":       []string{fmt.Sprintf(`(SECURITY_CODE="%s")`, code)},
		"pageNumber":   []string{"1"},
		"pageSize":     []string{"100"},
		"sortTypes":    []string{"-1"},
		"sortColumns":  []string{"PLAN_NOTICE_DATE"},
	})
	if err != nil {
		return nil, err
	}
	dividends := dividendDistributions{}
	for _, row := range rows {
		amount := optionalNumber(row, "PRETAX_BONUS_RMB")
		if amount == nil || *amount <= 0 {
			continue
		}
		year := financialYear(mapString(row, "REPORT_DATE"))
		dividends.add(year, *amount/10, "CNY")
	}
	return dividends, nil
}

func fetchAshareBuybacks(client *http.Client, secucode string) (map[string]*float64, error) {
	code := strings.TrimSuffix(normalizeSymbol(secucode), ".SH")
	code = strings.TrimSuffix(code, ".SZ")
	rows, err := fetchEastmoneyDataRows(client, "https://datacenter-web.eastmoney.com/api/data/v1/get", url.Values{
		"source":       []string{"WEB"},
		"client":       []string{"WEB"},
		"reportName":   []string{"RPTA_WEB_GETHGLIST_NEW"},
		"columns":      []string{"ALL"},
		"quoteColumns": []string{""},
		"filter":       []string{fmt.Sprintf(`(DIM_SCODE="%s")`, code)},
		"pageNumber":   []string{"1"},
		"pageSize":     []string{"100"},
		"sortTypes":    []string{"-1,-1,-1"},
		"sortColumns":  []string{"UPD,DIM_DATE,DIM_SCODE"},
	})
	if err != nil {
		return nil, err
	}
	byYear := map[string]*float64{}
	for _, row := range rows {
		amount := optionalNumber(row, "REPURAMOUNT", "ZJJE", "JEXX")
		if amount == nil || *amount <= 0 {
			continue
		}
		year := ashareBuybackFiscalYear(row)
		if year == "" {
			continue
		}
		if byYear[year] == nil {
			byYear[year] = ptr(0)
		}
		*byYear[year] += *amount
	}
	return byYear, nil
}

func ashareBuybackFiscalYear(row map[string]any) string {
	// REPORTDATE ties the buyback progress row to a financial period when present.
	// Otherwise use the execution/finish progress date, then the plan start date.
	return financialYear(firstNonEmpty(
		mapString(row, "REPORTDATE"),
		mapString(row, "FINISHDATE"),
		mapString(row, "EDATE"),
		mapString(row, "REPURADVANCEDATE"),
		mapString(row, "UPD"),
		mapString(row, "UPDATEDATE"),
		mapString(row, "NOTICEDATE"),
		mapString(row, "DIM_DATE3"),
		mapString(row, "REPURSTARTDATE"),
		mapString(row, "DIM_DATE"),
	))
}

func fetchHKDividendDistributions(client *http.Client, secucode string) (dividendDistributions, error) {
	code := strings.TrimSuffix(eastmoneyHKCode(secucode), ".HK")
	rows, err := fetchEastmoneyDataRows(client, "https://datacenter.eastmoney.com/securities/api/data/v1/get", url.Values{
		"source":       []string{"F10"},
		"client":       []string{"PC"},
		"reportName":   []string{"RPT_HKF10_MAIN_DIVBASIC"},
		"columns":      []string{"SECURITY_CODE,UPDATE_DATE,REPORT_TYPE,EX_DIVIDEND_DATE,DIVIDEND_DATE,TRANSFER_END_DATE,YEAR,PLAN_EXPLAIN,IS_BFP"},
		"quoteColumns": []string{""},
		"filter":       []string{fmt.Sprintf(`(SECURITY_CODE="%s")(IS_BFP="0")`, code)},
		"pageNumber":   []string{"1"},
		"pageSize":     []string{"200"},
		"sortTypes":    []string{"-1,-1"},
		"sortColumns":  []string{"NOTICE_DATE,EX_DIVIDEND_DATE"},
	})
	if err != nil {
		return nil, err
	}
	dividends := dividendDistributions{}
	for _, row := range rows {
		amount, currency := hkDividendAmount(mapString(row, "PLAN_EXPLAIN"))
		dividends.add(mapString(row, "YEAR"), amount, currency)
	}
	return dividends, nil
}

func eastmoneyHKCode(secucode string) string {
	code := strings.TrimSuffix(normalizeSymbol(secucode), ".HK")
	if value, err := strconvAtoi(code); err == nil {
		code = fmt.Sprintf("%05d", value)
	}
	return code + ".HK"
}

func hkDividendAmount(plan string) (float64, string) {
	plan = strings.TrimSpace(plan)
	currency := "HKD"
	switch {
	case strings.Contains(plan, "人民币"):
		currency = "CNY"
	case strings.Contains(plan, "港币"):
		currency = "HKD"
	}
	start := strings.Index(plan, "每股派")
	if start < 0 {
		return 0, currency
	}
	rest := plan[start+len("每股派"):]
	number := strings.Builder{}
	seenDigit := false
	for _, char := range rest {
		if (char >= '0' && char <= '9') || char == '.' {
			number.WriteRune(char)
			seenDigit = true
			continue
		}
		if seenDigit {
			break
		}
	}
	amount, err := strconv.ParseFloat(number.String(), 64)
	if err != nil || amount <= 0 {
		return 0, currency
	}
	return amount, currency
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

func hkRowsByReportDate(rows []map[string]any) map[string][]map[string]any {
	byDate := make(map[string][]map[string]any)
	for _, row := range rows {
		reportDate := cleanReportDate(firstNonEmpty(mapString(row, "REPORT_DATE"), mapString(row, "STD_REPORT_DATE")))
		if reportDate == "" {
			continue
		}
		byDate[reportDate] = append(byDate[reportDate], row)
	}
	return byDate
}

func hkBalanceAmount(rows []map[string]any, itemCodes ...string) *float64 {
	for _, code := range itemCodes {
		for _, row := range rows {
			if mapString(row, "STD_ITEM_CODE") != code {
				continue
			}
			if value := optionalNumber(row, "AMOUNT"); value != nil {
				return value
			}
		}
	}
	return nil
}

func hkBalanceCashAndShortInvestments(rows []map[string]any) *float64 {
	return hkBalanceSum(rows,
		[]string{"004002009", "004002010", "004002014"},
		[]string{"现金及等价物", "受限制存款", "交易性金融资产"},
	)
}

func hkBalanceInterestBearingDebt(rows []map[string]any) *float64 {
	return hkBalanceSum(rows,
		[]string{"004011006", "004020001", "004020005"},
		[]string{"贷款", "借款", "债券", "融资租赁负债"},
	)
}

func hkBalanceSum(rows []map[string]any, itemCodes []string, nameKeywords []string) *float64 {
	codeSet := make(map[string]bool, len(itemCodes))
	for _, code := range itemCodes {
		codeSet[code] = true
	}
	total := 0.0
	found := false
	for _, row := range rows {
		code := mapString(row, "STD_ITEM_CODE")
		name := mapString(row, "STD_ITEM_NAME")
		if !codeSet[code] && !containsAny(name, nameKeywords) {
			continue
		}
		value := optionalNumber(row, "AMOUNT")
		if value == nil {
			continue
		}
		total += *value
		found = true
	}
	if !found {
		return nil
	}
	return &total
}

func ashareCashAndShortInvestments(row map[string]any) *float64 {
	return sumOptionalNumbers(row, "MONETARYFUNDS", "TRADE_FINASSET", "FVTPL_FINASSET", "DERIVE_FINASSET")
}

func ashareInterestBearingDebt(row map[string]any) *float64 {
	return sumOptionalNumbers(row, "SHORT_LOAN", "LONG_LOAN", "BOND_PAYABLE", "SHORT_BOND_PAYABLE", "SHORT_FIN_PAYABLE", "NONCURRENT_LIAB_1YEAR", "LEASE_LIAB")
}

func sumOptionalNumbers(row map[string]any, keys ...string) *float64 {
	total := 0.0
	found := false
	for _, key := range keys {
		value := optionalNumber(row, key)
		if value == nil {
			continue
		}
		total += *value
		found = true
	}
	if !found {
		return nil
	}
	return &total
}

func netCashAmount(cash *float64, debt *float64) *float64 {
	if cash == nil || debt == nil {
		return nil
	}
	return ptr(*cash - *debt)
}

func containsAny(value string, keywords []string) bool {
	for _, keyword := range keywords {
		if keyword != "" && strings.Contains(value, keyword) {
			return true
		}
	}
	return false
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
