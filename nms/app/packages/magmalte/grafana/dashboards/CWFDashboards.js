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

import {gatewayTemplate, networkTemplate, variableTemplate} from './Dashboards';
import type {GrafanaDBData} from './Dashboards';

const msisdnTemplate = variableTemplate({
  labelName: 'msisdn',
  query: `label_values(msisdn)`,
  regex: `/.+/`,
  sort: 'num-asc',
  includeAll: false,
});

const apnTemplate = variableTemplate({
  labelName: 'apn',
  query: `label_values({networkID=~"$networkID",apn=~".+"},apn)`,
  regex: `/.+/`,
  sort: 'alpha-insensitive-asc',
  includeAll: true,
});

const dbDescription =
  'Do not edit: edits will be overwritten. Save this dashboard under another name to copy and edit.';

export const CWFSubscriberDBData: GrafanaDBData = {
  title: 'CWF - Subscribers',
  description: dbDescription,
  templates: [msisdnTemplate],
  rows: [
    {
      title: 'Traffic',
      panels: [
        {
          title: 'Traffic In',
          targets: [
            {
              expr: 'sum(octets_in{msisdn=~"$msisdn"}) by (imsi, msisdn)',
              legendFormat: '{{imsi}}, MSISDN: {{msisdn}}',
            },
          ],
          unit: 'decbytes',
          description: 'Inbound data per subscriber measured in bytes.',
        },
        {
          title: 'Throughput In',
          targets: [
            {
              expr:
                'avg(rate(octets_in{msisdn=~"$msisdn"}[5m])) by (imsi, msisdn)',
              legendFormat: '{{imsi}}, MSISDN: {{msisdn}}',
            },
          ],
          unit: 'Bps',
          description:
            'Inbound data rate per subscriber measured in bytes/second.',
        },
        {
          title: 'Traffic Out',
          targets: [
            {
              expr: 'sum(octets_out{msisdn=~"$msisdn"}) by (imsi, msisdn)',
              legendFormat: '{{imsi}}, MSISDN: {{msisdn}}',
            },
          ],
          unit: 'decbytes',
          description: 'Outbound data measured in bytes.',
        },
        {
          title: 'Throughput Out',
          targets: [
            {
              expr:
                'avg(rate(octets_out{msisdn=~"$msisdn"}[5m])) by (imsi, msisdn)',
              legendFormat: '{{imsi}}, MSISDN: {{msisdn}}',
            },
          ],
          unit: 'Bps',
          description:
            'Outbound data rate per subscriber measured in bytes/second.',
        },
      ],
    },
    {
      title: 'Session',
      panels: [
        {
          title: 'Active Sessions',
          targets: [
            {
              expr: 'active_sessions{msisdn=~"$msisdn"}',
              legendFormat:
                '{{imsi}} -- MSISDN: {{msisdn}} -- Session: {{id}} -- Network: {{networkID}} -- Gateway: {{gatewayID}}',
            },
          ],
          description:
            "Should just be 0 for subscribers without an active session and 1 for those with. Not sure what's up with this now.",
        },
      ],
    },
  ],
};

