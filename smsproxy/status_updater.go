package smsproxy

import (
	"fmt"
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
)

type statusUpdateError struct {
	err       error
	messageID fastsmsing.MessageID
	status    fastsmsing.MessageStatus
}

type statusUpdater struct {
	C          chan map[fastsmsing.MessageID]fastsmsing.MessageStatus
	Errors     chan statusUpdateError
	repository repository
}

func newStatusUpdater(repository repository) statusUpdater {
	return statusUpdater{
		C:          make(chan map[string]fastsmsing.MessageStatus),
		repository: repository,
		Errors:     make(chan statusUpdateError)}
}

func (u statusUpdater) Start() {
	// When started, statusUpdater should continue reading from statusUpdater.C channel, where updates will be delivered, and save them using repository.update(...)
	// fastssmsing.MessageStatus should be mapped to smsproxy.MessageStatus using `mapToInternalStatus` function before updating state using repository.update(...)
	// When mapping to internal status fails, or updating status using repository.update(...) fails - you should asynchronously send statusUpdateError to statusUpdater.Errors channel
	go func() {
		select {
		case statusMap := <-u.C:
			for k, v := range statusMap {
				status, err := mapToInternalStatus(v)
				if err != nil {
					u.Errors <- statusUpdateError{
						err:       err,
						messageID: k,
						status:    v,
					}
				} else {
					err := u.repository.update(k, status)
					if err != nil {
						u.Errors <- statusUpdateError{
							err:       err,
							messageID: k,
							status:    v,
						}
					}
				}
			}
		}
	}()
}

func mapToInternalStatus(status fastsmsing.MessageStatus) (MessageStatus, error) {
	for _, mappableStatus := range allStatuses {
		if string(status) == string(mappableStatus) {
			return mappableStatus, nil
		}
	}
	return "", fmt.Errorf("cannot map status %s to any known status", status)
}
