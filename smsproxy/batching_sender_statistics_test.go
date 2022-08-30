package smsproxy

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
	"testing"
	"time"
)

func TestSendingStatisticsOnFirstSuccess(t *testing.T) {
	// given
	clientMock := fastsmsing.NewClientMock()
	clientMock.On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Return(nil)
	statistics := newClientStatistics()
	sender := newBatchingClient(newRepository(), clientMock, newConfig().disableBatching(), statistics)
	// when
	err := sender.send(SendMessage{phoneNumber, message}, someMessageID)

	// then
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		return clientMock.AssertExpectations(t)
	}, 2*time.Second, 50*time.Millisecond)
	stats := statistics.GetStatistics()
	require.Equal(t, stats.failed, 0)
	require.Equal(t, stats.success, 1)
}

func TestShouldSendErrorStatisticsWhenRetryFailed(t *testing.T) {
	tests := []struct {
		attemptsFailed     int
		maxAttemptsSending int
		failedCount        int
		successCount       int
	}{
		{1, 3, 0, 1},
		{2, 3, 0, 1},
		{3, 3, 1, 0},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("attempt successful after %d/%d attempts", test.attemptsFailed+1, test.maxAttemptsSending), func(t *testing.T) {
			// given
			clientMock := fastsmsing.NewClientMock()
			clientMock.On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Return(errors.New("problem sending sms")).Times(test.attemptsFailed).
				On("Send", []fastsmsing.Message{{PhoneNumber: phoneNumber, Message: message, MessageID: someMessageID}}).Return(nil)
			statistics := newClientStatistics()
			sender := newBatchingClient(newRepository(), clientMock, newConfig().disableBatching().setMaxAttempts(test.maxAttemptsSending), statistics)

			// when
			err := sender.send(SendMessage{phoneNumber, message}, someMessageID)

			// then
			require.NoError(t, err)
			require.Eventually(t, func() bool {
				return clientMock.AssertExpectations(t)
			}, 2*time.Second, 50*time.Millisecond)
			stats := statistics.GetStatistics()
			require.Equal(t, stats.failed, test.failedCount)
			require.Equal(t, stats.success, test.successCount)
		})
	}
}
