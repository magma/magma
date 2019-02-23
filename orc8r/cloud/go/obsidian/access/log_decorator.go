/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package access

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// LogRequestDecorator closure, appends remote address, URI & certificate CN
// (if available from the passed http.Request) to every log string
func LogRequestDecorator(req *http.Request) func(f string, a ...interface{}) string {
	return func(f string, a ...interface{}) string {
		if req != nil {
			f += "; Remote: %s, URI: %s"
			a = append(a, req.RemoteAddr, req.RequestURI)
			ccn := req.Header.Get(CLIENT_CERT_CN_KEY)
			if len(ccn) > 0 {
				f += ", Cert CN: %s"
				a = append(a, ccn)
			}
		}
		return fmt.Sprintf(f, a...)
	}
}

// LogDecorator closure, appends remote address, URI & certificate CN
// (if available from the passed echo.Context) to every log string
func LogDecorator(c echo.Context) func(f string, a ...interface{}) string {
	var req *http.Request
	if c != nil {
		req = c.Request()
	}
	return LogRequestDecorator(req)
}
