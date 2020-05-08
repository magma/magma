/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import CellWifiIcon from '@material-ui/icons/CellWifi';
import KPITray from './KPITray';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React, {useEffect, useState} from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import type {KPIData} from './KPITray';
import type {lte_gateway} from '@fbcnms/magma-api';

const GATEWAY_KEEPALIVE_TIMEOUT_MS = 1000 * 5 * 60;

export default function GatewayKPIs() {
  const [gatewaySt, setGatewaySt] = useState<{[string]: lte_gateway}>({});
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchSt = async () => {
      try {
        const resp = await MagmaV1API.getLteByNetworkIdGateways({
          networkId: networkId,
        });
        setGatewaySt(resp);
      } catch (error) {
        enqueueSnackbar('Error getting gateway information', {
          variant: 'error',
        });
      }
    };
    fetchSt();
  }, [networkId, enqueueSnackbar]);

  const [upCount, downCount] = gatewayStatus(gatewaySt);
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
      const {status} = gateway;
      let isBackhaulDown = true;
      if (status != null) {
        const checkin = status.checkin_time;
        if (checkin != null) {
          const duration = Date.now() - checkin;
          isBackhaulDown = duration > GATEWAY_KEEPALIVE_TIMEOUT_MS;
        }
      }
      isBackhaulDown ? downCount++ : upCount++;
    });
  return [upCount, downCount];
}
