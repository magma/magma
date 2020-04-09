/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package service_health

// ServiceHealth defines an interface to fetch unhealthy services and enable
// functionality necessary for promotion/demotions of the gateway.
type ServiceHealth interface {
	// GetUnhealthyServices return a list of services found to be in an
	// unhealthy state.
	GetUnhealthyServices() ([]string, error)

	// Enable allows the enabling of service level functionality for the
	// provided service. It is up to implementors to determine
	// the specific functionality.
	Enable(service string) error
}
