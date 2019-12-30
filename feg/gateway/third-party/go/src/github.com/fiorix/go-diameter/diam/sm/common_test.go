// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package sm

import (
	"bytes"
	"net"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

func init() {
	dict.Default.Load(bytes.NewReader([]byte(acctDictionary)))
	dict.Default.Load(bytes.NewReader([]byte(authDictionary)))
}

var acctDictionary = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
	<application id="1001" type="acct">
	</application>
</diameter>
`

var authDictionary = `<?xml version="1.0" encoding="UTF-8"?>
<diameter>
	<application id="1002" type="auth">
	</application>
</diameter>
`

var (
	localhostAddress = datatype.Address(net.ParseIP("127.0.0.1"))

	serverSettings = &Settings{
		OriginHost:       "srv",
		OriginRealm:      "test",
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
	}

	serverSettings2 = &Settings{
		OriginHost:       "srv2",
		OriginRealm:      "test2",
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses:  []datatype.Address{localhostAddress},
	}

	clientSettings = &Settings{
		OriginHost:       "cli",
		OriginRealm:      "test",
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
	}

	clientSettings2 = &Settings{
		OriginHost:       "cli2",
		OriginRealm:      "test2",
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses:  []datatype.Address{localhostAddress},
	}
)
