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

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type RegistrationRequestGenerator struct{}

func (*RegistrationRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
	req := buildRegistrationRequest(cbsd)
	return []*Request{asRequest(Registration, req)}
}

func buildRegistrationRequest(cbsd *active_mode.Cbsd) *registrationRequest {
	settings := cbsd.GetSasSettings()
	if !settings.GetSingleStepEnabled() {
		return &registrationRequest{
			UserId:           settings.GetUserId(),
			FccId:            settings.GetFccId(),
			CbsdSerialNumber: settings.GetSerialNumber(),
		}
	}
	installation := cbsd.GetInstallationParams()
	return &registrationRequest{
		UserId:           settings.GetUserId(),
		FccId:            settings.GetFccId(),
		CbsdSerialNumber: settings.GetSerialNumber(),
		CbsdCategory:     strings.ToUpper(settings.GetCbsdCategory()),
		AirInterface:     &airInterface{RadioTechnology: "E_UTRA"},
		InstallationParam: &installationParam{
			Latitude:         installation.GetLatitudeDeg(),
			Longitude:        installation.GetLongitudeDeg(),
			Height:           installation.GetHeightM(),
			HeightType:       strings.ToUpper(installation.GetHeightType()),
			IndoorDeployment: installation.GetIndoorDeployment(),
			AntennaGain:      installation.GetAntennaGainDbi(),
		},
		MeasCapability: json.RawMessage("[]"),
	}
}

type registrationRequest struct {
	UserId            string             `json:"userId"`
	FccId             string             `json:"fccId"`
	CbsdSerialNumber  string             `json:"cbsdSerialNumber"`
	CbsdCategory      string             `json:"cbsdCategory,omitempty"`
	AirInterface      *airInterface      `json:"airInterface,omitempty"`
	InstallationParam *installationParam `json:"installationParam,omitempty"`
	MeasCapability    json.RawMessage    `json:"measCapability,omitempty"`
}

type airInterface struct {
	RadioTechnology string `json:"radioTechnology"`
}

type installationParam struct {
	Latitude         float32 `json:"latitude"`
	Longitude        float32 `json:"longitude"`
	Height           float32 `json:"height"`
	HeightType       string  `json:"heightType"`
	IndoorDeployment bool    `json:"indoorDeployment"`
	AntennaGain      float32 `json:"antennaGain"`
}
