#!/usr/bin/env bash
# Post an AI code review on a pull request using the OpenRouter API.
set -euo pipefail

if [[ -z "${OPENROUTER_API_KEY:-}" ]]; then
  echo "OPENROUTER_API_KEY not configured; skipping OpenRouter review."
  exit 0
fi

if [[ -z "${PR_NUMBER:-}" ]]; then
  echo "PR_NUMBER is required."
  exit 1
fi

MODEL="${OPENROUTER_MODEL:-anthropic/claude-3.5-haiku}"
REPO="${GITHUB_REPOSITORY}"
MAX_DIFF_CHARS="${OPENROUTER_MAX_DIFF_CHARS:-120000}"

pr_json="$(gh pr view "$PR_NUMBER" --repo "$REPO" --json title,body,author,baseRefName,headRefName,files)"
title="$(jq -r '.title' <<<"$pr_json")"
body="$(jq -r '.body // ""' <<<"$pr_json")"
author="$(jq -r '.author.login' <<<"$pr_json")"
base="$(jq -r '.baseRefName' <<<"$pr_json")"
head="$(jq -r '.headRefName' <<<"$pr_json")"
files="$(jq -r '[.files[].path] | join(", ")' <<<"$pr_json")"

diff="$(gh pr diff "$PR_NUMBER" --repo "$REPO" || true)"
if [[ -z "$diff" ]]; then
  echo "No diff found for PR #$PR_NUMBER."
  exit 0
fi

if (( ${#diff} > MAX_DIFF_CHARS )); then
  diff="${diff:0:MAX_DIFF_CHARS}

... [diff truncated at ${MAX_DIFF_CHARS} characters]"
fi

read -r -d '' system_prompt <<'EOF' || true
You are a senior Go reviewer for Tiny Go MCP Server (tinymcp), a thin helper on modelcontextprotocol/go-sdk.

Review the pull request diff. Focus on:
- Registration must not panic; return errors for nil server/handler at registration time
- MCP tool descriptions: when to use, when not to, sibling tools
- Handlers: context.Context first; success text via tinymcp.TextResult
- HTTP/SSE deploy: auth, localhost binding, cross-origin protection for public endpoints
- Race-safe tests for changed Go code; godoc on exports
- Dual imports (tinymcp + mcp) are intentional — not a design flaw

Output markdown only:
1. "## Review summary" — 2-3 sentences
2. "## Findings" — bullet list; prefix each with severity [critical|high|medium|low|info]
3. "## Suggested follow-ups" — optional checklist

Skip nitpicks on markdown-only or version-bump-only changes unless inconsistent.
Do not approve or request changes; provide review feedback only.
EOF

user_prompt="Repository: ${REPO}
PR #${PR_NUMBER}: ${title}
Author: ${author}
Base: ${base} <- Head: ${head}
Changed files: ${files}

PR description:
${body}

Additional reviewer notes:
${ADDITIONAL_CONTEXT:-none}

Diff:
${diff}"

request_json="$(jq -n \
  --arg model "$MODEL" \
  --arg system "$system_prompt" \
  --arg user "$user_prompt" \
  '{
    model: $model,
    temperature: 0.2,
    max_tokens: 4096,
    messages: [
      {role: "system", content: $system},
      {role: "user", content: $user}
    ]
  }')"

response_file="$(mktemp)"
http_code="$(curl -sS -o "$response_file" -w '%{http_code}' \
  https://openrouter.ai/api/v1/chat/completions \
  -H "Authorization: Bearer ${OPENROUTER_API_KEY}" \
  -H "Content-Type: application/json" \
  -H "HTTP-Referer: https://github.com/${REPO}" \
  -H "X-Title: tiny-go-mcp-server PR Review" \
  -d "$request_json")"

if [[ "$http_code" != "200" ]]; then
  echo "OpenRouter API failed with HTTP ${http_code}:"
  cat "$response_file"
  exit 1
fi

review_body="$(jq -r '.choices[0].message.content // empty' "$response_file")"
if [[ -z "$review_body" ]]; then
  echo "OpenRouter returned an empty review:"
  cat "$response_file"
  exit 1
fi

comment_body="$(cat <<EOF
## OpenRouter PR review

Model: \`${MODEL}\`

${review_body}

---
*Automated review via [OpenRouter](https://openrouter.ai). Comment \`/review\` on this PR to re-run.*
EOF
)"

gh pr comment "$PR_NUMBER" --repo "$REPO" --body "$comment_body"
