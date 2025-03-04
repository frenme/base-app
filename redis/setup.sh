#!/bin/sh
for n in redis-7700 redis-7701 redis-7702 redis-7703 redis-7704 redis-7705; do
  port=$(echo $n | cut -d'-' -f2)
  until redis-cli -h $n -p $port ping | grep -q PONG; do
    echo "Waiting for $n on port $port..."
    sleep 1
  done
done
echo yes | redis-cli --cluster create redis-7700:7700 redis-7701:7701 redis-7702:7702 redis-7703:7703 redis-7704:7704 redis-7705:7705 --cluster-replicas 1
echo "Redis cluster created"