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
import Text from '@fbcnms/ui/components/design-system/Text';
import TimeRangeSelector from '@fbcnms/ui/insights/TimeRangeSelector';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 150,
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
      <Grid container align="top" alignItems="flex-start">
        <Grid item xs={6}>
          <Text>
            <DataUsageIcon />
            Gateway Check-Ins
          </Text>
        </Grid>
        <Grid item xs={6}>
          <Grid container justify="flex-end" alignItems="center" spacing={1}>
            <Grid item>
              <Text>Filter By Time</Text>
            </Grid>
            <Grid item>
              <TimeRangeSelector
                className={classes.formControl}
                value={timeRange}
                onChange={setTimeRange}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <Card>
        <CardHeader
          title={<Text variant="h6">{state.title}</Text>}
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
