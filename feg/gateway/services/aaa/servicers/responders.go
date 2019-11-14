/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"magma/orc8r/gateway/directoryd"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/session_manager"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/cloud/go/protos"
)

const (
	MinIMSILen = 10
	MaxIMSILen = 16
)

// AbortSession is a method of AbortSessionResponder service.
func (srv *accountingService) AbortSession(
	ctx context.Context, req *lteprotos.AbortSessionRequest) (*lteprotos.AbortSessionResult, error) {

	res := &lteprotos.AbortSessionResult{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil Request")
	}
	imsi := req.GetUserName()
	if len(imsi) < MinIMSILen || len(imsi) > MaxIMSILen {
		return res, Errorf(codes.InvalidArgument, "Invalid IMSI: %s", imsi)
	}
	sid := srv.sessions.FindSession(imsi)
	if len(sid) == 0 {
		return res, Errorf(codes.NotFound, "Session for IMSI: %s is not found", imsi)
	}
	s := srv.sessions.GetSession(sid)
	if s == nil {
		return res, Errorf(codes.Internal, "Session for RadSID: %s and IMSI: %s is not found", sid, imsi)
	}
	s.Lock()
	sctx := proto.Clone(s.GetCtx()).(*protos.Context)
	asid := sctx.AcctSessionId
	s.Unlock()
	if len(req.GetSessionId()) > 0 &&
		len(asid) > 0 &&
		asid != req.GetSessionId() {
		return res, Errorf(codes.FailedPrecondition,
			"Accounting Session ID Mismatch for RadSID %s and IMSI: %s. Requested: %s, recorded: %s",
			sid, imsi, req.GetSessionId(), asid)
	}
	if srv.config.GetAccountingEnabled() {
		// ? can potentially end a new, valid session
		session_manager.EndSession(makeSID(imsi))
	} else {
		deleteRequest := &orcprotos.DeleteRecordRequest{
			Id: imsi,
		}
		directoryd.DeleteRecord(deleteRequest)

	}
	srv.sessions.RemoveSession(sid)
	conn, err := registry.GetConnection(registry.RADIUS)
	if err != nil {
		return res, Errorf(codes.Unavailable, "Error getting Radius RPC Connection: %v", err)
	}
	radcli := protos.NewAuthorizationClient(conn)
	_, err = radcli.Disconnect(ctx, &protos.DisconnectRequest{Ctx: sctx})
	if err != nil {
		err = Error(codes.Internal, err)
	}
	return res, err
}

// TerminateRegistration is a method of SWx Gateway Service Responder service.
func (srv *accountingService) TerminateRegistration(
	ctx context.Context, req *fegprotos.RegistrationTerminationRequest) (*fegprotos.RegistrationAnswer, error) {

	res := &fegprotos.RegistrationAnswer{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil Request")
	}
	imsi := req.GetUserName()
	if len(imsi) < MinIMSILen {
		return res, Errorf(codes.InvalidArgument, "Invalid IMSI: %s", imsi)
	}
	sid := srv.sessions.FindSession(imsi)
	if len(sid) == 0 {
		return res, Errorf(codes.NotFound, "Session for IMSI: %s is not found", imsi)
	}
	s := srv.sessions.GetSession(sid)
	if s == nil {
		return res, Errorf(codes.Internal, "Session for RadSID: %s and IMSI: %s is not found", sid, imsi)
	}
	s.Lock()
	sctx := proto.Clone(s.GetCtx()).(*protos.Context)
	authSid := sctx.AuthSessionId
	acctSid := sctx.AcctSessionId
	s.Unlock()
	if len(req.GetSessionId()) > 0 &&
		(len(authSid) > 0 || len(acctSid) > 0) &&
		(authSid != req.GetSessionId() && acctSid != req.GetSessionId()) {
		return res, Errorf(codes.FailedPrecondition,
			"Accounting Session ID Mismatch for RadSID %s and IMSI: %s. Requested: %s, recorded: auth: %s | acct: %s",
			sid, imsi, req.GetSessionId(), authSid, acctSid)
	}
	deleteRequest := &orcprotos.DeleteRecordRequest{
		Id: imsi,
	}
	directoryd.DeleteRecord(deleteRequest) // remove it from directoryd even if session manager will try to remove it again

	if srv.config.GetAccountingEnabled() {
		// ? can potentially end a new, valid session
		session_manager.EndSession(makeSID(imsi))
	}

	srv.sessions.RemoveSession(sid)
	conn, err := registry.GetConnection(registry.RADIUS)
	if err != nil {
		return res, Errorf(codes.Unavailable, "Error getting Radius RPC Connection: %v", err)
	}
	radcli := protos.NewAuthorizationClient(conn)
	_, err = radcli.Disconnect(ctx, &protos.DisconnectRequest{Ctx: sctx})
	if err != nil {
		err = Error(codes.Internal, err)
	}
	return res, err
}
