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
import AsyncMetric from '../../components/insights/AsyncMetric';
import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DataUsageIcon from '@mui/icons-material/DataUsage';
import Grid from '@mui/material/Grid';
import React, {useState} from 'react';
import Text from '../../theme/design-system/Text';
import TimeRangeSelector from '../../theme/design-system/TimeRangeSelector';
import {Theme} from '@mui/material/styles';
import {colors} from '../../theme/default';
import {makeStyles} from '@mui/styles';
import type {
  ChartStyle,
  TimeRange,
} from '../../components/insights/AsyncMetric';

const useStyles = makeStyles<Theme>(theme => ({
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
          maxTicksLimit: 3,
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
