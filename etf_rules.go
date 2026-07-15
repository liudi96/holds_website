package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ETFRuleStatus struct {
	Symbol            string          `json:"symbol"`
	Name              string          `json:"name"`
	Level             string          `json:"level,omitempty"`
	LevelLabel        string          `json:"levelLabel,omitempty"`
	MonthlyAmount     float64         `json:"monthlyAmount,omitempty"`
	WeeklyAmount      float64         `json:"weeklyAmount,omitempty"`
	Complete          bool            `json:"complete"`
	Reason            string          `json:"reason,omitempty"`
	AsOf              string          `json:"asOf,omitempty"`
	UpdatedAt         string          `json:"updatedAt,omitempty"`
	LevelUpdatedAt    string          `json:"levelUpdatedAt,omitempty"`
	PendingLevel      string          `json:"pendingLevel,omitempty"`
	PendingLevelLabel string          `json:"pendingLevelLabel,omitempty"`
	PendingSince      string          `json:"pendingSince,omitempty"`
	PendingAsOf       string          `json:"pendingAsOf,omitempty"`
	PendingDays       int             `json:"pendingDays,omitempty"`
	SignalHealth      string          `json:"signalHealth,omitempty"`
	ExecutionHealth   string          `json:"executionHealth,omitempty"`
	SignalAsOf        string          `json:"signalAsOf,omitempty"`
	ExecutionAsOf     string          `json:"executionAsOf,omitempty"`
	BlockingReasons   []string        `json:"blockingReasons,omitempty"`
	QualityUpdatedAt  string          `json:"qualityUpdatedAt,omitempty"`
	Metrics           []ETFRuleMetric `json:"metrics,omitempty"`
	Sources           []ETFRuleSource `json:"sources,omitempty"`
}

type ETFRuleMetric struct {
	Key               string   `json:"key"`
	Label             string   `json:"label"`
	Value             *float64 `json:"value,omitempty"`
	Unit              string   `json:"unit,omitempty"`
	AsOf              string   `json:"asOf,omitempty"`
	Available         bool     `json:"available"`
	Error             string   `json:"error,omitempty"`
	SourceTier        string   `json:"sourceTier,omitempty"`
	QualityState      string   `json:"qualityState,omitempty"`
	QualityMessage    string   `json:"qualityMessage,omitempty"`
	FetchedAt         string   `json:"fetchedAt,omitempty"`
	MaxAgeSeconds     int      `json:"maxAgeSeconds,omitempty"`
	SourceIDs         []string `json:"sourceIds,omitempty"`
	ValidationValue   *float64 `json:"validationValue,omitempty"`
	ValidationSource  string   `json:"validationSource,omitempty"`
	ConflictTolerance float64  `json:"conflictTolerance,omitempty"`
}

type ETFRuleSource struct {
	Name         string `json:"name"`
	URL          string `json:"url,omitempty"`
	Tier         string `json:"tier,omitempty"`
	QualityState string `json:"qualityState,omitempty"`
	AsOf         string `json:"asOf,omitempty"`
}

const (
	etfRuleDailyMetricMaxAgeDays      = 7
	etfRuleWeeklyMetricMaxAgeDays     = 14
	etfRuleMonthlyMetricMaxAgeDays    = 60
	etfRuleRuntimeTimestampDateLayout = "2006-01-02"
	historyOfMarketSP500PEURL         = "https://historyofmarket.com/api/sp500/pe.json"
	historyOfMarketNDXForwardPEURL    = "https://historyofmarket.com/api/ndx/forward-pe.json"
	multplShillerCAPEURL              = "https://www.multpl.com/shiller-pe/table/by-year"
	fundDBAPIHost                     = "https://api.jiucaishuo.com"
	fundDBIndexPageURL                = "https://funddb.cn/site/index"
	fundDBAPIVersion                  = "2.2.7"
	fundDBAPIReqKey                   = "EWf45rlv#kfsr@k#gfksgkr"
	primaryValuationMaxLagDays        = 3
	chinaBondHistoryURL               = "https://yield.chinabond.com.cn/cbweb-pbc-web/pbc/historyQuery"
	eastmoneyTreasuryYieldURL         = "https://datacenter-web.eastmoney.com/api/data/get"
	spChinaLowVolDividendIndexURL     = "https://www.spglobal.com/spdji/en/indices/dividends-factors/sp-china-a-share-largecap-low-volatility-high-dividend-50-index/"
)

type etfRuleLevel struct {
	Key   string
	Label string
}

type etfRuleConfig struct {
	Symbol              string
	Name                string
	PriceSymbol         string
	PriceSourceName     string
	PriceSourceURL      string
	ValuationMetricKey  string
	ValuationMetricName string
	ValuationSourceName string
	ValuationSourceURL  string
	Levels              map[string]etfRuleLevel
	Monthly             map[string]float64
	Weekly              map[string]float64
	Evaluate            func(etfRuleInputs) etfRuleEvaluation
}

type etfRuleInputs struct {
	Drawdown                      *float64
	DrawdownAsOf                  string
	ValuationPercentile           *float64
	ValuationZScore               *float64
	DividendYield                 *float64
	DividendYieldPercentile       *float64
	DividendSpreadPercentile      *float64
	PBPercentile                  *float64
	ValuationScore                *float64
	EarningsYieldSpreadPercentile *float64
	ValuationAsOf                 string
}

type etfRuleEvaluation struct {
	Level    string
	Complete bool
	Reason   string
}

var etfRuleLevels = map[string]etfRuleLevel{
	"quarter": {Key: "quarter", Label: "高估"},
	"half":    {Key: "half", Label: "偏高"},
	"one":     {Key: "one", Label: "中性"},
	"oneHalf": {Key: "oneHalf", Label: "偏低"},
	"two":     {Key: "two", Label: "低估"},
}

var etfRuleConfigs = []etfRuleConfig{
	{
		Symbol:              "022434",
		Name:                "南方中证A500ETF联接A",
		PriceSymbol:         a500TotalReturnIndexCode,
		PriceSourceName:     "中证指数中证A500全收益指数官方日线",
		PriceSourceURL:      a500TotalReturnIndexURL,
		ValuationMetricKey:  "pePercentile",
		ValuationMetricName: "中证A500滚动PE扩展窗口分位",
		ValuationSourceName: "中证指数中证A500滚动PE官方序列",
		ValuationSourceURL:  a500PEHistoryURL,
		Levels:              etfRuleLevels,
		Monthly:             fixedETFRuleAmounts(5000),
		Weekly:              fixedETFRuleAmounts(1250),
		Evaluate:            evaluateA500Rule,
	},
	{
		Symbol:              "018738",
		Name:                "博时标普500ETF联接E(人民币)",
		PriceSymbol:         sp500TotalReturnSymbol,
		PriceSourceName:     "S&P 500 Total Return（SPTR）公开日线",
		PriceSourceURL:      sp500TotalReturnDataURL,
		ValuationMetricKey:  "forwardPEPercentile",
		ValuationMetricName: "标普500未来PE十年分位",
		ValuationSourceName: "History of Market S&P 500 Forward PE（同一序列计算PE与盈利利差分位）",
		ValuationSourceURL:  sp500ForwardPEURL,
		Levels:              etfRuleLevels,
		Monthly:             fixedETFRuleAmounts(5000),
		Weekly:              fixedETFRuleAmounts(1250),
		Evaluate:            evaluateSP500Rule,
	},
	{
		Symbol:              "008163",
		Name:                "南方标普红利低波50ETF联接A",
		PriceSymbol:         "515450.SH",
		PriceSourceName:     "东方财富515450单位净值 + 每份分红（分红再投资总回报）",
		PriceSourceURL:      "https://fund.eastmoney.com/515450.html",
		ValuationMetricKey:  "valuationScore",
		ValuationMetricName: "红利低波估值得分V",
		ValuationSourceName: "南方基金515450申购赎回篮子 + 东方财富成分股PB与分红（场内代理）",
		ValuationSourceURL:  southernETFPCFPageURL,
		Levels:              etfRuleLevels,
		Monthly:             fixedETFRuleAmounts(5000),
		Weekly:              fixedETFRuleAmounts(1250),
		Evaluate:            evaluateDividendLowVolRule,
	},
	{
		Symbol:              "021000",
		Name:                "南方纳斯达克100指数发起(QDII)I",
		PriceSymbol:         "XNDX",
		PriceSourceName:     "Nasdaq官方纳斯达克100总收益指数日线",
		PriceSourceURL:      "https://indexes.nasdaq.com/Index/Overview/XNDX",
		ValuationMetricKey:  "forwardPEPercentile",
		ValuationMetricName: "纳指100未来PE十年分位",
		ValuationSourceName: "History of Market Nasdaq 100 Forward PE（同一序列计算PE与盈利利差分位）",
		ValuationSourceURL:  historyOfMarketNDXForwardPEURL,
		Levels:              etfRuleLevels,
		Monthly:             fixedETFRuleAmounts(5000),
		Weekly:              fixedETFRuleAmounts(1250),
		Evaluate:            evaluateNasdaq100Rule,
	},
}

func fixedETFRuleAmounts(amount float64) map[string]float64 {
	return map[string]float64{
		"quarter": amount,
		"half":    amount,
		"one":     amount,
		"oneHalf": amount,
		"two":     amount,
	}
}

func updateETFRuleStatuses(client *http.Client, now time.Time) ([]ETFRuleStatus, []QuoteSkip) {
	statuses := make([]ETFRuleStatus, 0, len(etfRuleConfigs))
	skipped := []QuoteSkip{}
	for _, config := range etfRuleConfigs {
		status, err := fetchETFRuleStatus(client, config, now)
		if err != nil {
			skipped = append(skipped, QuoteSkip{Type: "etf-rule", Symbol: config.Symbol, Name: config.Name, Error: err.Error()})
		}
		statuses = append(statuses, status)
	}
	return statuses, skipped
}

