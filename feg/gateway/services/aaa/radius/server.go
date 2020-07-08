/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package radius implements AAA server's radius interface for accounting & authentication
package radius

import "magma/feg/cloud/go/protos/mconfig"

// Server - radius server implementation
type Server struct {
	config      *mconfig.RadiusConfig
	authMethods []byte
}

// GetConfig returns server configs
func (s *Server) GetConfig() *mconfig.RadiusConfig {
	if s == nil {
		return defaultConfigs
	}
	return s.config
}

// New returns new radius server
func New(cfg *mconfig.RadiusConfig) *Server {
	return &Server{config: validateConfigs(cfg)}
}
