package smsproxy

import (
	"errors"
	"github.com/stretchr/testify/require"
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"

func TestSmsIsAccepted(t *testing.T) {
	// given
	smsProxy := ProdSmsProxy(fastsmsing.NewInMemoryClient())
	// when
	result, err := smsProxy.Send(SendMessage{phoneNumber, message})

	// then
	assert.NoError(t, err)
	messageStatus, err := smsProxy.GetStatus(result.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, Accepted, messageStatus)
}

func TestSmsIsNotAcceptedWhenErrorOccurredSavingMessage(t *testing.T) {
	// given
	repository := newFailingRepository()
	repository.saveError(someMessageID, errors.New("error saving sms"))
	sender := newBatchingClient(repository.build(), fastsmsing.NewInMemoryClient(), newConfig(), newClientStatistics())
	smsProxy := buildSmsProxy(repository.build(), sender, predefinedMessageID(someMessageID))

	// when
	result, err := smsProxy.Send(SendMessage{phoneNumber, message})

	// then
	require.Error(t, err)
	messageStatus, err := smsProxy.GetStatus(result.ID)
	require.NoError(t, err)
	require.EqualValues(t, NotFound, messageStatus)
}

func TestMessagesGetDeliveredWhenBatched(t *testing.T) {
	// given
	smsProxy := ProdSmsProxy(fastsmsing.NewInMemoryClient(), MinimumInBatchOption(2))
	smsProxy.Start()
	defer smsProxy.Stop()

	// when
	result, err := smsProxy.Send(SendMessage{phoneNumber, message})
	result2, err2 := smsProxy.Send(SendMessage{phoneNumber, message})

	// then
	assert.NoError(t, err)
	assert.NoError(t, err2)
	require.Eventually(t, func() bool {
		messageStatus, _ := smsProxy.GetStatus(result.ID)
		return messageStatus == Delivered
	}, 2*time.Second, 50*time.Millisecond)
	require.Eventually(t, func() bool {
		messageStatus, _ := smsProxy.GetStatus(result2.ID)
		return messageStatus == Delivered
	}, 2*time.Second, 50*time.Millisecond)
}

func TestMessagesGetDeliveredWithBatchingDisabled(t *testing.T) {
	// given
	smsProxy := ProdSmsProxy(fastsmsing.NewInMemoryClient(), DisableBatching())
	smsProxy.Start()
	defer smsProxy.Stop()

	// when
	result, err := smsProxy.Send(SendMessage{phoneNumber, message})

	// then
	assert.NoError(t, err)
	require.Eventually(t, func() bool {
		messageStatus, _ := smsProxy.GetStatus(result.ID)
		return messageStatus == Delivered
	}, 2*time.Second, 50*time.Millisecond)
}

func TestShouldReturnValidationErrorOnNewSms(t *testing.T) {
	tests := []struct {
		name    string
		message SendMessage
	}{
		{"empty phone", SendMessage{PhoneNumber: "", Message: "Some message"}},
		{"empty message", SendMessage{PhoneNumber: "1234", Message: ""}},
		{"phone number contains non-digits", SendMessage{PhoneNumber: "1234a", Message: "Some message"}},
		{"phone number contains country prefix", SendMessage{PhoneNumber: "+1234", Message: "Some message"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			smsProxy := ProdSmsProxy(fastsmsing.NewInMemoryClient(), MinimumInBatchOption(2))

			// when
			_, err := smsProxy.Send(test.message)

			// then
			require.Error(t, err)
			require.IsType(t, &ValidationError{}, err)
		})
	}
}
