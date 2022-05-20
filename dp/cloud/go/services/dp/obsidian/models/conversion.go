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
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
)

func CbsdToBackend(m *MutableCbsd) *protos.CbsdData {
	return &protos.CbsdData{
		UserId:            m.UserID,
		FccId:             m.FccID,
		SerialNumber:      m.SerialNumber,
		SingleStepEnabled: *m.SingleStepEnabled,
		CbsdCategory:      m.CbsdCategory,
		Capabilities: &protos.Capabilities{
			MaxPower:         *m.Capabilities.MaxPower,
			MinPower:         *m.Capabilities.MinPower,
			NumberOfAntennas: m.Capabilities.NumberOfAntennas,
		},
		Preferences: &protos.FrequencyPreferences{
			BandwidthMhz:   m.FrequencyPreferences.BandwidthMhz,
			FrequenciesMhz: m.FrequencyPreferences.FrequenciesMhz,
		},
		DesiredState:      m.DesiredState,
		InstallationParam: getProtoInstallationParam(m.InstallationParam),
	}
}

func CbsdFromBackend(details *protos.CbsdDetails) *Cbsd {
	return &Cbsd{
		Capabilities: &Capabilities{
			MaxPower:         &details.Data.Capabilities.MaxPower,
			MinPower:         &details.Data.Capabilities.MinPower,
			NumberOfAntennas: details.Data.Capabilities.NumberOfAntennas,
		},
		FrequencyPreferences: FrequencyPreferences{
			BandwidthMhz:   details.Data.Preferences.BandwidthMhz,
			FrequenciesMhz: makeSliceNotNil(details.Data.Preferences.FrequenciesMhz),
		},
		CbsdID:            details.CbsdId,
		DesiredState:      details.Data.DesiredState,
		FccID:             details.Data.FccId,
		Grant:             getGrant(details.Grant),
		ID:                details.Id,
		IsActive:          details.IsActive,
		SerialNumber:      details.Data.SerialNumber,
		State:             details.State,
		UserID:            details.Data.UserId,
		SingleStepEnabled: details.Data.SingleStepEnabled,
		CbsdCategory:      details.Data.CbsdCategory,
		InstallationParam: getModelInstallationParam(details.Data.GetInstallationParam()),
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

func getModelInstallationParam(params *protos.InstallationParam) *InstallationParam {
	if params == nil {
		return nil
	}
	return &InstallationParam{
		AntennaGain:      doubleValueToFloatOrNil(params.AntennaGain),
		Heightm:          doubleValueToFloatOrNil(params.HeightM),
		HeightType:       stringValueToStringOrNil(params.HeightType),
		IndoorDeployment: boolValueToBoolOrNil(params.IndoorDeployment),
		LatitudeDeg:      doubleValueToFloatOrNil(params.LatitudeDeg),
		LongitudeDeg:     doubleValueToFloatOrNil(params.LongitudeDeg),
	}
}

func getProtoInstallationParam(params *InstallationParam) *protos.InstallationParam {
	if params == nil {
		return nil
	}
	return &protos.InstallationParam{
		AntennaGain:      floatToDoubleValueOrNil(params.AntennaGain),
		HeightM:          floatToDoubleValueOrNil(params.Heightm),
		HeightType:       stringToStringValueOrNil(params.HeightType),
		IndoorDeployment: boolToBoolValueOrNil(params.IndoorDeployment),
		LatitudeDeg:      floatToDoubleValueOrNil(params.LatitudeDeg),
		LongitudeDeg:     floatToDoubleValueOrNil(params.LongitudeDeg),
	}
}

func doubleValueToFloatOrNil(v *wrappers.DoubleValue) *float64 {
	if v == nil {
		return nil
	}
	return &v.Value
}

func boolValueToBoolOrNil(v *wrappers.BoolValue) *bool {
	if v == nil {
		return nil
	}
	return &v.Value
}

func stringValueToStringOrNil(v *wrappers.StringValue) *string {
	if v == nil {
		return nil
	}
	return &v.Value
}

func floatToDoubleValueOrNil(v *float64) *wrappers.DoubleValue {
	if v == nil {
		return nil
	}
	return wrapperspb.Double(*v)
}

func boolToBoolValueOrNil(v *bool) *wrapperspb.BoolValue {
	if v == nil {
		return nil
	}
	return wrapperspb.Bool(*v)
}

func stringToStringValueOrNil(v *string) *wrapperspb.StringValue {
	if v == nil {
		return nil
	}
	return wrapperspb.String(*v)
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
