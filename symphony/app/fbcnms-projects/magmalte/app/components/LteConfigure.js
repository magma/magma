/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Configure from './network/Configure';
import DataPlanConfig from './network/DataPlanConfig';
import NetworkConfig from './network/NetworkConfig';
import PoliciesConfig from './network/PoliciesConfig';
import React from 'react';
import UpgradeConfig from './network/UpgradeConfig';

export default function LteConfigure() {
  const tabs = [
    {
      component: DataPlanConfig,
      label: 'Data Plans',
      path: 'dataplans',
    },
    {
      component: NetworkConfig,
      label: 'Network Configuration',
      path: 'network',
    },
    {
      component: UpgradeConfig,
      label: 'Upgrades',
      path: 'upgrades',
    },
    {
      component: PoliciesConfig,
      label: 'Policies',
      path: 'policies',
    },
  ];
  return <Configure tabRoutes={tabs} />;
}
