package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"order/internal/db"
	"os"
	"shared/pkg/models"
	"strings"
	"time"
)

func main() {
	db.InitConnections()
	defer db.CloseConnections()

	createKafkaTopicIfNotExists("example-topic", 2, 2)

	router := gin.Default()
	router.GET("/", pingHandler)
	router.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}

func pingHandler(c *gin.Context) {
	user := models.User{Name: "Alice"}
	ctx := context.Background()

	// postgres example
	var postgresData string
	db.PostgresMaster.QueryRow(ctx, "SELECT current_database()").Scan(&postgresData)

	// mongo example
	mongoDb := db.MongoClient.Database("exampleDb")
	collection := mongoDb.Collection("users")
	collection.InsertOne(ctx, user)
	mongoData, _ := collection.CountDocuments(ctx, bson.M{})

	// kafka example (order-service -> user-service)
	kafkaProducer()

	// http response
	c.JSON(http.StatusOK, gin.H{
		"object from another package": user,
		"postgresql data":             postgresData,
		"mongodb data":                mongoData,
	})
}

func kafkaProducer() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "example-topic",
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-1"),
			Value: []byte("Some message from Kafka"),
			Time:  time.Now(),
		},
	)
	if err != nil {
		log.Println("Kafka producer error: ", err)
	}
	log.Println("Kafka message has been sent")
}

func createKafkaTopicIfNotExists(topic string, numPartitions, replicationFactor int) {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	conn, _ := kafka.Dial("tcp", brokers[0])
	defer conn.Close()

	exists := kafkaTopicExists(conn, topic)
	if exists {
		return
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		},
	}

	err := conn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatalf("Kafka failed to create topic: %v", err)
	}
	log.Println("Kafka topic created successfully")
}

func kafkaTopicExists(conn *kafka.Conn, topic string) bool {
	partitions, err := conn.ReadPartitions()
	if err != nil {
		return false
	}
	for _, p := range partitions {
		if p.Topic == topic {
			return true
		}
	}
	return false
}
