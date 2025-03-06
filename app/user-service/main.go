package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	go func() {
		kafkaConsumer()
	}()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "user1",
		})
	})
	r.Run(":" + os.Getenv("USER_SERVICE_PORT"))
}

func kafkaConsumer() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    "example-topic",
		GroupID:  "my-group",
		MinBytes: 0,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka consumer error: ", err)
			continue
		}
		log.Printf("Kafka message got: key=%s, value=%s, offset=%d", string(msg.Key), string(msg.Value), msg.Offset)
	}
}
