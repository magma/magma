/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package definitions defines consts, vars & types common to gateway & cloud
package definitions

const (
	MconfigStreamName = "configs"

	// service names
	ControlProxyServiceName = "control_proxy"
	DispatcherServiceName   = "dispatcher"
	MagmadServiceName       = "magmad"
	MetricsdServiceName     = "metricsd"
	StateServiceName        = "state"
	StreamerServiceName     = "streamer"
)
