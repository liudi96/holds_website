# Local scheduled market updates on Windows

The local task mirrors the cloud server update flow: it keeps the backend
available, calls `POST /api/quotes/update`, writes market data to
`data/runtime/quotes.json`, and lets the backend maintain daily P&L history in
`data/portfolio.json`.

## Script

```powershell
PowerShell -NoProfile -ExecutionPolicy Bypass -File D:\program\oneshot\holds_website\scripts\update-market-data.ps1
```

The script:

- checks `http://127.0.0.1:8080/api/health`;
- starts the local backend with `.tools\go\bin\go.exe run .` if needed;
- calls `http://127.0.0.1:8080/api/quotes/update`;
- writes logs under `.run\`, which is ignored by Git.

## Scheduled tasks

Use two Windows Task Scheduler jobs, matching the cloud schedule:

```powershell
schtasks /Create /TN "HoldsWebsiteMarketUpdate-0830" /SC DAILY /ST 08:30 /TR "powershell.exe -NoProfile -ExecutionPolicy Bypass -File \"D:\program\oneshot\holds_website\scripts\update-market-data.ps1\"" /F

schtasks /Create /TN "HoldsWebsiteMarketUpdate-2030" /SC DAILY /ST 20:30 /TR "powershell.exe -NoProfile -ExecutionPolicy Bypass -File \"D:\program\oneshot\holds_website\scripts\update-market-data.ps1\"" /F
```

## Verify

```powershell
schtasks /Query /TN "HoldsWebsiteMarketUpdate-0830"
schtasks /Query /TN "HoldsWebsiteMarketUpdate-2030"
PowerShell -NoProfile -ExecutionPolicy Bypass -File D:\program\oneshot\holds_website\scripts\update-market-data.ps1
```

The task updates local data only. It does not sync holdings from the cloud
server.
