#!/bin/sh
set -e

CONFIG=/app/config.yml
for arg in "$@"; do
  case "$arg" in
    --config=*) CONFIG="${arg#*=}" ;;
  esac
done

echo "Filling countries (skip if already in DB)..."
/app/cli database fill --config="$CONFIG"

echo "Starting API server..."
exec /app/server "$@"
