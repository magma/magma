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

import type {prom_alert_config} from '@fbcnms/magma-api';

export default function getLteAlerts(
  networkID: string,
): {[string]: prom_alert_config} {
  return {
    'Service Restart Alert': {
      alert: 'Service Restart Alert',
      expr: `avg_over_time(service_restart_status{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts upon service restarts'},
    },
    'Cpu Percent Alert': {
      alert: 'Cpu Percent Alert',
      expr: `avg_over_time(cpu_percent{networkID=~"${networkID}"}[5m]) > 75`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts when cpu percent is greater than 75%'},
    },
    'Unexpected Service Restart Alert': {
      alert: 'Unexpected Service Restart Alert',
      expr: `increase(unexpected_service_restarts{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'major'},
      annotations: {
        description: 'Alerts when we have unexpected service restarts',
      },
    },
    'Bootstrap Exception Alert': {
      alert: 'Bootstrap Exception Alert',
      expr: `increase(bootstrap_exception{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description: 'Alerts when we have exceptions during bootstrapping',
      },
    },
    'S6A Auth Failure': {
      alert: 'S6A Auth Failure',
      expr: `increase(s6a_auth_failure{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts when we have auth failures'},
    },
    'S1 Setup Failure': {
      alert: 'S1 Setup Failure',
      expr: `increase(s1_setup{result="failure", networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'major'},
      annotations: {description: 'Alerts when we have S1 setup failures'},
    },
    'UE attach Failure': {
      alert: 'UE attach Failure',
      expr: `increase(ue_attach{result="failure", networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts when we have UE attach failures'},
    },
  };
}
