/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import AsyncMetric from '@fbcnms/ui/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import moment from 'moment';
import type {ChartStyle} from '@fbcnms/ui/insights/AsyncMetric';

type Props = {
  startEnd: [moment, moment],
};

const isValid = (start, end): boolean => {
  return start.isValid() && end.isValid() && moment.min(start, end) === start;
};

export default function({startEnd}: Props) {
  const [start, end] = startEnd;
  const state = {
    title: 'Frequency of Alerts and Events',
    legendLabels: ['Alerts', 'Events'],
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
              timeRange={'3_days'}
              startEnd={isValid(start, end) ? startEnd : undefined}
              legendLabels={state.legendLabels}
            />
          }
        />
      </Card>
    </Grid>
  );
}
