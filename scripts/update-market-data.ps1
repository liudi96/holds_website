param(
  [string]$BaseUrl = $(if ($env:HOLDS_WEBSITE_URL) { $env:HOLDS_WEBSITE_URL } else { "http://127.0.0.1:8080" }),
  [int]$TimeoutSeconds = $(if ($env:HOLDS_UPDATE_TIMEOUT) { [int]$env:HOLDS_UPDATE_TIMEOUT } else { 180 }),
  [switch]$NoStart
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Resolve-Path (Join-Path $scriptDir "..")
$runDir = Join-Path $repoRoot ".run"
New-Item -ItemType Directory -Force -Path $runDir | Out-Null

$logPath = Join-Path $runDir "local-market-update.log"
$appOutPath = Join-Path $runDir "local-app-scheduled.out.log"
$appErrPath = Join-Path $runDir "local-app-scheduled.err.log"
$base = $BaseUrl.TrimEnd("/")

function Write-RunLog {
  param([string]$Message)
  $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
  Add-Content -Path $logPath -Value "[$timestamp] $Message" -Encoding UTF8
}

function Test-Backend {
  try {
    Invoke-RestMethod -Uri "$base/api/health" -Method Get -TimeoutSec 5 | Out-Null
    return $true
  } catch {
    return $false
  }
}

function Start-LocalBackend {
  $goExe = Join-Path $repoRoot ".tools\go\bin\go.exe"
  if (!(Test-Path $goExe)) {
    $goExe = "go"
  }

  Write-RunLog "Starting local backend with $goExe"
  Start-Process `
    -FilePath $goExe `
    -ArgumentList @("run", ".") `
    -WorkingDirectory $repoRoot `
    -RedirectStandardOutput $appOutPath `
    -RedirectStandardError $appErrPath `
    -WindowStyle Hidden | Out-Null
}

Write-RunLog "Market update started for $base"

if (!(Test-Backend)) {
  if ($NoStart) {
    throw "Local backend is not available at $base"
  }
  Start-LocalBackend
}

$deadline = (Get-Date).AddSeconds($TimeoutSeconds)
while (!(Test-Backend)) {
  if ((Get-Date) -gt $deadline) {
    throw "Timed out waiting for local backend at $base"
  }
  Start-Sleep -Seconds 2
}

$response = Invoke-RestMethod `
  -Uri "$base/api/quotes/update" `
  -Method Post `
  -TimeoutSec $TimeoutSeconds

Write-RunLog "Market update finished: updated=$($response.updated) skipped=$($response.skipped.Count)"
