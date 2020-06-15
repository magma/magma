/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import moment from 'moment';

import {Bar} from 'react-chartjs-2';

export function getStep(start: moment, end: moment): [number, string, string] {
  const d = moment.duration(end.diff(start));
  if (d.asMinutes() <= 60.5) {
    return [5, 'minutes', 'HH::mm'];
  } else if (d.asHours() <= 3.5) {
    return [15, 'minutes', 'HH::mm'];
  } else if (d.asHours() <= 6.5) {
    return [15, 'minutes', 'HH::mm'];
  } else if (d.asHours() <= 12.5) {
    return [1, 'hours', 'HH::mm'];
  } else if (d.asHours() <= 24.5) {
    return [2, 'hours', 'HH::mm'];
  } else if (d.asDays() <= 3.5) {
    return [6, 'hours', 'DD-MM-YY HH::mm'];
  } else if (d.asDays() <= 7.5) {
    return [12, 'hours', 'DD-MM-YY HH::mm'];
  } else if (d.asDays() <= 14.5) {
    return [1, 'days', 'DD-MM-YYYY'];
  } else if (d.asDays() <= 30.5) {
    return [1, 'days', 'DD-MM-YYYY'];
  } else if (d.asMonths() <= 3.5) {
    return [7, 'days', 'DD-MM-YYYY'];
  }
  return [1, 'months', 'DD-MM-YYYY'];
}

export type Dataset = {
  label: string,
  borderWidth: number,
  backgroundColor: string,
  borderColor: string,
  hoverBorderColor: string,
  hoverBackgroundColor: string,
  data: Array<number>,
};

type HistogramProps = {
  labels: Array<string>,
  dataset: Array<Dataset>,
};

export default function CustomHistogram(props: HistogramProps) {
  return (
    <>
      <Bar
        data={{labels: props.labels, datasets: props.dataset}}
        options={{
          maintainAspectRatio: false,
          scaleShowValues: true,
          scales: {
            xAxes: [
              {
                gridLines: {
                  display: false,
                },
              },
            ],
            yAxes: [
              {
                gridLines: {
                  drawBorder: true,
                },
                ticks: {
                  maxTicksLimit: 3,
                },
              },
            ],
          },
        }}
      />
    </>
  );
}
