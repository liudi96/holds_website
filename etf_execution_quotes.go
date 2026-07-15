package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ETFExecutionQuoteUpdateResponse struct {
	Updated int         `json:"updated"`
	Skipped []QuoteSkip `json:"skipped"`
	State   AppState    `json:"state"`
}

func (s *Server) handleUpdateETFExecutionQuotes(w http.ResponseWriter, r *http.Request) {
	result, err := s.runETFExecutionQuoteUpdate(time.Now())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) runETFExecutionQuoteUpdate(now time.Time) (ETFExecutionQuoteUpdateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := loadState()
	if err != nil {
		return ETFExecutionQuoteUpdateResponse{}, fmt.Errorf("failed to load state: %w", err)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	statuses := make(map[string]ETFRuleStatus, len(state.ETFRuleStatuses))
	for _, status := range state.ETFRuleStatuses {
		statuses[normalizeFundSymbol(status.Symbol)] = status
	}
	type refreshResult struct {
		symbol  string
		metrics []ETFRuleMetric
		err     error
	}
	results := make(chan refreshResult, len(etfExecutionPlanConfigs))
	var wait sync.WaitGroup
	for _, config := range etfExecutionPlanConfigs {
		config := config
		wait.Add(1)
		go func() {
			defer wait.Done()
			metrics, refreshErr := fetchETFExecutionMetrics(client, config.TrackerSymbol, now)
			results <- refreshResult{symbol: config.TrackerSymbol, metrics: metrics, err: refreshErr}
		}()
	}
	wait.Wait()
	close(results)

	updated := 0
	skipped := []QuoteSkip{}
	for result := range results {
		status, ok := statuses[normalizeFundSymbol(result.symbol)]
		if !ok {
			config, _ := etfRuleConfigBySymbol(result.symbol)
			status = ETFRuleStatus{Symbol: result.symbol, Name: config.Name}
		}
		if result.err != nil {
			config, _ := etfExecutionConfig(result.symbol)
			skipped = append(skipped, QuoteSkip{Type: "etf-execution", Symbol: config.TacticalSymbol, Name: config.TacticalName, Error: result.err.Error()})
		}
		if len(result.metrics) == 0 {
			continue
		}
		for _, metric := range result.metrics {
			metric.FetchedAt = now.Format("2006-01-02 15:04:05")
			status.Metrics = upsertETFRuleMetric(status.Metrics, metric)
		}
		status.UpdatedAt = now.Format("2006-01-02 15:04:05")
		applyETFStatusDataQuality(&status, now)
		statuses[normalizeFundSymbol(result.symbol)] = status
		updated++
	}

	statusList := runtimeETFRuleStatusList(statuses)
	if err := saveRuntimeMarketData(nil, statusList, now.Format("2006-01-02 15:04:05")); err != nil {
		return ETFExecutionQuoteUpdateResponse{}, fmt.Errorf("failed to save ETF execution quotes: %w", err)
	}
	if err := hydrateState(&state); err != nil {
		return ETFExecutionQuoteUpdateResponse{}, fmt.Errorf("failed to hydrate ETF execution quotes: %w", err)
	}
	if err := saveState(state); err != nil {
		return ETFExecutionQuoteUpdateResponse{}, fmt.Errorf("failed to save ETF execution state: %w", err)
	}
	s.state = state
	return ETFExecutionQuoteUpdateResponse{Updated: updated, Skipped: skipped, State: state}, nil
}

func fetchETFExecutionMetrics(client *http.Client, trackerSymbol string, now time.Time) ([]ETFRuleMetric, error) {
	switch normalizeFundSymbol(trackerSymbol) {
	case "022434":
		snapshot, issues := fetchA500OpportunitySnapshot(client, now)
		if issues.Trading != nil {
			return unavailableETFExecutionMetrics("159352", "etfPremium", issues.Trading), issues.Trading
		}
		asOf := snapshot.MarketPriceDate
		return []ETFRuleMetric{
			{Key: "tacticalMarketPrice", Label: "159352场内价格", Value: floatMetric(snapshot.MarketPrice), AsOf: asOf, Available: true},
			{Key: "tacticalOfficialNAV", Label: "159352最近公布净值", Value: floatMetric(snapshot.OfficialNAV), AsOf: snapshot.OfficialNAVDate, Available: true},
			{Key: "tacticalEstimatedNAV", Label: "159352估算净值", Value: floatMetric(snapshot.EstimatedNAV), AsOf: asOf, Available: true},
			{Key: "etfPremium", Label: "159352估算溢价", Value: percentMetric(snapshot.Premium), Unit: "%", AsOf: asOf, Available: true},
			{Key: "bidAskSpread", Label: "159352买卖价差", Value: percentMetric(snapshot.BidAskSpread), Unit: "%", AsOf: asOf, Available: true},
			{Key: "openingGap", Label: "159352开盘涨跌", Value: percentMetric(snapshot.OpeningGap), Unit: "%", AsOf: asOf, Available: true},
		}, nil
	case "008163":
		market, err := fetchEastmoneyQuote(client, "515450.SH")
		if err != nil || market.Price <= 0 {
			if err == nil {
				err = fmt.Errorf("invalid 515450 market price")
			}
			return []ETFRuleMetric{{Key: "tacticalMarketPrice", Label: "515450场内价格", Available: false, Error: err.Error()}}, err
		}
		return []ETFRuleMetric{{Key: "tacticalMarketPrice", Label: "515450场内价格", Value: floatMetric(market.Price), AsOf: market.PriceDate, Available: true}}, nil
	case "018738":
		snapshot, issues := fetchSP500OpportunitySnapshot(client, now)
		if issues.Premium != nil {
			return unavailableETFExecutionMetrics("513650", "qdiiPremium", issues.Premium), issues.Premium
		}
		asOf := snapshot.MarketPriceDate
		return []ETFRuleMetric{
			{Key: "sp500FuturesChange", Label: "标普500期货当日变动", Value: percentMetric(snapshot.FuturesChange), Unit: "%", AsOf: snapshot.FuturesDate, Available: true},
			{Key: "usdCny", Label: "USD/CNY执行参考汇率", Value: floatMetric(snapshot.USDToCNY), AsOf: snapshot.FXDate, Available: snapshot.USDToCNY > 0},
			{Key: "tacticalMarketPrice", Label: "513650场内价格", Value: floatMetric(snapshot.MarketPrice), AsOf: asOf, Available: true},
			{Key: "tacticalOfficialNAV", Label: "513650最近公布净值", Value: floatMetric(snapshot.OfficialNAV), AsOf: snapshot.OfficialNAVDate, Available: true},
			{Key: "tacticalEstimatedNAV", Label: "513650估算实时净值", Value: floatMetric(snapshot.EstimatedNAV), AsOf: asOf, Available: true},
			{Key: "qdiiPremium", Label: "513650估算溢价", Value: percentMetric(snapshot.Premium), Unit: "%", AsOf: asOf, Available: true},
		}, nil
	case "021000":
		snapshot, issues := fetchNasdaqOpportunitySnapshot(client, now)
		if issues.Premium != nil {
			return unavailableETFExecutionMetrics("159659", "qdiiPremium", issues.Premium), issues.Premium
		}
		asOf := snapshot.MarketPriceDate
		return []ETFRuleMetric{
			{Key: "nasdaqFuturesChange", Label: "纳指期货当日变动", Value: percentMetric(snapshot.FuturesChange), Unit: "%", AsOf: snapshot.FuturesDate, Available: true},
			{Key: "usdCny", Label: "USD/CNY执行参考汇率", Value: floatMetric(snapshot.USDToCNY), AsOf: snapshot.FXDate, Available: snapshot.USDToCNY > 0},
			{Key: "tacticalMarketPrice", Label: "159659场内价格", Value: floatMetric(snapshot.MarketPrice), AsOf: asOf, Available: true},
			{Key: "tacticalOfficialNAV", Label: "159659最近公布净值", Value: floatMetric(snapshot.OfficialNAV), AsOf: snapshot.OfficialNAVDate, Available: true},
			{Key: "tacticalEstimatedNAV", Label: "159659估算实时净值", Value: floatMetric(snapshot.EstimatedNAV), AsOf: asOf, Available: true},
			{Key: "qdiiPremium", Label: "159659估算溢价", Value: percentMetric(snapshot.Premium), Unit: "%", AsOf: asOf, Available: true},
		}, nil
	default:
		return nil, fmt.Errorf("unknown ETF execution tracker %s", trackerSymbol)
	}
}

func unavailableETFExecutionMetrics(tacticalSymbol string, premiumKey string, err error) []ETFRuleMetric {
	errorText := "场内执行行情暂不可用"
	if err != nil {
		errorText = err.Error()
	}
	return []ETFRuleMetric{
		{Key: "tacticalMarketPrice", Label: tacticalSymbol + "场内价格", Available: false, Error: errorText},
		{Key: "tacticalEstimatedNAV", Label: tacticalSymbol + "估算实时净值", Available: false, Error: errorText},
		{Key: premiumKey, Label: tacticalSymbol + "估算溢价", Unit: "%", Available: false, Error: errorText},
	}
}

func upsertETFRuleMetric(metrics []ETFRuleMetric, metric ETFRuleMetric) []ETFRuleMetric {
	for i := range metrics {
		if metrics[i].Key == metric.Key {
			metrics[i] = metric
			return metrics
		}
	}
	return append(metrics, metric)
}
