package mqtt

import (
	"context"

	"github.com/hahahannes/e2e-go-utils/lib"
	"github.com/hahahannes/e2e-go-utils/lib/streaming"
)

func WaitForMQTTMessageReceived(ctx context.Context, regexTopic, regexMsg string, sendFnc func(context.Context) error, host, port string, logMessages bool) (lib.MessageReceived, error) {
	msgChannel := make(chan streaming.Message)
	topicConfig := map[string]byte{
		"#": byte(2),
	}
	client := NewMQTTClient(host, port, topicConfig, msgChannel, true)
	client.ConnectMQTTBroker(nil, nil)
	defer client.CloseConnection()
	result, err := streaming.WaitForMessageOnTopicReceived(ctx, regexTopic, regexMsg, sendFnc, msgChannel, logMessages)
	return result, err
}