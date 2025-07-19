package integration

import (
	"context"
	"encoding/json"
	"os"

	"github.com/segmentio/kafka-go"
)

type KafkaProducerImpl struct {
	writer     *kafka.Writer
	orderTopic string
}

func NewKafkaProducer() *KafkaProducerImpl {
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	if brokers[0] == "" {
		brokers = []string{"kafka:9092"}
	}
	topic := os.Getenv("KAFKA_ORDER_TOPIC")
	if topic == "" {
		topic = "order_placed"
	}
	return &KafkaProducerImpl{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		orderTopic: topic,
	}
}

type OrderPlacedBook struct {
	BookID   int `json:"book_id"`
	Quantity int `json:"quantity"`
}

type orderPlacedEvent struct {
	OrderID int               `json:"order_id"`
	UserID  string            `json:"user_id"`
	Books   []OrderPlacedBook `json:"books"`
}

func (k *KafkaProducerImpl) PublishOrderPlaced(ctx context.Context, orderID int, userID string, books []OrderPlacedBook) error {
	evt := orderPlacedEvent{OrderID: orderID, UserID: userID, Books: books}
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	return k.writer.WriteMessages(ctx, kafka.Message{Value: data})
}
