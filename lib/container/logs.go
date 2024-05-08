package container

import (
	"context"
	"fmt"
	"time"
	"regexp"
	"github.com/testcontainers/testcontainers-go"
	"github.com/hahahannes/e2e-go-utils/lib"
)

type LogConsumer struct {
	LogChannel chan string
}

func (c LogConsumer) Accept(rawLog testcontainers.Log) {
	log := string(rawLog.Content)
	c.LogChannel <- log
}

func WaitForContainerLog(regexMsg string, sendFnc func() error, container testcontainers.Container) (lib.MessageReceived, error) {
	var messageReceived lib.MessageReceived 
	resultChannel := make(chan lib.MessageReceived)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(60) * time.Second)
	defer cancel()
	
	logChannel := make(chan string)

	go func() {
		fmt.Println("Start waiting for logs")
		for {
			select {
			case log := <- logChannel:
				fmt.Println(log)
				msgMatch, _ := regexp.MatchString(regexMsg, log)

				if msgMatch {
					resultChannel <- lib.MessageReceived{
						Received: true,
						Message: log,
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

	logConsumer := LogConsumer{
		LogChannel: logChannel,
	}
	// Call FollowOutput after goroutine, so that there exist a receiver for the log channel 
	fmt.Println("Set Log Follower at container")
	container.FollowOutput(logConsumer)
	err := container.StartLogProducer(timeoutCtx)
	if err != nil {
		cancel()
	}

	err = sendFnc()
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