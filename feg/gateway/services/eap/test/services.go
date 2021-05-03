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

// Package test provides common definitions and function for eap related tests
package test

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos"
)

// SwxProxy - test SwxProxy proxy implementation
type SwxProxy struct{}

// Test SwxProxyServer implementation

// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s SwxProxy) Authenticate(
	ctx context.Context,
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {

	time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond * 10))))

	v, ok := Units[req.GetUserName()]
	if !ok {
		return &protos.AuthenticationAnswer{},
			status.Errorf(codes.PermissionDenied, "Unknown User: "+req.GetUserName())
	}
	res := &protos.AuthenticationAnswer{
		UserName: req.GetUserName(),
		SipAuthVectors: []*protos.AuthenticationAnswer_SIPAuthVector{
			{
				AuthenticationScheme: req.AuthenticationScheme,
				RandAutn:             v.RandAutn,
				Xres:                 v.Xres,
				ConfidentialityKey:   v.ConfidentialityKey,
				IntegrityKey:         v.IntegrityKey,
			},
		},
	}
	if req.RetrieveUserProfile {
		res.UserProfile = &protos.AuthenticationAnswer_UserProfile{Msisdn: v.MSISDN}
	}
	return res, nil
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s SwxProxy) Register(_ context.Context, _ *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{}, nil
}

// Deregister sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s SwxProxy) Deregister(_ context.Context, _ *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{}, nil
}

// NoUseSwxProxy - a dummu SwxProxy implementation which always returns an error & should not be called
type NoUseSwxProxy struct{}

//
// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s NoUseSwxProxy) Authenticate(
	ctx context.Context,
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {

	return nil, fmt.Errorf("Authenticate is NOT IMPLEMENTED")
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s NoUseSwxProxy) Register(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{}, fmt.Errorf("Register is NOT IMPLEMENTED")
}

// Deregister sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s NoUseSwxProxy) Deregister(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	return &protos.RegistrationAnswer{}, fmt.Errorf("Deregister is NOT IMPLEMENTED")
}
