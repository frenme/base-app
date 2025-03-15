#!/bin/bash

dirs=("app" "kong" "kafka" "redis" "mongo" "postgres")

for dir in "${dirs[@]}"; do
  echo "Starting service: $dir..."
  cd "$dir" || exit 1

  if [ "$dir" == "app" ]; then
    docker compose up -d --scale user-service=2
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