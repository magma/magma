/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package swx_proxy provides a thin client for using swx proxy service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package hlr_proxy

import (
	"context"
	"errors"
	"fmt"
	"log"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/hlr"
	"magma/feg/gateway/registry"

	"google.golang.org/grpc"
)

// Wrapper for GRPC Client
// functionality
type hlrProxyClient struct {
	hlr.HlrProxyClient
	cc *grpc.ClientConn
}

// getHlrProxyClient is a utility function to get a RPC connection to the
// HLR Proxy service
func getHlrProxyClient() (*hlrProxyClient, error) {
	var conn *grpc.ClientConn
	var err error
	conn, err = registry.GetConnection(registry.HLR_PROXY)
	if err != nil {
		errMsg := fmt.Sprintf("HLR Proxy client initialization error: %s", err)
		log.Print(errMsg)
		return nil, errors.New(errMsg)
	}
	return &hlrProxyClient{
		hlr.NewHlrProxyClient(conn),
		conn,
	}, err
}

const (
	resyncRandEnd = hlr.AuthInfoReq_ResyncInfo_RAND_LEN
	resyncAuthEnd = hlr.AuthInfoReq_ResyncInfo_RAND_LEN + hlr.AuthInfoReq_ResyncInfo_AUTH_LEN
)

// Authenticate - HLR equivalent of SWX Authenticate
func Authenticate(ctx context.Context, req *protos.AuthenticationRequest) (*protos.AuthenticationAnswer, error) {
	cli, err := getHlrProxyClient()
	if err != nil {
		return nil, err
	}
	hlrReq := &hlr.AuthInfoReq{
		UserName:                req.GetUserName(),
		NumRequestedUmtsVectors: req.GetSipNumAuthVectors(),
	}
	if rsLen := len(req.GetResyncInfo()); rsLen > int(resyncRandEnd) {
		hlrReq.ResyncInfo = &hlr.AuthInfoReq_ResyncInfo{
			Rand: req.GetResyncInfo()[:resyncRandEnd],
			Autn: req.GetResyncInfo()[resyncRandEnd:rsLen]}
	} else if rsLen > 0 {
		log.Printf("HLR Auth - Invalid ResyncInfo length: %d", rsLen)
	}
	hlrAns, err := cli.AuthInfo(ctx, hlrReq)

	res := &protos.AuthenticationAnswer{
		UserName:       req.GetUserName(),
		SipAuthVectors: []*protos.AuthenticationAnswer_SIPAuthVector{}}
	if err != nil {
		log.Printf("HLR RPC Error: %v", err)
		return res, err
	}
	if hlrAns.GetErrorCode() != hlr.ErrorCode_SUCCESS {
		msg := fmt.Sprintf("HLR Error: %s for User: %s", hlrAns.GetErrorCode().String(), req.GetUserName())
		log.Print(msg)
		return res, errors.New(msg)
	}
	for _, v := range hlrAns.GetUmtsVectors() {
		res.SipAuthVectors = append(res.SipAuthVectors, &protos.AuthenticationAnswer_SIPAuthVector{
			AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
			RandAutn:             append(v.GetRand(), v.GetAutn()...),
			Xres:                 v.GetXres(),
			ConfidentialityKey:   v.GetCk(),
			IntegrityKey:         v.GetIk(),
		})
	}
	return res, nil
}

// Register HLR equivalent of SWX register
func Register(_ context.Context, req *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{SessionId: req.SessionId}, nil
}

// Deregister HLR equivalent of SWX deregister
func Deregister(_ context.Context, req *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{SessionId: req.SessionId}, nil
}
