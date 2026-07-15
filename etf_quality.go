package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	etfSourceOfficial       = "official"
	etfSourceValidatedProxy = "validated_proxy"
	etfSourceFallback       = "fallback"

	etfQualityFresh       = "fresh"
	etfQualityDegraded    = "degraded"
	etfQualityStale       = "stale"
	etfQualityUnavailable = "unavailable"
	etfQualityDisputed    = "disputed"

	etfSignalHealthy  = "healthy"
	etfSignalDegraded = "degraded"
	etfSignalBlocked  = "blocked"

	etfExecutionReady     = "ready"
	etfExecutionReference = "reference"
	etfExecutionBlocked   = "blocked"

	etfExecutionQuoteMaxAgeSeconds  = 90
	etfExecutionInputMaxSkewSeconds = 30
)

func applyETFDataQuality(statuses []ETFRuleStatus, now time.Time) {
	for i := range statuses {
		applyETFStatusDataQuality(&statuses[i], now)
	}
}

func applyETFStatusDataQuality(status *ETFRuleStatus, now time.Time) {
	if status == nil {
		return
	}
	status.BlockingReasons = nil
	status.SignalHealth = etfSignalHealthy
	status.ExecutionHealth = etfExecutionReady
	status.QualityUpdatedAt = now.Format("2006-01-02 15:04:05")
	for i := range status.Metrics {
		metric := &status.Metrics[i]
		if strings.TrimSpace(metric.FetchedAt) == "" {
			metric.FetchedAt = status.UpdatedAt
		}
		metric.SourceTier = etfMetricSourceTier(*status, *metric)
		if len(metric.SourceIDs) == 0 {
			metric.SourceIDs = etfMetricSourceIDs(*status, *metric)
		}
		metric.MaxAgeSeconds = etfMetricMaxAgeSeconds(status.Symbol, metric.Key)
		metric.QualityState, metric.QualityMessage = etfMetricQuality(*metric, now)
	}

	drawdown := findETFRuleMetric(status.Metrics, "drawdown3y")
	status.SignalAsOf = metricAsOf(drawdown)
	if drawdown == nil || !drawdown.Available || drawdown.Value == nil {
		status.SignalHealth = etfSignalBlocked
		status.BlockingReasons = append(status.BlockingReasons, "核心全收益回撤缺失")
	} else {
		switch drawdown.QualityState {
		case etfQualityStale:
			status.SignalHealth = etfSignalBlocked
			status.BlockingReasons = append(status.BlockingReasons, "核心全收益回撤已过期")
		case etfQualityDisputed:
			status.SignalHealth = etfSignalBlocked
			status.BlockingReasons = append(status.BlockingReasons, "核心全收益回撤校验冲突")
		case etfQualityDegraded:
			status.SignalHealth = etfSignalDegraded
		}
	}

	executionMetrics := etfExecutionMetricKeys(status.Symbol)
	executionTimes := []time.Time{}
	for _, key := range executionMetrics {
		metric := findETFRuleMetric(status.Metrics, key)
		if metric == nil || !metric.Available || metric.Value == nil {
			status.ExecutionHealth = etfExecutionBlocked
			status.BlockingReasons = append(status.BlockingReasons, etfMetricDisplayName(status.Symbol, key)+"缺失")
			continue
		}
		if fetchedAt, ok := parseETFMetricTime(metric.FetchedAt); ok {
			executionTimes = append(executionTimes, fetchedAt)
		}
		if metric.QualityState == etfQualityDisputed {
			status.ExecutionHealth = etfExecutionBlocked
			status.BlockingReasons = append(status.BlockingReasons, etfMetricDisplayName(status.Symbol, key)+"校验冲突")
		}
	}
	if len(executionTimes) > 0 {
		oldest := executionTimes[0]
		for _, item := range executionTimes[1:] {
			if item.Before(oldest) {
				oldest = item
			}
		}
		status.ExecutionAsOf = oldest.Format("2006-01-02 15:04:05")
	}
	if status.ExecutionHealth != etfExecutionBlocked {
		if chinaAshareMarketOpen(now) {
			if len(executionTimes) == 0 || now.Sub(oldestETFTime(executionTimes)) > etfExecutionQuoteMaxAgeSeconds*time.Second {
				status.ExecutionHealth = etfExecutionBlocked
				status.BlockingReasons = append(status.BlockingReasons, "场内执行行情超过90秒")
			} else if !etfExecutionInputsAligned(*status) {
				status.ExecutionHealth = etfExecutionBlocked
				status.BlockingReasons = append(status.BlockingReasons, "场内执行数据时间未对齐")
			}
		} else {
			status.ExecutionHealth = etfExecutionReference
		}
	}

	if status.SignalHealth == etfSignalHealthy && etfStatusHasDegradedAuxiliary(*status) {
		status.SignalHealth = etfSignalDegraded
	}
	for i := range status.Sources {
		status.Sources[i].Tier = etfSourceTierForName(status.Sources[i].Name)
		status.Sources[i].QualityState = sourceQualityState(status.Sources[i].Tier)
		status.Sources[i].AsOf = firstNonEmpty(status.SignalAsOf, status.AsOf)
	}
}

