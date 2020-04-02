/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package health_client

import (
	"context"

	magmav1alpha1 "magma/cwf/k8s/cwf_operator/pkg/apis/magma/v1alpha1"
)

// TODO: Remove and replace with grpc proto client
type HealthServiceClient interface {
	GetHealthStatus(ctx context.Context, serviceAddr string) (magmav1alpha1.CarrierWifiAccessGatewayHealthCondition, error)

	Enable(ctx context.Context, serviceAddr string) error
}

type HealthClient struct{}

func (c *HealthClient) GetHealthStatus(ctx context.Context, serviceAddr string) (magmav1alpha1.CarrierWifiAccessGatewayHealthCondition, error) {
	return magmav1alpha1.Healthy, nil
}

func (c *HealthClient) Enable(ctx context.Context, serviceAddr string) error {
	return nil
}
