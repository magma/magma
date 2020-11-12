/*
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
import type {DataRows} from '../../components/DataGrid';

import DataGrid from '../../components/DataGrid';
import EnodebContext from '../../components/context/EnodebContext';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';

import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useContext} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

export function EnodebSummary() {
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);

  const kpiData: DataRows[] = [
    [
      {
        category: 'eNodeB Serial Number',
        value: enodebSerial,
      },
    ],
  ];
  return <DataGrid data={kpiData} />;
}

export function EnodebStatus() {
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];

  const isEnbHealthy = isEnodebHealthy(enbInfo);

  const kpiData: DataRows[] = [
    [
      {
        category: 'eNodeB Externally Managed',
        value:
          enbInfo.enb?.enodeb_config?.config_type === 'MANAGED'
            ? 'False'
            : 'True',
      },
      {
        category: 'Health',
        value: isEnbHealthy ? 'Good' : 'Bad',
        statusCircle: true,
        status: isEnbHealthy,
        tooltip: isEnbHealthy
          ? 'eNodeB transmit config and status match'
          : 'mismatch in eNodeB transmit config and status',
      },
      {
        category: 'Transmit Enabled',
        value: enbInfo.enb.enodeb_config?.managed_config?.transmit_enabled
          ? 'Enabled'
          : 'Disabled',
        statusCircle: true,
        status: enbInfo.enb.enodeb_config?.managed_config?.transmit_enabled,
        tooltip: 'current transmit configuration on the eNodeB',
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: enbInfo.enb_state.reporting_gateway_id ?? 'Not Available',
        statusCircle: true,
        status: enbInfo.enb_state.enodeb_connected,
      },
      {
        category: 'Mme Connected',
        value: enbInfo.enb_state.mme_connected ? 'Connected' : 'Disconnected',
        status: enbInfo.enb_state.mme_connected,
      },
      {
        category: 'IP Address',
        value: enbInfo.enb_state.ip_address ?? 'Not Available',
      },
    ],
  ];
  return <DataGrid data={kpiData} />;
}