func etfMetricQuality(metric ETFRuleMetric, now time.Time) (string, string) {
	if !metric.Available || metric.Value == nil {
		return etfQualityUnavailable, "本次未取得可用数据"
	}
	if metric.ValidationValue != nil && metric.ConflictTolerance > 0 && math.Abs(*metric.Value-*metric.ValidationValue) > metric.ConflictTolerance {
		source := strings.TrimSpace(metric.ValidationSource)
		if source == "" {
			source = "验证源"
		}
		return etfQualityDisputed, fmt.Sprintf("主数据与%s偏差超过容忍范围", source)
	}
	if strings.Contains(metric.Error, "冲突") || strings.Contains(metric.QualityMessage, "冲突") {
		return etfQualityDisputed, "主数据与验证数据不一致"
	}
	if metric.QualityState == etfQualityDegraded || strings.Contains(metric.QualityMessage, "沿用上次") {
		return etfQualityDegraded, firstNonEmpty(metric.QualityMessage, "本次更新失败，正在沿用上次成功值")
	}
	if strings.TrimSpace(metric.Error) != "" {
		return etfQualityDegraded, "本次更新失败，正在沿用上次成功值"
	}
	if metric.MaxAgeSeconds > 0 {
		value := metric.FetchedAt
		if metric.MaxAgeSeconds > etfExecutionQuoteMaxAgeSeconds {
			value = firstNonEmpty(metric.AsOf, metric.FetchedAt)
		}
		if parsed, ok := parseETFMetricTime(value); ok && now.Sub(parsed) > time.Duration(metric.MaxAgeSeconds)*time.Second {
			return etfQualityStale, "数据已超过允许时效"
		}
	}
	if metric.SourceTier == etfSourceFallback {
		return etfQualityDegraded, "当前使用经过校验的备用数据"
	}
	return etfQualityFresh, "数据可用"
}

func etfMetricSourceTier(status ETFRuleStatus, metric ETFRuleMetric) string {
	symbol := normalizeFundSymbol(status.Symbol)
	key := metric.Key
	if symbol == "018738" && key == "drawdown3y" && (strings.Contains(metric.Label, "SPY") || statusUsesSource(status, "SPY")) {
		return etfSourceFallback
	}
	if key == "forwardPE" || key == "forwardPEPercentile" || key == "earningsYieldSpreadPercentile" && symbol != "022434" {
		return etfSourceFallback
	}
	if key == "china10YBondYield" || key == "us10YBondYield" || key == "vix" || key == "vxn" {
		return etfSourceOfficial
	}
	if symbol == "022434" && containsETFMetricKey(key, "drawdown3y", "totalReturnClose", "totalReturnPeak", "indexPE", "pePercentile", "earningsYieldSpread", "earningsYieldSpreadPercentile", "rv20", "rv20Percentile", "fiveDayReturn") {
		return etfSourceOfficial
	}
	if symbol == "021000" && containsETFMetricKey(key, "drawdown3y") {
		return etfSourceOfficial
	}
	if symbol == "018738" && containsETFMetricKey(key, "drawdown3y") {
		return etfSourceOfficial
	}
	return etfSourceValidatedProxy
}

func etfMetricSourceIDs(status ETFRuleStatus, metric ETFRuleMetric) []string {
	result := []string{}
	for _, source := range status.Sources {
		name := source.Name
		switch metric.SourceTier {
		case etfSourceOfficial:
			if strings.Contains(name, "官方") || strings.Contains(name, "中证指数") || strings.Contains(name, "美国财政部") || strings.Contains(name, "Cboe") || strings.Contains(name, "Nasdaq") {
				result = append(result, name)
			}
		case etfSourceFallback:
			if strings.Contains(name, "备援") || strings.Contains(name, "History of Market") {
				result = append(result, name)
			}
		default:
			if !strings.Contains(name, "官方说明") {
				result = append(result, name)
			}
		}
	}
	if len(result) > 2 {
		return result[:2]
	}
	return result
}

func etfMetricMaxAgeSeconds(symbol string, key string) int {
	if containsETFMetricKey(key, "tacticalMarketPrice", "tacticalEstimatedNAV", "etfPremium", "bidAskSpread", "openingGap", "qdiiPremium", "sp500FuturesChange", "nasdaqFuturesChange") {
		return etfExecutionQuoteMaxAgeSeconds
	}
	if (normalizeFundSymbol(symbol) == "018738" || normalizeFundSymbol(symbol) == "021000") && key == "usdCny" {
		return etfExecutionQuoteMaxAgeSeconds
	}
	if containsETFMetricKey(key, "forwardPE", "forwardPEPercentile", "earningsYieldSpreadPercentile") && normalizeFundSymbol(symbol) != "022434" {
		return 31 * 24 * 60 * 60
	}
	if normalizeFundSymbol(symbol) == "008163" && containsETFMetricKey(key, "dividendYield", "dividendSpread", "dividendSpreadPercentile", "indexPB", "pbPercentile", "valuationScore", "basketCoverage") {
		return 14 * 24 * 60 * 60
	}
	if containsETFMetricKey(key, "drawdown3y", "totalReturnClose", "totalReturnPeak", "cnyTotalReturnDrawdown", "vix", "vxn", "usdCny") {
		return 7 * 24 * 60 * 60
	}
	return 14 * 24 * 60 * 60
}

