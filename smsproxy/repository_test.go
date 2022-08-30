package smsproxy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepositorySaveNewMessage(t *testing.T) {
	repository := newRepository()

	err := repository.save(someMessageID)
	assert.NoError(t, err)

	status, err := repository.get(someMessageID)
	assert.NoError(t, err)
	assert.EqualValues(t, Accepted, status)
}

func TestErrorWhenSavingSameIdTwice(t *testing.T) {
	repository := newRepository()

	err := repository.save(someMessageID)
	assert.NoError(t, err)

	err = repository.save(someMessageID)
	assert.Error(t, err)
}

func TestRepositoryUpdateState(t *testing.T) {
	repository := newRepository()
	err := repository.save(someMessageID)
	assert.NoError(t, err)
	err = repository.update(someMessageID, Failed)
	assert.NoError(t, err)

	status, getErr := repository.get(someMessageID)
	assert.EqualValues(t, Failed, status)
	assert.NoError(t, getErr)
}

func TestRepositoryUpdateNonExistingMessageState(t *testing.T) {
	repository := newRepository()
	err := repository.update(someMessageID, Failed)
	assert.Error(t, err)
}

func TestRepositoryGetNonExistingMessageState(t *testing.T) {
	repository := newRepository()
	status, _ := repository.get(someMessageID)
	assert.Equal(t, NotFound, status)
}

func TestCannotOverwriteFinalStatuses(t *testing.T) {
	tests := finalStatuses

	for _, finalStatus := range tests {
		t.Run(fmt.Sprintf("Test: cannot overwrite status %s", finalStatus), func(t *testing.T) {
			// given
			repository := newRepository()
			err := repository.save(someMessageID)
			require.NoError(t, err)
			err = repository.update(someMessageID, finalStatus)
			require.NoError(t, err)
			for _, otherStatus := range allStatuses {
				err = repository.update(someMessageID, otherStatus)
				require.Error(t, err)
			}
		})
	}
}

func TestChangingMessageStatuses(t *testing.T) {
	tests := finalStatuses
	for _, finalStatus := range tests {
		t.Run(fmt.Sprintf("Test: cannot overwrite status %s", finalStatus), func(t *testing.T) {
			// given
			repository := newRepository()
			// when
			err := repository.save(someMessageID)

			// then
			assert.NoError(t, err)
			status, _ := repository.get(someMessageID)
			assert.Equal(t, Accepted, status)

			// when
			err = repository.update(someMessageID, Confirmed)

			// then
			assert.NoError(t, err)
			status, _ = repository.get(someMessageID)
			assert.Equal(t, Confirmed, status)

			err = repository.update(someMessageID, finalStatus)
			assert.NoError(t, err)
			status, _ = repository.get(someMessageID)
			assert.Equal(t, finalStatus, status)
		})
	}
}
