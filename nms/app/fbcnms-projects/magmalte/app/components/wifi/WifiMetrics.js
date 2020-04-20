/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import type {MetricGraphConfig} from '@fbcnms/ui/insights/Metrics';
import type {TimeRange} from '@fbcnms/ui/insights/AsyncMetric';
import type {WifiGateway} from './WifiUtils';

import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from '@fbcnms/ui/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import FormControl from '@material-ui/core/FormControl';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import TimeRangeSelector from '@fbcnms/ui/insights/TimeRangeSelector';
import WifiSelectMesh from './WifiSelectMesh';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {Route} from 'react-router-dom';
import {buildWifiGatewayFromPayloadV1} from './WifiUtils';
import {find} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {map} from 'lodash';
import {resolveQuery} from '@fbcnms/ui/insights/Metrics';
import {useRouter, useSnackbar} from '@fbcnms/ui/hooks';
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
  } = useMagmaAPI(MagmaV1API.getWifiByNetworkIdMeshes, {
    networkId: nullthrows(match.params.networkId),
  });

  useSnackbar('Error fetching meshes', {variant: 'error'}, meshesError);

  const {
    error: devicesError,
    isLoading: devicesIsLoading,
    response: devicesResponse,
  } = useMagmaAPI(MagmaV1API.getWifiByNetworkIdGateways, {
    networkId: nullthrows(match.params.networkId),
  });

  useSnackbar('Error fetching devices', {variant: 'error'}, devicesError);

  if (
    devicesError ||
    meshesError ||
    devicesIsLoading ||
    !devicesResponse ||
    meshesIsLoading ||
    !meshesResponse
  ) {
    return <LoadingFiller />;
  }

  const meshes: string[] = meshesResponse || [];
  if (!meshes || meshes.length <= 0) {
    return <Text variant="h5">No meshes on this network</Text>;
  }

  const selectedMeshOrDefault = meshId || meshes[0];

  const devices: Array<WifiGateway> = map(devicesResponse || {})
    .filter(device => device.device)
    .map(device => buildWifiGatewayFromPayloadV1(device))
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

export default function () {
  const {relativePath, relativeUrl} = useRouter();
  return (
    <Route
      path={relativePath('/:meshId?/:deviceId?')}
      render={() => <Metrics parentRelativeUrl={relativeUrl} />}
    />
  );
}
