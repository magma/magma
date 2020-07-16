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

import CWFConfigure from './CWFConfigure';
import CWFGateways from './CWFGateways';
import CWFMetrics from './CWFMetrics';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import React from 'react';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import ShowChartIcon from '@material-ui/icons/ShowChart';

export function getCWFSections(): SectionsConfigs {
  const sections = [
    {
      path: 'gateways',
      label: 'Gateways',
      icon: <CellWifiIcon />,
      component: CWFGateways,
    },
    {
      path: 'configure',
      label: 'Configure',
      icon: <SettingsCellIcon />,
      component: CWFConfigure,
    },
    {
      path: 'metrics',
      label: 'Metrics',
      icon: <ShowChartIcon />,
      component: CWFMetrics,
    },
  ];

  return [
    'gateways', // landing path
    sections,
  ];
}
