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

import type {DataRows} from '../../components/DataGrid';

import DataGrid from '../../components/DataGrid';
import GatewayContext from '../../context/GatewayContext';
import LoadingFiller from '../../components/LoadingFiller';
import MagmaAPI from '../../api/MagmaAPI';
import React, {useContext} from 'react';
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPI';
import {DynamicServices} from '../../components/GatewayUtils';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {useInterval} from '../../hooks';
import {useParams} from 'react-router-dom';

export default function GatewayDetailStatus({refresh}: {refresh: boolean}) {
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const gatewayId: string = nullthrows(params.gatewayId);
  const gatewayContext = useContext(GatewayContext);

  useInterval(
    () => gatewayContext.refetch(gatewayId),
    refresh ? REFRESH_INTERVAL : null,
  );

  const gwInfo = gatewayContext.state[gatewayId];
  let checkInTime;
  if (
    gwInfo.status &&
    gwInfo.status.checkin_time !== undefined &&
    gwInfo.status.checkin_time !== null
  ) {
    checkInTime = new Date(gwInfo.status.checkin_time);
  }

  const startTime = Math.floor(Date.now() / 1000);
  const {response: cpuPercent, isLoading: isCpuPercentLoading} = useMagmaAPI(
    MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet,
    {
      networkId: networkId,
      query: `cpu_percent{gatewayID="${gwInfo.id}", service="magmad"}`,
      start: startTime.toString(),
    },
  );

  if (isCpuPercentLoading) {
    return <LoadingFiller />;
  }

  const logAggregation =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes(DynamicServices.TD_AGENT_BIT);

  const eventAggregation =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes(DynamicServices.EVENTD);

  const cpeMonitoring =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes(DynamicServices.MONITORD);

  const data: Array<DataRows> = [
    [
      {
        category: 'Health',
        value: gwInfo.checked_in_recently ? 'Good' : 'Bad',
        statusCircle: true,
        status: gwInfo.checked_in_recently,
        tooltip: gwInfo.checked_in_recently
          ? 'Gateway checked in recently'
          : 'Gateway has not checked in recently',
      },
      {
        category: 'Last Check in',
        value: checkInTime?.toLocaleString() ?? '-',
        statusCircle: false,
      },
      {
        category: 'CPU Usage',
        value: cpuPercent?.data?.result?.[0]?.values?.[0]?.[1] ?? 'Unknown',
        unit:
          cpuPercent?.data?.result?.[0]?.values?.[0]?.[1] ?? false ? '%' : '',
        statusCircle: false,
        tooltip: 'Current Gateway CPU %',
      },
    ],
    [
      {
        category: 'Event Aggregation',
        value: eventAggregation ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: eventAggregation,
      },
      {
        category: 'Log Aggregation',
        value: logAggregation ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: logAggregation,
      },
      {
        category: 'CPE Monitoring',
        value: cpeMonitoring ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: cpeMonitoring,
      },
    ],
  ];
  return <DataGrid data={data} />;
}
