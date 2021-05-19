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

import * as React from 'react';
import AlarmIcon from '@material-ui/icons/Alarm';
import AlarmsDashboard from '../../views/alarms/AlarmsDashboard';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DashboardIcon from '@material-ui/icons/Dashboard';
import Enodebs from './Enodebs';
import EquipmentDashboard from '../../views/equipment/EquipmentDashboard';
import Gateways from '../Gateways';
import Insights from '@fbcnms/ui/insights/Insights';
import LineStyleIcon from '@material-ui/icons/LineStyle';
import ListIcon from '@material-ui/icons/List';
import Logs from '@fbcnms/ui/insights/Logs/Logs';
import LteConfigure from '../LteConfigure';
import LteDashboard from '../../views/dashboard/lte/LteDashboard';
import LteMetrics from './LteMetrics';
import NetworkCheckIcon from '@material-ui/icons/NetworkCheck';
import NetworkDashboard from '../../views/network/NetworkDashboard';
import PeopleIcon from '@material-ui/icons/People';
import PublicIcon from '@material-ui/icons/Public';
import RouterIcon from '@material-ui/icons/Router';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import SubscriberDashboard from '../../views/subscriber/SubscriberOverview';
import Subscribers from '../Subscribers';
import TracingDashboard from '../../views/tracing/TracingDashboard';
import TrafficDashboard from '../../views/traffic/TrafficOverview';
import WifiTetheringIcon from '@material-ui/icons/WifiTethering';

export function getLteSections(
  alertsEnabled: boolean,
  logsEnabled: boolean,
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
      {
        path: 'alerts',
        label: 'Alerts',
        icon: <AlarmIcon />,
        component: AlarmsDashboard,
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
  return sections;
}

export function getLteSectionsV2(alertsEnabled: boolean): SectionsConfigs {
  const sections = [
    'dashboard', // landing path
    [
      {
        path: 'dashboard',
        label: 'Dashboard',
        icon: <DashboardIcon />,
        component: LteDashboard,
      },
      {
        path: 'equipment',
        label: 'Equipment',
        icon: <RouterIcon />,
        component: EquipmentDashboard,
      },
      {
        path: 'network',
        label: 'Network',
        icon: <NetworkCheckIcon />,
        component: NetworkDashboard,
      },
      {
        path: 'subscribers',
        label: 'Subscriber',
        icon: <PeopleIcon />,
        component: SubscriberDashboard,
      },
      {
        path: 'traffic',
        label: 'Traffic',
        icon: <WifiTetheringIcon />,
        component: TrafficDashboard,
      },
      {
        path: 'tracing',
        label: 'Call Tracing',
        icon: <LineStyleIcon />,
        component: TracingDashboard,
      },
      {
        path: 'metrics',
        label: 'Metrics',
        icon: <ShowChartIcon />,
        component: LteMetrics,
      },
    ],
  ];
  if (alertsEnabled) {
    sections[1].push({
      path: 'alerts',
      label: 'Alerts',
      icon: <AlarmIcon />,
      component: AlarmsDashboard,
    });
  }
  return sections;
}
