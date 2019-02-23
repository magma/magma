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

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/tools/commands"
)

// Delete command - removes given Operator, its ACLs & Certificates
func init() {
	cmd := CommandRegistry.Add(
		"delete",
		"Delete Given Operator, its ACL and all its Certificates",
		delete)
	cmd.Flags().Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s %s <Operator ID>\n", os.Args[0], cmd.Name())
	}
}

func delete(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	oid := strings.TrimSpace(f.Arg(0))
	if f.NArg() != 1 || len(oid) == 0 {
		f.Usage()
		log.Fatalf("A single Operator Id must be specified.")
	}
	// Create Operator Identity for the oid
	operator := identity.NewOperator(oid)
	fmt.Printf("Removing Operator: %s (%s)\n", oid, operator.HashString())
	certSNs, err := certifier.FindCertificates(operator)
	if err != nil {
		log.Printf("Error %s getting certificates for %s", err, oid)
	} else {
		fmt.Printf(
			"%d certificates associated with %s found\n", len(certSNs), oid)
		for _, csn := range certSNs {
			err = certifier.RevokeCertificateSN(csn)
			if err != nil {
				log.Printf(
					"Error %s deleting certificate SN:%s for %s", err, csn, oid)
			}
		}
	}
	err = accessd.DeleteOperator(operator)
	if err != nil {
		log.Printf("Error while removing Operator %s: %s", oid, err)
	}
	return 0
}
