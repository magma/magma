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
	"time"

	"magma/gateway/directoryd"
	orcprotos "magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/metrics"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/client"
)

type eapAuth struct {
	supportedMethods []byte
	sessions         aaa.SessionTable // AAA SessionTable, if Nil -> Auth only mode
	config           *mconfig.AAAConfig
	sessionTout      time.Duration // Idle Session Timeout
	accounting       *accountingService
}

const MacAddrKey = "mac_addr"

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
		log.Printf("EAP HandleIdentity Error: %v", err)
		err = nil
	}
	return resp, err
}

// Handle handles passed EAP payload & returns corresponding EAP result
func (srv *eapAuth) Handle(ctx context.Context, in *protos.Eap) (*protos.Eap, error) {
	resp, err := client.Handle(in)
	if resp == nil {
		return resp, Errorf(codes.Internal, "Auth Handle error: %v, <nil> response", err)
	}
	method := eap.Packet(resp.GetPayload()).Type()
	if method == uint8(protos.EapType_Reserved) {
		method = eap.Packet(in.GetPayload()).Type()
	}

	metrics.Auth.WithLabelValues(
		protos.EapCode(eap.Packet(resp.GetPayload()).Code()).String(),
		protos.EapType(method).String(),
		in.GetCtx().GetApn()).Inc()

	if err != nil && len(resp.GetPayload()) > 0 {
		// log error, but do not return it to Radius. EAP will carry its own error
		log.Printf("EAP Handle Error: %v", err)
		return resp, nil
	}
	if srv.sessions != nil && eap.Packet(resp.Payload).IsSuccess() {
		imsi := resp.Ctx.GetImsi()
		if srv.config.GetAccountingEnabled() && srv.config.GetCreateSessionOnAuth() {
			if srv.accounting == nil {
				resp.Payload[eap.EapMsgCode] = eap.FailureCode
				log.Printf("Cannot Create Session on Auth: accounting service is missing")
				return resp, nil
			}
			csResp, err := srv.accounting.CreateSession(ctx, resp.Ctx)
			if err != nil {
				resp.Payload[eap.EapMsgCode] = eap.FailureCode
				log.Printf("Failed to create session: %v", err)
				return resp, nil
			}
			resp.Ctx.AcctSessionId = csResp.GetSessionId()
		}
		// Add Session & overwrite an existing session with the same ID if present,
		// otherwise a UE can get stuck on buggy/non-unique AP or Radius session generation
		_, err := srv.sessions.AddSession(resp.Ctx, srv.sessionTout, srv.accounting.timeoutSessionNotifier, true)
		if err != nil {
			log.Printf("Error adding a new session for SID: %s: %v", resp.Ctx.GetSessionId(), err)
			return resp, nil // log error, but don't pass to caller, the auth only users will still be able to connect
		}
		updateRequest := &orcprotos.UpdateRecordRequest{
			Id:       imsi,
			Location: "hwid", // actual hwid will be filled in by directory service
			Fields:   map[string]string{MacAddrKey: resp.Ctx.GetMacAddr()},
		}
		go directoryd.UpdateRecord(updateRequest) // execute in a new goroutine in case calls to directoryd take long time
	}
	return resp, nil
}

// SupportedMethods returns sorted list (ascending, by type) of registered EAP Provider Methods
func (srv *eapAuth) SupportedMethods(ctx context.Context, in *protos.Void) (*protos.EapMethodList, error) {
	return &protos.EapMethodList{Methods: srv.supportedMethods}, nil
}