func fetchETFRuleStatus(client *http.Client, config etfRuleConfig, now time.Time) (ETFRuleStatus, error) {
	inputs := etfRuleInputs{}
	metrics := []ETFRuleMetric{}
	sources := []ETFRuleSource{{Name: config.PriceSourceName, URL: config.PriceSourceURL}}
	statusErrs := []string{}

	var nasdaqSnapshot nasdaqOpportunitySnapshot
	var nasdaqIssues nasdaqOpportunityErrors
	var sp500Snapshot sp500OpportunitySnapshot
	var sp500Issues sp500OpportunityErrors
	var a500Snapshot a500OpportunitySnapshot
	var a500Issues a500OpportunityErrors
	if config.Symbol == "022434" {
		a500Snapshot, a500Issues = fetchA500OpportunitySnapshot(client, now)
	}
	if config.Symbol == "018738" {
		sp500Snapshot, sp500Issues = fetchSP500OpportunitySnapshot(client, now)
	}
	if config.Symbol == "022434" {
		if a500Issues.Drawdown != nil {
			statusErrs = append(statusErrs, "中证A500全收益回撤："+a500Issues.Drawdown.Error())
			metrics = append(metrics, ETFRuleMetric{Key: "drawdown3y", Label: "中证A500全收益高点回撤", Unit: "%", Available: false, Error: a500Issues.Drawdown.Error()})
		} else {
			inputs.Drawdown = &a500Snapshot.Drawdown
			inputs.DrawdownAsOf = a500Snapshot.Date
			metrics = append(metrics,
				ETFRuleMetric{Key: "drawdown3y", Label: "中证A500全收益高点回撤", Value: percentMetric(a500Snapshot.Drawdown), Unit: "%", AsOf: a500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "totalReturnClose", Label: "中证A500全收益指数", Value: floatMetric(a500Snapshot.IndexClose), AsOf: a500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "totalReturnPeak", Label: "本轮全收益高点", Value: floatMetric(a500Snapshot.PeakClose), AsOf: a500Snapshot.PeakDate, Available: true},
			)
		}
	} else if config.Symbol == "021000" {
		nasdaqSnapshot, nasdaqIssues = fetchNasdaqOpportunitySnapshot(client, now)
		if nasdaqIssues.Drawdown != nil {
			statusErrs = append(statusErrs, "XNDX全收益回撤："+nasdaqIssues.Drawdown.Error())
			metrics = append(metrics, ETFRuleMetric{Key: "drawdown3y", Label: "XNDX十年全收益高点回撤", Unit: "%", Available: false, Error: nasdaqIssues.Drawdown.Error()})
		} else {
			inputs.Drawdown = &nasdaqSnapshot.Drawdown
			inputs.DrawdownAsOf = nasdaqSnapshot.Date
			metrics = append(metrics, ETFRuleMetric{Key: "drawdown3y", Label: "XNDX十年全收益高点回撤", Value: percentMetric(nasdaqSnapshot.Drawdown), Unit: "%", AsOf: nasdaqSnapshot.Date, Available: true})
		}
	} else if config.Symbol == "018738" {
		drawdownLabel := "SPTR全收益高点回撤"
		if !sp500Snapshot.DirectSPTR && strings.TrimSpace(sp500Snapshot.IndexSource) != "" {
			drawdownLabel = "标普500总回报回撤（SPY备援）"
		}
		if sp500Issues.Drawdown != nil {
			statusErrs = append(statusErrs, "SPTR全收益回撤："+sp500Issues.Drawdown.Error())
			metrics = append(metrics, ETFRuleMetric{Key: "drawdown3y", Label: drawdownLabel, Unit: "%", Available: false, Error: sp500Issues.Drawdown.Error()})
		} else {
			inputs.Drawdown = &sp500Snapshot.Drawdown
			inputs.DrawdownAsOf = sp500Snapshot.Date
			metrics = append(metrics, ETFRuleMetric{Key: "drawdown3y", Label: drawdownLabel, Value: percentMetric(sp500Snapshot.Drawdown), Unit: "%", AsOf: sp500Snapshot.Date, Available: true})
		}
	} else {
		drawdown, drawdownDate, err := fetchETFRuleDrawdown(client, config)
		drawdownLabel := "近3年总收益回撤"
		if config.Symbol == "008163" {
			drawdownLabel = "515450成立以来总回报回撤"
		}
		if err != nil {
			statusErrs = append(statusErrs, "回撤："+err.Error())
			metrics = append(metrics, ETFRuleMetric{Key: "drawdown3y", Label: drawdownLabel, Unit: "%", Available: false, Error: err.Error()})
		} else {
			inputs.Drawdown = &drawdown
			inputs.DrawdownAsOf = drawdownDate
			metrics = append(metrics, ETFRuleMetric{Key: "drawdown3y", Label: drawdownLabel, Value: percentMetric(drawdown), Unit: "%", AsOf: drawdownDate, Available: true})
		}
	}

	if config.Symbol == "022434" {
		if a500Issues.Valuation != nil {
			for _, metric := range []ETFRuleMetric{
				{Key: "indexPE", Label: "中证A500滚动PE", Available: false},
				{Key: "pePercentile", Label: "滚动PE扩展窗口分位", Unit: "%", Available: false},
				{Key: "china10YBondYield", Label: "中债10年期国债收益率", Unit: "%", Available: false},
				{Key: "earningsYieldSpread", Label: "A500股债利差", Unit: "%", Available: false},
				{Key: "earningsYieldSpreadPercentile", Label: "股债利差扩展窗口分位", Unit: "%", Available: false},
			} {
				metric.Error = a500Issues.Valuation.Error()
				metrics = append(metrics, metric)
			}
		} else {
			inputs.ValuationPercentile = &a500Snapshot.PEPercentile
			inputs.EarningsYieldSpreadPercentile = &a500Snapshot.SpreadPercentile
			inputs.ValuationAsOf = a500Snapshot.Date
			metrics = append(metrics,
				ETFRuleMetric{Key: "indexPE", Label: "中证A500滚动PE", Value: floatMetric(a500Snapshot.PE), AsOf: a500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "pePercentile", Label: fmt.Sprintf("滚动PE扩展窗口分位（%d日）", a500Snapshot.PEObservationCount), Value: percentMetric(a500Snapshot.PEPercentile), Unit: "%", AsOf: a500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "china10YBondYield", Label: "中债10年期国债收益率", Value: percentMetric(a500Snapshot.China10YBondYield), Unit: "%", AsOf: a500Snapshot.BondDate, Available: true},
				ETFRuleMetric{Key: "earningsYieldSpread", Label: "A500股债利差", Value: percentMetric(a500Snapshot.EarningsYieldSpread), Unit: "%", AsOf: a500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "earningsYieldSpreadPercentile", Label: fmt.Sprintf("股债利差扩展窗口分位（%d日）", a500Snapshot.SpreadObservationCount), Value: percentMetric(a500Snapshot.SpreadPercentile), Unit: "%", AsOf: a500Snapshot.Date, Available: true},
			)
		}
		if a500Issues.Panic != nil {
			for _, metric := range []ETFRuleMetric{
				{Key: "rv20", Label: "20日实现波动率", Unit: "%", Available: false},
				{Key: "rv20Percentile", Label: "RV20五年分位", Unit: "%", Available: false},
				{Key: "fiveDayReturn", Label: "近5个交易日全收益", Unit: "%", Available: false},
				{Key: "volumeRatio", Label: "成交额/20日均值", Unit: "倍", Available: false},
			} {
				metric.Error = a500Issues.Panic.Error()
				metrics = append(metrics, metric)
			}
		} else {
			metrics = append(metrics,
				ETFRuleMetric{Key: "rv20", Label: "20日实现波动率", Value: percentMetric(a500Snapshot.RV20), Unit: "%", AsOf: a500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "rv20Percentile", Label: fmt.Sprintf("RV20五年分位（%d日）", a500Snapshot.RV20ObservationCount), Value: percentMetric(a500Snapshot.RV20Percentile), Unit: "%", AsOf: a500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "fiveDayReturn", Label: "近5个交易日全收益", Value: percentMetric(a500Snapshot.FiveDayReturn), Unit: "%", AsOf: a500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "volumeRatio", Label: "成交额/20日均值", Value: floatMetric(a500Snapshot.VolumeRatio), Unit: "倍", AsOf: a500Snapshot.Date, Available: a500Snapshot.VolumeRatio > 0},
			)
		}
		metrics = append(metrics, ETFRuleMetric{Key: "breadthBelowMA20", Label: "成份股低于20日均线比例", Unit: "%", Available: false, Error: "免费官方成份股广度序列未接入，恐慌系数不使用此条件"})
		if a500Issues.Trading != nil {
			for _, metric := range []ETFRuleMetric{
				{Key: "tacticalMarketPrice", Label: a500TacticalETFCode + "场内价格", Available: false},
				{Key: "tacticalEstimatedNAV", Label: a500TacticalETFCode + "估算净值", Available: false},
				{Key: "etfPremium", Label: a500TacticalETFCode + "估算溢价", Unit: "%", Available: false},
				{Key: "bidAskSpread", Label: a500TacticalETFCode + "买卖价差", Unit: "%", Available: false},
				{Key: "openingGap", Label: a500TacticalETFCode + "开盘涨跌", Unit: "%", Available: false},
			} {
				metric.Error = a500Issues.Trading.Error()
				metrics = append(metrics, metric)
			}
		} else {
			metrics = append(metrics,
				ETFRuleMetric{Key: "tacticalMarketPrice", Label: a500TacticalETFCode + "场内价格", Value: floatMetric(a500Snapshot.MarketPrice), AsOf: a500Snapshot.MarketPriceDate, Available: true},
				ETFRuleMetric{Key: "tacticalOfficialNAV", Label: a500TacticalETFCode + "最近公布净值", Value: floatMetric(a500Snapshot.OfficialNAV), AsOf: a500Snapshot.OfficialNAVDate, Available: true},
				ETFRuleMetric{Key: "tacticalEstimatedNAV", Label: a500TacticalETFCode + "估算净值", Value: floatMetric(a500Snapshot.EstimatedNAV), AsOf: a500Snapshot.MarketPriceDate, Available: true},
				ETFRuleMetric{Key: "etfPremium", Label: a500TacticalETFCode + "估算溢价", Value: percentMetric(a500Snapshot.Premium), Unit: "%", AsOf: a500Snapshot.MarketPriceDate, Available: true},
				ETFRuleMetric{Key: "bidAskSpread", Label: a500TacticalETFCode + "买卖价差", Value: percentMetric(a500Snapshot.BidAskSpread), Unit: "%", AsOf: a500Snapshot.MarketPriceDate, Available: true},
				ETFRuleMetric{Key: "openingGap", Label: a500TacticalETFCode + "开盘涨跌", Value: percentMetric(a500Snapshot.OpeningGap), Unit: "%", AsOf: a500Snapshot.MarketPriceDate, Available: true},
			)
		}
	} else if config.Symbol == "008163" {
		valuation, valuationErr := fetchDividendLowVolIndexValuation(client, now)
		if valuationErr != nil {
			statusErrs = append(statusErrs, "红利低波场内估值代理："+valuationErr.Error())
			errorText := "515450篮子估值代理不可用：" + valuationErr.Error()
			metrics = append(metrics,
				ETFRuleMetric{Key: "dividendYield", Label: "515450篮子TTM股息率（代理）", Unit: "%", Available: false, Error: errorText},
				ETFRuleMetric{Key: "china10YBondYield", Label: "中债10年期国债收益率", Unit: "%", Available: false, Error: errorText},
				ETFRuleMetric{Key: "dividendSpread", Label: "篮子股债利差（代理）", Unit: "%", Available: false, Error: errorText},
				ETFRuleMetric{Key: "dividendSpreadPercentile", Label: "篮子股债利差5年分位（代理）", Unit: "%", Available: false, Error: errorText},
				ETFRuleMetric{Key: "indexPB", Label: "515450篮子PB（代理）", Available: false, Error: errorText},
				ETFRuleMetric{Key: "pbPercentile", Label: "篮子PB五年历史分位（代理）", Unit: "%", Available: false, Error: errorText},
				ETFRuleMetric{Key: "valuationScore", Label: "红利估值得分V", Unit: "%", Available: false, Error: errorText},
				ETFRuleMetric{Key: "basketCoverage", Label: "篮子有效覆盖", Unit: "%", Available: false, Error: errorText},
			)
		} else {
			inputs.DividendYield = &valuation.DividendYield
			inputs.DividendSpreadPercentile = &valuation.SpreadPercentile
			inputs.PBPercentile = &valuation.PBPercentile
			inputs.ValuationScore = &valuation.ValuationScore
			inputs.ValuationAsOf = valuation.Date
			metrics = append(metrics,
				ETFRuleMetric{Key: "dividendYield", Label: "515450篮子TTM股息率（代理）", Value: percentMetric(valuation.DividendYield), Unit: "%", AsOf: valuation.Date, Available: true},
				ETFRuleMetric{Key: "china10YBondYield", Label: "中债10年期国债收益率", Value: percentMetric(valuation.BondYield), Unit: "%", AsOf: valuation.BondDate, Available: true},
				ETFRuleMetric{Key: "dividendSpread", Label: "篮子股债利差（代理）", Value: percentMetric(valuation.Spread), Unit: "%", AsOf: valuation.Date, Available: true},
				ETFRuleMetric{Key: "dividendSpreadPercentile", Label: "篮子股债利差5年分位（代理）", Value: percentMetric(valuation.SpreadPercentile), Unit: "%", AsOf: valuation.Date, Available: true},
				ETFRuleMetric{Key: "indexPB", Label: "515450篮子PB（代理）", Value: floatMetric(valuation.PB), AsOf: valuation.Date, Available: true},
				ETFRuleMetric{Key: "pbPercentile", Label: "篮子PB五年历史分位（代理）", Value: percentMetric(valuation.PBPercentile), Unit: "%", AsOf: valuation.Date, Available: true},
				ETFRuleMetric{Key: "valuationScore", Label: "红利估值得分V", Value: percentMetric(valuation.ValuationScore), Unit: "%", AsOf: valuation.Date, Available: true},
				ETFRuleMetric{Key: "basketCoverage", Label: fmt.Sprintf("篮子有效覆盖（%d/%d）", valuation.ValidComponentCount, valuation.ComponentCount), Value: percentMetric(valuation.Coverage), Unit: "%", AsOf: valuation.Date, Available: true},
			)
		}
		market, marketErr := fetchEastmoneyQuote(client, "515450.SH")
		if marketErr != nil || market.Price <= 0 {
			errorText := "515450场内行情暂不可用"
			if marketErr != nil {
				errorText = marketErr.Error()
			}
			metrics = append(metrics, ETFRuleMetric{Key: "tacticalMarketPrice", Label: "515450场内价格", Available: false, Error: errorText})
			statusErrs = append(statusErrs, "515450场内行情："+errorText)
		} else {
			metrics = append(metrics, ETFRuleMetric{Key: "tacticalMarketPrice", Label: "515450场内价格", Value: floatMetric(market.Price), AsOf: market.PriceDate, Available: true})
		}
	} else if config.Symbol == "018738" {
		if sp500Issues.Valuation != nil {
			statusErrs = append(statusErrs, "标普估值："+sp500Issues.Valuation.Error())
			for _, metric := range []ETFRuleMetric{
				{Key: "forwardPE", Label: "标普500未来12个月PE", Available: false},
				{Key: "forwardPEPercentile", Label: "未来PE十年周度分位", Unit: "%", Available: false},
				{Key: "us10YBondYield", Label: "美国10年期国债收益率", Unit: "%", Available: false},
				{Key: "earningsYieldSpread", Label: "盈利收益率利差", Unit: "%", Available: false},
				{Key: "earningsYieldSpreadPercentile", Label: "盈利利差十年周度分位", Unit: "%", Available: false},
			} {
				metric.Error = sp500Issues.Valuation.Error()
				metrics = append(metrics, metric)
			}
		} else {
			inputs.ValuationPercentile = &sp500Snapshot.ForwardPEPercentile
			inputs.EarningsYieldSpreadPercentile = &sp500Snapshot.SpreadPercentile
			inputs.ValuationAsOf = sp500Snapshot.ValuationDate
			metrics = append(metrics,
				ETFRuleMetric{Key: "forwardPE", Label: "标普500未来12个月PE", Value: floatMetric(sp500Snapshot.ForwardPE), AsOf: sp500Snapshot.ValuationDate, Available: true},
				ETFRuleMetric{Key: "forwardPEPercentile", Label: "未来PE十年周度分位", Value: percentMetric(sp500Snapshot.ForwardPEPercentile), Unit: "%", AsOf: sp500Snapshot.ValuationDate, Available: true},
				ETFRuleMetric{Key: "us10YBondYield", Label: "美国10年期国债收益率", Value: percentMetric(sp500Snapshot.US10YBondYield), Unit: "%", AsOf: sp500Snapshot.ValuationDate, Available: true},
				ETFRuleMetric{Key: "earningsYieldSpread", Label: "盈利收益率利差", Value: percentMetric(sp500Snapshot.EarningsYieldSpread), Unit: "%", AsOf: sp500Snapshot.ValuationDate, Available: true},
				ETFRuleMetric{Key: "earningsYieldSpreadPercentile", Label: fmt.Sprintf("盈利利差十年周度分位（%d周）", sp500Snapshot.ValuationObservationCnt), Value: percentMetric(sp500Snapshot.SpreadPercentile), Unit: "%", AsOf: sp500Snapshot.ValuationDate, Available: true},
			)
		}
		metrics = append(metrics, ETFRuleMetric{
			Key:       "forwardEarningsRevision3m",
			Label:     "未来盈利预期三个月修正",
			Unit:      "%",
			Available: false,
			Error:     sp500Issues.EarningsRevision.Error(),
		})
		if sp500Issues.VIX != nil {
			statusErrs = append(statusErrs, "VIX："+sp500Issues.VIX.Error())
			metrics = append(metrics, ETFRuleMetric{Key: "vix", Label: "Cboe标普500波动率VIX", Available: false, Error: sp500Issues.VIX.Error()})
		} else {
			metrics = append(metrics, ETFRuleMetric{Key: "vix", Label: "Cboe标普500波动率VIX", Value: floatMetric(sp500Snapshot.VIX), AsOf: sp500Snapshot.VIXDate, Available: true})
		}
		if sp500Issues.CNYDrawdown != nil {
			statusErrs = append(statusErrs, "人民币全收益回撤："+sp500Issues.CNYDrawdown.Error())
			metrics = append(metrics, ETFRuleMetric{Key: "cnyTotalReturnDrawdown", Label: "人民币口径全收益回撤", Unit: "%", Available: false, Error: sp500Issues.CNYDrawdown.Error()})
		} else {
			metrics = append(metrics,
				ETFRuleMetric{Key: "cnyTotalReturnDrawdown", Label: "人民币口径全收益回撤", Value: percentMetric(sp500Snapshot.CNYDrawdown), Unit: "%", AsOf: sp500Snapshot.Date, Available: true},
				ETFRuleMetric{Key: "usdCny", Label: "USD/CNY参考汇率", Value: floatMetric(sp500Snapshot.USDToCNY), AsOf: sp500Snapshot.FXDate, Available: true},
			)
		}
		if sp500Issues.Premium != nil {
			statusErrs = append(statusErrs, sp500TacticalETFCode+"估算溢价："+sp500Issues.Premium.Error())
			for _, metric := range []ETFRuleMetric{
				{Key: "sp500FuturesChange", Label: "标普500期货当日变动", Unit: "%", Available: false},
				{Key: "tacticalMarketPrice", Label: sp500TacticalETFCode + "场内价格", Available: false},
				{Key: "tacticalEstimatedNAV", Label: sp500TacticalETFCode + "估算实时净值", Available: false},
				{Key: "qdiiPremium", Label: sp500TacticalETFCode + "估算溢价", Unit: "%", Available: false},
			} {
				metric.Error = sp500Issues.Premium.Error()
				metrics = append(metrics, metric)
			}
		} else {
			metrics = append(metrics,
				ETFRuleMetric{Key: "sp500FuturesChange", Label: "标普500期货当日变动", Value: percentMetric(sp500Snapshot.FuturesChange), Unit: "%", AsOf: sp500Snapshot.FuturesDate, Available: true},
				ETFRuleMetric{Key: "tacticalMarketPrice", Label: sp500TacticalETFCode + "场内价格", Value: floatMetric(sp500Snapshot.MarketPrice), AsOf: sp500Snapshot.MarketPriceDate, Available: true},
				ETFRuleMetric{Key: "tacticalOfficialNAV", Label: sp500TacticalETFCode + "最近公布净值", Value: floatMetric(sp500Snapshot.OfficialNAV), AsOf: sp500Snapshot.OfficialNAVDate, Available: true},
				ETFRuleMetric{Key: "tacticalEstimatedNAV", Label: sp500TacticalETFCode + "估算实时净值", Value: floatMetric(sp500Snapshot.EstimatedNAV), AsOf: sp500Snapshot.MarketPriceDate, Available: true},
				ETFRuleMetric{Key: "qdiiPremium", Label: sp500TacticalETFCode + "估算溢价", Value: percentMetric(sp500Snapshot.Premium), Unit: "%", AsOf: sp500Snapshot.MarketPriceDate, Available: true},
			)
		}
	} else if config.Symbol == "021000" {
		if nasdaqIssues.Valuation != nil {
			statusErrs = append(statusErrs, "纳指估值："+nasdaqIssues.Valuation.Error())
			for _, metric := range []ETFRuleMetric{
				{Key: "forwardPE", Label: "纳指100未来12个月PE", Available: false},
				{Key: "forwardPEPercentile", Label: "未来PE十年周度分位", Unit: "%", Available: false},
				{Key: "us10YBondYield", Label: "美国10年期国债收益率", Unit: "%", Available: false},
				{Key: "earningsYieldSpread", Label: "盈利收益率利差", Unit: "%", Available: false},
				{Key: "earningsYieldSpreadPercentile", Label: "盈利利差十年周度分位", Unit: "%", Available: false},
			} {
				metric.Error = nasdaqIssues.Valuation.Error()
				metrics = append(metrics, metric)
			}
		} else {
			inputs.ValuationPercentile = &nasdaqSnapshot.ForwardPEPercentile
			inputs.EarningsYieldSpreadPercentile = &nasdaqSnapshot.SpreadPercentile
			inputs.ValuationAsOf = nasdaqSnapshot.ValuationDate
			metrics = append(metrics,
				ETFRuleMetric{Key: "forwardPE", Label: "纳指100未来12个月PE", Value: floatMetric(nasdaqSnapshot.ForwardPE), AsOf: nasdaqSnapshot.ValuationDate, Available: true},
				ETFRuleMetric{Key: "forwardPEPercentile", Label: "未来PE十年周度分位", Value: percentMetric(nasdaqSnapshot.ForwardPEPercentile), Unit: "%", AsOf: nasdaqSnapshot.ValuationDate, Available: true},
				ETFRuleMetric{Key: "us10YBondYield", Label: "美国10年期国债收益率", Value: percentMetric(nasdaqSnapshot.US10YBondYield), Unit: "%", AsOf: nasdaqSnapshot.ValuationDate, Available: true},
				ETFRuleMetric{Key: "earningsYieldSpread", Label: "盈利收益率利差", Value: percentMetric(nasdaqSnapshot.EarningsYieldSpread), Unit: "%", AsOf: nasdaqSnapshot.ValuationDate, Available: true},
				ETFRuleMetric{Key: "earningsYieldSpreadPercentile", Label: fmt.Sprintf("盈利利差十年周度分位（%d周）", nasdaqSnapshot.ValuationObservationCnt), Value: percentMetric(nasdaqSnapshot.SpreadPercentile), Unit: "%", AsOf: nasdaqSnapshot.ValuationDate, Available: true},
			)
		}
		if nasdaqIssues.VXN != nil {
			statusErrs = append(statusErrs, "VXN："+nasdaqIssues.VXN.Error())
			metrics = append(metrics, ETFRuleMetric{Key: "vxn", Label: "Cboe纳指100波动率VXN", Available: false, Error: nasdaqIssues.VXN.Error()})
		} else {
			metrics = append(metrics, ETFRuleMetric{Key: "vxn", Label: "Cboe纳指100波动率VXN", Value: floatMetric(nasdaqSnapshot.VXN), AsOf: nasdaqSnapshot.VXNDate, Available: true})
		}
		if nasdaqIssues.CNYDrawdown != nil {
			statusErrs = append(statusErrs, "人民币全收益回撤："+nasdaqIssues.CNYDrawdown.Error())
			metrics = append(metrics, ETFRuleMetric{Key: "cnyTotalReturnDrawdown", Label: "人民币口径全收益回撤", Unit: "%", Available: false, Error: nasdaqIssues.CNYDrawdown.Error()})
		} else {
			metrics = append(metrics,
				ETFRuleMetric{Key: "cnyTotalReturnDrawdown", Label: "人民币口径全收益回撤", Value: percentMetric(nasdaqSnapshot.CNYDrawdown), Unit: "%", AsOf: nasdaqSnapshot.Date, Available: true},
				ETFRuleMetric{Key: "usdCny", Label: "USD/CNY参考汇率", Value: floatMetric(nasdaqSnapshot.USDToCNY), AsOf: nasdaqSnapshot.FXDate, Available: true},
			)
		}
		if nasdaqIssues.Premium != nil {
			statusErrs = append(statusErrs, nasdaqTacticalETFCode+"估算溢价："+nasdaqIssues.Premium.Error())
			for _, metric := range []ETFRuleMetric{
				{Key: "nasdaqFuturesChange", Label: "纳指期货当日变动", Unit: "%", Available: false},
				{Key: "tacticalMarketPrice", Label: nasdaqTacticalETFCode + "场内价格", Available: false},
				{Key: "tacticalEstimatedNAV", Label: nasdaqTacticalETFCode + "估算实时净值", Available: false},
				{Key: "qdiiPremium", Label: nasdaqTacticalETFCode + "估算溢价", Unit: "%", Available: false},
			} {
				metric.Error = nasdaqIssues.Premium.Error()
				metrics = append(metrics, metric)
			}
		} else {
			metrics = append(metrics,
				ETFRuleMetric{Key: "nasdaqFuturesChange", Label: "纳指期货当日变动", Value: percentMetric(nasdaqSnapshot.FuturesChange), Unit: "%", AsOf: nasdaqSnapshot.FuturesDate, Available: true},
				ETFRuleMetric{Key: "tacticalMarketPrice", Label: nasdaqTacticalETFCode + "场内价格", Value: floatMetric(nasdaqSnapshot.MarketPrice), AsOf: nasdaqSnapshot.MarketPriceDate, Available: true},
				ETFRuleMetric{Key: "tacticalOfficialNAV", Label: nasdaqTacticalETFCode + "最近公布净值", Value: floatMetric(nasdaqSnapshot.OfficialNAV), AsOf: nasdaqSnapshot.OfficialNAVDate, Available: true},
				ETFRuleMetric{Key: "tacticalEstimatedNAV", Label: nasdaqTacticalETFCode + "估算实时净值", Value: floatMetric(nasdaqSnapshot.EstimatedNAV), AsOf: nasdaqSnapshot.MarketPriceDate, Available: true},
				ETFRuleMetric{Key: "qdiiPremium", Label: nasdaqTacticalETFCode + "估算溢价", Value: percentMetric(nasdaqSnapshot.Premium), Unit: "%", AsOf: nasdaqSnapshot.MarketPriceDate, Available: true},
			)
		}
	} else {
		valuation, valuationErr := fetchETFRuleValuation(client, config)
		if valuationErr != nil {
			statusErrs = append(statusErrs, config.ValuationMetricName+"："+valuationErr.Error())
			metrics = append(metrics, ETFRuleMetric{Key: config.ValuationMetricKey, Label: config.ValuationMetricName, Unit: configValuationMetricUnit(config), Available: false, Error: valuationErr.Error()})
		} else {
			inputs.ValuationAsOf = valuation.Date
			metricValue := valuation.Value
			if valuation.Kind == "zScore" {
				inputs.ValuationZScore = &valuation.Value
			} else {
				inputs.ValuationPercentile = &valuation.Value
				metricValue = valuation.Value * 100
			}
			metrics = append(metrics, ETFRuleMetric{Key: config.ValuationMetricKey, Label: config.ValuationMetricName, Value: floatMetric(metricValue), Unit: valuation.Unit, AsOf: valuation.Date, Available: true})
		}
	}
	if strings.TrimSpace(config.ValuationSourceName) != "" {
		sources = append(sources, ETFRuleSource{Name: config.ValuationSourceName, URL: config.ValuationSourceURL})
	}
	if config.Symbol == "022434" {
		sources = append(sources,
			ETFRuleSource{Name: "中证指数中证A500全收益指数", URL: a500TotalReturnIndexURL},
			ETFRuleSource{Name: "中证指数中证A500滚动PE", URL: a500PEHistoryURL},
			ETFRuleSource{Name: "中债10年期国债收益率（官网校验）", URL: chinaBondHistoryURL},
			ETFRuleSource{Name: "东方财富中债收益率历史序列", URL: "https://data.eastmoney.com/cjsj/zmgzsyl.html"},
			ETFRuleSource{Name: "腾讯" + a500TacticalETFCode + "五档盘口", URL: "http://qt.gtimg.cn/q=sz" + a500TacticalETFCode},
			ETFRuleSource{Name: "东方财富" + a500TacticalETFCode + "基金净值", URL: a500TacticalETFNetValueURL},
		)
	} else if config.Symbol == "008163" {
		sources = append(sources,
			ETFRuleSource{Name: "东方财富A股历史PB与收盘价", URL: eastmoneyValueAnalysisPageURL},
			ETFRuleSource{Name: "东方财富A股分红送配", URL: eastmoneyShareBonusPageURL},
			ETFRuleSource{Name: "标普中国A股大盘红利低波50指数", URL: spChinaLowVolDividendIndexURL},
			ETFRuleSource{Name: "中债10年期国债收益率（官网校验）", URL: chinaBondHistoryURL},
			ETFRuleSource{Name: "东方财富中债收益率历史序列", URL: "https://data.eastmoney.com/cjsj/zmgzsyl.html"},
			ETFRuleSource{Name: "东方财富515450场内行情", URL: "https://quote.eastmoney.com/sh515450.html"},
		)
	} else if config.Symbol == "018738" {
		sources = append(sources,
			ETFRuleSource{Name: "S&P Dow Jones Indices标普500官方说明", URL: sp500OfficialIndexURL},
			ETFRuleSource{Name: "美国财政部10年期国债收益率", URL: nasdaqUS10YTreasuryURL},
			ETFRuleSource{Name: "Cboe VIX官方历史数据", URL: sp500VIXHistoryURL},
			ETFRuleSource{Name: "Frankfurter/欧洲央行USD-CNY参考汇率", URL: nasdaqFXHistoryBaseURL},
			ETFRuleSource{Name: "东方财富" + sp500TacticalETFCode + "场内行情", URL: sp500TacticalETFQuoteURL},
			ETFRuleSource{Name: "东方财富" + sp500TacticalETFCode + "基金净值", URL: sp500TacticalETFNetValuePageURL},
			ETFRuleSource{Name: "新浪标普500期货行情（Yahoo Finance备援，估算净值辅助）", URL: sp500SinaFuturesURL},
		)
		if !sp500Snapshot.DirectSPTR && strings.TrimSpace(sp500Snapshot.IndexSource) != "" {
			sources = append(sources, ETFRuleSource{Name: sp500Snapshot.IndexSource, URL: sp500Snapshot.IndexSourceURL})
		}
	} else if config.Symbol == "021000" {
		sources = append(sources,
			ETFRuleSource{Name: "美国财政部10年期国债收益率", URL: nasdaqUS10YTreasuryURL},
			ETFRuleSource{Name: "Cboe VXN官方历史数据", URL: nasdaqVXNHistoryURL},
			ETFRuleSource{Name: "Frankfurter/欧洲央行USD-CNY参考汇率", URL: nasdaqFXHistoryBaseURL},
			ETFRuleSource{Name: "东方财富" + nasdaqTacticalETFCode + "场内行情", URL: nasdaqTacticalETFQuoteURL},
			ETFRuleSource{Name: "东方财富" + nasdaqTacticalETFCode + "基金净值", URL: nasdaqTacticalETFNetValuePageURL},
			ETFRuleSource{Name: "新浪纳指期货行情（Yahoo Finance备援，估算净值辅助）", URL: nasdaqSinaFuturesURL},
		)
	}

	evaluation := config.Evaluate(inputs)
	level := config.Levels[evaluation.Level]
	status := ETFRuleStatus{
		Symbol:        config.Symbol,
		Name:          config.Name,
		Level:         evaluation.Level,
		LevelLabel:    level.Label,
		MonthlyAmount: config.Monthly["one"],
		WeeklyAmount:  config.Weekly["one"],
		Complete:      evaluation.Complete,
		Reason:        evaluation.Reason,
		AsOf:          firstNonEmpty(inputs.DrawdownAsOf, inputs.ValuationAsOf),
		UpdatedAt:     now.Format("2006-01-02 15:04:05"),
		Metrics:       metrics,
		Sources:       sources,
	}
	if status.Level == "" {
		status.LevelLabel = "待数据"
		status.Reason = firstNonEmpty(evaluation.Reason, strings.Join(statusErrs, "；"))
	}
	if len(statusErrs) > 0 && status.Complete {
		status.Reason = strings.TrimSpace(status.Reason + "；部分辅助指标未取到：" + strings.Join(statusErrs, "；"))
	}
	if status.AsOf == "" {
		status.AsOf = now.Format("2006-01-02")
	}
	status = enforceETFRuleStatusConfidence(status, config, now)
	if len(statusErrs) > 0 {
		return status, errors.New(strings.Join(statusErrs, "；"))
	}
	return status, nil
}

