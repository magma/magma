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

import {apnTemplate, getNetworkTemplate} from './Dashboards';
import type {GrafanaDBData} from './Dashboards';

export const AnalyticsDBData = (networkIDs: Array<string>): GrafanaDBData => {
  return {
    title: 'Aggregated Analyses',
    description:
      'Various KPIs aggregated and calculated from other existing metrics.',
    templates: [getNetworkTemplate(networkIDs), apnTemplate],
    rows: [
      {
        title: 'Active Users',
        panels: [
          {
            title: 'Active Users Over Time',
            targets: [
              {
                expr: 'active_users_over_time{networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - Days: {{days}}',
              },
            ],
            description: 'Number of unique active users over the past N days.',
          },
        ],
      },
      {
        title: 'User Throughput',
        panels: [
          {
            title: 'User Throughput (Download)',
            targets: [
              {
                expr: 'user_throughput{direction="in",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - Days: {{days}}',
              },
            ],
            unit: 'Bps',
            description:
              'Average user download throughput over the network in the last N days.',
          },
          {
            title: 'User Throughput (Upload)',
            targets: [
              {
                expr:
                  'user_throughput{direction="out",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - Days: {{days}}',
              },
            ],
            unit: 'Bps',
            description:
              'Average user upload throughput over the network in the last N days.',
          },
          {
            title: 'Throughput Per APN (Upload)',
            targets: [
              {
                expr: 'throughput_per_ap{direction="out",apn=~"$apn"}',
                legendFormat: '{{apn}}',
              },
            ],
            unit: 'Bps',
            description:
              'Average user upload throughput for a given APN in the last N days',
          },
          {
            title: 'Throughput Per APN (Download)',
            targets: [
              {
                expr: 'throughput_per_ap{direction="in",apn=~"$apn"}',
                legendFormat: '{{apn}}',
              },
            ],
            unit: 'Bps',
            description:
              'Average user download throughput for a given APN in the last N days',
          },
        ],
      },
      {
        title: 'User Consumption',
        panels: [
          {
            title: 'User Consumption (Upload)',
            targets: [
              {
                expr:
                  'user_consumption{direction="out",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - Days: {{days}}',
              },
            ],
            unit: 'decbytes',
            description:
              'Total user upload consumption over the network in the past N days.',
          },
          {
            title: 'User Consumption (Download)',
            targets: [
              {
                expr:
                  'user_consumption{direction="in",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - Days: {{days}}',
              },
            ],
            unit: 'decbytes',
            description:
              'Total user download consumption over the network in the past N days.',
          },
          {
            title: 'User Consumption Hourly (Upload)',
            targets: [
              {
                expr:
                  'user_consumption_hourly{direction="out",networkID=~"$networkID"}',
                legendFormat: '{{networkID}}',
              },
            ],
            unit: 'decbytes',
            description:
              'Total user upload consumption over the network in the past hour.',
          },
          {
            title: 'User Consumption Hourly (Download)',
            targets: [
              {
                expr:
                  'user_consumption_hourly{direction="in",networkID=~"$networkID"}',
                legendFormat: '{{networkID}}',
              },
            ],
            unit: 'decbytes',
            description:
              'Total user download consumption over the network in the past hour.',
          },
        ],
      },
      {
        title: 'Authentications',
        panels: [
          {
            title: 'Authentications Over Time',
            targets: [
              {
                expr: 'authentications_over_time{networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - Code: {{code}}',
              },
            ],
            description:
              'Total number of authentication failures or successes in the network in the last N days.',
          },
        ],
      },
    ],
  };
};
