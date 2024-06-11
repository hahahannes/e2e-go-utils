package kafka

import (
	"time"
	"github.com/hahahannes/e2e-go-utils/lib"
	"github.com/hahahannes/e2e-go-utils/lib/streaming"
	"context"
)


func WaitForKafkaMessageReceived(ctx context.Context, topic, regexMsg string, sendFnc func(context.Context) error, timeout time.Duration, host, port string, logMessages bool) (lib.MessageReceived, error) {
	msgChannel := make(chan streaming.Message)
	url := host + ":" + port
	NewConsumer(ctx, url, topic, msgChannel)
	return streaming.WaitForMessageOnTopicReceived(ctx, topic, regexMsg, sendFnc, msgChannel, timeout, logMessages)
}