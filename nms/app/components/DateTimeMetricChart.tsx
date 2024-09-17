/*
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

import AsyncMetric from './insights/AsyncMetric';
import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import CardTitleRow from './layout/CardTitleRow';
import DataUsageIcon from '@mui/icons-material/DataUsage';
import Grid from '@mui/material/Grid';
import React from 'react';
import Text from '../theme/design-system/Text';

import TextField from '@mui/material/TextField';
import {DateTimePicker} from '@mui/x-date-pickers/DateTimePicker';
import {colors} from '../theme/default';
import {makeStyles} from '@mui/styles';
import {subHours} from 'date-fns';
import {useState} from 'react';

export type DateTimeMetricChartProps = {
  title: string;
  queries: Array<string>;
  legendLabels: Array<string>;
  unit?: string;
  startDate?: Date;
  endDate?: Date;
};

const useStyles = makeStyles({
  dateTimeText: {
    color: colors.primary.comet,
  },
});

const CHART_COLORS = [colors.secondary.dodgerBlue, colors.data.flamePea];

export default function DateTimeMetricChart(props: DateTimeMetricChartProps) {
  const classes = useStyles();
  const [startDate, setStartDate] = useState(subHours(new Date(), 3));
  const [endDate, setEndDate] = useState(new Date());

  function Filter() {
    return (
      <Grid container justifyContent="flex-end" alignItems="center" spacing={1}>
        <Grid item>
          <Text variant="body3" className={classes.dateTimeText}>
            Filter By Date
          </Text>
        </Grid>
        <Grid item>
          <DateTimePicker
            renderInput={props => <TextField {...props} />}
            maxDate={endDate}
            disableFuture
            value={startDate}
            onChange={date => setStartDate(date!)}
          />
        </Grid>
        <Grid item>
          <Text variant="body3" className={classes.dateTimeText}>
            to
          </Text>
        </Grid>
        <Grid item>
          <DateTimePicker
            renderInput={props => <TextField {...props} />}
            disableFuture
            value={endDate}
            onChange={date => setEndDate(date!)}
          />
        </Grid>
      </Grid>
    );
  }

  return (
    <>
      {!(props.startDate && props.endDate) && (
        <CardTitleRow
          icon={DataUsageIcon}
          label={props.title}
          filter={Filter}
        />
      )}
      <Card elevation={0}>
        <CardHeader
          title={<Text variant="body2">{props.title}</Text>}
          subheader={
            <AsyncMetric
              height={300}
              style={{
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
              }}
              label={`Frequency of ${props.title}`}
              unit={props.unit ?? ''}
              queries={props.queries}
              timeRange={'3_hours'}
              startEnd={[
                props.startDate ?? startDate,
                props.endDate ?? endDate,
              ]}
              legendLabels={props.legendLabels}
              chartColors={CHART_COLORS}
            />
          }
        />
      </Card>
    </>
  );
}
