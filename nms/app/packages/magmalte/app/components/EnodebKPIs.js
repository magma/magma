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
import type {EnodebInfo} from './lte/EnodebUtils';

import DataGrid from './DataGrid';
import EnodebContext from './context/EnodebContext';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';

import {useContext} from 'react';

export default function EnodebKPIs() {
  const ctx = useContext(EnodebContext);
  const [total, transmitting] = enodebStatus(ctx.state.enbInfo);

  const data: DataRows[] = [
    [
      {
        icon: SettingsInputAntennaIcon,
        value: 'eNodeBs',
      },
      {
        category: 'Severe Events',
        value: 0,
      },
      {
        category: 'Total',
        value: total || 0,
      },
      {
        category: 'Transmitting',
        value: transmitting || 0,
      },
    ],
  ];

  return <DataGrid data={data} />;
}

function enodebStatus(enbInfo: {[string]: EnodebInfo}): [number, number] {
  let transmitCnt = 0;
  Object.keys(enbInfo)
    .map((k: string) => enbInfo[k])
    .map((e: EnodebInfo) => {
      if (e.enb_state.rf_tx_on) {
        transmitCnt++;
      }
    });
  return [Object.keys(enbInfo).length, transmitCnt];
}
