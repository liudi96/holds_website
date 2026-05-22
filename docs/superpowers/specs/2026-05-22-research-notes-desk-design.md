# Research Notes Desk Design

## Goal

改造研究页，使它成为以个人心得为主轴的研究界面。首页优先支持快速记录人的判断，个股和行业作为可选关联对象与研究上下文，而不是第一屏主角。

## Confirmed Decisions

- 首页首要任务：方便记录和复盘个人心得。
- 最小记录单元：快速随手记。
- 关联规则：心得默认独立保存，可选关联个股和行业。
- 改造路径：心得优先工作台。

## Product Shape

研究页首页改为三块：

1. 快速输入区：顶部输入框，正文必填，其余字段可选。
2. 心得流：按时间倒序展示心得，支持筛选全部、未整理、有个股、有行业、已转原则。
3. 研究上下文：右侧展示关联对象入口、待整理提示、导入分析和更新行业数据动作。

现有个股详情页、行业详情页、导入分析、行业数据更新能力保留。它们在研究页上降为上下文和工具，不再主导第一屏。

## Data Model

在 `AppState` 中新增 `researchNotes`：

```go
type ResearchNote struct {
    ID                int64    `json:"id"`
    Date              string   `json:"date"`
    UpdatedAt         string   `json:"updatedAt,omitempty"`
    Content           string   `json:"content"`
    Tags              []string `json:"tags,omitempty"`
    Mood              string   `json:"mood,omitempty"`
    Confidence        string   `json:"confidence,omitempty"`
    Status            string   `json:"status"`
    LinkedSymbols     []string `json:"linkedSymbols,omitempty"`
    LinkedIndustryIDs []string `json:"linkedIndustryIds,omitempty"`
}
```

Status first version:

- `unfiled`: 未整理。
- `organized`: 已整理。
- `principle`: 已转原则。

不复用 `decisionLogs` 或 `researchUpdates`。原因：`decisionLogs` 是系统/交易/行情事件，`researchUpdates` 是导入分析事件，心得是人的判断片段，可以没有任何个股或行业。

## API

新增 API：

- `POST /api/research-notes`: 创建心得。
- `PUT /api/research-notes/{id}`: 编辑心得内容、标签、关联对象和状态。
- `DELETE /api/research-notes/{id}`: 删除误记。

保存沿用现有 `saveStateWithBackup` 路径，继续写入 `portfolio.json`，保留备份和原子替换机制。

Validation:

- 创建和编辑时 `content` 去空白后不能为空。
- `status` 只能是允许值；空值归一化为 `unfiled`。
- `linkedSymbols` 使用现有 `normalizeSymbol` 去重。
- `linkedIndustryIds` 使用现有行业 ID 归一化；无法匹配时允许保存，但前端显示为未匹配，避免阻断随手记。

## Frontend

研究页仍使用现有 `data-page="industry"` 入口，页面标题仍为“研究台”，但内容改为心得优先。

主要组件：

- Quick note composer：正文输入、可选个股/行业/标签/置信度/状态。
- Research note feed：心得卡片列表。
- Feed filters：全部、未整理、有个股、有行业、已转原则。
- Context rail：关联研究入口、待整理提示、导入分析、更新行业数据。

交互规则：

- 只填正文即可保存。
- 保存成功后，心得出现在列表顶部，输入区清空。
- 保存失败时保留输入并显示错误。
- 空内容在前端直接提示，不发请求。
- 编辑和删除在卡片内完成。

## Error Handling

- Backend load/save failure returns structured API error and does not mutate in-memory state.
- Frontend request failure keeps current draft and shows status text.
- Unknown linked symbol or industry does not block save; note remains independent and can be linked later.

## Testing

Go tests:

- Create research note trims content and assigns ID/date/status.
- Empty content is rejected.
- Update can change content, tags, links, and status.
- Invalid status is rejected; empty status is normalized to `unfiled`.
- Delete removes the note.
- Notes sort newest first after create.

Frontend/static checks:

- `node --check app.js`

Browser verification:

- Start the app.
- Open research page.
- Save a note with only content.
- Refresh state/page and confirm the note remains.
- Add optional stock/industry association and confirm linked labels render.

## Out Of Scope

- Rich text editor.
- Attachments.
- Automatic principle generation.
- Changing ChatGPT research import JSON.
- Reworking stock detail page.
- Reworking industry JSON format.
