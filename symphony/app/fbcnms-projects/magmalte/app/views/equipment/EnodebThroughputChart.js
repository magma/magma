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
import moment from 'moment';
import React from 'react';
import Text from '../../theme/design-system/Text';

import {colors} from '../../theme/default';
import {DateTimePicker} from '@material-ui/pickers';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  cardTitleRow: {
    marginBottom: theme.spacing(1),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
  dateTimeText: {
    color: colors.primary.comet,
  },
  input: {
    color: colors.primary.brightGray,
    backgroundColor: colors.button.fill,
    border: `1px solid ${colors.primary.brightGray}`,
    borderRadius: '4px',
    textAlign: 'center',
    padding: `${theme.spacing(0.5)}px 0`,
  },
}));

export type EnbThroughputChartProps = {
  title: string,
  queries: Array<string>,
  legendLabels: Array<string>,
};

export default function EnodebThroughputChart(props: EnbThroughputChartProps) {
  const classes = useStyles();
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [endDate, setEndDate] = useState(moment());

  return (
    <>
      <Grid container alignItems="center" className={classes.cardTitleRow}>
        <Grid container xs={6}>
          <DataUsageIcon className={classes.cardTitleIcon} />
          <Text variant="body1">{props.title}</Text>
        </Grid>
        <Grid item xs={6} className={classes.dateFilters}>
          <Grid container justify="flex-end" alignItems="center" spacing={1}>
            <Grid item>
              <Text variant="body3" className={classes.dateTimeText}>
                Filter By Date
              </Text>
            </Grid>
            <Grid item>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                maxDate={endDate}
                disableFuture
                value={startDate}
                inputProps={{className: classes.input}}
                onChange={setStartDate}
              />
            </Grid>
            <Grid item>
              <Text variant="body3" className={classes.dateTimeText}>
                to
              </Text>
            </Grid>
            <Grid item>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                disableFuture
                value={endDate}
                inputProps={{className: classes.input}}
                onChange={setEndDate}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <Card>
        <CardHeader
          title={<Text variant="body2">Frequency of {props.title}</Text>}
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
              label={`Frequency of ${props.title}`}
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
