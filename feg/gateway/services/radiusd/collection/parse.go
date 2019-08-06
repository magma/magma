/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

package collection

import (
	"errors"
	"fmt"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// ParsePrometheusText parses the HTTP response body from common exporters
// that expose prometheus metrics in a common text format.
// Returns metric families keyed by the metric name.
func ParsePrometheusText(prometheusText string) (map[string]*dto.MetricFamily, error) {
	reader := strings.NewReader(prometheusText)

	parser := expfmt.TextParser{}
	metricFamilies, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing metric families from text: %s", err))
	}

	return metricFamilies, nil
}
