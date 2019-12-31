// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smparser

import (
	"bytes"

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
