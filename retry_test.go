package retrycommand

import (
	"testing"
	"time"
)

func TestRetryCommand(t *testing.T) {
	tests := []struct {
		name          string
		commandName   string
		commandArgs   []string
		options       []RetryCommandOption
		expectedError string
	}{
		{
			name:          "Command succeeds immediately",
			commandName:   "echo",
			commandArgs:   []string{"hello"},
			options:       nil,
			expectedError: "",
		},
		{
			name:          "Command fails and retries",
			commandName:   "hello there",
			commandArgs:   []string{},
			options:       []RetryCommandOption{WithMaxRetries(3), WithTimeBetweenAttempts(1 * time.Second)},
			expectedError: `Attempt 2 finished with error: exec: "hello there": executable file not found in $PATH`,
		},
		{
			name:          "Command with expected duration timeout",
			commandName:   "sleep",
			commandArgs:   []string{"2"},
			options:       []RetryCommandOption{WithExpectedDuration(1 * time.Second), WithMaxRetries(2), WithTimeBetweenAttempts(1 * time.Second)},
			expectedError: `Attempt 1 finished with error: signal: killed`,
		},
		{
			name:          "Command retries but eventually succeeds",
			commandName:   "sh",
			commandArgs:   []string{"-c", "exit $(($RANDOM % 2))"},
			options:       []RetryCommandOption{WithMaxRetries(5), WithTimeBetweenAttempts(500 * time.Millisecond)},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RetryCommand(tt.commandName, tt.commandArgs, tt.options...)
			if err == nil {
				return
			}

			if err.Error() != tt.expectedError {
				t.Errorf("RetryCommand() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}
