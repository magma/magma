/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package servicers

import (
	"fmt"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/registry"
)

const (
	// MaxDiamRTRetries - number of retries for responding to RTR
	MaxDiamRTRetries = 1
)

type fegRelayClient struct {
	registry registry.CloudRegistry
}

type CloseableSwxGatewayServiceResponderClient struct {
	protos.SwxGatewayServiceClient
	conn *grpc.ClientConn
}

func (client *CloseableSwxGatewayServiceResponderClient) Close() {
	client.conn.Close()
}

// GetSwxGatewayServiceResponderClient returns a client to the local terminate registration client
func GetSwxGatewayServiceResponderClient(
	cloudRegistry registry.CloudRegistry) (*CloseableSwxGatewayServiceResponderClient, error) {

	conn, err := cloudRegistry.GetCloudConnection(feg_relay.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to SWx Gateway Relay: %s", err)
	}
	return &CloseableSwxGatewayServiceResponderClient{
		SwxGatewayServiceClient: protos.NewSwxGatewayServiceClient(conn),
		conn:                    conn,
	}, nil
}

func (r *fegRelayClient) RelayRTR(rtr *RTR) (protos.ErrorCode, error) {
	var err error
	if r == nil || r.registry == nil {
		err = fmt.Errorf("No relay registry for RTR")
		return protos.ErrorCode_UNABLE_TO_DELIVER, err
	}
	client, err := GetSwxGatewayServiceResponderClient(r.registry)
	if err != nil {
		return protos.ErrorCode_UNABLE_TO_DELIVER, err
	}
	defer client.Close()

	_, err = client.TerminateRegistration(context.Background(), &protos.RegistrationTerminationRequest{
		UserName:   string(rtr.UserName),
		ReasonCode: protos.RegistrationTerminationRequest_ReasonCode(rtr.DeregistrationReason.ReasonCode),
		ReasonInfo: string(rtr.DeregistrationReason.ReasonInfo),
		SessionId:  string(rtr.SessionID),
	})
	if err != nil {
		return protos.ErrorCode_LIMITED_SUCCESS, err
	}
	return protos.ErrorCode_SUCCESS, nil
}

func (r *fegRelayClient) RelayASR(*diameter.ASR) (protos.ErrorCode, error) {
	return protos.ErrorCode_COMMAND_UNSUPORTED, fmt.Errorf("Relay for ASR is not implemented")
}

func handleRTR(s *swxProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("handling RTR %v\n", m)
		var rtr RTR
		err := m.Unmarshal(&rtr)
		if err != nil {
			glog.Errorf("RTR Unmarshal failed for remote %s & message %s: %v", c.RemoteAddr(), m, err)
			return
		}
		imsi := string(rtr.UserName)
		if len(imsi) == 0 {
			imsi, err = diameter.ExtractImsiFromSessionID(string(rtr.SessionID))
			if err != nil {
				err = fmt.Errorf("Error retreiving IMSI from Session ID %s: %s", rtr.SessionID, err)
				glog.Error(err)
				err = s.sendRTA(c, m, protos.ErrorCode_UNKNOWN_SESSION_ID, &rtr, MaxDiamRTRetries)
				return
			}
			rtr.UserName = datatype.UTF8String(imsi)
		}
		if s.cache != nil {
			s.cache.Remove(imsi)
		}
		go func() {
			code, err := s.Relay.RelayRTR(&rtr)
			if err != nil {
				glog.Error(err)
			}

			err = s.sendRTA(c, m, code, &rtr, MaxDiamRTRetries)
			if err != nil {
				glog.Errorf("Failed to send RTA: %v", err)
			}
		}()
	}
}

func (s *swxProxy) sendRTA(c diam.Conn, m *diam.Message, code protos.ErrorCode, rtr *RTR, retries uint) error {
	ans := m.Answer(uint32(code))
	// SessionID is required to be the AVP in position 1
	ans.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(rtr.SessionID)))
	ans.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(rtr.AuthSessionState))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.config.ClientCfg.Host))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.config.ClientCfg.Realm))
	if s.originStateID != 0 {
		m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(s.originStateID))
	}
	_, err := ans.WriteToWithRetry(c, retries)
	return err
}
