package smsproxy

import (
	"fmt"
	"sync"
)

type repository interface {
	update(id MessageID, newStatus MessageStatus) error
	save(id MessageID) error
	get(id MessageID) (MessageStatus, error)
}

type inMemoryRepository struct {
	db   map[MessageID]MessageStatus
	lock sync.RWMutex
}

func (r *inMemoryRepository) save(id MessageID) error {
	// save given MessageID with ACCEPTED status. If given MessageID already exists, return an error

	if _, ok := r.db[id]; ok {
		return fmt.Errorf("message with the given id %s already exist", id)
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	r.db[id] = Accepted
	return nil
}

func (r *inMemoryRepository) get(id MessageID) (MessageStatus, error) {
	// return status of given message, by it's MessageID. If not found, return NOT_FOUND status
	r.lock.Lock()
	defer r.lock.Unlock()

	if val, ok := r.db[id]; ok {
		return val, nil
	}
	return NotFound, nil
}

func (r *inMemoryRepository) update(id MessageID, newStatus MessageStatus) error {
	// Set new status for a given message.
	// If message is not in ACCEPTED state already - return an error.
	// If current status is FAILED or DELIVERED - don't update it and return an error. Those are final statuses and cannot be overwritten.
	r.lock.Lock()
	defer r.lock.Unlock()

	val := r.db[id]
	acceptedIndex := -1
	valIndex := -1
	for i := range allStatuses {
		if allStatuses[i] == Accepted {
			acceptedIndex = i
		}
		if allStatuses[i] == val {
			valIndex = i
		}
	}

	if valIndex < 0 {
		return fmt.Errorf("invalid status")
	}

	if valIndex < acceptedIndex {
		return fmt.Errorf("to update the status should be at ACCEPTED")
	}

	if val == Failed || val == Delivered {
		return fmt.Errorf("unable to overwrite the status")
	}

	r.db[id] = newStatus
	return nil
}

func newRepository() repository {
	return &inMemoryRepository{db: make(map[MessageID]MessageStatus), lock: sync.RWMutex{}}
}
