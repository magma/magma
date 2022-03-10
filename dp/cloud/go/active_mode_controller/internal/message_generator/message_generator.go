package message_generator

import (
	"context"
	"fmt"
	"time"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/message"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas_helpers"
	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type messageGenerator struct {
	heartbeatTimeout  time.Duration
	inactivityTimeout time.Duration
	indexProvider     ranges.IndexProvider
}

func NewMessageGenerator(
	heartbeatTimeout time.Duration,
	inactivityTimeout time.Duration,
	indexProvider ranges.IndexProvider,
) *messageGenerator {
	return &messageGenerator{
		heartbeatTimeout:  heartbeatTimeout,
		inactivityTimeout: inactivityTimeout,
		indexProvider:     indexProvider,
	}
}

type Message interface {
	fmt.Stringer
	Send(context.Context, message.ClientProvider) error
}

func (m *messageGenerator) GenerateMessages(state *active_mode.State, now time.Time) []Message {
	var requests []*sas.Request
	var msgs []Message
	for _, cbsd := range state.Cbsds {
		g := m.getPerCbsdMessageGenerator(cbsd, now)
		reqs := g.sas.GenerateRequests(cbsd)
		requests = append(requests, reqs...)
		msgs = append(msgs, g.crud.generateMessages(cbsd.GetDbData())...)
	}
	payloads := sas_helpers.Build(requests)
	for _, payload := range payloads {
		msgs = append(msgs, message.NewSasMessage(payload))
	}
	return msgs
}

func (m *messageGenerator) getPerCbsdMessageGenerator(cbsd *active_mode.Cbsd, now time.Time) *perCbsdMessageGenerator {
	if cbsd.GetState() == active_mode.CbsdState_Unregistered {
		data := cbsd.GetDbData()
		if data.GetIsDeleted() {
			return &perCbsdMessageGenerator{
				sas:  &noRequestGenerator{},
				crud: &deleteMessageGenerator{},
			}
		}
		if data.GetIsUpdated() {
			return &perCbsdMessageGenerator{
				sas:  &noRequestGenerator{},
				crud: &updateMessageGenerator{},
			}
		}
	}
	return &perCbsdMessageGenerator{
		sas:  m.getPerCbsdSasRequestGenerator(cbsd, now),
		crud: &noMessageGenerator{},
	}
}

func (m *messageGenerator) getPerCbsdSasRequestGenerator(cbsd *active_mode.Cbsd, now time.Time) sasRequestGenerator {
	if requiresDeregistration(cbsd.GetDbData()) && cbsd.GetState() == active_mode.CbsdState_Registered {
		return sas.NewDeregistrationRequestGenerator()
	}
	if cbsd.GetDesiredState() == active_mode.CbsdState_Unregistered {
		if cbsd.GetState() == active_mode.CbsdState_Registered {
			return sas.NewDeregistrationRequestGenerator()
		}
		return &noRequestGenerator{}
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
	return &firstNotNullRequestGenerator{
		generators: []sasRequestGenerator{
			sas.NewGrantRequestGenerator(m.indexProvider),
			sas.NewSpectrumInquiryRequestGenerator(),
		},
	}
}

func requiresDeregistration(data *active_mode.DatabaseCbsd) bool {
	return data.GetIsDeleted() || data.GetIsUpdated()
}

type perCbsdMessageGenerator struct {
	sas  sasRequestGenerator
	crud crudMessageGenerator
}

type sasRequestGenerator interface {
	GenerateRequests(cbsd *active_mode.Cbsd) []*sas.Request
}

type crudMessageGenerator interface {
	generateMessages(*active_mode.DatabaseCbsd) []Message
}

type noRequestGenerator struct{}

func (*noRequestGenerator) GenerateRequests(_ *active_mode.Cbsd) []*sas.Request {
	return nil
}

type firstNotNullRequestGenerator struct {
	generators []sasRequestGenerator
}

func (f *firstNotNullRequestGenerator) GenerateRequests(config *active_mode.Cbsd) []*sas.Request {
	for _, g := range f.generators {
		messages := g.GenerateRequests(config)
		if len(messages) > 0 {
			return messages
		}
	}
	return nil
}

type noMessageGenerator struct{}

func (*noMessageGenerator) generateMessages(_ *active_mode.DatabaseCbsd) []Message {
	return nil
}

type deleteMessageGenerator struct{}

func (*deleteMessageGenerator) generateMessages(data *active_mode.DatabaseCbsd) []Message {
	return []Message{message.NewDeleteMessage(data.GetId())}
}

type updateMessageGenerator struct{}

func (*updateMessageGenerator) generateMessages(data *active_mode.DatabaseCbsd) []Message {
	return []Message{message.NewUpdateMessage(data.GetId())}
}
