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

import {
  apnTemplate,
  gatewayTemplate,
  getNetworkTemplate,
  msisdnTemplate,
} from './Dashboards';
import type {GrafanaDBData} from './Dashboards';

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
              expr:
                'sum(ue_reported_usage{msisdn=~"$msisdn", direction="down"}) by (IMSI, msisdn)',
              legendFormat: '{{IMSI}}, MSISDN: {{msisdn}}',
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
                'avg(rate(ue_reported_usage{msisdn=~"$msisdn", direction="down"}[5m])) by (IMSI, msisdn)',
              legendFormat: '{{IMSI}}, MSISDN: {{msisdn}}',
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
              expr:
                'sum(ue_reported_usage{msisdn=~"$msisdn", direction="up"}) by (IMSI, msisdn)',
              legendFormat: '{{IMSI}}, MSISDN: {{msisdn}}',
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
                'avg(rate(ue_reported_usage{msisdn=~"$msisdn", direction="up"}[5m])) by (IMSI, msisdn)',
              legendFormat: '{{IMSI}}, MSISDN: {{msisdn}}',
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

export const CWFAccessPointDBData = (
  networkIDs: Array<string>,
): GrafanaDBData => {
  return {
    title: 'CWF - Access Points',
    description: dbDescription,
    templates: [getNetworkTemplate(networkIDs), apnTemplate],
    rows: [
      {
        title: 'Message Stats',
        panels: [
          {
            title: 'Accounting Stops (Rate)',
            targets: [
              {
                expr: 'sum(rate(accounting_stop{apn=~"$apn"}[5m])) by (apn)',
                legendFormat: '{{apn}}',
              },
            ],
            description: 'Radius accounting stops received from AP/WLC',
          },
          {
            title: 'Authorization (Rate)',
            targets: [
              {
                expr: 'sum(rate(eap_auth{apn=~"$apn"}[5m])) by (code, apn)',
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
                expr:
                  'sum(ue_reported_usage{apn=~"$apn", direction="down"}) by (apn)',
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
                expr:
                  'sum(ue_reported_usage{apn=~"$apn", direction="up"}) by (apn)',
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
                expr:
                  'avg(rate(ue_reported_usage{apn=~"$apn", direction="down"}[5m])) by (apn)',
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
                expr:
                  'avg(rate(ue_reported_usage{apn=~"$apn", direction="up"}[5m])) by (apn)',
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
            title: 'Session Stop (Rate)',
            targets: [
              {
                expr: 'sum(rate(session_stop{apn=~"$apn"}[5m])) by (apn)',
                legendFormat: '{{apn}}',
              },
            ],
            description: 'Rate of number of sessions removed for any reason',
          },
          {
            title: 'Session Timeout (Rate)',
            targets: [
              {
                expr: 'sum(rate(session_timeouts{apn=~"$apn"}[5m])) by (apn)',
                legendFormat: '{{apn}}',
              },
            ],
            description:
              'Subset of session_stop. Count of any session that times out from aaa server',
          },
          {
            title: 'Session Terminate (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(session_manager_terminate{apn=~"$apn"}[5m])) by (apn)',
                legendFormat: '{{apn}}',
              },
            ],
            description: 'Session terminations rate initiated by sessiond',
          },
        ],
      },
    ],
  };
};

