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
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/metrics"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/session_manager"
	lte_protos "magma/lte/cloud/go/protos"
)

type accountingService struct {
	sessions    aaa.SessionTable
	config      *mconfig.AAAConfig
	sessionTout time.Duration // Idle Session Timeout
}

const (
	imsiPrefix = "IMSI"
)

// NewEapAuthenticator returns a new instance of EAP Auth service
func NewAccountingService(sessions aaa.SessionTable, cfg *mconfig.AAAConfig) (*accountingService, error) {
	return &accountingService{
		sessions:    sessions,
		config:      cfg,
		sessionTout: GetIdleSessionTimeout(cfg),
	}, nil
}

// Start implements Radius Acct-Status-Type: Start endpoint
func (srv *accountingService) Start(ctx context.Context, aaaCtx *protos.Context) (*protos.AcctResp, error) {
	if aaaCtx == nil {
		return &protos.AcctResp{}, status.Errorf(codes.InvalidArgument, "Nil AAA Context")
	}
	sid := aaaCtx.GetSessionId()
	s := srv.sessions.GetSession(sid)
	if s == nil {
		return &protos.AcctResp{}, Errorf(
			codes.FailedPrecondition, "Accounting Start: Session %s was not authenticated", sid)
	}
	var err error
	if srv.config.GetAccountingEnabled() && !srv.config.GetCreateSessionOnAuth() {
		_, err = srv.CreateSession(ctx, aaaCtx)
	} else {
		srv.sessions.SetTimeout(sid, srv.sessionTout, srv.timeoutSessionNotifier)
	}
	return &protos.AcctResp{}, err
}

// InterimUpdate implements Radius Acct-Status-Type: Interim-Update endpoint
func (srv *accountingService) InterimUpdate(_ context.Context, ur *protos.UpdateRequest) (*protos.AcctResp, error) {
	if ur == nil {
		return &protos.AcctResp{}, Errorf(codes.InvalidArgument, "Nil Update Request")
	}
	sid := ur.GetCtx().GetSessionId()
	s := srv.sessions.GetSession(sid)
	if s == nil {
		return &protos.AcctResp{}, Errorf(
			codes.FailedPrecondition, "Accounting Update: Session %s was not authenticated", sid)
	}
	srv.sessions.SetTimeout(sid, srv.sessionTout, srv.timeoutSessionNotifier)

	metrics.OctetsIn.WithLabelValues(s.GetCtx().GetApn(), s.GetCtx().GetImsi()).Add(float64(ur.GetOctetsIn()))
	metrics.OctetsOut.WithLabelValues(s.GetCtx().GetApn(), s.GetCtx().GetImsi()).Add(float64(ur.GetOctetsOut()))

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
		return &protos.AcctResp{}, Errorf(codes.FailedPrecondition, "Accounting Stop: Session %s is not found", sid)
	}

	s.Lock()
	sessionImsi := s.GetCtx().GetImsi()
	apn := s.GetCtx().GetApn()
	s.Unlock()

	var err error
	if srv.config.GetAccountingEnabled() {
		_, err = session_manager.EndSession(makeSID(sessionImsi))
		if err != nil {
			err = Error(codes.Unavailable, err)
		}
	}
	metrics.AcctStop.WithLabelValues(apn, sessionImsi)

	return &protos.AcctResp{}, err
}

