#!/usr/bin/env bash
# Render healthcheck ワークフローの実行を gh で監視する。
# 使い方（リポジトリルート）: ./scripts/monitor-render-healthcheck.sh
# オプション: INTERVAL_SEC=120 DURATION_MIN=35 ./scripts/monitor-render-healthcheck.sh
set -euo pipefail

INTERVAL_SEC="${INTERVAL_SEC:-60}"
DURATION_MIN="${DURATION_MIN:-35}"
# 10 分 cron の次の発火を一度は拾えるよう、デフォルトは 35 分監視

REPO="$(gh repo view --json nameWithOwner -q .nameWithOwner)"
WF_PATH=".github/workflows/render-healthcheck.yml"
WF_ID="$(gh api "repos/${REPO}/actions/workflows" -q ".workflows[] | select(.path==\"${WF_PATH}\") | .id" | head -1)"
if [[ -z "${WF_ID}" ]]; then
  echo "エラー: ワークフロー ${WF_PATH} が見つかりません" >&2
  exit 1
fi

echo "repo=${REPO} workflow_id=${WF_ID} path=${WF_PATH}"
echo "間隔 ${INTERVAL_SEC} 秒、最大 ${DURATION_MIN} 分間ポーリングします（Ctrl+C で中断）。"
echo "---"

end=$((SECONDS + DURATION_MIN * 60))
last_seen=""

while (( SECONDS < end )); do
  # 直近 5 件（schedule / workflow_dispatch の両方）
  mapfile -t lines < <(gh api "repos/${REPO}/actions/workflows/${WF_ID}/runs?per_page=5" \
    --jq '.workflow_runs[] | "\(.created_at)\t\(.event)\t\(.status)\t\(.conclusion // "-")\t\(.html_url)"')
  if [[ ${#lines[@]} -eq 0 ]]; then
    echo "$(date -Iseconds) 実行履歴なし"
  else
    blob="$(printf '%s\n' "${lines[@]}")"
    if [[ "${blob}" != "${last_seen}" ]]; then
      echo "$(date -Iseconds) 直近の実行:"
      printf '%s\n' "${lines[@]}" | while IFS= read -r row; do echo "  ${row}"; done
      last_seen="${blob}"
    else
      echo "$(date -Iseconds) 変化なし（最新は上記と同じ）"
    fi
  fi
  sleep "${INTERVAL_SEC}"
done

echo "---"
echo "監視終了。GitHub の schedule は UTC で、初回や間欠時は数分ずれることがあります。"
echo "一覧: gh run list --workflow=${WF_PATH} --limit 20"
