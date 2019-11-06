/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from '@fbcnms/magmalte/app/components/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import FormControl from '@material-ui/core/FormControl';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import TimeRangeSelector from '@fbcnms/magmalte/app/components/insights/TimeRangeSelector';
import WifiSelectMesh from './WifiSelectMesh';
import {Route} from 'react-router-dom';
import type {MetricGraphConfig} from '@fbcnms/magmalte/app/components/insights/Metrics';
import type {TimeRange} from '@fbcnms/magmalte/app/components/insights/AsyncMetric';
import type {WifiGateway} from './WifiUtils';

import {buildWifiGatewayFromPayload, meshesURL} from './WifiUtils';

import {MagmaAPIUrls} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {find} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {resolveQuery} from '@fbcnms/magmalte/app/components/insights/Metrics';
import {useAxios, useRouter, useSnackbar} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  appBar: {
    display: 'inline-block',
  },
  chartRow: {
    display: 'flex',
  },
  formControl: {
    minWidth: '400px',
    padding: theme.spacing(),
  },
  formControlPeriod: {
    minWidth: '200px',
    padding: theme.spacing(),
  },
}));

const CHART_CONFIGS: MetricGraphConfig[] = [
  {
    id: 'ap_hops_to_gateway',
    label: 'Hops to Gateway',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: gw =>
          `ap_hops_to_gateway{gatewayID=~"${gw}", service="linkstatsd"}`,
      },
    ],
  },
  {
    id: 'gateway_cpu',
    label: 'CPU (%)',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: gw => `cpu_percent{gatewayID=~"${gw}", service="magmad"}`,
      },
    ],
    unit: '%',
  },
  {
    id: 'mem',
    label: 'Memory',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: gw =>
          `mem_available{gatewayID=~"${gw}", service="magmad"}`,
      },
    ],
    unit: 'bytes',
    // transform: 'formula(/ $1 1048576)',
  },
  {
    id: 'disk',
    label: 'Disk (%)',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: gw =>
          `disk_percent{gatewayID=~"${gw}", service="magmad"}`,
      },
    ],
    unit: '%',
  },
];

function Metrics(props: {parentRelativeUrl: string => string}) {
  const {history, match} = useRouter();
  const classes = useStyles();
  const {meshId, deviceId} = match.params;
  const [timeRange, setTimeRange] = useState<TimeRange>('24_hours');
  const {parentRelativeUrl} = props;

  const {
    error: meshesError,
    isLoading: meshesIsLoading,
    response: meshesResponse,
  } = useAxios({
    method: 'get',
    url: meshesURL(match),
  });

  useSnackbar('Error fetching meshes', {variant: 'error'}, meshesError);

  const {
    error: devicesError,
    isLoading: devicesIsLoading,
    response: devicesResponse,
  } = useAxios({
    method: 'get',
    url: MagmaAPIUrls.gateways(match, true),
  });

  useSnackbar('Error fetching devices', {variant: 'error'}, devicesError);

  if (
    devicesError ||
    meshesError ||
    devicesIsLoading ||
    !devicesResponse ||
    !devicesResponse.data ||
    meshesIsLoading ||
    !meshesResponse ||
    !meshesResponse.data
  ) {
    return <LoadingFiller />;
  }

  const meshes: string[] = meshesResponse.data || [];
  if (!meshes || meshes.length <= 0) {
    return <Text variant="h5">No meshes on this network</Text>;
  }

  const selectedMeshOrDefault = meshId || meshes[0];

  const devices: Array<WifiGateway> = (devicesResponse.data || [])
    // TODO: skip filter when magma API bug fixed t34643616
    .filter(device => device.record && device.config)
    .map(rawGateway => buildWifiGatewayFromPayload(rawGateway))
    // filter by selected meshID
    .filter(
      device =>
        selectedMeshOrDefault === null ||
        device.meshid === selectedMeshOrDefault,
    );
  devices.sort((d1, d2) =>
    d1.info.toLowerCase() > d2.info.toLowerCase() ? 1 : -1,
  );
  const selectedDevice = find(devices, device => device.id === deviceId);
  const selectedDeviceId = selectedDevice?.id || '';

  const onMeshSelectorChange = targetMesh => {
    if (targetMesh === '') {
      history.push(parentRelativeUrl('/'));
    } else {
      history.push(parentRelativeUrl(`/${targetMesh}`));
    }
  };
  const onDeviceSelectorChange = ({target}) => {
    if (!target.value) {
      history.push(parentRelativeUrl(`/${selectedMeshOrDefault}`));
    } else {
      history.push(
        parentRelativeUrl(`/${selectedMeshOrDefault}/${target.value}`),
      );
    }
  };

  const selectors = (
    <AppBar className={classes.appBar} position="static" color="default">
      <WifiSelectMesh
        meshes={meshes}
        onChange={onMeshSelectorChange}
        selectedMeshID={selectedMeshOrDefault}
        helperText=""
        disallowEmpty={true}
        classes={{formControl: classes.formControl}}
      />
      <FormControl variant="filled" className={classes.formControl}>
        <InputLabel htmlFor="devices">Device</InputLabel>
        <Select
          inputProps={{id: 'devices'}}
          value={selectedDeviceId}
          onChange={onDeviceSelectorChange}>
          <MenuItem value={''} key={''}>
            All
          </MenuItem>
          {devices.map(device => (
            <MenuItem value={device.id} key={device.id}>
              {device.info} : {device.id}
            </MenuItem>
          ))}
        </Select>
      </FormControl>
      <TimeRangeSelector
        className={classes.formControl}
        value={timeRange}
        onChange={setTimeRange}
      />
    </AppBar>
  );

  if (!devices || devices.length <= 0) {
    return (
      <>
        {selectors}
        <Text variant="h5">No devices on this mesh</Text>
      </>
    );
  }

  return (
    <>
      {selectors}
      <GridList cols={2} cellHeight={300}>
        {CHART_CONFIGS.map((config, i) => (
          <GridListTile key={i} cols={1}>
            <Card>
              <CardContent>
                <Text component="h6" variant="h6">
                  {config.label}
                </Text>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    queries={resolveQuery(
                      config,
                      'gatewayID',
                      selectedDeviceId || `${selectedMeshOrDefault}_id_.*`,
                    )}
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
  const {relativePath, relativeUrl} = useRouter();
  return (
    <Route
      path={relativePath('/:meshId?/:deviceId?')}
      render={() => <Metrics parentRelativeUrl={relativeUrl} />}
    />
  );
}
