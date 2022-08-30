package smsproxy

import (
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
	"strings"
)

type batchingSmsProxy struct {
	repository     repository
	batchingClient batchingClient
	generator      idGenerator
	updater        statusUpdater
	fastclient     fastsmsing.FastSmsingClient
}

func ProdSmsProxy(client fastsmsing.FastSmsingClient, options ...ConfigOption) SmsProxy {
	config := newConfig()
	for _, option := range options {
		option(&config)
	}
	repository := newRepository()
	updater := newStatusUpdater(repository)
	sender := newBatchingClient(repository, client, config, newClientStatistics())
	proxy := batchingSmsProxy{repository: repository, batchingClient: sender, generator: uuidGenerate, updater: updater, fastclient: client}
	return &proxy
}

func buildSmsProxy(repository repository, sender batchingClient, generator idGenerator) SmsProxy {
	proxy := batchingSmsProxy{repository: repository, batchingClient: sender, generator: generator}
	return &proxy
}

func (p *batchingSmsProxy) Start() {
	p.updater.Start()
	p.fastclient.Subscribe(p.updater.C)
}

func (p *batchingSmsProxy) Stop() {
	p.fastclient.Stop()
}

func (p *batchingSmsProxy) GetStatus(messageID string) (MessageStatus, error) {
	return p.repository.get(messageID)
}

func (p *batchingSmsProxy) Send(message SendMessage) (SendingResult, error) {
	err := validate(message)
	if err != nil {
		return SendingResult{}, err
	}
	messageID := p.generator()
	if err := p.batchingClient.send(message, messageID); err != nil {
		return SendingResult{}, err
	}
	return SendingResult{ID: messageID}, nil
}

func validate(message SendMessage) *ValidationError {
	if len(message.Message) == 0 {
		return &ValidationError{"empty message"}
	}
	if len(message.PhoneNumber) == 0 {
		return &ValidationError{"empty phone number"}
	}
	if !validPhoneNumber(message.PhoneNumber) {
		return &ValidationError{"phone number can only contain digits"}
	}
	return nil
}

func validPhoneNumber(number PhoneNumber) bool {
	var isNotDigit = func(c rune) bool { return c < '0' || c > '9' }
	return strings.IndexFunc(number, isNotDigit) == -1
}
