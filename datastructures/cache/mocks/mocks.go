package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type LayerMock struct {
    mock.Mock
}

func (m *LayerMock) Get(ctx context.Context, key string) (string, error) {
    args := m.Called(ctx, key)
    return args.String(0), args.Error(1)
}

func (m *LayerMock) Set(ctx context.Context, key string, value string) error {
    args := m.Called(ctx, key, value)
    return args.Error(0)
}

