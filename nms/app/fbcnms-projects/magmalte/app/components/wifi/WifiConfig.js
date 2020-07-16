/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Configure from '@fbcnms/magmalte/app/components/network/Configure';
import React from 'react';
import WifiNetworkConfig from './WifiNetworkConfig';

export default function WifiConfig() {
  const tabs = [
    {
      component: WifiNetworkConfig,
      label: 'Network Configuration',
      path: '',
    },
  ];
  return <Configure tabRoutes={tabs} />;
}
