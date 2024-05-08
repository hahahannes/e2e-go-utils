package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	Context context.Context
	KafkaUrl string
}

func NewKafkaProducer(ctx context.Context, kafkaUrl string) (producer *Producer) {
	return &Producer{
		Context: ctx,
		KafkaUrl: kafkaUrl,
	}
}

func (producer *Producer) Produce(topic, message, key string) error {
	writer := kafka.Writer{
		Addr:     kafka.TCP(producer.KafkaUrl),
		Topic:   topic,
		Balancer: &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	// Retry logic as it can happen that kafka is in leadership handling while producing
	// or that there is a delay between auto topic creation and beeing ready to accept messages
	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(producer.Context, 10*time.Second)
		defer cancel()
		err := writer.WriteMessages(ctx, kafka.Message{
			Key: []byte(key),
			Value: []byte(message),
		})
		if err != nil {
			return err 
		}

		time.Sleep(time.Millisecond * 250)
		continue
	}

	err = writer.Close()
	return err
}
