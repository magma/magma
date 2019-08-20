/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

package collection

import (
	"errors"

	dto "github.com/prometheus/client_model/go"
)

func getLabelNames(metricFamily *dto.MetricFamily) ([]string, error) {
	metricArr := metricFamily.GetMetric()
	if len(metricArr) == 0 {
		return nil, errors.New("No metric samples available to get labels")
	}
	sample := metricArr[0]
	labelPairArr := sample.GetLabel()
	labelNameArr := make([]string, len(labelPairArr))
	for i, labelPair := range labelPairArr {
		labelNameArr[i] = *labelPair.Name
	}
	return labelNameArr, nil
}

func getLabelVals(metric *dto.Metric) []string {
	labelPairArr := metric.GetLabel()
	labelValArr := make([]string, len(labelPairArr))
	for i, labelPair := range labelPairArr {
		labelValArr[i] = labelPair.GetValue()
	}
	return labelValArr
}
