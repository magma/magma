/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package reverse_proxy

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/access/tests"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	RegisterNetworkV1 = "/magma/v1/networks"
)

func TestReverseProxyMiddleware(t *testing.T) {
	e := startTestServer(t)
	listener := tests.WaitForTestServer(t, e)
	if listener == nil {
		return // WaitForTestServer should have 'logged' error already
	}
	_, port, err := net.SplitHostPort(listener.Addr().String())
	assert.NoError(t, err)

	obsidian.Port, err = strconv.Atoi(port)
	assert.NoError(t, err)

	urlPrefix := "http://" + listener.Addr().String()

	// Explicitly set request IP to non-localhost to ensure reverse proxy
	// middleware is tested
	s, err := sendRequest(
		"GET", // READ
		urlPrefix+RegisterNetworkV1,
		"10.10.10.10",
	)
	assert.NoError(t, err)
	assert.Equal(t, 200, s)

	// Test localhost request IP
	s, err = sendRequest(
		"GET", // READ
		urlPrefix+RegisterNetworkV1,
		"127.0.0.1",
	)
	assert.NoError(t, err)
	assert.Equal(t, 200, s)
}

func startTestServer(t *testing.T) *echo.Echo {
	e := echo.New()
	assert.NotNil(t, e)

	e.GET(RegisterNetworkV1, func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	})
	e.Use(ReverseProxy) // inject obsidian reverse proxy middleware

	go func(t *testing.T) {
		assert.NoError(t, e.Start(""))
	}(t)

	return e
}

func sendRequest(method string, url string, addr string) (int, error) {
	var body io.Reader = nil
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(echo.HeaderXRealIP, addr)
	var client = &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	return response.StatusCode, err
}
