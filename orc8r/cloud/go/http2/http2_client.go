/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package http2 contains a minimal implementation of non-TLS http/2 server
// and client
package http2

import (
	"crypto/tls"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

// H2CClient is a http2 client supports non-SSL only
type H2CClient struct {
	*http.Client
}

// NewH2CClient creates a new h2cclient.
func NewH2CClient() *H2CClient {
	return &H2CClient{&http.Client{
		// Skip TLS dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				// dial the addr from url,
				// or :80 if no legitimate ip retrieved from url
				return net.Dial(netw, addr)
			},
		},
	}}
}
