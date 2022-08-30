package restapi

import (
	"gitlab.com/devskiller-tasks/messaging-app-golang/smsproxy"
)

type SendSmsRequest struct {
	PhoneNumber string
	Content     string
}

type SmsSendResponse struct {
	ID smsproxy.MessageID
}

type SmsStatusResponse struct {
	Status smsproxy.MessageStatus
}

type HttpErrorResponse struct {
	Error string
}
