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

import CircularProgress from '@material-ui/core/CircularProgress';
import MagmaV1API from '../../../generated/WebClient';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import moment from 'moment';
import {Line} from 'react-chartjs-2';

import {makeStyles} from '@material-ui/styles';
import {useEffect, useMemo, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

type AxesOptions = {
  gridLines: {
    drawBorder?: boolean,
    display?: boolean,
  },
  ticks: {
    maxTicksLimit: number,
  },
};

export type ChartStyle = {
  options: {
    xAxes: AxesOptions,
    yAxes: AxesOptions,
  },
  data: {
    lineTension: number,
    pointRadius: number,
  },
  legend: {
    position: string,
    align: string,
  },
};

type Props = {
  label: string,
  unit: string,
  queries: Array<string>,
  legendLabels?: Array<string>,
  timeRange: TimeRange,
  startEnd?: [moment, moment],
  networkId?: string,
  style?: ChartStyle,
  height?: number,
  chartColors?: Array<string>,
};

const useStyles = makeStyles(() => ({
  loadingContainer: {
    paddingTop: 100,
    textAlign: 'center',
  },
}));

export type TimeRange =
  | '3_hours'
  | '6_hours'
  | '12_hours'
  | '24_hours'
  | '7_days'
  | '14_days'
  | '30_days';

type RangeValue = {
  days?: number,
  hours?: number,
  step: string,
  unit: string,
};

type Dataset = {
  label: string,
  unit: string,
  fill: boolean,
  lineTension: number,
  pointRadius: number,
  borderWidth: number,
  backgroundColor: string,
  borderColor: string,
  data: Array<{t: number, y: number | string}>,
};

const RANGE_VALUES: {[TimeRange]: RangeValue} = {
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

interface DatabaseHelper<T> {
  getLegendLabel(data: T, tagSets: Array<{[string]: string}>): string;
  datapointFieldName: string;
}

class PrometheusHelper implements DatabaseHelper<PrometheusResponse> {
  getLegendLabel(
    result: PrometheusResponse,
    tagSets: Array<{[string]: string}>,
  ): string {
    const {metric} = result;

    const tags = [];
    const droppedTags = ['networkID', '__name__'];
    const droppedIfSameTags = ['gatewayID', 'service'];

    const uniqueTagValues = {};
    droppedIfSameTags.forEach(tagName => {
      uniqueTagValues[tagName] = Array.from(
        new Set(tagSets.map(item => item[tagName])),
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

  datapointFieldName = 'values';
}

type PrometheusResponse = {
  metric: {[key: string]: string},
};

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
  const end = moment();
  const endUnix = end.unix() * 1000;
  const start = end.clone().subtract({days, hours});
  const startUnix = start.unix() * 1000;
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

function getStepUnit(startEnd: [moment, moment]): [string, string] {
  const [start, end] = startEnd;
  const d = moment.duration(end.diff(start));
  const hrs = d.asHours();
  const days = d.asDays();
  let r: RangeValue;
  if (hrs <= 24) {
    if (hrs <= 3) {
      r = RANGE_VALUES['3_hours'];
    } else if (hrs <= 6) {
      r = RANGE_VALUES['6_hours'];
    } else if (hrs <= 12) {
      r = RANGE_VALUES['12_hours'];
    } else {
      r = RANGE_VALUES['24_hours'];
    }
  } else {
    if (days <= 7) {
      r = RANGE_VALUES['7_days'];
    } else if (days <= 14) {
      r = RANGE_VALUES['14_days'];
    } else {
      r = RANGE_VALUES['30_days'];
    }
  }
  return [r.step, r.unit];
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
        startUnix: start.unix() * 1000,
        end: end.toISOString(),
        endUnix: end.unix() * 1000,
        step,
      };
    } else {
      return getStartEnd(props.timeRange);
    }
  }, [props.timeRange, props.startEnd]);

  const [allDatasets, setAllDatasets] = useState<?Array<Dataset>>(null);
  const enqueueSnackbar = useEnqueueSnackbar();
  const stringedQueries = JSON.stringify(props.queries);

  const dbHelper = useMemo(() => new PrometheusHelper(), []);

  useEffect(
    () => {
      const queries = props.queries;
      const requests = queries.map(async (query, index) => {
        try {
          // eslint-disable-next-line max-len
          const response = await MagmaV1API.getNetworksByNetworkIdPrometheusQueryRange(
            {
              // $FlowFixMe[sketchy-null-string] TODO(andreilee): from fbcnms-ui
              networkId: props.networkId || params.networkId,
              start: startEnd.start,
              end: startEnd.end,
              step: startEnd.step,
              query,
            },
          );
          const label = props.legendLabels ? props.legendLabels[index] : null;
          return {response, label};
        } catch (error) {
          enqueueSnackbar('Error getting metric ' + props.label, {
            variant: 'error',
          });
        }
        return null;
      });

      Promise.all(requests).then(allResponses => {
        let index = 0;
        const datasets = [];
        const {style} = props;
        allResponses.filter(Boolean).forEach(r => {
          const response = r.response;
          const label = r.label;
          const result = response.data.result;
          if (result) {
            const tagSets = result.map(it => it.metric);
            result.map(it =>
              datasets.push({
                // $FlowFixMe[sketchy-null-string] TODO(andreilee): from fbcnms-ui
                label: label || dbHelper.getLegendLabel(it, tagSets),
                unit: props.unit || '',
                fill: false,
                lineTension: style ? style.data.lineTension : 0,
                pointHitRadius: 10,
                pointRadius: style ? style.data.pointRadius : 0,
                borderWidth: 2,
                backgroundColor: getColorForIndex(index, props.chartColors),
                borderColor: getColorForIndex(index++, props.chartColors),
                data: it[dbHelper.datapointFieldName].map(i => ({
                  t: parseInt(i[0]) * 1000,
                  y: parseFloat(i[1]),
                })),
              }),
            );
          }
        });
        // Add "NaN" to the beginning/end of each dataset to force the chart to
        // display the whole time frame requested
        datasets.forEach(dataset => {
          if (dataset.data[0].t > startEnd.startUnix) {
            dataset.data.unshift({t: startEnd.startUnix, y: 'NaN'});
          }
          if (dataset.data[dataset.data.length - 1].t < startEnd.endUnix) {
            dataset.data.push({t: startEnd.endUnix, y: 'NaN'});
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
  let unit: string;
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
        scaleShowValues: true,
        scales: {
          xAxes: [
            {
              gridLines: style ? style.options.xAxes.gridLines : {},
              ticks: style ? style.options.xAxes.ticks : {},
              type: 'time',
              time: {
                unit,
                round: 'second',
                tooltipFormat: ' YYYY/MM/DD h:mm:ss a',
              },
              scaleLabel: {
                display: true,
                labelString: 'Date',
              },
            },
          ],
          yAxes: [
            {
              gridLines: style ? style.options.yAxes.gridLines : {},
              ticks: style ? style.options.yAxes.ticks : {},
              position: 'left',
              scaleLabel: {
                display: true,
                labelString: props.unit,
              },
            },
          ],
        },
        tooltips: {
          enabled: true,
          mode: 'nearest',
          callbacks: {
            label: (tooltipItem, data) =>
              data.datasets[tooltipItem.datasetIndex].label +
              ': ' +
              tooltipItem.yLabel +
              ' ' +
              data.datasets[tooltipItem.datasetIndex].unit,
          },
        },
      }}
      legend={{
        display: allDatasets.length < 5,
        position: style ? style.legend.position : 'bottom',
        align: style ? style.legend.align : 'center',
        labels: {
          boxWidth: 12,
        },
      }}
      data={{datasets: allDatasets}}
    />
  );
}
