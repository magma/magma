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

package models

import (
	"time"

	"github.com/go-openapi/strfmt"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
)

func CbsdToBackend(m *MutableCbsd) (*protos.CbsdData, error) {
	return &protos.CbsdData{
		UserId:                    m.UserID,
		FccId:                     m.FccID,
		SerialNumber:              m.SerialNumber,
		SingleStepEnabled:         *m.SingleStepEnabled,
		CbsdCategory:              m.CbsdCategory,
		CarrierAggregationEnabled: *m.CarrierAggregationEnabled,
		GrantRedundancy:           *m.GrantRedundancy,
		Capabilities: &protos.Capabilities{
			MaxPower:         *m.Capabilities.MaxPower,
			MinPower:         *m.Capabilities.MinPower,
			NumberOfAntennas: m.Capabilities.NumberOfAntennas,
			MaxIbwMhz:        m.Capabilities.MaxIbwMhz,
		},
		Preferences: &protos.FrequencyPreferences{
			BandwidthMhz:   m.FrequencyPreferences.BandwidthMhz,
			FrequenciesMhz: m.FrequencyPreferences.FrequenciesMhz,
		},
		DesiredState: m.DesiredState,
		InstallationParam: &protos.InstallationParam{
			AntennaGain: to_pointer.FloatToDoubleValue(m.InstallationParam.AntennaGain),
		},
	}, nil
}

func CbsdFromBackend(details *protos.CbsdDetails) *Cbsd {
	return &Cbsd{
		Capabilities: Capabilities{
			MaxPower:         &details.Data.Capabilities.MaxPower,
			MinPower:         &details.Data.Capabilities.MinPower,
			NumberOfAntennas: details.Data.Capabilities.NumberOfAntennas,
			MaxIbwMhz:        details.Data.Capabilities.MaxIbwMhz,
		},
		FrequencyPreferences: FrequencyPreferences{
			BandwidthMhz:   details.Data.Preferences.BandwidthMhz,
			FrequenciesMhz: makeSliceNotNil(details.Data.Preferences.FrequenciesMhz),
		},
		GrantRedundancy:           details.Data.GrantRedundancy,
		CarrierAggregationEnabled: details.Data.CarrierAggregationEnabled,
		CbsdID:                    details.CbsdId,
		DesiredState:              details.Data.DesiredState,
		FccID:                     details.Data.FccId,
		Grants:                    getGrants(details.Grants),
		ID:                        details.Id,
		IsActive:                  details.IsActive,
		SerialNumber:              details.Data.SerialNumber,
		State:                     details.State,
		UserID:                    details.Data.UserId,
		SingleStepEnabled:         details.Data.SingleStepEnabled,
		CbsdCategory:              details.Data.CbsdCategory,
		InstallationParam:         getModelInstallationParam(details.Data.InstallationParam),
	}
}

func makeSliceNotNil(s []int64) []int64 {
	if len(s) == 0 {
		return []int64{}
	}
	return s
}

func getGrants(grants []*protos.GrantDetails) []*Grant {
	res := make([]*Grant, len(grants))
	for i, g := range grants {
		res[i] = &Grant{
			BandwidthMhz:       g.BandwidthMhz,
			FrequencyMhz:       g.FrequencyMhz,
			GrantExpireTime:    to_pointer.TimeToDateTime(g.GrantExpireTimestamp),
			MaxEirp:            g.MaxEirp,
			State:              g.State,
			TransmitExpireTime: to_pointer.TimeToDateTime(g.TransmitExpireTimestamp),
		}
	}
	return res
}

func getModelInstallationParam(params *protos.InstallationParam) InstallationParam {
	return InstallationParam{
		AntennaGain:      to_pointer.DoubleValueToFloat(params.AntennaGain),
		Heightm:          to_pointer.DoubleValueToFloat(params.HeightM),
		HeightType:       to_pointer.StringValueToString(params.HeightType),
		IndoorDeployment: to_pointer.BoolValueToBool(params.IndoorDeployment),
		LatitudeDeg:      to_pointer.DoubleValueToFloat(params.LatitudeDeg),
		LongitudeDeg:     to_pointer.DoubleValueToFloat(params.LongitudeDeg),
	}
}

type LogInterface struct {
	Body         string `json:"log_message"`
	FccID        string `json:"fcc_id"`
	From         string `json:"log_from"`
	SerialNumber string `json:"cbsd_serial_number"`
	Time         int64  `json:"event_timestamp"`
	To           string `json:"log_to"`
	Type         string `json:"log_name"`
}

func LogInterfaceToLog(i *LogInterface) *Log {
	return &Log{
		Body:         i.Body,
		FccID:        i.FccID,
		From:         i.From,
		SerialNumber: i.SerialNumber,
		Time:         strfmt.DateTime(time.Unix(i.Time, 0).UTC()),
		To:           i.To,
		Type:         i.Type,
	}
}
