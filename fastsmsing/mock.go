package fastsmsing

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
)

type FastSmsingMock struct {
	mock.Mock
	lock sync.RWMutex
}

func (m *FastSmsingMock) Send(messages []Message) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	args := m.Called(messages)
	return args.Error(0)
}

func (m *FastSmsingMock) Subscribe(chan map[string]MessageStatus) {
	panic("not used for this task")
}

func (m *FastSmsingMock) Stop() {
	panic("not used for this task")
}
func (m *FastSmsingMock) AssertExpectations(t *testing.T) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.Mock.AssertExpectations(t)
}

func NewClientMock() *FastSmsingMock {
	return &FastSmsingMock{lock: sync.RWMutex{}}
}
