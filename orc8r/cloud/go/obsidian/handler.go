/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"magma/orc8r/cloud/go/util"

	"github.com/labstack/echo"
	"google.golang.org/grpc"
)

type (
	HttpMethod             byte
	handlerRegistry        map[string]echo.HandlerFunc
	echoHandlerInitializer func(*echo.Echo, string, echo.HandlerFunc, ...echo.MiddlewareFunc) *echo.Route
)

// Handler wraps a function which serves a specified path and http method.
type Handler struct {
	Path string

	// Methods is a bitmask so one Handler can support multiple http methods.
	// See consts defined below.
	Methods HttpMethod

	HandlerFunc echo.HandlerFunc
}

const (
	GET HttpMethod = 1 << iota
	POST
	PUT
	DELETE
	ALL = GET | POST | PUT | DELETE
)

const (
	wildcard        = "*"
	networkWildcard = "N*"
)

var registries = map[HttpMethod]handlerRegistry{
	GET:    {},
	POST:   {},
	PUT:    {},
	DELETE: {},
}

var echoHandlerInitializers = map[HttpMethod]echoHandlerInitializer{
	GET:    (*echo.Echo).GET,
	POST:   (*echo.Echo).POST,
	PUT:    (*echo.Echo).PUT,
	DELETE: (*echo.Echo).DELETE,
}

// nopWriter wraps an http.ResponseWriter to no-op the Write() method.
// We need this to prevent multiplexed handlers from writing the same return
// value to the context response twice.
type nopWriter struct {
	writer http.ResponseWriter
}

func (n *nopWriter) Header() http.Header {
	return n.writer.Header()
}