func etfExecutionMetricKeys(symbol string) []string {
	switch normalizeFundSymbol(symbol) {
	case "022434":
		return []string{"tacticalMarketPrice", "etfPremium", "bidAskSpread"}
	case "008163":
		return []string{"tacticalMarketPrice"}
	case "018738":
		return []string{"sp500FuturesChange", "usdCny", "tacticalMarketPrice", "tacticalEstimatedNAV", "qdiiPremium"}
	case "021000":
		return []string{"nasdaqFuturesChange", "usdCny", "tacticalMarketPrice", "tacticalEstimatedNAV", "qdiiPremium"}
	default:
		return nil
	}
}

func etfExecutionInputsAligned(status ETFRuleStatus) bool {
	keys := etfExecutionMetricKeys(status.Symbol)
	if len(keys) < 2 {
		return true
	}
	var earliest time.Time
	var latest time.Time
	for _, key := range keys {
		metric := findETFRuleMetric(status.Metrics, key)
		if metric == nil {
			return false
		}
		observedAt, ok := parseETFMetricTime(metric.FetchedAt)
		if !ok {
			return false
		}
		if earliest.IsZero() || observedAt.Before(earliest) {
			earliest = observedAt
		}
		if latest.IsZero() || observedAt.After(latest) {
			latest = observedAt
		}
	}
	return latest.Sub(earliest) <= etfExecutionInputMaxSkewSeconds*time.Second
}

func etfStatusHasDegradedAuxiliary(status ETFRuleStatus) bool {
	for _, metric := range status.Metrics {
		if metric.Key == "drawdown3y" || containsETFMetricKey(metric.Key, etfExecutionMetricKeys(status.Symbol)...) {
			continue
		}
		if metric.QualityState == etfQualityDegraded || metric.QualityState == etfQualityUnavailable || metric.QualityState == etfQualityStale {
			return true
		}
	}
	return false
}

func etfSourceTierForName(name string) string {
	if strings.Contains(name, "备援") || strings.Contains(name, "History of Market") {
		return etfSourceFallback
	}
	if strings.Contains(name, "中证指数") || strings.Contains(name, "美国财政部") || strings.Contains(name, "Cboe") || strings.Contains(name, "Nasdaq官方") || strings.Contains(name, "中债") || strings.Contains(name, "S&P Dow Jones") {
		return etfSourceOfficial
	}
	return etfSourceValidatedProxy
}

func sourceQualityState(tier string) string {
	if tier == etfSourceFallback {
		return etfQualityDegraded
	}
	return etfQualityFresh
}

func statusUsesSource(status ETFRuleStatus, fragment string) bool {
	for _, source := range status.Sources {
		if strings.Contains(source.Name, fragment) {
			return true
		}
	}
	return false
}

func containsETFMetricKey(key string, keys ...string) bool {
	for _, candidate := range keys {
		if key == candidate {
			return true
		}
	}
	return false
}

func parseETFMetricTime(value string) (time.Time, bool) {
	text := strings.TrimSpace(value)
	for _, layout := range []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"} {
		if parsed, err := time.ParseInLocation(layout, text, time.FixedZone("Asia/Shanghai", 8*60*60)); err == nil {
			if layout == "2006-01-02" {
				parsed = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			}
			return parsed, true
		}
	}
	return time.Time{}, false
}

func metricAsOf(metric *ETFRuleMetric) string {
	if metric == nil {
		return ""
	}
	return metric.AsOf
}

func oldestETFTime(values []time.Time) time.Time {
	if len(values) == 0 {
		return time.Time{}
	}
	oldest := values[0]
	for _, value := range values[1:] {
		if value.Before(oldest) {
			oldest = value
		}
	}
	return oldest
}

func chinaAshareMarketOpen(now time.Time) bool {
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	local := now.In(location)
	if local.Weekday() == time.Saturday || local.Weekday() == time.Sunday {
		return false
	}
	minutes := local.Hour()*60 + local.Minute()
	return (minutes >= 9*60+30 && minutes <= 11*60+30) || (minutes >= 13*60 && minutes <= 15*60)
}

func etfMetricDisplayName(symbol string, key string) string {
	labels := map[string]string{
		"tacticalMarketPrice":  "场内价格",
		"tacticalEstimatedNAV": "估算实时净值",
		"etfPremium":           "估算溢价",
		"bidAskSpread":         "买卖价差",
		"qdiiPremium":          "估算溢价",
	}
	if label := labels[key]; label != "" {
		return label
	}
	return fmt.Sprintf("%s %s", normalizeFundSymbol(symbol), key)
}
