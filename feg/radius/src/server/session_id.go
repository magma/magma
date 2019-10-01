/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package server

import (
	"fmt"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
)

// GetSessionID Extracts the radius session id from the given radius request
func (s *Server) GetSessionID(r *radius.Request) string {
	calledStationIDAttr, _ := rfc2865.CalledStationID_Lookup(r.Packet)
	callingStationIDAttr, _ := rfc2865.CallingStationID_Lookup(r.Packet)

	return s.ComposeSessionID(
		string(calledStationIDAttr),
		string(callingStationIDAttr),
	)
}

func (s *Server) ComposeSessionID(calling string, called string) string {
	return fmt.Sprintf("%s__%s", string(calling), string(called))
}
