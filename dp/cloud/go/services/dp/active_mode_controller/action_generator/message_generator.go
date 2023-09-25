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

package action_generator

import (
	"time"

	"magma/dp/cloud/go/services/dp/active_mode_controller/action_generator/sas"
	"magma/dp/cloud/go/services/dp/storage"
)

type ActionGenerator struct {
	HeartbeatTimeout  time.Duration
	InactivityTimeout time.Duration
	Rng               RNG
}

type RNG interface {
	Int() int
}

func (a *ActionGenerator) GenerateActions(cbsds []*storage.DetailedCbsd, now time.Time) []Action {
	actions := make([]Action, 0, len(cbsds))
	for _, cbsd := range cbsds {
		g := a.getPerCbsdMessageGenerator(cbsd, now)
		actions = append(actions, g.generateActions(cbsd)...)
	}
	return actions
}

// TODO make this more readable
func (a *ActionGenerator) getPerCbsdMessageGenerator(cbsd *storage.DetailedCbsd, now time.Time) actionGeneratorPerCbsd {
	isActive := now.Sub(cbsd.Cbsd.LastSeen.Time) <= a.InactivityTimeout
	if cbsd.CbsdState.Name.String == "unregistered" {
		if cbsd.Cbsd.IsDeleted.Bool {
			return &deleteGenerator{}
		} else if cbsd.Cbsd.ShouldDeregister.Bool {
			return &acknowledgeDeregisterGenerator{}
		} else if isActive && cbsd.DesiredState.Name.String == "registered" {
			return &sasRequestGenerator{g: &sas.RegistrationRequestGenerator{}}
		}
	} else if cbsd.Cbsd.IsDeleted.Bool ||
		cbsd.Cbsd.ShouldDeregister.Bool ||
		cbsd.DesiredState.Name.String == "unregistered" {
		return &sasRequestGenerator{g: &sas.DeregistrationRequestGenerator{}}
	} else if cbsd.Cbsd.ShouldRelinquish.Bool {
		if len(cbsd.Grants) == 0 {
			return &acknowledgeRelinquishGenerator{}
		} else {
			return &sasRequestGenerator{g: &sas.RelinquishmentRequestGenerator{}}
		}
	} else if !isActive {
		return &sasRequestGenerator{g: &sas.RelinquishmentRequestGenerator{}}
	} else if len(cbsd.Cbsd.Channels) == 0 {
		return &sasRequestGenerator{g: &sas.SpectrumInquiryRequestGenerator{}}
	} else if len(cbsd.Cbsd.AvailableFrequencies) == 0 {
		return &storeAvailableFrequenciesGenerator{}
	} else {
		nextSend := now.Add(a.HeartbeatTimeout).Unix()
		gm := &grantManager{
			nextSendTimestamp: nextSend,
			rng:               a.Rng,
		}
		return &sasRequestGenerator{g: gm}
	}
	return &nothingGenerator{}
}
