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
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';

function fetchGatewayKPIs() {
  const mockKPIs = [
    {category: 'Severe Events', value: 10},
    {category: 'Connected', value: 20},
    {category: 'Disconnected', value: 30},
  ];
  return mockKPIs;
}

export default function GatewayKPIs() {
  const mockKPIs = fetchGatewayKPIs();
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
      {mockKPIs.map((kpi, i) => (
        <Grid item key={i}>
          <Card variant="outlined">
            <CardHeader title={kpi.category} />
            <CardContent>
              <Text variant="h6">{kpi.value}</Text>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
}
