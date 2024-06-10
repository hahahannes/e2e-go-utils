package kafka

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/hahahannes/e2e-go-utils/lib/streaming"
	"github.com/segmentio/kafka-go"
)

func NewConsumer(ctx context.Context, kafkaUrl string, topic string, messageChannel chan streaming.Message) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:                []string{kafkaUrl},
		Topic:                  topic,
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("close kafka reader")
				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}
				return
			default:
				m, err := r.ReadMessage(ctx)
				if err == io.EOF || err == context.Canceled {
					log.Println("close consumer for topic ")
					return
				}
				if err != nil {
					log.Println("ERROR: while consuming topic ", err)
					return
				}
				msg := string(m.Value)
				log.Println(fmt.Sprintf("Received %s on Topic: %s", msg, m.Topic))
				messageChannel <- streaming.Message{
					Topic: m.Topic,
					Value: msg,
				}
			}
		}
	}()
}