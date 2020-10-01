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
import Alarms from '@fbcnms/ui/insights/Alarms/Alarms';
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
