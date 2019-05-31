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

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap/client"
)

type eapAuth struct {
	supportedMethods []byte
}

// NewEapAuthenticator returns a new instance of EAP Auth service
func NewEapAuthenticator() (protos.AuthenticatorServer, error) {
	return &eapAuth{supportedMethods: client.SupportedTypes()}, nil
}

// HandleIdentity passes Identity EAP payload to corresponding method provider & returns corresponding
// EAP result
// NOTE: Identity Request is handled by APs & does not involve EAP Authenticator's support
func (s *eapAuth) HandleIdentity(ctx context.Context, in *protos.EapIdentity) (*protos.Eap, error) {
	resp, err := client.HandleIdentityResponse(uint8(in.GetMethod()), &protos.Eap{Payload: in.Payload, Ctx: in.Ctx})
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		log.Printf("EAP HandleIdentity Error: %v", err)
		err = nil
	}
	return resp, err
}

// Handle handles passed EAP payload & returns corresponding EAP result
func (s *eapAuth) Handle(ctx context.Context, in *protos.Eap) (*protos.Eap, error) {
	resp, err := client.Handle(in)
	if err != nil && resp != nil && len(resp.GetPayload()) > 0 {
		log.Printf("EAP Handle Error: %v", err)
		err = nil
	}
	return resp, err
}

// SupportedMethods returns sorted list (ascending, by type) of registered EAP Provider Methods
func (s *eapAuth) SupportedMethods(ctx context.Context, in *protos.Void) (*protos.EapMethodList, error) {
	return &protos.EapMethodList{Methods: s.supportedMethods}, nil
}
