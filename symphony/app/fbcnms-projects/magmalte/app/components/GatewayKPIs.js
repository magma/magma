/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {KPIData} from './KPITray';
import type {lte_gateway} from '@fbcnms/magma-api';

import CellWifiIcon from '@material-ui/icons/CellWifi';
import isGatewayHealthy from './GatewayUtils';
import KPITray from './KPITray';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import nullthrows from '@fbcnms/util/nullthrows';
import React from 'react';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {useRouter} from '@fbcnms/ui/hooks';

export default function GatewayKPIs() {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const {response: lteGateways, isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdGateways,
    {
      networkId: networkId,
    },
  );

  if (isLoading || !lteGateways) {
    return <LoadingFiller />;
  }
  const [upCount, downCount] = gatewayStatus(lteGateways);
  const kpiData: KPIData[] = [
    {category: 'Severe Events', value: 0},
    {category: 'Connected', value: upCount || 0},
    {category: 'Disconnected', value: downCount || 0},
  ];
  return <KPITray icon={CellWifiIcon} description="Gateways" data={kpiData} />;
}

function gatewayStatus(gatewaySt: {[string]: lte_gateway}): [number, number] {
  let upCount = 0;
  let downCount = 0;
  Object.keys(gatewaySt)
    .map((k: string) => gatewaySt[k])
    .filter((g: lte_gateway) => g.cellular && g.id)
    .map(function(gateway: lte_gateway) {
      isGatewayHealthy(gateway) ? upCount++ : downCount++;
    });
  return [upCount, downCount];
}
