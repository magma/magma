/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/client"
)

type eapAuth struct {
	supportedMethods []byte
	sessions         aaa.SessionTable // AAA SessionTable, if Nil -> Auth only mode
	config           *mconfig.AAAConfig
	accounting       protos.AccountingServer
}

// NewEapAuthenticator returns a new instance of EAP Auth service
func NewEapAuthenticator(
	sessions aaa.SessionTable,
	cfg *mconfig.AAAConfig,
	acct protos.AccountingServer) (protos.AuthenticatorServer, error) {

	return &eapAuth{supportedMethods: client.SupportedTypes(), sessions: sessions, config: cfg, accounting: acct}, nil
}

// HandleIdentity passes Identity EAP payload to corresponding method provider & returns corresponding
// EAP result
// NOTE: Identity Request is handled by APs & does not involve EAP Authenticator's support
func (srv *eapAuth) HandleIdentity(ctx context.Context, in *protos.EapIdentity) (*protos.Eap, error) {
	resp, err := client.HandleIdentityResponse(uint8(in.GetMethod()), &protos.Eap{Payload: in.Payload, Ctx: in.Ctx})
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		log.Printf("EAP HandleIdentity Error: %v", err)
		err = nil
	}
	return resp, err
}

// Handle handles passed EAP payload & returns corresponding EAP result
func (srv *eapAuth) Handle(ctx context.Context, in *protos.Eap) (*protos.Eap, error) {
	resp, err := client.Handle(in)
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		log.Printf("EAP Handle Error: %v", err)
		err = nil
	} else if srv.sessions != nil && resp != nil && eap.Packet(resp.Payload).Type() == eap.SuccessCode {
		s, err := srv.sessions.AddSession(resp.Ctx, aaa.DefaultSessionTimeout)
		if err != nil {
			log.Printf("Error adding a new session for SID: %s: %v", resp.Ctx.GetSessionId(), err)
			if s != nil {
				s.Lock()
				defer s.Unlock()
				if resp.Ctx.GetImsi() != s.GetCtx().GetImsi() {
					// same user, just overwrite it
					// different IMSI, same Session ID - likely Radius server issue, LOG & overwrite
					log.Printf(
						"Radius Request Error: Same Session ID (%s) is used for different IMSIs, old: %s, new: %s",
						resp.Ctx.GetSessionId(), s.GetCtx().GetImsi(), resp.Ctx.GetImsi())
				}
				s.SetCtx(resp.Ctx)
			}
			if srv.config.GetAccountingEnabled() && srv.config.GetCreateSessionOnAuth() &&
				resp.Payload[eap.EapMsgCode] == eap.SuccessCode {

				if srv.accounting == nil {
					return resp, status.Errorf(
						codes.Unavailable,
						"Cannot Create Session on Auth: accounting service is missing")
				}
				_, err = srv.accounting.CreateSession(ctx, in.Ctx)
				if err != nil {
					resp.Payload[eap.EapMsgCode] = eap.FailureCode
					return resp, err
				}
			}
		}
	}
	return resp, err
}

// SupportedMethods returns sorted list (ascending, by type) of registered EAP Provider Methods
func (srv *eapAuth) SupportedMethods(ctx context.Context, in *protos.Void) (*protos.EapMethodList, error) {
	return &protos.EapMethodList{Methods: srv.supportedMethods}, nil
}