func fetchETFRuleDrawdown(client *http.Client, config etfRuleConfig) (float64, string, error) {
	const tradingDays = 3 * 252
	var (
		closes []dailyClose
		err    error
	)
	switch config.Symbol {
	case "022434":
		closes, err = fetchCSIIndexPerformance(client, "000510CNY010", time.Now().AddDate(-3, 0, 0), time.Now())
	case "018738":
		closes, err = fetchSP500TotalReturnCloses(client)
	case "008163":
		closes, err = fetchDividendLowVolTotalReturnCloses(client, 7*252+80)
	case "021000":
		closes, err = fetchNasdaqIndexHistoryChart(client, "XNDX", time.Now().AddDate(-nasdaqTacticalHistoryYears, 0, -14), time.Now())
	default:
		closes, err = fetchRuleDailyCloses(client, config.PriceSymbol, tradingDays+40)
	}
	if err != nil {
		return 0, "", err
	}
	if config.Symbol == "008163" || config.Symbol == "018738" {
		return drawdownFromRecentHigh(closes, len(closes))
	}
	return drawdownFromRecentHigh(closes, tradingDays)
}

func fetchCSIIndexPerformance(client *http.Client, indexCode string, start time.Time, end time.Time) ([]dailyClose, error) {
	values := url.Values{}
	values.Set("indexCode", strings.TrimSpace(indexCode))
	values.Set("startDate", start.Format("20060102"))
	values.Set("endDate", end.Format("20060102"))
	endpoint := "https://www.csindex.com.cn/csindex-home/perf/index-perf?" + values.Encode()
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.csindex.com.cn/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("CSI index request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var payload struct {
		Code string `json:"code"`
		Data []struct {
			TradeDate string  `json:"tradeDate"`
			Close     float64 `json:"close"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Code != "200" {
		return nil, fmt.Errorf("CSI index response code %s", payload.Code)
	}
	closes := make([]dailyClose, 0, len(payload.Data))
	for _, item := range payload.Data {
		date, err := time.Parse("20060102", strings.TrimSpace(item.TradeDate))
		if err == nil && item.Close > 0 {
			closes = append(closes, dailyClose{Date: date.Format(etfRuleRuntimeTimestampDateLayout), Price: item.Close})
		}
	}
	if len(closes) == 0 {
		return nil, errors.New("CSI index performance is empty")
	}
	sort.Slice(closes, func(i, j int) bool { return closes[i].Date < closes[j].Date })
	return closes, nil
}

func fetchNasdaqIndexHistoryChart(client *http.Client, symbol string, start time.Time, end time.Time) ([]dailyClose, error) {
	values := url.Values{}
	values.Set("id", strings.ToUpper(strings.TrimSpace(symbol)))
	values.Set("startDate", start.Format("2006-01-02"))
	values.Set("endDate", end.Format("2006-01-02"))
	req, err := http.NewRequest(http.MethodPost, "https://indexes.nasdaq.com/Index/HistoryChartData", strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://indexes.nasdaq.com/Index/Overview/"+strings.ToUpper(strings.TrimSpace(symbol)))
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("Nasdaq index history request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var points []struct {
		Timestamp int64   `json:"x"`
		Close     float64 `json:"y"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&points); err != nil {
		return nil, err
	}
	closes := make([]dailyClose, 0, len(points))
	for _, point := range points {
		if point.Timestamp > 0 && point.Close > 0 {
			closes = append(closes, dailyClose{Date: time.UnixMilli(point.Timestamp).UTC().Format(etfRuleRuntimeTimestampDateLayout), Price: point.Close})
		}
	}
	if len(closes) == 0 {
		return nil, errors.New("Nasdaq index history is empty")
	}
	sort.Slice(closes, func(i, j int) bool { return closes[i].Date < closes[j].Date })
	return closes, nil
}

func fetchSPYTotalReturnCloses(client *http.Client, limit int) ([]dailyClose, error) {
	closes, err := fetchNasdaqHistoricalCloses(client, "SPY", "etf", limit)
	if err != nil {
		return nil, err
	}
	events, stateStreetErr := fetchStateStreetSPYDividends(client)
	if stateStreetErr != nil {
		events, err = fetchStockAnalysisDividends(client, "SPY")
		if err != nil {
			return nil, fmt.Errorf("SPY dividend sources unavailable (State Street: %v; StockAnalysis: %v)", stateStreetErr, err)
		}
	}
	return totalReturnCloses(closes, events)
}

func fetchDividendLowVolTotalReturnCloses(client *http.Client, limit int) ([]dailyClose, error) {
	closes, events, err := fetchEastmoneyFundNAVTrend(client, "515450")
	if err != nil {
		return nil, err
	}
	if limit > 0 && len(closes) > limit {
		closes = closes[len(closes)-limit:]
	}
	return totalReturnCloses(closes, events)
}

func fetchEastmoneyFundNAVTrend(client *http.Client, code string) ([]dailyClose, []cashDividendEvent, error) {
	normalized := normalizeFundSymbol(code)
	endpoint := "https://fund.eastmoney.com/pingzhongdata/" + url.PathEscape(normalized) + ".js"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/javascript,text/javascript,*/*")
	req.Header.Set("Referer", "https://fund.eastmoney.com/"+url.PathEscape(normalized)+".html")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, nil, fmt.Errorf("fund NAV trend request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, nil, err
	}
	return parseEastmoneyFundNAVTrend(body)
}

func parseEastmoneyFundNAVTrend(body []byte) ([]dailyClose, []cashDividendEvent, error) {
	const marker = "var Data_netWorthTrend = "
	text := string(body)
	start := strings.Index(text, marker)
	if start < 0 {
		return nil, nil, errors.New("missing fund NAV trend")
	}
	start += len(marker)
	end := strings.Index(text[start:], ";")
	if end < 0 {
		return nil, nil, errors.New("incomplete fund NAV trend")
	}
	var points []struct {
		Timestamp int64   `json:"x"`
		NAV       float64 `json:"y"`
		UnitMoney string  `json:"unitMoney"`
	}
	if err := json.Unmarshal([]byte(text[start:start+end]), &points); err != nil {
		return nil, nil, err
	}
	closes := make([]dailyClose, 0, len(points))
	events := []cashDividendEvent{}
	for _, point := range points {
		if point.Timestamp <= 0 || point.NAV <= 0 {
			continue
		}
		date := time.UnixMilli(point.Timestamp).In(time.FixedZone("Asia/Shanghai", 8*60*60)).Format(etfRuleRuntimeTimestampDateLayout)
		closes = append(closes, dailyClose{Date: date, Price: point.NAV})
		if amount, err := fundDividendAmount(point.UnitMoney); err == nil && amount > 0 {
			events = append(events, cashDividendEvent{Date: date, Amount: amount})
		}
	}
	if len(closes) < 2 {
		return nil, nil, errors.New("insufficient fund NAV trend")
	}
	sort.Slice(closes, func(i, j int) bool { return closes[i].Date < closes[j].Date })
	return closes, events, nil
}

func totalReturnCloses(closes []dailyClose, events []cashDividendEvent) ([]dailyClose, error) {
	if len(closes) == 0 {
		return nil, errors.New("missing close prices")
	}
	ordered := append([]dailyClose(nil), closes...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Date < ordered[j].Date })
	dividends := map[string]float64{}
	for _, event := range events {
		if event.Date != "" && event.Amount > 0 {
			dividends[event.Date] += event.Amount
		}
	}
	if ordered[0].Price <= 0 {
		return nil, errors.New("invalid initial close price")
	}
	result := make([]dailyClose, 0, len(ordered))
	value := ordered[0].Price
	result = append(result, dailyClose{Date: ordered[0].Date, Price: value})
	for i := 1; i < len(ordered); i++ {
		previous := ordered[i-1].Price
		current := ordered[i].Price
		if previous <= 0 || current <= 0 {
			continue
		}
		value *= (current + dividends[ordered[i].Date]) / previous
		result = append(result, dailyClose{Date: ordered[i].Date, Price: value})
	}
	if len(result) < 2 {
		return nil, errors.New("insufficient total-return closes")
	}
	return result, nil
}

func fetchTencentRawDailyCloses(client *http.Client, symbol string, limit int) ([]dailyClose, error) {
	sourceSymbol, _, err := tencentSymbol(symbol)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 30
	}
	const batchSize = 640
	closes := make([]dailyClose, 0, limit)
	endDate := ""
	for len(closes) < limit {
		batchLimit := min(batchSize, limit-len(closes))
		batch, err := fetchTencentRawDailyCloseBatch(client, sourceSymbol, endDate, batchLimit)
		if err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}
		closes = append(batch, closes...)
		if len(batch) < batchLimit {
			break
		}
		oldest, err := time.Parse(etfRuleRuntimeTimestampDateLayout, batch[0].Date)
		if err != nil {
			return nil, err
		}
		endDate = oldest.AddDate(0, 0, -1).Format(etfRuleRuntimeTimestampDateLayout)
	}
	if len(closes) == 0 {
		return nil, errors.New("missing tencent raw close prices")
	}
	if len(closes) > limit {
		closes = closes[len(closes)-limit:]
	}
	return closes, nil
}

func fetchTencentRawDailyCloseBatch(client *http.Client, sourceSymbol string, endDate string, limit int) ([]dailyClose, error) {
	param := fmt.Sprintf("%s,day,,%s,%d", sourceSymbol, strings.TrimSpace(endDate), limit)
	endpoint := "https://web.ifzq.gtimg.cn/appstock/app/kline/kline?param=" + url.QueryEscape(param)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json,text/plain,*/*")
	req.Header.Set("Referer", "https://gu.qq.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("tencent raw kline request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	var payload struct {
		Code int                        `json:"code"`
		Data map[string]json.RawMessage `json:"data"`
		Msg  string                     `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Code != 0 {
		return nil, fmt.Errorf("tencent raw kline response: %s", payload.Msg)
	}
	raw, ok := payload.Data[sourceSymbol]
	if !ok {
		return nil, fmt.Errorf("missing tencent raw kline symbol: %s", sourceSymbol)
	}
	var symbolPayload map[string]json.RawMessage
	if err := json.Unmarshal(raw, &symbolPayload); err != nil {
		return nil, err
	}
	return parseTencentKlineRows(symbolPayload["day"])
}

func fetchRuleDailyCloses(client *http.Client, symbol string, limit int) ([]dailyClose, error) {
	normalized := normalizeSymbol(symbol)
	if strings.HasSuffix(normalized, ".SH") || strings.HasSuffix(normalized, ".SZ") || strings.HasSuffix(normalized, ".HK") {
		closes, err := fetchTencentDailyCloses(client, normalized, limit)
		if err == nil && len(closes) > 0 {
			return closes, nil
		}
		return fetchEastmoneyDailyCloses(client, normalized, limit)
	}
	if secID := eastmoneyGlobalIndexSecID(normalized); secID != "" {
		return fetchGlobalIndexDailyCloses(client, normalized, secID, limit)
	}
	closes, err := fetchYahooDailyCloses(client, normalized, "1y")
	if err == nil && len(closes) > 0 {
		return closes, nil
	}
	if strings.EqualFold(normalized, "^GSPC") || strings.EqualFold(normalized, "SPX") {
		nasdaqCloses, nasdaqErr := fetchNasdaqHistoricalCloses(client, "SPY", "etf", limit)
		if nasdaqErr == nil && len(nasdaqCloses) > 0 {
			return nasdaqCloses, nil
		}
	}
	stooqCloses, stooqErr := fetchStooqDailyCloses(client, normalized)
	if stooqErr == nil && len(stooqCloses) > 0 {
		if limit > 0 && len(stooqCloses) > limit {
			return stooqCloses[len(stooqCloses)-limit:], nil
		}
		return stooqCloses, nil
	}
	return nil, fmt.Errorf("yahoo: %v; stooq: %v", err, stooqErr)
}

type dailyCloseCandidate struct {
	Source string
	Closes []dailyClose
}

func fetchGlobalIndexDailyCloses(client *http.Client, symbol string, eastmoneySecID string, limit int) ([]dailyClose, error) {
	candidates := []dailyCloseCandidate{}
	errs := []string{}

	addCandidate := func(source string, closes []dailyClose, err error) {
		if err != nil {
			errs = append(errs, source+": "+err.Error())
			return
		}
		if len(closes) == 0 {
			errs = append(errs, source+": empty daily closes")
			return
		}
		candidates = append(candidates, dailyCloseCandidate{Source: source, Closes: closes})
	}

	if strings.EqualFold(symbol, "^NDX") || strings.EqualFold(symbol, "NDX") {
		closes, err := fetchNasdaqHistoricalCloses(client, "NDX", "index", limit)
		if err == nil && len(closes) > 0 {
			if latest, latestErr := fetchNasdaqLatestQuoteClose(client, "NDX", "index"); latestErr == nil {
				closes = appendOrReplaceLatestDailyClose(closes, latest)
			}
		}
		addCandidate("nasdaq", closes, err)
	}
	if strings.TrimSpace(eastmoneySecID) != "" {
		closes, err := fetchEastmoneyDailyClosesBySecID(client, eastmoneySecID, limit)
		addCandidate("eastmoney", closes, err)
	}
	closes, err := fetchYahooDailyCloses(client, symbol, "2y")
	addCandidate("yahoo", closes, err)
	if strings.EqualFold(symbol, "^GSPC") || strings.EqualFold(symbol, "SPX") {
		closes, err := fetchNasdaqHistoricalCloses(client, "SPY", "etf", limit)
		addCandidate("nasdaq-spy", closes, err)
	}
	closes, err = fetchStooqDailyCloses(client, symbol)
	addCandidate("stooq", closes, err)

	if closes, _, ok := selectLatestDailyCloseCandidate(candidates, limit); ok {
		return closes, nil
	}
	return nil, errors.New(strings.Join(errs, "; "))
}

func selectLatestDailyCloseCandidate(candidates []dailyCloseCandidate, limit int) ([]dailyClose, string, bool) {
	var best dailyCloseCandidate
	bestDate := ""
	for _, candidate := range candidates {
		date := latestDailyCloseDate(candidate.Closes)
		if date == "" {
			continue
		}
		if bestDate == "" || date > bestDate {
			best = candidate
			bestDate = date
		}
	}
	if bestDate == "" {
		return nil, "", false
	}
	closes := best.Closes
	if limit > 0 && len(closes) > limit {
		closes = closes[len(closes)-limit:]
	}
	return closes, best.Source, true
}

func latestDailyCloseDate(closes []dailyClose) string {
	if len(closes) == 0 {
		return ""
	}
	return strings.TrimSpace(closes[len(closes)-1].Date)
}

func appendOrReplaceLatestDailyClose(closes []dailyClose, latest dailyClose) []dailyClose {
	if latest.Date == "" || latest.Price <= 0 {
		return closes
	}
	if len(closes) == 0 {
		return []dailyClose{latest}
	}
	result := append([]dailyClose(nil), closes...)
	last := result[len(result)-1]
	if latest.Date < last.Date {
		return result
	}
	if latest.Date == last.Date {
		result[len(result)-1] = latest
		return result
	}
	result = append(result, latest)
	return result
}

func eastmoneyGlobalIndexSecID(symbol string) string {
	switch strings.ToUpper(strings.TrimSpace(symbol)) {
	case "^GSPC", "SPX":
		return "100.SPX"
	case "^NDX", "NDX":
		return "100.NDX"
	default:
		return ""
	}
}

func fetchYahooDailyCloses(client *http.Client, symbol string, rangeParam string) ([]dailyClose, error) {
	sourceSymbol := yahooSymbol(symbol)
	if strings.TrimSpace(rangeParam) == "" {
		rangeParam = "1y"
	}
	endpoint := "https://query2.finance.yahoo.com/v8/finance/chart/" + url.PathEscape(sourceSymbol) + "?range=" + url.QueryEscape(rangeParam) + "&interval=1d"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://finance.yahoo.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("daily close request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload yahooChartResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Chart.Error != nil || len(payload.Chart.Result) == 0 {
		return nil, errors.New("empty daily close response")
	}
	result := payload.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return nil, errors.New("missing close series")
	}
	location := loadLocation(result.Meta.ExchangeTimezone)
	closes := result.Indicators.Quote[0].Close
	validCloses := make([]dailyClose, 0, len(closes))
	for i, closePrice := range closes {
		if closePrice > 0 {
			validCloses = append(validCloses, dailyClose{Price: closePrice, Date: closeDate(result.Timestamp, i, location)})
		}
	}
	if len(validCloses) == 0 {
		return nil, errors.New("no valid close prices")
	}
	return validCloses, nil
}

func fetchNasdaqHistoricalCloses(client *http.Client, symbol string, assetClass string, limit int) ([]dailyClose, error) {
	if limit <= 0 {
		limit = 280
	}
	toDate := time.Now().Format("2006-01-02")
	calendarDays := int(math.Ceil(float64(limit)*365/252)) + 60
	fromDate := time.Now().AddDate(0, 0, -calendarDays).Format("2006-01-02")
	values := url.Values{}
	values.Set("assetclass", assetClass)
	values.Set("fromdate", fromDate)
	values.Set("todate", toDate)
	values.Set("limit", strconv.Itoa(limit+80))
	endpoint := "https://api.nasdaq.com/api/quote/" + url.PathEscape(strings.ToUpper(strings.TrimSpace(symbol))) + "/historical?" + values.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.nasdaq.com")
	req.Header.Set("Referer", "https://www.nasdaq.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("nasdaq historical request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	closes, err := parseNasdaqHistoricalCloses(body)
	if err != nil {
		return nil, err
	}
	if len(closes) > limit {
		return closes[len(closes)-limit:], nil
	}
	return closes, nil
}

func fetchNasdaqLatestQuoteClose(client *http.Client, symbol string, assetClass string) (dailyClose, error) {
	values := url.Values{}
	values.Set("assetclass", assetClass)
	endpoint := "https://api.nasdaq.com/api/quote/" + url.PathEscape(strings.ToUpper(strings.TrimSpace(symbol))) + "/info?" + values.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return dailyClose{}, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.nasdaq.com")
	req.Header.Set("Referer", "https://www.nasdaq.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return dailyClose{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return dailyClose{}, fmt.Errorf("nasdaq quote request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		Data struct {
			PrimaryData struct {
				LastSalePrice      string `json:"lastSalePrice"`
				LastTradeTimestamp string `json:"lastTradeTimestamp"`
			} `json:"primaryData"`
		} `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil {
		return dailyClose{}, err
	}
	price, err := parseMarketNumber(payload.Data.PrimaryData.LastSalePrice)
	if err != nil || price <= 0 {
		return dailyClose{}, errors.New("missing nasdaq quote close price")
	}
	date := normalizeNasdaqQuoteDate(payload.Data.PrimaryData.LastTradeTimestamp)
	if date == "" {
		return dailyClose{}, errors.New("missing nasdaq quote close date")
	}
	return dailyClose{Date: date, Price: price}, nil
}

func parseNasdaqHistoricalCloses(body []byte) ([]dailyClose, error) {
	var payload struct {
		Data struct {
			TradesTable struct {
				Rows []struct {
					Date  string `json:"date"`
					Close string `json:"close"`
				} `json:"rows"`
			} `json:"tradesTable"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	rows := payload.Data.TradesTable.Rows
	if len(rows) == 0 {
		return nil, errors.New("missing nasdaq historical rows")
	}
	closes := make([]dailyClose, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		date := normalizeNasdaqHistoricalDate(rows[i].Date)
		if date == "" {
			continue
		}
		closePrice, err := parseMarketNumber(rows[i].Close)
		if err != nil || closePrice <= 0 {
			continue
		}
		closes = append(closes, dailyClose{Date: date, Price: closePrice})
	}
	if len(closes) == 0 {
		return nil, errors.New("missing nasdaq historical close prices")
	}
	return closes, nil
}

func normalizeNasdaqHistoricalDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"01/02/2006", "1/2/2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format("2006-01-02")
		}
	}
	return ""
}

func normalizeNasdaqQuoteDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"Jan 2, 2006", "Jan 02, 2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format("2006-01-02")
		}
	}
	return ""
}

func parseMarketNumber(value string) (float64, error) {
	cleaned := strings.NewReplacer(",", "", "$", "", " ", "").Replace(strings.TrimSpace(value))
	return strconv.ParseFloat(cleaned, 64)
}

func drawdownFromRecentHigh(closes []dailyClose, window int) (float64, string, error) {
	if len(closes) == 0 {
		return 0, "", errors.New("missing close prices")
	}
	if window <= 0 || window > len(closes) {
		window = len(closes)
	}
	recent := closes[len(closes)-window:]
	latest := recent[len(recent)-1]
	high := 0.0
	for _, close := range recent {
		if close.Price > high {
			high = close.Price
		}
	}
	if high <= 0 || latest.Price <= 0 {
		return 0, "", errors.New("invalid close prices")
	}
	drawdown := (high - latest.Price) / high
	return drawdown, latest.Date, nil
}

func fetchStooqDailyCloses(client *http.Client, symbol string) ([]dailyClose, error) {
	sourceSymbol, err := stooqSymbol(symbol)
	if err != nil {
		return nil, err
	}
	endpoint := "https://stooq.com/q/d/l/?s=" + url.QueryEscape(sourceSymbol) + "&i=d"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "holds-website etf rule updater")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("stooq request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	return parseStooqDailyCSV(body)
}

func stooqSymbol(symbol string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(symbol)) {
	case "^GSPC", "SPX":
		return "^spx", nil
	case "^NDX", "NDX":
		return "^ndx", nil
	default:
		return "", fmt.Errorf("unsupported stooq symbol: %s", symbol)
	}
}

func parseStooqDailyCSV(body []byte) ([]dailyClose, error) {
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(lines) < 2 {
		return nil, errors.New("empty stooq csv")
	}
	closes := make([]dailyClose, 0, len(lines)-1)
	for _, line := range lines[1:] {
		fields := strings.Split(strings.TrimSpace(line), ",")
		if len(fields) < 5 || strings.EqualFold(fields[4], "null") {
			continue
		}
		date := strings.TrimSpace(fields[0])
		if _, err := time.Parse("2006-01-02", date); err != nil {
			continue
		}
		price, err := strconv.ParseFloat(strings.TrimSpace(fields[4]), 64)
		if err != nil || price <= 0 {
			continue
		}
		closes = append(closes, dailyClose{Date: date, Price: price})
	}
	if len(closes) == 0 {
		return nil, errors.New("missing stooq close prices")
	}
	return closes, nil
}

type etfRuleValuation struct {
	Value float64
	Date  string
	Unit  string
	Kind  string
}

func fetchETFRuleValuation(client *http.Client, config etfRuleConfig) (etfRuleValuation, error) {
	switch config.Symbol {
	case "022434":
		value, date, err := fetchA500PEPercentile(client, time.Now())
		return etfRuleValuation{Value: value, Date: date, Unit: "%", Kind: "percentile"}, err
	case "018738":
		value, date, err := fetchSP500PEPercentile(client)
		return etfRuleValuation{Value: value, Date: date, Unit: "%", Kind: "percentile"}, err
	case "021000":
		value, date, err := fetchNasdaq100PEPercentile(client)
		return etfRuleValuation{Value: value, Date: date, Unit: "%", Kind: "percentile"}, err
	default:
		return etfRuleValuation{}, errors.New("valuation source not configured")
	}
}

type fundDBIndexPEPoint struct {
	Percentile float64
	Date       string
}

func fetchFundDBIndexPBPercentile(client *http.Client, indexCode string, category string) (float64, string, error) {
	body, err := fetchFundDBIndexValuationPayload(client, indexCode, category)
	if err != nil {
		return 0, "", err
	}
	point, err := parseFundDBIndexPercentile(body, "pb", "PB")
	if err != nil {
		return 0, "", err
	}
	return point.Percentile, point.Date, nil
}

func fetchFundDBIndexPEPercentile(client *http.Client, indexCode string, category string) (float64, string, error) {
	body, err := fetchFundDBIndexValuationPayload(client, indexCode, category)
	if err != nil {
		return 0, "", err
	}
	point, err := parseFundDBIndexPEPercentile(body)
	if err != nil {
		return 0, "", err
	}
	return point.Percentile, point.Date, nil
}

func fetchFundDBIndexValuationPayload(client *http.Client, indexCode string, category string) ([]byte, error) {
	payload := map[string]any{
		"gu_code":     strings.TrimSpace(indexCode),
		"pe_category": "pe",
		"year":        10,
		"category":    strings.TrimSpace(category),
		"ver":         "new",
	}
	body, err := fundDBPost(client, "/v2/guzhi/newtubiaodata", payload)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parseFundDBIndexPEPercentile(body []byte) (fundDBIndexPEPoint, error) {
	return parseFundDBIndexPercentile(body, "pe", "PE")
}

func parseFundDBIndexPBPercentile(body []byte) (fundDBIndexPEPoint, error) {
	return parseFundDBIndexPercentile(body, "pb", "PB")
}

func parseFundDBIndexPercentile(body []byte, attribute string, label string) (fundDBIndexPEPoint, error) {
	var payload struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			UpdateTime string `json:"update_time"`
			TopData    []struct {
				Attribute       string `json:"attribute"`
				Name            string `json:"name"`
				NewPercentValue struct {
					Value string `json:"value"`
				} `json:"new_percent_value"`
			} `json:"top_data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return fundDBIndexPEPoint{}, err
	}
	if payload.Code != 0 {
		return fundDBIndexPEPoint{}, fmt.Errorf("funddb guzhi response code %d: %s", payload.Code, strings.TrimSpace(payload.Message))
	}
	date := strings.TrimSpace(payload.Data.UpdateTime)
	if _, err := time.Parse(etfRuleRuntimeTimestampDateLayout, date); err != nil {
		return fundDBIndexPEPoint{}, fmt.Errorf("funddb missing update date: %s", date)
	}
	for _, item := range payload.Data.TopData {
		if strings.TrimSpace(item.Attribute) != strings.TrimSpace(attribute) {
			continue
		}
		percentile, err := parseFundDBPercentile(item.NewPercentValue.Value)
		if err != nil {
			return fundDBIndexPEPoint{}, err
		}
		return fundDBIndexPEPoint{Percentile: percentile, Date: date}, nil
	}
	return fundDBIndexPEPoint{}, fmt.Errorf("funddb missing %s percentile", label)
}

func parseFundDBPercentile(value string) (float64, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "--" {
		return 0, errors.New("funddb percentile is empty")
	}
	number, err := firstTextNumber(trimmed)
	if err != nil {
		return 0, err
	}
	if strings.Contains(trimmed, "%") || number > 1 {
		number /= 100
	}
	if number < 0 || number > 1 || math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, fmt.Errorf("funddb percentile out of range: %s", value)
	}
	return number, nil
}

func fundDBPost(client *http.Client, path string, payload map[string]any) ([]byte, error) {
	signed := fundDBSignedPayload(payload, time.Now())
	body, err := json.Marshal(signed)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fundDBAPIHost+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://funddb.cn")
	req.Header.Set("Referer", fundDBIndexPageURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	requestClient := client
	if requestClient == nil {
		requestClient = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := requestClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("funddb request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return body, nil
}

func fundDBSignedPayload(payload map[string]any, now time.Time) map[string]any {
	signed := map[string]any{}
	for key, value := range payload {
		signed[key] = value
	}
	signed["type"] = "pc"
	signed["version"] = fundDBAPIVersion
	if _, ok := signed["authtoken"]; !ok {
		signed["authtoken"] = ""
	}
	signed["act_time"] = now.UnixMilli()

	keys := make([]string, 0, len(signed))
	for key := range signed {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	signatureText := strings.Builder{}
	for _, key := range keys {
		valueText, include := fundDBSignatureValue(signed[key])
		if include {
			signatureText.WriteString(valueText)
		}
	}
	signatureText.WriteString(fundDBAPIReqKey)
	sum := md5.Sum([]byte(signatureText.String()))
	fundDBApplySignatureFields(signed, fmt.Sprintf("%x", sum))
	return signed
}

func fundDBSignatureValue(value any) (string, bool) {
	switch typed := value.(type) {
	case nil:
		return "", false
	case string:
		return typed, strings.TrimSpace(typed) != ""
	case bool:
		if !typed {
			return "", false
		}
		return "true", true
	case int:
		return strconv.Itoa(typed), true
	case int64:
		return strconv.FormatInt(typed, 10), true
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64), true
	case map[string]any, []any:
		return "", false
	default:
		return fmt.Sprint(typed), true
	}
}

func fundDBApplySignatureFields(payload map[string]any, signature string) {
	if len(signature) < 32 {
		return
	}
	part := func(start int, end int) string { return signature[start:end] }
	payload["tirgkjfs"] = part(0, 2)
	payload["abiokytke"] = part(21, 23)
	payload["u54rg5d"] = part(2, 4)
	payload["kf54ge7"] = part(31, 32)
	payload["tiklsktr4"] = part(1, 2)
	payload["lksytkjh"] = part(17, 21)
	payload["sbnoywr"] = part(23, 25)
	payload["bgd7h8tyu54"] = part(6, 8)
	payload["y654b5fs3tr"] = part(11, 12)
	payload["bioduytlw"] = part(5, 6)
	payload["bd4uy742"] = part(26, 27)
	payload["h67456y"] = part(16, 19)
	payload["bvytikwqjk"] = part(6, 8)
	payload["ngd4uy551"] = part(17, 19)
	payload["bgiuytkw"] = part(9, 11)
	payload["nd354uy4752"] = part(30, 31)
	payload["ghtoiutkmlg"] = part(11, 14)
	payload["bd24y6421f"] = part(24, 26)
	payload["tbvdiuytk"] = part(16, 17)
	payload["ibvytiqjek"] = part(14, 16)
	payload["jnhf8u5231"] = part(9, 11)
	payload["fjlkatj"] = part(2, 5)
	payload["hy5641d321t"] = part(25, 27)
	payload["iogojti"] = part(25, 26)
	payload["ngd4yut78"] = part(12, 14)
	payload["nkjhrew"] = part(26, 27)
	payload["yt447e13f"] = part(8, 9)
	payload["n3bf4uj7y7"] = part(18, 19)
	payload["nbf4uj7y432"] = part(21, 23)
	payload["yi854tew"] = part(29, 31)
	payload["h13ey474"] = part(29, 32)
	payload["quikgdky"] = part(27, 29)
}

type leguleguIndexPERow struct {
	Date          string
	TtmPE         float64
	TtmPEQuantile *float64
}

func fetchA500PEPercentile(client *http.Client, now time.Time) (float64, string, error) {
	percentile, date, fundDBErr := fetchFundDBIndexPEPercentile(client, "000510.SH", "")
	if fundDBErr == nil && !valuationDateStale(date, now, primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	if fundDBErr == nil {
		fundDBErr = fmt.Errorf("funddb PE stale as of %s", date)
	}
	rows, err := fetchLeguleguA500PERows(client, now)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("funddb: %v; legulegu: %v", fundDBErr, err)
	}
	fallbackPercentile, fallbackDate, err := a500PEPercentileFromRows(rows)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("funddb: %v; legulegu: %v", fundDBErr, err)
	}
	if date != "" && !quoteDateAfter(fallbackDate, date) {
		return percentile, date, nil
	}
	return fallbackPercentile, fallbackDate, nil
}

func fetchSP500PEPercentile(client *http.Client) (float64, string, error) {
	now := time.Now()
	percentile, date, fundDBErr := fetchFundDBIndexPEPercentile(client, "SPX.GI", "5")
	if fundDBErr == nil && !valuationDateStale(date, now, primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	if fundDBErr == nil {
		fundDBErr = fmt.Errorf("funddb PE stale as of %s", date)
	}
	fallbackPercentile, fallbackDate, err := fetchSP500CAPEPercentile(client)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("funddb: %v; historyofmarket CAPE: %v", fundDBErr, err)
	}
	if date != "" && !quoteDateAfter(fallbackDate, date) {
		return percentile, date, nil
	}
	return fallbackPercentile, fallbackDate, nil
}

func fetchNasdaq100PEPercentile(client *http.Client) (float64, string, error) {
	now := time.Now()
	percentile, date, fundDBErr := fetchFundDBIndexPEPercentile(client, "NDX.GI", "5")
	if fundDBErr == nil && !valuationDateStale(date, now, primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	if fundDBErr == nil {
		fundDBErr = fmt.Errorf("funddb PE stale as of %s", date)
	}
	fallbackPercentile, fallbackDate, err := fetchNasdaq100ForwardPEPercentile(client)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("funddb: %v; historyofmarket forward PE: %v", fundDBErr, err)
	}
	if date != "" && !quoteDateAfter(fallbackDate, date) {
		return percentile, date, nil
	}
	return fallbackPercentile, fallbackDate, nil
}

func fetchLeguleguA500PEPercentile(client *http.Client, now time.Time) (float64, string, error) {
	rows, err := fetchLeguleguA500PERows(client, now)
	if err != nil {
		return 0, "", err
	}
	return a500PEPercentileFromRows(rows)
}

func fetchLeguleguA500PERows(client *http.Client, now time.Time) ([]leguleguIndexPERow, error) {
	pageURL := "https://legulegu.com/stockdata/index-ttm-lyr-pe?indexCode=000510.CSI"
	pageReq, err := http.NewRequest(http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, err
	}
	pageReq.Header.Set("Accept", "text/html,*/*")
	pageReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	pageReq.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	pageResp, err := client.Do(pageReq)
	if err != nil {
		return nil, err
	}
	cookieHeader := responseCookieHeader(pageResp)
	io.Copy(io.Discard, io.LimitReader(pageResp.Body, 1<<20))
	pageResp.Body.Close()
	if pageResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("legulegu page request failed: %s", pageResp.Status)
	}
	if strings.TrimSpace(cookieHeader) == "" {
		return nil, errors.New("missing legulegu session cookies")
	}

	var lastErr error
	for _, tokenDate := range []time.Time{now, now.AddDate(0, 0, -1), now.AddDate(0, 0, -2)} {
		rows, err := fetchLeguleguA500PERowsWithToken(client, cookieHeader, tokenDate)
		if err == nil && len(rows) > 0 {
			return rows, nil
		}
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("empty legulegu A500 PE response")
}

func fetchLeguleguA500PERowsWithToken(client *http.Client, cookieHeader string, tokenDate time.Time) ([]leguleguIndexPERow, error) {
	values := url.Values{}
	values.Set("indexCode", "000510.CSI")
	values.Set("token", leguleguToken(tokenDate))
	endpoint := "https://legulegu.com/api/stockdata/index-basic-pe?" + values.Encode()
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", cookieHeader)
	req.Header.Set("Referer", "https://legulegu.com/stockdata/index-ttm-lyr-pe?indexCode=000510.CSI")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("legulegu A500 PE request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return nil, errors.New("empty legulegu A500 PE response")
	}
	return parseLeguleguIndexPERows(body)
}

func responseCookieHeader(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	parts := []string{}
	for _, cookie := range resp.Cookies() {
		if strings.TrimSpace(cookie.Name) != "" {
			parts = append(parts, cookie.Name+"="+cookie.Value)
		}
	}
	return strings.Join(parts, "; ")
}

func leguleguToken(date time.Time) string {
	sum := md5.Sum([]byte(date.Format("2006-01-02")))
	return fmt.Sprintf("%x", sum)
}

func parseLeguleguIndexPERows(body []byte) ([]leguleguIndexPERow, error) {
	var payload struct {
		Data []struct {
			Date          string   `json:"date"`
			TtmPE         float64  `json:"ttmPe"`
			TtmPEQuantile *float64 `json:"ttmPeQuantile"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	rows := make([]leguleguIndexPERow, 0, len(payload.Data))
	for _, row := range payload.Data {
		if _, err := time.Parse("2006-01-02", row.Date); err != nil {
			continue
		}
		if row.TtmPE <= 0 && !validPercentilePointer(row.TtmPEQuantile) {
			continue
		}
		rows = append(rows, leguleguIndexPERow{Date: row.Date, TtmPE: row.TtmPE, TtmPEQuantile: row.TtmPEQuantile})
	}
	if len(rows) == 0 {
		return nil, errors.New("missing legulegu A500 PE rows")
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Date < rows[j].Date })
	return rows, nil
}

func a500PEPercentileFromRows(rows []leguleguIndexPERow) (float64, string, error) {
	if len(rows) == 0 {
		return 0, "", errors.New("missing A500 PE rows")
	}
	latest := rows[len(rows)-1]
	latestDate, err := time.Parse("2006-01-02", latest.Date)
	if err != nil {
		return 0, "", err
	}
	if validPercentilePointer(latest.TtmPEQuantile) {
		return *latest.TtmPEQuantile, latest.Date, nil
	}
	cutoff := latestDate.AddDate(-5, 0, 0)
	values := make([]float64, 0, len(rows))
	for _, row := range rows {
		rowDate, err := time.Parse("2006-01-02", row.Date)
		if err != nil || rowDate.Before(cutoff) || row.TtmPE <= 0 {
			continue
		}
		values = append(values, row.TtmPE)
	}
	if len(values) == 0 || latest.TtmPE <= 0 {
		return 0, "", errors.New("missing five-year A500 PE values")
	}
	return percentileRank(latest.TtmPE, values), latest.Date, nil
}

func validPercentilePointer(value *float64) bool {
	return value != nil && !math.IsNaN(*value) && !math.IsInf(*value, 0) && *value >= 0 && *value <= 1
}

func fetchSP500CAPEPercentile(client *http.Client) (float64, string, error) {
	percentile, date, err := fetchHistoryOfMarketSP500CAPEPercentile(client)
	if err == nil && !valuationDateStale(date, time.Now(), primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	historyErr := err
	if historyErr == nil {
		historyErr = fmt.Errorf("historyofmarket CAPE stale as of %s", date)
	}
	observations, err := fetchMultplMonthlyValues(client, multplShillerCAPEURL)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("historyofmarket: %v; multpl: %v", historyErr, err)
	}
	fallbackPercentile, fallbackDate, err := capePercentileFromMonthlyValues(observations, 10)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("historyofmarket: %v; multpl: %v", historyErr, err)
	}
	if date != "" && !quoteDateAfter(fallbackDate, date) {
		return percentile, date, nil
	}
	return fallbackPercentile, fallbackDate, nil
}

func fetchNasdaq100ForwardPEPercentile(client *http.Client) (float64, string, error) {
	percentile, date, err := fetchHistoryOfMarketNasdaq100ForwardPEPercentile(client)
	if err == nil && !valuationDateStale(date, time.Now(), primaryValuationMaxLagDays) {
		return percentile, date, nil
	}
	historyErr := err
	if historyErr == nil {
		historyErr = fmt.Errorf("historyofmarket Nasdaq 100 forward PE stale as of %s", date)
	}
	snapshot, err := fetchWorldPERatioNasdaq100(client, worldPERatioNasdaq100URL)
	if err != nil {
		if date != "" {
			return percentile, date, nil
		}
		return 0, "", fmt.Errorf("historyofmarket: %v; worldperatio: %v", historyErr, err)
	}
	return snapshot.Percentile, snapshot.Date, nil
}

type historyOfMarketPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type historyOfMarketCurrentValuation struct {
	Forward float64 `json:"forward"`
}

func fetchHistoryOfMarketSP500CAPEPercentile(client *http.Client) (float64, string, error) {
	var payload struct {
		CAPE []historyOfMarketPoint `json:"cape"`
	}
	if err := fetchHistoryOfMarketJSON(client, historyOfMarketSP500PEURL, &payload); err != nil {
		return 0, "", err
	}
	return percentileFromHistoryOfMarketPoints(payload.CAPE, 10, "CAPE")
}

func fetchHistoryOfMarketNasdaq100ForwardPEPercentile(client *http.Client) (float64, string, error) {
	var payload struct {
		Updated string                          `json:"updated"`
		Current historyOfMarketCurrentValuation `json:"current"`
		Forward []historyOfMarketPoint          `json:"forward"`
	}
	if err := fetchHistoryOfMarketJSON(client, historyOfMarketNDXForwardPEURL, &payload); err != nil {
		return 0, "", err
	}
	points := historyOfMarketPointsWithCurrentForward(payload.Forward, payload.Updated, payload.Current)
	return percentileFromHistoryOfMarketPoints(points, 10, "Nasdaq 100 forward PE")
}

func historyOfMarketPointsWithCurrentForward(points []historyOfMarketPoint, updated string, current historyOfMarketCurrentValuation) []historyOfMarketPoint {
	combined := append([]historyOfMarketPoint(nil), points...)
	if current.Forward > 0 && strings.TrimSpace(updated) != "" {
		combined = append(combined, historyOfMarketPoint{
			Date:  strings.TrimSpace(updated),
			Value: current.Forward,
		})
	}
	return combined
}

func fetchHistoryOfMarketJSON(client *http.Client, endpoint string, target any) error {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "holds-website etf rule updater")
	requestClient := client
	if requestClient == nil {
		requestClient = &http.Client{Timeout: 20 * time.Second}
	}
	resp, err := requestClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("historyofmarket request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(target)
}

func percentileFromHistoryOfMarketPoints(points []historyOfMarketPoint, years int, label string) (float64, string, error) {
	observations := make([]dailyClose, 0, len(points))
	for _, point := range points {
		if _, err := time.Parse(etfRuleRuntimeTimestampDateLayout, point.Date); err != nil {
			continue
		}
		if point.Value <= 0 {
			continue
		}
		observations = append(observations, dailyClose{Date: point.Date, Price: point.Value})
	}
	return percentileFromDatedValues(observations, years, label)
}

func capePercentileFromMonthlyValues(observations []dailyClose, years int) (float64, string, error) {
	return percentileFromDatedValues(observations, years, "CAPE")
}

func valuationDateStale(date string, now time.Time, maxLagDays int) bool {
	parsed, err := time.Parse(etfRuleRuntimeTimestampDateLayout, strings.TrimSpace(date))
	if err != nil {
		return true
	}
	reference := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, parsed.Location())
	if reference.Before(parsed) {
		return false
	}
	return reference.Sub(parsed) > time.Duration(maxLagDays)*24*time.Hour
}

func percentileFromDatedValues(observations []dailyClose, years int, label string) (float64, string, error) {
	if len(observations) < 5 {
		return 0, "", fmt.Errorf("not enough %s observations", label)
	}
	ordered := append([]dailyClose(nil), observations...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Date < ordered[j].Date })
	latest := dailyClose{}
	for i := len(ordered) - 1; i >= 0; i-- {
		if ordered[i].Price > 0 && strings.TrimSpace(ordered[i].Date) != "" {
			latest = ordered[i]
			break
		}
	}
	if latest.Price <= 0 {
		return 0, "", fmt.Errorf("missing %s values", label)
	}
	latestDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, latest.Date)
	if err != nil {
		return 0, "", err
	}
	cutoff := latestDate.AddDate(-years, 0, 0)
	values := make([]float64, 0, len(ordered))
	for _, observation := range ordered {
		observationDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, observation.Date)
		if err != nil || observationDate.Before(cutoff) || observation.Price <= 0 {
			continue
		}
		values = append(values, observation.Price)
	}
	if len(values) == 0 {
		return 0, "", fmt.Errorf("missing ten-year %s values", label)
	}
	return percentileRank(latest.Price, values), latest.Date, nil
}

func fetchMultplMonthlyValues(client *http.Client, endpoint string) ([]dailyClose, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "holds-website etf rule updater")
	req.Header.Set("Referer", "https://www.multpl.com/")

	requestClient := client
	if requestClient == nil || (requestClient.Timeout > 0 && requestClient.Timeout < 30*time.Second) {
		requestClient = &http.Client{Timeout: 30 * time.Second}
	}

	resp, err := requestClient.Do(req)
	if err != nil {
		time.Sleep(500 * time.Millisecond)
		req, retryReqErr := http.NewRequest(http.MethodGet, endpoint, nil)
		if retryReqErr != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "holds-website etf rule updater")
		req.Header.Set("Referer", "https://www.multpl.com/")
		resp, err = requestClient.Do(req)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("valuation request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	return parseMultplTable(body)
}

func parseMultplTable(body []byte) ([]dailyClose, error) {
	text := string(body)
	rowPattern := regexp.MustCompile(`(?is)<tr[^>]*>.*?</tr>`)
	cellPattern := regexp.MustCompile("(?is)<td[^>]*>(.*?)</td>")
	rows := rowPattern.FindAllString(text, -1)
	values := make([]dailyClose, 0, len(rows))
	for _, row := range rows {
		cells := cellPattern.FindAllStringSubmatch(row, -1)
		if len(cells) < 2 {
			continue
		}
		dateText := htmlPlainText(cells[0][1])
		valueText := htmlPlainText(cells[1][1])
		value, err := firstTextNumber(valueText)
		if err != nil || value <= 0 {
			continue
		}
		date := normalizeMultplDate(dateText)
		if date == "" {
			continue
		}
		values = append(values, dailyClose{Date: date, Price: value})
	}
	if len(values) == 0 {
		return nil, errors.New("missing valuation rows")
	}
	return values, nil
}

type cashDividendEvent struct {
	Date   string
	Amount float64
}

func fetchDividendLowVolYield(client *http.Client) (float64, string, error) {
	snapshot, err := fetchDividendLowVolValuationSnapshot(client)
	if err != nil {
		return 0, "", err
	}
	return snapshot.Yield, snapshot.Date, nil
}

type dividendLowVolValuationSnapshot struct {
	Yield              float64
	Percentile         float64
	Date               string
	BondYield          float64
	BondDate           string
	Spread             float64
	SpreadPercentile   float64
	SpreadDate         string
	SpreadAvailable    bool
	SpreadError        string
	SpreadObservations int
}

type datedRate struct {
	Date  string
	Value float64
}

// fetchDividendLowVolValuationSnapshot is a legacy fund-payout diagnostic. It
// must not be used as the target index's TTM dividend yield or trading signal.
func fetchDividendLowVolValuationSnapshot(client *http.Client) (dividendLowVolValuationSnapshot, error) {
	closes, err := fetchTencentRawDailyCloses(client, "515450.SH", 10*252+80)
	if err != nil {
		return dividendLowVolValuationSnapshot{}, err
	}
	if len(closes) == 0 || closes[len(closes)-1].Price <= 0 {
		return dividendLowVolValuationSnapshot{}, errors.New("missing 515450 close price")
	}
	events, err := fetchEastmoneyFundDividends(client, "515450")
	if err != nil {
		return dividendLowVolValuationSnapshot{}, err
	}
	latest := closes[len(closes)-1]
	trailingAmount, err := trailingFundDividendAmount(events, latest.Date)
	if err != nil {
		return dividendLowVolValuationSnapshot{}, err
	}
	currentYield := trailingAmount / latest.Price
	monthly := monthEndCloses(closes)
	yieldHistory := make([]datedRate, 0, len(monthly))
	for _, close := range monthly {
		amount, err := trailingFundDividendAmount(events, close.Date)
		if err != nil || amount <= 0 || close.Price <= 0 {
			continue
		}
		yieldHistory = append(yieldHistory, datedRate{Date: close.Date, Value: amount / close.Price})
	}
	if len(yieldHistory) < 12 {
		return dividendLowVolValuationSnapshot{}, fmt.Errorf("insufficient dividend-yield history: %d monthly observations", len(yieldHistory))
	}
	historyValues := make([]float64, 0, len(yieldHistory))
	for _, point := range yieldHistory {
		historyValues = append(historyValues, point.Value)
	}
	snapshot := dividendLowVolValuationSnapshot{
		Yield:      currentYield,
		Percentile: percentileRank(currentYield, historyValues),
		Date:       latest.Date,
	}
	oldestDate := yieldHistory[0].Date
	bondHistory, bondErr := fetchEastmoneyChina10YBondYieldHistory(client, oldestDate)
	if bondErr != nil {
		snapshot.SpreadError = bondErr.Error()
		return snapshot, nil
	}
	latestDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, latest.Date)
	if err != nil {
		snapshot.SpreadError = err.Error()
		return snapshot, nil
	}
	officialBond, officialErr := fetchChinaBondOfficial10YYield(client, latestDate)
	if officialErr != nil {
		snapshot.SpreadError = officialErr.Error()
		return snapshot, nil
	}
	spread, spreadPercentile, observationCount, spreadErr := calculateDividendSpread(
		datedRate{Date: latest.Date, Value: currentYield},
		yieldHistory,
		bondHistory,
		officialBond,
	)
	if spreadErr != nil {
		snapshot.SpreadError = spreadErr.Error()
		return snapshot, nil
	}
	snapshot.BondYield = officialBond.Value
	snapshot.BondDate = officialBond.Date
	snapshot.Spread = spread
	snapshot.SpreadPercentile = spreadPercentile
	snapshot.SpreadDate = latest.Date
	snapshot.SpreadAvailable = true
	snapshot.SpreadObservations = observationCount
	return snapshot, nil
}

func monthEndCloses(closes []dailyClose) []dailyClose {
	byMonth := map[string]dailyClose{}
	months := []string{}
	for _, close := range closes {
		if len(close.Date) < 7 || close.Price <= 0 {
			continue
		}
		month := close.Date[:7]
		if _, ok := byMonth[month]; !ok {
			months = append(months, month)
		}
		if previous, ok := byMonth[month]; !ok || close.Date > previous.Date {
			byMonth[month] = close
		}
	}
	sort.Strings(months)
	result := make([]dailyClose, 0, len(months))
	for _, month := range months {
		result = append(result, byMonth[month])
	}
	return result
}

func fetchEastmoneyChina10YBondYieldHistory(client *http.Client, startDate string) ([]datedRate, error) {
	if _, err := time.Parse(etfRuleRuntimeTimestampDateLayout, startDate); err != nil {
		return nil, fmt.Errorf("invalid treasury-yield start date: %s", startDate)
	}
	const pageSize = 500
	pointsByDate := map[string]datedRate{}
	pageCount := 1
	for page := 1; page <= pageCount && page <= 20; page++ {
		values := url.Values{}
		values.Set("type", "RPTA_WEB_TREASURYYIELD")
		values.Set("sty", "ALL")
		values.Set("st", "SOLAR_DATE")
		values.Set("sr", "-1")
		values.Set("p", strconv.Itoa(page))
		values.Set("ps", strconv.Itoa(pageSize))
		req, err := http.NewRequest(http.MethodGet, eastmoneyTreasuryYieldURL+"?"+values.Encode(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json,text/plain,*/*")
		req.Header.Set("Referer", "https://data.eastmoney.com/cjsj/zmgzsyl.html")
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("eastmoney treasury-yield request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
		}
		pagePoints, pages, err := parseEastmoneyChina10YBondYields(body)
		if err != nil {
			return nil, err
		}
		if pages > 0 {
			pageCount = pages
		}
		oldest := ""
		for _, point := range pagePoints {
			pointsByDate[point.Date] = point
			if oldest == "" || point.Date < oldest {
				oldest = point.Date
			}
		}
		if oldest != "" && oldest <= startDate {
			break
		}
	}
	points := make([]datedRate, 0, len(pointsByDate))
	for _, point := range pointsByDate {
		if point.Date >= startDate {
			points = append(points, point)
		}
	}
	sort.Slice(points, func(i, j int) bool { return points[i].Date < points[j].Date })
	if len(points) == 0 {
		return nil, errors.New("missing eastmoney China 10Y bond-yield history")
	}
	return points, nil
}

func parseEastmoneyChina10YBondYields(body []byte) ([]datedRate, int, error) {
	var payload struct {
		Result struct {
			Pages int `json:"pages"`
			Data  []struct {
				Date  string   `json:"SOLAR_DATE"`
				Yield *float64 `json:"EMM00166466"`
			} `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, 0, err
	}
	points := make([]datedRate, 0, len(payload.Result.Data))
	for _, row := range payload.Result.Data {
		date := normalizeTreasuryYieldDate(row.Date)
		if date == "" || row.Yield == nil || *row.Yield <= 0 {
			continue
		}
		points = append(points, datedRate{Date: date, Value: *row.Yield / 100})
	}
	if len(points) == 0 {
		return nil, payload.Result.Pages, errors.New("missing eastmoney China 10Y bond-yield rows")
	}
	return points, payload.Result.Pages, nil
}

func fetchChinaBondOfficial10YYield(client *http.Client, endDate time.Time) (datedRate, error) {
	values := url.Values{}
	values.Set("startDate", endDate.AddDate(0, 0, -14).Format(etfRuleRuntimeTimestampDateLayout))
	values.Set("endDate", endDate.Format(etfRuleRuntimeTimestampDateLayout))
	values.Set("gjqx", "10")
	values.Set("qxId", "hzsylqx")
	values.Set("locale", "en_US")
	req, err := http.NewRequest(http.MethodGet, chinaBondHistoryURL+"?"+values.Encode(), nil)
	if err != nil {
		return datedRate{}, err
	}
	req.Header.Set("Accept", "text/html,*/*")
	req.Header.Set("Referer", "https://yield.chinabond.com.cn/cbweb-pbc-web/pbc/more?locale=en_US")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return datedRate{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return datedRate{}, fmt.Errorf("ChinaBond history request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return datedRate{}, err
	}
	return parseChinaBondOfficial10YYield(body)
}

func parseChinaBondOfficial10YYield(body []byte) (datedRate, error) {
	rowPattern := regexp.MustCompile(`(?is)<tr[^>]*>.*?</tr>`)
	cellPattern := regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)
	latest := datedRate{}
	for _, row := range rowPattern.FindAllString(string(body), -1) {
		cells := cellPattern.FindAllStringSubmatch(row, -1)
		if len(cells) < 9 {
			continue
		}
		date := normalizeTreasuryYieldDate(htmlPlainText(cells[1][1]))
		yield, err := firstTextNumber(htmlPlainText(cells[8][1]))
		if date == "" || err != nil || yield <= 0 {
			continue
		}
		if latest.Date == "" || date > latest.Date {
			latest = datedRate{Date: date, Value: yield / 100}
		}
	}
	if latest.Date == "" {
		return datedRate{}, errors.New("missing ChinaBond official 10Y yield")
	}
	return latest, nil
}

func normalizeTreasuryYieldDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format(etfRuleRuntimeTimestampDateLayout)
		}
	}
	return ""
}

func calculateDividendSpread(currentYield datedRate, yieldHistory []datedRate, bondHistory []datedRate, officialBond datedRate) (float64, float64, int, error) {
	if currentYield.Date == "" || currentYield.Value <= 0 || officialBond.Date == "" || officialBond.Value <= 0 {
		return 0, 0, 0, errors.New("invalid current dividend or bond yield")
	}
	officialHistoryPoint, ok := datedRateOnOrBefore(bondHistory, officialBond.Date, 0)
	if !ok {
		return 0, 0, 0, fmt.Errorf("missing Eastmoney bond yield for ChinaBond date %s", officialBond.Date)
	}
	if math.Abs(officialHistoryPoint.Value-officialBond.Value)/officialBond.Value > 0.01 {
		return 0, 0, 0, fmt.Errorf("China 10Y bond-yield sources differ: ChinaBond %.4f%%, Eastmoney %.4f%%", officialBond.Value*100, officialHistoryPoint.Value*100)
	}
	currentDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, currentYield.Date)
	if err != nil {
		return 0, 0, 0, err
	}
	officialDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, officialBond.Date)
	if err != nil {
		return 0, 0, 0, err
	}
	if currentDate.Sub(officialDate).Hours()/24 > 10 || officialDate.After(currentDate) {
		return 0, 0, 0, fmt.Errorf("ChinaBond 10Y yield is not aligned with dividend yield: %s vs %s", officialBond.Date, currentYield.Date)
	}
	spreads := make([]float64, 0, len(yieldHistory))
	for _, dividendYield := range yieldHistory {
		bondYield, ok := datedRateOnOrBefore(bondHistory, dividendYield.Date, 10)
		if !ok {
			continue
		}
		spreads = append(spreads, dividendYield.Value-bondYield.Value)
	}
	if len(spreads) < 12 {
		return 0, 0, len(spreads), fmt.Errorf("insufficient dividend-spread history: %d monthly observations", len(spreads))
	}
	currentSpread := currentYield.Value - officialBond.Value
	return currentSpread, percentileRank(currentSpread, spreads), len(spreads), nil
}

func datedRateOnOrBefore(points []datedRate, targetDate string, maxLagDays int) (datedRate, bool) {
	target, err := time.Parse(etfRuleRuntimeTimestampDateLayout, targetDate)
	if err != nil {
		return datedRate{}, false
	}
	for i := len(points) - 1; i >= 0; i-- {
		point := points[i]
		date, err := time.Parse(etfRuleRuntimeTimestampDateLayout, point.Date)
		if err != nil || date.After(target) {
			continue
		}
		lagDays := int(target.Sub(date).Hours() / 24)
		if lagDays > maxLagDays {
			return datedRate{}, false
		}
		return point, true
	}
	return datedRate{}, false
}

func fetchStockAnalysisDividends(client *http.Client, symbol string) ([]cashDividendEvent, error) {
	endpoint := "https://stockanalysis.com/etf/" + strings.ToLower(strings.TrimSpace(symbol)) + "/dividend/"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://stockanalysis.com/etf/"+strings.ToLower(strings.TrimSpace(symbol))+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/126 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stockanalysis dividend request failed: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
	}
	return parseStockAnalysisDividends(body)
}

