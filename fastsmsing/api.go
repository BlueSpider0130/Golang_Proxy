package fastsmsing

type PhoneNumber = string
type MessageContent = string
type MessageID = string

type Message struct {
	PhoneNumber PhoneNumber
	Message     MessageContent
	MessageID   MessageID
}

type FastSmsingClient interface {
	Send(messages []Message) error
	Subscribe(chan map[MessageID]MessageStatus)
	Stop()
}

type MessageStatus string

var CONFIRMED = MessageStatus("CONFIRMED")
var FAILED = MessageStatus("FAILED")
var DELIVERED = MessageStatus("DELIVERED")
