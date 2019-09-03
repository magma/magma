/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SectionsConfigs} from '../layout/Section';

import AlarmIcon from '@material-ui/icons/Alarm';
import Alarms from '../insights/Alarms/Alarms';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import Enodebs from './Enodebs';
import Gateways from '../Gateways';
import Insights from '../insights/Insights';
import LteConfigure from '../LteConfigure';
import LteMetrics from './LteMetrics';
import PeopleIcon from '@material-ui/icons/People';
import PublicIcon from '@material-ui/icons/Public';
import React from 'react';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import Subscribers from '../Subscribers';

export function getLteSections(alertsEnabled: boolean): SectionsConfigs {
  const sections = [
    'map', // landing path
    [
      {
        path: 'map',
        label: 'Map',
        icon: <PublicIcon />,
        component: Insights,
      },
      {
        path: 'metrics',
        label: 'Metrics',
        icon: <ShowChartIcon />,
        component: LteMetrics,
      },
      {
        path: 'subscribers',
        label: 'Subscribers',
        icon: <PeopleIcon />,
        component: Subscribers,
      },
      {
        path: 'gateways',
        label: 'Gateways',
        icon: <CellWifiIcon />,
        component: Gateways,
      },
      {
        path: 'enodebs',
        label: 'eNodeB Devices',
        icon: <SettingsInputAntennaIcon />,
        component: Enodebs,
      },
      {
        path: 'configure',
        label: 'Configure',
        icon: <SettingsCellIcon />,
        component: LteConfigure,
      },
    ],
  ];
  if (alertsEnabled) {
    sections[1].push({
      path: 'alerts',
      label: 'Alerts',
      icon: <AlarmIcon />,
      component: Alarms,
    });
  }
  return sections;
}
