/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import Grid from '@material-ui/core/Grid';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React, {useEffect, useState} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import type {lte_gateway} from '@fbcnms/magma-api';

const GATEWAY_KEEPALIVE_TIMEOUT_MS = 1000 * 5 * 60;

export default function GatewayKPIs() {
  const [gatewaySt, setGatewaySt] = useState({});
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
  return (
    <Grid container alignItems="center">
      <Grid item>
        <Card elevation={0}>
          <CardHeader title="Gateways" />
          <CardContent>
            <CellWifiIcon fontSize="large" />
          </CardContent>
        </Card>
      </Grid>
      <Grid item>
        <Card variant="outlined">
          <CardHeader title="Severe Events" />
          <CardContent>
            <Text variant="h6" data-testid="Severe Events">
              0
            </Text>
          </CardContent>
        </Card>
      </Grid>
      <Grid item>
        <Card variant="outlined">
          <CardHeader title="Connected" />
          <CardContent>
            <Text variant="h6" data-testid="Connected">
              {upCount ?? 0}
            </Text>
          </CardContent>
        </Card>
      </Grid>
      <Grid item>
        <Card variant="outlined">
          <CardHeader title="Disconnected" />
          <CardContent>
            <Text variant="h6" data-testid="Disconnected">
              {downCount ?? 0}
            </Text>
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  );
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
