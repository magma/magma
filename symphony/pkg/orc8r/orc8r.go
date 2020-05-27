// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orc8r

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Config configures orc8r http client.
type Config struct {
	Host string
	Cert string
	PKey string
}

// NewClient returns an http client to use for orc8r communication.
func NewClient(cfg Config) (*http.Client, error) {
	cert, err := tls.LoadX509KeyPair(cfg.Cert, cfg.PKey)
	if err != nil {
		return nil, fmt.Errorf("cannot load certificates: %w", err)
	}
	var rt http.RoundTripper = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	rt = &ochttp.Transport{Base: rt}
	rt = &Transport{Base: rt, Host: cfg.Host}
	return &http.Client{Transport: rt}, nil
}

// Transport sends http requests to orc8r host.
type Transport struct {
	Base http.RoundTripper
	Host string
}

// RoundTrip implements http.RoundTripper interface.
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	req := *r
	req.Header = r.Header.Clone()
	req.URL.Host = t.Host
	req.URL.Scheme = "https"
	return t.Base.RoundTrip(r)
}

// AddFlagsVar adds the flags used by this package to the Kingpin application.
func AddFlagsVar(a *kingpin.Application, config *Config) {
	a.Flag("orc8r.host", "orchestrator host").
		Envar("ORC8R_HOST").
		StringVar(&config.Host)
	a.Flag("orc8r.cert", "orchestrator certificate").
		Envar("ORC8R_CERT").
		StringVar(&config.Cert)
	a.Flag("orc8r.pkey", "orchestrator private key").
		Envar("ORC8R_PKEY").
		StringVar(&config.PKey)
}

// AddFlags adds the flags used by this package to the Kingpin application.
func AddFlags(a *kingpin.Application) *Config {
	config := &Config{}
	AddFlagsVar(a, config)
	return config
}
