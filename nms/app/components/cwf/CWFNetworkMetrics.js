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

import NetworkMetrics from '../insights/NetworkMetrics';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {MetricGraphConfig} from '../insights/Metrics';

const chartConfigs: MetricGraphConfig[] = [
  {
    label: 'Authorization',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `sum(eap_auth) by (code)`,
      },
    ],
  },
  {
    label: 'Active Sessions',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `sum(active_sessions)`,
      },
    ],
  },
  {
    label: 'Traffic In',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `sum(octets_in)`,
      },
    ],
  },
  {
    label: 'Traffic Out',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `sum(octets_out)`,
      },
    ],
  },
  {
    label: 'Throughput In',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `avg(rate(octets_in[5m]))`,
      },
    ],
  },
  {
    label: 'Throughput Out',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `avg(rate(octets_out[5m]))`,
      },
    ],
  },
  {
    label: 'Accounting Stops',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `sum(accounting_stop)`,
      },
    ],
  },
  {
    label: 'Session Create Latency',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `avg(create_session_lat)`,
      },
    ],
  },
  {
    label: 'Session Stop',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: _ => `sum(session_stop)`,
      },
    ],
  },
];

export default function () {
  return <NetworkMetrics configs={chartConfigs} />;
}
