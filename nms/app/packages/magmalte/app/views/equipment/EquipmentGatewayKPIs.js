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
import type {DataRows} from '../../components/DataGrid';
import type {lte_gateway, promql_return_object} from '@fbcnms/magma-api';

import DataGrid from '../../components/DataGrid';
import GatewayContext from '../../components/context/GatewayContext';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import isGatewayHealthy from '../../components/GatewayUtils';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {useContext} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

export function getLatency(
  resp: ?promql_return_object,
  fn: (...args: Array<number>) => number,
) {
  const respArr = resp?.data?.result
    ?.map(item => {
      return parseFloat(item?.value?.[1]);
    })
    .filter(Boolean);
  return respArr && respArr.length ? fn(...respArr).toFixed(2) : 0;
}

export default function EquipmentGatewayKPIs() {
  const {match} = useRouter();
  const ctx = useContext(GatewayContext);
  const lteGateways = ctx.state;

  const networkId: string = nullthrows(match.params.networkId);
  const timeRange = '3h';
  const {response: maxResponse} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusQuery,
    {
      networkId: networkId,
      query: `max_over_time(magmad_ping_rtt_ms{service="magmad",metric="rtt_ms"}[${timeRange}])`,
    },
  );

  const {response: minResponse} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusQuery,
    {
      networkId: networkId,
      query: `min_over_time(magmad_ping_rtt_ms{service="magmad",metric="rtt_ms"}[${timeRange}])`,
    },
  );

  const {response: avgResponse} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusQuery,
    {
      networkId: networkId,
      query: `avg_over_time(magmad_ping_rtt_ms{service="magmad",metric="rtt_ms"}[${timeRange}])`,
    },
  );

  const maxLatency = getLatency(maxResponse, Math.max);
  const minLatency = getLatency(minResponse, Math.min);

  const avgLatencyArr = avgResponse?.data?.result
    ?.map(item => {
      return parseFloat(item?.value?.[1]);
    })
    .filter(Boolean);

  let avgLatency = 0;
  if (avgLatencyArr && avgLatencyArr.length) {
    const sum = avgLatencyArr.reduce(function (a, b) {
      return a + b;
    }, 0);
    avgLatency = sum / avgLatencyArr.length;
    avgLatency = avgLatency.toFixed(2);
  }

  let upCount = 0;
  let downCount = 0;
  Object.keys(lteGateways)
    .map((gwId: string) => lteGateways[gwId])
    .filter((g: lte_gateway) => g.cellular && g.id)
    .map((gateway: lte_gateway) => {
      isGatewayHealthy(gateway) ? upCount++ : downCount++;
    });
  let pctHealthyGw = 0;
  if (upCount > 0 && upCount + downCount > 0) {
    pctHealthyGw = ((upCount * 100) / (upCount + downCount)).toFixed(2);
  }

  const kpiData: DataRows[] = [
    [
      {
        category: 'Max Latency',
        value: maxLatency,
        unit: 'ms',
        tooltip:
          'Max ping latency(for host 8.8.8.8) observed across all gateways',
      },
      {
        category: 'Min Latency',
        value: minLatency,
        unit: 'ms',
        tooltip:
          'Min ping latency(for host 8.8.8.8) observed across all gateways',
      },
      {
        category: 'Avg Latency',
        value: avgLatency,
        unit: 'ms',
        tooltip:
          'Avg ping latency(for host 8.8.8.8) observed across all gateways',
      },
      {
        category: '% Healthy Gateways',
        value: pctHealthyGw,
        tooltip: '% of gateways which have checked in within last 5 minutes',
      },
    ],
  ];
  return <DataGrid data={kpiData} />;
}
