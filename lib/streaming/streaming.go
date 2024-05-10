package streaming 

import (
	"context"
	"fmt"
	"time"
	"regexp"
	"github.com/hahahannes/e2e-go-utils/lib"
)


type Message struct {
	Value string
	Topic string
}

func WaitForMessageOnTopicReceived(regexTopic, regexMsg string, sendFnc func(context.Context) error, messageChannel chan Message, timeout time.Duration) (lib.MessageReceived, error) {
	// Start listening on the message channgel where incoming MQTT messages will land
	// then start command which will eventually lead to a message published
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return lib.WaitForMessageReceived(ctx, sendFnc, messageChannel, func (msg any) (error, bool) {
		value := msg.(Message).Value
		topic := msg.(Message).Topic
		fmt.Println(topic + " - " + value)
		msgMatch, _ := regexp.MatchString(regexMsg, value)
		topicMatch, _ := regexp.MatchString(regexTopic, topic)

		if msgMatch && topicMatch {
			return nil, true
		}
		return nil, false
	})
}