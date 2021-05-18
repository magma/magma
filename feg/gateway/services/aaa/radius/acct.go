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

// package radius implements AAA server's radius interface for accounting & authentication
package radius

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"layeh.com/radius/rfc6911"

	"magma/feg/gateway/services/aaa/client"
	"magma/feg/gateway/services/aaa/protos"
)

type clientAcct struct {
	protos.UnimplementedAccountingServer
}

// Acct-Status-Type Start
func (*clientAcct) Start(_ context.Context, in *protos.Context) (*protos.AcctResp, error) {
	return client.Start(in)
}

// Accounting Interim-Update
func (*clientAcct) InterimUpdate(_ context.Context, in *protos.UpdateRequest) (*protos.AcctResp, error) {
	return client.InterimUpdate(in)
}

// Acct-Status-Type Stop
func (*clientAcct) Stop(_ context.Context, in *protos.StopRequest) (*protos.AcctResp, error) {
	return client.Stop(in)
}

// StartAuth starts Auth Radius server, it'll block on success & needs to be executed in its own thread/routine
func (s *Server) StartAcct() error {
	if s == nil {
		return fmt.Errorf("nil radius accounting server")
	}
	if s.accounting == nil { // accounting is not provided, use AAA RPC client
		s.accounting = &clientAcct{}
	}
	server := radius.PacketServer{
		Addr:         s.GetConfig().AcctAddr,
		Network:      s.GetConfig().Network,
		SecretSource: radius.StaticSecretSource(s.GetConfig().Secret),
		Handler:      &s.AcctServer,
	}
	glog.Infof("Starting Radius Accounting server on %s::%s", server.Network, server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		glog.Errorf("failed to start radius accounting server @ %s::%s - %v", server.Network, server.Addr, err)
	}
	return err
}

// ServeRADIUS - radius handler interface implementation for EAP server
func (s *AcctServer) ServeRADIUS(w radius.ResponseWriter, r *radius.Request) {
	var err error
	if w == nil || r == nil || r.Packet == nil {
		glog.Errorf("invalid request: %v", r)
		return
	}
	p := r.Packet
	if p.Code != radius.CodeAccountingRequest {
		glog.Errorf("unexpected request code: %s, dropping request from: %s", p.Code.String(), r.RemoteAddr.String())
		return
	}
	acctType := rfc2866.AcctStatusType_Get(p)
	acctSessionID := rfc2866.AcctSessionID_GetString(p)
	callingsid := rfc2865.CallingStationID_GetString(p)
	calledsid := rfc2865.CalledStationID_GetString(p)

	if len(acctSessionID) == 0 {
		acctSessionID = GenSessionID(callingsid, calledsid)
		glog.Warningf(
			"missing Acct-Session-Id for acct status type %s in message from %s; generated ID: '%s'",
			acctType, r.RemoteAddr.String(), acctSessionID)
		if len(acctSessionID) == 0 {
			glog.Errorf("empty session ID, dropping request from: %s", r.RemoteAddr.String())
			return
		}
	}

	aaaCtx := &protos.Context{
		SessionId: acctSessionID,
		MacAddr:   callingsid,
		Apn:       calledsid,
	}
	ip := rfc2865.FramedIPAddress_Get(p)
	if ip == nil {
		ip = rfc6911.FramedIPv6Address_Get(p)
	}
	if ip != nil {
		aaaCtx.IpAddr = ip.String()
	}
	switch acctType {
	case rfc2866.AcctStatusType_Value_AccountingOn:
	case rfc2866.AcctStatusType_Value_Start:
		_, err = s.accounting.Start(context.Background(), aaaCtx)
		if err != nil {
			glog.Errorf("accounting start error for remote %s, UE %s: %v", r.RemoteAddr.String(), aaaCtx.String(), err)
			return
		}
	case rfc2866.AcctStatusType_Value_InterimUpdate:
		_, err = s.accounting.InterimUpdate(r.Context(), makeUpdateReq(aaaCtx, p))
		if err != nil {
			glog.Errorf("interim update error for remote %s, UE %s: %v", r.RemoteAddr.String(), aaaCtx.String(), err)
			return
		}
	case rfc2866.AcctStatusType_Value_AccountingOff:
	case rfc2866.AcctStatusType_Value_Stop:
		updateReq := makeUpdateReq(aaaCtx, p)
		if (updateReq.OctetsIn | updateReq.OctetsOut | updateReq.PacketsIn | updateReq.PacketsOut) != 0 {
			_, err = s.accounting.InterimUpdate(r.Context(), updateReq)
			if err != nil {
				glog.Errorf("stop update error for remote %s, UE %s: %v", r.RemoteAddr.String(), aaaCtx.String(), err)
			}
		}
		stopRequest := &protos.StopRequest{
			Cause:     protos.StopRequest_NAS_REQUEST,
			Ctx:       aaaCtx,
			OctetsIn:  uint32(rfc2866.AcctInputOctets_Get(p)),
			OctetsOut: uint32(rfc2866.AcctOutputOctets_Get(p)),
		}
		_, err = s.accounting.Stop(r.Context(), stopRequest)
		if err != nil {
			glog.Errorf("accounting stop error for remote %s, UE %s: %v", r.RemoteAddr.String(), aaaCtx.String(), err)
			return
		}
	default:
		glog.Errorf("unknown Acct-Status-Type received for %s: %s", aaaCtx.String(), acctType.String())
		return
	}
	glog.V(2).Infof("successfully handled %s Accounting Request for %v", acctType.String(), aaaCtx)

	// Build & send response
	resp := p.Response(radius.CodeAccountingResponse)
	resp.Add(rfc2866.AcctSessionID_Type, radius.Attribute(acctSessionID))

	err = w.Write(resp)
	if err != nil {
		glog.Errorf("error sending %s response to %s: %v", r.Code.String(), r.RemoteAddr, err)
	}
}

func makeUpdateReq(aaaCtx *protos.Context, p *radius.Packet) *protos.UpdateRequest {
	u := &protos.UpdateRequest{
		OctetsIn:   uint32(rfc2866.AcctInputOctets_Get(p)),
		OctetsOut:  uint32(rfc2866.AcctOutputOctets_Get(p)),
		PacketsIn:  uint32(rfc2866.AcctInputPackets_Get(p)),
		PacketsOut: uint32(rfc2866.AcctOutputPackets_Get(p)),
		Ctx:        aaaCtx,
	}
	glog.V(2).Infof("Interim Update: %v", u)
	return u
}
