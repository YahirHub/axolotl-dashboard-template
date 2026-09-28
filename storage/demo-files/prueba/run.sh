#!/usr/bin/env bash
set -euo pipefail

APP_NAME="Axolotl"
PORT=3000

echo "Iniciando $APP_NAME en puerto $PORT"

for item in go php javascript; do
  echo "Probando: $item"
done
