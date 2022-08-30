package smsproxy

import (
	"github.com/stretchr/testify/require"
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
	"testing"
	"time"
)

func TestShouldSendMessageImmediately(t *testing.T) {
	// given
	minimumMessages := 0
	clientMock := fastsmsing.NewClientMock()
	clientMock.On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Return(nil)
	sender := newBatchingClient(newRepository(), clientMock, newConfig().setMinimumInBatch(minimumMessages), newClientStatistics())
	// when
	err := sender.send(SendMessage{phoneNumber, message}, someMessageID)

	// then
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		return clientMock.AssertExpectations(t)
	}, 2*time.Second, 50*time.Millisecond)
}

func TestShouldSendMessagesInBatch(t *testing.T) {
	// given
	minimumMessages := 2
	clientMock := fastsmsing.NewClientMock()
	clientMock.On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}, {PhoneNumber: phoneNumber2, Message: message2, MessageID: someMessageID2}}).Return(nil)
	sender := newBatchingClient(newRepository(), clientMock, newConfig().setMinimumInBatch(minimumMessages), newClientStatistics())

	// when
	err := sender.send(SendMessage{phoneNumber, message}, someMessageID)

	// then
	require.NoError(t, err)
	require.True(t, clientMock.AssertNumberOfCalls(t, "Send", 0))

	// when
	err = sender.send(SendMessage{phoneNumber2, message2}, someMessageID2)

	// then
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		return clientMock.AssertExpectations(t)
	}, 2*time.Second, 50*time.Millisecond)
}

func TestShouldSendMessageInIndependentBatches(t *testing.T) {
	// given
	minimumMessages := 0
	clientMock := fastsmsing.NewClientMock()
	clientMock.On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Return(nil)
	clientMock.On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber2, Message: message2, MessageID: someMessageID2}}).Return(nil)
	sender := newBatchingClient(newRepository(), clientMock, newConfig().setMinimumInBatch(minimumMessages), newClientStatistics())
	// when
	err := sender.send(SendMessage{phoneNumber, message}, someMessageID)
	err2 := sender.send(SendMessage{phoneNumber2, message2}, someMessageID2)

	// then
	require.NoError(t, err)
	require.NoError(t, err2)
	require.Eventually(t, func() bool {
		return clientMock.AssertExpectations(t)
	}, 2*time.Second, 50*time.Millisecond)
}
