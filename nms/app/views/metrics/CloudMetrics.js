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

import type {MetricGraphConfig} from '../../components/insights/Metrics';
import type {TimeRange} from '../../components/insights/AsyncMetric';

import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from '../../components/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import TimeRangeSelector from '../../components/insights/TimeRangeSelector';

import {makeStyles} from '@material-ui/styles';
import {resolveQuery} from '../../components/insights/Metrics';
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
  const [timeRange, setTimeRange] = useState<TimeRange>('24_hours');

  const chartConfigs: MetricGraphConfig[] = [
    {
      label: 'REST API (2xx status code)',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => "sum(response_status{code=~'2..'})",
        },
      ],
      unit: '',
    },
    {
      label: 'REST API (3xx status code)',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => "sum(response_status{code=~'3..'})",
        },
      ],
      unit: '',
    },
    {
      label: 'REST API (4xx status code)',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => "sum(response_status{code=~'4..'})",
        },
      ],
      unit: '',
    },
    {
      label: 'REST API (5xx status code)',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => "sum(response_status{code=~'5..'})",
        },
      ],
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
                <Text component="h6" variant="h6">
                  {config.label}
                </Text>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    queries={resolveQuery(config, 'gatewayID', '')}
                    timeRange={timeRange}
                    networkId="cloud"
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
