package lib

import (
	"context"
	"fmt"
	"time"
	"regexp"
)

func WaitForMessageReceived[T any] (ctx context.Context, sendFnc func(context.Context) error, messageChannel chan T, matchFnc func(msg interface{}) (error, bool), logMessages bool) (MessageReceived, error) {
	// Start listening on the message channgel where incoming MQTT messages will land
	// then start command which will eventually lead to a message published

	resultChannel := make(chan MessageReceived)
	var messageReceived MessageReceived 
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case msg := <- messageChannel:
				if logMessages {
					fmt.Println(msg)
				}
				err, matched := matchFnc(msg)
				if logMessages {
					fmt.Println(matched)
				}
				if err != nil {
					resultChannel <- MessageReceived{
						Received: false,
						Error: err,
					}
				}

				if matched {
					resultChannel <- MessageReceived{
						Received: true,
						Message: msg,
					}
					return
				}
			case <- subCtx.Done(): 
				resultChannel <- MessageReceived{
					Received: false,
					Message: fmt.Sprintf("Wait loop finished due to Context: %s", subCtx.Err()),
				}
				return
			}		
		}
	}()

	go func() {
		// Send async, could be an application running
		fmt.Println("Start Func")
		err := sendFnc(ctx)
		if err != nil {
			fmt.Println("Error occured Func: " + err.Error())
			resultChannel <- MessageReceived{
				Received: false,
				Message: "sendFnc errored",
				Error: err,
			}
			fmt.Println("Cancel the waiting loop")
			cancel()
			return
		}
		fmt.Println("Func was executed successfully")
	}()

	messageReceived = <- resultChannel
	fmt.Printf("Received message: %+v", messageReceived)
	return messageReceived, nil
}

func WaitForStringReceived(regexMsg string, sendFnc func(context.Context) error, channel chan string, timeout time.Duration, logMessages bool) (MessageReceived, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return WaitForMessageReceived[string](ctx, sendFnc, channel, func (log any) (error, bool) {
		msgMatch, err := regexp.MatchString(regexMsg, log.(string))
		return err, msgMatch
	}, logMessages)
}