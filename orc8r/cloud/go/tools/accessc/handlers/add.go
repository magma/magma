/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers implements individual accessc commands as well as common
// across multiple commands functionality
package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"magma/orc8r/cloud/go/tools/commands"
)

// Add command - Creates a new Operator and its ACL from specified Entities with
// required permissions, also creates a new Operator Certificate and associates
// it with the Operator

func init() {
	cmd := CommandRegistry.Add(
		"add",
		"Add a new Operator with a new Certificate",
		add)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s %s [OPTIONS] <OperatorID>\n", os.Args[0], cmd.Name())
		f.PrintDefaults()
	}
	f.StringVar(&certFName, "cert", "",
		"Name used for Operator's certificate files: Name.pem, Name.key.pem. "+
			"Defaults to: '<Operator ID>_cert' if not provided")

	addInit(f) // Bind flags common to add & add-admin,  see common_add.go

	entHelp := "%s with required permissions in the form: <network Id|*>:" +
		"R|W|RW. Use '*:RW' for admin permissions. At least one permission " +
		"must be given."

	f.Var(&networks, "n", fmt.Sprintf(entHelp, "Networks"))
	f.Var(&operators, "o", fmt.Sprintf(entHelp, "Operators"))
	f.Var(&gateways, "g", fmt.Sprintf(entHelp, "Gateways"))
}

// add Handler
func add(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	oid := strings.TrimSpace(f.Arg(0))
	if f.NArg() != 1 || len(oid) == 0 {
		f.Usage()
		log.Fatalf("A single Operator Id must be specified.")
	}
	acl := BuildACLForEntities(networks, operators, gateways)
	if len(acl) == 0 {
		f.Usage()
		log.Fatal("At least one ACL entity must be provided")
	}
	return addACL(oid, acl) // see common_add.go
}
