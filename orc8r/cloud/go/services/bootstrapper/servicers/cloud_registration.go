package servicers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const nonceLength = 12

const timeoutDuration = 30 * time.Minute

type cloudRegistrationServicer struct {
	store storage.Store
}

func NewCloudRegistrationServicer(store storage.Store) (protos.CloudRegistrationServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Storage store is nil")
	}
	return &cloudRegistrationServicer{store}, nil
}

func (crs *cloudRegistrationServicer) GetToken(c context.Context, request *protos.GetTokenRequest) (*protos.GetTokenResponse, error) {
	n := request.Gateway.NetworkId
	l := request.Gateway.LogicalId.Id

	tokenInfo, err := crs.store.GetTokenInfoFromLogicalID(n, l)
	if err != nil {
		glog.Infof("Could not get tokenInfo for networkID %v and logicalID %v: %v", n, l, err)
	}

	refresh := request.Refresh || tokenInfo == nil || tokenTimedOut(tokenInfo)

	if refresh {
		// TODO(reginawang3495) add a test with tokenInfo = nil to make sure it works
		tokenInfo, err = crs.generateAndSaveToken(n, l, tokenInfo.Nonce)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not generate and save tokenInfo for networkID %v and logicalID %v: %v", n, l, err)
		}
	}

	return &protos.GetTokenResponse{Timeout: tokenInfo.Timeout, Token: nonceToToken(tokenInfo.Nonce)}, nil
}

func (crs *cloudRegistrationServicer) GetGatewayRegistrationInfo(c context.Context, request *protos.GetGatewayRegistrationInfoRequest) (*protos.GetGatewayRegistrationInfoResponse, error) {
	return nil, nil
}

func (crs *cloudRegistrationServicer) generateAndSaveToken(networkID string, logicalID string, oldNonce string) (*protos.TokenInfo, error) {
	nonce := generateSecureToken(nonceLength)
	t := time.Now().Add(timeoutDuration)

	timeout :=
		&timestamp.Timestamp{
			Seconds: int64(t.Second()),
			Nanos:   int32(t.Unix()),
		}
	tokenInfo := protos.TokenInfo{
		Gateway: &protos.GatewayInfo{
			NetworkId: networkID,
			LogicalId: &protos.AccessGatewayID{
				Id: logicalID,
			},
		},
		Nonce:   nonce,
		Timeout: timeout,
	}

	err := crs.store.SetTokenInfo(oldNonce, tokenInfo)
	if err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

func nonceToToken(nonce string) string {
	return bootstrapper.TokenPrepend + nonce
}

func tokenTimedOut(info *protos.TokenInfo) bool {
	return time.Now().Before(time.Unix(info.Timeout.Seconds, 0))
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
