package kafka

import (
	"time"
	"github.com/hahahannes/e2e-go-utils/lib"
	"github.com/hahahannes/e2e-go-utils/lib/streaming"
	"context"
)


func WaitForKafkaMessageReceived(topic, regexMsg string, sendFnc func() error, timeout time.Duration, host, port string) (lib.MessageReceived, error) {
	msgChannel := make(chan streaming.Message)
	url := "http://" + host + ":" + port
	ctx, cancel := context.WithCancel(context.Background())
	NewConsumer(ctx, url, topic, msgChannel)
	defer cancel()
	return streaming.WaitForMessageOnTopicReceived(topic, regexMsg, sendFnc, msgChannel, timeout)
}