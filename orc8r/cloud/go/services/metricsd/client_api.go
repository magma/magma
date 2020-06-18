/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package metricsd

import (
	"context"

	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	service_registry "magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// PushMetrics pushes a set of metrics to the metricsd service.
func PushMetrics(metrics protos.PushedMetricsContainer) error {
	client, err := getMetricsdClient()
	if err != nil {
		return err
	}
	_, err = client.Push(context.Background(), &metrics)
	return err
}

// getMetricsdClient is a utility function to get a RPC connection to the
// metricsd service
func getMetricsdClient() (protos.MetricsControllerClient, error) {
	conn, err := service_registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewMetricsControllerClient(conn), err
}
