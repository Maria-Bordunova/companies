package kafka

import (
	"companies/internal/company_ctx"
	"companies/internal/config"
	"companies/internal/entity/event"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type MessageBody struct {
	MessageType event.EventType `json:"message_type"`
	Payload     payload         `json:"payload"`
}

type payload struct {
	UId string `json:"uid"`
}

type EventProducer struct {
	config config.Kafka
}

func NewEventProducer(config config.Kafka) *EventProducer {
	return &EventProducer{config: config}
}

func (o *EventProducer) Produce(ctx context.Context, uid string, eventType event.EventType) error {
	sConfig := sarama.NewConfig()

	sConfig.Producer.Return.Successes = true
	sConfig.Producer.Return.Errors = true
	producer, err := sarama.NewSyncProducer([]string{o.config.Host}, sConfig)
	if err != nil {
		return err
	}
	defer producer.Close()

	messageBody := MessageBody{
		MessageType: eventType,
		Payload: payload{
			UId: uid,
		},
	}
	messageBodyJson, err := jsonToMsg(messageBody)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{Topic: o.config.Topic, Key: nil, Value: sarama.StringEncoder(messageBodyJson)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	log := company_ctx.Logger(ctx)
	log.With("partition_id", partition).
		With("offset", offset).
		With("msg", msg).
		Info("event produced")

	return nil
}

func jsonToMsg(b MessageBody) (string, error) {
	// Convert payload to JSON string
	jsonString, err := json.Marshal(b)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return "", err
	}
	return string(jsonString), nil
}
