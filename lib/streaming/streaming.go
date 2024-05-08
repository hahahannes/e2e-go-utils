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

func WaitForMessageOnTopicReceived(regexTopic, regexMsg string, sendFnc func() error, messageChannel chan Message, timeout time.Duration) (lib.MessageReceived, error) {
	// Start listening on the message channgel where incoming MQTT messages will land
	// then start command which will eventually lead to a message published

	var messageReceived lib.MessageReceived 
	resultChannel := make(chan lib.MessageReceived)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)

	go func() {
		fmt.Println("Start waiting for messages on: " + regexTopic)
		for {
			select {
			case msg := <- messageChannel:
				value := msg.Value
				topic := msg.Topic
				fmt.Println(topic + " - " + value)
				msgMatch, _ := regexp.MatchString(regexMsg, value)
				topicMatch, _ := regexp.MatchString(regexTopic, topic)

				if msgMatch && topicMatch {
					resultChannel <- lib.MessageReceived{
						Received: true,
						Message: value,
					}
					return
				}
			case <- timeoutCtx.Done(): 
				resultChannel <- lib.MessageReceived{
					Received: false,
				}
				return
			}		
		}
	}()

	err := sendFnc()
	if err != nil {
		fmt.Printf("Error occured during setup: " + err.Error())
		cancel()
	}

	messageReceived = <- resultChannel

	if err != nil {
		return messageReceived, err
	}
	return messageReceived, nil

}