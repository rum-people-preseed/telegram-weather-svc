package mocks

import "github.com/stretchr/testify/mock"

type MockSender struct {
	mock.Mock
}

func (m *MockSender) SendGetRequest(url string) ([]byte, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
