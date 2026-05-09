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

## 更新行情

推荐在总览页点击“更新行情”。网站会优先拉取 Yahoo Finance 日线收盘价；如果云服务器 IP 被 Yahoo 限流，会自动切换到腾讯实时行情，必要时再尝试东方财富。更新会写入 `data/portfolio.json`，持仓和已拉取行情的候选股安全边际会按 `(内在价值 - 最新价) / 内在价值` 同步重算。

也可以用命令行更新：

```bash
go run ./cmd/update-quotes
```

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
