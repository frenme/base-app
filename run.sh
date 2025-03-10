#!/bin/bash

dirs=("app" "kafka" "kong" "grafana-prometheus")

for dir in "${dirs[@]}"; do
  echo "Starting service: $dir..."
  cd "$dir" || exit 1

  if [ "$dir" == "service_with_scale" ]; then
    docker compose up -d --scale shared-service=3
  elif [ "$dir" == "elk" ]; then
    docker compose up setup
    docker compose up -d
  else
    docker compose up -d
  fi

  cd ..
done

for dir in "${dirs[@]}"; do
  ( cd "$dir" && docker compose logs -f ) &
done

wait