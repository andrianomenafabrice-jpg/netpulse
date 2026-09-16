#!/usr/bin/env bash
# Lance netpulse en tâche de fond et enregistre son PID.
set -euo pipefail

BIN="./netpulse"
PIDFILE="netpulse.pid"

if [ -f "$PIDFILE" ] && kill -0 "$(cat "$PIDFILE")" 2>/dev/null; then
    echo "netpulse tourne déjà (PID $(cat "$PIDFILE"))"
    exit 1
fi

nohup "$BIN" "$@" > netpulse.out 2>&1 &
echo $! > "$PIDFILE"
echo "netpulse lancé en arrière-plan (PID $(cat "$PIDFILE"))"