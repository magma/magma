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
 */

import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
import EnodebContext from '../../context/EnodebContext';
import React from 'react';
import SettingsInputAntennaIcon from '@mui/icons-material/SettingsInputAntenna';
import nullthrows from '../../../shared/util/nullthrows';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useContext} from 'react';
import {useInterval} from '../../hooks';
import {useParams} from 'react-router-dom';
import type {DataRows} from '../../components/DataGrid';

export function EnodebSummary() {
  const ctx = useContext(EnodebContext);
  const params = useParams();
  const enodebSerial: string = nullthrows(params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const kpiData: Array<DataRows> = [
    [
      {
        category: 'eNodeB Serial Number',
        value: enodebSerial,
      },
    ],
  ];
  return (
    <>
      <CardTitleRow icon={SettingsInputAntennaIcon} label={enbInfo.enb.name} />
      <DataGrid data={kpiData} />
    </>
  );
}

export function EnodebStatus({refresh}: {refresh: boolean}) {
  const params = useParams();
  const enodebSerial: string = nullthrows(params.enodebSerial);
  const enodebContext = useContext(EnodebContext);

  useInterval(
    () => enodebContext.refetch(enodebSerial),
    refresh ? REFRESH_INTERVAL : null,
  );

  const enbInfo = enodebContext.state.enbInfo?.[enodebSerial];
  const isEnbHealthy = enbInfo ? isEnodebHealthy(enbInfo) : false;
  const isEnbManaged = enbInfo?.enb?.enodeb_config?.config_type === 'MANAGED';

  const kpiData: Array<DataRows> = [
    [
      {
        category: 'eNodeB Externally Managed',
        value: isEnbManaged ? 'False' : 'True',
      },
      {
        category: 'Health',
        value: isEnbManaged ? (isEnbHealthy ? 'Good' : 'Bad') : '-',
        statusCircle: isEnbManaged,
        status: isEnbHealthy,
        tooltip: isEnbManaged
          ? isEnbHealthy
            ? 'eNodeB transmit config and status match'
            : 'mismatch in eNodeB transmit config and status'
          : 'Health information unavailable on externally managed eNodeBs',
      },
      {
        category: 'Transmit Enabled',
        value: enbInfo?.enb.enodeb_config?.managed_config?.transmit_enabled
          ? 'Enabled'
          : 'Disabled',
        statusCircle: true,
        status: enbInfo?.enb.enodeb_config?.managed_config?.transmit_enabled,
        tooltip: 'current transmit configuration on the eNodeB',
      },
      {
        category: 'Subscribers',
        value: enbInfo?.enb_state?.ues_connected ?? 0,
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: enbInfo?.enb_state.reporting_gateway_id ?? 'Not Available',
        statusCircle: true,
        status: enbInfo?.enb_state.enodeb_connected,
      },
      {
        category: 'Mme Connected',
        value: enbInfo?.enb_state.mme_connected ? 'Connected' : 'Disconnected',
        status: enbInfo?.enb_state.mme_connected,
      },
      {
        category: 'IP Address',
        value: enbInfo?.enb_state.ip_address ?? 'Not Available',
      },
    ],
  ];
  return (
    <>
      <DataGrid data={kpiData} />
    </>
  );
}
