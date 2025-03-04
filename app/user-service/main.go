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
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "user1",
		})
	})
	go func() {
		r.Run(":" + os.Getenv("USER_SERVICE_PORT"))
	}()

	kafkaConsumer()
}

func kafkaConsumer() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "example-topic",
		GroupID: "my-group",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka consumer error: ", err)
		}
		log.Printf("Kafka message got: key=%s, value=%s, offset=%d", string(msg.Key), string(msg.Value), msg.Offset)
	}
}
