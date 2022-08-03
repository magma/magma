/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package message_generator

import (
	"context"
	"fmt"
	"time"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/message"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas/eirp"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas/grant"
	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas_helpers"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type messageGenerator struct {
	heartbeatTimeout  time.Duration
	inactivityTimeout time.Duration
	rng               RNG
}

type RNG interface {
	Int() int
}

func NewMessageGenerator(
	heartbeatTimeout time.Duration,
	inactivityTimeout time.Duration,
	rng RNG,
) *messageGenerator {
	return &messageGenerator{
		heartbeatTimeout:  heartbeatTimeout,
		inactivityTimeout: inactivityTimeout,
		rng:               rng,
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
		msgs = append(msgs, g.action.generateActions(cbsd)...)
	}
	payloads := sas_helpers.Build(requests)
	for _, payload := range payloads {
		msgs = append(msgs, message.NewSasMessage(payload))
	}
	return msgs
}

func (m *messageGenerator) getPerCbsdMessageGenerator(cbsd *active_mode.Cbsd, now time.Time) *perCbsdMessageGenerator {
	generator := &perCbsdMessageGenerator{
		sas:    &noRequestGenerator{},
		action: &noMessageGenerator{},
	}
	isActive := cbsd.LastSeenTimestamp >= now.Add(-m.inactivityTimeout).Unix()
	if cbsd.State == active_mode.CbsdState_Unregistered {
		if cbsd.DbData.IsDeleted {
			generator.action = &deleteMessageGenerator{}
		} else if cbsd.DbData.ShouldDeregister {
			generator.action = &updateMessageGenerator{}
		} else if isActive && cbsd.DesiredState == active_mode.CbsdState_Registered {
			generator.sas = &sas.RegistrationRequestGenerator{}
		}
	} else if cbsd.DbData.IsDeleted ||
		cbsd.DbData.ShouldDeregister ||
		cbsd.DesiredState == active_mode.CbsdState_Unregistered {
		generator.sas = &sas.DeregistrationRequestGenerator{}
	} else if !isActive {
		generator.sas = &sas.RelinquishmentRequestGenerator{}
	} else if len(cbsd.Channels) == 0 {
		generator.sas = &sas.SpectrumInquiryRequestGenerator{}
	} else if len(cbsd.GrantSettings.AvailableFrequencies) == 0 {
		generator.action = &availableFrequenciesMessageGenerator{}
	} else {
		nextSend := now.Add(m.heartbeatTimeout).Unix()
		generator.sas = &grantManager{
			nextSendTimestamp: nextSend,
			rng:               m.rng,
		}
	}
	return generator
}

type grantManager struct {
	nextSendTimestamp int64
	rng               RNG
}

func (g *grantManager) GenerateRequests(cbsd *active_mode.Cbsd) []*sas.Request {
	grants := grant.GetFrequencyGrantMapping(cbsd.Grants)
	calc := eirp.NewCalculator(cbsd.InstallationParams.AntennaGainDbi, cbsd.EirpCapabilities)
	processors := grant.Processors[*sas.Request]{
		Del: &sas.RelinquishmentProcessor{
			CbsdId: cbsd.CbsdId,
			Grants: grants,
		},
		Keep: &sas.HeartbeatProcessor{
			NextSendTimestamp: g.nextSendTimestamp,
			CbsdId:            cbsd.CbsdId,
			Grants:            grants,
		},
		Add: &sas.GrantProcessor{
			CbsdId:   cbsd.CbsdId,
			Calc:     calc,
			Channels: cbsd.Channels,
		},
	}
	requests := grant.ProcessGrants(
		cbsd.Grants, cbsd.Preferences, cbsd.GrantSettings,
		processors, g.rng.Int(),
	)
	if len(requests) > 0 {
		return requests
	}
	gen := sas.SpectrumInquiryRequestGenerator{}
	return gen.GenerateRequests(cbsd)
}

type perCbsdMessageGenerator struct {
	sas    sasRequestGenerator
	action actionMessageGenerator
}

type sasRequestGenerator interface {
	GenerateRequests(cbsd *active_mode.Cbsd) []*sas.Request
}

type actionMessageGenerator interface {
	generateActions(cbsd *active_mode.Cbsd) []Message
}

type noRequestGenerator struct{}

func (*noRequestGenerator) GenerateRequests(_ *active_mode.Cbsd) []*sas.Request {
	return nil
}

type noMessageGenerator struct{}

func (*noMessageGenerator) generateActions(_ *active_mode.Cbsd) []Message {
	return nil
}

type deleteMessageGenerator struct{}

func (*deleteMessageGenerator) generateActions(data *active_mode.Cbsd) []Message {
	return []Message{message.NewDeleteMessage(data.DbData.Id)}
}

type updateMessageGenerator struct{}

func (*updateMessageGenerator) generateActions(data *active_mode.Cbsd) []Message {
	return []Message{message.NewUpdateMessage(data.DbData.Id)}
}

type availableFrequenciesMessageGenerator struct{}

func (*availableFrequenciesMessageGenerator) generateActions(data *active_mode.Cbsd) []Message {
	calc := eirp.NewCalculator(data.InstallationParams.AntennaGainDbi, data.EirpCapabilities)
	frequencies := grant.CalcAvailableFrequencies(data.Channels, calc)
	msg := message.NewStoreAvailableFrequenciesMessage(data.DbData.Id, frequencies)
	return []Message{msg}
}
