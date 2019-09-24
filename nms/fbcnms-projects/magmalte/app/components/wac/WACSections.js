/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SectionsConfigs} from '@fbcnms/magmalte/app/components/layout/Section';

import LeakAddIcon from '@material-ui/icons/LeakAdd';
import React from 'react';
import WACDevices from './WACDevices';

export function getWACSections(): SectionsConfigs {
  return [
    'devices', // landing path
    [
      {
        path: 'devices',
        label: 'Devices',
        icon: <LeakAddIcon />,
        component: WACDevices,
      },
    ],
  ];
}
