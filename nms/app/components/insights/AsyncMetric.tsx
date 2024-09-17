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

import CircularProgress from '@mui/material/CircularProgress';
import MagmaAPI from '../../api/MagmaAPI';
import React from 'react';
import Text from '../../theme/design-system/Text';
import {LayoutPosition, TimeUnit} from 'chart.js';
import {Line} from 'react-chartjs-2';
import {PromqlMetric, PromqlMetricValue} from '../../../generated';
import {defaultTooltip} from '../CustomMetrics';
import {differenceInDays, differenceInHours, getUnixTime, sub} from 'date-fns';
import {makeStyles} from '@mui/styles';
import {useEffect, useMemo, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';

type AxesOptions = {
  gridLines: {
    drawBorder?: boolean;
    display?: boolean;
  };
  ticks: {
    maxTicksLimit: number;
  };
};

export type ChartStyle = {
  options: {
    xAxes: AxesOptions;
    yAxes: AxesOptions;
  };
  data: {
    lineTension: number;
    pointRadius: number;
  };
  legend: {
    position: LayoutPosition;
    align: 'center' | 'end' | 'start';
  };
};

type Props = {
  label: string;
  unit: string;
  queries: Array<string>;
  legendLabels?: Array<string>;
  timeRange: TimeRange;
  startEnd?: [Date, Date];
  networkId?: string;
  style?: ChartStyle;
  height?: number;
  chartColors?: Array<string>;
};

const useStyles = makeStyles({
  loadingContainer: {
    paddingTop: 100,
    textAlign: 'center',
  },
});

export type TimeRange =
  | '3_hours'
  | '6_hours'
  | '12_hours'
  | '24_hours'
  | '7_days'
  | '14_days'
  | '30_days';

type RangeValue = {
  days?: number;
  hours?: number;
  step: string;
  unit: TimeUnit;
};

type Dataset = {
  label: string;
  unit: string;
  fill: boolean;
  lineTension: number;
  pointRadius: number;
  pointHitRadius: number;
  borderWidth: number;
  backgroundColor: string;
  borderColor: string;
  data: Array<{x: number; y: number | string}>;
};

const RANGE_VALUES: Record<TimeRange, RangeValue> = {
  '3_hours': {
    hours: 3,
    step: '30s',
    unit: 'minute',
  },
  '6_hours': {
    hours: 6,
    step: '1m',
    unit: 'hour',
  },
  '12_hours': {
    hours: 12,
    step: '5m',
    unit: 'hour',
  },
  '24_hours': {
    days: 1,
    step: '15m',
    unit: 'hour',
  },
  '7_days': {
    days: 7,
    step: '2h',
    unit: 'day',
  },
  '14_days': {
    days: 14,
    step: '4h',
    unit: 'day',
  },
  '30_days': {
    days: 30,
    step: '8h',
    unit: 'day',
  },
};

const COLORS = ['blue', 'red', 'green', 'yellow', 'purple', 'black'];

class PrometheusHelper {
  getLegendLabel(
    result: PromqlMetricValue,
    tagSets: Array<PromqlMetric>,
  ): string {
    const {metric} = result as {metric: Record<string, string>};

    const tags = [];
    const droppedTags = ['networkID', '__name__'];
    const droppedIfSameTags = ['gatewayID', 'service'];

    const uniqueTagValues: Record<string, Array<string>> = {};
    droppedIfSameTags.forEach(tagName => {
      uniqueTagValues[tagName] = Array.from(
        new Set(
          (tagSets as Array<Record<string, string>>).map(item => item[tagName]),
        ),
      );
    });

    for (const key in metric) {
      if (
        metric.hasOwnProperty(key) &&
        !droppedTags.includes(key) &&
        (!uniqueTagValues[key] || uniqueTagValues[key].length !== 1)
      ) {
        tags.push(key + '=' + metric[key]);
      }
    }
    return tags.length === 0
      ? metric['__name__']
      : `${metric['__name__']} (${tags.join(', ')})`;
  }

  datapointFieldName = 'values' as const;
}

function Progress() {
  const classes = useStyles();
  return (
    <div className={classes.loadingContainer}>
      <CircularProgress />
    </div>
  );
}

function getStartEnd(timeRange: TimeRange) {
  const {days, hours, step} = RANGE_VALUES[timeRange];
  const end = new Date();
  const endUnix = getUnixTime(end) * 1000;
  const start = sub(end, {days, hours});
  const startUnix = getUnixTime(start) * 1000;
  return {
    start: start.toISOString(),
    startUnix: startUnix,
    end: end.toISOString(),
    endUnix: endUnix,
    step,
  };
}

function getUnit(timeRange: TimeRange) {
  return RANGE_VALUES[timeRange].unit;
}

function getStepUnit(startEnd: [Date, Date]): [string, TimeUnit] {
  const [start, end] = startEnd;
  const durationInHours = differenceInHours(end, start);
  const durationInDays = differenceInDays(end, start);

  let range: RangeValue;
  if (durationInHours <= 24) {
    if (durationInHours <= 3) {
      range = RANGE_VALUES['3_hours'];
    } else if (durationInHours <= 6) {
      range = RANGE_VALUES['6_hours'];
    } else if (durationInHours <= 12) {
      range = RANGE_VALUES['12_hours'];
    } else {
      range = RANGE_VALUES['24_hours'];
    }
  } else {
    if (durationInDays <= 7) {
      range = RANGE_VALUES['7_days'];
    } else if (durationInDays <= 14) {
      range = RANGE_VALUES['14_days'];
    } else {
      range = RANGE_VALUES['30_days'];
    }
  }

  return [range.step, range.unit];
}

function getColorForIndex(index: number, customChartColors?: Array<string>) {
  if (customChartColors != null) {
    return customChartColors[index % customChartColors.length];
  }
  return COLORS[index % COLORS.length];
}

function useDatasetsFetcher(props: Props) {
  const params = useParams();
  const startEnd = useMemo(() => {
    if (props.startEnd) {
      const [start, end] = props.startEnd;
      const [step] = getStepUnit(props.startEnd);
      return {
        start: start.toISOString(),
        startUnix: getUnixTime(start) * 1000,
        end: end.toISOString(),
        endUnix: getUnixTime(end) * 1000,
        step,
      };
    } else {
      return getStartEnd(props.timeRange);
    }
  }, [props.timeRange, props.startEnd]);

  const [allDatasets, setAllDatasets] = useState<Array<Dataset> | null>(null);
  const enqueueSnackbar = useEnqueueSnackbar();
  const stringedQueries = JSON.stringify(props.queries);

  const dbHelper = useMemo(() => new PrometheusHelper(), []);

  useEffect(
    () => {
      const queries = props.queries;
      const requests = queries.map(async (query, index) => {
        try {
          // eslint-disable-next-line max-len
          const response = (
            await MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet({
              networkId: props.networkId || params.networkId!,
              start: startEnd.start,
              end: startEnd.end,
              step: startEnd.step,
              query,
            })
          ).data;
          const label = props.legendLabels ? props.legendLabels[index] : null;
          return {response, label};
        } catch (error) {
          enqueueSnackbar('Error getting metric ' + props.label, {
            variant: 'error',
          });
        }
        return null;
      });

      void Promise.all(requests).then(allResponses => {
        let index = 0;
        const datasets: Array<Dataset> = [];
        const {style} = props;
        allResponses.filter(Boolean).forEach(r => {
          const response = r!.response;
          const label = r!.label;
          const result = response.data.result;
          if (result) {
            const tagSets = result.map(it => it.metric);
            result.map(it =>
              datasets.push({
                label: label || dbHelper.getLegendLabel(it, tagSets),
                unit: props.unit || '',
                fill: false,
                lineTension: style ? style.data.lineTension : 0,
                pointHitRadius: 10,
                pointRadius: style ? style.data.pointRadius : 0,
                borderWidth: 2,
                backgroundColor: getColorForIndex(index, props.chartColors),
                borderColor: getColorForIndex(index++, props.chartColors),
                data: it[dbHelper.datapointFieldName]!.map(i => ({
                  x: parseInt(i[0]) * 1000,
                  y: parseFloat(i[1]),
                })),
              }),
            );
          }
        });
        // Add "NaN" to the beginning/end of each dataset to force the chart to
        // display the whole time frame requested
        datasets.forEach(dataset => {
          if (dataset.data[0].x > startEnd.startUnix) {
            dataset.data.unshift({x: startEnd.startUnix, y: 'NaN'});
          }
          if (dataset.data[dataset.data.length - 1].x < startEnd.endUnix) {
            dataset.data.push({x: startEnd.endUnix, y: 'NaN'});
          }
        });
        setAllDatasets(datasets);
      });
    } /* eslint-disable react-hooks/exhaustive-deps */,
    [
      stringedQueries,
      params,
      props.networkId,
      props.unit,
      startEnd,
      props.label,
      props.legendLabels,
      enqueueSnackbar,
      dbHelper,
    ],
  );
  /* eslint-enable react-hooks/exhaustive-deps */

  return allDatasets;
}

export default function AsyncMetric(props: Props) {
  const allDatasets = useDatasetsFetcher(props);
  if (allDatasets === null) {
    return <Progress />;
  }

  if (!allDatasets || allDatasets?.length === 0) {
    return <Text variant="body2">No Data</Text>;
  }
  const {style} = props;
  const {startEnd} = props;
  let unit: TimeUnit;
  if (startEnd) {
    [, unit] = getStepUnit(startEnd);
  } else {
    unit = getUnit(props.timeRange);
  }

  return (
    <Line
      height={props.height}
      options={{
        maintainAspectRatio: false,
        scales: {
          x: {
            grid: style ? style.options.xAxes.gridLines : {},
            ticks: style ? style.options.xAxes.ticks : {},
            type: 'time',
            time: {
              unit,
              round: 'second',
              tooltipFormat: ' yyyy/MM/dd h:mm:ss a',
            },
            title: {
              display: true,
              text: 'Date',
            },
          },
          y: {
            grid: style ? style.options.yAxes.gridLines : {},
            ticks: style ? style.options.yAxes.ticks : {},
            position: 'left',
            title: {
              display: true,
              text: props.unit,
            },
          },
        },
        plugins: {
          tooltip: {
            enabled: true,
            mode: 'nearest',
            callbacks: {
              label(tooltipItem) {
                return defaultTooltip(tooltipItem, props);
              },
            },
          },
          legend: {
            display: allDatasets.length < 5,
            position: style ? style.legend.position : 'bottom',
            align: style ? style.legend.align : 'center',
            labels: {
              boxWidth: 12,
            },
          },
        },
      }}
      data={{datasets: allDatasets}}
    />
  );
}
