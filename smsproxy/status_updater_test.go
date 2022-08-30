package smsproxy

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
	"testing"
	"time"
)

func TestUpdatesStatusOfSms(t *testing.T) {
	// given
	repository := newRepository()
	err := repository.save(someMessageID)
	require.NoError(t, err)
	updater := newStatusUpdater(repository)
	update := map[string]fastsmsing.MessageStatus{
		someMessageID: fastsmsing.DELIVERED,
	}
	updater.Start()

	// when
	go func() {
		updater.C <- update
		close(updater.C)
	}()

	// then
	require.Eventually(t, func() bool {
		status, _ := repository.get(someMessageID)
		return status == Delivered
	}, 2*time.Second, time.Millisecond*50)
}

func TestNotUpdatingStatusIfUnrecognisedStatus(t *testing.T) {
	// given
	repository := newRepository()
	err := repository.save(someMessageID)
	require.NoError(t, err)
	updater := newStatusUpdater(repository)
	unrecognisedStatus := fastsmsing.MessageStatus("something")
	update := map[string]fastsmsing.MessageStatus{
		someMessageID: unrecognisedStatus,
	}
	updater.Start()
	timeout := time.NewTimer(500 * time.Millisecond)
	defer timeout.Stop()

	// when
	go func() {
		updater.C <- update
		close(updater.C)
	}()

	// then
	// then
	select {
	case <-timeout.C:
		t.Fatalf("Timeout - didn't receive expected error")
	case err := <-updater.Errors:
		timeout.Stop()
		require.Equal(t, unrecognisedStatus, err.status, "statuses does not match")
		require.Equal(t, someMessageID, err.messageID, "message ID does not match")
		require.NotNil(t, err.err, "error should not be nil")
	}
	finalStatus, err := repository.get(someMessageID)
	require.NoError(t, err)
	require.Equal(t, Accepted, finalStatus, "Should not overwrite status after failed status mapping")
}

func TestErrorWhenRepositoryFailedToUpdateStatus(t *testing.T) {
	// given
	repository := newFailingRepository()
	err := repository.save(someMessageID)
	require.NoError(t, err)
	repository.updateError(someMessageID, fmt.Errorf("cannot update messageID: %s", someMessageID))
	updater := newStatusUpdater(repository.build())
	newStatus := fastsmsing.CONFIRMED
	update := map[string]fastsmsing.MessageStatus{
		someMessageID: newStatus,
	}
	updater.Start()
	timeout := time.NewTimer(500 * time.Millisecond)
	defer timeout.Stop()

	// when
	go func() {
		updater.C <- update
		close(updater.C)
	}()

	// then
	select {
	case <-timeout.C:
		t.Fatalf("Timeout - didn't receive expected error")
	case err := <-updater.Errors:
		timeout.Stop()
		require.Equal(t, newStatus, err.status, "statuses does not match")
		require.Equal(t, someMessageID, err.messageID, "message ID does not match")
		require.NotNil(t, err.err, "error should not be nil")
	}
	finalStatus, err := repository.get(someMessageID)
	require.NoError(t, err)
	require.Equal(t, Accepted, finalStatus, "Should not overwrite newStatus after failed newStatus mapping")
}

func TestUpdateStatusForAllMessages(t *testing.T) {
	// given
	repository := newRepository()
	err := repository.save(someMessageID)
	require.NoError(t, err)
	err = repository.save(someMessageID2)
	require.NoError(t, err)
	updater := newStatusUpdater(repository)
	update := map[string]fastsmsing.MessageStatus{
		someMessageID: fastsmsing.DELIVERED, someMessageID2: fastsmsing.DELIVERED,
	}
	updater.Start()
	timeout := time.NewTimer(500 * time.Millisecond)
	defer timeout.Stop()

	// when
	go func() {
		updater.C <- update
		close(updater.C)
	}()

	// then
	require.Eventually(t, func() bool {
		status, _ := repository.get(someMessageID)
		return status == Delivered
	}, 2*time.Second, time.Millisecond*50)
	require.Eventually(t, func() bool {
		status, _ := repository.get(someMessageID2)
		return status == Delivered
	}, 2*time.Second, time.Millisecond*50)

}
