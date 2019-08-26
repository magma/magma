/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian

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

// configs
var (
	TLS                bool
	Port               int
	ServerCertPemPath  string
	ServerKeyPemPath   string
	ClientCAPoolPath   string
	AllowAnyClientCert bool
	StaticFolder       string
)
