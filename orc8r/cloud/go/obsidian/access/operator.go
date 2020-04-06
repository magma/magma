/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package access

import (
	"fmt"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

// RequestOperator returns Identity of request's Operator (client)
// If either the request is missing TLS certificate headers or the certificate's
// SN is not found by Certifier or one of certificate & its identity checks fail
// - nil will be returned & the corresponding error logged
func RequestOperator(c echo.Context) (*protos.Identity, error) {
	if c == nil {
		glog.Error("Nil Echo Context")
		return nil, fmt.Errorf("Internal Server Error (Context)") // nil CTX, no useful info to log here
	}
	req := c.Request()
	if req == nil {
		glog.Error("Nil HTTP Request")
		return nil, fmt.Errorf("Internal Server Error (Request)")
	}

	// Get Certificate SN header value
	// TBD: to optimize - use map directly
	csn := req.Header.Get(CLIENT_CERT_SN_KEY)
	if len(csn) == 0 {
		glog.Warning(LogDecorator(c)("Missing REST Client Certificate"))
		return nil, fmt.Errorf("Missing Client Certificate")
	}
	certInfo, err := certifier.GetCertificateIdentity(csn)
	if err != nil {
		glog.Error(LogDecorator(c)(
			"Certificate SN '%s' lookup error '%s'", csn, err))
		if _, ok := err.(errors.ClientInitError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("Unknown Client Certificate SN: %s, err: %v", csn, err)
	}
	if certInfo == nil {
		glog.Error(LogDecorator(c)("No Certificate Info for SN: %s", csn))
		return nil, fmt.Errorf("Unregistered Client Certificate, SN: %s", csn)
	}
	// Check if certificate time is not expired/not active yet
	err = certifier.VerifyDateRange(certInfo)
	if err != nil {
		glog.Error(LogDecorator(c)(
			"Certificate Validation Error '%s' for SN: %s", err, csn))
		return nil, fmt.Errorf("Certificate Validation Error: %s", err)
	}
	opId := certInfo.Id
	if opId == nil {
		glog.Error(LogDecorator(c)("Nil Identity for Certificate SN: %s", csn))
		return nil, fmt.Errorf("Internal Server Error (Identity)")
	}
	// Check if it's operator identity
	if !identity.IsOperator(opId) {
		glog.Error(LogDecorator(c)(
			"Identity (%s) of CSN %s is not Operator", opId.HashString(), csn))
		return nil, fmt.Errorf("Internal Server Error (Operator)")
	}
	// all checks are OK, return it
	return certInfo.Id, nil
}
