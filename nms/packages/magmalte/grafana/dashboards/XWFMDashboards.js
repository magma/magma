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

import {gatewayTemplate, getNetworkTemplate} from './Dashboards';
import type {GrafanaDBData} from './Dashboards';

export const XWFMDBData = (networks: Array<string>): GrafanaDBData => {
  return {
    title: 'Container/Node Stats',
    description:
      'This dashboard shows stats from node_exporter and cAdvisor if those services are installed on your gateways.',
    templates: [getNetworkTemplate(networks), gatewayTemplate],
    rows: [
      {
        title: 'Containers',
        panels: [
          {
            title: 'Container Uptime',
            targets: [
              {
                expr:
                  'time() - container_start_time_seconds{gatewayID=~"$gatewayID",networkID=~"$networkID",name=~".+",name!="cadvisor"}',
                legendFormat: '{{networkID}}-{{gatewayID}} - {{name}}',
              },
            ],
            unit: 's',
            description: 'Time since last container restart',
          },
          {
            title: 'Container CPU Usage',
            targets: [
              {
                expr:
                  'sum(rate(container_cpu_usage_seconds_total{networkID=~"$networkID",gatewayID=~"$gatewayID",name=~".+"}[1m])) by (name,networkID,gatewayID)',
                legendFormat: '{{networkID}}-{{gatewayID}} - {{name}}',
              },
            ],
            description: 'CPU usage per container (averaged over 1 minute)',
          },
          {
            title: 'Container Memory Usage',
            targets: [
              {
                expr:
                  'container_memory_usage_bytes{networkID=~"$networkID",gatewayID=~"$gatewayID",name=~".+"}',
                legendFormat: '{{networkID}}-{{name}}',
              },
            ],
            unit: 'decbytes',
            description: 'Memory Usage (bytes) per container',
          },
          {
            title: '7-Day Rolling Availability',
            targets: [
              {
                expr:
                  'sum_over_time(min by(networkID,gatewayID)(time() - container_last_seen{gatewayID=~"$gatewayID",networkID=~"$networkID",name=~".+",name!="cadvisor"} <= bool 30)[7d:1m]) / 10080',
                legendFormat: '{{networkID}}-{{gatewayID}}',
              },
            ],
          },
        ],
      },
      {
        title: 'Node',
        panels: [
          {
            title: 'Packet Transfer Up',
            targets: [
              {
                expr:
                  'sum(rate(node_network_transmit_packets_total{networkID=~"$networkID",gatewayID=~"$gatewayID"}[1m])) by (networkID,gatewayID)',
                legendFormat: '{{networkID}}-{{gatewayID}} - UP',
              },
              {
                expr:
                  'sum(rate(node_network_receive_packets_total{networkID=~"$networkID",gatewayID=~"$gatewayID"}[1m])) by (networkID,gatewayID)',
                legendFormat: '{{networkID}}-{{gatewayID}} - DOWN',
              },
            ],
            description: 'Rate of packet transfer per gateway',
          },
          {
            title: 'Node Bytes Received',
            targets: [
              {
                expr:
                  'sum(rate(node_network_receive_bytes_total{networkID=~"$networkID",gatewayID=~"$gatewayID"}[5m])) by (networkID,gatewayID)',
                legendFormat: '{{networkID}}-{{gatewayID}} - SENT',
              },
              {
                expr:
                  'sum(rate(node_network_transmit_bytes_total{networkID=~"$networkID",gatewayID=~"$gatewayID"}[5m])) by (networkID,gatewayID)',
                legendFormat: '{{networkID}}-{{gatewayID}} - RECEIVED',
              },
            ],
            unit: 'Bps',
            description: 'Rate of bytes sent/received per gateway',
          },
          {
            title: 'Memory Usage',
            targets: [
              {
                expr:
                  '(1 - (node_memory_Active_bytes{networkID=~"$networkID",gatewayID=~"$gatewayID"}) / node_memory_MemTotal_bytes{networkID=~"$networkID",gatewayID=~"$gatewayID"}) * 100',
                legendFormat: '{{networkID}}-{{gatewayID}}',
              },
            ],
            unit: 'percent',
            description: 'Memory usage of the node per gateway',
            yMax: 100,
            yMin: 0,
          },
          {
            title: 'CPU Idle Percent',
            targets: [
              {
                expr:
                  'avg(rate(node_cpu_seconds_total{networkID=~"$networkID",mode="idle",gatewayID=~"$gatewayID"}[1m])) by (networkID,gatewayID) * 100',
                legendFormat: '{{networkID}}-{{gatewayID}}',
              },
            ],
            description: 'Percent of time spent idle by CPU per gateway',
            yMax: 100,
            yMin: 0,
          },
        ],
      },
    ],
  };
};
