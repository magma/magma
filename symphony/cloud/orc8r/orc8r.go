// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orc8r

type Config struct {
	Hostname   string `env:"API_HOSTNAME" long:"api-hostname" description:"the api host for orchestrator"`
	Cert       string `env:"API_CERT" long:"api-cert" description:"the cert for connecting to orchestrator api"`
	PrivateKey string `env:"API_PRIVATE_KEY" long:"api-private-key" description:"the private key for connecting to orchestrator api"`
}
