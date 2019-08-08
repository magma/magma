/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package mconfig provides gateway Go support for cloud managed configuration (mconfig)
// Managed Configs are stored in proto JSON marshaled form in gateway.mconfig by external process (magmad)
// and periodically (MCONFOG_REFRESH_INTERVAL) refreshed by a dedicated routine
//
//go:generate bash -c "MAGMA_MODULES='$MAGMA_ROOT/orc8r $MAGMA_ROOT/lte $MAGMA_ROOT/feg' make -C $MAGMA_ROOT/orc8r/cloud gen"
//
package mconfig

import "time"

const (
	DefaultConfigFileDir        = "/etc/magma"
	DefaultDynamicConfigFileDir = "/var/opt/magma/configs"
	ConfigFileDirEnv            = "MAGMA_CONFIG_LOCATION"
	MconfigFileName             = "gateway.mconfig"
	MconfigRefreshInterval      = time.Second * 120
)
