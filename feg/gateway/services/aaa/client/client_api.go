/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package client provides a thin API client for communicating with AAA Server.
// This can be used by apps to discover and contact the service, without knowing about
// the underlying RPC implementation.
package client

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
)

type aaaClient struct {
	protos.AuthenticatorClient
	protos.AccountingClient
}

// getAaaClient is a utility function to get a RPC connection to the AAA service providing
// Authenticator & Accounting RPCs
func getAaaClient() (*aaaClient, error) {
	conn, err := registry.GetConnection(registry.AAA_SERVER)
	if err != nil {
		errMsg := fmt.Sprintf("AAA Server client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &aaaClient{
		protos.NewAuthenticatorClient(conn),
		protos.NewAccountingClient(conn),
	}, err
}

// HandleIdentity passes Identity EAP payload to corresponding method provider & returns corresponding
// EAP result
// NOTE: Identity Request is handled by APs & does not involve EAP Authenticator's support
func HandleIdentity(in *protos.EapIdentity) (*protos.Eap, error) {
	if in == nil {
		return nil, errors.New("Nil EapIdentity Parameter")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.HandleIdentity(context.Background(), in)
}

// Handle handles passed EAP payload & returns corresponding EAP result
func Handle(in *protos.Eap) (*protos.Eap, error) {
	if in == nil {
		return nil, errors.New("Nil Eap Parameter")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.Handle(context.Background(), in)
}

// Start implements Radius Acct-Status-Type: Start endpoint
func Start(aaaCtx *protos.Context) (*protos.AcctResp, error) {
	if aaaCtx == nil {
		return nil, errors.New("Nil AAA Ctx")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.Start(context.Background(), aaaCtx)
}

// InterimUpdate implements Radius Acct-Status-Type: Interim-Update endpoint
func InterimUpdate(ur *protos.UpdateRequest) (*protos.AcctResp, error) {
	if ur == nil {
		return nil, errors.New("Nil Interim Update Request")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.InterimUpdate(context.Background(), ur)
}

// Stop implements Radius Acct-Status-Type: Stop endpoint
func Stop(req *protos.StopRequest) (*protos.AcctResp, error) {
	if req == nil {
		return nil, errors.New("Nil Stop Request")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.Stop(context.Background(), req)
}
