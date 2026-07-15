package main

import (
	"testing"
	"time"
)

func TestETFDataQualityBlocksStaleIntradayQuoteButKeepsCloseReference(t *testing.T) {
	value := 7.85
	price := 1.29
	premium := -0.07
	spread := 0.08
	status := ETFRuleStatus{
		Symbol: "022434", UpdatedAt: "2026-07-14 10:00:00",
		Metrics: []ETFRuleMetric{
			{Key: "drawdown3y", Label: "中证A500全收益高点回撤", Value: &value, AsOf: "2026-07-13", Available: true},
			{Key: "tacticalMarketPrice", Label: "159352场内价格", Value: &price, AsOf: "2026-07-14", FetchedAt: "2026-07-14 10:00:00", Available: true},
			{Key: "etfPremium", Label: "159352估算溢价", Value: &premium, AsOf: "2026-07-14", FetchedAt: "2026-07-14 10:00:00", Available: true},
			{Key: "bidAskSpread", Label: "159352买卖价差", Value: &spread, AsOf: "2026-07-14", FetchedAt: "2026-07-14 10:00:00", Available: true},
		},
	}
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	applyETFStatusDataQuality(&status, time.Date(2026, 7, 14, 10, 0, 30, 0, location))
	if status.SignalHealth != etfSignalHealthy || status.ExecutionHealth != etfExecutionReady {
		t.Fatalf("fresh status should be executable: %+v", status)
	}
	applyETFStatusDataQuality(&status, time.Date(2026, 7, 14, 10, 2, 0, 0, location))
	if status.ExecutionHealth != etfExecutionBlocked {
		t.Fatalf("stale intraday quote should block execution: %+v", status)
	}
	applyETFStatusDataQuality(&status, time.Date(2026, 7, 14, 16, 0, 0, 0, location))
	if status.ExecutionHealth != etfExecutionReference {
		t.Fatalf("outside market hours should be a close reference: %+v", status)
	}
}

func TestETFDataQualityUsesNeutralAuxiliaryDegradation(t *testing.T) {
	drawdown := 8.0
	price := 1.2
	status := ETFRuleStatus{
		Symbol: "008163", UpdatedAt: "2026-07-14 16:00:00",
		Metrics: []ETFRuleMetric{
			{Key: "drawdown3y", Label: "515450成立以来总回报回撤", Value: &drawdown, AsOf: "2026-07-14", Available: true},
			{Key: "valuationScore", Label: "估值得分V", Available: false, Error: "source unavailable"},
			{Key: "tacticalMarketPrice", Label: "515450场内价格", Value: &price, AsOf: "2026-07-14", FetchedAt: "2026-07-14 16:00:00", Available: true},
		},
	}
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	applyETFStatusDataQuality(&status, time.Date(2026, 7, 14, 16, 1, 0, 0, location))
	if status.SignalHealth != etfSignalDegraded {
		t.Fatalf("missing auxiliary valuation should degrade, not block: %+v", status)
	}
	if status.ExecutionHealth != etfExecutionReference {
		t.Fatalf("outside-session execution data should remain reference: %+v", status)
	}
}

func TestETFDataQualityBlocksConflictingCoreSignal(t *testing.T) {
	primary := 7.85
	validation := 10.10
	status := ETFRuleStatus{
		Symbol: "022434",
		Metrics: []ETFRuleMetric{{
			Key: "drawdown3y", Label: "中证A500全收益高点回撤", Value: &primary, AsOf: "2026-07-14", Available: true,
			ValidationValue: &validation, ValidationSource: "校验序列", ConflictTolerance: 1.0,
		}},
	}
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	applyETFStatusDataQuality(&status, time.Date(2026, 7, 14, 16, 0, 0, 0, location))
	if status.SignalHealth != etfSignalBlocked || status.Metrics[0].QualityState != etfQualityDisputed {
		t.Fatalf("conflicting core signal must block tactical execution: %+v", status)
	}
}

func TestETFDataQualityBlocksMisalignedQDIIExecutionInputs(t *testing.T) {
	drawdown := 10.0
	futures := -1.0
	fx := 7.2
	price := 1.2
	nav := 1.19
	premium := 0.84
	metrics := []ETFRuleMetric{
		{Key: "drawdown3y", Value: &drawdown, AsOf: "2026-07-13", Available: true},
		{Key: "sp500FuturesChange", Value: &futures, FetchedAt: "2026-07-14 10:00:00", Available: true},
		{Key: "usdCny", Value: &fx, FetchedAt: "2026-07-14 10:00:00", Available: true},
		{Key: "tacticalMarketPrice", Value: &price, FetchedAt: "2026-07-14 10:00:00", Available: true},
		{Key: "tacticalEstimatedNAV", Value: &nav, FetchedAt: "2026-07-14 09:59:20", Available: true},
		{Key: "qdiiPremium", Value: &premium, FetchedAt: "2026-07-14 10:00:00", Available: true},
	}
	status := ETFRuleStatus{Symbol: "018738", Metrics: metrics}
	location := time.FixedZone("Asia/Shanghai", 8*60*60)
	applyETFStatusDataQuality(&status, time.Date(2026, 7, 14, 10, 0, 30, 0, location))
	if status.ExecutionHealth != etfExecutionBlocked {
		t.Fatalf("misaligned QDII inputs must block execution: %+v", status)
	}
	foundReason := false
	for _, reason := range status.BlockingReasons {
		if reason == "场内执行数据时间未对齐" {
			foundReason = true
		}
	}
	if !foundReason {
		t.Fatalf("missing alignment reason: %+v", status.BlockingReasons)
	}
}
