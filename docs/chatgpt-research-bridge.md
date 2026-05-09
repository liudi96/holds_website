# ChatGPT Research Bridge

Use this bridge when ChatGPT performs the stock analysis and Codex imports the result into the website.

## Workflow

1. Ask ChatGPT to analyze one stock and output only the JSON schema below.
2. Open the website and click `导入分析`.
3. Paste the JSON, run preview validation, then confirm the import.

The website writes the confirmed result into `data/portfolio.json` and creates a backup under `data/backups/`.

CLI import is still available for saved files:

```bash
go run ./cmd/import-research data/research/0700.HK.json
```

## ChatGPT Prompt

```text
Analyze the following stock for a long-term value-investing portfolio.

Requirements:
- Output only valid JSON. No Markdown fences, no explanation outside JSON.
- Use the exact schema below.
- Do not add extra fields. Unknown fields are rejected by the importer.
- Use decimal ratios for percentages, for example 0.09 means 9%.
- Use null when a numeric field is unknown.
- `valuation.marginOfSafety` is the analysis-time estimate. For existing holdings, the website recalculates displayed safety margin from `intrinsicValue` and the latest close price.
- `currency` should match the listing: `.HK` uses `HKD`, `.SH`/`.SZ`/`.SS` use `CNY`.
- `quality.totalScore` should equal `businessModel + moat + governance + financialQuality`.
- Keep action, risk, notes, advice, and discipline concise but specific.
- asOf must be YYYY-MM-DD.

Stock:
- Symbol:
- Name:
- Market/currency:
- User context:

JSON schema:
{
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
  "killCriteria": [
    "若核心业务增长和自由现金流连续两个季度明显恶化，应重新评估内在价值",
    "若监管、治理或财报可信度出现重大风险，应暂停新增资金"
  ],
  "notes": "Summarize the key financial facts, assumptions, and source period here."
}
```

## Import Rules

- If `symbol` matches an existing holding, the holding analysis fields are updated.
- If `symbol` matches an existing candidate, the candidate fields are updated.
- If `symbol` is new, it is added to the candidate pool.
- `plan` is upserted by the top-level `symbol` when possible, then by stock name for old data.
- Do not include `symbol` inside `plan`; the importer derives Plan identity from the top-level `symbol`.
- `plan.rank` may be approximate. The importer normalizes Plan ranks into a unique sequence after import.
- `valuation.intrinsicValue` is the core estimate from ChatGPT. `targetBuyPrice`, `priceLevels`, `dividend`, `dividendYield`, and `estimatedAnnualCash` do not need to be provided.
- `valuationConfidence` and `killCriteria` are optional. If omitted, the website derives valuation confidence from quality score and risk text, and derives the detail-page bear case from existing risk/status fields.
- The site computes the first-buy price as `intrinsicValue * 75%`, watch price as `firstBuyPrice * 105%`, and aggressive buy price as `firstBuyPrice * 90%`.
- Dividend data is fetched by the quote update flow when the data source provides it. Dividend yield is calculated as latest full fiscal-year cash dividend total divided by company market capitalization; comprehensive shareholder return is calculated as cash dividends plus buybacks divided by market capitalization.
- `dividend.reliability` is optional. If omitted, the website derives `stable/review/risk` from dividend data completeness, valuation confidence, and major risk text.
- The investment committee views are computed locally from the same imported data:
  - Howard Marks focuses on risk compensation, safety margin, valuation confidence, major risk words, and position concentration.
  - Benjamin Graham focuses on safety margin, financial quality, dividend reliability, and whether the valuation is defensive enough.
  - Warren Buffett focuses on quality score, business model, moat, governance, financial quality, valuation confidence, and long-term risk.
- Future optional analysis fields may be useful but are not required: `circleOfCompetence`, `ownerEarnings`, `roeHistory`, `debtRatio`, `dividendCoverage`, and capital allocation notes. Missing fields should not block import.
- The website preview validates the JSON before writing. Confirmed imports update `data/portfolio.json` and first create a backup under `data/backups/`.
- Holding safety margin is calculated as `(intrinsicValue - currentPrice) / intrinsicValue`. Candidate-pool stocks use the same formula after the overview `更新行情` button has fetched quote data; otherwise they continue to use the imported `valuation.marginOfSafety`.
- Quote fields such as `currentPrice`, `previousClose`, and close dates are owned by the overview `更新行情` button or `cmd/update-quotes`, not this import.
