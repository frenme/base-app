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
	"log/slog"
	"net/http"
	"order/internal/db"
	"os"
	"shared/pkg/models"
	"shared/pkg/utils"
	"strings"
	"time"
)

var logger *slog.Logger

func init() {
	baseHandler := slog.NewJSONHandler(os.Stdout, nil)
	handler := utils.LoggerHandler{Handler: baseHandler}
	logger = slog.New(handler)
}

func main() {
	db.InitConnections()
	defer db.CloseConnections()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", pingHandler)
	router.GET("/redis", redisHandler)
	router.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}

func pingHandler(c *gin.Context) {
	user := models.User{Name: "Alice"}
	ctx := context.Background()

	logger.Info("Order created #123")

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
var brokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
var writer = &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: "example-topic", Balancer: &kafka.RoundRobin{}}

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
