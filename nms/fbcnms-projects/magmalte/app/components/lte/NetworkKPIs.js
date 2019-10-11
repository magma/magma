/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MetricGraphConfig} from '../insights/Metrics';
import type {TimeRange} from '../insights/AsyncMetric';

import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from '../insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import React from 'react';
import TimeRangeSelector from '../insights/TimeRangeSelector';
import Typography from '@material-ui/core/Typography';

import {makeStyles} from '@material-ui/styles';
import {resolveQuery} from '../insights/Metrics';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  formControl: {
    minWidth: '200px',
    padding: theme.spacing(),
  },
  appBar: {
    display: 'inline-block',
  },
}));

export default function CloudMetrics() {
  const classes = useStyles();
  const [timeRange, setTimeRange] = useState<TimeRange>('3_hours');

  const chartConfigs: MetricGraphConfig[] = [
    {
      label: 'Disk Percent',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(disk_percent)',
        },
      ],
      legendLabels: ['Disk Percent'],
      unit: '%',
    },
    {
      label: 'Number of Connected eNBs',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(enb_connected)',
        },
      ],
      legendLabels: ['Connected'],
      unit: '',
    },
    {
      label: 'Number of Connected UEs',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(ue_connected)',
        },
      ],
      legendLabels: ['Connected'],
      unit: '',
    },
    {
      label: 'Number of Registered UEs',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(ue_registered)',
        },
      ],
      legendLabels: ['Registered'],
      unit: '',
    },
    {
      label: 'S1 Setup',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(s1_setup)',
        },
        {
          resolveQuery: _ => "sum(s1_setup{result='success'})",
        },
        {
          resolveQuery: _ => "sum(s1_setup) - sum(s1_setup{result='success'})",
        },
      ],
      legendLabels: ['Total', 'Success', 'Failure'],
      unit: '',
    },
    {
      label: 'Attach/Reg Attempts',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(ue_attach)',
        },
        {
          resolveQuery: _ => "sum(ue_attach{result='attach_proc_successful'})",
        },
        {
          resolveQuery: _ =>
            "sum(ue_attach) - sum(ue_attach{result='attach_proc_successful'})",
        },
      ],
      legendLabels: ['Total', 'Success', 'Failure'],
      unit: '',
    },
    {
      label: 'Detach/Dereg Attempts',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(ue_detach)',
        },
        {
          resolveQuery: _ => "sum(ue_detach{result='attach_proc_successful'})",
        },
        {
          resolveQuery: _ =>
            "sum(ue_detach) - sum(ue_detach{result='attach_proc_successful'})",
        },
      ],
      legendLabels: ['Total', 'Success', 'Failure'],
      unit: '',
    },
    {
      label: 'GPS Connection Uptime',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'avg(enodeb_gps_connected)',
        },
      ],
      legendLabels: ['Uptime'],
      unit: '',
    },
    {
      label: 'Device Transmitting Status',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'avg(enodeb_rf_tx_enabled)',
        },
      ],
      legendLabels: ['Transmitting Status'],
      unit: '',
    },
    {
      label: 'Service Requests',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(service_request)',
        },
        {
          resolveQuery: _ => "sum(service_request{result='success'})",
        },
        {
          resolveQuery: _ =>
            "sum(service_request) - sum(service_request{result='success'})",
        },
      ],
      legendLabels: ['Total', 'Success', 'Failure'],
      unit: '',
    },
  ];

  return (
    <>
      <AppBar className={classes.appBar} position="static" color="default">
        <TimeRangeSelector
          className={classes.formControl}
          value={timeRange}
          onChange={setTimeRange}
        />
      </AppBar>
      <GridList cols={2} cellHeight={300}>
        {chartConfigs.map((config, i) => (
          <GridListTile key={i} cols={1}>
            <Card>
              <CardContent>
                <Typography component="h6" variant="h6">
                  {config.label}
                </Typography>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    queries={resolveQuery(config, '')}
                    timeRange={timeRange}
                    legendLabels={config.legendLabels}
                  />
                </div>
              </CardContent>
            </Card>
          </GridListTile>
        ))}
      </GridList>
    </>
  );
}
