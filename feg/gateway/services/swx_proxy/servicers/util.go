/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *swxProxy) sendDiameterMsg(msg *diam.Message, retryCount uint) error {
	conn, err := s.connMan.GetConnection(s.smClient, s.config.ServerCfg)
	if err != nil {
		return err
	}
	err = conn.SendRequest(msg, retryCount)
	if err != nil {
		err = status.Errorf(codes.DataLoss, err.Error())
	}
	return err
}

func (s *swxProxy) IsHlrClient(imsi string) bool {
	if s != nil {
		return s.config.IsHlrClient(imsi)
	}
	return false
}
