/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/tools/mconfig_views_cli/handlers"
)

func main() {
	flag.Parse()
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Printf("\nUsage: %s command [OPTIONS]\n", filepath.Base(cmd))
		flag.PrintDefaults()
		fmt.Println("\nCommands:")
		handlers.Commands.Usage()
	}

	plugin.LoadAllPluginsFatalOnError(&plugin.DefaultOrchestratorPluginLoader{})

	cmdName := flag.Arg(0)
	if len(flag.Args()) < 1 || cmdName == "" || cmdName == "help" {
		flag.Usage()
		os.Exit(1)
	}

	cmd := handlers.Commands.Get(cmdName)
	if cmd == nil {
		fmt.Println("\nInvalid Command: ", cmdName)
		flag.Usage()
		os.Exit(1)
	}
	args := os.Args[2:]
	cmd.Flags().Parse(args)
	os.Exit(cmd.Handle(args))
}
