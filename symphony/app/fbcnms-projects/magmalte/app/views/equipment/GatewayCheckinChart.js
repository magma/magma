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
import TimeRangeSelector from '../../theme/design-system/TimeRangeSelector';

import {CardTitleFilterRow} from '../../components/layout/CardTitleRow';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  dateTimeText: {
    color: colors.primary.comet,
  },
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

  function Filter() {
    return (
      <Grid container justify="flex-end" alignItems="center" spacing={1}>
        <Grid item>
          <Text variant="body3" className={classes.dateTimeText}>
            Filter By Time
          </Text>
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
    );
  }

  return (
    <>
      <CardTitleFilterRow
        icon={DataUsageIcon}
        label="Gateway Check-Ins"
        filter={Filter}
      />
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
