package smsproxy

type SmsProxy interface {
	Send(message SendMessage) (SendingResult, error)
	GetStatus(messageID string) (MessageStatus, error)
	Start()
	Stop()
}

type SendMessage struct {
	PhoneNumber PhoneNumber
	Message     Message
}

type PhoneNumber = string
type Message = string

type MessageID = string

type SendingResult struct {
	ID  MessageID
	Err error
}

type ValidationError struct {
	value string
}

func (v *ValidationError) Error() string {
	return v.value
}
