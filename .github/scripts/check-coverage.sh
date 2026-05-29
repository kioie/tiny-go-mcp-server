#!/usr/bin/env bash
# Enforces minimum statement coverage for core packages (tinymcp + reference server).
set -euo pipefail

MIN="${1:-70}"
shift || true
PKGS=("$@")
if [ ${#PKGS[@]} -eq 0 ]; then
  PKGS=("./tinymcp/..." "./cmd/tiny-go-mcp/...")
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT}"

go test -race -coverprofile=coverage.out -covermode=atomic "${PKGS[@]}"
PCT="$(go tool cover -func=coverage.out | awk '/^total:/ {gsub(/%/,"",$3); print $3}')"
echo "Coverage: ${PCT}% (minimum ${MIN}%)"

awk -v pct="${PCT}" -v min="${MIN}" 'BEGIN {
  if (pct + 0 < min + 0) {
    printf("Coverage %.1f%% is below minimum %s%%\n", pct, min) > "/dev/stderr"
    exit 1
  }
}'
