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

package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
)

func TestCbsdToBackend(t *testing.T) {
	cbsd := models.MutableCbsd{
		Capabilities: models.Capabilities{
			AntennaGain:      to_pointer.Float(5),
			MaxPower:         to_pointer.Float(24),
			MinPower:         to_pointer.Float(0),
			NumberOfAntennas: 1,
		},
		FrequencyPreferences: models.FrequencyPreferences{
			BandwidthMhz:   10,
			FrequenciesMhz: []int64{3600},
		},
		FccID:        "barID",
		SerialNumber: "12345",
		UserID:       "fooUser",
	}
	data := models.CbsdToBackend(&cbsd)
	assert.Equal(t, data.UserId, cbsd.UserID)
	assert.Equal(t, data.FccId, cbsd.FccID)
	assert.Equal(t, data.SerialNumber, cbsd.SerialNumber)
	assert.Equal(t, data.Capabilities.AntennaGain, *cbsd.Capabilities.AntennaGain)
	assert.Equal(t, data.Capabilities.MaxPower, *cbsd.Capabilities.MaxPower)
	assert.Equal(t, data.Capabilities.MinPower, *cbsd.Capabilities.MinPower)
	assert.Equal(t, data.Capabilities.NumberOfAntennas, cbsd.Capabilities.NumberOfAntennas)
	assert.Equal(t, data.Preferences.BandwidthMhz, cbsd.FrequencyPreferences.BandwidthMhz)
	assert.Equal(t, data.Preferences.FrequenciesMhz, cbsd.FrequencyPreferences.FrequenciesMhz)
}

func TestCbsdFromBackendWithoutGrant(t *testing.T) {
	details := getCbsdDetails()
	data := models.CbsdFromBackend(details)
	assert.Nil(t, data.Grant)
	assert.Equal(t, data.ID, details.Id)
	assert.Equal(t, data.IsActive, details.IsActive)
	assert.Equal(t, data.CbsdID, details.CbsdId)
	assert.Equal(t, data.UserID, details.Data.UserId)
	assert.Equal(t, data.FccID, details.Data.FccId)
	assert.Equal(t, data.SerialNumber, details.Data.SerialNumber)
	assert.Equal(t, *data.Capabilities.MinPower, details.Data.Capabilities.MinPower)
	assert.Equal(t, *data.Capabilities.MaxPower, details.Data.Capabilities.MaxPower)
	assert.Equal(t, data.Capabilities.NumberOfAntennas, details.Data.Capabilities.NumberOfAntennas)
	assert.Equal(t, *data.Capabilities.AntennaGain, details.Data.Capabilities.AntennaGain)
	assert.Equal(t, data.FrequencyPreferences.BandwidthMhz, details.Data.Preferences.BandwidthMhz)
	assert.Equal(t, data.FrequencyPreferences.FrequenciesMhz, details.Data.Preferences.FrequenciesMhz)
}

func TestCbsdFromBackendWithGrant(t *testing.T) {
	details := getCbsdDetails()
	details.Grant = getGrant()
	data := models.CbsdFromBackend(details)
	assert.Equal(t, data.Grant.BandwidthMhz, details.Grant.BandwidthMhz)
	assert.Equal(t, data.Grant.FrequencyMhz, details.Grant.FrequencyMhz)
	assert.Equal(t, data.Grant.GrantExpireTime, to_pointer.TimeToDateTime(details.Grant.GrantExpireTimestamp))
	assert.Equal(t, data.Grant.TransmitExpireTime, to_pointer.TimeToDateTime(details.Grant.TransmitExpireTimestamp))
	assert.Equal(t, data.Grant.MaxEirp, details.Grant.MaxEirp)
	assert.Equal(t, data.Grant.State, details.Grant.State)
}

func TestCbsdFromBackendWithEmptyFrequencies(t *testing.T) {
	details := getCbsdDetails()
	details.Data.Preferences.FrequenciesMhz = nil
	data := models.CbsdFromBackend(details)
	assert.Equal(t, []int64{}, data.FrequencyPreferences.FrequenciesMhz)
}

func getCbsdDetails() *protos.CbsdDetails {
	return &protos.CbsdDetails{
		Id:       1,
		IsActive: false,
		Data: &protos.CbsdData{
			UserId:       "barId",
			FccId:        "bazId",
			SerialNumber: "12345",
			Capabilities: &protos.Capabilities{
				MinPower:         0,
				MaxPower:         24,
				NumberOfAntennas: 1,
				AntennaGain:      5,
			},
			Preferences: &protos.FrequencyPreferences{
				BandwidthMhz:   20,
				FrequenciesMhz: []int64{3600},
			},
		},
		CbsdId: "someCbsdId",
	}
}

func getGrant() *protos.GrantDetails {
	return &protos.GrantDetails{
		BandwidthMhz:            10,
		FrequencyMhz:            12345,
		MaxEirp:                 123,
		State:                   "someState",
		TransmitExpireTimestamp: 12345678,
		GrantExpireTimestamp:    12345678,
	}
}
