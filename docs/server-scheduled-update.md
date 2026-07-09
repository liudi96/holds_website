# 云服务器定时更新行情

目标：页面只读取已经保存好的数据，不在进入总览页时触发行情更新。云服务器负责定时调用后端更新接口，更新股票行情、基金净值、ETF 追踪和从 2026-07-07 开始的日盈亏记录。

## 安装 systemd timer

仓库路径按当前云服务器约定使用：

```bash
cd /root/oneshot/holds_website
sudo cp deploy/holds-market-update.service /etc/systemd/system/
sudo cp deploy/holds-market-update.timer /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now holds-market-update.timer
```

默认每天执行两次：

- 08:30：早盘前刷新一次行情、基金净值和 ETF 追踪状态。
- 20:30：晚间补一次收盘后行情、基金净值和 ETF 追踪状态。

如果只想每天 20:30 执行一次，删除 `deploy/holds-market-update.timer` 里的 `OnCalendar=*-*-* 08:30:00` 后重新复制并 `daemon-reload`。

## 手动验证

```bash
cd /root/oneshot/holds_website
bash scripts/update-market-data.sh
systemctl list-timers --all | grep holds-market-update
journalctl -u holds-market-update.service -n 50 --no-pager
```

## 行为说明

- 脚本调用 `http://127.0.0.1:8080/api/quotes/update`，所以网站后端服务必须正在运行。
- 后端会把 runtime 行情写入 `data/runtime/quotes.json`，把日盈亏历史写入 `data/portfolio.json`。
- 日盈亏记录只从 `2026-07-07` 开始，不回填更早日期。
- 如果某几天没有打开网站，下一次定时更新会按日线补齐缺失交易日。
- 页面进入总览页不会再主动更新行情，只会读取服务端已保存的数据。
