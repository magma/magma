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
 */

import type {MetricGraphConfig} from './Metrics';
import type {TimeRange} from './AsyncMetric';

import AppBar from '@mui/material/AppBar';
import AsyncMetric from './AsyncMetric';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import ImageList from '@mui/material/ImageList';
import ImageListItem from '@mui/material/ImageListItem';
import React from 'react';
import Text from '../../theme/design-system/Text';
import TimeRangeSelector from '../insights/TimeRangeSelector';

import {Theme} from '@mui/material/styles';
import {makeStyles} from '@mui/styles';
import {resolveQuery} from './Metrics';
import {useState} from 'react';

const useStyles = makeStyles<Theme>(theme => ({
  formControl: {
    minWidth: '200px',
    padding: theme.spacing(),
  },
  appBar: {
    display: 'inline-block',
  },
}));

export default function (props: {configs: Array<MetricGraphConfig>}) {
  const classes = useStyles();
  const [timeRange, setTimeRange] = useState<TimeRange>('3_hours');

  return (
    <>
      <AppBar className={classes.appBar} position="static" color="default">
        <TimeRangeSelector
          className={classes.formControl}
          value={timeRange}
          onChange={setTimeRange}
        />
      </AppBar>
      <ImageList cols={2} rowHeight={300}>
        {props.configs.map((config, i) => (
          <ImageListItem key={i} cols={1}>
            <Card>
              <CardContent>
                <Text variant="h6">{config.label}</Text>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    queries={resolveQuery(config, '', '')}
                    timeRange={timeRange}
                    legendLabels={config.legendLabels}
                  />
                </div>
              </CardContent>
            </Card>
          </ImageListItem>
        ))}
      </ImageList>
    </>
  );
}
