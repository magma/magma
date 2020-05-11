/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from '@fbcnms/ui/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import Grid from '@material-ui/core/Grid';
import React, {useState} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import TimeRangeSelector from '@fbcnms/ui/insights/TimeRangeSelector';
import {makeStyles} from '@material-ui/styles';
import type {ChartStyle} from '@fbcnms/ui/insights/AsyncMetric';
import type {TimeRange} from '@fbcnms/ui/insights/AsyncMetric';

const useStyles = makeStyles(_ => ({
  formControl: {
    minWidth: '50px',
  },
  appBar: {
    display: 'inline-block',
  },
}));

export default function() {
  const state = {
    title: 'Frequency of Alerts and Events',
    legendLabels: ['Alerts', 'Events'],
  };
  const classes = useStyles();
  const [timeRange, setTimeRange] = useState<TimeRange>('3_hours');

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
      <AppBar className={classes.appBar} position="static" color="default">
        <TimeRangeSelector
          className={classes.formControl}
          value={timeRange}
          onChange={setTimeRange}
        />
      </AppBar>
      <Grid>
        <Card>
          <CardHeader
            title={<Text variant="h6">{state.title}</Text>}
            subheader={
              <AsyncMetric
                style={chartStyle}
                label={state.title}
                unit=""
                queries={['sum(ALERTS)']}
                timeRange={timeRange}
                legendLabels={state.legendLabels}
              />
            }
          />
        </Card>
      </Grid>
    </>
  );
}
