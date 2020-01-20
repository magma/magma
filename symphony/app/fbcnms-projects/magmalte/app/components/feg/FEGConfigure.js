/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Configure from '../network/Configure';
import React from 'react';
import UpgradeConfig from '../network/UpgradeConfig';

export default function CWFConfigure() {
  const tabs = [
    {
      component: UpgradeConfig,
      label: 'Upgrades',
      path: 'upgrades',
    },
  ];
  return <Configure tabRoutes={tabs} />;
}