export const CWFAccessPointDBData: GrafanaDBData = {
  title: 'CWF - Access Points',
  description: dbDescription,
  templates: [networkTemplate, apnTemplate],
  rows: [
    {
      title: 'Message Stats',
      panels: [
        {
          title: 'Accounting Stops',
          targets: [
            {
              expr: 'sum(accounting_stop{apn=~"$apn"}) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          description: 'Radius accounting stops received from AP/WLC',
        },
        {
          title: 'Authorization',
          targets: [
            {
              expr: 'sum(eap_auth{apn=~"$apn"}) by (code, apn)',
              legendFormat: '{{apn}}-{{code}}',
            },
          ],
          description:
            'EAP Authorization responses, partitioned by response type (Failure, Success) where request is the sum of success and failures',
        },
      ],
    },
    {
      title: 'Traffic',
      panels: [
        {
          title: 'Traffic In',
          targets: [
            {
              expr: 'sum(octets_in{apn=~"$apn"}) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          unit: 'decbytes',
          description: 'Inbound data measured in bytes.',
        },
        {
          title: 'Traffic Out',
          targets: [
            {
              expr: 'sum(octets_out{apn=~"$apn"}) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          unit: 'decbytes',
          description: 'Outbound data measured in bytes.',
        },
        {
          title: 'Throughput In',
          targets: [
            {
              expr: 'avg(rate(octets_in{apn=~"$apn"}[5m])) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          unit: 'Bps',
          description: 'Inbound data rate measured in bytes/second.',
        },
        {
          title: 'Throughput Out',
          targets: [
            {
              expr: 'avg(rate(octets_out{apn=~"$apn"}[5m])) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          unit: 'Bps',
          description: 'Outbound data rate measured in bytes/second.',
        },
      ],
    },
    {
      title: 'Session',
      panels: [
        {
          title: 'Active Sessions',
          targets: [
            {
              expr: 'sum(active_sessions{apn=~"$apn"}) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          description: 'Number of active user sessions in the network',
        },
        {
          title: 'Session Stop',
          targets: [
            {
              expr: 'sum(session_stop{apn=~"$apn"}) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          description: 'Number of sessions removed for any reason',
        },
        {
          title: 'Session Timeout',
          targets: [
            {
              expr: 'sum(session_timeouts{apn=~"$apn"}) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          description:
            'Subset of session_stop. Count of any session that times out from aaa server',
        },
        {
          title: 'Session Terminate',
          targets: [
            {
              expr: 'sum(session_manager_terminate{apn=~"$apn"}) by (apn)',
              legendFormat: '{{apn}}',
            },
          ],
          unit: 's',
          description:
            'When sessiond terminates a session to aaa server. Reasons: Mostly due to out of quota, TODO: Fill in later?',
        },
      ],
    },
  ],
};

export const CWFNetworkDBData: GrafanaDBData = {
  title: 'CWF - Networks',
  description: dbDescription,
  templates: [networkTemplate],
  rows: [
    {
      title: 'Message Stats',
      panels: [
        {
          title: 'Authorization',
          targets: [
            {
              expr:
                'sum(eap_auth{networkID=~"$networkID"}) by (code, networkID)',
              legendFormat: '{{networkID}}-{{code}}',
            },
          ],
          description:
            'EAP Authorization responses, partitioned by response type (Failure, Success) where request is the sum of success and failures',
        },
        {
          title: 'Accounting Stops',
          targets: [
            {
              expr:
                'sum(accounting_stop{networkID=~"$networkID"}) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Radius accounting stops received from AP/WLC',
        },
      ],
    },
    {
      title: 'Traffic',
      panels: [
        {
          title: 'Traffic In',
          targets: [
            {
              expr: 'sum(octets_in{networkID=~"$networkID"}) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          unit: 'decbytes',
          description: 'Inbound data measured in bytes.',
        },
        {
          title: 'Traffic Out',
          targets: [
            {
              expr: 'sum(octets_out{networkID=~"$networkID"}) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          unit: 'decbytes',
          description: 'Outbound data measured in bytes.',
        },
        {
          title: 'Throughput In',
          targets: [
            {
              expr:
                'avg(rate(octets_in{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          unit: 'Bps',
          description: 'Inbound data rate measured in bytes/second.',
        },
        {
          title: 'Throughput Out',
          targets: [
            {
              expr:
                'avg(rate(octets_out{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          unit: 'Bps',
          description: 'Outbound data rate measured in bytes/second.',
        },
      ],
    },
    {
      title: 'Latency',
      panels: [
        {
          title: 'Session Create Latency',
          targets: [
            {
              expr:
                'avg(create_session_lat{networkID=~"$networkID"}) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          unit: 's',
          description:
            'Average time taken to create a session over the network.',
        },
      ],
    },
    {
      title: 'Session',
      panels: [
        {
          title: 'Active Sessions',
          targets: [
            {
              expr:
                'sum(active_sessions{networkID=~"$networkID"}) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Number of active user sessions in the network',
        },
        {
          title: 'Session Stop',
          targets: [
            {
              expr: 'sum(session_stop{networkID=~"$networkID"}) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Number of sessions removed for any reason',
        },
        {
          title: 'Session Timeouts',
          targets: [
            {
              expr:
                'sum(session_timeouts{networkID=~"$networkID"}) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Subset of session_stop. Count of any session that times out from aaa server',
        },
        {
          title: 'Session Terminate',
          targets: [
            {
              expr:
                'sum(session_manager_terminate{networkID=~"$networkID"}) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'When sessiond terminates a session to aaa server. Reasons: Mostly due to out of quota, TODO: Fill in later?',
        },
      ],
    },
    {
      title: 'Diameter Result Codes',
      panels: [
        {
          title: 'Gx Result Codes',
          targets: [
            {
              expr:
                'sum(rate(gx_result_codes{networkID=~"$networkID"}[5m])) by (networkID, code)',
              legendFormat: '{{networkID}} - {{code}}',
            },
          ],
          description: 'Rate of Gx responses segmented by code',
        },
        {
          title: 'Gy Result Codes',
          targets: [
            {
              expr:
                'sum(rate(gy_result_codes{networkID=~"$networkID"}[5m])) by (networkID, code)',
              legendFormat: '{{networkID}} - {{code}}',
            },
          ],
          description: 'Rate of Gy responses segmented by code',
        },
        {
          title: 'SWX Result Codes',
          targets: [
            {
              expr:
                'sum(rate(swx_result_codes{networkID=~"$networkID"}[5m])) by (networkID, code)',
              legendFormat: '{{networkID}} - {{code}}',
            },
          ],
          description: 'Rate of SWx responses segmented by diameter base code',
        },
        {
          title: 'SWX Experimental Result Codes',
          targets: [
            {
              expr:
                'sum(rate(swx_experimental_result_codes{networkID=~"$networkID"}[5m])) by (networkID, code)',
              legendFormat: '{{networkID}} - {{code}}',
            },
          ],
          description: 'Rate of SWx responses segmented by SWx-specific code',
        },
      ],
    },
    {
      title: 'Diameter Timeouts',
      panels: [
        {
          title: 'Gx Timeouts',
          targets: [
            {
              expr:
                'sum(rate(gx_timeouts_total{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of Gx requests that did not receive a response (and thus timed out)',
        },
        {
          title: 'Gy Timeouts',
          targets: [
            {
              expr:
                'sum(rate(gy_timeouts_total{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of Gy requests that did not receive a response (and thus timed out)',
        },
        {
          title: 'SWX Timeouts',
          targets: [
            {
              expr:
                'sum(rate(swx_timeouts_total{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of SWx requests that did not receive a response (and thus timed out)',
        },
      ],
    },
    {
      title: 'OCS CCR Requests',
      panels: [
        {
          title: 'Initializations',
          targets: [
            {
              expr: 'sum(rate(ocs_ccr_init_requests_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of Gy CCR-I requests',
        },
        {
          title: 'Terminations',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_terminate_requests_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of Gy CCR-T requests',
        },
        {
          title: 'Updates',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_update_requests_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of Gy CCR-U requests',
        },
      ],
    },
    {
      title: 'OCS Send Failures',
      panels: [
        {
          title: 'Initialization Failures',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_init_send_failures_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of Gy CCR-I messages that were unable to be sent due to diameter connection errors',
        },
        {
          title: 'Temination Failures',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_terminate_send_failures_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of Gy CCR-T messages that were unable to be sent due to diameter connection errors',
        },
        {
          title: 'Update Failures',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_update_send_failures_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of Gy CCR-U messages that were unable to be sent due to diameter connection errors',
        },
      ],
    },
    {
      title: 'PCRF CCR Requests',
      panels: [
        {
          title: 'Initializations',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_init_requests_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of Gx CCR-I requests',
        },
        {
          title: 'Teminations',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_terminate_requests_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of Gx CCR-T requests',
        },
        {
          title: 'Updates',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_update_requests_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of Gx CCR-U requests',
        },
      ],
    },
    {
      title: 'PCRF CCR Send Failures',
      panels: [
        {
          title: 'Initialization Failures',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_init_send_failures_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of Gx CCR-I messages that were unable to be sent due to diameter connection errors',
        },
        {
          title: 'Temination Failures',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_terminate_send_failures_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of Gx CCR-T messages that were unable to be sent due to diameter connection errors',
        },
        {
          title: 'Update Failures',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_update_send_failures_total[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description:
            'Rate of Gx CCR-U messages that were unable to be sent due to diameter connection errors',
        },
      ],
    },
    {
      title: 'HSS Requests/Failures',
      panels: [
        {
          title: 'MAR Requests',
          targets: [
            {
              expr:
                'sum(rate(mar_requests_total{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of SWx MAR requests',
        },
        {
          title: 'SAR Requests',
          targets: [
            {
              expr:
                'sum(rate(sar_requests_total{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of SWx SAR requests',
        },
        {
          title: 'MAR Failures',
          targets: [
            {
              expr:
                'sum(rate(mar_send_failures_total{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of SWx MAR request failures',
        },
        {
          title: 'SAR Failures',
          targets: [
            {
              expr:
                'sum(rate(sar_send_failures_total{networkID=~"$networkID"}[5m])) by (networkID)',
              legendFormat: '{{networkID}}',
            },
          ],
          description: 'Rate of SWx SAR request failures',
        },
      ],
    },
  ],
};

export const CWFGatewayDBData: GrafanaDBData = {
  title: 'CWF - Gateways',
  description: dbDescription,
  templates: [networkTemplate, gatewayTemplate],
  rows: [
    {
      title: 'Diameter Result Codes',
      panels: [
        {
          title: 'Gx Result Codes',
          targets: [
            {
              expr:
                'sum(rate(gx_result_codes{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID, code)',
              legendFormat: '{{gatewayID}} - {{code}}',
            },
          ],
          description: 'Rate of Gx responses segmented by code',
        },
        {
          title: 'Gy Result Codes',
          targets: [
            {
              expr:
                'sum(rate(gy_result_codes{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID, code)',
              legendFormat: '{{gatewayID}} - {{code}}',
            },
          ],
          description: 'Rate of Gy responses segmented by code',
        },
        {
          title: 'SWX Result Codes',
          targets: [
            {
              expr:
                'sum(rate(swx_result_codes{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID, code)',
              legendFormat: '{{gatewayID}} - {{code}}',
            },
          ],
          description: 'Rate of SWx responses segmented by diameter base code',
        },
        {
          title: 'SWX Experimental Result Codes',
          targets: [
            {
              expr:
                'sum(rate(swx_experimental_result_codes{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID, code)',
              legendFormat: '{{gatewayID}} - {{code}}',
            },
          ],
          description: 'Rate of SWx responses segmented by SWx-specific code',
        },
      ],
    },
    {
      title: 'Diameter Timeouts',
      panels: [
        {
          title: 'Gx Timeouts',
          targets: [
            {
              expr:
                'sum(rate(gx_timeouts_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description:
            'Rate of Gx requests that did not receive a response (and thus timed out)',
        },
        {
          title: 'Gy Timeouts',
          targets: [
            {
              expr:
                'sum(rate(gy_timeouts_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description:
            'Rate of Gy requests that did not receive a response (and thus timed out)',
        },
        {
          title: 'SWX Timeouts',
          targets: [
            {
              expr:
                'sum(rate(swx_timeouts_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description:
            'Rate of SWx requests that did not receive a response (and thus timed out)',
        },
      ],
    },
    {
      title: 'OCS CCR Requests',
      panels: [
        {
          title: 'Initializations',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_init_requests_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of Gy CCR-I requests',
        },
        {
          title: 'Terminations',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_terminate_requests_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of Gy CCR-T requests',
        },
        {
          title: 'Updates',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_update_requests_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of Gy CCR-U requests',
        },
      ],
    },
    {
      title: 'OCS Send Failures',
      panels: [
        {
          title: 'Initialization Failures',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_init_send_failures_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description:
            'Rate of Gy CCR-I messages that were unable to be sent due to diameter connection errors',
        },
        {
          title: 'Termination Failures',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_terminate_send_failures_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description:
            'Rate of Gy CCR-T messages that were unable to be sent due to diameter connection errors',
        },
        {
          title: 'Update Failures',
          targets: [
            {
              expr:
                'sum(rate(ocs_ccr_update_send_failures_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description:
            'Rate of Gy CCR-U messages that were unable to be sent due to diameter connection errors',
        },
      ],
    },
    {
      title: 'PCRF CCR Requests',
      panels: [
        {
          title: 'Initializations',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_init_requests_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of Gx CCR-I requests',
        },
        {
          title: 'Terminations',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_terminate_requests_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of Gx CCR-T requests',
        },
        {
          title: 'Updates',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_update_requests_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of Gx CCR-U requests',
        },
      ],
    },
    {
      title: 'PCRF CCR Send Failures',
      panels: [
        {
          title: 'Initialization Failures',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_init_send_failures_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description:
            'Rate of Gx CCR-I messages that were unable to be sent due to diameter connection errors',
        },
        {
          title: 'Termination Failures',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_terminate_send_failures_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description:
            'Rate of Gx CCR-T messages that were unable to be sent due to diameter connection errors',
        },
        {
          title: 'Update Failures',
          targets: [
            {
              expr:
                'sum(rate(pcrf_ccr_update_send_failures_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
        },
      ],
      description:
        'Rate of Gx CCR-U messages that were unable to be sent due to diameter connection errors',
    },
    {
      title: 'HSS Requests/Failures',
      panels: [
        {
          title: 'MAR Requests',
          targets: [
            {
              expr:
                'sum(rate(mar_requests_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of SWx MAR requests',
        },
        {
          title: 'SAR Requests',
          targets: [
            {
              expr:
                'sum(rate(sar_requests_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of SWx SAR requests',
        },
        {
          title: 'MAR Failures',
          targets: [
            {
              expr:
                'sum(rate(mar_send_failures_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of SWx MAR request failures',
        },
        {
          title: 'SAR Failures',
          targets: [
            {
              expr:
                'sum(rate(sar_send_failures_total{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID)',
              legendFormat: '{{gatewayID}}',
            },
          ],
          description: 'Rate of SWx SAR request failures',
        },
      ],
    },
  ],
};
