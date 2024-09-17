package registration

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/go-openapi/strfmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/lib/go/protos"
)

// RegistrationService is public for ease of testing and mocking out functions
type RegistrationService struct {
	GetGatewayDeviceInfo func(ctx context.Context, token string) (*protos.GatewayDeviceInfo, error)
	RegisterDevice       func(deviceInfo *protos.GatewayDeviceInfo, hwid *protos.AccessGatewayID, challengeKey *protos.ChallengeKey) error
	GetControlProxy      func(networkID string) (string, error)
}

func NewRegistrationServicer() protos.RegistrationServer {
	return &RegistrationService{
		GetGatewayDeviceInfo: bootstrapper.GetGatewayDeviceInfo,
		RegisterDevice:       RegisterDevice,
		GetControlProxy:      GetControlProxy,
	}
}

func (r *RegistrationService) Register(c context.Context, request *protos.RegisterRequest) (*protos.RegisterResponse, error) {
	deviceInfo, err := r.GetGatewayDeviceInfo(context.Background(), request.Token)
	if err != nil {
		clientErr := makeErr(fmt.Sprintf("could not get device info from token %v: %v", request.Token, err))
		return clientErr, nil
	}

	err = updateGatewayDevice(c, deviceInfo, request.Hwid, request.ChallengeKey)
	if err != nil {
		clientErr := makeErr(fmt.Sprintf("error updating gateway: %v", err))
		return clientErr, nil
	}

	err = r.RegisterDevice(deviceInfo, request.Hwid, request.ChallengeKey)
	if err != nil {
		clientErr := makeErr(fmt.Sprintf("error registering device: %v", err))
		return clientErr, nil
	}

	controlProxy, err := r.GetControlProxy(deviceInfo.NetworkId)
	if err != nil {
		clientErr := makeErr(fmt.Sprintf("error getting control proxy: %v", err))
		return clientErr, nil
	}

	res := &protos.RegisterResponse{
		Response: &protos.RegisterResponse_ControlProxy{
			ControlProxy: controlProxy,
		},
	}
	return res, nil
}

func RegisterDevice(deviceInfo *protos.GatewayDeviceInfo, hwid *protos.AccessGatewayID, challengeKey *protos.ChallengeKey) error {
	gatewayRecord, err := createGatewayDevice(hwid, challengeKey)
	if err != nil {
		return err
	}
	err = device.RegisterDevice(context.Background(), deviceInfo.GetNetworkId(), orc8r.AccessGatewayRecordType, hwid.Id, gatewayRecord, serdes.Device)
	return err
}

func GetControlProxy(networkID string) (string, error) {
	controlProxy, err := tenants.GetControlProxyFromNetworkID(context.Background(), networkID)
	if err != nil {
		clientErr := fmt.Errorf("could not get control-proxy from tenant with network ID %s: %w", networkID, err)
		return "", clientErr
	}

	return controlProxy.GetControlProxy(), nil
}

// makeErr makes a protos.RegisterResponse_Error for protos.RegisterResponse
func makeErr(errString string) *protos.RegisterResponse {
	errRes := &protos.RegisterResponse{
		Response: &protos.RegisterResponse_Error{
			Error: errString,
		},
	}
	return errRes
}

// createGatewayDevice creates the gateway device model
func createGatewayDevice(hwID *protos.AccessGatewayID, challengeKey *protos.ChallengeKey) (*models.GatewayDevice, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(string(challengeKey.Key))
	if err != nil {
		return nil, err
	}
	challengeKeyBase64 := strfmt.Base64(decodedKey)
	gatewayDevice := &models.GatewayDevice{
		HardwareID: hwID.Id,
		Key:        &models.ChallengeKey{KeyType: challengeKey.KeyType.String(), Key: &challengeKeyBase64},
	}
	return gatewayDevice, nil
}

// updateGatewayDevice writes to the device information to the gateway entity
func updateGatewayDevice(ctx context.Context, deviceInfo *protos.GatewayDeviceInfo, hwID *protos.AccessGatewayID, challengeKey *protos.ChallengeKey) error {
	networkID := deviceInfo.NetworkId
	gatewayID := deviceInfo.LogicalId

	ent, err := configurator.LoadEntity(
		ctx,
		networkID, orc8r.MagmadGatewayType, gatewayID,
		configurator.EntityLoadCriteria{},
		serdes.Entity,
	)
	if err != nil {
		return err
	}

	device, err := createGatewayDevice(hwID, challengeKey)
	if err != nil {
		return err
	}

	gw := (&models.MagmadGateway{}).FromBackendModels(ent, device, nil)

	_, err = configurator.UpdateEntities(ctx, networkID, gw.ToEntityUpdateCriteria(ent), serdes.Entity)
	if err != nil {
		return err
	}

	return nil
}
