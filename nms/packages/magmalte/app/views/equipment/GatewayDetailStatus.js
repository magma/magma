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
import GatewayContext from '../../components/context/GatewayContext';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import isGatewayHealthy, {DynamicServices} from '../../components/GatewayUtils';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {
  REFRESH_INTERVAL,
  useRefreshingContext,
} from '../../components/context/RefreshContext';
import {useRouter} from '@fbcnms/ui/hooks';

export default function GatewayDetailStatus({refresh}: {refresh: boolean}) {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const gatewayId: string = nullthrows(match.params.gatewayId);
  // Auto refresh gateways every 30 seconds
  const state = useRefreshingContext({
    context: GatewayContext,
    networkId: networkId,
    type: 'gateway',
    interval: REFRESH_INTERVAL,
    id: gatewayId,
    refresh: refresh,
  });
  const gwInfo = state[gatewayId];
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
    MagmaV1API.getNetworksByNetworkIdPrometheusQueryRange,
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

  const isHealthy = isGatewayHealthy(gwInfo);
  const data: DataRows[] = [
    [
      {
        category: 'Health',
        value: isHealthy ? 'Good' : 'Bad',
        statusCircle: true,
        status: isGatewayHealthy(gwInfo),
        tooltip: isHealthy
          ? 'Gateway checked in recently'
          : "Gateway hasn't checked in within last 5 minutes",
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
