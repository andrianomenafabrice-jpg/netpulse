#!/usr/bin/env bash
# Affiche les dernières lignes du log netpulse.log en continu.
LOGFILE="${1:-netpulse.log}"

if [ ! -f "$LOGFILE" ]; then
    echo "fichier de log introuvable : $LOGFILE"
    exit 1
fi

tail -f "$LOGFILE"