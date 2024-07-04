package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	KafkaUrl string
}

func NewKafkaProducer(kafkaUrl string) (producer *Producer) {
	return &Producer{
		KafkaUrl: kafkaUrl,
	}
}

func (producer *Producer) Produce(topic, message, key string, ctx context.Context) error {
	writer := kafka.Writer{
		Addr:     kafka.TCP(producer.KafkaUrl),
		Topic:   topic,
		AllowAutoTopicCreation: true,
	}

	// Retry logic as it can happen that kafka is in leadership handling while producing
	// or that there is a delay between auto topic creation and beeing ready to accept messages
	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		err = writer.WriteMessages(ctx, kafka.Message{
			Key: []byte(key),
			Value: []byte(message),
		})

		time.Sleep(time.Millisecond * 250)
		continue
	}

	err = writer.Close()
	return err
}
