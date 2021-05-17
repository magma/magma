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
import type {SectionsConfigs} from '@fbcnms/magmalte/app/components/layout/Section';

import AlarmIcon from '@material-ui/icons/Alarm';
import AlarmsDashboard from '../../views/alarms/AlarmsDashboard';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import FEGConfigure from './FEGConfigure';
import FEGGateways from './FEGGateways';
import FEGMetrics from './FEGMetrics';
import React from 'react';
import SettingsCellIcon from '@material-ui/icons/SettingsCell';
import ShowChartIcon from '@material-ui/icons/ShowChart';

export function getFEGSections(dashboardV2Enabled: boolean): SectionsConfigs {
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

  if (dashboardV2Enabled) {
    // TODO add equipment, policy and subscriber section
  }

  return [
    'gateways', // landing path
    sections,
  ];
}
