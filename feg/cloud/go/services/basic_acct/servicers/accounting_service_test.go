/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import (
	"context"
	"net"
	"sync"
	"testing"

	anpb "github.com/magma/augmented-networks/accounting/protos"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/basic_acct/tests"
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/service/middleware/unary"
	"magma/orc8r/cloud/go/service/middleware/unary/test_utils"
)

func TestAugmentedNetworkAccounting(t *testing.T) {
	var testSessionId = "xYz2183028327839028e302"

	tests.SetupNetworks(t)

	// Test mapping of provider & consume from CTX
	ctx, _, _, err := unary.SetIdentityFromContext(getIncomingContextWithCertificate(t, tests.AgwHwId), nil, nil)
	session := &protos.AcctSession{
		User:      &protos.AcctSession_IMSI{IMSI: tests.NhImsi},
		SessionId: testSessionId,
	}
	provider, consumer, gw, err := RetrieveParticipants(ctx, session)
	assert.NoError(t, err)
	assert.Equal(t, tests.FederatedLteNetworkID, provider)
	assert.Equal(t, tests.ServingFegNetworkID, consumer)
	assert.Equal(t, tests.AgwId, gw)

	// Start mock ANSC accounting service
	anscServer := &anContractServer{}
	lis, err := net.Listen("tcp", "localhost:")
	assert.NoError(t, err)

	grpcServer := grpc.NewServer()
	anpb.RegisterAccountingServer(grpcServer, anscServer)
	go grpcServer.Serve(lis)

	// Create FeG Acct Service
	fegAccSrvr := &BaseAccService{}
	fegAccSrvr.SetConfig(&Config{
		RemoteAddr: lis.Addr().String(),
		NoTls:      true,
	})
	_, err = fegAccSrvr.Start(ctx, session)
	assert.NoError(t, err)

	_, err = fegAccSrvr.Update(ctx, &protos.AcctUpdateReq{
		Session:   session,
		OctetsIn:  1024,
		OctetsOut: 50120,
	})
	assert.NoError(t, err)

	// Duplicate acct
	go func() {
		_, err = fegAccSrvr.Update(ctx, &protos.AcctUpdateReq{
			Session:   session,
			OctetsIn:  1024 + 512,
			OctetsOut: 50120 + 2560,
		})
		assert.NoError(t, err)
	}()
	_, err = fegAccSrvr.Update(ctx, &protos.AcctUpdateReq{
		Session:   session,
		OctetsIn:  1024 + 512,
		OctetsOut: 50120 + 2560,
	})
	assert.NoError(t, err)

	_, err = fegAccSrvr.Stop(ctx, &protos.AcctUpdateReq{
		Session:   session,
		OctetsIn:  1024 + 512 + 128,
		OctetsOut: 50120 + 2560 + 1280,
	})
	assert.NoError(t, err)

	contr := anscServer.contracts[contractKey{
		consumer: tests.ServingFegNetworkID,
		provider: tests.FederatedLteNetworkID}]

	inOctets, outOctets := contr.totalIn(), contr.totalOut()
	assert.Equal(t, uint64(1664), inOctets)
	assert.Equal(t, uint64(53960), outOctets)
}

func getIncomingContextWithCertificate(t *testing.T, hwID string) context.Context {
	csn := test_utils.StartMockGwAccessControl(t, []string{hwID})
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(identity.CLIENT_CERT_SN_KEY, csn[0]))
}

//
// Mock ANSC Remote Service
//
type contractKey struct {
	consumer, provider string
}

type contract map[string]struct {
	in, out uint64
}

// totalIn summs up all IN octets per contract
func (c contract) totalIn() (sum uint64) {
	for _, v := range c {
		sum += v.in
	}
	return
}

// totalIn summs up all OUT octets per contract
func (c contract) totalOut() (sum uint64) {
	for _, v := range c {
		sum += v.out
	}
	return
}

type anContractServer struct {
	anpb.UnimplementedAccountingServer
	sync.Mutex
	contracts map[contractKey]contract
}

func (s *anContractServer) Start(_ context.Context, in *anpb.Session) (*anpb.SessionResp, error) {
	if s == nil {
		return &anpb.SessionResp{}, nil
	}
	s.Lock()
	defer s.Unlock()

	if s.contracts == nil {
		s.contracts = map[contractKey]contract{}
	}
	key := contractKey{
		consumer: in.GetConsumerId(),
		provider: in.GetProviderId(),
	}
	sessionKey := in.GetIMSI() + in.GetSessionId()
	s.contracts[key] = contract{sessionKey: struct{ in, out uint64 }{0, 0}}
	return &anpb.SessionResp{}, nil
}

func (s *anContractServer) Update(_ context.Context, in *anpb.UpdateReq) (*anpb.SessionResp, error) {
	if s == nil || s.contracts == nil {
		return &anpb.SessionResp{}, nil
	}
	s.Lock()
	defer s.Unlock()

	key := contractKey{
		consumer: in.GetSession().GetConsumerId(),
		provider: in.GetSession().GetProviderId(),
	}
	if contr := s.contracts[key]; contr != nil {
		sessionKey := in.GetSession().GetIMSI() + in.GetSession().GetSessionId()
		if sessionAcct, ok := contr[sessionKey]; ok {
			if in.GetOctetsIn() > sessionAcct.in {
				sessionAcct.in = in.GetOctetsIn()
			}
			if in.GetOctetsOut() > sessionAcct.out {
				sessionAcct.out = in.GetOctetsOut()
			}
			contr[sessionKey] = sessionAcct
		}
	}
	return &anpb.SessionResp{}, nil
}

func (s *anContractServer) Stop(ctx context.Context, in *anpb.UpdateReq) (*anpb.StopResp, error) {
	s.Update(ctx, in)
	return &anpb.StopResp{}, nil
}
