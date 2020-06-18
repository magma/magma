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
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/accessd"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/lib/go/errors"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

// Access Middleware:
// 1) determines request's access type (READ/WRITE)
// 2) finds Operator & Entities of the request
// 3) verifies Operator's access permissions for the entities

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c == nil || c.Request() == nil {
			return handleError(c, http.StatusBadRequest, "Invalid Request")
		}
		// find out request's access type (READ|WRITE|READ & WRITE)
		perm := requestPermissions(c)

		// Get Request's Operator
		oper, err := RequestOperator(c)
		if err != nil {
			if _, ok := err.(errors.ClientInitError); ok {
				return handleError(c, http.StatusServiceUnavailable, "Service Unavailable")
			}
			return handleError(
				c,
				http.StatusUnauthorized,
				"Client Credentials Error: %s", err)
		}
		if oper == nil {
			return handleError(
				c,
				http.StatusUnauthorized,
				"Missing Client Credentials")
		}

		// Bypass farther identity Checks for static docs GET having an
		// operator cert should be enough
		if urlPath := c.Path(); perm != accessprotos.AccessControl_READ || !(strings.HasPrefix(urlPath, obsidian.StaticURLPrefix)) {
			// Get Request's Entities' Ids
			ids := FindRequestedIdentities(c)

			// Check Operator's ACL for required entity permissions
			ents := make([]*accessprotos.AccessControl_Entity, len(ids))
			for i, e := range ids {
				ents[i] = &accessprotos.AccessControl_Entity{Id: e, Permissions: perm}
			}
			err = accessd.CheckPermissions(oper, ents...)
			if err != nil {
				if _, ok := err.(errors.ClientInitError); ok {
					return handleError(c, http.StatusServiceUnavailable, "Service Unavailable")
				}
				return handleError(
					c, http.StatusForbidden, "Access Denied (%s)", err)
			}
		}
		// all good, call next handler
		if next != nil {
			return next(c)
		}
		return nil
	}
}

// Return required request permission (READ, WRITE or READ & WRITE)
// corresponding to the request method
func requestPermissions(c echo.Context) accessprotos.AccessControl_Permission {
	// As a default - require ALL permissions (Read AND Write) for all
	// 'unclassified' methods
	perm := accessprotos.AccessControl_READ | accessprotos.AccessControl_WRITE
	// Find out if it's READ or WRITE
	switch c.Request().Method {
	case "GET", "HEAD":
		perm = accessprotos.AccessControl_READ
	case "PUT", "POST", "DELETE":
		perm = accessprotos.AccessControl_WRITE
	default:
		glog.Error(LogDecorator(c)("Unclassified HTTP Method: %s", c.Request().Method))
	}
	return perm
}

func handleError(
	c echo.Context,
	status int,
	f string, a ...interface{},
) error {
	glog.Error(LogDecorator(c)(f, a...))
	return echo.NewHTTPError(status, fmt.Sprintf(f, a...))
}
