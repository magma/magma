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
import type {ChartStyle} from '@fbcnms/ui/insights/AsyncMetric';

import AsyncMetric from '@fbcnms/ui/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '../theme/design-system/Text';
import moment from 'moment';

import {colors} from '../theme/default';

type Props = {
  startEnd: [moment, moment],
};
const CHART_COLORS = [
  colors.secondary.dodgerBlue,
  colors.data.flamePea,
  'green',
  'yellow',
  'purple',
  'black',
];

const isValid = (start, end): boolean => {
  return start.isValid() && end.isValid() && moment.min(start, end) === start;
};

export default function ({startEnd}: Props) {
  const [start, end] = startEnd;
  const state = {
    title: 'Frequency of Alerts and Events',
    legendLabels: ['Alerts', 'Events'],
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
  return (
    <Grid>
      <Card elevation={0}>
        <CardHeader
          title={<Text variant="body2">{state.title}</Text>}
          subheader={
            <AsyncMetric
              style={chartStyle}
              label={state.title}
              unit=""
              queries={['sum(ALERTS)']}
              timeRange={'3_hours'}
              startEnd={isValid(start, end) ? startEnd : undefined}
              legendLabels={state.legendLabels}
              chartColors={CHART_COLORS}
            />
          }
        />
      </Card>
    </Grid>
  );
}
