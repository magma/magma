/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// package initflag initializes (parses) Go flag if needed, it allows the noise free use of golog & other packages
// relying on flags being parsed
package initflag

import (
	"flag"
)

func init() {
	// only if not already parsed
	if !flag.Parsed() {
		// save original settings
		orgUsage := flag.CommandLine.Usage
		origOut := flag.CommandLine.Output()
		origErrorHandling := flag.CommandLine.ErrorHandling()

		// set to 'silent'
		flag.CommandLine.Init(flag.CommandLine.Name(), flag.ContinueOnError)
		flag.CommandLine.Usage = func() {}
		flag.CommandLine.SetOutput(devNull{})
		flag.Parse()

		// restore original settings
		flag.CommandLine.Init(flag.CommandLine.Name(), origErrorHandling)
		flag.CommandLine.Usage = orgUsage
		flag.CommandLine.SetOutput(origOut)
	}
}

type devNull struct{}

func (devNull) Write(b []byte) (int, error) {
	return len(b), nil
}
