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

export type KPIData = {
  category: string,
  value: number | string,
  unit?: string,
  icon?: ComponentType<SvgIconExports>,
};
type Props = {
  icon?: ComponentType<SvgIconExports>,
  description?: string,
  data: KPIData[],
};

function KPIIcon(Icon: ComponentType<SvgIconExports>) {
  return <Icon fontSize="large" />;
}

export default function KPITray(props: Props) {
  return (
    <Grid container alignItems="center" justify="center">
      {props.icon ? (
        <>
          <Grid item>
            <Card elevation={0}>
              <CardContent>{KPIIcon(props.icon)}</CardContent>
            </Card>
          </Grid>
          <Grid item>
            <Card elevation={0}>
              <CardContent>
                <Text variant="h6">{props.description}</Text>
              </CardContent>
            </Card>
          </Grid>
        </>
      ) : (
        ''
      )}
      {props.data.map((kpi, i) => (
        <Grid item xs key={'data-' + i}>
          <Card>
            <CardHeader
              title={kpi.category}
              subheader={
                <Grid
                  container
                  justify="center"
                  alignItems="center"
                  spacing={3}>
                  <Grid item>{kpi.icon ? KPIIcon(kpi.icon) : ''}</Grid>
                  <Grid item>
                    <Text variant="h5">
                      {kpi.value} {kpi.unit ?? ''}
                    </Text>
                  </Grid>
                </Grid>
              }
              titleTypographyProps={{align: 'center', variant: 'body1'}}
              data-testid={kpi.category}
            />
          </Card>
        </Grid>
      ))}
    </Grid>
  );
}