export const CWFNetworkDBData = (networkIDs: Array<string>): GrafanaDBData => {
  return {
    title: 'CWF - Networks',
    description: dbDescription,
    templates: [getNetworkTemplate(networkIDs)],
    rows: [
      {
        title: 'Message Stats',
        panels: [
          {
            title: 'Authorization (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(eap_auth{networkID=~"$networkID"}[5m])) by (code, networkID)',
                legendFormat: '{{networkID}}-{{code}}',
              },
            ],
            description:
              'EAP Authorization responses, partitioned by response type (Failure, Success) where request is the sum of success and failures',
          },
          {
            title: 'Accounting Stops (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(accounting_stop{networkID=~"$networkID"}[5m])) by (networkID)',
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
                expr:
                  'sum(ue_reported_usage{networkID=~"$networkID", direction="down"}) by (networkID)',
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
                expr:
                  'sum(ue_reported_usage{networkID=~"$networkID", direction="up"}) by (networkID)',
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
                  'avg(rate(ue_reported_usage{networkID=~"$networkID", direction="down"}[5m])) by (networkID)',
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
                  'avg(rate(ue_reported_usage{networkID=~"$networkID", direction="up"}[5m])) by (networkID)',
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
                  'sum(active_sessions{networkID=~"$networkID"}) by (networkID, apn)',
                legendFormat: '{{networkID}}-{{apn}}',
              },
            ],
            description: 'Number of active user sessions in the network',
          },
          {
            title: 'Session Stop (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(session_stop{networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
            description: 'Number of sessions removed for any reason',
          },
          {
            title: 'Session Timeouts (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(session_timeouts{networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
            description:
              'Subset of session_stop. Count of any session that times out from aaa server',
          },
          {
            title: 'Session Terminate (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(session_manager_terminate{networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
            description: 'Session terminations initiated by sessiond',
          },
        ],
      },
      {
        title: 'Diameter Result Codes',
        panels: [
          {
            title: 'Gx Result Codes (Rate)',
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
            title: 'Gy Result Codes (Rate)',
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
            title: 'SWX Result Codes (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(swx_result_codes{networkID=~"$networkID"}[5m])) by (networkID, code)',
                legendFormat: '{{networkID}} - {{code}}',
              },
            ],
            description:
              'Rate of SWx responses segmented by diameter base code',
          },
          {
            title: 'SWX Experimental Result Codes (Rate)',
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
            title: 'Gx Timeouts (Rate)',
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
            title: 'Gy Timeouts (Rate)',
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
            title: 'SWX Timeouts (Rate)',
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
            title: 'Initializations (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(ocs_ccr_init_requests_total[5m])) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
            description: 'Rate of Gy CCR-I requests',
          },
          {
            title: 'Terminations (Rate)',
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
            title: 'Updates (Rate)',
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
            title: 'Initialization Failures (Rate)',
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
            title: 'Temination Failures (Rate)',
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
            title: 'Update Failures (Rate)',
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
            title: 'Initializations (Rate)',
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
            title: 'Teminations (Rate)',
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
            title: 'Updates (Rate)',
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
            title: 'Initialization Failures (Rate)',
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
            title: 'Temination Failures (Rate)',
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
            title: 'Update Failures (Rate)',
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
            title: 'MAR Requests (Rate)',
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
            title: 'SAR Requests (Rate)',
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
            title: 'MAR Failures (Rate)',
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
            title: 'SAR Failures (Rate)',
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
};

export const CWFGatewayDBData = (networkIDs: Array<string>): GrafanaDBData => {
  return {
    title: 'CWF - Gateways',
    description: dbDescription,
    templates: [getNetworkTemplate(networkIDs), gatewayTemplate],
    rows: [
      {
        title: 'Diameter Result Codes',
        panels: [
          {
            title: 'Gx Result Codes (Rate)',
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
            title: 'Gy Result Codes (Rate)',
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
            title: 'SWX Result Codes (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(swx_result_codes{networkID=~"$networkID", gatewayID=~"$gatewayID"}[5m])) by (gatewayID, code)',
                legendFormat: '{{gatewayID}} - {{code}}',
              },
            ],
            description:
              'Rate of SWx responses segmented by diameter base code',
          },
          {
            title: 'SWX Experimental Result Codes (Rate)',
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
            title: 'Gx Timeouts (Rate)',
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
            title: 'Gy Timeouts (Rate)',
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
            title: 'SWX Timeouts (Rate)',
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
            title: 'Initializations (Rate)',
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
            title: 'Terminations (Rate)',
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
            title: 'Updates (Rate)',
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
        title: 'OCS Send Failures (Rate)',
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
            title: 'Termination Failures (Rate)',
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
            title: 'Update Failures (Rate)',
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
            title: 'Initializations (Rate)',
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
            title: 'Terminations (Rate)',
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
            title: 'Updates (Rate)',
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
            title: 'Initialization Failures (Rate)',
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
            title: 'Termination Failures (Rate)',
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
            title: 'Update Failures (Rate)',
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
            title: 'MAR Requests (Rate)',
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
            title: 'SAR Requests (Rate)',
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
            title: 'MAR Failures (Rate)',
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
            title: 'SAR Failures (Rate)',
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
};
