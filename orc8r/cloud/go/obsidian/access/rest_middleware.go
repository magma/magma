/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package access

import (
	"fmt"
	"net/http"
	"strings"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/accessd"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/golang/glog"
	"github.com/labstack/echo"
)

// Access Middleware:
// 1) determines request's access type (READ/WRITE)
// 2) finds Operator & Entities of the request
// 3) verifies Operator's access permissions for the entities
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		decorate := getDecorator(c.Request())
		req := c.Request()
		if req == nil {
			return makeErr(decorate, http.StatusBadRequest, "invalid request")
		}
		glog.V(1).Infof("Received request in access middleware: %+v", req)

		operator, err := getOperator(req, decorate)
		if err != nil {
			return transformErr(decorate, err, http.StatusUnauthorized, "Invalid client credentials: %s", err)
		}
		if operator == nil {
			return makeErr(decorate, http.StatusUnauthorized, "missing client credentials")
		}

		perms := getRequestedPermissions(req, decorate)
		isStatic := strings.HasPrefix(c.Path(), obsidian.StaticURLPrefix) || strings.HasPrefix(c.Path(), obsidian.StaticURLPrefixLegacy)
		isStaticReadOnly := isStatic && perms == accessprotos.AccessControl_READ
		if !isStaticReadOnly {
			// Get Request's Entities' Ids
			ids := FindRequestedIdentities(c)

			// Check Operator's ACL for required entity permissions
			ents := make([]*accessprotos.AccessControl_Entity, 0, len(ids))
			for _, id := range ids {
				ents = append(ents, &accessprotos.AccessControl_Entity{Id: id, Permissions: perms})
			}
			err = accessd.CheckPermissions(operator, ents...)
			if err != nil {
				return transformErr(decorate, err, http.StatusForbidden, "access denied (%s)", err)
			}
		}

		if next != nil {
			glog.V(4).Info("Access middleware successfully verified permissions. Sending request to the next middleware.")
			return next(c)
		}

		return nil
	}
}

// getRequestedPermissions returns the required request permission (READ, WRITE
// or READ+WRITE) corresponding to the request method.
func getRequestedPermissions(req *http.Request, decorate logDecorator) accessprotos.AccessControl_Permission {
	switch req.Method {
	case "GET", "HEAD":
		return accessprotos.AccessControl_READ
	case "PUT", "POST", "DELETE":
		return accessprotos.AccessControl_WRITE
	default:
		glog.Info(decorate("Unclassified HTTP method '%s', defaulting to read+write requested permissions", req.Method))
		return accessprotos.AccessControl_READ | accessprotos.AccessControl_WRITE
	}
}

func transformErr(decorate logDecorator, err error, status int, errFmt string, errArgs ...interface{}) error {
	if _, ok := err.(merrors.ClientInitError); ok {
		return makeErr(decorate, http.StatusServiceUnavailable, "service unavailable")
	}
	return makeErr(decorate, status, errFmt, errArgs...)
}

func makeErr(decorate logDecorator, status int, errFmt string, errArgs ...interface{}) error {
	glog.V(1).Infof("REST middleware (obsidian) rejected request: %s", decorate(errFmt, errArgs...))
	return echo.NewHTTPError(status, fmt.Sprintf(errFmt, errArgs...))
}
