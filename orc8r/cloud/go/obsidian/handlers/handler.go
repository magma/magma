/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers provides common glue & functionality for all API handlers
// implemented by "magma/handlers/*" packages
package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"magma/orc8r/cloud/go/obsidian/config"
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
	URL_SEP                   = "/"
	MAGMA_URL_ROOT            = "magma"
	MAGMA_NETWORKS_URL_PART   = "networks"
	MAGMA_OPERATORS_URL_PART  = "operators"
	MAGMA_CHANNELS_URL_PART   = "channels"
	MAGMA_PROMETHEUS_URL_PART = "prometheus"
	// "/magma"
	REST_ROOT = URL_SEP + MAGMA_URL_ROOT
	// "/magma/networks"
	NETWORKS_ROOT = REST_ROOT + URL_SEP + MAGMA_NETWORKS_URL_PART
	// "/magma/operators"
	OPERATORS_ROOT = REST_ROOT + URL_SEP + MAGMA_OPERATORS_URL_PART
	// "/magma/channels"
	CHANNELS_ROOT = REST_ROOT + URL_SEP + MAGMA_CHANNELS_URL_PART
	// "/magma/network/{network_id}/prometheus
	PROMETHEUS_ROOT  = REST_ROOT + URL_SEP + "networks" + URL_SEP + ":network_id" + URL_SEP + MAGMA_PROMETHEUS_URL_PART
	WILDCARD         = "*"
	NETWORK_WILDCARD = "N*"
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

func register(registry handlerRegistry, path string, handler echo.HandlerFunc) error {
	_, registered := registry[path]
	if registered {
		return fmt.Errorf("HandlerFunc[s] already registered for path: %q", path)
	}
	registry[path] = handler
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
			err := register(registries[method], handler.Path, handler.HandlerFunc)
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

func CheckNetworkAccess(c echo.Context, networkId string) *echo.HTTPError {
	if !config.TLS {
		return nil
	}
	if c != nil {
		if r := c.Request(); r != nil {
			if len(r.TLS.PeerCertificates) > 0 {
				var cert = r.TLS.PeerCertificates[0]
				if cert != nil {
					if cert.Subject.CommonName == WILDCARD ||
						cert.Subject.CommonName == NETWORK_WILDCARD ||
						cert.Subject.CommonName == networkId {
						return nil
					}
					for _, san := range cert.DNSNames {
						if san == WILDCARD ||
							san == NETWORK_WILDCARD ||
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

func GetLogicalGwId(c echo.Context) (string, *echo.HTTPError) {
	logicalGwId := c.Param("logical_ag_id")
	if logicalGwId == "" {
		return logicalGwId, HttpError(
			fmt.Errorf("Invalid/Missing Gateway ID"),
			http.StatusBadRequest)
	}
	return logicalGwId, nil
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
