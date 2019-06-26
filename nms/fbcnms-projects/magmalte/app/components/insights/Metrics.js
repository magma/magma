/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TimeRange} from './AsyncMetric';

import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from './AsyncMetric';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import FormControl from '@material-ui/core/FormControl';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFiller from '../LoadingFiller';
import MenuItem from '@material-ui/core/MenuItem';
import MagmaTopBar from '../MagmaTopBar';
import Typography from '@material-ui/core/Typography';
import {Route, Switch} from 'react-router-dom';
import Select from '@material-ui/core/Select';

import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {find} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useSnackbar, useRouter} from '@fbcnms/ui/hooks';
import {useCallback, useState} from 'react';

const useStyles = makeStyles(theme => ({
  appBar: {
    display: 'inline-block',
  },
  chartRow: {
    display: 'flex',
  },
  formControl: {
    minWidth: '200px',
    padding: theme.spacing.unit,
  },
}));

type Config = {
  id: string,
  filters: string[],
  label: string,
  metric: string,
  unit?: string,
};

const CHART_CONFIGS: Config[] = [
  {
    id: 'enodeb_rf_tx_enabled',
    filters: ['service=enodebd'],
    label: 'E-Node B Status',
    metric: 'enodeb_rf_tx_enabled',
  },
  {
    id: 'connected_subscribers',
    filters: ['service="mme"'],
    label: 'Connected Subscribers',
    metric: 'ue_connected',
  },
  {
    id: 'download_throughput',
    filters: ['service=enodebd'],
    label: 'Download Throughput',
    metric: 'pdcp_user_plane_bytes_dl',
    // 'transform' => 'formula(* $1 26.667)',
    unit: ' Kbps',
  },
  {
    id: 'upload_throughput',
    filters: ['service=enodebd'],
    label: 'Upload Throughput',
    metric: 'pdcp_user_plane_bytes_ul',
    // 'transform' => 'formula(* $1 26.667)',
    unit: ' Kbps',
  },
  {
    id: 'latency',
    filters: ['service=magmad'],
    label: 'Latency',
    metric: 'magmad_ping_rtt_ms_8_8_8_8_metric_rtt_ms',
    unit: ' ms',
  },
  {
    id: 'gateway_cpu',
    filters: ['service=magmad'],
    label: 'Gateway CPU (%)',
    metric: 'cpu_percent',
    unit: '%',
  },
  {
    id: 'temperature_coretemp_0',
    filters: ['service=magmad'],
    label: 'Temperature (℃)',
    metric: 'temperature_.+_coretemp_0',
    unit: '℃',
  },
  {
    id: 'disk',
    filters: ['service=magmad'],
    label: 'Disk (%)',
    metric: 'disk_percent',
    unit: '%',
  },
  {
    id: 's6a_auth_success',
    filters: ['service=subscriberdb'],
    label: 's6a Auth Success',
    metric: 's6a_auth_success',
    // 'transform' => 'rate(1m, duration=900)',
    unit: '',
  },
  {
    id: 's6a_auth_failure_code_ResultCode_DIAMETER_AUTHORIZATION_REJECTED',
    filters: ['service=subscriberdb'],
    label: 's6a Auth Failure',
    metric: 's6a_auth_failure',
    // 'transform' => 'rate(1m, duration=900)',
    unit: '',
  },
];

function resolveQuery(config: Config, gatewayId: string) {
  const filters = [...config.filters, `gatewayID=${gatewayId}`].join(',');
  return `${config.metric},${filters}`;
}

function Metrics() {
  const {history, match} = useRouter();
  const classes = useStyles();
  const selectedGateway = match.params.selectedGatewayId;
  const [timeRange, setTimeRange] = useState<TimeRange>('24_hours');

  const {error, isLoading, response: response} = useAxios({
    method: 'get',
    url: MagmaAPIUrls.gateways(match, true),
  });

  const onGatewayChanged = useCallback(
    event => {
      const gatewayId = event.target.value;
      history.push(`/${match.params.networkId}/metrics/${gatewayId}`);
    },
    [match],
  );

  useSnackbar('Error fetching devices', {variant: 'error'}, error);

  if (error) {
    return <LoadingFiller />;
  }
  if (error || isLoading || !response || !response.data) {
    return <LoadingFiller />;
  }

  const gateways = response.data.filter(state => state.record);
  const defaultGateway = find(
    gateways,
    gateway => gateway.status?.hardware_id !== null,
  );
  const selectedGatewayOrDefault =
    selectedGateway || defaultGateway?.gateway_id;

  return (
    <>
      <AppBar className={classes.appBar} position="static" color="default">
        <FormControl variant="filled" className={classes.formControl}>
          <InputLabel htmlFor="devices">Device</InputLabel>
          <Select
            inputProps={{id: 'devices'}}
            value={selectedGatewayOrDefault}
            onChange={onGatewayChanged}>
            {gateways.map(device => (
              <MenuItem value={device.gateway_id} key={device.gateway_id}>
                {device.record.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <FormControl variant="filled" className={classes.formControl}>
          <InputLabel htmlFor="time_range">Period</InputLabel>
          <Select
            inputProps={{id: 'time_range'}}
            value={timeRange}
            onChange={event => setTimeRange((event.target.value: any))}>
            <MenuItem value="24_hours">Last 24 hours</MenuItem>
            <MenuItem value="7_days">Last 7 days</MenuItem>
            <MenuItem value="14_days">Last 14 days</MenuItem>
            <MenuItem value="30_days">Last 30 days</MenuItem>
          </Select>
        </FormControl>
      </AppBar>
      <GridList cols={2} cellHeight={300}>
        {CHART_CONFIGS.map(config => (
          <GridListTile key={config.id} cols={1}>
            <Card>
              <CardContent>
                <Typography component="h6" variant="h6">
                  {config.label}
                </Typography>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    query={resolveQuery(config, selectedGatewayOrDefault)}
                    timeRange={timeRange}
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

export default function() {
  const {match} = useRouter();
  return (
    <>
      <MagmaTopBar />
      <Switch>
        <Route path={`${match.path}/:selectedGatewayId`} component={Metrics} />
        <Route path={`${match.path}/`} component={Metrics} />
      </Switch>
    </>
  );
}
