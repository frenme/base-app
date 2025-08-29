#!/bin/bash

echo "=== Redis Cluster Diagnostic Script ==="
echo

# Проверка контейнеров
echo "=== Docker Containers Status ==="
docker ps | grep redis
echo

# Проверка отдельных узлов
echo "=== Individual Node Ping Tests ==="
nodes="redis-7700:7700 redis-7701:7701 redis-7702:7702 redis-7703:7703 redis-7704:7704 redis-7705:7705"

for node in $nodes; do
    host=$(echo $node | cut -d: -f1)
    port=$(echo $node | cut -d: -f2)
    echo -n "Testing $node: "
    
    result=$(docker exec $host redis-cli -p $port ping 2>/dev/null)
    if [ "$result" = "PONG" ]; then
        echo "✅ OK"
    else
        echo "❌ FAILED"
    fi
done
echo

# Проверка cluster info
echo "=== Cluster Information ==="
docker exec redis-7700 redis-cli -p 7700 cluster info | grep cluster_state
docker exec redis-7700 redis-cli -p 7700 cluster info | grep cluster_slots
docker exec redis-7700 redis-cli -p 7700 cluster info | grep cluster_known_nodes
echo

# Проверка cluster nodes
echo "=== Cluster Nodes ==="
docker exec redis-7700 redis-cli -p 7700 cluster nodes
echo

# Проверка подключения из temp-service
echo "=== Connectivity from temp-service ==="
temp_container=$(docker ps | grep temp-service | awk '{print $NF}')
if [ -n "$temp_container" ]; then
    echo "Testing from temp-service container: $temp_container"
    
    for node in $nodes; do
        host=$(echo $node | cut -d: -f1)
        port=$(echo $node | cut -d: -f2)
        echo -n "  $node: "
        
        result=$(docker exec $temp_container nc -zv $host $port 2>&1)
        if echo "$result" | grep -q "open"; then
            echo "✅ Connected"
        else
            echo "❌ Failed"
        fi
    done
else
    echo "temp-service container not found"
fi
echo

# Проверка DNS резолва из temp-service
echo "=== DNS Resolution from temp-service ==="
if [ -n "$temp_container" ]; then
    for host in redis-7700 redis-7701 redis-7702; do
        echo -n "  $host: "
        result=$(docker exec $temp_container nslookup $host 2>/dev/null | grep "Address:" | tail -1 | awk '{print $2}')
        if [ -n "$result" ]; then
            echo "✅ $result"
        else
            echo "❌ Failed to resolve"
        fi
    done
else
    echo "temp-service container not found"
fi
echo

echo "=== Diagnostic Complete ==="
