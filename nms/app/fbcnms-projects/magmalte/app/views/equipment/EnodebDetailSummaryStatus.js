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
import type {KPIRows} from '../../components/KPIGrid';

import Card from '@material-ui/core/Card';
import EnodebContext from '../../components/context/EnodebContext';
import KPIGrid from '../../components/KPIGrid';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';

import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useContext} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

export function EnodebSummary() {
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);

  const kpiData: KPIRows[] = [
    [
      {
        category: 'eNodeB Serial Number',
        value: enodebSerial,
        statusCircle: false,
      },
    ],
  ];
  return (
    <Card elevation={0}>
      <KPIGrid data={kpiData} />
    </Card>
  );
}

export function EnodebStatus() {
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];

  const isEnbHealthy = isEnodebHealthy(enbInfo);

  const kpiData: KPIRows[] = [
    [
      {
        category: 'Health',
        value: isEnbHealthy ? 'Good' : 'Bad',
        statusCircle: true,
        status: isEnbHealthy,
      },
      {
        category: 'Transmit Enabled',
        value: enbInfo.enb.config.transmit_enabled ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: enbInfo.enb.config.transmit_enabled,
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: enbInfo.enb_state.reporting_gateway_id ?? '',
        statusCircle: true,
        status: enbInfo.enb_state.enodeb_connected,
      },
      {
        category: 'Mme Connected',
        value: enbInfo.enb_state.mme_connected ? 'Connected' : 'Disconnected',
        statusCircle: false,
        status: enbInfo.enb_state.mme_connected,
      },
    ],
  ];
  return <KPIGrid data={kpiData} />;
}
