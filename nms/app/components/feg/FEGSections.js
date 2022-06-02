/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import type {SectionsConfigs} from '../layout/Section';

import AlarmIcon from '@material-ui/icons/Alarm';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AlarmsDashboard from '../../views/alarms/AlarmsDashboard';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DashboardIcon from '@material-ui/icons/Dashboard';
import FEGConfigure from './FEGConfigure';
import FEGDashboard from '../../views/dashboard/feg/FEGDashboard';
import FEGEquipmentDashboard from '../../views/equipment/FEGEquipmentDashboard';
import FEGGateways from './FEGGateways';
import FEGMetrics from './FEGMetrics';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import FEGNetworkDashboard from '../../views/network/FEGNetworkDashboard';
import NetworkCheckIcon from '@material-ui/icons/NetworkCheck';
import React from 'react';
import RouterIcon from '@material-ui/icons/Router';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import ShowChartIcon from '@material-ui/icons/ShowChart';

export function getFEGSections(): SectionsConfigs {
  const sections = [
    {
      path: 'dashboard',
      label: 'Dashboard',
      icon: <DashboardIcon />,
      component: FEGDashboard,
    },
    {
      path: 'gateways',
      label: 'Gateways',
      icon: <CellWifiIcon />,
      component: FEGGateways,
    },
    {
      path: 'network',
      label: 'Network',
      icon: <NetworkCheckIcon />,
      component: FEGNetworkDashboard,
    },
    {
      path: 'equipment',
      label: 'Equipment',
      icon: <RouterIcon />,
      component: FEGEquipmentDashboard,
    },
    {
      path: 'configure',
      label: 'Configure',
      icon: <SettingsCellIcon />,
      component: FEGConfigure,
    },
    {
      path: 'alerts',
      label: 'Alerts',
      icon: <AlarmIcon />,
      component: AlarmsDashboard,
    },
    {
      path: 'metrics',
      label: 'Metrics',
      icon: <ShowChartIcon />,
      component: FEGMetrics,
    },
  ];

  return [
    'dashboard', // landing path
    sections,
  ];
}
