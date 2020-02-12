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
	srvRegistry "magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

const ServiceName = "METRICSD"

// getMetricsdClient is a utility function to get a RPC connection to the
// metricsd service
func getMetricsdClient() (protos.MetricsControllerClient, error) {
	conn, err := srvRegistry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewMetricsControllerClient(conn), err
}

func PushMetrics(metrics protos.PushedMetricsContainer) error {
	md, err := getMetricsdClient()
	if err != nil {
		return err
	}
	_, err = md.Push(context.Background(), &metrics)
	return err
}
