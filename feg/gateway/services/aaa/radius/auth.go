/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package radius implements AAA server's radius interface for accounting & authentication
package radius

import (
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

// StartAuth starts Auth Radius server, it'll block on success & needs to be executed in its own thread/routine
func (s *Server) StartAuth() error {
	if s == nil {
		return fmt.Errorf("nil radius auth server")
	}
	methods, err := client.SupportedMethods()
	if err != nil {
		glog.Errorf("radius EAP server failed to get supported EAP methods: %v", err)
		return err
	}
	s.authMethods = methods.GetMethods()
	server := radius.PacketServer{
		Addr:         s.GetConfig().AuthAddr,
		Network:      s.GetConfig().Network,
		SecretSource: radius.StaticSecretSource(s.GetConfig().Secret),
		Handler:      s,
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
func (s *Server) ServeRADIUS(w radius.ResponseWriter, r *radius.Request) {
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
	ip := rfc2865.FramedIPAddress_Get(p)
	sessionCtx := &protos.Context{
		SessionId: rfc2866.AcctSessionID_GetString(p),
		Apn:       rfc2865.CalledStationID_GetString(p),
		MacAddr:   rfc2865.CallingStationID_GetString(p),
	}
	if len(sessionCtx.AcctSessionId) == 0 {
		sessionCtx.AcctSessionId = GenSessionID(sessionCtx.MacAddr, sessionCtx.Apn)
	}
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
		for _, method := range s.authMethods {
			eapRes, err = client.HandleIdentity(&protos.EapIdentity{
				Payload: eapp,
				Ctx:     sessionCtx,
				Method:  uint32(method),
			})
			if err == nil {
				break
			}
		}
	} else {
		eapRes, err = client.Handle(&protos.Eap{Payload: eapp, Ctx: sessionCtx})
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
	resp.Add(rfc2869.EAPMessage_Type, radius.Attribute(eapRes.Payload))

	// Add key material for Access-Accept/EAP-Success message
	if resp.Code == radius.CodeAccessAccept {
		userNameAttr := p.Get(rfc2865.UserName_Type)
		if userNameAttr == nil {
			userNameAttr = radius.Attribute([]byte(postHandlerCtx.GetIdentity()))
		}
		resp.Add(rfc2865.UserName_Type, userNameAttr)
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

// GenSessionID creates syntetic radius session ID if none is supplied by the client
func GenSessionID(calling string, called string) string {
	return fmt.Sprintf("%s__%s", string(calling), string(called))
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
