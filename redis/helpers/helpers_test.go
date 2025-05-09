package helpers

import (
	"errors"
	"testing"

	"github.com/splitio/go-toolkit/v6/redis/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEnsureConnected(t *testing.T) {
	var resMock mocks.MockResultOutput
	resMock.On("String").Return(pong).Once()
	resMock.On("Err").Return(nil).Once()

	var clientMock mocks.MockClient
	clientMock.On("Ping").Return(&resMock).Once()
	EnsureConnected(&clientMock)
}

func TestEnsureConnectedError(t *testing.T) {
	var resMock mocks.MockResultOutput
	resMock.On("String").Return("").Once()
	resMock.On("Err").Return(errors.New("someError")).Once()

	var clientMock mocks.MockClient
	clientMock.On("Ping").Return(&resMock).Once()

	assert.Panics(t, func() { EnsureConnected(&clientMock) })
}

func TestEnsureConnectedNotPong(t *testing.T) {
	var resMock mocks.MockResultOutput
	resMock.On("String").Return("PANG").Once()
	resMock.On("Err").Return(nil).Once()

	var clientMock mocks.MockClient
	clientMock.On("Ping").Return(&resMock).Once()

	assert.Panics(t, func() { EnsureConnected(&clientMock) })
}
