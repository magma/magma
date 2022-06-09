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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {ChartStyle} from '../../components/insights/AsyncMetric';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {TimeRange} from '../../components/insights/AsyncMetric';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AsyncMetric from '../../components/insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import DataUsageIcon from '@material-ui/icons/DataUsage';
import Grid from '@material-ui/core/Grid';
import React, {useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import TimeRangeSelector from '../../theme/design-system/TimeRangeSelector';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
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

const CHART_COLORS = [colors.secondary.dodgerBlue];
const TITLE = 'Frequency of Gateway Check-Ins';

export default function () {
  const classes = useStyles();
  const [timeRange, setTimeRange] = useState<TimeRange>('3_hours');

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
      <Grid container justifyContent="flex-end" alignItems="center" spacing={1}>
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
      <CardTitleRow
        icon={DataUsageIcon}
        label="Gateway Check-Ins"
        filter={Filter}
      />
      <Card elevation={0}>
        <CardHeader
          title={<Text variant="body2">{TITLE}</Text>}
          subheader={
            <AsyncMetric
              height={300}
              style={chartStyle}
              label={TITLE}
              unit="Count"
              queries={['sum(checkin_status)']}
              timeRange={timeRange}
              legendLabels={['Check-Ins']}
              chartColors={CHART_COLORS}
            />
          }
        />
      </Card>
    </>
  );
}
