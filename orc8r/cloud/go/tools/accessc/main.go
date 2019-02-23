/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Command Line Tool to create & manage Operators, ACLs and Certificates
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/tools/accessc/handlers"
)

func main() {
	flag.Parse()
	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})

	// Init help for all commands
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Printf(
			"\nUsage: \033[1m%s command [OPTIONS]\033[0m\n\n",
			filepath.Base(cmd))
		flag.PrintDefaults()
		fmt.Println("\nCommands:")
		handlers.CommandRegistry.Usage()
	}
	cmdName := flag.Arg(0)
	if len(flag.Args()) < 1 || cmdName == "" || cmdName == "help" {
		flag.Usage()
		os.Exit(1)
	}

	cmd := handlers.CommandRegistry.Get(cmdName)
	if cmd == nil {
		fmt.Println("\nInvalid Command: ", cmdName)
		flag.Usage()
		os.Exit(1)
	}
	args := os.Args[2:]
	cmd.Flags().Parse(args)
	os.Exit(cmd.Handle(args))
}
