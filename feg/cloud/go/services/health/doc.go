/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package health contains the health service, a feg-orchestrator microservice which
// manages feg health and HA clusters
package health

const (
	ServiceName       = "HEALTH"
	DBTableName       = "health"
	HealthStatusType  = "health_status"
	ClusterStatusType = "cluster_status"
)
