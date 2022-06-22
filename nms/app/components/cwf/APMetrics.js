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

import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import SelectorMetrics from '../insights/SelectorMetrics';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {MetricGraphConfig} from '../insights/Metrics';

const APN_CONFIGS: Array<MetricGraphConfig> = [
  {
    label: 'Authorization',
    basicQueryConfigs: [],
    filters: [],
    customQueryConfigs: [
      {
        resolveQuery: apn => `sum(eap_auth{apn="${apn}"}) by (code)`,
      },
    ],
  },
  {
    label: 'Active Sessions',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: apn => `sum(active_sessions{apn="${apn}"})`,
      },
    ],
  },
  {
    label: 'Traffic In',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: apn =>
          `label_replace(sum(octets_in{apn="${apn}"}), "__name__", "octets_in", "__name__", ".*")`,
      },
    ],
  },
  {
    label: 'Traffic Out',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: apn => `sum(octets_out{apn="${apn}"})`,
      },
    ],
  },
  {
    label: 'Throughput In',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: apn => `avg(rate(octets_in{apn="${apn}"}[5m]))`,
      },
    ],
  },
  {
    label: 'Throughput Out',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: apn => `avg(rate(octets_out{apn="${apn}"}[5m]))`,
      },
    ],
  },
  {
    label: 'Accounting Stops',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: apn => `sum(accounting_stop{apn="${apn}"})`,
      },
    ],
  },
  {
    label: 'Session Terminate ',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: apn => `sum(session_terminate{apn="${apn}"})`,
      },
    ],
  },
];

export default function () {
  return <SelectorMetrics configs={APN_CONFIGS} selectorKey="apn" />;
}
