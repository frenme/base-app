# format: name:partitions:replications
TOPICS_LIST=(
  "example-topic:3:3"
  "another-topic:2:3"
)

wait_for_kafka() {
  HOST=$(echo "$BOOTSTRAP_SERVER" | cut -d: -f1)
  PORT=$(echo "$BOOTSTRAP_SERVER" | cut -d: -f2)
  echo "Waiting for Kafka on $HOST:$PORT..."
  while ! bash -c "echo > /dev/tcp/$HOST/$PORT" 2>/dev/null; do
    echo "Kafka is unavailable, wait 5 seconds..."
    sleep 5
  done
  echo "Kafka is available!"
}

wait_for_kafka

create_topic() {
  local topic_name=$1
  local partitions=$2
  local replication=$3

  if /opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server "$BOOTSTRAP_SERVER" --list | grep -q "^${topic_name}$"; then
    echo "Topic ${topic_name} already exists"
  else
    echo "Creating topic ${topic_name} (partitions=${partitions}, replication=${replication})..."
    /opt/bitnami/kafka/bin/kafka-topics.sh \
      --create \
      --bootstrap-server "$BOOTSTRAP_SERVER" \
      --replication-factor "$replication" \
      --partitions "$partitions" \
      --topic "$topic_name"
  fi
}

for config in "${TOPICS_LIST[@]}"; do
  IFS=":" read -r topic partitions replication <<< "$config"
  create_topic "$topic" "$partitions" "$replication"
done