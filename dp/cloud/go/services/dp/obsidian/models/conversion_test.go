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

func TestMessageFromBackend(t *testing.T) {
	details := &protos.Log{
		From:           "CBSD",
		To:             "DP",
		Name:           "some log",
		Message:        "some msg",
		SerialNumber:   "123456",
		FccId:          "some fcc id",
		TimestampMilli: 123,
	}
	data := models.MessageFromBackend(details)
	assert.Equal(t, data.From, details.From)
	assert.Equal(t, data.To, details.To)
	assert.Equal(t, data.Type, details.Name)
	assert.Equal(t, data.Body, details.Message)
	assert.Equal(t, data.SerialNumber, details.SerialNumber)
	assert.Equal(t, data.FccID, details.FccId)
	assert.Equal(t, data.Time, *to_pointer.TimeMilliToDate(details.TimestampMilli))
}

func TestCbsdToBackend(t *testing.T) {
	cbsd := models.Cbsd{
		Capabilities: &models.Capabilities{
			AntennaGain:      to_pointer.Float(5),
			MaxPower:         to_pointer.Float(24),
			MinPower:         to_pointer.Float(0),
			NumberOfAntennas: to_pointer.Int(1),
		},
		FccID:        to_pointer.Str("barID"),
		IsActive:     to_pointer.Bool(false),
		SerialNumber: to_pointer.Str("12345"),
		UserID:       to_pointer.Str("fooUser"),
	}
	data := models.CbsdToBackend(&cbsd)
	assert.Equal(t, data.UserId, *cbsd.UserID)
	assert.Equal(t, data.FccId, *cbsd.FccID)
	assert.Equal(t, data.SerialNumber, *cbsd.SerialNumber)
	assert.Equal(t, data.Capabilities.AntennaGain, *cbsd.Capabilities.AntennaGain)
	assert.Equal(t, data.Capabilities.MaxPower, *cbsd.Capabilities.MaxPower)
	assert.Equal(t, data.Capabilities.MinPower, *cbsd.Capabilities.MinPower)
	assert.Equal(t, data.Capabilities.NumberOfAntennas, *cbsd.Capabilities.NumberOfAntennas)
}

func TestCbsdFromBackendWithoutGrant(t *testing.T) {
	details := getCbsdDetails(false)
	data := models.CbsdFromBackend(details)
	assert.Nil(t, data.Grant)
	assert.Equal(t, data.ID, details.Id)
	assert.Equal(t, data.CbsdID, details.CbsdId)
	assert.Equal(t, *data.UserID, details.Data.UserId)
	assert.Equal(t, *data.FccID, details.Data.FccId)
	assert.Equal(t, *data.SerialNumber, details.Data.SerialNumber)
	assert.Equal(t, *data.Capabilities.MinPower, details.Data.Capabilities.MinPower)
	assert.Equal(t, *data.Capabilities.MaxPower, details.Data.Capabilities.MaxPower)
	assert.Equal(t, *data.Capabilities.NumberOfAntennas, details.Data.Capabilities.NumberOfAntennas)
	assert.Equal(t, *data.Capabilities.AntennaGain, details.Data.Capabilities.AntennaGain)
}

func TestCbsdFromBackendWithGrant(t *testing.T) {
	details := getCbsdDetails(true)
	data := models.CbsdFromBackend(details)
	assert.Equal(t, data.Grant.BandwidthMhz, details.Grant.BandwidthMhz)
	assert.Equal(t, data.Grant.FrequencyMhz, details.Grant.FrequencyMhz)
	assert.Equal(t, data.Grant.GrantExpireTime, *to_pointer.TimeToDateTime(details.Grant.GrantExpireTimestamp))
	assert.Equal(t, data.Grant.TransmitExpireTime, *to_pointer.TimeToDateTime(details.Grant.TransmitExpireTimestamp))
	assert.Equal(t, *data.Grant.MaxEirp, details.Grant.MaxEirp)
	assert.Equal(t, data.Grant.State, details.Grant.State)
}

func getCbsdDetails(withGrant bool) *protos.CbsdDetails {
	details := protos.CbsdDetails{
		Id: 1,
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
		},
		CbsdId: "someCbsdId",
	}
	if withGrant {
		details.Grant = &protos.GrantDetails{
			BandwidthMhz:            10,
			FrequencyMhz:            12345,
			MaxEirp:                 123,
			State:                   "someState",
			TransmitExpireTimestamp: 12345678,
			GrantExpireTimestamp:    12345678,
		}
	}
	return &details
}
