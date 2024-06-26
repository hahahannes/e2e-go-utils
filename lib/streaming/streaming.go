package streaming 

import (
	"context"
	"regexp"
	"github.com/hahahannes/e2e-go-utils/lib"
)


type Message struct {
	Value string
	Topic string
}

func WaitForMessageOnTopicReceived(ctx context.Context, regexTopic, regexMsg string, sendFnc func(context.Context) error, messageChannel chan Message, logMessages bool) (lib.MessageReceived, error) {
	// Start listening on the message channgel where incoming MQTT messages will land
	// then start command which will eventually lead to a message published

	return lib.WaitForMessageReceived(ctx, sendFnc, messageChannel, func (msg any) (error, bool) {
		value := msg.(Message).Value
		topic := msg.(Message).Topic
		msgMatch, err := regexp.MatchString(regexMsg, value)
		if err != nil {
			return err, false
		}
		topicMatch, err := regexp.MatchString(regexTopic, topic)
		if err != nil {
			return err, false
		}
		
		if msgMatch && topicMatch {
			return nil, true
		}
		return nil, false
	}, logMessages)
}