const stateStreetSPYDistributionsURL = "https://www.ssga.com/library-content/products/fund-data/etfs/us/spdr-etf-historical-distributions.xlsx"

func fetchStateStreetSPYDividends(client *http.Client) ([]cashDividendEvent, error) {
	req, err := http.NewRequest(http.MethodGet, stateStreetSPYDistributionsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("State Street dividend request failed: %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	return parseStateStreetSPYDividends(body)
}

func parseStateStreetSPYDividends(workbook []byte) ([]cashDividendEvent, error) {
	archive, err := zip.NewReader(bytes.NewReader(workbook), int64(len(workbook)))
	if err != nil {
		return nil, fmt.Errorf("open State Street dividend workbook: %w", err)
	}
	var sharedStringsFile *zip.File
	var dividendSheetFile *zip.File
	for _, file := range archive.File {
		switch file.Name {
		case "xl/sharedStrings.xml":
			sharedStringsFile = file
		case "xl/worksheets/sheet1.xml":
			dividendSheetFile = file
		}
	}
	if sharedStringsFile == nil || dividendSheetFile == nil {
		return nil, errors.New("State Street dividend workbook is missing required worksheets")
	}

	sharedReader, err := sharedStringsFile.Open()
	if err != nil {
		return nil, err
	}
	sharedStrings, err := parseXLSXSharedStrings(sharedReader)
	sharedReader.Close()
	if err != nil {
		return nil, fmt.Errorf("parse State Street shared strings: %w", err)
	}

	sheetReader, err := dividendSheetFile.Open()
	if err != nil {
		return nil, err
	}
	events, err := parseStateStreetDividendSheet(sheetReader, sharedStrings, "SPY")
	sheetReader.Close()
	if err != nil {
		return nil, fmt.Errorf("parse State Street dividend sheet: %w", err)
	}
	if len(events) == 0 {
		return nil, errors.New("State Street dividend workbook has no SPY rows")
	}
	sort.Slice(events, func(i, j int) bool { return events[i].Date < events[j].Date })
	return events, nil
}

func parseXLSXSharedStrings(reader io.Reader) ([]string, error) {
	decoder := xml.NewDecoder(reader)
	values := []string{}
	inSharedString := false
	inText := false
	var value strings.Builder
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch token := token.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "si":
				inSharedString = true
				value.Reset()
			case "t":
				inText = inSharedString
			}
		case xml.EndElement:
			switch token.Name.Local {
			case "t":
				inText = false
			case "si":
				values = append(values, value.String())
				inSharedString = false
			}
		case xml.CharData:
			if inText {
				value.Write([]byte(token))
			}
		}
	}
	return values, nil
}

