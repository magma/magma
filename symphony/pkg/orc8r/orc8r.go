// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orc8r

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"go.opencensus.io/plugin/ochttp"
)

// Config configures orc8r http client.
type Config struct {
	Host string `env:"HOST" long:"host" description:"orchestrator hostname"`
	Cert string `env:"CERT" long:"cert" description:"orchestrator certificate"`
	PKey string `env:"PKEY" long:"pkey" description:"orchestrator private key"`
}

// NewClient returns an http client to use for orc8r communication.
func NewClient(cfg Config) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.PKey)
	if err != nil {
		return nil, fmt.Errorf("cannot load certificates: %w", err)
	}
	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	transport = &ochttp.Transport{Base: transport}
	transport = Transport{Base: transport, Host: cfg.Host}
	return &http.Client{Transport: transport}, nil
}

// Transport sends http requests to orc8r host.
type Transport struct {
	Base http.RoundTripper
	Host string
}

// RoundTrip implements http.RoundTripper interface.
func (t Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	req := *r
	req.Header = r.Header.Clone()
	req.URL.Host = t.Host
	req.URL.Scheme = "https"
	return t.Base.RoundTrip(r)
}
