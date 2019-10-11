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

import AlarmIcon from '@material-ui/icons/Alarm';
import Alarms from '@fbcnms/magmalte/app/components/insights/Alarms/Alarms';
import EditIcon from '@material-ui/icons/Edit';
import PublicIcon from '@material-ui/icons/Public';
import React from 'react';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import ViewListIcon from '@material-ui/icons/ViewList';
import WifiConfig from './WifiConfig';
import WifiMap from './WifiMap';
import WifiMeshesDevicesTable from './WifiMeshesDevicesTable';
import WifiMetrics from './WifiMetrics';

export function getMeshSections(alertsEnabled: boolean): SectionsConfigs {
  const sections = [
    {
      path: 'map',
      label: 'Map',
      icon: <PublicIcon />,
      component: WifiMap,
    },
    {
      path: 'metrics',
      label: 'Metrics',
      icon: <ShowChartIcon />,
      component: WifiMetrics,
    },
    {
      path: 'devices',
      label: 'Devices',
      icon: <ViewListIcon />,
      component: WifiMeshesDevicesTable,
    },
    {
      path: 'configure',
      label: 'Configure',
      icon: <EditIcon />,
      component: WifiConfig,
    },
  ];

  if (alertsEnabled) {
    sections.splice(2, 0, {
      path: 'alerts',
      label: 'Alerts',
      icon: <AlarmIcon />,
      component: Alarms,
    });
  }

  return [
    'map', // landing path
    sections,
  ];
}