func parseStateStreetDividendSheet(reader io.Reader, sharedStrings []string, ticker string) ([]cashDividendEvent, error) {
	decoder := xml.NewDecoder(reader)
	events := []cashDividendEvent{}
	row := map[string]string{}
	cellReference := ""
	cellType := ""
	cellValue := ""
	inCellValue := false
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch token := token.(type) {
		case xml.StartElement:
			switch token.Name.Local {
			case "row":
				row = map[string]string{}
			case "c":
				cellReference = ""
				cellType = ""
				cellValue = ""
				for _, attribute := range token.Attr {
					switch attribute.Name.Local {
					case "r":
						cellReference = attribute.Value
					case "t":
						cellType = attribute.Value
					}
				}
			case "v":
				inCellValue = cellReference != ""
			}
		case xml.EndElement:
			switch token.Name.Local {
			case "v":
				inCellValue = false
			case "c":
				column := xlsxColumn(cellReference)
				if column == "" {
					continue
				}
				resolved := strings.TrimSpace(cellValue)
				if cellType == "s" {
					index, err := strconv.Atoi(resolved)
					if err != nil || index < 0 || index >= len(sharedStrings) {
						continue
					}
					resolved = strings.TrimSpace(sharedStrings[index])
				}
				row[column] = resolved
			case "row":
				if !strings.EqualFold(strings.TrimSpace(row["B"]), ticker) {
					continue
				}
				date := normalizeStateStreetDividendDate(row["D"])
				amount, err := parseMarketNumber(row["G"])
				if date == "" || err != nil || amount <= 0 {
					continue
				}
				events = append(events, cashDividendEvent{Date: date, Amount: amount})
			}
		case xml.CharData:
			if inCellValue {
				cellValue += string(token)
			}
		}
	}
	return events, nil
}

