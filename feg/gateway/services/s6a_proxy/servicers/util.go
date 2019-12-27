/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package servicers

import (
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"

	"github.com/fiorix/go-diameter/v4/diam"
)

func (s *s6aProxy) addDiamOriginAVPs(m *diam.Message) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.clientCfg.Host))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.clientCfg.Realm))
	if s.originStateID != 0 {
		m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(s.originStateID))
	}
}
