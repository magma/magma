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
import FEGGatewayContext from '../../components/context/FEGGatewayContext';
import LoadingFiller from '../../components/LoadingFiller';
import MagmaAPI from '../../api/MagmaAPI';
import React, {useContext} from 'react';
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPI';
import {
  DynamicServices,
  GatewayTypeEnum,
  HEALTHY_STATUS,
} from '../../components/GatewayUtils';
import {REFRESH_INTERVAL} from '../../components/context/RefreshContext';
import {useInterval} from '../../hooks';
import {useParams} from 'react-router-dom';

export default function GatewayDetailStatus({refresh}: {refresh: boolean}) {
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const gatewayId: string = nullthrows(params.gatewayId);
  const gwCtx = useContext(FEGGatewayContext);

  // Auto refresh gateways every 30 seconds
  useInterval(
    () => gwCtx.refetch(gatewayId),
    refresh ? REFRESH_INTERVAL : null,
  );
  const gwInfo = gwCtx.state[gatewayId] || {};
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

  const logAggregationEnabled =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes(DynamicServices.TD_AGENT_BIT);

  const eventAggregationEnabled =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes(DynamicServices.EVENTD);

  const cpeMonitoringEnabled =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes(DynamicServices.MONITORD);

  const gwHealth = gwCtx.health[gwInfo?.id].status
    ? gwCtx.health[gwInfo?.id].status === HEALTHY_STATUS
      ? GatewayTypeEnum.HEALTHY_GATEWAY
      : GatewayTypeEnum.UNHEALTHY_GATEWAY
    : 'N/A';

  if (isCpuPercentLoading) {
    return <LoadingFiller />;
  }

  const data: Array<DataRows> = [
    [
      {
        category: 'Health',
        value: gwHealth,
        statusCircle: true,
        // make kpi inactive if health status had error (health service not enabled)
        statusInactive: !gwCtx.health[gwInfo?.id]?.status,
        status: gwHealth === GatewayTypeEnum.HEALTHY_GATEWAY,
        tooltip:
          "Federation gateway's health as reported by the health service",
      },
      {
        category: 'Last Check in',
        value: checkInTime?.toLocaleString() ?? '-',
        statusCircle: false,
        tooltip: 'The last Time the gateway checked in',
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
        value: eventAggregationEnabled ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: eventAggregationEnabled,
      },
      {
        category: 'Log Aggregation',
        value: logAggregationEnabled ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: logAggregationEnabled,
      },
      {
        category: 'CPE Monitoring',
        value: cpeMonitoringEnabled ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: cpeMonitoringEnabled,
      },
    ],
  ];
  return <DataGrid data={data} />;
}
