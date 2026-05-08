# holds_website

股票持仓管理台。

## 本地运行

```bash
go run .
```

启动后访问：

```text
http://127.0.0.1:8080
```

## 更新行情

拉取 Yahoo Finance 日线收盘价，并更新 `data/portfolio.json`：

```bash
go run ./cmd/update-quotes
```

只校验不写入：

```bash
go run ./cmd/update-quotes -dry-run
```

## 导入 ChatGPT 股票分析

先按 [docs/chatgpt-research-bridge.md](docs/chatgpt-research-bridge.md) 的格式让 ChatGPT 输出 JSON，保存到 `data/research/`，例如：

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
