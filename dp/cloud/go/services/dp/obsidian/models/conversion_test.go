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

	b "magma/dp/cloud/go/services/dp/builders"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
)

func TestCbsdToBackend(t *testing.T) {
	cbsd := b.NewMutableCbsdModelPayloadBuilder().Payload
	data, _ := models.CbsdToBackend(cbsd)
	assert.Equal(t, data.UserId, cbsd.UserID)
	assert.Equal(t, data.FccId, cbsd.FccID)
	assert.Equal(t, data.SerialNumber, cbsd.SerialNumber)
	assert.Equal(t, data.Capabilities.MaxPower, *cbsd.Capabilities.MaxPower)
	assert.Equal(t, data.Capabilities.MinPower, *cbsd.Capabilities.MinPower)
	assert.Equal(t, data.Capabilities.NumberOfAntennas, cbsd.Capabilities.NumberOfAntennas)
	assert.Equal(t, data.Capabilities.MaxIbwMhz, cbsd.Capabilities.MaxIbwMhz)
	assert.Equal(t, data.Preferences.BandwidthMhz, cbsd.FrequencyPreferences.BandwidthMhz)
	assert.Equal(t, data.Preferences.FrequenciesMhz, cbsd.FrequencyPreferences.FrequenciesMhz)
	assert.Equal(t, data.DesiredState, cbsd.DesiredState)
	assert.Equal(t, data.SingleStepEnabled, *cbsd.SingleStepEnabled)
	assert.Equal(t, data.CbsdCategory, cbsd.CbsdCategory)
	assert.Equal(t, data.CarrierAggregationEnabled, *cbsd.CarrierAggregationEnabled)
	assert.Equal(t, data.GrantRedundancy, *cbsd.GrantRedundancy)
}

func TestCbsdFromBackendWithoutGrants(t *testing.T) {
	details := b.NewDetailedProtoCbsdBuilder(
		b.NewCbsdProtoPayloadBuilder()).Details
	data := models.CbsdFromBackend(details)
	assert.Empty(t, data.Grants)
	assert.Equal(t, data.ID, details.Id)
	assert.Equal(t, data.IsActive, details.IsActive)
	assert.Equal(t, data.CbsdID, details.CbsdId)
	assert.Equal(t, data.UserID, details.Data.UserId)
	assert.Equal(t, data.FccID, details.Data.FccId)
	assert.Equal(t, data.SerialNumber, details.Data.SerialNumber)
	assert.Equal(t, *data.Capabilities.MinPower, details.Data.Capabilities.MinPower)
	assert.Equal(t, *data.Capabilities.MaxPower, details.Data.Capabilities.MaxPower)
	assert.Equal(t, data.Capabilities.NumberOfAntennas, details.Data.Capabilities.NumberOfAntennas)
	assert.Equal(t, data.Capabilities.MaxIbwMhz, details.Data.Capabilities.MaxIbwMhz)
	assert.Equal(t, data.FrequencyPreferences.BandwidthMhz, details.Data.Preferences.BandwidthMhz)
	assert.Equal(t, data.FrequencyPreferences.FrequenciesMhz, details.Data.Preferences.FrequenciesMhz)
	assert.Equal(t, data.DesiredState, details.Data.DesiredState)
	assert.Equal(t, data.SingleStepEnabled, details.Data.SingleStepEnabled)
	assert.Equal(t, data.CbsdCategory, details.Data.CbsdCategory)
	assert.Equal(t, data.CarrierAggregationEnabled, details.Data.CarrierAggregationEnabled)
	assert.Equal(t, data.GrantRedundancy, details.Data.GrantRedundancy)
}

func TestCbsdFromBackendWithGrants(t *testing.T) {
	details := b.NewDetailedProtoCbsdBuilder(
		b.NewCbsdProtoPayloadBuilder()).WithGrant().Details
	data := models.CbsdFromBackend(details)
	assert.Len(t, data.Grants, len(details.Grants))
	expected, actual := data.Grants[0], details.Grants[0]
	assert.Equal(t, expected.BandwidthMhz, actual.BandwidthMhz)
	assert.Equal(t, expected.FrequencyMhz, actual.FrequencyMhz)
	assert.Equal(t, expected.GrantExpireTime, to_pointer.TimeToDateTime(actual.GrantExpireTimestamp))
	assert.Equal(t, expected.TransmitExpireTime, to_pointer.TimeToDateTime(actual.TransmitExpireTimestamp))
	assert.Equal(t, expected.MaxEirp, actual.MaxEirp)
	assert.Equal(t, expected.State, actual.State)
}

func TestCbsdFromBackendWithEmptyInstallationParam(t *testing.T) {
	details := b.NewDetailedProtoCbsdBuilder(
		b.NewCbsdProtoPayloadBuilder().WithEmptyInstallationParam()).Details
	data := models.CbsdFromBackend(details)
	assert.Nil(t, data.InstallationParam.LatitudeDeg)
	assert.Nil(t, data.InstallationParam.LongitudeDeg)
	assert.Nil(t, data.InstallationParam.IndoorDeployment)
	assert.Nil(t, data.InstallationParam.Heightm)
	assert.Nil(t, data.InstallationParam.HeightType)
	assert.Nil(t, data.InstallationParam.AntennaGain)
}

func TestCbsdFromBackendWithoutInstallationParam(t *testing.T) {
	details := b.NewDetailedProtoCbsdBuilder(
		b.NewCbsdProtoPayloadBuilder()).Details
	data := models.CbsdFromBackend(details)
	assert.Empty(t, data.InstallationParam)
}

func TestCbsdFromBackendWithInstallationParam(t *testing.T) {
	details := b.NewDetailedProtoCbsdBuilder(
		b.NewCbsdProtoPayloadBuilder().WithFullInstallationParam()).Details
	data := models.CbsdFromBackend(details)
	assert.Equal(t, data.InstallationParam.LatitudeDeg, &details.Data.InstallationParam.LatitudeDeg.Value)
	assert.Equal(t, data.InstallationParam.LongitudeDeg, &details.Data.InstallationParam.LongitudeDeg.Value)
	assert.Equal(t, data.InstallationParam.IndoorDeployment, &details.Data.InstallationParam.IndoorDeployment.Value)
	assert.Equal(t, data.InstallationParam.Heightm, &details.Data.InstallationParam.HeightM.Value)
	assert.Equal(t, data.InstallationParam.HeightType, &details.Data.InstallationParam.HeightType.Value)
	assert.Equal(t, data.InstallationParam.AntennaGain, &details.Data.InstallationParam.AntennaGain.Value)
}

func TestCbsdFromBackendWithEmptyFrequencies(t *testing.T) {
	details := b.NewDetailedProtoCbsdBuilder(
		b.NewCbsdProtoPayloadBuilder().WithEmptyPreferences()).Details
	data := models.CbsdFromBackend(details)
	assert.Equal(t, []int64{}, data.FrequencyPreferences.FrequenciesMhz)
	assert.Equal(t, int64(0), data.FrequencyPreferences.BandwidthMhz)
}
