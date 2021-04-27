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
	"layeh.com/radius/rfc2869"
	"layeh.com/radius/rfc6911"

	"magma/feg/gateway/services/aaa/client"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
)

type clientAuthenticator struct{}

// HandleIdentity passes Identity EAP payload to corresponding method provider & returns corresponding EAP result
func (clientAuthenticator) HandleIdentity(_ context.Context, in *protos.EapIdentity) (*protos.Eap, error) {
	return client.HandleIdentity(in)
}

// Handle handles passed EAP payload & returns corresponding EAP result
func (clientAuthenticator) Handle(_ context.Context, in *protos.Eap) (*protos.Eap, error) {
	return client.Handle(in)
}

// SupportedMethods returns sorted list (ascending, by type) of registered EAP Provider Methods
func (clientAuthenticator) SupportedMethods(context.Context, *protos.Void) (*protos.EapMethodList, error) {
	return client.SupportedMethods()
}

// StartAuth starts Auth Radius server, it'll block on success & needs to be executed in its own thread/routine
func (s *Server) StartAuth() error {
	if s == nil {
		return fmt.Errorf("nil radius auth server")
	}
	if s.authenticator == nil { // authenticator is not provided, use AAA RPC client
		s.authenticator = clientAuthenticator{}
	}
	methods, err := s.authenticator.SupportedMethods(context.Background(), &protos.Void{})
	if err != nil {
		glog.Errorf("radius EAP server failed to get supported EAP methods: %v", err)
		return err
	}
	s.authMethods = methods.GetMethods()
	server := radius.PacketServer{
		Addr:         s.GetConfig().AuthAddr,
		Network:      s.GetConfig().Network,
		SecretSource: radius.StaticSecretSource(s.GetConfig().Secret),
		Handler:      &s.AuthServer,
		// InsecureSkipVerify: false,
	}
	glog.Infof("Starting Radius EAP server on %s::%s", server.Network, server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		glog.Errorf("failed to start radius EAP server @ %s::%s - %v", server.Network, server.Addr, err)
	}
	return err
}

// ServeRADIUS - radius handler interface implementation for EAP server
func (s *AuthServer) ServeRADIUS(w radius.ResponseWriter, r *radius.Request) {
	if w == nil || r == nil || r.Packet == nil {
		glog.Errorf("invalid request: %v", r)
		return
	}
	p := r.Packet
	e := p.Get(rfc2869.EAPMessage_Type)
	if e == nil {
		glog.Errorf("%s request from %s is missing EAP Attribute", p.Code.String(), r.RemoteAddr.String())
		err := w.Write(p.Response(radius.CodeAccessReject))
		if err != nil {
			glog.Errorf("error sending access reject to %s: %v", r.RemoteAddr, err)
		}
		return
	}
	sessionCtx := &protos.Context{
		SessionId: rfc2866.AcctSessionID_GetString(p),
		Apn:       rfc2865.CalledStationID_GetString(p),
		MacAddr:   rfc2865.CallingStationID_GetString(p),
	}
	if len(sessionCtx.SessionId) == 0 {
		sessionCtx.SessionId = GenSessionID(sessionCtx.MacAddr, sessionCtx.Apn)
		glog.Warningf(
			"missing radius Session Id in message from %s; using generated ID: %s",
			r.RemoteAddr.String(), sessionCtx.AcctSessionId)
	}
	ip := rfc2865.FramedIPAddress_Get(p)
	if ip == nil {
		ip = rfc6911.FramedIPv6Address_Get(p)
	}
	if ip != nil {
		sessionCtx.IpAddr = ip.String()
	}
	eapp := eap.Packet(e)
	var (
		eapRes *protos.Eap
		err    error
	)
	if eapp.Code() == eap.ResponseCode && eapp.Type() == eap.MethodIdentity {
		// First, try to let authenticator choose matching provider
		eapRes, err = s.authenticator.HandleIdentity(r.Context(), &protos.EapIdentity{
			Payload: eapp,
			Ctx:     sessionCtx,
		})
		if err != nil {
			// couldn't find matching provider, iterate over all available providers
			for _, method := range s.authMethods {
				eapRes, err = s.authenticator.HandleIdentity(r.Context(), &protos.EapIdentity{
					Payload: eapp,
					Ctx:     sessionCtx,
					Method:  uint32(method),
				})
				if err == nil {
					break
				}
			}
		}
	} else {
		eapRes, err = s.authenticator.Handle(r.Context(), &protos.Eap{Payload: eapp, Ctx: sessionCtx})
	}
	if err != nil {
		glog.Errorf("EAP Handle error: %s", err)
		resp := p.Response(radius.CodeAccessReject)
		AddMessageAuthenticatorAttr(resp)
		err = w.Write(resp)
		if err != nil {
			glog.Errorf("error sending access reject to %s: %v", r.RemoteAddr, err)
		}
		return
	}
	postHandlerCtx := eapRes.GetCtx()
	eapPacket := eap.Packet(eapRes.Payload)
	eapCode := eapPacket.Code()
	resp := p.Response(ToRadiusCode(eapCode))
	resp.Add(rfc2869.EAPMessage_Type, eapRes.Payload)

	// Add key material for Access-Accept/EAP-Success message
	if resp.Code == radius.CodeAccessAccept {
		userNameAttr := p.Get(rfc2865.UserName_Type)
		if userNameAttr == nil {
			userNameAttr = []byte(postHandlerCtx.GetIdentity())
		}
		resp.Add(rfc2865.UserName_Type, userNameAttr)
		// Add optional Acct-Interim-Interval AVP to indicate that we want periodic Interim Updates from the client
		// If the client does not implement Acct-Interim-Interval AVP, has an alternative configuration or
		// its accounting is disabled - it'll ignore the AVP (see: https://tools.ietf.org/rfc/rfc2869.html#section-2.1)
		rfc2869.AcctInterimInterval_Add(resp, defaultAcctInterimUpdateInterval)
		// Add MPPE keys
		if rcv, snd, err := GetKeyingAttributes(postHandlerCtx.GetMsk(), r.Secret, r.Authenticator[:]); err != nil {
			glog.Errorf("keying material generate error for client at %s: %v", r.RemoteAddr, err)
		} else {
			resp.Add(rfc2865.VendorSpecific_Type, rcv)
			resp.Add(rfc2865.VendorSpecific_Type, snd)
		}
		glog.V(1).Infof("successfully authenticated user: %s", string(userNameAttr))
	}
	err = AddMessageAuthenticatorAttr(resp)
	if err != nil {
		glog.Errorf("failed to add Message-Authenticator AVP: %v", err)
	}
	err = w.Write(resp)
	if err != nil {
		glog.Errorf("error sending %s response to %s: %v", r.Code.String(), r.RemoteAddr, err)
	}
}

// GenSessionID creates synthetic radius session ID if none is supplied by the client
func GenSessionID(calling string, called string) string {
	return fmt.Sprintf("%s__%s", calling, called)
}

// ToRadiusCode returs the RADIUS packet code which, as per RFCxxxx
// should carry the EAP payload of the given EAP Code
func ToRadiusCode(eapCode uint8) radius.Code {
	switch eapCode {
	case eap.SuccessCode:
		return radius.CodeAccessAccept
	case eap.ResponseCode, eap.RequestCode:
		return radius.CodeAccessChallenge
	default:
		return radius.CodeAccessReject
	}
}
