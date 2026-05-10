package main

import (
	"encoding/json"
	"testing"
)

func floatPtr(value float64) *float64 {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func TestValidateResearchAcceptsLegacyJSON(t *testing.T) {
	raw := []byte(`{
		"symbol": "0700.HK",
		"name": "腾讯控股",
		"asOf": "2026-05-01",
		"currency": "HKD",
		"industry": "互联网平台/游戏/广告/金融科技",
		"status": "过渡观察",
		"action": "旧仓持有，新资金等待股息和DCF达标",
		"risk": "政策、地缘和AI投入回报周期需折价",
		"valuation": {
			"intrinsicValue": 508,
			"fairValueRange": "HK$480-560",
			"targetBuyPrice": null,
			"marginOfSafety": 0.09
		},
		"quality": {
			"totalScore": 89,
			"businessModel": 28,
			"moat": 23,
			"governance": 17,
			"financialQuality": 21
		},
		"plan": {
			"rank": 1,
			"priority": "观察",
			"advice": "等待主策略条件达标",
			"discipline": "未达标不追买"
		},
		"notes": "旧导入格式不包含 dividend 和 netCash。"
	}`)

	var research ResearchImport
	if err := json.Unmarshal(raw, &research); err != nil {
		t.Fatalf("unmarshal legacy research: %v", err)
	}
	research = normalizeResearch(research)
	if _, err := validateResearch(research); err != nil {
		t.Fatalf("validate legacy research: %v", err)
	}
	if research.Dividend != nil {
		t.Fatalf("legacy research should keep dividend nil")
	}
	if research.NetCash != nil {
		t.Fatalf("legacy research should keep netCash nil")
	}
	if research.OwnerCashFlowAudit != nil {
		t.Fatalf("legacy research should keep ownerCashFlowAudit nil")
	}
}

func TestNormalizeResearchKeepsDividendForecastAndComputesNetCash(t *testing.T) {
	research := normalizeResearch(ResearchImport{
		Symbol:   "000333.SZ",
		Name:     "美的集团",
		AsOf:     "2026-05-01",
		Currency: "CNY",
		Industry: "家电",
		Status:   "主策略候选",
		Action:   "等待A股股息率和DCF边际同时达标",
		Risk:     "原材料和海外需求波动",
		Valuation: Valuation{
			IntrinsicValue: floatPtr(100),
			MarginOfSafety: floatPtr(0.16),
		},
		Quality: Quality{
			TotalScore:       floatPtr(86),
			BusinessModel:    floatPtr(27),
			Moat:             floatPtr(22),
			Governance:       floatPtr(17),
			FinancialQuality: floatPtr(20),
		},
		Plan: PlanInput{
			Rank:       1,
			Priority:   "主策略复核",
			Advice:     "仅在股息率达标后新增",
			Discipline: "A股股息率≥6%，DCF边际≥15%",
		},
		Dividend: &Dividend{
			FiscalYear:         "FY2025",
			DividendPerShare:   floatPtr(5.0),
			ForecastFiscalYear: "FY2026E",
			ForecastPerShare:   floatPtr(6.2),
			ForecastYield:      floatPtr(0.062),
			Reliability:        "stable",
		},
		NetCash: &NetCashProfile{
			CashAndShortInvestments: floatPtr(100),
			InterestBearingDebt:     floatPtr(20),
			Haircut:                 floatPtr(0.7),
			ExCashPE:                floatPtr(9.5),
			ExCashPFCF:              floatPtr(8.8),
			FCFYield:                floatPtr(0.12),
			FCFPositiveYears:        intPtr(5),
		},
		OwnerCashFlowAudit: &OwnerCashFlowAudit{
			TenYearDemand:                  OwnerAuditItem{Status: " PASS ", Note: "十年需求稳定"},
			AssetDurability:                OwnerAuditItem{Status: "pass", Note: "品牌资产耐久"},
			MaintenanceCapexLight:          OwnerAuditItem{Status: "review", Note: "维持性资本开支待复核"},
			DividendFCFSupport:             OwnerAuditItem{Status: "pass", Note: "分红由FCF覆盖"},
			DividendReinvestmentEfficiency: OwnerAuditItem{Status: "review", Note: "估值仍需复核"},
			RoeRoicDurability:              OwnerAuditItem{Status: "pass", Note: "长期ROIC稳定"},
			ValuationSystemRisk:            OwnerAuditItem{Status: "pass", Note: "未见估值体系永久改变"},
		},
		Notes: "验证双策略导入字段。",
	})

	if _, err := validateResearch(research); err != nil {
		t.Fatalf("validate research with dual strategy fields: %v", err)
	}
	if research.Dividend == nil {
		t.Fatalf("expected dividend profile")
	}
	if research.Dividend.ForecastCurrency != "CNY" {
		t.Fatalf("forecast currency = %q, want CNY", research.Dividend.ForecastCurrency)
	}
	if research.NetCash == nil {
		t.Fatalf("expected net cash profile")
	}
	if research.NetCash.Currency != "CNY" {
		t.Fatalf("net cash currency = %q, want CNY", research.NetCash.Currency)
	}
	if research.NetCash.NetCash == nil || *research.NetCash.NetCash != 80 {
		t.Fatalf("net cash = %v, want 80", research.NetCash.NetCash)
	}
	if research.NetCash.AdjustedNetCash == nil || *research.NetCash.AdjustedNetCash != 56 {
		t.Fatalf("adjusted net cash = %v, want 56", research.NetCash.AdjustedNetCash)
	}
	if research.OwnerCashFlowAudit == nil {
		t.Fatalf("expected owner cash flow audit")
	}
	if research.OwnerCashFlowAudit.TenYearDemand.Status != "pass" {
		t.Fatalf("normalized owner audit status = %q, want pass", research.OwnerCashFlowAudit.TenYearDemand.Status)
	}
}

func TestValidateResearchRejectsInvalidNetCashHaircut(t *testing.T) {
	research := normalizeResearch(ResearchImport{
		Symbol:   "0999.HK",
		Name:     "测试公司",
		AsOf:     "2026-05-01",
		Currency: "HKD",
		Valuation: Valuation{
			MarginOfSafety: floatPtr(0.2),
		},
		NetCash: &NetCashProfile{
			Haircut: floatPtr(1.2),
		},
	})

	if _, err := validateResearch(research); err == nil {
		t.Fatalf("expected invalid netCash.haircut to be rejected")
	}
}

func TestValidateResearchAllowsZeroExCashMultiples(t *testing.T) {
	research := normalizeResearch(ResearchImport{
		Symbol:   "0999.HK",
		Name:     "测试公司",
		AsOf:     "2026-05-01",
		Currency: "HKD",
		Valuation: Valuation{
			MarginOfSafety: floatPtr(0.2),
		},
		NetCash: &NetCashProfile{
			ExCashPE:   floatPtr(0),
			ExCashPFCF: floatPtr(0),
		},
	})

	if _, err := validateResearch(research); err != nil {
		t.Fatalf("zero ex-cash multiples should be accepted: %v", err)
	}
}

func TestValidateResearchRejectsInvalidOwnerAuditStatus(t *testing.T) {
	research := normalizeResearch(ResearchImport{
		Symbol:   "0999.HK",
		Name:     "测试公司",
		AsOf:     "2026-05-01",
		Currency: "HKD",
		Valuation: Valuation{
			MarginOfSafety: floatPtr(0.2),
		},
		OwnerCashFlowAudit: &OwnerCashFlowAudit{
			TenYearDemand: OwnerAuditItem{Status: "unknown", Note: "非法状态"},
		},
	})

	if _, err := validateResearch(research); err == nil {
		t.Fatalf("expected invalid ownerCashFlowAudit status to be rejected")
	}
}

func TestApplyResearchKeepsExistingOwnerAuditWhenOmitted(t *testing.T) {
	existingAudit := &OwnerCashFlowAudit{
		TenYearDemand: OwnerAuditItem{Status: "pass", Note: "需求长期存在"},
	}
	holding := Holding{
		Symbol:             "000333.SZ",
		Name:               "美的集团",
		Currency:           "CNY",
		OwnerCashFlowAudit: existingAudit,
	}
	research := normalizeResearch(ResearchImport{
		Symbol:   "000333.SZ",
		Name:     "美的集团",
		AsOf:     "2026-05-01",
		Currency: "CNY",
		Valuation: Valuation{
			MarginOfSafety: floatPtr(0.2),
		},
	})

	applyHoldingResearch(&holding, research, "test")

	if holding.OwnerCashFlowAudit == nil {
		t.Fatalf("existing owner audit should be preserved")
	}
	if holding.OwnerCashFlowAudit.TenYearDemand.Note != "需求长期存在" {
		t.Fatalf("owner audit was overwritten when omitted")
	}
}
