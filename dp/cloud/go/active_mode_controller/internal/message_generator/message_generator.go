package message_generator

import (
	"context"
	"fmt"
	"time"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/crud"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/message"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas_helpers"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
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

type Message interface {
	fmt.Stringer
	Send(context.Context, message.ClientProvider) error
}

func (m *messageGenerator) GenerateMessages(state *active_mode.State, now time.Time) []Message {
	var requests []*sas.Request
	var actions []*crud.Action
	for _, config := range state.ActiveModeConfigs {
		g := m.getPerCbsdMessageGenerator(config, now)
		reqs := g.sas.GenerateRequests(config)
		reqs = sas_helpers.Filter(config.GetCbsd().GetPendingRequests(), reqs)
		requests = append(requests, reqs...)
		acts := g.crud.GenerateActions(config)
		actions = append(actions, acts...)
	}
	var msgs []Message
	payloads := sas_helpers.Build(requests)
	for _, payload := range payloads {
		msgs = append(msgs, message.NewSasMessage(payload))
	}
	for _, action := range actions {
		msgs = append(msgs, message.NewDeleteMessage(action.SerialNumber))
	}
	return msgs
}

func (m *messageGenerator) getPerCbsdMessageGenerator(config *active_mode.ActiveModeConfig, now time.Time) *perCbsdMessageGenerator {
	cbsd := config.GetCbsd()
	if pendingStateChange(cbsd.GetPendingRequests()) {
		return &perCbsdMessageGenerator{
			sas:  &noMessageGenerator{},
			crud: &noActionsGenerator{},
		}
	}
	if cbsd.GetIsDeleted() && cbsd.GetState() == active_mode.CbsdState_Unregistered {
		return &perCbsdMessageGenerator{
			sas:  &noMessageGenerator{},
			crud: &deleteActionGenerator{},
		}
	}
	return &perCbsdMessageGenerator{
		sas:  m.getPerCbsdSasRequestGenerator(config, now),
		crud: &noActionsGenerator{},
	}
}

var stateChangeRequestType = map[active_mode.RequestsType]bool{
	active_mode.RequestsType_RegistrationRequest:   true,
	active_mode.RequestsType_DeregistrationRequest: true,
}

func pendingStateChange(requests []*active_mode.Request) bool {
	for _, r := range requests {
		if stateChangeRequestType[r.Type] {
			return true
		}
	}
	return false
}

func (m *messageGenerator) getPerCbsdSasRequestGenerator(config *active_mode.ActiveModeConfig, now time.Time) sasRequestGenerator {
	cbsd := config.GetCbsd()
	if cbsd.GetIsDeleted() && cbsd.GetState() == active_mode.CbsdState_Registered {
		return sas.NewDeregistrationRequestGenerator()
	}
	if config.GetDesiredState() == active_mode.CbsdState_Unregistered {
		if cbsd.GetState() == active_mode.CbsdState_Registered {
			return sas.NewDeregistrationRequestGenerator()
		}
		return &noMessageGenerator{}
	}
	if cbsd.GetLastSeenTimestamp() < now.Add(-m.inactivityTimeout).Unix() {
		return sas.NewRelinquishmentRequestGenerator()
	}
	if cbsd.GetState() == active_mode.CbsdState_Unregistered {
		return sas.NewRegistrationRequestGenerator()
	}
	if len(cbsd.GetGrants()) != 0 {
		nextSend := now.Add(m.heartbeatTimeout).Unix()
		return sas.NewHeartbeatRequestGenerator(nextSend)
	}
	return &firstNotNullMessageGenerator{
		generators: []sasRequestGenerator{
			sas.NewGrantRequestGenerator(),
			sas.NewSpectrumInquiryRequestGenerator(),
		},
	}
}

type perCbsdMessageGenerator struct {
	sas  sasRequestGenerator
	crud crudActionGenerator
}

type sasRequestGenerator interface {
	GenerateRequests(*active_mode.ActiveModeConfig) []*sas.Request
}

type crudActionGenerator interface {
	GenerateActions(config *active_mode.ActiveModeConfig) []*crud.Action
}

type noMessageGenerator struct{}

func (*noMessageGenerator) GenerateRequests(_ *active_mode.ActiveModeConfig) []*sas.Request {
	return nil
}

type firstNotNullMessageGenerator struct {
	generators []sasRequestGenerator
}

func (f *firstNotNullMessageGenerator) GenerateRequests(config *active_mode.ActiveModeConfig) []*sas.Request {
	for _, g := range f.generators {
		messages := g.GenerateRequests(config)
		if len(messages) > 0 {
			return messages
		}
	}
	return nil
}

type noActionsGenerator struct{}

func (*noActionsGenerator) GenerateActions(_ *active_mode.ActiveModeConfig) []*crud.Action {
	return nil
}

type deleteActionGenerator struct{}

func (*deleteActionGenerator) GenerateActions(config *active_mode.ActiveModeConfig) []*crud.Action {
	return []*crud.Action{{
		Type:         crud.Delete,
		SerialNumber: config.GetCbsd().GetSerialNumber(),
	}}
}
