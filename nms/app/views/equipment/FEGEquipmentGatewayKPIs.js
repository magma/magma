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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {DataRows} from '../../components/DataGrid';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataGrid from '../../components/DataGrid';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import FEGGatewayContext from '../../components/context/FEGGatewayContext';
import MagmaV1API from '../../../generated/WebClient';
import React from 'react';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';

// $FlowFixMe migrated to typescript
import {HEALTHY_STATUS, UNHEALTHY_STATUS} from '../../components/GatewayUtils';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {getLatency} from './EquipmentGatewayKPIs';
import {useContext} from 'react';
import {useParams} from 'react-router-dom';

/**
 * Displays the maximum latency, minimum latency, average latency,
 * total federation gateway count, healthy federation gateway count,
 * and the percentage of healthy federation gateways.
 */
export default function FEGEquipmentGatewayKPIs() {
  const params = useParams();
  const ctx = useContext(FEGGatewayContext);
  const fegGatewaysHealthStatus = ctx.health;
  const networkId: string = nullthrows(params.networkId);
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
  const avgLatencies: Array<number> =
    avgResponse?.data?.result
      ?.map(item => {
        return parseFloat(item?.value?.[1]);
      })
      .filter(Boolean) ?? [];
  let avgLatency = 0;
  if (avgLatencies && avgLatencies.length) {
    const sum = avgLatencies.reduce(function (a, b) {
      return a + b;
    }, 0);
    avgLatency = sum / avgLatencies.length;
    avgLatency = avgLatency.toFixed(2);
  }
  const fegGatewayCount = Object.keys(ctx.state).filter(Boolean).length;
  const upCount = Object.keys(fegGatewaysHealthStatus).filter(
    fegGatewayId =>
      fegGatewaysHealthStatus[fegGatewayId].status === HEALTHY_STATUS,
  ).length;
  const downCount = Object.keys(fegGatewaysHealthStatus).filter(
    fegGatewayId =>
      fegGatewaysHealthStatus[fegGatewayId].status === UNHEALTHY_STATUS,
  ).length;
  let percentHealthyGw = 0;
  if (upCount > 0 && upCount + downCount > 0) {
    percentHealthyGw = ((upCount * 100) / (upCount + downCount)).toFixed(2);
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
        category: 'Federation Gateway Count',
        value: fegGatewayCount,
        tooltip: 'Total number of federation gateways',
      },
      {
        category: 'Healthy Federation Gateway Count',
        value: upCount,
        tooltip: 'Total number of healthy federation gateways',
      },
      {
        category: '% Healthy Gateways',
        value: percentHealthyGw,
        tooltip: '% of gateways which are healthy',
      },
    ],
  ];

  return <DataGrid data={kpiData} />;
}
