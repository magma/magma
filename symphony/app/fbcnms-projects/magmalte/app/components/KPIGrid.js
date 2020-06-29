/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import Grid from '@material-ui/core/Grid';
import React from 'react';

type KPIData = {
  category: string,
  value: number | string,
  unit?: string,
  statusCircle: boolean,
  status?: 'Disabled' | 'Up' | 'Down',
};

export type KPIRows = KPIData[];

type Props = {data: KPIRows[]};

export default function KPIGrid(props: Props) {
  const kpiGrid = props.data.map((row, i) => (
    <Grid key={i} container direction="row">
      {row.map((kpi, j) => (
        <Grid item xs key={`data-${i}-${j}`}>
          <Card>
            <CardHeader
              title={kpi.category}
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body1'}}
              subheader={
                <Grid container>
                  {kpi.statusCircle && (
                    <Grid item>
                      <DeviceStatusCircle
                        isGrey={kpi.status === 'Disabled'}
                        isActive={kpi.status === 'Up'}
                        isFilled={true}
                      />
                    </Grid>
                  )}
                  <Grid item>{kpi.value + (kpi.unit ?? '')}</Grid>
                </Grid>
              }
            />
          </Card>
        </Grid>
      ))}
    </Grid>
  ));
  return (
    <Grid container alignItems="center" justify="center">
      {kpiGrid}
    </Grid>
  );
}
