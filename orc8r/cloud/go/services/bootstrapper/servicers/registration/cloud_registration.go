package registration

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	rootCAFilePath = "/var/opt/magma/certs/rootCA.pem"

	// length of timeout for the nonce
	timeoutDuration = 30 * time.Minute

	// number of characters that the nonce will have
	nonceLength = 30
)

type cloudRegistrationServicer struct {
	store Store
}

func NewCloudRegistrationServicer(store Store) (protos.CloudRegistrationServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Storage store is nil")
	}
	return &cloudRegistrationServicer{store}, nil
}

func (crs *cloudRegistrationServicer) GetToken(c context.Context, request *protos.GetTokenRequest) (*protos.GetTokenResponse, error) {
	networkId := request.GatewayDeviceInfo.NetworkId
	logicalId := request.GatewayDeviceInfo.LogicalId

	tokenInfo, err := crs.store.GetTokenInfoFromLogicalID(networkId, logicalId)
	if err != nil {
		glog.Infof("Could not get tokenInfo for networkID %v and logicalID %v: %v", networkId, logicalId, err)
	}

	refresh := request.Refresh || tokenInfo == nil || isTokenExpired(tokenInfo)
	if refresh {
		// TODO(reginawang3495) add a test with tokenInfo = nil to make sure it works
		tokenInfo, err = crs.generateAndSaveTokenInfo(networkId, logicalId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not generate and save tokenInfo for networkID %v and logicalID %v: %v", networkId, logicalId, err)
		}
	}

	return &protos.GetTokenResponse{Timeout: tokenInfo.Timeout, Token: nonceToToken(tokenInfo.Nonce)}, nil
}

func (crs *cloudRegistrationServicer) GetGatewayRegistrationInfo(c context.Context, request *protos.GetGatewayRegistrationInfoRequest) (*protos.GetGatewayRegistrationInfoResponse, error) {
	rootCA, err := getRootCA()
	if err != nil {
		return nil, err
	}
	domainName := getDomainName()

	return &protos.GetGatewayRegistrationInfoResponse{
		RootCa:               rootCA,
		DomainName:           domainName,
	}, nil
}

func (crs *cloudRegistrationServicer) GetGatewayDeviceInfo(c context.Context, request *protos.GetGatewayDeviceInfoRequest) (*protos.GetGatewayDeviceInfoResponse, error) {
	nonce := nonceFromToken(request.Token)

	tokenInfo, err := crs.store.GetTokenInfoFromNonce(nonce)
	if err != nil {
		return &protos.GetGatewayDeviceInfoResponse{
			Response: &protos.GetGatewayDeviceInfoResponse_Error{
				Error: fmt.Sprintf("Could not get token info from nonce %v: %v", nonce, err),
			},
		}, nil
	}

	return &protos.GetGatewayDeviceInfoResponse{
		Response:             &protos.GetGatewayDeviceInfoResponse_GatewayDeviceInfo{
			GatewayDeviceInfo: &protos.GatewayDeviceInfo{
				NetworkId:           tokenInfo.GatewayDeviceInfo.NetworkId,
				LogicalId:            tokenInfo.GatewayDeviceInfo.LogicalId,
			},
		},
	}, nil
}

func (crs *cloudRegistrationServicer) generateAndSaveTokenInfo(networkID string, logicalID string) (*protos.TokenInfo, error) {
	nonce := generateSecureNonce(nonceLength)

	t := time.Now().Add(timeoutDuration)
	timeout :=
		&timestamp.Timestamp{
			Seconds: int64(t.Second()),
			Nanos:   int32(t.Unix()),
		}

	tokenInfo := protos.TokenInfo{
		GatewayDeviceInfo: &protos.GatewayDeviceInfo{
			NetworkId: networkID,
			LogicalId: logicalID,
		},
		Nonce:   nonce,
		Timeout: timeout,
	}

	err := crs.store.SetTokenInfo(tokenInfo)
	if err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

func getRootCA() (string, error) {
	body, err := ioutil.ReadFile(rootCAFilePath)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// TODO(#10437)
func getDomainName() string {
	return "localhost"
}
