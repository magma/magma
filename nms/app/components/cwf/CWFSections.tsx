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
 */
import type {SectionsConfigs} from '../layout/Section';

import AlarmIcon from '@material-ui/icons/Alarm';
import Alarms from '../../views/alarms/legacy/Alarms';
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
    {
      path: 'alerts',
      label: 'Alerts',
      icon: <AlarmIcon />,
      component: Alarms,
    },
  ];

  return [
    'gateways', // landing path
    sections,
  ];
}
