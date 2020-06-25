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
	"fmt"
	"os"
)

func init() {
	flag.Var(
		&syslogDest,
		syslogFlag,
		"Redirect stderr to syslog, optional syslog destination in network::address format (system default otherwise)")

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
	// Check if the process needs to redirect stderr to syslog
	if flag.Lookup(syslogFlag) != nil {
		if err := redirectToSyslog(); err != nil {
			// Cannot use glog here, it should not be initialized yet
			fmt.Fprintf(os.Stderr, "ERROR redirecting to syslog: %v\n", err)
		}
	}
	// Check if the process needs to redirect stdout to stderr
	if *stdoutToStderr {
		stdout, os.Stdout = os.Stderr, os.Stdout
	}
}

type devNull struct{}

func (devNull) Write(b []byte) (int, error) {
	return len(b), nil
}
