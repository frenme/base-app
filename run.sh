#!/bin/bash

# dev or prod
dirs=("app" "kong" "postgres")
# dirs=("app" "kong" "kafka" "redis" "mongo" "postgres")

# make files runnable
chmod +x kafka/setup.sh redis/setup.sh mongo/setup.sh mongo/generate-keyfile.sh
cd mongo || { exit 1; }
./generate-keyfile.sh
cd ..

# launching
for dir in "${dirs[@]}"; do
  echo "Starting service: $dir..."
  cd "$dir" || exit 1

  if [ "$dir" == "app" ]; then
    # configure scale parameter (example)
    docker compose up -d \
    --scale user-service=1 \
    --scale temp-service=1
  elif [ "$dir" == "elk" ]; then
    docker compose up setup
    docker compose up -d
  else
    docker compose up -d
  fi

  cd ..
done

# print logs
for dir in "${dirs[@]}"; do
  ( cd "$dir" && docker compose logs -f ) &
done

wait