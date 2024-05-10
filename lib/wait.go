package lib

import (
	"context"
	"fmt"
	"time"
	"regexp"
)

func WaitForMessageReceived[T any] (ctx context.Context, sendFnc func() error, messageChannel chan T, matchFnc func(msg interface{}) (error, bool)) (MessageReceived, error) {
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
				fmt.Println(msg)
				err, matched := matchFnc(msg)
				if err != nil {
					return
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
				}
				return
			}		
		}
	}()

	err := sendFnc()
	if err != nil {
		fmt.Printf("Error occured during send: " + err.Error())
		cancel()
		return messageReceived, err
	}

	messageReceived = <- resultChannel
	return messageReceived, nil

}

func WaitForStringReceived(regexMsg string, sendFnc func() error, channel chan string, timeout time.Duration) (MessageReceived, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return WaitForMessageReceived[string](ctx, sendFnc, channel, func (log any) (error, bool) {
		msgMatch, err := regexp.MatchString(regexMsg, log.(string))
		fmt.Println(fmt.Sprintf("Msg: %s, %t", log, msgMatch))
		return err, msgMatch
	})
}