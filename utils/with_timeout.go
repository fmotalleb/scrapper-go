package utils

import (
	"errors"
	"time"
)

// WithDeadline waits for a value from the provided channel until the specified timeout duration.
// It returns a pointer to the received value and nil error if successful, or nil and an error if the timeout is reached.
func WithDeadline[T any](channel <-chan T, timeout time.Duration) (*T, error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil, errors.New("timeout reached when waiting for the first result from the channel")
	case value := <-channel:
		return &value, nil
	}
}
