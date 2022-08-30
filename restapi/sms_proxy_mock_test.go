package restapi

import (
	"errors"
	"gitlab.com/devskiller-tasks/messaging-app-golang/smsproxy"
)

type errorSmsProxy struct {
	errorContent string
}

func (s *errorSmsProxy) Send(message smsproxy.SendMessage) (smsproxy.SendingResult, error) {
	return smsproxy.SendingResult{}, errors.New(s.errorContent)
}

func (s *errorSmsProxy) GetStatus(messageID string) (smsproxy.MessageStatus, error) {
	return "", errors.New(s.errorContent)
}

func (s *errorSmsProxy) Start() {
}

func (s *errorSmsProxy) Stop() {
}

func newErrorSmsProxy(errorContent string) smsproxy.SmsProxy {
	proxy := errorSmsProxy{errorContent: errorContent}
	return &proxy
}
