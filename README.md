# 纸牌屋

股票投资组合管理台。

## 本地运行

```bash
go run .
```

启动后访问：

```text
http://127.0.0.1:8080
```

服务默认监听 `0.0.0.0:8080`，部署到云服务器后可通过服务器公网 IP 和 8080 端口访问。

如果 8080 已被占用，可指定端口：

```bash
PORT=8081 go run .
```

生产或云服务器建议把组合数据放到可持久化目录，并通过环境变量指定：

```bash
PORTFOLIO_DATA_DIR=/data/holds go run .
```

服务会从该目录读写 `portfolio.json`、`runtime/quotes.json` 和 `backups/`。如果目录里没有 `portfolio.json`，启动时会用仓库内置数据初始化一份。所有组合数据写入都会先写临时文件再原子替换，交易、持仓编辑、基金编辑和研究导入会自动备份旧文件。

## 更新行情

进入总览页时，网站会按冷却窗口在后台自动刷新股票行情；完整的股票行情、基金净值、股息数据和 ETF 追踪状态刷新由本地或云服务器定时任务调用 `/api/quotes/update` 完成。网站会优先拉取 Yahoo Finance 日线收盘价；如果云服务器 IP 被 Yahoo 限流，会自动切换到腾讯实时行情，必要时再尝试东方财富。基金里，场内 ETF 复用证券行情路径，场外公募基金拉取天天基金净值。

行情和基金净值优先写入运行时文件 `runtime/quotes.json`，不会污染 `portfolio.json` 里的静态研究档案。页面读取状态时会自动把 `portfolio.json` 和 runtime 行情合并，安全边际按 `(内在价值 - 最新价) / 内在价值` 实时重算。

ETF权益池由中证A500、红利低波、标普500和纳指100四类组成，各占权益池 25%，当前各自目标 35 万元；四只场外基金统一按每个交易日 250 元执行基础定投。

### 中证A500场内机会仓口径

- `022434` 保持每日 250 元场外定投；`159352` 只承担场内机会仓，二者持仓合并计入中证A500目标。
- 主触发只使用中证官方 `000510CNY010` 全收益指数相对历史高点的回撤，不使用价格指数、ETF K线或基金净值代替。档位为 `-7% / -12% / -18% / -25% / -35% / -45%`，对应本轮固定机会资金 `P₀` 的 `10% / 20% / 25% / 25% / 15% / 5%`。
- PE使用中证官网自 2024-09-02 起的官方滚动PE扩展窗口，股债利差按 `1/PE-中债10年期收益率` 计算。PE与利差缺失时估值系数按1处理，不阻断回撤主信号。
- 恐慌系数由RV20五年分位或近5个交易日全收益跌幅确认；成份股广度没有可靠免费官方序列时保持中性，不根据新闻或情绪加速。
- `159352` 估算净值由最近公布净值按中证A500全收益指数滚动到行情日。溢价和买卖价差都不高于0.15%才允许执行；溢价高于0.30%或价差高于0.20%时不追价。
- 本轮最长18个月，每3个月检查最低完成率；低于进度线的差额从尚未使用的 `P₀` 扣除并分4周补足。

### 红利低波场内建仓口径

- `008163` 每日定投 250 元；`515450` 承担场内机动建仓，现有 `563020` 仅计入红利低波持仓。
- 回撤信号使用 `515450` 单位净值与每份现金分红重建的分红再投资总回报，不使用场内收盘价、单位净值原始跌幅或累计净值直接代替。
- 估值确认使用 `515450` 官方历史申购赎回篮子重算：南方基金提供每日成分与数量，东方财富提供成分股历史收盘价、PB 和分红，中债提供 10 年期国债收益率。该值明确标记为场内篮子代理，不冒充标普官方指数值。
- 系统保存至少五年的周频篮子股息率与 PB 历史，计算股债利差分位、PB 分位及 `V=75%×股债利差分位+25%×PB便宜度`。有效成分或权重覆盖不足 95% 时拒绝产生场内信号。
- 股债利差分位或 PB 分位缺失、历史少于 230 个周频观测、数据超过 14 天未更新时，场外定投继续执行，场内分批计划保持等待，不沿用旧的基金分红代理值。

### 标普500场内机会仓口径

