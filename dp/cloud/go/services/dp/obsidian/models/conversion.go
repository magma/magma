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
	"github.com/go-openapi/strfmt"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
)

func CbsdToBackend(m *MutableCbsd) *protos.CbsdData {
	return &protos.CbsdData{
		UserId:       m.UserID,
		FccId:        m.FccID,
		SerialNumber: m.SerialNumber,
		Capabilities: &protos.Capabilities{
			AntennaGain:      *m.Capabilities.AntennaGain,
			MaxPower:         *m.Capabilities.MaxPower,
			MinPower:         *m.Capabilities.MinPower,
			NumberOfAntennas: m.Capabilities.NumberOfAntennas,
		},
		Preferences: &protos.FrequencyPreferences{
			BandwidthMhz:   m.FrequencyPreferences.BandwidthMhz,
			FrequenciesMhz: m.FrequencyPreferences.FrequenciesMhz,
		},
	}
}

func CbsdFromBackend(details *protos.CbsdDetails) *Cbsd {
	return &Cbsd{
		Capabilities: Capabilities{
			AntennaGain:      &details.Data.Capabilities.AntennaGain,
			MaxPower:         &details.Data.Capabilities.MaxPower,
			MinPower:         &details.Data.Capabilities.MinPower,
			NumberOfAntennas: details.Data.Capabilities.NumberOfAntennas,
		},
		FrequencyPreferences: FrequencyPreferences{
			BandwidthMhz:   details.Data.Preferences.BandwidthMhz,
			FrequenciesMhz: makeSliceNotNil(details.Data.Preferences.FrequenciesMhz),
		},
		CbsdID:       details.CbsdId,
		FccID:        details.Data.FccId,
		Grant:        getGrant(details.Grant),
		ID:           details.Id,
		IsActive:     details.IsActive,
		SerialNumber: details.Data.SerialNumber,
		State:        details.State,
		UserID:       details.Data.UserId,
	}
}

func makeSliceNotNil(s []int64) []int64 {
	if len(s) == 0 {
		return []int64{}
	}
	return s
}

func getGrant(grant *protos.GrantDetails) *Grant {
	if grant == nil {
		return nil
	}
	return &Grant{
		BandwidthMhz:       grant.BandwidthMhz,
		FrequencyMhz:       grant.FrequencyMhz,
		GrantExpireTime:    to_pointer.TimeToDateTime(grant.GrantExpireTimestamp),
		MaxEirp:            grant.MaxEirp,
		State:              grant.State,
		TransmitExpireTime: to_pointer.TimeToDateTime(grant.TransmitExpireTimestamp),
	}
}

type LogInterface struct {
	Body         string          `json:"log_message"`
	FccID        string          `json:"fcc_id"`
	From         string          `json:"log_from"`
	SerialNumber string          `json:"cbsd_serial_number"`
	Time         strfmt.DateTime `json:"@timestamp"`
	To           string          `json:"log_to"`
	Type         string          `json:"log_name"`
}

func LogInterfaceToLog(i *LogInterface) *Log {
	return &Log{
		Body:         i.Body,
		FccID:        i.FccID,
		From:         i.From,
		SerialNumber: i.SerialNumber,
		Time:         i.Time,
		To:           i.To,
		Type:         i.Type,
	}
}
