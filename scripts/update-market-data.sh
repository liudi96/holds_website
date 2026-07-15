#!/usr/bin/env bash
set -euo pipefail

base_url="${HOLDS_WEBSITE_URL:-http://127.0.0.1:8080}"
timeout_seconds="${HOLDS_UPDATE_TIMEOUT:-180}"
response_file="$(mktemp)"
trap 'rm -f "${response_file}"' EXIT

curl --fail --silent --show-error \
  --max-time "${timeout_seconds}" \
  --request POST \
  "${base_url%/}/api/quotes/update" >"${response_file}"

python3 - "${response_file}" "${base_url}" <<'PY'
import json
import sys
from datetime import datetime

path, base_url = sys.argv[1:]
with open(path, encoding="utf-8") as handle:
    response = json.load(handle)

updated = int(response.get("updated") or 0)
skipped = response.get("skipped") or []
print(
    f"market data update finished through {base_url} at "
    f"{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}: "
    f"updated={updated} skipped={len(skipped)}"
)
for item in skipped:
    error = " ".join(str(item.get("error") or "unknown error").split())[:240]
    print(
        "market data partial failure: "
        f"type={item.get('type') or '-'} symbol={item.get('symbol') or '-'} "
        f"name={item.get('name') or '-'} error={error}"
    )

if updated == 0 and skipped:
    raise SystemExit(2)
PY
