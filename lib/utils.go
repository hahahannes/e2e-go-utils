package lib

import (
	"time"
)

func CheckConditionWithRetry(funcCondition func() bool, numberRetries int, timeout int) bool {
	// Use to check for a condition that is expected to be true within some time
	loops := 0

	for loops < numberRetries {
		result := funcCondition()
		if result {
			return true
		}
		loops++
		time.Sleep(time.Duration(timeout) * time.Second)
	}

	return false
}

type MessageReceived struct {
	Received bool
	Message interface{}
	Error error
}