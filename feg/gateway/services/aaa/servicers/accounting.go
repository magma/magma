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

	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/protos"
)

type accountingService struct {
	sessions aaa.SessionTable
}

// NewEapAuthenticator returns a new instance of EAP Auth service
func NewAccountingService(sessions aaa.SessionTable) (protos.AccountingServer, error) {
	return &accountingService{sessions: sessions}, nil
}

// Start implements Radius Acct-Status-Type: Start endpoint
func (srv *accountingService) Start(_ context.Context, aaaCtx *protos.Context) (*protos.AcctResp, error) {
	if aaaCtx == nil {
		return &protos.AcctResp{}, status.Errorf(codes.InvalidArgument, "Nil AAA Context")
	}
	sid := aaaCtx.GetSessionId()
	s := srv.sessions.GetSession(sid)
	if s == nil {
		return &protos.AcctResp{}, status.Errorf(
			codes.FailedPrecondition, "Accounting Start: Session %s was not authenticated", sid)
	}
	return &protos.AcctResp{}, status.Errorf(codes.Unimplemented, "Not Implemented")
}

// InterimUpdate implements Radius Acct-Status-Type: Interim-Update endpoint
func (srv *accountingService) InterimUpdate(_ context.Context, ur *protos.UpdateRequest) (*protos.AcctResp, error) {
	if ur == nil {
		return &protos.AcctResp{}, status.Errorf(codes.InvalidArgument, "Nil Update Request")
	}
	sid := ur.GetCtx().GetSessionId()
	s := srv.sessions.GetSession(sid)
	if s == nil {
		return &protos.AcctResp{}, status.Errorf(
			codes.FailedPrecondition, "Accounting Update: Session %s was not authenticated", sid)
	}
	return &protos.AcctResp{}, status.Errorf(codes.Unimplemented, "Not Implemented")
}

// Stop implements Radius Acct-Status-Type: Stop endpoint
func (srv *accountingService) Stop(_ context.Context, req *protos.StopRequest) (*protos.AcctResp, error) {
	if req == nil {
		return &protos.AcctResp{}, status.Errorf(codes.InvalidArgument, "Nil Stop Request")
	}
	sid := req.GetCtx().GetSessionId()
	s := srv.sessions.RemoveSession(sid)
	if s == nil {
		return &protos.AcctResp{}, status.Errorf(
			codes.FailedPrecondition, "Accounting Stop: Session %s cannot be found", sid)
	}
	return &protos.AcctResp{}, status.Errorf(codes.Unimplemented, "Not Implemented")
}

// CreateSession is an "outbound" RPC for session manager which can be called from start()
func (srv *accountingService) CreateSession(context.Context, *protos.Context) (*protos.AcctResp, error) {
	return &protos.AcctResp{}, status.Errorf(codes.Unimplemented, "Not Implemented")
}

// TerminateSession is an "inbound" RPC from session manager to notify accounting of a client session termination
func (srv *accountingService) TerminateSession(
	context.Context, *protos.TerminateSessionRequest) (*protos.AcctResp, error) {

	return &protos.AcctResp{}, status.Errorf(codes.Unimplemented, "Not Implemented")
}
