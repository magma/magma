/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {KPIData} from '../../components/KPITray';
import type {lte_gateway} from '@fbcnms/magma-api';

import isGatewayHealthy from '../../components/GatewayUtils';
import KPITray from '../../components/KPITray';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import nullthrows from '@fbcnms/util/nullthrows';
import React from 'react';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {useRouter} from '@fbcnms/ui/hooks';

const getLatency = (resp, fn) => {
  const respArr = resp?.data?.result
    ?.map(item => {
      return parseFloat(item?.value?.[1]);
    })
    .filter(Boolean);
  return respArr && respArr.length ? fn(...respArr).toFixed(2) : 0;
};

export default function EquipmentGatewayKPIs({
  lteGateways,
}: {
  lteGateways: {[string]: lte_gateway},
}) {
  const {match} = useRouter();
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
    const sum = avgLatencyArr.reduce(function(a, b) {
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

  const kpiData: KPIData[] = [
    {category: 'Max Latency', value: maxLatency, unit: 'ms'},
    {category: 'Min Latency', value: minLatency, unit: 'ms'},
    {category: 'Avg Latency', value: avgLatency, unit: 'ms'},
    {
      category: '% Healthy Gateways',
      value: pctHealthyGw,
    },
  ];
  return <KPITray data={kpiData} />;
}
