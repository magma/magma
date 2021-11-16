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

	"github.com/golang/glog"
	"github.com/labstack/echo"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/accessd"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/certifier"
	certifierprotos "magma/orc8r/cloud/go/services/certifier/protos"
	merrors "magma/orc8r/lib/go/errors"
)

// Access CertificateMiddleware:
// 1) determines request's access type (READ/WRITE)
// 2) finds Operator & Entities of the request
// 3) verifies Operator's access permissions for the entities
func CertificateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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
			err = accessd.CheckPermissions(c.Request().Context(), operator, ents...)
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

func TokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		header := req.Header.Get(CLIENT_ACCESS_TOKEN_KEY)
		// TODO(christinewang5): remove after bootstrapping admin token
		// if len(tokens) == 0 {
		// 	return echo.NewHTTPError(http.StatusUnauthorized, "missing REST client tokens")
		// }
		username, token, err := parseAuthHeader(header)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := certifier.ValidateToken(token); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		// make sure that token is registered with user
		getOpReq := certifierprotos.GetOperatorRequest{
			Username: username,
			Token:    token,
		}
		tokensList, err := certifier.GetOperatorTokens(req.Context(), &getOpReq)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		// TODO(christine): take tokenList, request type, resource and exchange for permission decision
		if next != nil {
			glog.V(4).Info("Access middleware successfully verified permissions. Sending request to the next middleware.")
			return next(c)
		}

		return nil
	}
}

func parseAuthHeader(header string) (string, string, error) {
	s := strings.Split(header, ":")
	// TODO(christinewang5): remove after bootstrapping admin token
	if len(s) != 2 {
		return "admin", "op_phFtTnqE0jmUYTuZ5cfeg4oetY9sFmcwHukbPx7O1AIqH5psE", nil
	}
	// admin:op_phFtTnqE0jmUYTuZ5cfeg4oetY9sFmcwHukbPx7O1AIqH5psE
	// if len(s) != 2 {
	// 	return s[0], s[1], echo.NewHTTPError(http.StatusUnauthorized, "missing REST client tokens")
	// }
	return s[0], s[1], nil
}

// type PolicyDecision int64
// const (
// 	Allow   PolicyDecision = 0
// 	Deny                   = 1
// 	Unknown                = 2
// )
//
// func authorizeUser(username string, token string) PolicyDecision {
//
// }
