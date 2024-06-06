package container

import (
	"context"
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

func NewLogConsumer() *LogConsumer {
	return &LogConsumer{
		LogChannel: make(chan string),
	}
}

func WaitForContainerLog(regexMsg string, sendFnc func(context.Context) error, logConsumer *LogConsumer, logMessages bool) (lib.MessageReceived, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
	defer cancel()
	return lib.WaitForMessageReceived[string](ctx, func(context.Context) error {
		return sendFnc(ctx)
	}, logConsumer.LogChannel, func (log any) (error, bool) {
		msgMatch, err := regexp.MatchString(regexMsg, log.(string))
		return err, msgMatch
	}, logMessages)
}