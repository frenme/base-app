#!/bin/bash

dirs=("app" "kong" "kafka" "redis" "postgres" "mongo" "elk" "grafana-prometheus")

for dir in "${dirs[@]}"; do
  echo "Stopping service: $dir..."
  cd "$dir" || exit 1

  docker compose down -v

  cd ..
done
