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
import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import type {ComponentType} from 'react';

export type KPIData = {category: string, value: number | string, unit?: string};
type Props = {
  icon?: ComponentType<SvgIconExports>,
  description?: string,
  data: KPIData[],
};

export default function KPITray(props: Props) {
  const kpiTray = [];
  if (props.icon) {
    const KpiIcon = props.icon;
    kpiTray.push(
      <Grid item key="kpiIcon">
        <Card elevation={0}>
          <CardContent>
            <KpiIcon fontSize="large" />
          </CardContent>
        </Card>
      </Grid>,
      <Grid item key="kpiTitle">
        <Card elevation={0}>
          <CardContent>
            <Text variant="h6">{props.description}</Text>
          </CardContent>
        </Card>
      </Grid>,
    );
  }

  kpiTray.push(
    props.data.map((kpi, i) => (
      <Grid item xs key={'data-' + i}>
        <Card>
          <CardHeader
            title={kpi.category}
            subheader={kpi.value + (kpi.unit ?? '')}
            titleTypographyProps={{align: 'center', variant: 'body1'}}
            subheaderTypographyProps={{variant: 'h5', align: 'center'}}
            data-testid={kpi.category}
          />
        </Card>
      </Grid>
    )),
  );
  return (
    <Grid container alignItems="center" justify="center">
      {kpiTray}
    </Grid>
  );
}
