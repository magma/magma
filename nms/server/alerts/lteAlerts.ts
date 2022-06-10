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

export default function getLteAlerts(
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
    'Service Restart Alert': {
      alert: 'Service Restart Alert',
      expr: `increase(service_restart_status{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts upon service restarts'},
    },
    'Cpu Percent Alert': {
      alert: 'Cpu Percent Alert',
      expr: `avg_over_time(cpu_percent{networkID=~"${networkID}"}[5m]) > 70`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts when cpu percent is greater than 70%'},
    },
    'High Disk Usage Alert': {
      alert: 'High Disk Usage Alert',
      expr: `avg_over_time(disk_percent{networkID="${networkID}"}[5m]) > 70`,
      labels: {severity: 'major'},
      annotations: {
        description: 'Alerts when disk percent is greater than 70%',
      },
    },
    'High Memory Usage Alert': {
      alert: 'High Memory Usage Alert',
      expr: `((1 - mem_available{networkID="${networkID}"} / mem_total{networkID=~"${networkID}"}) * 100) > 70`,
      labels: {severity: 'major'},
      annotations: {
        description: 'Alerts when memory used is greater than 70%',
      },
    },
    'Unexpected Service Restart Alert': {
      alert: 'Unexpected Service Restart Alert',
      expr: `increase(unexpected_service_restarts{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'major'},
      annotations: {
        description: 'Alerts when we have unexpected service restarts',
      },
    },
    'Service Crashlooping Alert': {
      alert: 'Service Crashlooping Alert',
      expr: `increase(unexpected_service_restarts{networkID=~"${networkID}"}[5m]) > 3`,
      labels: {severity: 'critical'},
      annotations: {
        description: 'Alerts when we have services crashlooping',
      },
    },
    'Sctpd Crashlooping Alert': {
      alert: 'Sctpd Crashlooping Alert',
      expr: `increase(unexpected_service_restarts{networkID=~"${networkID}", service_name="sctpd"}[5m]) > 3`,
      labels: {severity: 'critical'},
      annotations: {
        description: 'Alerts when we have sctpd service crashlooping',
        remediation: 'Reboot the gateway',
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
      expr: `increase(s1_setup{result="failure", networkID=~"${networkID}"}[1h]) > 0`,
      labels: {severity: 'major'},
      annotations: {
        description: 'Alerts when we have S1 setup failures',
        remediation: 'Restart sctpd service',
      },
    },
    'UE attach Failure': {
      alert: 'UE attach Failure',
      expr: `increase(ue_attach{result="failure", networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts when we have UE attach failures'},
    },
    'High duplicate attach requests': {
      alert: 'High duplicate attach requests',
      expr: `increase(duplicate_attach_request{networkID="${networkID}"}[5m]) > 200`,
      labels: {severity: 'critical'},
      annotations: {
        description:
          'Alert when there are a large number of duplicate attach requests',
      },
    },
    'Gateway Checkin Failure': {
      alert: 'Gateway Checkin Failure',
      expr: `checkin_status{networkID=~"${networkID}"} < 1`,
      labels: {severity: 'critical'},
      annotations: {
        description: 'Alerts when we have gateway checkin failure',
        troubleshooting:
          'Run checkin_cli.py script on the gateway and follow resolution steps suggested',
      },
    },
    'Dip in Connected UEs': {
      alert: 'Dip in Connected UEs',
      expr: `(ue_connected{networkID=~"${networkID}"} - ue_connected{networkID=~"${networkID}"} offset 5m) / (ue_connected{networkID=~"${networkID}"}) < -0.5`,
      labels: {severity: 'critical'},
      annotations: {
        description:
          'Alerts when there is a 50% dip in connected UEs in last 5 minutes',
      },
    },
  };
}
