package retrycommand

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type RetryCommandOption func(*RetryCommandConfig)

type RetryCommandConfig struct {
	MaxRetries          int
	ExpectedDuration    time.Duration
	TimeBetweenAttempts time.Duration
}

func WithMaxRetries(maxRetries int) RetryCommandOption {
	return func(config *RetryCommandConfig) {
		config.MaxRetries = maxRetries
	}
}

func WithExpectedDuration(expectedDuration time.Duration) RetryCommandOption {
	return func(config *RetryCommandConfig) {
		config.ExpectedDuration = expectedDuration
	}
}

func WithTimeBetweenAttempts(timeBetweenAttempts time.Duration) RetryCommandOption {
	return func(config *RetryCommandConfig) {
		config.TimeBetweenAttempts = timeBetweenAttempts
	}
}

func RetryCommand(commandName string, commandArgs []string, opts ...RetryCommandOption) error {
	config := &RetryCommandConfig{
		MaxRetries:          3,
		ExpectedDuration:    15 * time.Second,
		TimeBetweenAttempts: 15 * time.Second,
	}

	for _, opt := range opts {
		opt(config)
	}

	var err error
	for i := 0; i < config.MaxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), config.ExpectedDuration)
		defer cancel()

		cmd := exec.CommandContext(ctx, commandName, commandArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err == nil {
			return nil
		}

		fmt.Printf("Attempt %d finished with error: %v\n", i, err)
		time.Sleep(config.TimeBetweenAttempts)
	}
	return err
}
