/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {SectionsConfigs} from '../layout/Section';

import AlarmIcon from '@material-ui/icons/Alarm';
import Alarms from '@fbcnms/ui/insights/Alarms/Alarms';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import Enodebs from './Enodebs';
import EquipmentDashboard from '../../views/equipment/EquipmentDashboard';
import Gateways from '../Gateways';
import Insights from '@fbcnms/ui/insights/Insights';
import ListIcon from '@material-ui/icons/List';
import Logs from '@fbcnms/ui/insights/Logs/Logs';
import LteConfigure from '../LteConfigure';
import LteDashboard from './LteDashboard';
import LteMetrics from './LteMetrics';
import NetworkCheckIcon from '@material-ui/icons/NetworkCheck';
import NetworkDashboard from '../../views/network/NetworkDashboard';
import PeopleIcon from '@material-ui/icons/People';
import PublicIcon from '@material-ui/icons/Public';
import React from 'react';
import RouterIcon from '@material-ui/icons/Router';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import SubscriberDashboard from '../../views/subscriber/SubscriberOverview';
import Subscribers from '../Subscribers';

export function getLteSections(
  alertsEnabled: boolean,
  logsEnabled: boolean,
  dashboardV2Enabled: boolean,
): SectionsConfigs {
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
  if (logsEnabled) {
    sections[1].splice(2, 0, {
      path: 'logs',
      label: 'Logs',
      icon: <ListIcon />,
      component: Logs,
    });
  }
  if (alertsEnabled) {
    sections[1].splice(2, 0, {
      path: 'alerts',
      label: 'Alerts',
      icon: <AlarmIcon />,
      component: Alarms,
    });
  }
  if (dashboardV2Enabled) {
    sections[1].splice(2, 0, {
      path: 'dashboard',
      label: 'Dashboard',
      icon: <ShowChartIcon />,
      component: LteDashboard,
    });
    sections[1].splice(3, 0, {
      path: 'equipment',
      label: 'EquipmentV2',
      icon: <RouterIcon />,
      component: EquipmentDashboard,
    });
    sections[1].splice(4, 0, {
      path: 'network',
      label: 'NetworkV2',
      icon: <NetworkCheckIcon />,
      component: NetworkDashboard,
    });
    sections[1].splice(5, 0, {
      path: 'subscriberv2',
      label: 'SubscriberV2',
      icon: <PeopleIcon />,
      component: SubscriberDashboard,
    });
  }
  return sections;
}
