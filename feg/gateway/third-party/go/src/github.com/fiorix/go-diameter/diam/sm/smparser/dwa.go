// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smparser

import "github.com/fiorix/go-diameter/diam"

// DWA is a Device-Watchdog-Answer message.
// See RFC 6733 section 5.5.2 for details.
type DWA struct {
	ResultCode    uint32 `avp:"Result-Code"`
	OriginStateID uint32 `avp:"Origin-State-Id"`
}

// Parse parses the given message.
func (dwa *DWA) Parse(m *diam.Message) error {
	if err := m.Unmarshal(dwa); err != nil {
		return err
	}
	return nil
}
