#!/usr/bin/env bash
set -euo pipefail

ADDR="127.0.0.1:18081"
BASE_URL="http://${ADDR}"
WORKSPACE="$(mktemp -d)"
SERVER_PID=""

cleanup() {
  if [[ -n "${SERVER_PID}" ]] && kill -0 "${SERVER_PID}" 2>/dev/null; then
    kill "${SERVER_PID}" 2>/dev/null || true
    wait "${SERVER_PID}" 2>/dev/null || true
  fi
  rm -rf "${WORKSPACE}"
}
trap cleanup EXIT

start_server() {
  TODOOPEN_SERVER_ADDR="${ADDR}" TODOOPEN_WORKSPACE_ROOT="${WORKSPACE}" \
    go run ./cmd/todoopen-server >/dev/null 2>&1 &
  SERVER_PID=$!

  for _ in {1..50}; do
    if curl -fsS "${BASE_URL}/healthz" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.1
  done

  echo "server did not become healthy at ${BASE_URL}" >&2
  return 1
}

stop_server() {
  if [[ -n "${SERVER_PID}" ]] && kill -0 "${SERVER_PID}" 2>/dev/null; then
    kill "${SERVER_PID}" 2>/dev/null || true
    wait "${SERVER_PID}" 2>/dev/null || true
  fi
  SERVER_PID=""
}

echo "[1/4] starting server"
start_server

echo "[2/4] creating task"
go run ./cmd/todoopen task create --server "${BASE_URL}" --title "persist me" >/dev/null

stop_server

echo "[3/4] restarting server"
start_server

echo "[4/4] verifying task still exists"
if ! go run ./cmd/todoopen task list --server "${BASE_URL}" | grep -q "persist me"; then
  echo "persistence smoke test failed: task missing after restart" >&2
  exit 1
fi

echo "PASS: task persisted across restart"
