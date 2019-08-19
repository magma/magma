/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Section} from '../layout/Section';

import BarChartIcon from '@material-ui/icons/BarChart';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import Configure from '../network/Configure';
import Enodebs from './Enodebs';
import Gateways from '../Gateways';
import Insights from '../insights/Insights';
import Metrics from '../insights/Metrics';
import PeopleIcon from '@material-ui/icons/People';
import PublicIcon from '@material-ui/icons/Public';
import React from 'react';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import Subscribers from '../Subscribers';

export function getLteSections(): Section[] {
  return [
    {
      path: 'map',
      label: 'Map',
      icon: <PublicIcon />,
      component: Insights,
    },
    {
      path: 'metrics',
      label: 'Metrics',
      icon: <BarChartIcon />,
      component: Metrics,
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
      component: Configure,
    },
  ];
}
