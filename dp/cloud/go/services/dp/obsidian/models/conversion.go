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
	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
)

func MessageFromBackend(details *protos.Log) *Message {
	return &Message{
		From:         details.From,
		To:           details.To,
		Body:         details.Message,
		FccID:        details.FccId,
		SerialNumber: details.SerialNumber,
		Time:         *to_pointer.TimeMilliToDate(details.TimestampMilli),
		Type:         details.Name,
	}
}

func CbsdToBackend(m *MutableCbsd) *protos.CbsdData {
	return &protos.CbsdData{
		UserId:       *m.UserID,
		FccId:        *m.FccID,
		SerialNumber: *m.SerialNumber,
		Capabilities: &protos.Capabilities{
			AntennaGain:      *m.Capabilities.AntennaGain,
			MaxPower:         *m.Capabilities.MaxPower,
			MinPower:         *m.Capabilities.MinPower,
			NumberOfAntennas: *m.Capabilities.NumberOfAntennas,
		},
	}
}

func CbsdFromBackend(details *protos.CbsdDetails) *Cbsd {
	return &Cbsd{
		Capabilities: &Capabilities{
			AntennaGain:      &details.Data.Capabilities.AntennaGain,
			MaxPower:         &details.Data.Capabilities.MaxPower,
			MinPower:         &details.Data.Capabilities.MinPower,
			NumberOfAntennas: &details.Data.Capabilities.NumberOfAntennas,
		},
		CbsdID:       details.CbsdId,
		FccID:        &details.Data.FccId,
		Grant:        getGrant(details),
		ID:           details.Id,
		IsActive:     details.IsActive,
		SerialNumber: &details.Data.SerialNumber,
		State:        details.State,
		UserID:       &details.Data.UserId,
	}
}

func getGrant(details *protos.CbsdDetails) *Grant {
	grant := details.Grant
	if grant == nil {
		return nil
	}

	return &Grant{
		BandwidthMhz:       grant.BandwidthMhz,
		FrequencyMhz:       grant.FrequencyMhz,
		GrantExpireTime:    *to_pointer.TimeToDateTime(grant.GrantExpireTimestamp),
		MaxEirp:            &grant.MaxEirp,
		State:              grant.State,
		TransmitExpireTime: *to_pointer.TimeToDateTime(grant.TransmitExpireTimestamp),
	}
}
