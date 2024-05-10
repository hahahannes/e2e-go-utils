package container

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hahahannes/e2e-go-utils/lib"
	"github.com/testcontainers/testcontainers-go"
)

type LogConsumer struct {
	LogChannel chan string
}

func (c LogConsumer) Accept(rawLog testcontainers.Log) {
	log := string(rawLog.Content)
	c.LogChannel <- log
}

func startLoggingAndSend(ctx context.Context, logChannel chan string, container testcontainers.Container, sendFnc func() error) error {
	logConsumer := LogConsumer{
		LogChannel: logChannel,
	}
	// Call FollowOutput after goroutine, so that there exist a receiver for the log channel 
	fmt.Println("Set Log Follower at container")
	container.FollowOutput(logConsumer)
	err := container.StartLogProducer(ctx)
	if err != nil {
		return err
	}
	err = sendFnc()
	return err
}

func WaitForContainerLog(regexMsg string, sendFnc func() error, container testcontainers.Container) (lib.MessageReceived, error) {
	logChannel := make(chan string)

	ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
	defer cancel()
	return lib.WaitForMessageReceived(ctx, func() error {
		return startLoggingAndSend(ctx, logChannel, container, sendFnc)
	}, logChannel, func (log any) (error, bool) {
		msgMatch, err := regexp.MatchString(regexMsg, log.(string))
		return err, msgMatch
	})
}