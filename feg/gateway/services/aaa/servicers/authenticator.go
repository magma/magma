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

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"context"
	"fmt"
	"time"

	"github.com/emakeev/snowflake"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/events"
	"magma/feg/gateway/services/aaa/metrics"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/client"
	"magma/gateway/directoryd"
	orcprotos "magma/orc8r/lib/go/protos"
)

type eapAuth struct {
	supportedMethods []byte
	sessions         aaa.SessionTable // AAA SessionTable, if Nil -> Auth only mode
	config           *mconfig.AAAConfig
	sessionTout      time.Duration // Idle Session Timeout
	accounting       *accountingService
}

const (
	MacAddrKey  = "mac_addr"
	NanoInMilli = int64(time.Millisecond / time.Nanosecond)
)

var gatewayHardwareId string

func init() {
	uuid, _ := snowflake.Get()
	gatewayHardwareId = uuid.String()
	if len(gatewayHardwareId) == 0 {
		gatewayHardwareId = "hwid"
	}
}

// NewEapAuthenticator returns a new instance of EAP Auth service
func NewEapAuthenticator(
	sessions aaa.SessionTable,
	cfg *mconfig.AAAConfig,
	acct *accountingService) (protos.AuthenticatorServer, error) {

	return &eapAuth{
		supportedMethods: client.SupportedTypes(),
		sessions:         sessions,
		config:           cfg,
		sessionTout:      GetIdleSessionTimeout(cfg),
		accounting:       acct}, nil
}

// HandleIdentity passes Identity EAP payload to corresponding method provider & returns corresponding
// EAP result
// NOTE: Identity Request is handled by APs & does not involve EAP Authenticator's support
func (srv *eapAuth) HandleIdentity(ctx context.Context, in *protos.EapIdentity) (*protos.Eap, error) {
	resp, err := client.HandleIdentityResponse(uint8(in.GetMethod()), &protos.Eap{Payload: in.Payload, Ctx: in.Ctx})
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		errMsg := fmt.Sprintf("EAP HandleIdentity Error for Identity '%s', APN '%s': %v", resp.GetCtx().GetIdentity(), resp.GetCtx().GetApn(), err)
		glog.Error(errMsg)
		if srv.config.GetEventLoggingEnabled() {
			events.LogAuthenticationFailedEvent(in.GetCtx(), errMsg)
		}
		err = nil
	}
	return resp, err
}

// Handle handles passed EAP payload & returns corresponding EAP result
func (srv *eapAuth) Handle(ctx context.Context, in *protos.Eap) (*protos.Eap, error) {
	resp, err := client.Handle(in)
	if resp == nil {
		errMsg := fmt.Sprintf("Auth Handle error: %v, <nil> response", err)
		if srv.config.GetEventLoggingEnabled() {
			events.LogAuthenticationFailedEvent(in.GetCtx(), errMsg)
		}
		return resp, Errorf(codes.Internal, errMsg)
	}
	method := eap.Packet(resp.GetPayload()).Type()
	if method == uint8(protos.EapType_Reserved) {
		method = eap.Packet(in.GetPayload()).Type()
	}

	metrics.Auth.WithLabelValues(
		protos.EapCode(eap.Packet(resp.GetPayload()).Code()).String(),
		protos.EapType(method).String(),
		in.GetCtx().GetApn(),
		resp.Ctx.GetImsi()).Inc()

	if err != nil && len(resp.GetPayload()) > 0 {
		// log error, but do not return it to Radius. EAP will carry its own error
		errMsg := fmt.Sprintf("EAP Handle Error for Identity '%s', APN '%s': %v", resp.GetCtx().GetIdentity(), resp.GetCtx().GetApn(), err)
		if srv.config.GetEventLoggingEnabled() {
			events.LogAuthenticationFailedEvent(resp.GetCtx(), errMsg)
		}
		glog.Error(errMsg)
		return resp, nil
	}
	if srv.sessions != nil && eap.Packet(resp.Payload).IsSuccess() {
		if srv.config.GetEventLoggingEnabled() {
			events.LogAuthenticationSuccessEvent(resp.GetCtx())
		}

		imsi := resp.Ctx.GetImsi()
		if srv.config.GetAccountingEnabled() && srv.config.GetCreateSessionOnAuth() {
			if srv.accounting == nil {
				resp.Payload[eap.EapMsgCode] = eap.FailureCode
				glog.Errorf("Cannot Create Session on Auth: accounting service is missing")
				return resp, nil
			}
			csResp, err := srv.accounting.CreateSession(ctx, resp.Ctx)
			if err != nil {
				resp.Payload[eap.EapMsgCode] = eap.FailureCode
				glog.Errorf("Failed to create session: %v", err)
				return resp, nil
			}
			resp.Ctx.AcctSessionId = csResp.GetSessionId()
		}
		if srv.accounting != nil {
			_, err := srv.accounting.baseAccountingStart(resp.Ctx)
			if err != nil {
				resp.Payload[eap.EapMsgCode] = eap.FailureCode
				glog.Errorf("Accounting session start error: %v", err)
				return resp, nil
			}
		}
		// Add Session & overwrite an existing session with the same ID if present,
		// otherwise a UE can get stuck on buggy/non-unique AP or Radius session generation
		resp.Ctx.CreatedTimeMs = uint64(time.Now().UnixNano() / NanoInMilli)
		_, err := srv.sessions.AddSession(resp.Ctx, srv.sessionTout, srv.accounting.timeoutSessionNotifier, true)
		if err != nil {
			glog.Errorf("Error adding a new session for SID: %s: %v", resp.Ctx.GetSessionId(), err)
			return resp, nil // log error, but don't pass to caller, the auth only users will still be able to connect
		}
		updateRequest := &orcprotos.UpdateRecordRequest{
			Id:       imsi,
			Location: gatewayHardwareId,
			Fields:   map[string]string{MacAddrKey: resp.Ctx.GetMacAddr()},
		}
		// execute in a new goroutine in case calls to directoryd take long time
		go directoryd.UpdateRecord(updateRequest)
	}
	return resp, nil
}

// SupportedMethods returns sorted list (ascending, by type) of registered EAP Provider Methods
func (srv *eapAuth) SupportedMethods(ctx context.Context, in *protos.Void) (*protos.EapMethodList, error) {
	return &protos.EapMethodList{Methods: srv.supportedMethods}, nil
}
