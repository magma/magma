/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package metricsd

// Constants to represent the keys in the metricsd.yml config file
const (
	Profile                 = "profile"
	PrometheusPushAddresses = "prometheusPushAddresses"
	PrometheusQueryAddress  = "prometheusQueryAddress"

	PrometheusConfigServiceURL   = "prometheusConfigServiceURL"
	AlertmanagerConfigServiceURL = "alertmanagerConfigServiceURL"
	AlertmanagerApiURL           = "alertmanagerApiURL"
)
