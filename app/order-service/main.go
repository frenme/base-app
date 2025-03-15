package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

	go func() {
		createKafkaTopicIfNotExists("example-topic", 2, 2)
	}()

	router := gin.Default()
	router.GET("/", pingHandler)
	router.GET("/redis", redisHandler)
	router.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}

func pingHandler(c *gin.Context) {
	user := models.User{Name: "Alice"}
	ctx := context.Background()

	log.Printf("ORDER_SERVICE: Order created #123")

	// postgres example
	var postgresData string
	db.PostgresMaster.QueryRow(ctx, "SELECT current_database()").Scan(&postgresData)

	// mongo example
	mongoDb := db.MongoClient.Database("exampleDb")
	collection := mongoDb.Collection("users")
	collection.InsertOne(ctx, user)
	mongoData, _ := collection.CountDocuments(ctx, bson.M{})

	// kafka example (order-service -> user-service)
	go func() {
		kafkaProducer()
	}()

	// http response
	c.JSON(http.StatusOK, gin.H{
		"object from another package": user,
		"postgresql data":             postgresData,
		"mongodb data":                mongoData,
	})
}

// REDIS --------------------------
func redisHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Redis": "OK",
	})

	var ctx = context.Background()
	db1 := db.MongoClient.Database("exampleDb")
	collection := db1.Collection("users")
	count, err := collection.CountDocuments(ctx, bson.D{})
	if count < 100000 {
		fmt.Println("Insert 100000 rows in MongoDB...")
		var docs []interface{}
		for i := 0; i < 100000; i++ {
			docs = append(docs, bson.D{
				{"index", i},
				{"value", fmt.Sprintf("data_%d", i)},
			})
		}
		collection.InsertMany(ctx, docs)
	}

	start := time.Now()
	cursor, err := collection.Find(ctx, bson.D{})
	var results []bson.M
	cursor.All(ctx, &results)
	durationMongo := time.Since(start)
	fmt.Printf("Read time without Redis: %v\n", durationMongo)

	redisNodes := strings.Split(os.Getenv("REDIS_CLUSTER"), ",")
	redisCluster := redis.NewClusterClient(&redis.ClusterOptions{Addrs: redisNodes})
	cacheKey := "mongo_data"
	cachedData, err := redisCluster.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		var mongoResults []bson.M
		dataBytes, _ := json.Marshal(mongoResults)
		redisCluster.Set(ctx, cacheKey, dataBytes, 10*time.Minute)
	} else {
		start = time.Now()
		var cachedResults []bson.M
		err = json.Unmarshal([]byte(cachedData), &cachedResults)
		duration := time.Since(start)
		fmt.Printf("Read time with Redis: %v\n", duration)
	}
}

// KAFKA --------------------------
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
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	defer conn.Close()

	controller, _ := conn.Controller()
	controllerAddr := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	controllerConn, _ := kafka.Dial("tcp", controllerAddr)
	defer controllerConn.Close()
	controllerConn.SetDeadline(time.Now().Add(10 * time.Second))

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		},
	}

	err := controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		if strings.Contains(err.Error(), "TopicExists") {
			log.Printf("Topic %s already exists", topic)
			return
		}
		log.Fatalf("Failed to create topic: %v", err)
	}
	log.Printf("Topic %s created successfully", topic)
}
