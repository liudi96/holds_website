# ChatGPT Research Bridge

Use this bridge when ChatGPT performs the stock analysis and Codex imports the result into the website.

## Workflow

1. Ask ChatGPT to analyze one stock and output only the JSON schema below.
2. Save the JSON under `data/research/<SYMBOL>.json`, for example `data/research/0700.HK.json`.
3. Import it:

```bash
go run ./cmd/import-research data/research/0700.HK.json
```

4. Refresh the website.

## ChatGPT Prompt

```text
Analyze the following stock for a long-term value-investing portfolio.

Requirements:
- Output only valid JSON. No Markdown fences, no explanation outside JSON.
- Use the exact schema below.
- Use decimal ratios for percentages, for example 0.09 means 9%.
- Use null when a numeric field is unknown.
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
  "action": "继续持有；新资金暂不追买，等待目标价附近再分批",
  "risk": "政策、地缘、AI投入回报周期和广告/游戏周期波动需折价",
  "valuation": {
    "intrinsicValue": 508,
    "fairValueRange": "HK$480-560",
    "targetBuyPrice": 432,
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
    "advice": "等待≤HK$432，HK$400-430可分批",
    "discipline": "优秀资产要求≥15%安全边际；未达标不追买"
  },
  "notes": "Summarize the key financial facts, assumptions, and source period here."
}
```

## Import Rules

- If `symbol` matches an existing holding, the holding analysis fields are updated.
- If `symbol` matches an existing candidate, the candidate fields are updated.
- If `symbol` is new, it is added to the candidate pool.
- `plan` is upserted by stock name.
- Quote fields such as `currentPrice`, `previousClose`, and close dates are owned by `cmd/update-quotes`, not this import.
