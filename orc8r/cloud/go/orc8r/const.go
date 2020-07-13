/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package orc8r

const (
	ModuleName = "orc8r"

	NetworkFeaturesConfig   = "orc8r_features"
	MagmadGatewayType       = "magmad_gateway"
	AccessGatewayRecordType = "access_gateway_record"
	GatewayStateType        = "gw_state"
	DirectoryRecordType     = "directory_record"

	UpgradeTierEntityType           = "upgrade_tier"
	UpgradeReleaseChannelEntityType = "upgrade_release_channel"

	DnsdNetworkType = "dnsd_network"

	MetricsExporterLabel = "orc8r.io/metrics_exporter"
	StateIndexerLabel    = "orc8r.io/state_indexer"
	StreamProviderLabel  = "orc8r.io/stream_provider"

	StateIndexerVersionAnnotation   = "orc8r.io/state_indexer_version"
	StateIndexerTypesAnnotation     = "orc8r.io/state_indexer_types"
	StreamProviderStreamsAnnotation = "orc8r.io/stream_provider_streams"

	AnnotationListSeparator = ","
)
