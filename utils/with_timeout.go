package utils

import (
	"errors"
	"time"
)

func WithDeadline[T any](channel <-chan T, timeout time.Duration) (*T, error) {
	timer := time.NewTimer(timeout)
	select {
	case <-timer.C:
		return nil, errors.New("timeout reached when waiting for the first result from the channel")
	case value := <-channel:
		return &value, nil
	}
}
