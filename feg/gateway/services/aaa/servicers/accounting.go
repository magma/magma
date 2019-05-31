/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/gateway/services/aaa/protos"
)

type accountingService struct {
}

// NewEapAuthenticator returns a new instance of EAP Auth service
func NewAccountingService() (protos.AccountingServer, error) {
	return &accountingService{}, nil
}

// Start implements Radius Acct-Status-Type: Start endpoint
func (s *accountingService) Start(_ context.Context, aaaCtx *protos.Context) (*protos.AcctResp, error) {
	return &protos.AcctResp{}, status.Errorf(codes.Unimplemented, "Not Implemented")
}

// InterimUpdate implements Radius Acct-Status-Type: Interim-Update endpoint
func (s *accountingService) InterimUpdate(_ context.Context, ur *protos.UpdateRequest) (*protos.AcctResp, error) {
	return &protos.AcctResp{}, status.Errorf(codes.Unimplemented, "Not Implemented")
}

// Stop implements Radius Acct-Status-Type: Stop endpoint
func (s *accountingService) Stop(_ context.Context, _ *protos.Context) (*protos.AcctResp, error) {
	return &protos.AcctResp{}, status.Errorf(codes.Unimplemented, "Implemented")
}
