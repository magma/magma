/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {SectionsConfigs} from '@fbcnms/magmalte/app/components/layout/Section';

import CellWifiIcon from '@material-ui/icons/CellWifi';
import FEGConfigure from './FEGConfigure';
import FEGGateways from './FEGGateways';
import React from 'react';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';

export function getFEGSections(): SectionsConfigs {
  const sections = [
    {
      path: 'gateways',
      label: 'Gateways',
      icon: <CellWifiIcon />,
      component: FEGGateways,
    },
    {
      path: 'configure',
      label: 'Configure',
      icon: <SettingsCellIcon />,
      component: FEGConfigure,
    },
  ];

  return [
    'gateways', // landing path
    sections,
  ];
}
