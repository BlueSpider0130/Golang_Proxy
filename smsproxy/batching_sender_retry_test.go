package smsproxy

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
	"testing"
	"time"
)

func TestShouldAttemptSendingBatchOnlyOnce(t *testing.T) {
	tests := []struct {
		attemptsCount int
	}{
		{1},
		{0},
		{-1},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("attempts count: %d", test.attemptsCount), func(t *testing.T) {
			// given
			clientMock := fastsmsing.NewClientMock()
			clientMock.On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Return(errors.New("error sending messages"))
			clientMock.On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Maybe().Return(nil)
			sender := newBatchingClient(newRepository(), clientMock, newConfig().disableBatching().setMaxAttempts(test.attemptsCount), newClientStatistics())

			// when
			_ = sender.send(SendMessage{phoneNumber, message}, someMessageID)

			// then
			require.Eventually(t, func() bool {
				return clientMock.AssertExpectations(t)
			}, 2*time.Second, 50*time.Millisecond)
			require.True(t, clientMock.AssertNumberOfCalls(t, "Send", 1))
		})
	}
}

func TestSendingShouldBeRetriedAfterError(t *testing.T) {
	// given
	clientMock := fastsmsing.NewClientMock()
	clientMock.
		On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Return(errors.New("error sending messages")).Twice().
		On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Return(nil)

	sender := newBatchingClient(newRepository(), clientMock, newConfig().disableBatching().setMaxAttempts(3), newClientStatistics())

	// when
	_ = sender.send(SendMessage{phoneNumber, message}, someMessageID)

	// then
	require.Eventually(t, func() bool {
		return clientMock.AssertExpectations(t)
	}, 2*time.Second, 50*time.Millisecond)
	require.True(t, clientMock.AssertNumberOfCalls(t, "Send", 3))
}
