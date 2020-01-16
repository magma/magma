/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"fmt"
	"log"

	"magma/feg/gateway/services/aaa/metrics"
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
		res.Code = lteprotos.AbortSessionResult_USER_NOT_FOUND
		res.ErrorMessage = fmt.Sprintf("Session for IMSI: %s is not found", imsi)
		log.Print(res.ErrorMessage)
		return res, nil
	}
	s := srv.sessions.GetSession(sid)
	if s == nil {
		res.Code = lteprotos.AbortSessionResult_SESSION_NOT_FOUND
		res.ErrorMessage = fmt.Sprintf("Session for Radius Session ID: %s and IMSI: %s is not found", sid, imsi)
		log.Print(res.ErrorMessage)
		return res, nil
	}
	s.Lock()
	sctx := proto.Clone(s.GetCtx()).(*protos.Context)
	asid := sctx.AcctSessionId
	s.Unlock()
	if len(req.GetSessionId()) > 0 &&
		len(asid) > 0 &&
		asid != req.GetSessionId() {

		res.Code = lteprotos.AbortSessionResult_SESSION_NOT_FOUND
		res.ErrorMessage = fmt.Sprintf(
			"Accounting Session ID Mismatch for RadSID %s and IMSI: %s. Requested: %s, recorded: %s",
			sid, imsi, req.GetSessionId(), asid)
		log.Print(res.ErrorMessage)
		return res, nil
	}
	if srv.config.GetAccountingEnabled() {
		// ? can potentially end a new, valid session
		req := &lteprotos.LocalEndSessionRequest{
			Sid: makeSID(imsi),
			Apn: sctx.GetApn(),
		}
		session_manager.EndSession(req)
		metrics.EndSession.WithLabelValues(sctx.GetApn(), metrics.DecorateIMSI(sctx.GetImsi())).Inc()
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
		res.Code = lteprotos.AbortSessionResult_RADIUS_SERVER_ERROR
		res.ErrorMessage = fmt.Sprintf(
			"Radius Disconnect Error: %v for IMSI: %s, Acct SID: %s, Radius SID: %s", err, imsi, asid, sid)
		log.Print(res.ErrorMessage)
		return res, nil
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
		sid := makeSID(imsi)
		req := &lteprotos.LocalEndSessionRequest{
			Sid: sid,
			Apn: sctx.GetApn(),
		}
		session_manager.EndSession(req)
		metrics.EndSession.WithLabelValues(sctx.GetApn(), sid.Id).Inc()
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
