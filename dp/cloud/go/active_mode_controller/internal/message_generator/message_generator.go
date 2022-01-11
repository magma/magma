package message_generator

import (
	"encoding/json"
	"time"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
	"magma/dp/cloud/go/active_mode_controller/protos/requests"
)

type messageGenerator struct {
	heartbeatTimeout  time.Duration
	inactivityTimeout time.Duration
}

func NewMessageGenerator(heartbeatTimeout time.Duration, inactivityTimeout time.Duration) *messageGenerator {
	return &messageGenerator{
		heartbeatTimeout:  heartbeatTimeout,
		inactivityTimeout: inactivityTimeout,
	}
}

func (m *messageGenerator) GenerateMessages(state *active_mode.State, now time.Time) []*requests.RequestPayload {
	var messages []message
	for _, config := range state.ActiveModeConfigs {
		generator := m.getPerCbsdMessageGenerator(config, now)
		perCbsd := generator.generateMessages(config)
		perCbsd = filterMessages(config.GetCbsd().GetPendingRequests(), perCbsd)
		messages = append(messages, perCbsd...)
	}
	payloads := make([]*requests.RequestPayload, len(messages))
	for i, m := range messages {
		payloads[i] = messageToRequest(m)
	}
	return payloads
}

func (m *messageGenerator) getPerCbsdMessageGenerator(config *active_mode.ActiveModeConfig, now time.Time) perCbsdMessageGenerator {
	if config.GetDesiredState() == active_mode.CbsdState_Unregistered {
		if config.GetCbsd().GetState() == active_mode.CbsdState_Registered {
			return &deregistrationMessageGenerator{}
		}
		return &noMessageGenerator{}
	}
	if config.GetCbsd().GetLastSeenTimestamp() < now.Add(-m.inactivityTimeout).Unix() {
		return &relinquishmentMessageGenerator{}
	}
	if config.GetCbsd().GetState() == active_mode.CbsdState_Unregistered {
		return &registerMessageGenerator{}
	}
	if len(config.GetCbsd().GetGrants()) != 0 {
		deadlineUnix := now.Add(m.heartbeatTimeout).Unix()
		return &heartbeatMessageGenerator{deadlineUnix: deadlineUnix}
	}
	return &firstNotNullMessageGenerator{
		generators: []perCbsdMessageGenerator{
			&grantMessageGenerator{},
			&spectrumInquiryMessageGenerator{},
		},
	}
}

type perCbsdMessageGenerator interface {
	generateMessages(config *active_mode.ActiveModeConfig) []message
}

type message interface {
	name() string
}

func messageToRequest(m message) *requests.RequestPayload {
	data := map[string][]interface{}{
		m.name() + "Request": {m},
	}
	payload, _ := json.Marshal(data)
	return &requests.RequestPayload{Payload: string(payload)}
}

type noMessageGenerator struct{}

func (*noMessageGenerator) generateMessages(_ *active_mode.ActiveModeConfig) []message {
	return nil
}

type firstNotNullMessageGenerator struct {
	generators []perCbsdMessageGenerator
}

func (f *firstNotNullMessageGenerator) generateMessages(config *active_mode.ActiveModeConfig) []message {
	for _, g := range f.generators {
		messages := g.generateMessages(config)
		if len(messages) > 0 {
			return messages
		}
	}
	return nil
}
