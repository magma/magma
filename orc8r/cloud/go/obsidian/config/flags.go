/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package config 'owns' all configuration parameters settable via command line
// application flags
package config

var (
	TLS                bool
	Port               int
	MagmadDBDriver     string
	MagmadDBSource     string
	ServerCertPemPath  string
	ServerKeyPemPath   string
	ClientCAPoolPath   string
	AllowAnyClientCert bool
	StaticFolder       string
)
