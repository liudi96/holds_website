# ChatGPT Research Bridge

Use this bridge when ChatGPT performs the stock analysis and Codex imports the result into the website.

## Workflow

1. Ask ChatGPT to analyze one stock and output only the JSON schema below.
2. Open the website and click `导入分析`.
3. Paste the JSON, run preview validation, then confirm the import.

The website writes the confirmed result into `data/portfolio.json` and creates a backup under `data/backups/`.

Use `updateType: "fullReview"` for annual/full re-underwriting. Use `updateType: "eventUpdate"` for quarterly results, dividends, buybacks, regulatory events, earnings previews, or other incremental updates. If `updateType` is omitted, the importer treats the JSON as `fullReview`.

For institutional research, ChatGPT should first write a full Markdown research report, then output the import JSON. The full report should cover the old ZIP archive, latest facts, profit driver map, sensitivity analysis, scenario valuation, and operation discipline. The website JSON remains a compact summary for execution.

CLI import is still available for saved files:

```bash
go run ./cmd/import-research data/research/0700.HK.json
```

## ChatGPT Prompt

```text
Analyze the following stock for a dual-strategy value-investing portfolio.

Requirements:
- Output only valid JSON. No Markdown fences, no explanation outside JSON.
- Use the exact schema below.
- Do not add extra fields. Unknown fields are rejected by the importer.
- Use decimal ratios for percentages, for example 0.09 means 9%.
- Use null when a numeric field is unknown.
- Main strategy: self-selected blue chips, A-share comprehensive shareholder return yield >= 6% or H-share comprehensive shareholder return yield >= 8%, DCF margin >= 15%, no major risk, no low-confidence valuation, no clearly unsustainable dividend.
- The website converts `ownerCashFlowAudit` into a 100-point long-owner cash-flow score. Main strategy requires >= 75/100; `review` receives partial credit, so the audit is no longer an all-or-nothing hard gate.
- Side strategy: net-cash cigar butts, adjusted net cash after haircut, A-share ex-cash PE <= 10 or H-share ex-cash PE <= 8, with positive/free-cash-flow support.
- Existing holdings that fail the new thresholds should usually be marked as transition observation, not forced sale, unless major risk is present.
- `valuation.marginOfSafety` is the analysis-time estimate. For existing holdings, the website recalculates displayed safety margin from `intrinsicValue` and the latest close price.
- `currency` should match the listing: `.HK` uses `HKD`, `.SH`/`.SZ`/`.SS` use `CNY`.
- `quality.totalScore` should equal `businessModel + moat + governance + financialQuality`.
- Keep action, risk, notes, advice, and discipline concise but specific.
- asOf must be YYYY-MM-DD.

Institutional research steps before JSON:
- Read the old ZIP stock file first: cost, position size, old conclusion, target buy price, intrinsic value, safety margin, plan discipline, risk triggers, and recent research updates.
- Verify latest price, market cap, latest financial report, dividend, buyback, material news, industry data, and policy changes.
- Build a profit driver map: revenue side, cost side, policy side, competition side, and capital allocation side.
- Run sensitivity checks on variables that can change value: raw materials, volume, gross margin, expense ratio, tax/regulatory policy, or industry-specific factors.
- Provide bearish/base/bullish intrinsic value ranges before choosing the JSON intrinsicValue.
- State whether the research quality checklist is fully covered. If not, state the gaps in the report, not in the JSON.

Stock:
- Symbol:
- Name:
- Market/currency:
- User context:

JSON schema:
{
  "updateType": "fullReview",
  "symbol": "0700.HK",
  "name": "腾讯控股",
  "asOf": "2026-05-07",
  "currency": "HKD",
  "industry": "互联网平台/游戏/广告/金融科技",
  "status": "未达标（安全边际<15%）",
  "action": "继续持有；新资金暂不追买，等待安全边际达标后再分批",
  "risk": "政策、地缘、AI投入回报周期和广告/游戏周期波动需折价",
  "valuationConfidence": "high",
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
    "priority": "观察/低优先级",
    "advice": "等待安全边际达标后再分批，未达标不追买",
    "discipline": "优秀资产要求≥15%安全边际；未达标不追买"
  },
  "dividend": {
    "fiscalYear": "FY2025",
    "dividendPerShare": 4.5,
    "dividendCurrency": "HKD",
    "payoutRatio": 0.16,
    "reliability": "stable",
    "forecastFiscalYear": "FY2026E",
    "forecastPerShare": 5.2,
    "forecastCurrency": "HKD",
    "forecastYield": 0.083
  },
  "netCash": {
    "cashAndShortInvestments": 320000000000,
    "interestBearingDebt": 120000000000,
    "netCash": 200000000000,
    "currency": "HKD",
    "haircut": 0.7,
    "haircutReason": "平台现金流稳定但需保留监管和再投资折扣",
    "adjustedNetCash": 140000000000,
    "exCashPe": 13.5,
    "exCashPfcf": 14.2,
    "fcfYield": 0.065,
    "fcfPositiveYears": 5,
    "note": "净现金、FCF 和估值口径使用 FY2025 年报与当前市值。"
  },
  "ownerCashFlowAudit": {
    "tenYearDemand": { "status": "pass", "note": "核心产品/服务十年后仍有稳定需求。" },
    "assetDurability": { "status": "pass", "note": "品牌、资源或网络资产不易折旧。" },
    "maintenanceCapexLight": { "status": "review", "note": "需继续核实维持性资本开支。" },
    "dividendFcfSupport": { "status": "pass", "note": "分红由真实自由现金流覆盖。" },
    "dividendReinvestmentEfficiency": { "status": "review", "note": "当前估值对分红再投入效率一般。" },
    "roeRoicDurability": { "status": "pass", "note": "ROE/ROIC 有长期维持基础。" },
    "valuationSystemRisk": { "status": "pass", "note": "暂未发现行业估值体系永久改变。" }
  },
  "killCriteria": [
    "若核心业务增长和自由现金流连续两个季度明显恶化，应重新评估内在价值",
    "若监管、治理或财报可信度出现重大风险，应暂停新增资金"
  ],
  "notes": "Summarize the key financial facts, assumptions, and source period here."
}
```

