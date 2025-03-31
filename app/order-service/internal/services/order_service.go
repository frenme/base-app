package services

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
	"order/internal/repository"
	"os"
	"shared/pkg/models"
	"strings"
	"time"
)

var writer *kafka.Writer

func init() {
	var brokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	writer = &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: "example-topic", Balancer: &kafka.RoundRobin{}}
}

func GetOrderData(c *gin.Context) {
	user := models.User{Name: "Alice"}
	postgresData := repository.GetPostgresData()
	mongoData := repository.GetMongoData(user)

	go func() {
		kafkaProducer()
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":                     "e123xample of v1 route",
		"object from another package": user,
		"postgresql data":             postgresData,
		"mongodb data":                mongoData,
	})
}

func GetOrderAnotherData(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "example of v2 route",
	})
}

// KAFKA --------------------------
// order-service -> user-service
func kafkaProducer() {
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte("Some message from Kafka"),
			Time:  time.Now(),
		},
	)
	if err != nil {
		log.Println("Kafka producer error: ", err)
	}
	log.Println("Kafka message has been sent")
}

// OrderCache REDIS --------------------------
func OrderCache(c *gin.Context) {
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
