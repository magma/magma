/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package config 'owns' all configuration parameters settable via command line
// application flags
package config

const (
	Product              = "Obsidian Server"
	Version              = "0.1"
	DefaultPort          = 9081
	DefaultHttpsPort     = 9443
	DefaultServerCert    = "server_cert.pem"
	DefaultServerCertKey = "server_cert.key.pem"
	DefaultClientCAs     = "ca_cert.pem"
	DefaultStaticFolder  = "/var/opt/magma/static"
	StaticURLPrefix      = "/apidocs"
	ServiceName          = "OBSIDIAN"
)
