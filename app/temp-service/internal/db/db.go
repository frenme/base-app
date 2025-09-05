package db

import (
	"context"
	"fmt"
	"os"
	sdb "shared/pkg/db"
	"shared/pkg/logger"
	"strings"
)

var PostgresDB sdb.Postgres
var MongoDB sdb.Mongo
var RedisDB sdb.Redis
var KafkaPublisher *sdb.KafkaPublisher
var KafkaConsumers map[string]*sdb.KafkaConsumer
var log *logger.Logger

func init() {
	log = logger.New()
}

func InitConnections() {
	PostgresDB = sdb.GetPostgresClient(
		os.Getenv("POSTGRES_MASTER_CONNECTION"),
		os.Getenv("POSTGRES_REPLICAS_CONNECTION"),
		"temp-service",
	)

	MongoDB = sdb.CreateMongoClient(os.Getenv("MONGO_CONNECTION"), "temp-service")

	RedisDB = sdb.CreateRedisClient(os.Getenv("REDIS_CLUSTER"), "temp-service")

	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaTopics := strings.Split(os.Getenv("KAFKA_TOPICS"), ",")
	KafkaPublisher = sdb.NewKafkaPublisher(kafkaBrokers, kafkaTopics)
	KafkaConsumers = make(map[string]*sdb.KafkaConsumer, len(kafkaTopics))
	for _, t := range kafkaTopics {
		KafkaConsumers[t] = sdb.NewKafkaConsumer(kafkaBrokers, "temp-service-group", []string{t})
	}
}

func StartKafkaConsumers(ctx context.Context) {
	for topic, c := range KafkaConsumers {
		go func(tp string, consumer *sdb.KafkaConsumer) {
			if err := consumer.RunKafkaOutput(ctx); err != nil {
				fmt.Println("kafka consumer stopped:", tp, err)
			}
		}(topic, c)
	}
}