- `018738` 保持每日 250 元场外定投；`513650` 作为可配置的场内机会仓标的，两者持仓共同计入标普500目标 35 万元。
- 触发信号优先使用 S&P 500 Total Return（`SPTR`，机器行情代码 `^SP500TR`）相对历史高点的回撤；若 Yahoo 对服务器限流，则透明降级到 State Street 官方 SPY 价格与分红重建的美元总回报代理，页面会明确标记“SPY备援”，不会改用 `513650` K线。档位为 `-8% / -12% / -18% / -25% / -35%`，分别对应机会资金 `P` 的 `10% / 20% / 25% / 25% / 15%`，剩余 5% 为机动或时间仓。
- 未来 PE 与盈利收益率利差使用 History of Market 同一条 S&P 500 forward PE 历史；美国10年期收益率取美国财政部官方日线，统一转成十年周频分位。估值只调整档位金额，不单独触发买入。
- VIX 取 Cboe 官方历史数据，只决定分两笔、一次完成或提前半档，不能绕过 SPTR `-8%` 首档。免费同口径的三个月盈利预期修正没有确认时，系统关闭“提前半档”，但不影响基础定投和已触发档位。
- 人民币口径回撤使用 `SPTR×USD/CNY`；人民币回撤明显不足时，场内候选金额减半。
- `513650` 溢价使用最新场内价格和基金净值，并以 SPTR、USD/CNY、标普500期货变动估算实时净值：估算溢价不高于 0.5% 正常执行，0.5%—1% 减半，高于 1% 暂停。
- 任一核心执行指标缺失或过期时，`018738` 场外定投继续，`513650` 场内机会仓保持等待。

### 纳指100场内机会仓口径

- `021000` 保持每日 250 元场外定投；`159659` 作为可配置的场内机会仓标的，两者持仓共同计入纳指100目标仓位。
- 触发信号只使用 Nasdaq 官方 `XNDX` 全收益指数相对十年区间内历史高点的回撤，档位为 `-10% / -15% / -20% / -30% / -40%`，分别对应机会资金 `P` 的 `10% / 20% / 25% / 25% / 15%`，剩余 5% 为机动或时间仓。
- 未来 PE 与盈利收益率利差使用 History of Market 同一条 Nasdaq 100 forward PE 历史；美国10年期收益率取美国财政部官方日线，统一转成十年周频分位。估值只调整档位金额，不单独触发买入。
- VXN 取 Cboe 官方历史数据，只决定分两笔、一次完成或提前半档，不能绕过 XNDX `-10%` 首档。
- 人民币口径回撤使用 `XNDX×USD/CNY`；USD/CNY 历史取 Frankfurter/欧洲央行参考汇率。人民币回撤明显不足时，场内候选金额减半。
- `159659` 溢价使用东方财富最新场内价格和基金净值，并以 XNDX、USD/CNY、新浪纳指期货变动（Yahoo Finance 备援）估算实时净值：估算溢价不高于 0.5% 正常执行，0.5%—1% 减半，高于 1% 暂停。
- 任一核心执行指标缺失或过期时，`021000` 场外定投继续，`159659` 场内机会仓保持等待，不用 VXN、表面溢价或旧指标替代。

也可以用命令行更新股票行情和基金净值：

```bash
go run ./cmd/update-quotes
```

指定 runtime 行情文件：

```bash
go run ./cmd/update-quotes -quotes data/runtime/quotes.json
```

命令行工具同样支持 `PORTFOLIO_DATA_DIR` 作为默认数据目录，也可以继续用 `-data`、`-quotes` 显式指定文件。

只校验不写入：

```bash
go run ./cmd/update-quotes -dry-run
```

## 导入 ChatGPT 股票分析

先按 [docs/chatgpt-research-bridge.md](docs/chatgpt-research-bridge.md) 的格式让 ChatGPT 输出 JSON。

推荐在网站右上角点击“导入分析”，粘贴 JSON，先校验预览，再确认写入。确认导入会更新 `data/portfolio.json`，并自动备份旧文件到 `data/backups/`。

也可以保存到 `data/research/` 后用命令导入，例如：

```text
data/research/0700.HK.json
```

导入到网站数据：

```bash
go run ./cmd/import-research data/research/0700.HK.json
```

只校验不写入：

```bash
go run ./cmd/import-research -dry-run data/research/0700.HK.json
```

## 检查

```bash
go test ./...
node --check app.js
```

## 杀死程序

```bash
kill $(lsof -tiTCP:8080 -sTCP:LISTEN)
```

## Server scheduled updates

Use a systemd timer on the cloud server to call the backend update endpoint. See [docs/server-scheduled-update.md](docs/server-scheduled-update.md). The overview page may refresh stock quotes opportunistically, while scheduled tasks remain responsible for the full market-data, fund NAV, dividend, and ETF update flow.

## Local scheduled updates

Use Windows Task Scheduler to call the same backend update endpoint locally. See [docs/local-scheduled-update.md](docs/local-scheduled-update.md).
