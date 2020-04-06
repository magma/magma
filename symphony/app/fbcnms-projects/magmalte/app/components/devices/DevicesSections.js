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

import AlarmIcon from '@material-ui/icons/Alarm';
import Alarms from '@fbcnms/magmalte/app/components/insights/Alarms/Alarms';
import CustomMetrics from '../insights/CustomMetrics';
import DeviceHub from '@material-ui/icons/DeviceHub';
import DevicesAgents from './DevicesAgents';
import DevicesStatusTable from './DevicesStatusTable';
import React from 'react';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import ViewListIcon from '@material-ui/icons/ViewList';

export function getDevicesSections(alertsEnabled: boolean): SectionsConfigs {
  const sections = [
    {
      path: 'devices',
      label: 'Devices',
      icon: <ViewListIcon />,
      component: DevicesStatusTable,
    },
    {
      path: 'metrics',
      label: 'Metrics',
      icon: <ShowChartIcon />,
      component: CustomMetrics,
    },
    {
      path: 'agents',
      label: 'Agents',
      icon: <DeviceHub />,
      component: DevicesAgents,
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
    'devices', // landing path
    sections,
  ];
}
