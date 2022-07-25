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

package builders

import "magma/dp/cloud/go/active_mode_controller/protos/active_mode"

const (
	Now     = 1000
	DbId    = 123
	CbsdId  = "some_cbsd_id"
	GrantId = "some_grant_id"
)

var (
	SomeChannel = &active_mode.Channel{
		LowFrequencyHz:  3550e6,
		HighFrequencyHz: 3700e6,
	}
	NoAvailableFrequencies = []uint32{0, 0, 0, 0}
)

type cbsdBuilder struct {
	cbsd *active_mode.Cbsd
}

func NewCbsdBuilder() *cbsdBuilder {
	return &cbsdBuilder{
		cbsd: &active_mode.Cbsd{
			CbsdId:            CbsdId,
			State:             active_mode.CbsdState_Registered,
			DesiredState:      active_mode.CbsdState_Registered,
			LastSeenTimestamp: Now,
			SasSettings: &active_mode.SasSettings{
				SingleStepEnabled: false,
				CbsdCategory:      "A",
				SerialNumber:      "some_serial_number",
				FccId:             "some_fcc_id",
				UserId:            "some_user_id",
			},
			InstallationParams: &active_mode.InstallationParams{
				AntennaGainDbi: 15,
			},
			EirpCapabilities: &active_mode.EirpCapabilities{
				MinPower:      0,
				MaxPower:      30,
				NumberOfPorts: 1,
			},
			DbData: &active_mode.DatabaseCbsd{
				Id: DbId,
			},
			Preferences: &active_mode.FrequencyPreferences{
				BandwidthMhz: 20,
			},
			GrantSettings: &active_mode.GrantSettings{
				MaxIbwMhz: 150,
			},
		},
	}
}

func (c *cbsdBuilder) Build() *active_mode.Cbsd {
	return c.cbsd
}

func (c *cbsdBuilder) Inactive() *cbsdBuilder {
	c.cbsd.LastSeenTimestamp = 0
	return c
}

func (c *cbsdBuilder) WithState(state active_mode.CbsdState) *cbsdBuilder {
	c.cbsd.State = state
	return c
}

func (c *cbsdBuilder) WithDesiredState(state active_mode.CbsdState) *cbsdBuilder {
	c.cbsd.DesiredState = state
	return c
}

func (c *cbsdBuilder) Deleted() *cbsdBuilder {
	c.cbsd.DbData.IsDeleted = true
	return c
}

func (c *cbsdBuilder) ForDeregistration() *cbsdBuilder {
	c.cbsd.DbData.ShouldDeregister = true
	return c
}

func (c *cbsdBuilder) WithChannel(channel *active_mode.Channel) *cbsdBuilder {
	c.cbsd.Channels = append(c.cbsd.Channels, channel)
	return c
}

func (c *cbsdBuilder) WithGrant(grant *active_mode.Grant) *cbsdBuilder {
	c.cbsd.Grants = append(c.cbsd.Grants, grant)
	return c
}

func (c *cbsdBuilder) WithAvailableFrequencies(frequencies []uint32) *cbsdBuilder {
	c.cbsd.GrantSettings.AvailableFrequencies = frequencies
	return c
}

func (c *cbsdBuilder) WithCarrierAggregation() *cbsdBuilder {
	c.cbsd.GrantSettings.GrantRedundancyEnabled = true
	c.cbsd.GrantSettings.CarrierAggregationEnabled = true
	return c
}

func (c *cbsdBuilder) WithName(name string) *cbsdBuilder {
	c.cbsd.CbsdId = name
	c.cbsd.SasSettings.SerialNumber = name
	c.cbsd.SasSettings.FccId = name
	c.cbsd.SasSettings.UserId = name
	return c
}
