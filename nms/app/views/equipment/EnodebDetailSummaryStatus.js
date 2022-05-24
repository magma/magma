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

import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
// $FlowFixMe migrated to typescript
import EnodebContext from '../../components/context/EnodebContext';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {
  REFRESH_INTERVAL,
  useRefreshingContext,
} from '../../components/context/RefreshContext';
// $FlowFixMe migrated to typescript
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useContext} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

export function EnodebSummary() {
  const ctx = useContext(EnodebContext);
  const params = useParams();
  const enodebSerial: string = nullthrows(params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const kpiData: DataRows[] = [
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
  const networkId: string = nullthrows(params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();

  // Auto refresh enodeb every 30 seconds
  const state = useRefreshingContext({
    context: EnodebContext,
    networkId: networkId,
    type: 'enodeb',
    interval: REFRESH_INTERVAL,
    id: enodebSerial,
    enqueueSnackbar,
    refresh: refresh,
  });

  // $FlowIgnore
  const enbInfo = state.enbInfo?.[enodebSerial];
  const isEnbHealthy = enbInfo ? isEnodebHealthy(enbInfo) : false;
  const isEnbManaged = enbInfo?.enb?.enodeb_config?.config_type === 'MANAGED';

  const kpiData: DataRows[] = [
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
