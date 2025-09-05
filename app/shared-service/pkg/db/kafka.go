package db

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Publisher interface {
	Write(ctx context.Context, topic string, key, value []byte) error
	Close() error
}

type KafkaPublisher struct {
	writers map[string]*kafka.Writer
}

func NewKafkaPublisher(brokers []string, topics []string) *KafkaPublisher {
	ws := make(map[string]*kafka.Writer, len(topics))
	for _, t := range topics {
		ws[t] = &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        t,
			Balancer:     &kafka.RoundRobin{},
			RequiredAcks: kafka.RequireAll,
			Async:        true,
		}
	}
	return &KafkaPublisher{writers: ws}
}

func (p *KafkaPublisher) Write(ctx context.Context, topic string, key, value []byte) error {
	w, ok := p.writers[topic]
	if !ok {
		return fmt.Errorf("unknown topic: %s", topic)
	}
	return w.WriteMessages(ctx, kafka.Message{Key: key, Value: value})
}

func (p *KafkaPublisher) Close() error {
	var firstErr error
	for _, w := range p.writers {
		if err := w.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

type KafkaConsumer struct {
	r *kafka.Reader
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           brokers,
		GroupID:           groupID,
		GroupTopics:       topics,
		MinBytes:          1e3,
		MaxBytes:          10e6,
		CommitInterval:    0,
		HeartbeatInterval: time.Second,
		SessionTimeout:    10 * time.Second,
		RebalanceTimeout:  30 * time.Second,
	})
	return &KafkaConsumer{r: r}
}

func (c *KafkaConsumer) RunKafkaOutput(ctx context.Context) error {
	defer c.r.Close()

	backoff := time.Second
	for {
		m, err := c.r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			fmt.Println("kafka fetch error:", err)
			time.Sleep(backoff)
			if backoff < 10*time.Second {
				backoff *= 2
			}
			continue
		}
		backoff = time.Second

		fmt.Printf("kafka message topic=%s partition=%d offset=%d key=%s value=%s\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		if err := c.r.CommitMessages(ctx, m); err != nil {
			fmt.Println("kafka commit error:", err)
		}
	}
}