func (*nopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (n *nopWriter) WriteHeader(statusCode int) {
	n.writer.WriteHeader(statusCode)
}

func register(registry handlerRegistry, handler Handler) error {
	_, registered := registry[handler.Path]
	if registered {
		return fmt.Errorf("HandlerFunc[s] already registered for path: %q", handler.Path)
	}
	registry[handler.Path] = handler.HandlerFunc
	return nil
}

// Register registers a given handler for given path and HTTP methods
// Note: the handlers won't become active until they are 'attached' to the echo
// server, see AttachAll below
func Register(handler Handler) error {
	if (handler.Methods & ^ALL) != 0 {
		return fmt.Errorf("Invalid handler method[s]: %b", handler.Methods)
	}

	if len(handler.Path) == 0 {
		return errors.New("Empty path is not supported")
	}
	for method := GET; method < ALL; method <<= 1 {
		if (method & handler.Methods) != 0 {
			err := register(registries[method], handler)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Unregister unregisters the handler for the specified path and HttpMethod if
// it is registered. No action will be taken if no such handler is registered.
func Unregister(path string, methods HttpMethod) {
	reg, regExists := registries[methods]
	if regExists {
		_, handlerExists := reg[path]
		if handlerExists {
			delete(reg, path)
		}
	}
}

// RegisterAll registers an array of Handlers. If an error is encountered while
// registering any handler, RegisterAll will exit early with that error and
// rollback any handlers which were already registered.
func RegisterAll(handlers []Handler) error {
	for i, handler := range handlers {
		if err := Register(handler); err != nil {
			for rollbackIdx := 0; rollbackIdx < i; rollbackIdx++ {
				Unregister(handlers[rollbackIdx].Path, handlers[rollbackIdx].Methods)
			}
			return err
		}
	}
	return nil
}

// AttachAll activates all registered (see: Register above) handlers
// Main package should call AttachAll after all handlers were registered
func AttachAll(e *echo.Echo, m ...echo.MiddlewareFunc) {
	for method, registry := range registries {
		ei := echoHandlerInitializers[method]
		if ei != nil {
			for path, handler := range registry {
				ei(e, path, handler, m...)
			}
		}
	}
}

func HttpError(err error, code ...int) *echo.HTTPError {
	var status = http.StatusInternalServerError
	if len(code) > 0 && code[0] >= http.StatusContinue &&
		code[0] <= http.StatusNetworkAuthenticationRequired {
		status = code[0]
	}
	log.Printf("REST HTTP Error: %s, Status: %d", err, status)
	return echo.NewHTTPError(status, grpc.ErrorDesc(err))
}

func CheckWildcardNetworkAccess(c echo.Context) *echo.HTTPError {
	return CheckNetworkAccess(c, networkWildcard)
}

func CheckNetworkAccess(c echo.Context, networkId string) *echo.HTTPError {
	if !TLS {
		return nil
	}
	if c != nil {
		if r := c.Request(); r != nil {
			if len(r.TLS.PeerCertificates) > 0 {
				var cert = r.TLS.PeerCertificates[0]
				if cert != nil {
					if cert.Subject.CommonName == wildcard ||
						cert.Subject.CommonName == networkWildcard ||
						cert.Subject.CommonName == networkId {
						return nil
					}
					for _, san := range cert.DNSNames {
						if san == wildcard ||
							san == networkWildcard ||
							san == networkId {
							return nil
						}
					}
					log.Printf(
						"Client Cert %s is not authorized for network: %s",
						util.FormatPkixSubject(&cert.Subject), networkId)
					return echo.NewHTTPError(http.StatusForbidden,
						"Client Certificate is not authorized")
				}
			}
		}
	}
	log.Printf("Client Certificate With valid SANs is required for network: %s",
		networkId)
	return echo.NewHTTPError(http.StatusForbidden,
		"Client Certificate With valid SANs is required")
}

func GetNetworkId(c echo.Context) (string, *echo.HTTPError) {
	nid := c.Param("network_id")
	if nid == "" {
		return nid, NetworkIdHttpErr()
	}
	return nid, CheckNetworkAccess(c, nid)
}

func GetTenantID(c echo.Context) (uint64, *echo.HTTPError) {
	oid := c.Param("tenant_id")
	if oid == "" {
		return 0, TenantIdHttpErr()
	}
	intTenantID, err := strconv.ParseUint(oid, 10, 64)
	if err != nil {
		return 0, TenantIdHttpErr()
	}
	return intTenantID, CheckTenantAccess(c)
}

// CheckTenantAccess checks that the context has network wildcard access
// i.e. is admin
func CheckTenantAccess(c echo.Context) *echo.HTTPError {
	if !TLS {
		return nil
	}
	if c != nil {
		if r := c.Request(); r != nil {
			if len(r.TLS.PeerCertificates) > 0 {
				var cert = r.TLS.PeerCertificates[0]
				if cert != nil {
					if cert.Subject.CommonName == wildcard ||
						cert.Subject.CommonName == networkWildcard {
						return nil
					}
					for _, san := range cert.DNSNames {
						if san == wildcard ||
							san == networkWildcard {
							return nil
						}
					}
					log.Printf(
						"Client Cert %s does not have wildcard access",
						util.FormatPkixSubject(&cert.Subject))
					return echo.NewHTTPError(http.StatusForbidden,
						"Client Certificate is not authorized")
				}
			}
		}
	}
	log.Printf("Client Certificate With valid SANs is required for tenant access")
	return echo.NewHTTPError(http.StatusForbidden,
		"Client Certificate With valid SANs is required")
}

// DEPRECATED - use GetGatewayID, and use :gateway_id as path param
func GetLogicalGwId(c echo.Context) (string, *echo.HTTPError) {
	logicalGwId := c.Param("logical_ag_id")
	if logicalGwId == "" {
		return logicalGwId, HttpError(
			fmt.Errorf("Invalid/Missing Gateway ID"),
			http.StatusBadRequest)
	}
	return logicalGwId, nil
}

// DEPRECATED - use GetNetworkAndGatewayIDs, and use :gateway_id as path param
func GetNetworkAndGWID(c echo.Context) (string, string, error) {
	networkID, err := GetNetworkId(c)
	if err != nil {
		return "", "", err
	}
	gatewayID, err := GetLogicalGwId(c)
	if err != nil {
		return "", "", err
	}
	return networkID, gatewayID, nil
}

func GetNetworkAndGatewayIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := GetParamValues(c, "network_id", "gateway_id")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

// GetParamValues returns a list of the value for each param provided in
// `paramNames`. Returns a status bad request HTTP error if any param value
// is blank.
func GetParamValues(c echo.Context, paramNames ...string) ([]string, *echo.HTTPError) {
	ret := make([]string, 0, len(paramNames))
	for _, paramName := range paramNames {
		val := c.Param(paramName)
		if val == "" {
			return []string{}, echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid/missing param %s", paramName))
		}
		ret = append(ret, val)
	}
	return ret, nil
}

func GetOperatorId(c echo.Context) (string, *echo.HTTPError) {
	operId := c.Param("operator_id")
	if operId == "" {
		return operId, HttpError(
			fmt.Errorf("Invalid/Missing Operator ID"),
			http.StatusBadRequest)
	}
	return operId, nil
}

func NetworkIdHttpErr() *echo.HTTPError {
	return HttpError(fmt.Errorf("Missing Network ID"), http.StatusBadRequest)
}

func TenantIdHttpErr() *echo.HTTPError {
	return HttpError(fmt.Errorf("Missing Tenant ID"), http.StatusBadRequest)
}
