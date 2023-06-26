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

package sas

import (
	"encoding/json"
	"strings"

	"magma/dp/cloud/go/services/dp/storage"
)

type RegistrationRequestGenerator struct{}

func (*RegistrationRequestGenerator) GenerateRequests(cbsd *storage.DetailedCbsd) []*storage.MutableRequest {
	payload := buildRegistrationRequest(cbsd.Cbsd)
	req := makeRequest(Registration, payload)
	return []*storage.MutableRequest{req}
}

func buildRegistrationRequest(cbsd *storage.DBCbsd) *RegistrationRequest {
	if !cbsd.SingleStepEnabled.Bool {
		return &RegistrationRequest{
			UserId:           cbsd.UserId.String,
			FccId:            cbsd.FccId.String,
			CbsdSerialNumber: cbsd.CbsdSerialNumber.String,
		}
	}
	return &RegistrationRequest{
		UserId:           cbsd.UserId.String,
		FccId:            cbsd.FccId.String,
		CbsdSerialNumber: cbsd.CbsdSerialNumber.String,
		CbsdCategory:     strings.ToUpper(cbsd.CbsdCategory.String),
		AirInterface:     &AirInterface{RadioTechnology: "E_UTRA"},
		InstallationParam: &InstallationParam{
			Latitude:         cbsd.LatitudeDeg.Float64,
			Longitude:        cbsd.LongitudeDeg.Float64,
			Height:           cbsd.HeightM.Float64,
			HeightType:       strings.ToUpper(cbsd.HeightType.String),
			IndoorDeployment: cbsd.IndoorDeployment.Bool,
			AntennaGain:      cbsd.AntennaGainDbi.Float64,
		},
		MeasCapability: json.RawMessage("[]"),
	}
}

type RegistrationRequest struct {
	UserId            string             `json:"userId"`
	FccId             string             `json:"fccId"`
	CbsdSerialNumber  string             `json:"cbsdSerialNumber"`
	CbsdCategory      string             `json:"cbsdCategory,omitempty"`
	AirInterface      *AirInterface      `json:"airInterface,omitempty"`
	InstallationParam *InstallationParam `json:"installationParam,omitempty"`
	MeasCapability    json.RawMessage    `json:"measCapability,omitempty"`
}

type AirInterface struct {
	RadioTechnology string `json:"radioTechnology"`
}

type InstallationParam struct {
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Height           float64 `json:"height"`
	HeightType       string  `json:"heightType"`
	IndoorDeployment bool    `json:"indoorDeployment"`
	AntennaGain      float64 `json:"antennaGain"`
}
