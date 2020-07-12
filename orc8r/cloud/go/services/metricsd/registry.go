/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// File registry.go provides a metrics exporter registry by forwarding calls to
// the service registry.

package metricsd

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/lib/go/registry"
)

// GetMetricsExporters returns all registered metrics exporters.
func GetMetricsExporters() []exporters.Exporter {
	services := registry.FindServices(orc8r.MetricsExporterLabel)

	var ret []exporters.Exporter
	for _, s := range services {
		ret = append(ret, exporters.NewRemoteExporter(s))
	}

	return ret
}
