/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package ocstats

import (
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// HTTPServerResponseCountByStatusAndPath is an additional view for server response status code and path.
var HTTPServerResponseCountByStatusAndPath = &view.View{
	Name:        "opencensus.io/http/server/response_count_by_status_code_path",
	Description: "Server response count by status code and path",
	TagKeys:     []tag.Key{ochttp.StatusCode, ochttp.KeyServerRoute},
	Measure:     ochttp.ServerLatency,
	Aggregation: view.Count(),
}
