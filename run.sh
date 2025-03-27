#!/bin/bash

# dev or prod
dirs=("app" "kong" "kafka" "redis" "mongo" "postgres")

# make files runnable
chmod +x kafka/setup.sh redis/setup.sh mongo/setup.sh mongo/create-keyfile.sh
cd mongo || { exit 1; }
./create-keyfile.sh
cd ..

# launching
for dir in "${dirs[@]}"; do
  echo "Starting service: $dir..."
  cd "$dir" || exit 1

  if [ "$dir" == "app" ]; then
    # configure scale parameter
    docker compose up -d \
    --scale user-service=2 \
    --scale order-service=1
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