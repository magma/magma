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
import DeviceStatusCircle from '../theme/design-system/DeviceStatusCircle';
import Grid from '@material-ui/core/Grid';
import React from 'react';

import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  kpiHeaderBlock: {
    display: 'flex',
    alignItems: 'center',
    padding: 0,
  },
  kpiHeaderContent: {
    display: 'flex',
    alignItems: 'center',
  },
  kpiHeaderIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
  kpiBlock: {
    boxShadow: `0 0 0 1px ${colors.primary.concrete}`,
  },
  kpiLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  kpiValue: {
    color: colors.primary.brightGray,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: props => (props.hasStatus ? 'calc(100% - 16px)' : '100%'),
  },
  kpiBox: {
    width: '100%',
    '& > div': {
      width: '100%',
    },
  },
}));

// Status Indicator displays a small text with an DeviceStatusCircle icon
// disabled indicates if the status color is to be grayed out
// up/down indicates if we have to display status to be in green or in red
function StatusIndicator(disabled: boolean, up: boolean, val: string) {
  const props = {hasStatus: true};
  const classes = useStyles(props);
  return (
    <Grid container zeroMinWidth alignItems="center" xs={12}>
      <Grid item>
        <DeviceStatusCircle isGrey={disabled} isActive={up} isFilled={true} />
      </Grid>
      <Grid item className={classes.kpiValue}>
        {val}
      </Grid>
    </Grid>
  );
}

type KPIData = {
  category: string,
  value: number | string,
  unit?: string,
  statusCircle: boolean,
  statusInactive: boolean,
  status?: 'Disabled' | 'Up' | 'Down',
};

export type KPIRows = KPIData[];

type Props = {data: KPIRows[]};

export default function KPIGrid(props: Props) {
  const classes = useStyles();
  const kpiGrid = props.data.map((row, i) => (
    <Grid key={i} container direction="row" zeroMinWidth>
      {row.map((kpi, j) => (
        <Grid
          item
          xs={12}
          md
          key={`data-${i}-${j}`}
          zeroMinWidth
          className={classes.kpiBlock}>
          <CardHeader
            className={classes.kpiBox}
            title={kpi.category}
            titleTypographyProps={{
              variant: 'body3',
              className: classes.kpiLabel,
              title: kpi.category,
            }}
            subheaderTypographyProps={{
              variant: 'body1',
              className: classes.kpiValue,
              title: kpi.value + (kpi.unit ?? ''),
            }}
            subheader={
              kpi.statusCircle
                ? StatusIndicator(
                    kpi.statusInactive || false,
                    kpi.status,
                    kpi.value + (kpi.unit ?? ''),
                  )
                : kpi.value + (kpi.unit ?? '')
            }
          />
        </Grid>
      ))}
    </Grid>
  ));
  return (
    <Card elevation={0}>
      <Grid container alignItems="center" justify="center">
        {kpiGrid}
      </Grid>
    </Card>
  );
}
