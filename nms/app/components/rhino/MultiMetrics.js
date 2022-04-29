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

import AsyncMetric from '../insights/AsyncMetric';
import {resolveQuery} from '../insights/Metrics';
import type {MetricGraphConfig} from '../insights/Metrics';
import type {TimeRange} from '../insights/AsyncMetric';

import AppBar from '@material-ui/core/AppBar';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import FormControl from '@material-ui/core/FormControl';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFiller from '../../../fbc_js_core/ui/components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useState} from 'react';
import Select from '@material-ui/core/Select';
import Text from '../../theme/design-system/Text';
import TimeRangeSelector from '../insights/TimeRangeSelector';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';

import useMagmaAPI from '../../../api/useMagmaAPI';
import {find, map} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useSnackbar} from '../../../fbc_js_core/ui/hooks';

const useStyles = makeStyles(theme => ({
  appBar: {
    display: 'inline-block',
  },
  chartRow: {
    display: 'flex',
  },
  formControl: {
    minWidth: '200px',
    padding: theme.spacing(),
  },
}));

function MultiMetrics(props: {
  onGatewaySelectorChange: (SyntheticInputEvent<EventTarget>) => void,
  configs: Array<MetricGraphConfig>,
}) {
  const params = useParams();
  const classes = useStyles();
  const selectedGateway = params.selectedGatewayId;
  const [timeRange, setTimeRange] = useState<TimeRange>('24_hours');

  const {error, isLoading, response: paginated_gateways} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdGateways,
    {
      networkId: params.networkId,
    },
  );

  useSnackbar('Error fetching devices', {variant: 'error'}, error);

  if (error || isLoading || !paginated_gateways) {
    return <LoadingFiller />;
  }

  const defaultGateway = find(
    paginated_gateways.gateways,
    gateway => gateway.status?.hardware_id !== null,
  );
  const selectedGatewayOrDefault = selectedGateway || defaultGateway?.id;

  return (
    <>
      <AppBar className={classes.appBar} position="static" color="default">
        <FormControl variant="filled" className={classes.formControl}>
          <InputLabel htmlFor="devices">Device</InputLabel>
          <Select
            inputProps={{id: 'devices'}}
            value={selectedGatewayOrDefault}
            onChange={props.onGatewaySelectorChange}>
            {map(paginated_gateways.gateways, device => (
              <MenuItem value={device.id} key={device.id}>
                {device.name}
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
      <GridList cols={2} cellHeight={400}>
        {props.configs.map((config, i) => (
          <GridListTile key={i} cols={1}>
            <Card>
              <CardContent>
                <Text component="h6" variant="h6">
                  {config.label}
                </Text>
                <div style={{height: 350}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    queries={resolveQuery(
                      config,
                      'gatewayID',
                      selectedGatewayOrDefault,
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

export default function (props: {configs: Array<MetricGraphConfig>}) {
  const navigate = useNavigate();

  const metrics = (
    <MultiMetrics
      configs={props.configs}
      onGatewaySelectorChange={({target}) => navigate(target.value)}
    />
  );

  return (
    <Routes>
      <Route path=":selectedGatewayId" element={metrics} />
      <Route index element={metrics} />
    </Routes>
  );
}