Event update schema:

```json
{
  "updateType": "eventUpdate",
  "symbol": "0700.HK",
  "name": "腾讯控股",
  "asOf": "2026-05-15",
  "event": {
    "type": "earnings",
    "title": "2026Q1 财报更新",
    "date": "2026-05-15",
    "source": "公司公告",
    "summary": "收入和自由现金流好于原假设，AI capex 继续上升。"
  },
  "impact": {
    "thesisChange": "minor",
    "valuationChange": "raise",
    "riskChange": "unchanged",
    "actionChange": "unchanged"
  },
  "updates": {
    "valuation": {
      "intrinsicValue": 550,
      "fairValueRange": "HK$500-610",
      "marginOfSafety": 0.15
    },
    "risk": "AI资本开支仍需跟踪，但短期现金流韧性增强",
    "notesAppend": "2026Q1 证实广告和游戏恢复，暂不改变买入纪律。"
  }
}
```

## Import Rules

- If `symbol` matches an existing holding, the holding analysis fields are updated.
- If `symbol` matches an existing candidate, the candidate fields are updated.
- If `symbol` is new, it is added to the candidate pool.
- `fullReview` replaces the stock's current research fields. `eventUpdate` only applies fields explicitly present under `updates`; missing fields keep their previous website values.
- `eventUpdate` records a durable `researchUpdates` timeline entry so the next ChatGPT export can build on prior event/earnings updates.
- Full profit driver maps, sensitivity tables, industry policy analysis, competition analysis, and scenario valuation belong in the Markdown research report. JSON only stores the compact decision summary; use `notes` or `notesAppend` for the most important changes.
- `plan` is upserted by the top-level `symbol` when possible, then by stock name for old data.
- Do not include `symbol` inside `plan`; the importer derives Plan identity from the top-level `symbol`.
- `plan.rank` may be approximate. The importer normalizes Plan ranks into a unique sequence after import.
- `valuation.intrinsicValue` is the core DCF estimate from ChatGPT. The main strategy requires displayed DCF margin >= 15%.
- `valuationConfidence` and `killCriteria` are optional. If omitted, the website derives valuation confidence from quality score and risk text, and derives the detail-page bear case from existing risk/status fields.
- The site computes the first-buy price as `intrinsicValue * 75%`, watch price as `firstBuyPrice * 105%`, and aggressive buy price as `firstBuyPrice * 90%`.
- Dividend data is fetched by the quote update flow when the data source provides it, but research may also provide `dividend.forecastFiscalYear`, `forecastPerShare`, `forecastCurrency`, and `forecastYield` for reference. Main-strategy comprehensive shareholder return passes only when latest full fiscal-year comprehensive shareholder return yield reaches A-share 6% / H-share 8%.
- `ownerCashFlowAudit` is required for a main-strategy buy. Each item uses `status: pass|review|fail` plus `note`; the website scores the seven items with weights: ten-year demand 18, asset durability 14, light reinvestment 12, dividend FCF support 18, reinvestment efficiency 12, ROE/ROIC durability 14, valuation-system risk 12. Missing fields default to review only when some audit evidence exists.
- If `valuationSystemRisk.status` is `fail`, the website treats the stock as risk exclusion. Other review/fail items reduce the long-owner score but no longer automatically block main-strategy buying when the total score still reaches 75/100.
- Dividend yield is calculated as latest full fiscal-year cash dividend total divided by company market capitalization; comprehensive shareholder return is calculated as cash dividends plus buybacks divided by market capitalization.
- `dividend.reliability` is optional. If omitted, the website derives `stable/review/risk` from dividend data completeness, valuation confidence, and major risk text.
- `netCash.cashAndShortInvestments`, `interestBearingDebt`, `netCash`, `currency`, `haircut`, `haircutReason`, `adjustedNetCash`, `exCashPe`, `exCashPfcf`, `fcfYield`, `shareholderFcf`, `shareholderFcfCurrency`, `shareholderFcfBasis`, `consolidatedFcf`, `minorityFcfAdjustment`, and `fcfPositiveYears` are optional but should be supplied for cigar-butt candidates. For companies with material minority interests, `shareholderFcf` should be the ordinary-shareholder free cash flow after minority-interest leakage.
- Net-cash haircut convention: stable dividend 100%, normal 70%, weak/cyclical 40%, major risk 0%. If `haircut` is omitted, the website estimates it from dividend reliability and risk text.
- The website computes dual-strategy grouping locally: main strategy, side-strategy cigar butt, transition observation, or risk exclusion.
- Future optional analysis fields may be useful but are not required: `circleOfCompetence`, `ownerEarnings`, `roeHistory`, `debtRatio`, `dividendCoverage`, and capital allocation notes. Missing fields should not block import.
- The website preview validates the JSON before writing. Confirmed imports update `data/portfolio.json` and first create a backup under `data/backups/`.
- Holding safety margin is calculated as `(intrinsicValue - currentPrice) / intrinsicValue`. Candidate-pool stocks use the same formula after the overview `更新行情` button has fetched quote data into the local runtime quote file; otherwise they continue to use the imported `valuation.marginOfSafety`.
- Quote fields such as `currentPrice`, `previousClose`, and close dates are owned by `data/runtime/quotes.json`, which is generated by the overview `更新行情` button or `cmd/update-quotes`. Research import must not write these fields into `data/portfolio.json`.
