package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"github.com/hahahannes/e2e-go-utils/lib/streaming"
)

func NewConsumer(ctx context.Context, kafkaUrl string, topic string, messageChannel chan streaming.Message) {
	r := kafka.NewReader(kafka.ReaderConfig{
		CommitInterval:         0, //synchronous commits
		Brokers:                []string{kafkaUrl},
		Topic:                  topic,
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("close kafka reader ")
				return
			default:
				m, err := r.FetchMessage(ctx)
				if err == io.EOF || err == context.Canceled {
					log.Println("close consumer for topic ")
					return
				}
				if err != nil {
					log.Println("ERROR: while consuming topic ", err)
					return
				}
				
				messageChannel <- streaming.Message{
					Topic: m.Topic,
					Value: string(m.Value),
				}
			}
		}
	}()
}