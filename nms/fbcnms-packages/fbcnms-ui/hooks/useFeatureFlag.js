/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Context} from 'react';

export type FeatureID =
  | 'lte_network_metrics'
  | 'sso_example_feature'
  | 'audit_log_example_feature'
  | 'third_party_devices'
  | 'network_topology'
  | 'prometheus_metrics_database';

import {useContext} from 'react';

type ContextType = {enabledFeatures: FeatureID[]} & {[string]: any};

export default function(
  appContext: Context<ContextType>,
  featureId: FeatureID,
): boolean {
  const {enabledFeatures} = useContext(appContext);
  return enabledFeatures.indexOf(featureId) !== -1;
}
