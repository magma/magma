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

import type {PromAlertConfig} from '../../generated-ts';

export default function getCwfAlerts(
  networkID: string,
): {[name: string]: PromAlertConfig} {
  return {
    'Certificate Expiring Soon': {
      alert: 'Certificate Expiring Soon',
      expr: `cert_expires_in_hours < 720`,
      labels: {severity: 'major'},
      annotations: {
        description: `Alerts when certificate necessary for Orc8r function is expiring soon`,
      },
    },
    'Container Restarts': {
      alert: 'Container Restarts',
      expr: `time() - container_start_time_seconds{name=~".+",name!="cadvisor", networkID=~"${networkID}"} < 300`,
      labels: {severity: 'major'},
      for: '6m',
      annotations: {
        description: `container is restarting often`,
      },
    },
    'Container Cpu Usage': {
      alert: 'Container CPU Usage High',
      expr: `sum(rate(container_cpu_usage_seconds_total{name=~".+", networkID=~"${networkID}"}[5m])) by (name,networkID,gatewayID) > 0.9`,
      labels: {severity: 'major'},
      annotations: {
        description: 'A container has had very high CPU usage for 5 minutes',
      },
    },
    'SWx Proxy RPC Failures': {
      alert: 'SWx Proxy RPC Failures',
      expr: `increase(swx_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts when we have SWx Proxy RPC Failures'},
    },
    'EAP SIM SWx Proxy RPC Failures': {
      alert: 'EAP SIM SWx Proxy RPC Failures',
      expr: `increase(eap_sim_swx_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description: 'Alerts when we have EAP SIM SWx Proxy RPC Failures',
      },
    },
    'EAP SIM Errors/Failures originated from peers': {
      alert: 'EAP SIM Errors/Failures originated from peers',
      expr: `increase(eap_sim_peer_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description:
          'Alerts when we have SIM Errors/Failures originated from peers',
      },
    },
  };
}
