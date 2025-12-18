#!/bin/sh
set -e

# Railway (and many other PaaS) inject a `PORT` env var.
# Homebox expects `HBOX_WEB_PORT`, so map it at runtime.
if [ -n "${PORT:-}" ] && [ -z "${HBOX_WEB_PORT:-}" ]; then
  export HBOX_WEB_PORT="$PORT"
fi

exec /app/api "$@"
