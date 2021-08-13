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
import type {ChartData, ChartTooltipItem} from 'react-chartjs-2';

import React from 'react';
import moment from 'moment';

import {Bar, Line} from 'react-chartjs-2';

export function getStepString(delta: number, unit: string) {
  return delta.toString() + unit[0];
}

export function getStep(start: moment, end: moment): [number, string, string] {
  const d = moment.duration(end.diff(start));
  if (d.asMinutes() <= 60.5) {
    return [5, 'minute', 'HH:mm'];
  } else if (d.asHours() <= 3.5) {
    return [15, 'minute', 'HH:mm'];
  } else if (d.asHours() <= 6.5) {
    return [15, 'minute', 'HH:mm'];
  } else if (d.asHours() <= 12.5) {
    return [1, 'hour', 'HH:mm'];
  } else if (d.asHours() <= 24.5) {
    return [2, 'hour', 'HH:mm'];
  } else if (d.asDays() <= 1.5) {
    return [3, 'hour', 'DD-MM-YY HH:mm'];
  } else if (d.asDays() <= 3.5) {
    return [6, 'hour', 'DD-MM-YY HH:mm'];
  } else if (d.asDays() <= 7.5) {
    return [12, 'hour', 'DD-MM-YY HH:mm'];
  }
  return [24, 'hour', 'DD-MM-YYYY'];
}

// for querying event and log count, the api doesn't have a step attribute
// hence we have to split the start and end window into several sets of
// [start, end] queries which can then be queried in parallel
export function getQueryRanges(
  start: moment,
  end: moment,
  delta: number,
  unit: string,
): Array<[moment, moment]> {
  const queries = [];
  let s = start.clone();
  // go back delta time so that we get the total number of events
  // or logs at that 's' point of time
  s = s.subtract(delta, unit);
  while (end.diff(s, unit) >= delta) {
    const e = s.clone();
    e.add(delta, unit);
    queries.push([s, e]);
    s = e;
  }
  return queries;
}

export type DatasetType = {
  t: number,
  y: number,
};

export type Dataset = {
  label: string,
  borderWidth: number,
  backgroundColor: string,
  borderColor: string,
  hoverBorderColor: string,
  hoverBackgroundColor: string,
  data: Array<DatasetType>,
};

type Props = {
  dataset: Array<Dataset>,
  unit?: string,
  yLabel?: string,
  // $FlowFixMe[value-as-type] Unresolved types
  tooltipHandler?: (ChartTooltipItem, ChartData) => string,
};

export default function CustomHistogram(props: Props) {
  return (
    <>
      <Bar
        height={300}
        data={{datasets: props.dataset}}
        options={{
          maintainAspectRatio: false,
          scaleShowValues: true,
          scales: {
            xAxes: [
              {
                stacked: true,
                gridLines: {
                  display: false,
                },
                type: 'time',
                ticks: {
                  source: 'data',
                },
                time: {
                  unit: props?.unit,
                  round: 'second',
                  tooltipFormat: 'YYYY/MM/DD h:mm:ss a',
                },
                scaleLabel: {
                  display: true,
                  labelString: 'Date',
                },
              },
            ],
            yAxes: [
              {
                stacked: true,
                gridLines: {
                  drawBorder: true,
                },
                ticks: {
                  maxTicksLimit: 3,
                },
                scaleLabel: {
                  display: true,
                  labelString: props?.yLabel ?? '',
                },
              },
            ],
          },
          tooltips: {
            enabled: true,
            mode: 'nearest',
            callbacks: {
              label: (tooltipItem, data) => {
                return (
                  props.tooltipHandler?.(tooltipItem, data) ??
                  data.datasets[tooltipItem.datasetIndex].label +
                    ': ' +
                    tooltipItem.yLabel +
                    ' ' +
                    (data.datasets[tooltipItem.datasetIndex].unit ?? '')
                );
              },
            },
          },
        }}
      />
    </>
  );
}

export function CustomLineChart(props: Props) {
  return (
    <>
      <Line
        height={300}
        data={{
          datasets: props.dataset,
        }}
        legend={{
          display: true,
          position: 'top',
          align: 'end',
          labels: {
            boxWidth: 12,
          },
        }}
        options={{
          maintainAspectRatio: false,
          scaleShowValues: true,
          scales: {
            xAxes: [
              {
                gridLines: {
                  display: false,
                },
                ticks: {
                  maxTicksLimit: 10,
                },
                type: 'time',
                time: {
                  unit: props?.unit,
                  round: 'second',
                  tooltipFormat: 'YYYY/MM/DD h:mm:ss a',
                },
                scaleLabel: {
                  display: true,
                  labelString: 'Date',
                },
              },
            ],
            yAxes: [
              {
                gridLines: {
                  drawBorder: true,
                },
                ticks: {
                  maxTicksLimit: 5,
                },
                scaleLabel: {
                  display: true,
                  labelString: props?.yLabel ?? '',
                },
                position: 'left',
              },
            ],
          },
          tooltips: {
            enabled: true,
            mode: 'nearest',
            callbacks: {
              label: (tooltipItem, data) => {
                return (
                  props.tooltipHandler?.(tooltipItem, data) ??
                  data.datasets[tooltipItem.datasetIndex].label +
                    ': ' +
                    tooltipItem.yLabel +
                    ' ' +
                    (data.datasets[tooltipItem.datasetIndex].unit ?? '')
                );
              },
            },
          },
        }}
      />
    </>
  );
}
