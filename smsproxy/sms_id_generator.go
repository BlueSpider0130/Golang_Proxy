package smsproxy

import "github.com/google/uuid"

type idGenerator = func() MessageID

func uuidGenerate() MessageID {
	return uuid.New().String()
}

func predefinedMessageID(id string) idGenerator {
	return func() MessageID {
		return id
	}
}
