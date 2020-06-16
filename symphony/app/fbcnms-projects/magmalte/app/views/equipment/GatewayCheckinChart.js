/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {ChartStyle} from '@fbcnms/ui/insights/AsyncMetric';
import type {TimeRange} from '@fbcnms/ui/insights/AsyncMetric';

import AsyncMetric from '@fbcnms/ui/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import DataUsageIcon from '@material-ui/icons/DataUsage';
import Grid from '@material-ui/core/Grid';
import React, {useState} from 'react';
import Text from '../../theme/design-system/Text';
import {colors} from '../../theme/default';
import TimeRangeSelector from '../../theme/design-system/TimeRangeSelector';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  cardTitleRow: {
    marginBottom: theme.spacing(1),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 150,
  },
  cardFilters: {
    marginBottom: theme.spacing(1),
  },
}));

export default function() {
  const classes = useStyles();
  const [timeRange, setTimeRange] = useState<TimeRange>('3_hours');
  const state = {
    title: 'Frequency of Gateway Check-Ins',
    legendLabels: ['Check-Ins', 'Events'],
  };
  const chartStyle: ChartStyle = {
    data: {
      lineTension: 0.2,
      pointRadius: 0.1,
    },
    options: {
      xAxes: {
        gridLines: {
          display: false,
        },
        ticks: {
          maxTicksLimit: 10,
        },
      },
      yAxes: {
        gridLines: {
          drawBorder: true,
        },
        ticks: {
          maxTicksLimit: 1,
        },
      },
    },
    legend: {
      position: 'top',
      align: 'end',
    },
  };
  return (
    <>
      <Grid container alignItems="center">
        <Grid item xs={6}>
          <Grid container alignItems="center" className={classes.cardTitleRow}>
            <DataUsageIcon className={classes.cardTitleIcon} />
            <Text variant="body1">Gateway Check-Ins</Text>
          </Grid>
        </Grid>
        <Grid item xs={6}>
          <Grid
            container
            justify="flex-end"
            alignItems="center"
            className={classes.cardFilters}>
            <Grid item>
              <Text variant="body3">Filter By Time</Text>
            </Grid>
            <Grid item>
              <TimeRangeSelector
                variant="outlined"
                className={classes.formControl}
                value={timeRange}
                onChange={setTimeRange}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <Card elevation={0}>
        <CardHeader
          title={<Text variant="body2">{state.title}</Text>}
          subheader={
            <AsyncMetric
              style={chartStyle}
              label={state.title}
              unit=""
              queries={['sum(checkin_status)']}
              timeRange={timeRange}
              startEnd={undefined}
              legendLabels={state.legendLabels}
            />
          }
        />
      </Card>
    </>
  );
}