// CreateSession is an "outbound" RPC for session manager which can be called from start()
func (srv *accountingService) CreateSession(grpcCtx context.Context, aaaCtx *protos.Context) (*protos.AcctResp, error) {

	startime := time.Now()

	mac, err := net.ParseMAC(aaaCtx.GetMacAddr())
	if err != nil {
		return &protos.AcctResp{}, Errorf(codes.InvalidArgument, "Invalid MAC Address: %v", err)
	}
	req := &lte_protos.LocalCreateSessionRequest{
		Sid:             makeSID(aaaCtx.GetImsi()),
		UeIpv4:          aaaCtx.GetIpAddr(),
		Msisdn:          ([]byte)(aaaCtx.GetMsisdn()),
		RatType:         lte_protos.RATType_TGPP_WLAN,
		HardwareAddr:    mac,
		RadiusSessionId: aaaCtx.GetSessionId(),
	}
	_, err = session_manager.CreateSession(req)
	if err == nil {
		srv.sessions.SetTimeout(req.GetRadiusSessionId(), srv.sessionTout, srv.timeoutSessionNotifier)
		metrics.CreateSessionLatency.Observe(time.Since(startime).Seconds())
	} else {
		err = Errorf(codes.Internal, "Create Session Error: %v", err)
	}

	return &protos.AcctResp{}, err
}

// TerminateSession is an "inbound" RPC from session manager to notify accounting of a client session termination
func (srv *accountingService) TerminateSession(
	ctx context.Context, req *protos.TerminateSessionRequest) (*protos.AcctResp, error) {

	sid := req.GetRadiusSessionId()
	s := srv.sessions.RemoveSession(sid)

	if s == nil {
		return &protos.AcctResp{}, Errorf(codes.FailedPrecondition, "Session %s is not found", sid)
	}

	s.Lock()
	sctx := s.GetCtx()
	imsi := sctx.GetImsi()
	apn := sctx.GetApn()
	s.Unlock()

	metrics.SessionTerminate.WithLabelValues(apn, imsi)

	if !strings.HasPrefix(imsi, imsiPrefix) {
		imsi = imsiPrefix + imsi
	}
	if imsi != req.GetImsi() {
		return &protos.AcctResp{}, Errorf(
			codes.InvalidArgument, "Mismatched IMSI: %s != %s of session %s", req.GetImsi(), imsi, sid)
	}
	conn, err := registry.GetConnection(registry.RADIUS)
	if err != nil {
		return &protos.AcctResp{}, Errorf(codes.Unavailable, "Error getting Radius RPC Connection: %v", err)
	}
	radcli := protos.NewAuthorizationClient(conn)
	_, err = radcli.Disconnect(ctx, &protos.DisconnectRequest{Ctx: sctx})
	if err != nil {
		err = Error(codes.Internal, err)
	}
	return &protos.AcctResp{}, err
}

// EndTimedOutSession is an "inbound" -> session manager AND "outbound" -> Radius server notification of a timed out
// session. It should be called for a timed out and recently removed from the sessions table session.
func (srv *accountingService) EndTimedOutSession(aaaCtx *protos.Context) error {
	if aaaCtx == nil {
		return status.Errorf(codes.InvalidArgument, "Nil AAA Context")
	}
	var err, radErr error

	if srv.config.GetAccountingEnabled() {
		_, err = session_manager.EndSession(makeSID(aaaCtx.GetImsi()))
	}

	conn, radErr := registry.GetConnection(registry.RADIUS)
	if radErr != nil {
		radErr = status.Errorf(codes.Unavailable, "Session Timeout Notification Radius Connection Error: %v", radErr)
	} else {
		_, radErr = protos.NewAuthorizationClient(conn).Disconnect(
			context.Background(), &protos.DisconnectRequest{Ctx: aaaCtx})
	}
	if radErr != nil {
		if err != nil {
			err = Errorf(
				codes.Internal, "Session Timeout Notification errors; session manager: %v, Radius: %v", err, radErr)
		} else {
			err = Error(codes.Unavailable, radErr)
		}
	}
	return err
}

func (srv *accountingService) timeoutSessionNotifier(s aaa.Session) error {
	if srv != nil && s != nil {
		return srv.EndTimedOutSession(s.GetCtx())
	}
	return nil
}

func makeSID(imsi string) *lte_protos.SubscriberID {
	if !strings.HasPrefix(imsi, imsiPrefix) {
		imsi = imsiPrefix + imsi
	}
	return &lte_protos.SubscriberID{Id: imsi, Type: lte_protos.SubscriberID_IMSI}
}
