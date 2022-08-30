package smsproxy

import (
	"github.com/google/uuid"
)

var (
	someMessageID = uuid.New().String()
	phoneNumber   = "123456789"
	message       = "You got a reward!"

	someMessageID2 = uuid.New().String()
	phoneNumber2   = "123456576526"
	message2       = "Welcome to the club"
)
