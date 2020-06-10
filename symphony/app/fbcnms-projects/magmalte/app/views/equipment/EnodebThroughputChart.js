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

import AsyncMetric from '@fbcnms/ui/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import DataUsageIcon from '@material-ui/icons/DataUsage';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import moment from 'moment';

import {DateTimePicker} from '@material-ui/pickers';
import {useState} from 'react';

export type EnbThroughputChartProps = {
  title: string,
  queries: Array<string>,
  legendLabels: Array<string>,
};

export default function EnodebThroughputChart(props: EnbThroughputChartProps) {
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [endDate, setEndDate] = useState(moment());

  return (
    <>
      <Grid container align="top" alignItems="flex-start">
        <Grid item xs={6}>
          <Text>
            <DataUsageIcon />
            {props.title}
          </Text>
        </Grid>
        <Grid item xs={6}>
          <Grid container justify="flex-end" alignItems="center" spacing={1}>
            <Grid item>
              <Text>Filter By Date</Text>
            </Grid>
            <Grid item>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                maxDate={endDate}
                disableFuture
                value={startDate}
                onChange={setStartDate}
              />
            </Grid>
            <Grid item>
              <Text>To</Text>
            </Grid>
            <Grid item>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                disableFuture
                value={endDate}
                onChange={setEndDate}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <Card>
        <CardHeader
          title={<Text variant="h6">{props.title}</Text>}
          subheader={
            <AsyncMetric
              style={{
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
              }}
              label={props.title}
              unit=""
              queries={props.queries}
              timeRange={'3_hours'}
              startEnd={[startDate, endDate]}
              legendLabels={props.legendLabels}
            />
          }
        />
      </Card>
    </>
  );
}
