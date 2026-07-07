#!/usr/bin/env bash
set -euo pipefail

base_url="${HOLDS_WEBSITE_URL:-http://127.0.0.1:8080}"
timeout_seconds="${HOLDS_UPDATE_TIMEOUT:-180}"

curl --fail --silent --show-error \
  --max-time "${timeout_seconds}" \
  --request POST \
  "${base_url%/}/api/quotes/update" >/dev/null

printf 'market data updated through %s at %s\n' "${base_url}" "$(date '+%Y-%m-%d %H:%M:%S')"
