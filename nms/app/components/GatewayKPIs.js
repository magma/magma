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

import type {DataRows} from './DataGrid';
import type {lte_gateway} from '../../generated/MagmaAPIBindings';

import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataGrid from './DataGrid';
// $FlowFixMe migrated to typescript
import GatewayContext from './context/GatewayContext';
import React from 'react';

import {useContext} from 'react';

function gatewayStatus(gatewaySt: {[string]: lte_gateway}): [number, number] {
  let upCount = 0;
  let downCount = 0;
  Object.keys(gatewaySt)
    .map((k: string) => gatewaySt[k])
    .filter((g: lte_gateway) => g.cellular && g.id)
    .map(function (gateway: lte_gateway) {
      gateway.checked_in_recently ? upCount++ : downCount++;
    });
  return [upCount, downCount];
}

export default function GatewayKPIs() {
  const gwCtx = useContext(GatewayContext);
  const [upCount, downCount] = gatewayStatus(gwCtx.state);

  const data: DataRows[] = [
    [
      {
        icon: CellWifiIcon,
        value: 'Gateways',
      },
      {
        category: 'Severe Events',
        value: 0,
        tooltip: 'Severe Events reported by the gateway',
      },
      {
        category: 'Connected',
        value: upCount || 0,
        tooltip: 'Number of gateways checked in recently',
      },
      {
        category: 'Disconnected',
        value: downCount || 0,
        tooltip: 'Number of gateways not checked in recently',
      },
    ],
  ];

  return <DataGrid data={data} />;
}
