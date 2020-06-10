/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// package initflag initializes (parses) Go flag if needed, it allows the noise free use of golog & other packages
// relying on flag being parsed
package initflag

import (
	"flag"
	"os"
)

func init() {
	if !flag.Parsed() {
		name := ""
		if len(os.Args) > 0 {
			name = os.Args[0]
		}
		flag.CommandLine.Init(name, flag.ContinueOnError)
		flag.CommandLine.Parse([]string{})
	}
}