func xlsxColumn(reference string) string {
	end := 0
	for end < len(reference) {
		character := reference[end]
		if character < 'A' || character > 'Z' {
			break
		}
		end++
	}
	return reference[:end]
}

func normalizeStateStreetDividendDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"01/02/2006", "1/2/2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format(etfRuleRuntimeTimestampDateLayout)
		}
	}
	serial, err := strconv.ParseFloat(value, 64)
	if err != nil || serial <= 0 {
		return ""
	}
	date := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC).Add(time.Duration(math.Round(serial*24)) * time.Hour)
	return date.Format(etfRuleRuntimeTimestampDateLayout)
}

func parseStockAnalysisDividends(body []byte) ([]cashDividendEvent, error) {
	page := string(body)
	historyIndex := strings.Index(page, "Dividend History")
	if historyIndex < 0 {
		return nil, errors.New("missing stockanalysis dividend history")
	}
	tableBodyStart := strings.Index(page[historyIndex:], "<tbody")
	if tableBodyStart < 0 {
		return nil, errors.New("missing stockanalysis dividend table")
	}
	tableBodyStart += historyIndex
	tableBodyEnd := strings.Index(page[tableBodyStart:], "</tbody>")
	if tableBodyEnd < 0 {
		return nil, errors.New("incomplete stockanalysis dividend table")
	}
	tableBody := page[tableBodyStart : tableBodyStart+tableBodyEnd]
	rowPattern := regexp.MustCompile(`(?is)<tr[^>]*>.*?</tr>`)
	cellPattern := regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)
	events := []cashDividendEvent{}
	for _, row := range rowPattern.FindAllString(tableBody, -1) {
		cells := cellPattern.FindAllStringSubmatch(row, -1)
		if len(cells) < 2 {
			continue
		}
		date := normalizeStockAnalysisDate(htmlPlainText(cells[0][1]))
		amount, err := parseMarketNumber(htmlPlainText(cells[1][1]))
		if date == "" || err != nil || amount <= 0 {
			continue
		}
		events = append(events, cashDividendEvent{Date: date, Amount: amount})
	}
	if len(events) == 0 {
		return nil, errors.New("missing stockanalysis dividend rows")
	}
	return events, nil
}

