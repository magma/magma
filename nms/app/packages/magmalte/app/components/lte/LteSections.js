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

import DashboardIcon from '@material-ui/icons/Dashboard';
import EquipmentDashboard from '../../views/equipment/EquipmentDashboard';
import LteDashboard from './LteDashboard';
import LteMetrics from './LteMetrics';
import NetworkCheckIcon from '@material-ui/icons/NetworkCheck';
import NetworkDashboard from '../../views/network/NetworkDashboard';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import RouterIcon from '@material-ui/icons/Router';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import SubscriberDashboard from '../../views/subscriber/SubscriberOverview';
import TrafficDashboard from '../../views/traffic/TrafficOverview';
import WifiTetheringIcon from '@material-ui/icons/WifiTethering';

export function getLteSections(): SectionsConfigs {
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
        label: 'EquipmentV2',
        icon: <RouterIcon />,
        component: EquipmentDashboard,
      },
      {
        path: 'network',
        label: 'NetworkV2',
        icon: <NetworkCheckIcon />,
        component: NetworkDashboard,
      },
      {
        path: 'subscribers',
        label: 'SubscriberV2',
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
        path: 'metrics',
        label: 'Metrics',
        icon: <ShowChartIcon />,
        component: LteMetrics,
      },
    ],
  ];
  return sections;
}
