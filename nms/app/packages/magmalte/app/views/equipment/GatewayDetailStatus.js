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
import type {lte_gateway} from '@fbcnms/magma-api';

import DataGrid from '../../components/DataGrid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import isGatewayHealthy from '../../components/GatewayUtils';
import {useRouter} from '@fbcnms/ui/hooks';

export default function GatewayDetailStatus({gwInfo}: {gwInfo: lte_gateway}) {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  let checkInTime = new Date(0);
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
    gwInfo.magmad.dynamic_services.includes('td-agent-bit');

  const eventAggregation =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes('eventd');

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
        value: checkInTime.toLocaleString(),
        statusCircle: false,
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
        category: 'CPU Usage',
        value: cpuPercent?.data?.result?.[0]?.values?.[0]?.[1] ?? 'Unknown',
        unit:
          cpuPercent?.data?.result?.[0]?.values?.[0]?.[1] ?? false ? '%' : '',
        statusCircle: false,
        tooltip: 'Current Gateway CPU %',
      },
    ],
  ];
  return <DataGrid data={data} />;
}
