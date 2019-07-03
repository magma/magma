/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"net"
	"time"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/session_manager"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/protos"
	lte_protos "magma/lte/cloud/go/protos"
)

type accountingService struct {
	sessions    aaa.SessionTable
	config      *mconfig.AAAConfig
	sessionTout time.Duration // Idle Session Timeout
}

// NewEapAuthenticator returns a new instance of EAP Auth service
func NewAccountingService(sessions aaa.SessionTable, cfg *mconfig.AAAConfig) (protos.AccountingServer, error) {
	tout := aaa.DefaultSessionTimeout
	if cfg != nil {
		tout = time.Millisecond * time.Duration(cfg.GetIdleSessionTimeoutMs())
	}
	return &accountingService{sessions: sessions, config: cfg, sessionTout: tout}, nil
}

// Start implements Radius Acct-Status-Type: Start endpoint
func (srv *accountingService) Start(ctx context.Context, aaaCtx *protos.Context) (*protos.AcctResp, error) {
	if aaaCtx == nil {
		return &protos.AcctResp{}, status.Errorf(codes.InvalidArgument, "Nil AAA Context")
	}
	sid := aaaCtx.GetSessionId()
	s := srv.sessions.GetSession(sid)
	if s == nil {
		return &protos.AcctResp{}, status.Errorf(
			codes.FailedPrecondition, "Accounting Start: Session %s was not authenticated", sid)
	}
	var err error
	if srv.config.GetAccountingEnabled() && !srv.config.GetCreateSessionOnAuth() {
		_, err = srv.CreateSession(ctx, aaaCtx)
	} else {
		srv.sessions.SetTimeout(sid, srv.sessionTout)
	}
	return &protos.AcctResp{}, err
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
	srv.sessions.SetTimeout(sid, srv.sessionTout)
	return &protos.AcctResp{}, nil
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
			codes.FailedPrecondition, "Accounting Stop: Session %s is not found", sid)
	}
	_, err := session_manager.EndSession(
		&lte_protos.SubscriberID{
			Id:   req.GetCtx().GetImsi(),
			Type: lte_protos.SubscriberID_IMSI})
	return &protos.AcctResp{}, err
}

// CreateSession is an "outbound" RPC for session manager which can be called from start()
func (srv *accountingService) CreateSession(grpcCtx context.Context, aaaCtx *protos.Context) (*protos.AcctResp, error) {

	mac, err := net.ParseMAC(aaaCtx.GetMacAddr())
	if err != nil {
		return &protos.AcctResp{}, status.Errorf(
			codes.InvalidArgument,
			"Invalid MAC Address: %v", err)
	}
	req := &lte_protos.LocalCreateSessionRequest{
		Sid:             &lte_protos.SubscriberID{Id: aaaCtx.GetImsi(), Type: lte_protos.SubscriberID_IMSI},
		Msisdn:          ([]byte)(aaaCtx.GetMsisdn()),
		RatType:         lte_protos.RATType_TGPP_WLAN,
		HardwareAddr:    mac,
		RadiusSessionId: aaaCtx.GetSessionId(),
	}
	_, err = session_manager.CreateSession(req)
	if err == nil {
		srv.sessions.SetTimeout(req.GetRadiusSessionId(), srv.sessionTout)
	}
	return &protos.AcctResp{}, err
}

// TerminateSession is an "inbound" RPC from session manager to notify accounting of a client session termination
func (srv *accountingService) TerminateSession(
	ctx context.Context, req *protos.TerminateSessionRequest) (*protos.AcctResp, error) {

	sid := req.GetRadiusSessionId()
	s := srv.sessions.RemoveSession(sid)
	if s == nil {
		return &protos.AcctResp{}, status.Errorf(codes.FailedPrecondition, "Session %s is not found", sid)
	}
	s.Lock()
	defer s.Unlock()
	imsi := s.GetCtx().GetImsi()
	if imsi != req.GetImsi() {
		return &protos.AcctResp{}, status.Errorf(
			codes.InvalidArgument, "Mismatched IMSI: %s != %s of session %s", req.GetImsi(), imsi, sid)
	}
	conn, err := registry.GetConnection(registry.RADIUS)
	if err != nil {
		return &protos.AcctResp{}, status.Errorf(
			codes.Unavailable, "Error getting Radius RPC Connection: %v", err)
	}
	radcli := protos.NewAuthorizationClient(conn)
	_, err = radcli.Disconnect(context.Background(), &protos.DisconnectRequest{Ctx: s.GetCtx()})
	return &protos.AcctResp{}, err
}
