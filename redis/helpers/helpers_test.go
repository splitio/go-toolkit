package helpers

import (
	"errors"
	"testing"

	"github.com/splitio/go-toolkit/redis"
	"github.com/splitio/go-toolkit/redis/mocks"
)

func TestEnsureConnected(t *testing.T) {
	redisClient := mocks.MockClient{
		PingCall: func() redis.Result {
			return &mocks.MockResultOutput{
				ErrCall:    func() error { return nil },
				StringCall: func() string { return "PONG" },
			}
		},
	}
	EnsureConnected(&redisClient)
}

func TestEnsureConnectedError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "Couldn't connect to redis: someError" {
				t.Error("Expected \"Couldn't connect to redis: someError\". Got: ", r)
			}
		}
	}()
	redisClient := mocks.MockClient{
		PingCall: func() redis.Result {
			return &mocks.MockResultOutput{
				ErrCall:    func() error { return errors.New("someError") },
				StringCall: func() string { return "" },
			}
		},
	}
	EnsureConnected(&redisClient)
	t.Error("Should not reach this line")
}

func TestEnsureConnectedNotPong(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "Invalid redis ping response when connecting: PANG" {
				t.Error("Invalid redis ping response when connecting: PANG", r)
			}
		}
	}()
	redisClient := mocks.MockClient{
		PingCall: func() redis.Result {
			return &mocks.MockResultOutput{
				ErrCall:    func() error { return nil },
				StringCall: func() string { return "PANG" },
			}
		},
	}
	EnsureConnected(&redisClient)
	t.Error("Should not reach this line")
}
