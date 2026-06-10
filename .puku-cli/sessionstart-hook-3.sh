#!/usr/bin/env bash
# sessionstart-hook-3.sh
# Runs on SessionStart — restores session context for npbs_project

set -e

# Derive project root from script location (.claude/ is at project root)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SESSION_DIR="$SCRIPT_DIR/sessions/c7c59082-733b-448a-968f-e80390c0e1a5"
PLAN_FILE="$SESSION_DIR/ethereal-sparking-wind.md"

if [ ! -f "$PLAN_FILE" ]; then
  exit 0
fi

# Detect if plan file has been updated since last run
LAST_RUN_FILE="/tmp/npbs_session_last_run"
CURRENT_MD5=$(md5sum "$PLAN_FILE" 2>/dev/null | cut -d' ' -f1)
PREVIOUS_MD5=$(cat "$LAST_RUN_FILE" 2>/dev/null || echo "")

if [ "$CURRENT_MD5" != "$PREVIOUS_MD5" ]; then
  echo ""
  echo "=== Session Context Restored ==="
  echo "Project: Network Performance Benchmarking Suite"
  echo "Plan: $PLAN_FILE"
  echo ""

  # Extract key decisions from plan
  if grep -q "Hello world first" "$PLAN_FILE"; then
    echo "- Approach: Hello world first (single mode end-to-end, then 6-mode parallel)"
    echo "- Statistics: Median + Stddev"
    echo "- Report: Mode | Throughput (Gbps) | Latency (μs) | Jitter (μs)"
    echo ""
    echo "Current phase: Phase 1 — Hello World"
  fi

  echo "=============================="
  echo ""

  # Persist md5 for next run
  echo "$CURRENT_MD5" > "$LAST_RUN_FILE"
fi

exit 0