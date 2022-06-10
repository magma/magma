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

export default function getFegAlerts(
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
    'Gateway Cert Expiring Alert': {
      alert: 'Gateway Certificate is expiring',
      expr: `increase(gateway_expiring_cert{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'major'},
      annotations: {
        description: 'Alerts when we have expiring gateway certificate',
      },
    },
    'OCS CCR_INIT Send Failed': {
      alert: 'OCS CCR_INIT Send Failed',
      expr: `increase(ocs_ccr_init_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description: 'Alerts when we have ocs ccs init send failures',
      },
    },
    'OCS CCR_UPDATE Send Failed': {
      alert: 'OCS CCR_UPDATE Send Failed',
      expr: `increase(ocs_ccr_update_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description: 'Alerts when we have ocs ccs update send failures',
      },
    },
    'OCS CCR_TERMINATE Send Failed': {
      alert: 'OCS CCR_TERMINATE Send Failed',
      expr: `increase(ocs_ccr_terminate_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description: 'Alerts when we have ocs ccs terminate send failures',
      },
    },
    'PCRF CCR_INIT Send Failed': {
      alert: 'PCRF CCR_INIT Send Failed',
      expr: `increase(pcrf_ccr_init_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description: 'Alerts when we have pcrf ccs init send failures',
      },
    },
    'PCRF CCR_UPDATE Send Failed': {
      alert: 'PCRF CCR_UPDATE Send Failed',
      expr: `increase(pcrf_ccr_update_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description: 'Alerts when we have pcrf ccs update send failures',
      },
    },
    'PCRF CCR_TERMINATE Send Failed': {
      alert: 'PCRF CCR_TERMINATE Send Failed',
      expr: `increase(pcrf_ccr_terminate_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description: 'Alerts when we have pcrf ccs terminate send failures',
      },
    },
    'AIR requests that failed to send to HSS': {
      alert: 'AIR requests that failed to send to HSS',
      expr: `increase(air_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description:
          'Alerts when we have air requests that failed to send to HSS',
      },
    },
    'MAR requests that failed to send to HSS': {
      alert: 'MAR requests that failed to send to HSS',
      expr: `increase(mar_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description:
          'Alerts when we have MAR requests that failed to send to HSS',
      },
    },
    'SAR requests that failed to send to HSS': {
      alert: 'SAR requests that failed to send to HSS',
      expr: `increase(sar_send_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {
        description:
          'Alerts when we have SAR requests that failed to send to HSS',
      },
    },
    'SWx Proxy RPC Failures': {
      alert: 'SWx Proxy RPC Failures',
      expr: `increase(swx_failures_total{networkID=~"${networkID}"}[5m]) > 0`,
      labels: {severity: 'minor'},
      annotations: {description: 'Alerts when we have SWx Proxy RPC Failures'},
    },
    'EAP SIM SWx Proxy RPC Failures': {
      alert: 'SEAP SIM SWx Proxy RPC Failures',
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
