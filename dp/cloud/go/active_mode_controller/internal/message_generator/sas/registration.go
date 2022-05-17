package sas

import (
	"encoding/json"
	"strings"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

type registrationRequestGenerator struct{}

func NewRegistrationRequestGenerator() *registrationRequestGenerator {
	return &registrationRequestGenerator{}
}

func (*registrationRequestGenerator) GenerateRequests(cbsd *active_mode.Cbsd) []*Request {
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
			HeightType:       installation.GetHeightType(),
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
