// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orc8r

import (
	"crypto/tls"
	"net/http"

	"github.com/pkg/errors"
)

type Config struct {
	Hostname   string `env:"API_HOSTNAME" long:"api-hostname" description:"the api host for orchestrator"`
	Cert       string `env:"API_CERT" long:"api-cert" description:"the cert for connecting to orchestrator api"`
	PrivateKey string `env:"API_PRIVATE_KEY" long:"api-private-key" description:"the private key for connecting to orchestrator api"`
}

type Client struct {
	*http.Client
	Hostname string
}

func (config Config) Build() (*Client, error) {
	cert, err := tls.LoadX509KeyPair(config.Cert, config.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load orc8r certificates")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	client := Client{
		Hostname: config.Hostname,
		Client:   httpClient,
	}

	return &client, nil
}
