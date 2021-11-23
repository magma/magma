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
	rootCA string
}

func NewCloudRegistrationServicer(store Store) (protos.CloudRegistrationServer, error) {
	if store == nil {
		return nil, fmt.Errorf("storage store is nil")
	}
	rootCA, err := getRootCA()
	if err != nil {
		return nil, err
	}
	return &cloudRegistrationServicer{store, rootCA}, nil
}

func (crs *cloudRegistrationServicer) GetToken(c context.Context, request *protos.GetTokenRequest) (*protos.GetTokenResponse, error) {
	networkId := request.GatewayDeviceInfo.NetworkId
	logicalId := request.GatewayDeviceInfo.LogicalId

	tokenInfo, err := crs.store.GetTokenInfoFromLogicalID(networkId, logicalId)
	if err != nil {
		// error is not bubbled up since the tokenInfo is only important if the token exists and is expired
		// if GetTokenInfoFromLogicalID fails, we can just ignore and continue
		glog.V(2).Infof("could not get tokenInfo for networkID %v and logicalID %v: %v", networkId, logicalId, err)
	}

	refresh := request.Refresh || tokenInfo == nil || isTokenExpired(tokenInfo)
	if refresh {
		// TODO(reginawang3495) add a test with tokenInfo = nil to make sure it works
		tokenInfo, err = crs.generateAndSaveTokenInfo(networkId, logicalId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not generate and save tokenInfo for networkID %v and logicalID %v: %v", networkId, logicalId, err)
		}
	}

	res := &protos.GetTokenResponse{Timeout: tokenInfo.Timeout, Token: nonceToToken(tokenInfo.Nonce)}
	return res, nil
}

func (crs *cloudRegistrationServicer) GetGatewayRegistrationInfo(c context.Context, request *protos.GetGatewayRegistrationInfoRequest) (*protos.GetGatewayRegistrationInfoResponse, error) {
	domainName := getDomainName()

	res := &protos.GetGatewayRegistrationInfoResponse{
		RootCa:               crs.rootCA,
		DomainName:           domainName,
	}
	return res, nil
}

func (crs *cloudRegistrationServicer) GetGatewayDeviceInfo(c context.Context, request *protos.GetGatewayDeviceInfoRequest) (*protos.GetGatewayDeviceInfoResponse, error) {
	nonce := nonceFromToken(request.Token)

	tokenInfo, err := crs.store.GetTokenInfoFromNonce(nonce)
	if err != nil {
		res := &protos.GetGatewayDeviceInfoResponse{
			Response: &protos.GetGatewayDeviceInfoResponse_Error{
				Error: fmt.Sprintf("could not get token info from nonce %v: %v", nonce, err),
			},
		}
		return res, nil
	}

	res := &protos.GetGatewayDeviceInfoResponse{
		Response:             &protos.GetGatewayDeviceInfoResponse_GatewayDeviceInfo{
			GatewayDeviceInfo: &protos.GatewayDeviceInfo{
				NetworkId:           tokenInfo.GatewayDeviceInfo.NetworkId,
				LogicalId:            tokenInfo.GatewayDeviceInfo.LogicalId,
			},
		},
	}
	return res, nil
}

func (crs *cloudRegistrationServicer) generateAndSaveTokenInfo(networkID string, logicalID string) (*protos.TokenInfo, error) {
	nonce := generateNonce(nonceLength)

	t := time.Now().Add(timeoutDuration)

	tokenInfo := &protos.TokenInfo{
		GatewayDeviceInfo: &protos.GatewayDeviceInfo{
			NetworkId: networkID,
			LogicalId: logicalID,
		},
		Nonce:   nonce,
		Timeout: &timestamp.Timestamp{
			Seconds: int64(t.Second()),
			Nanos:   int32(t.Unix()),
		},
	}

	err := crs.store.SetTokenInfo(*tokenInfo)
	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
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
	return "warning: not implemented"
}
