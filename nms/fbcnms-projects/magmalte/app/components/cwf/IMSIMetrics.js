/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MetricGraphConfig} from '../insights/Metrics';

import React from 'react';
import SelectorMetrics from '../insights/SelectorMetrics';

const IMSI_CONFIGS: Array<MetricGraphConfig> = [
  {
    label: 'Traffic In',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: imsi => `sum(octets_in{imsi="${imsi}"})`,
      },
    ],
  },
  {
    label: 'Traffic Out',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: imsi => `sum(octets_out{imsi="${imsi}"})`,
      },
    ],
  },
  {
    label: 'Throughput In',
    basicQueryConfigs: [],
    unit: '',
    customQueryConfigs: [
      {
        resolveQuery: imsi => `avg(rate(octets_in{imsi="${imsi}"}[5m]))`,
      },
    ],
  },
  {
    label: 'Throughput Out',
    basicQueryConfigs: [],
    unit: '',
    customQueryConfigs: [
      {
        resolveQuery: imsi => `avg(rate(octets_out{imsi="${imsi}"}[5m]))`,
      },
    ],
  },
  {
    label: 'Active Sessions',
    basicQueryConfigs: [],
    unit: '',
    customQueryConfigs: [
      {
        resolveQuery: imsi => `active_sessions{imsi="${imsi}"}`,
      },
    ],
  },
];

export default function() {
  return <SelectorMetrics configs={IMSI_CONFIGS} selectorKey="imsi" />;
}
