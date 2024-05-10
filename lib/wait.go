package lib

import (
	"context"
	"fmt"
	"time"
)

func WaitForMessageReceived(ctx context.Context, sendFnc func() error, messageChannel any, matchFnc func(msg interface{}) (error, bool)) (MessageReceived, error) {
	// Start listening on the message channgel where incoming MQTT messages will land
	// then start command which will eventually lead to a message published

	resultChannel := make(chan MessageReceived)
	var messageReceived MessageReceived 
	
	go func() {
		for {
			select {
			case msg := <- messageChannel.(chan any):
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
			case <- ctx.Done(): 
				resultChannel <- MessageReceived{
					Received: false,
				}
				return
			}		
		}
	}()

	err := sendFnc()
	if err != nil {
		fmt.Printf("Error occured during setup: " + err.Error())
	}

	messageReceived = <- resultChannel

	if err != nil {
		return messageReceived, err
	}
	return messageReceived, nil

}

func WaitForStringReceived(expectedMsg string, sendFnc func() error, channel chan string) (MessageReceived, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
	defer cancel()
	return WaitForMessageReceived(ctx, sendFnc, channel, func (msg any) (error, bool) {
		if msg == expectedMsg {
			return nil, true
		}
		return nil, false
	})
}