func normalizeStockAnalysisDate(value string) string {
	for _, layout := range []string{"Jan 2, 2006", "Jan 02, 2006", "2006-01-02"} {
		if date, err := time.Parse(layout, strings.TrimSpace(value)); err == nil {
			return date.Format(etfRuleRuntimeTimestampDateLayout)
		}
	}
	return ""
}

func fetchEastmoneyFundDividends(client *http.Client, code string) ([]cashDividendEvent, error) {
	endpoint := "https://fundf10.eastmoney.com/fhsp_" + url.PathEscape(normalizeFundSymbol(code)) + ".html"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,*/*")
	req.Header.Set("Referer", "https://fundf10.eastmoney.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; holds-website etf rule updater)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("fund dividend request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	return parseEastmoneyFundDividends(body)
}

func parseEastmoneyFundDividends(body []byte) ([]cashDividendEvent, error) {
	rowPattern := regexp.MustCompile(`(?is)<tr[^>]*>.*?</tr>`)
	cellPattern := regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)
	rows := rowPattern.FindAllString(string(body), -1)
	events := make([]cashDividendEvent, 0, len(rows))
	for _, row := range rows {
		cells := cellPattern.FindAllStringSubmatch(row, -1)
		if len(cells) < 4 {
			continue
		}
		exDate := htmlPlainText(cells[2][1])
		if _, err := time.Parse("2006-01-02", exDate); err != nil {
			continue
		}
		amountText := htmlPlainText(cells[3][1])
		amount, err := fundDividendAmount(amountText)
		if err != nil || amount <= 0 {
			continue
		}
		events = append(events, cashDividendEvent{Date: exDate, Amount: amount})
	}
	if len(events) == 0 {
		return nil, errors.New("missing fund dividend rows")
	}
	return events, nil
}

func fundDividendAmount(value string) (float64, error) {
	amount, err := firstTextNumber(value)
	if err != nil {
		return 0, errors.New("dividend amount not found")
	}
	return amount, nil
}

func trailingFundDividendAmount(events []cashDividendEvent, referenceDate string) (float64, error) {
	reference, err := time.Parse("2006-01-02", referenceDate)
	if err != nil {
		return 0, err
	}
	cutoff := reference.AddDate(-1, 0, 0)
	total := 0.0
	for _, event := range events {
		eventDate, err := time.Parse("2006-01-02", event.Date)
		if err != nil || event.Amount <= 0 {
			continue
		}
		if eventDate.Before(cutoff) || eventDate.After(reference) {
			continue
		}
		total += event.Amount
	}
	if total <= 0 {
		return 0, errors.New("missing trailing fund dividend")
	}
	return total, nil
}

func normalizeMultplDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"Jan 2, 2006", "January 2, 2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format("2006-01-02")
		}
	}
	return ""
}

func percentileRank(value float64, values []float64) float64 {
	if len(values) == 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	sortedValues := append([]float64(nil), values...)
	sort.Float64s(sortedValues)
	count := 0
	for _, item := range sortedValues {
		if item <= value {
			count++
		}
	}
	return float64(count) / float64(len(sortedValues))
}

const worldPERatioNasdaq100URL = "https://worldperatio.com/index/nasdaq-100/"

type worldPERatioSnapshot struct {
	CurrentPE  float64
	Average10Y float64
	StdDev10Y  float64
	ZScore     float64
	Percentile float64
	Date       string
}

func fetchWorldPERatioNasdaq100(client *http.Client, endpoint string) (worldPERatioSnapshot, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return worldPERatioSnapshot{}, err
	}
	req.Header.Set("User-Agent", "holds-website etf rule updater")
	req.Header.Set("Referer", "https://worldperatio.com/")

	resp, err := client.Do(req)
	if err != nil {
		return worldPERatioSnapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return worldPERatioSnapshot{}, fmt.Errorf("worldperatio request failed: %s %s", resp.Status, strings.TrimSpace(string(body)))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return worldPERatioSnapshot{}, err
	}
	return parseWorldPERatioNasdaq100(body)
}

func parseWorldPERatioNasdaq100(body []byte) (worldPERatioSnapshot, error) {
	pageText := htmlPlainText(string(body))
	currentPE, date, err := parseWorldPERatioCurrentPE(pageText)
	if err != nil {
		return worldPERatioSnapshot{}, err
	}

	row := worldPERatioPeriodRow(string(body), "Last 10Y")
	if strings.TrimSpace(row) == "" {
		return worldPERatioSnapshot{}, errors.New("missing Last 10Y row")
	}
	cellPattern := regexp.MustCompile("(?is)<td[^>]*>(.*?)</td>")
	cellMatches := cellPattern.FindAllStringSubmatch(row, -1)
	if len(cellMatches) < 6 {
		return worldPERatioSnapshot{}, errors.New("incomplete Last 10Y row")
	}
	cells := make([]string, 0, len(cellMatches))
	for _, match := range cellMatches {
		cells = append(cells, htmlPlainText(match[1]))
	}

	average, err := firstTextNumber(cells[1])
	if err != nil {
		return worldPERatioSnapshot{}, fmt.Errorf("missing Last 10Y average: %w", err)
	}
	stdDev, err := firstTextNumber(cells[2])
	if err != nil {
		return worldPERatioSnapshot{}, fmt.Errorf("missing Last 10Y standard deviation: %w", err)
	}
	if stdDev <= 0 {
		return worldPERatioSnapshot{}, errors.New("invalid Last 10Y standard deviation")
	}
	zScore, err := firstSigmaValue(cells[5])
	if err != nil {
		zScore = (currentPE - average) / stdDev
	}
	return worldPERatioSnapshot{
		CurrentPE:  currentPE,
		Average10Y: average,
		StdDev10Y:  stdDev,
		ZScore:     zScore,
		Percentile: normalPercentileFromZ(zScore),
		Date:       date,
	}, nil
}

func worldPERatioPeriodRow(pageHTML string, period string) string {
	rowPattern := regexp.MustCompile("(?is)<tr[^>]*>.*?</tr>")
	periodPattern := regexp.MustCompile("(?i)\\b" + regexp.QuoteMeta(period) + "\\b")
	for _, row := range rowPattern.FindAllString(pageHTML, -1) {
		if periodPattern.MatchString(htmlPlainText(row)) {
			return row
		}
	}
	return ""
}

func parseWorldPERatioCurrentPE(pageText string) (float64, string, error) {
	pattern := regexp.MustCompile("(?i)Price-to-Earnings\\s*\\(P/E\\)\\s+Ratio\\s+for\\s+Nasdaq\\s+100\\s+Index\\s+is\\s+([0-9]+(?:\\.[0-9]+)?)\\s*,\\s+calculated\\s+on\\s+([0-9]{1,2}\\s+[A-Za-z]+\\s+[0-9]{4})")
	match := pattern.FindStringSubmatch(pageText)
	if len(match) < 3 {
		return 0, "", errors.New("missing current P/E paragraph")
	}
	currentPE, err := strconv.ParseFloat(match[1], 64)
	if err != nil || currentPE <= 0 {
		return 0, "", errors.New("invalid current P/E")
	}
	date := normalizeWorldPERatioDate(match[2])
	if date == "" {
		return 0, "", fmt.Errorf("unsupported current P/E date: %s", match[2])
	}
	return currentPE, date, nil
}

func normalizeWorldPERatioDate(value string) string {
	value = strings.TrimSpace(value)
	for _, layout := range []string{"02 January 2006", "2 January 2006", "02 Jan 2006", "2 Jan 2006", "2006-01-02"} {
		if date, err := time.Parse(layout, value); err == nil {
			return date.Format("2006-01-02")
		}
	}
	return ""
}

func htmlPlainText(value string) string {
	withoutScripts := regexp.MustCompile("(?is)<script[^>]*>.*?</script>|<style[^>]*>.*?</style>").ReplaceAllString(value, " ")
	withoutTags := regexp.MustCompile("(?is)<[^>]+>").ReplaceAllString(withoutScripts, " ")
	decoded := html.UnescapeString(withoutTags)
	return strings.Join(strings.Fields(decoded), " ")
}

func firstTextNumber(value string) (float64, error) {
	pattern := regexp.MustCompile("[-+]?[0-9]+(?:\\.[0-9]+)?")
	match := pattern.FindString(value)
	if match == "" {
		return 0, errors.New("number not found")
	}
	return strconv.ParseFloat(match, 64)
}

func firstSigmaValue(value string) (float64, error) {
	pattern := regexp.MustCompile("([-+]?[0-9]+(?:\\.[0-9]+)?)\\s*(?:σ|sigma)")
	match := pattern.FindStringSubmatch(value)
	if len(match) < 2 {
		return 0, errors.New("sigma value not found")
	}
	return strconv.ParseFloat(match[1], 64)
}

func normalPercentileFromZ(zScore float64) float64 {
	if math.IsNaN(zScore) || math.IsInf(zScore, 0) {
		return 0
	}
	percentile := 0.5 * (1 + math.Erf(zScore/math.Sqrt2))
	if percentile < 0 {
		return 0
	}
	if percentile > 1 {
		return 1
	}
	return percentile
}

func evaluateA500Rule(inputs etfRuleInputs) etfRuleEvaluation {
	pePercentile := valueOrNaN(inputs.ValuationPercentile)
	spreadPercentile := valueOrNaN(inputs.EarningsYieldSpreadPercentile)
	if !known(pePercentile) || !known(spreadPercentile) {
		return completeRule("one", "估值数据缺失，按中性系数V=1处理；全收益回撤仍是唯一主触发")
	}
	switch {
	case pePercentile <= 0.30 && spreadPercentile >= 0.70:
		return completeRule("oneHalf", "估值便宜：PE不高于30%分位且股债利差不低于70%分位，机会仓金额系数V=1.25")
	case pePercentile >= 0.75 && spreadPercentile <= 0.30:
		return completeRule("quarter", "估值偏贵：PE不低于75%分位且股债利差不高于30%分位，机会仓金额系数V=0.5")
	default:
		return completeRule("one", "估值中性：机会仓金额系数V=1，全收益回撤决定是否触发")
	}
}

func evaluateSP500Rule(inputs etfRuleInputs) etfRuleEvaluation {
	pePercentile := valueOrNaN(inputs.ValuationPercentile)
	spreadPercentile := valueOrNaN(inputs.EarningsYieldSpreadPercentile)
	if !known(pePercentile) || !known(spreadPercentile) {
		return pendingRule("需要同源未来PE分位和盈利收益率利差分位判断标普机会仓估值")
	}
	switch {
	case pePercentile < 0.40 && spreadPercentile > 0.60:
		return completeRule("oneHalf", "估值便宜：未来PE低于40%分位且盈利利差高于60%分位；只调整机会仓金额，不单独触发买入")
	case pePercentile > 0.80 && spreadPercentile < 0.20:
		return completeRule("quarter", "估值昂贵：未来PE高于80%分位且盈利利差低于20%分位；-8%和-12%档金额减半")
	default:
		return completeRule("one", "估值中性：SPTR全收益回撤决定档位，估值只调整本档金额")
	}
}

func evaluateDividendLowVolRule(inputs etfRuleInputs) etfRuleEvaluation {
	valuationScore := valueOrNaN(inputs.ValuationScore)
	spreadPercentile := valueOrNaN(inputs.DividendSpreadPercentile)
	if known(valuationScore) {
		base := dividendAttractivenessBaseLevel(valuationScore)
		return completeRule(base, "按75%股债利差分位+25%PB便宜度计算估值得分V；场外基础定投不择时")
	}
	if known(spreadPercentile) {
		base := dividendAttractivenessBaseLevel(spreadPercentile)
		return completeRule(base, "PB分位待数据，暂按股债利差分位显示水位；V未确认前不触发场内大额买入")
	}
	return pendingRule("需要标的指数股债利差分位判断水位；不单独使用PE、绝对股息率或恐慌指数")
}

func dividendLowVolValuationScore(spreadPercentile float64, pbPercentile float64) float64 {
	return 0.75*spreadPercentile + 0.25*(1-pbPercentile)
}

func evaluateNasdaq100Rule(inputs etfRuleInputs) etfRuleEvaluation {
	pePercentile := valueOrNaN(inputs.ValuationPercentile)
	spreadPercentile := valueOrNaN(inputs.EarningsYieldSpreadPercentile)
	if !known(pePercentile) || !known(spreadPercentile) {
		return pendingRule("需要同源未来PE分位和盈利收益率利差分位判断纳指机会仓估值")
	}
	switch {
	case pePercentile < 0.30 || spreadPercentile > 0.70:
		return completeRule("oneHalf", "估值便宜：未来PE低于30%分位或盈利利差高于70%分位；只调整机会仓金额，不单独触发买入")
	case pePercentile > 0.80 && spreadPercentile < 0.20:
		return completeRule("quarter", "估值昂贵：未来PE高于80%分位且盈利利差低于20%分位；-10%和-15%档金额减半")
	default:
		return completeRule("one", "估值中性：XNDX全收益回撤决定档位，估值只调整本档金额")
	}
}

func evaluatePEPercentileWaterLevel(inputs etfRuleInputs, name string) etfRuleEvaluation {
	valuation := valueOrNaN(inputs.ValuationPercentile)
	if !known(valuation) {
		return pendingRule("需要" + name + " PE分位判断水位")
	}
	level := percentileBaseLevel(valuation, 0.80, 0.60, 0.40, 0.20)
	return completeRule(level, "按PE分位判断水位；基础定投不根据估值择时")
}

func percentileBaseLevel(value float64, quarterThreshold float64, halfThreshold float64, oneThreshold float64, oneHalfThreshold float64) string {
	switch {
	case value > quarterThreshold:
		return "quarter"
	case value >= halfThreshold:
		return "half"
	case value >= oneThreshold:
		return "one"
	case value >= oneHalfThreshold:
		return "oneHalf"
	default:
		return "two"
	}
}

func dividendYieldBaseLevel(yield float64) string {
	switch {
	case yield < 0.047:
		return "quarter"
	case yield <= 0.050:
		return "half"
	case yield <= 0.058:
		return "one"
	case yield <= 0.062:
		return "oneHalf"
	default:
		return "two"
	}
}

func dividendAttractivenessBaseLevel(percentile float64) string {
	switch {
	case percentile < 0.20:
		return "quarter"
	case percentile < 0.40:
		return "half"
	case percentile < 0.60:
		return "one"
	case percentile <= 0.80:
		return "oneHalf"
	default:
		return "two"
	}
}

func lowerETFRuleLevel(left string, right string) string {
	order := map[string]int{"quarter": 0, "half": 1, "one": 2, "oneHalf": 3, "two": 4}
	if order[left] <= order[right] {
		return left
	}
	return right
}

func zScoreBaseLevel(zScore float64) string {
	switch {
	case zScore > 2:
		return "quarter"
	case zScore >= 1:
		return "half"
	case zScore >= -1:
		return "one"
	case zScore >= -2:
		return "oneHalf"
	default:
		return "two"
	}
}

func downshiftLevel(level string) string {
	switch level {
	case "two":
		return "oneHalf"
	case "oneHalf":
		return "one"
	case "one":
		return "half"
	case "half":
		return "quarter"
	default:
		return "quarter"
	}
}

func completeRule(level string, reason string) etfRuleEvaluation {
	return etfRuleEvaluation{Level: level, Complete: true, Reason: reason}
}

func partialRule(level string, reason string) etfRuleEvaluation {
	return etfRuleEvaluation{Level: level, Complete: false, Reason: reason}
}

func pendingRule(reason string) etfRuleEvaluation {
	return etfRuleEvaluation{Complete: false, Reason: reason}
}

func valueOrNaN(value *float64) float64 {
	if value == nil {
		return math.NaN()
	}
	return *value
}

func known(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func percentMetric(value float64) *float64 {
	percent := value * 100
	return &percent
}

func floatMetric(value float64) *float64 {
	return &value
}

func configValuationMetricUnit(config etfRuleConfig) string {
	if config.ValuationMetricKey == "peZScore" {
		return "σ"
	}
	return "%"
}

func runtimeETFRuleStatusList(records map[string]ETFRuleStatus) []ETFRuleStatus {
	keys := make([]string, 0, len(records))
	for key := range records {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	list := make([]ETFRuleStatus, 0, len(keys))
	for _, key := range keys {
		list = append(list, records[key])
	}
	return list
}

func mergeETFRuleStatusWithExisting(next ETFRuleStatus, existing ETFRuleStatus, now time.Time) ETFRuleStatus {
	if strings.TrimSpace(existing.Symbol) == "" {
		if config, ok := etfRuleConfigBySymbol(next.Symbol); ok {
			next = enforceETFRuleStatusConfidence(next, config, now)
			return stabilizeETFRuleLevel(next, ETFRuleStatus{}, config)
		}
		return next
	}
	existingMetrics := map[string]ETFRuleMetric{}
	for _, metric := range existing.Metrics {
		if strings.TrimSpace(metric.Key) != "" {
			existingMetrics[metric.Key] = metric
		}
	}
	usedFallback := false
	for i := range next.Metrics {
		if next.Metrics[i].Available {
			continue
		}
		if retiredDividendLowVolProxyMetric(next.Symbol, next.Metrics[i]) {
			continue
		}
		previous, ok := existingMetrics[next.Metrics[i].Key]
		if !ok || !previous.Available {
			continue
		}
		if next.Symbol == "008163" && dividendLowVolValuationMetric(previous.Key) && !dividendLowVolStatusUsesBasketProxy(existing) {
			continue
		}
		previous.Error = ""
		previous.QualityState = etfQualityDegraded
		previous.QualityMessage = "本次更新失败，沿用上次成功值"
		next.Metrics[i] = previous
		usedFallback = true
	}
	if usedFallback {
		next = refreshETFRuleStatusFromMetrics(next)
	}
	if config, ok := etfRuleConfigBySymbol(next.Symbol); ok {
		next = enforceETFRuleStatusConfidence(next, config, now)
		next = stabilizeETFRuleLevel(next, existing, config)
	}
	return next
}

func dividendLowVolValuationMetric(key string) bool {
	switch strings.TrimSpace(key) {
	case "dividendYield", "china10YBondYield", "dividendSpread", "dividendSpreadPercentile", "indexPB", "pbPercentile", "valuationScore", "basketCoverage":
		return true
	default:
		return false
	}
}

func dividendLowVolStatusUsesBasketProxy(status ETFRuleStatus) bool {
	for _, source := range status.Sources {
		if strings.Contains(source.Name, "515450申购赎回篮子") {
			return true
		}
	}
	return false
}

func retiredDividendLowVolProxyMetric(symbol string, metric ETFRuleMetric) bool {
	if symbol != "008163" {
		return false
	}
	if metric.Key != "dividendSpread" && metric.Key != "dividendSpreadPercentile" {
		return false
	}
	return strings.Contains(metric.Error, "基金现金分红不代替指数股息率")
}

func stabilizeETFRuleLevel(next ETFRuleStatus, existing ETFRuleStatus, config etfRuleConfig) ETFRuleStatus {
	if !next.Complete || strings.TrimSpace(next.Level) == "" {
		return next
	}
	observationDate := etfRuleStatusObservationDate(next)
	if strings.TrimSpace(existing.Level) == "" || !existing.Complete {
		next.LevelUpdatedAt = firstNonEmpty(observationDate, next.AsOf)
		return clearETFRulePendingLevel(next)
	}
	if next.Level == existing.Level {
		next.LevelUpdatedAt = firstNonEmpty(existing.LevelUpdatedAt, existing.AsOf, observationDate)
		return clearETFRulePendingLevel(next)
	}
	pendingDays := 1
	pendingSince := observationDate
	pendingAsOf := observationDate
	if existing.PendingLevel == next.Level {
		pendingDays = existing.PendingDays
		pendingSince = firstNonEmpty(existing.PendingSince, observationDate)
		pendingAsOf = firstNonEmpty(existing.PendingAsOf, observationDate)
		if observationDate != "" && observationDate > pendingAsOf {
			pendingDays++
			pendingAsOf = observationDate
		}
	}
	if pendingDays >= 5 {
		next.LevelUpdatedAt = firstNonEmpty(observationDate, next.AsOf)
		next.Reason = strings.TrimSpace(next.Reason + "；跨档边界已连续5个交易日确认")
		return clearETFRulePendingLevel(next)
	}
	candidateLevel := next.Level
	candidateLabel := config.Levels[candidateLevel].Label
	next.PendingLevel = candidateLevel
	next.PendingLevelLabel = candidateLabel
	next.PendingSince = pendingSince
	next.PendingAsOf = pendingAsOf
	next.PendingDays = pendingDays
	next.Level = existing.Level
	next.LevelLabel = config.Levels[existing.Level].Label
	next.MonthlyAmount = config.Monthly[existing.Level]
	next.WeeklyAmount = config.Weekly[existing.Level]
	next.LevelUpdatedAt = firstNonEmpty(existing.LevelUpdatedAt, existing.AsOf)
	next.Reason = fmt.Sprintf("候选水位%s跨档确认中，继续显示%s", candidateLabel, next.LevelLabel)
	return next
}

func clearETFRulePendingLevel(status ETFRuleStatus) ETFRuleStatus {
	status.PendingLevel = ""
	status.PendingLevelLabel = ""
	status.PendingSince = ""
	status.PendingAsOf = ""
	status.PendingDays = 0
	return status
}

func etfRuleStatusObservationDate(status ETFRuleStatus) string {
	for _, metric := range status.Metrics {
		if metric.Key == "drawdown3y" && metric.Available && strings.TrimSpace(metric.AsOf) != "" {
			return strings.TrimSpace(metric.AsOf)
		}
	}
	return strings.TrimSpace(status.AsOf)
}

func refreshETFRuleStatusFromMetrics(status ETFRuleStatus) ETFRuleStatus {
	config, ok := etfRuleConfigBySymbol(status.Symbol)
	if !ok {
		return status
	}
	inputs := etfRuleInputs{}
	for _, metric := range status.Metrics {
		if !metric.Available || metric.Value == nil {
			continue
		}
		switch metric.Key {
		case "drawdown3y":
			value := *metric.Value / 100
			inputs.Drawdown = &value
			inputs.DrawdownAsOf = metric.AsOf
		case "dividendYield":
			value := *metric.Value / 100
			inputs.DividendYield = &value
			inputs.ValuationAsOf = metric.AsOf
		case "dividendYieldPercentile":
			value := *metric.Value / 100
			inputs.DividendYieldPercentile = &value
			inputs.ValuationAsOf = metric.AsOf
		case "dividendSpreadPercentile":
			value := *metric.Value / 100
			inputs.DividendSpreadPercentile = &value
			inputs.ValuationAsOf = metric.AsOf
		case "pbPercentile":
			value := *metric.Value / 100
			inputs.PBPercentile = &value
			inputs.ValuationAsOf = metric.AsOf
		case "valuationScore":
			value := *metric.Value / 100
			inputs.ValuationScore = &value
			inputs.ValuationAsOf = metric.AsOf
		case "earningsYieldSpreadPercentile":
			value := *metric.Value / 100
			inputs.EarningsYieldSpreadPercentile = &value
			inputs.ValuationAsOf = metric.AsOf
		case config.ValuationMetricKey:
			if metric.Key == "peZScore" || strings.TrimSpace(metric.Unit) == "σ" {
				value := *metric.Value
				inputs.ValuationZScore = &value
			} else {
				value := *metric.Value / 100
				inputs.ValuationPercentile = &value
			}
			inputs.ValuationAsOf = metric.AsOf
		}
	}
	evaluation := config.Evaluate(inputs)
	level := config.Levels[evaluation.Level]
	status.Level = evaluation.Level
	status.LevelLabel = level.Label
	status.MonthlyAmount = config.Monthly["one"]
	status.WeeklyAmount = config.Weekly["one"]
	status.Complete = evaluation.Complete
	status.Reason = evaluation.Reason
	status.AsOf = firstNonEmpty(inputs.DrawdownAsOf, inputs.ValuationAsOf, status.AsOf)
	if status.Level == "" {
		status.LevelLabel = "待数据"
	}
	return status
}

func etfRuleConfigBySymbol(symbol string) (etfRuleConfig, bool) {
	normalized := normalizeFundSymbol(symbol)
	for _, config := range etfRuleConfigs {
		if normalizeFundSymbol(config.Symbol) == normalized {
			return config, true
		}
	}
	return etfRuleConfig{}, false
}

func enforceETFRuleStatusConfidence(status ETFRuleStatus, config etfRuleConfig, now time.Time) ETFRuleStatus {
	if !status.Complete {
		return status
	}
	if len(etfRuleStatusConfidenceIssues(status, config, now)) == 0 {
		return status
	}
	status.Complete = false
	if strings.TrimSpace(status.Reason) == "" {
		status.Reason = "等待指标刷新"
	}
	return status
}

func etfRuleStatusConfidenceIssues(status ETFRuleStatus, config etfRuleConfig, now time.Time) []string {
	issues := []string{}
	metricsByKey := map[string]ETFRuleMetric{}
	for _, metric := range status.Metrics {
		if strings.TrimSpace(metric.Key) == "" {
			continue
		}
		metricsByKey[metric.Key] = metric
	}
	for _, key := range etfRuleRequiredValuationMetricKeys(status, config) {
		metric, ok := metricsByKey[key]
		if !ok {
			issues = append(issues, key+"缺失")
			continue
		}
		issues = append(issues, etfRuleMetricConfidenceIssues(metric, config, now)...)
	}
	if len(status.Sources) < 2 {
		issues = append(issues, "数据源不足")
	}
	return issues
}

func etfRuleRequiredValuationMetricKeys(status ETFRuleStatus, config etfRuleConfig) []string {
	if config.Symbol == "022434" {
		return []string{"drawdown3y", "etfPremium", "bidAskSpread"}
	}
	if config.Symbol == "018738" {
		return []string{"forwardPE", "forwardPEPercentile", "us10YBondYield", "earningsYieldSpreadPercentile", "vix", "cnyTotalReturnDrawdown", "qdiiPremium"}
	}
	if config.Symbol == "021000" {
		return []string{"forwardPE", "forwardPEPercentile", "us10YBondYield", "earningsYieldSpreadPercentile", "vxn", "cnyTotalReturnDrawdown", "qdiiPremium"}
	}
	if config.Symbol != "008163" {
		return []string{config.ValuationMetricKey}
	}
	return []string{"dividendYield", "china10YBondYield", "dividendSpreadPercentile", "pbPercentile", "valuationScore", "basketCoverage"}
}

func etfRuleMetricConfidenceIssues(metric ETFRuleMetric, config etfRuleConfig, now time.Time) []string {
	issues := []string{}
	if !metric.Available || metric.Value == nil {
		return append(issues, firstNonEmpty(metric.Label, metric.Key)+"不可用")
	}
	value := *metric.Value
	if math.IsNaN(value) || math.IsInf(value, 0) || !etfRuleMetricValueInExpectedRange(metric, config, value) {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"数值异常")
	}
	if strings.TrimSpace(metric.Error) != "" {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"沿用旧值")
	}
	metricDate, err := time.Parse(etfRuleRuntimeTimestampDateLayout, strings.TrimSpace(metric.AsOf))
	if err != nil {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"日期缺失")
		return issues
	}
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, metricDate.Location())
	if metricDate.After(nowDate.Add(24 * time.Hour)) {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"日期异常")
		return issues
	}
	maxAgeDays := etfRuleMetricMaxAgeDays(metric, config)
	if nowDate.Sub(metricDate).Hours()/24 > float64(maxAgeDays) {
		issues = append(issues, firstNonEmpty(metric.Label, metric.Key)+"过期")
	}
	return issues
}

func etfRuleMetricValueInExpectedRange(metric ETFRuleMetric, config etfRuleConfig, value float64) bool {
	switch metric.Key {
	case "drawdown3y", "cnyTotalReturnDrawdown":
		return value >= 0 && value <= 100
	case "dividendYield", "china10YBondYield", "us10YBondYield":
		return value > 0 && value <= 20
	case "dividendSpread", "earningsYieldSpread":
		return value >= -20 && value <= 20
	case "forwardPE", "indexPE":
		return value >= 5 && value <= 100
	case "vix", "vxn":
		return value >= 5 && value <= 150
	case "usdCny":
		return value >= 4 && value <= 12
	case "nasdaqFuturesChange", "sp500FuturesChange", "qdiiPremium", "etfPremium", "bidAskSpread", "openingGap", "fiveDayReturn", "forwardEarningsRevision3m":
		return value >= -50 && value <= 100
	case "tacticalMarketPrice", "tacticalOfficialNAV", "tacticalEstimatedNAV":
		return value > 0 && value <= 100
	case "indexPB":
		return value >= 0.05 && value <= 20
	case "basketCoverage":
		return value >= dividendLowVolMinimumCoverage*100 && value <= 100
	case "dividendYieldPercentile", "dividendSpreadPercentile", "pbPercentile", "valuationScore", "forwardPEPercentile", "pePercentile", "earningsYieldSpreadPercentile", "rv20Percentile", "breadthBelowMA20":
		return value >= 0 && value <= 100
	case "rv20":
		return value >= 0 && value <= 200
	case "volumeRatio":
		return value >= 0 && value <= 20
	case "totalReturnClose", "totalReturnPeak":
		return value > 0
	case config.ValuationMetricKey:
		if metric.Key == "peZScore" || strings.TrimSpace(metric.Unit) == "σ" {
			return value >= -6 && value <= 6
		}
		return value >= 0 && value <= 100
	default:
		return true
	}
}

func etfRuleMetricMaxAgeDays(metric ETFRuleMetric, config etfRuleConfig) int {
	if config.Symbol == "008163" && dividendLowVolValuationMetric(metric.Key) {
		return etfRuleWeeklyMetricMaxAgeDays
	}
	if metric.Key == "dividendYieldPercentile" || metric.Key == "dividendSpreadPercentile" || metric.Key == "pbPercentile" || metric.Key == "valuationScore" || metric.Key == "forwardPE" || metric.Key == "forwardPEPercentile" || metric.Key == "pePercentile" || metric.Key == "earningsYieldSpreadPercentile" {
		return etfRuleMonthlyMetricMaxAgeDays
	}
	if metric.Key == config.ValuationMetricKey && (config.Symbol == "018738" || config.Symbol == "021000") {
		return etfRuleMonthlyMetricMaxAgeDays
	}
	return etfRuleDailyMetricMaxAgeDays
}
