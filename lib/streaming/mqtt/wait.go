package mqtt

import (
	"time"
	"github.com/hahahannes/e2e-go-utils/lib"
	"github.com/hahahannes/e2e-go-utils/lib/streaming"
)

func WaitForMQTTMessageReceived(regexTopic, regexMsg string, sendFnc func() error, timeout time.Duration, host, port string) (lib.MessageReceived, error) {
	msgChannel := make(chan streaming.Message)
	topicConfig := map[string]byte{
		"#": byte(2),
	}
	client := NewMQTTClient(host, port, topicConfig, msgChannel)
	client.ConnectMQTTBroker(nil, nil)
	defer client.CloseConnection()
	result, err := streaming.WaitForMessageOnTopicReceived(regexTopic, regexMsg, sendFnc, msgChannel, timeout)
	return result, err
}