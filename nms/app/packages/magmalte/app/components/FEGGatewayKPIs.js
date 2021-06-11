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

import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataGrid from './DataGrid';
import FEGGatewayContext from './context/FEGGatewayContext';
import React from 'react';

import {useContext} from 'react';

export default function FEGGatewayKPIs() {
  const gwCtx = useContext(FEGGatewayContext);
  //gwCtx not used now but will be used
  gwCtx.state;
  //TODO: Get the actual message counts.
  const messageCounts = 0;
  const data: DataRows[] = [
    [
      {
        icon: CellWifiIcon,
        value: 'Federation Gateway',
      },
      {
        category: 'Message counts',
        value: messageCounts,
        tooltip: 'Number of messages reported by the gateway',
      },
    ],
  ];

  return <DataGrid data={data} />;
}
