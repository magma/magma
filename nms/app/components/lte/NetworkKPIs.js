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

export default function () {
  const chartConfigs: MetricGraphConfig[] = [
    {
      label: 'Disk Percent',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(disk_percent)',
        },
      ],
      legendLabels: ['Disk Percent'],
      unit: '%',
    },
    {
      label: 'Number of Connected eNBs',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(enb_connected)',
        },
      ],
      legendLabels: ['Connected'],
      unit: '',
    },
    {
      label: 'Number of Connected UEs',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(ue_connected)',
        },
      ],
      legendLabels: ['Connected'],
      unit: '',
    },
    {
      label: 'Number of Registered UEs',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(ue_registered)',
        },
      ],
      legendLabels: ['Registered'],
      unit: '',
    },
    {
      label: 'S1 Setup',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(s1_setup)',
        },
        {
          resolveQuery: _ => "sum(s1_setup{result='success'})",
        },
        {
          resolveQuery: _ => "sum(s1_setup) - sum(s1_setup{result='success'})",
        },
      ],
      legendLabels: ['Total', 'Success', 'Failure'],
      unit: '',
    },
    {
      label: 'Attach/Reg Attempts',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(ue_attach)',
        },
        {
          resolveQuery: _ => "sum(ue_attach{result='attach_proc_successful'})",
        },
        {
          resolveQuery: _ =>
            "sum(ue_attach) - sum(ue_attach{result='attach_proc_successful'})",
        },
      ],
      legendLabels: ['Total', 'Success', 'Failure'],
      unit: '',
    },
    {
      label: 'Detach/Dereg Attempts',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(ue_detach)',
        },
        {
          resolveQuery: _ => "sum(ue_detach{result='attach_proc_successful'})",
        },
        {
          resolveQuery: _ =>
            "sum(ue_detach) - sum(ue_detach{result='attach_proc_successful'})",
        },
      ],
      legendLabels: ['Total', 'Success', 'Failure'],
      unit: '',
    },
    {
      label: 'GPS Connection Uptime',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'avg(enodeb_gps_connected)',
        },
      ],
      legendLabels: ['Uptime'],
      unit: '',
    },
    {
      label: 'Device Transmitting Status',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'avg(enodeb_rf_tx_enabled)',
        },
      ],
      legendLabels: ['Transmitting Status'],
      unit: '',
    },
    {
      label: 'Service Requests',
      basicQueryConfigs: [],
      customQueryConfigs: [
        {
          resolveQuery: _ => 'sum(service_request)',
        },
        {
          resolveQuery: _ => "sum(service_request{result='success'})",
        },
        {
          resolveQuery: _ =>
            "sum(service_request) - sum(service_request{result='success'})",
        },
      ],
      legendLabels: ['Total', 'Success', 'Failure'],
      unit: '',
    },
  ];

  return <NetworkMetrics configs={chartConfigs} />;
}
