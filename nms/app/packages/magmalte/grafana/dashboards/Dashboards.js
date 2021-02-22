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

//$FlowFixMe TODO: type this package
import * as Grafana from 'grafana-dash-gen';

const netIDVar = 'networkID';
const gwIDVar = 'gatewayID';

const variableSortNumbers: {[VariableSortOption]: number} = {
  none: 0,
  'alpha-asc': 1,
  'alpha-desc': 2,
  'num-asc': 3,
  'num-desc': 4,
  'alpha-insensitive-asc': 5,
  'alpha-insensitive-desc': 6,
};

export const getNetworkTemplate = (
  networkIDs: Array<string>,
): TemplateConfig => {
  return customVariableTemplate({
    name: netIDVar,
    type: 'custom',
    options: networkIDs,
    sort: 'alpha-insensitive-asc',
    includeAll: true,
  });
};

// This templating schema will produce a variable in the dashboard
// named gatewayID which is a multi-selectable option of all the
// gateways associated with this organization that exist for the
// currently selected $networkID. $networkID variable must also
// be configured for this dashboard in order for it to work
export const gatewayTemplate: TemplateConfig = variableTemplate({
  name: gwIDVar,
  query: `label_values({networkID=~"$networkID",gatewayID=~".+"}, ${gwIDVar})`,
  regex: `/.+/`,
  sort: 'alpha-insensitive-asc',
  includeAll: true,
});

export const msisdnTemplate = variableTemplate({
  name: 'msisdn',
  query: `label_values(msisdn)`,
  regex: `/.+/`,
  sort: 'num-asc',
  includeAll: false,
});

export const apnTemplate = variableTemplate({
  name: 'apn',
  query: `label_values({networkID=~"$networkID",apn=~".+"},apn)`,
  regex: `/.+/`,
  sort: 'alpha-insensitive-asc',
  includeAll: true,
});

export const NetworkDBData = (networkIDs: Array<string>): GrafanaDBData => {
  return {
    title: 'Networks',
    description:
      'Metrics relevant to the whole network. Do not edit: edits will be overwritten. Save this dashboard under another name to copy and edit.',
    templates: [getNetworkTemplate(networkIDs)],
    rows: [
      {
        title: 'Connections',
        panels: [
          {
            title: 'Number of Connected UEs',
            targets: [
              {
                expr:
                  'sum(ue_connected{networkID=~"$networkID"}) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
            aggregates: {avg: true, max: true},
          },
          {
            title: 'Number of Registered UEs',
            targets: [
              {
                expr:
                  'sum(ue_registered{networkID=~"$networkID"}) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
            aggregates: {avg: true, max: true},
          },
          {
            title: 'Attach Success Rate',
            targets: [
              {
                expr:
                  '(sum by(networkID) (increase(ue_attach{action="attach_accept_sent",networkID=~"$networkID"}[3h]))) * 100 / (sum by(networkID) (increase(ue_attach{action=~"attach_accept_sent|attach_reject_sent|attach_abort",networkID=~"$networkID"}[3h])))',
                legendFormat: '{{networkID}}',
              },
            ],
            yMax: 100,
          },
          {
            title: 'Duplicate Attach Requests (1h Increase)',
            targets: [
              {
                expr:
                  'sum by(networkID) (increase(duplicate_attach_request{networkID=~"$networkID"}[1h]))',
                legend: '{{networkID}}',
              },
            ],
          },
          {
            title: 'Number of Connected eNBs',
            targets: [
              {
                expr:
                  'sum(enb_connected{networkID=~"$networkID"}) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
            aggregates: {avg: true, max: true},
          },
        ],
      },
      {
        title: 'S1',
        panels: [
          {
            title: 'S1 Setup (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(s1_setup{networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(s1_setup{networkID=~"$networkID",result="success"}[5m])) by (networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(s1_setup{networkID=~"$networkID"}[5m]))by(networkID)-sum(rate(s1_setup{result="success",networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'S1 Setup (1h Increase)',
            targets: [
              {
                expr:
                  'sum(increase(s1_setup{networkID=~"$networkID"}[1h])) by (networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(s1_setup{networkID=~"$networkID",result="success"}[1h])) by (networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(s1_setup{networkID=~"$networkID"}[1h]))by(networkID)-sum(increase(s1_setup{result="success",networkID=~"$networkID"}[1h])) by (networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
        ],
      },
      {
        title: 'Attach/Detach',
        panels: [
          {
            title: 'Attach/Reg Attempts (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(ue_attach{networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(ue_attach{networkID=~"$networkID",result="attach_proc_successful"}[5m])) by (networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(ue_attach{networkID=~"$networkID"}[5m])) by (networkID) -sum(rate(s1_setup{result="attach_proc_successful",networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'Attach/Reg Attempts (1h Increase)',
            targets: [
              {
                expr:
                  'sum(increase(ue_attach{networkID=~"$networkID"}[1h])) by (networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(ue_attach{networkID=~"$networkID",result="attach_proc_successful"}[1h])) by (networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(ue_attach{networkID=~"$networkID"}[1h])) by (networkID) -sum(increase(s1_setup{result="attach_proc_successful",networkID=~"$networkID"}[1h])) by (networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'Detach/Dereg Attempts (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(ue_detach{networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(ue_detach{networkID=~"$networkID",result="attach_proc_successful"}[5m])) by (networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(ue_detach{networkID=~"$networkID"}[5m])) by (networkID) -sum(rate(s1_setup{result="attach_proc_successful",networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'Detach/Dereg Attempts (1h Increase)',
            targets: [
              {
                expr:
                  'sum(increase(ue_detach{networkID=~"$networkID"}[1h])) by (networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(ue_detach{networkID=~"$networkID",result="attach_proc_successful"}[1h])) by (networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(ue_detach{networkID=~"$networkID"}[1h])) by (networkID) -sum(increase(s1_setup{result="attach_proc_successful",networkID=~"$networkID"}[1h])) by (networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'Attach Success Rate',
            targets: [
              {
                expr:
                  '(sum by(networkID) (increase(ue_attach{action="attach_accept_sent",networkID=~"$networkID"}[3h]))) * 100 / (sum by(networkID) (increase(ue_attach{action=~"attach_accept_sent|attach_reject_sent|attach_abort",networkID=~"$networkID"}[3h])))',
                legendFormat: '{{networkID}}',
              },
            ],
            yMax: 100,
          },
        ],
      },
      {
        title: 'Connection Status',
        panels: [
          {
            title: 'GPS Connection Uptime',
            targets: [
              {
                expr:
                  'avg(enodeb_gps_connected{networkID=~"$networkID"}) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
            unit: 's',
          },
          {
            title: 'Device Transmitting Status',
            targets: [
              {
                expr:
                  'avg(enodeb_rf_tx_enabled{networkID=~"$networkID"}) by (networkID)',
                legendFormat: '{{networkID}}',
              },
            ],
          },
        ],
      },
      {
        title: 'Service Requests',
        panels: [
          {
            title: 'Service Requests (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(service_request{networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(service_request{networkID=~"$networkID",result="success"}[5m])) by (networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(service_request{networkID=~"$networkID"}[5m])) by (networkID)-sum(rate(s1_setup{result="success",networkID=~"$networkID"}[5m])) by (networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'Service Requests (1h Increase)',
            targets: [
              {
                expr:
                  'sum(increase(service_request{networkID=~"$networkID"}[1h])) by (networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(service_request{networkID=~"$networkID",result="success"}[1h])) by (networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(service_request{networkID=~"$networkID"}[1h])) by (networkID)-sum(increase(s1_setup{result="success",networkID=~"$networkID"}[1h])) by (networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'Service Request Success Rate',
            targets: [
              {
                expr:
                  'round((sum by(networkID) (increase(service_request{networkID=~"$networkID", result="success"}[3h]))) *100 / ((sum by(networkID) (increase(service_request{networkID=~"$networkID", result="failure"}[3h]))) + (sum by(networkID) (increase(service_request{networkID=~"$networkID", result="success"}[3h])))))',
                legendFormat: '{{networkID}}',
              },
            ],
            yMax: 100,
          },
        ],
      },
      {
        title: 's6a',
        panels: [
          {
            title: 's6a Auth Failure (Rate)',
            targets: [
              {
                expr: 'rate(s6a_auth_failure{networkID=~"$networkID"}[5m])',
                legendFormat: '{{networkID}}',
              },
            ],
          },
          {
            title: 's6a Auth Failure (1h Increase)',
            targets: [
              {
                expr: 'increase(s6a_auth_failure{networkID=~"$networkID"}[1h])',
                legendFormat: '{{networkID}}',
              },
            ],
          },
          {
            title: 's6a Auth Success (Rate)',
            targets: [
              {
                expr: 'rate(s6a_auth_success{networkID=~"$networkID"}[5m])',
                legendFormat: '{{networkID}}',
              },
            ],
          },
          {
            title: 's6a Auth Success (1h Increase)',
            targets: [
              {
                expr: 'increase(s6a_auth_success{networkID=~"$networkID"}[1h])',
                legendFormat: '{{networkID}}',
              },
            ],
          },
          {
            title: 's6a Authentication Success Rate',
            targets: [
              {
                expr:
                  'sum by(networkID) (increase(s6a_auth_success{networkID=~"$networkID"}[3h]) * 100) / (sum by (networkID) (increase(s6a_auth_success{networkID=~"$networkID"}[3h])) + sum by(networkID) (increase(s6a_auth_failure{networkID=~"$networkID"}[3h])))',
                legendFormat: '{{networkID}}',
              },
            ],
            yMax: 100,
          },
        ],
      },
      {
        title: 'Sessions',
        panels: [
          {
            title: 'Session Create Success Rate',
            targets: [
              {
                expr:
                  '(sum by(networkID) (increase(mme_spgw_create_session_rsp{networkID=~"$networkID", result="success"}[3h]) * 100)) / (sum by(networkID) (increase(mme_spgw_create_session_rsp{networkID=~"$networkID"}[3h])))',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            yMax: 100,
          },
        ],
      },
    ],
  };
};

export const GatewayDBData = (networkIDs: Array<string>): GrafanaDBData => {
  return {
    title: 'Gateways',
    description:
      'Metrics relevant to the gateways. Do not edit: edits will be overwritten. Save this dashboard under another name to copy and edit.',
    templates: [getNetworkTemplate(networkIDs), gatewayTemplate],
    rows: [
      {
        title: 'Connections',
        panels: [
          {
            title: 'E-Node B Status',
            targets: [
              {
                expr:
                  'enodeb_rf_tx_enabled{gatewayID=~"$gatewayID",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
          },
          {
            title: 'Connected Subscribers',
            targets: [
              {
                expr:
                  'ue_connected{gatewayID=~"$gatewayID",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
          },
          {
            title: 'Attach Success Rate',
            targets: [
              {
                expr:
                  '(sum by(networkID,gatewayID) (increase(ue_attach{action="attach_accept_sent",networkID=~"$networkID",gatewayID=~"$gatewayID"}[3h]))) * 100 / (sum by(networkID,gatewayID) (increase(ue_attach{action=~"attach_accept_sent|attach_reject_sent|attach_abort",networkID=~"$networkID",gatewayID=~"$gatewayID"}[3h])))',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            yMax: 100,
          },
          {
            title: 'Duplicate Attach Requests (1h Increase)',
            targets: [
              {
                expr:
                  'sum by(networkID,gatewayID) (increase(duplicate_attach_request{networkID=~"$networkID",gatewayID=~"$gatewayID"}[1h]))',
                legend: '{{networkID}} - {{gatewayID}}',
              },
            ],
          },
        ],
      },
      {
        title: 'Traffic',
        panels: [
          {
            title: 'Download Throughput',
            targets: [
              {
                expr:
                  'rate(gtp_port_user_plane_dl_bytes{gatewayID=~"$gatewayID",networkID=~"$networkID"}[5m])',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            unit: 'Bps',
          },
          {
            title: 'Upload Throughput',
            targets: [
              {
                expr:
                  'rate(gtp_port_user_plane_ul_bytes{gatewayID=~"$gatewayID",networkID=~"$networkID"}[5m])',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            unit: 'Bps',
          },
        ],
      },
      {
        title: 'Service Requests',
        panels: [
          {
            title: 'Service Requests (Rate)',
            targets: [
              {
                expr:
                  'sum(rate(service_request{gatewayID=~"$gatewayID",networkID=~"$networkID"}[5m])) by (gatewayID,networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(service_request{gatewayID=~"$gatewayID",networkID=~"$networkID",result="success"}[5m])) by (gatewayID,networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(rate(service_request{gatewayID=~"$gatewayID",networkID=~"$networkID"}[5m])) by (gatewayID,networkID)-sum(rate(s1_setup{gatewayID=~"$gatewayID",result="success",networkID=~"$networkID"}[5m])) by (gatewayID,networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'Service Requests (1h Increase)',
            targets: [
              {
                expr:
                  'sum(increase(service_request{gatewayID=~"$gatewayID",networkID=~"$networkID"}[1h])) by (gatewayID,networkID)',
                legendFormat: 'Total: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(service_request{gatewayID=~"$gatewayID",networkID=~"$networkID",result="success"}[1h])) by (gatewayID,networkID)',
                legendFormat: 'Success: {{networkID}}',
              },
              {
                expr:
                  'sum(increase(service_request{gatewayID=~"$gatewayID",networkID=~"$networkID"}[1h])) by (gatewayID,networkID)-sum(increase(s1_setup{gatewayID=~"$gatewayID",result="success",networkID=~"$networkID"}[1h])) by (gatewayID,networkID)',
                legendFormat: 'Failure: {{networkID}}',
              },
            ],
          },
          {
            title: 'Service Request Success Rate',
            targets: [
              {
                expr:
                  'round((sum by(gatewayID,networkID) (increase(service_request{gatewayID=~"$gatewayID",networkID=~"$networkID", result="success"}[3h]))) *100 / ((sum by(gatewayID,networkID) (increase(service_request{gatewayID=~"$gatewayID",networkID=~"$networkID", result="failure"}[3h]))) + (sum by(gatewayID,networkID) (increase(service_request{networkID=~"$networkID", result="success"}[3h])))))',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            yMax: 100,
          },
        ],
      },
      {
        title: 'Miscellaneous',
        panels: [
          {
            title: 'Latency',
            targets: [
              {
                expr:
                  'magmad_ping_rtt_ms{gatewayID=~"$gatewayID",networkID=~"$networkID",metric="rtt_ms"}',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            unit: 's',
          },
          {
            title: 'Gateway CPU %',
            targets: [
              {
                expr:
                  'cpu_percent{gatewayID=~"$gatewayID",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            unit: 'percent',
          },
          {
            title: 'Temperature',
            targets: [
              {
                expr:
                  'temperature{gatewayID=~"$gatewayID",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}} - {{sensor}}',
              },
            ],
            yMin: null,
            unit: 'celsius',
          },
          {
            title: 'Disk %',
            targets: [
              {
                expr:
                  'disk_percent{gatewayID=~"$gatewayID",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            unit: 'percent',
          },
        ],
      },
      {
        title: 's6a',
        panels: [
          {
            title: 's6a Auth Failure (Rate)',
            targets: [
              {
                expr:
                  'rate(s6a_auth_failure{gatewayID=~"$gatewayID",networkID=~"$networkID"}[5m])',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
          },
          {
            title: 's6a Auth Failure (1h Increase)',
            targets: [
              {
                expr:
                  'increase(s6a_auth_failure{gatewayID=~"$gatewayID",networkID=~"$networkID"}[1h])',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
          },
          {
            title: 's6a Auth Success (Rate)',
            targets: [
              {
                expr:
                  'rate(s6a_auth_success{gatewayID=~"$gatewayID",networkID=~"$networkID"}[5m])',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
          },
          {
            title: 's6a Auth Success (1h Increase)',
            targets: [
              {
                expr:
                  'increase(s6a_auth_success{gatewayID=~"$gatewayID",networkID=~"$networkID"}[1h])',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
          },
          {
            title: 's6a Authentication Success Rate',
            targets: [
              {
                expr:
                  'sum by(networkID,gatewayID) (increase(s6a_auth_success{gatewayID=~"$gatewayID",networkID=~"$networkID"}[3h]) * 100) / (sum by (gatewayID,networkID) (increase(s6a_auth_success{gatewayID=~"$gatewayID",networkID=~"$networkID"}[3h])) + sum by(gatewayID,networkID) (increase(s6a_auth_failure{gatewayID=~"$gatewayID",networkID=~"$networkID"}[3h])))',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            yMax: 100,
          },
        ],
      },
      {
        title: 'Sessions',
        panels: [
          {
            title: 'Session Create Success Rate',
            targets: [
              {
                expr:
                  '(sum by(gatewayID,networkID) (increase(mme_spgw_create_session_rsp{gatewayID=~"$gatewayID",networkID=~"$networkID", result="success"}[3h]) * 100)) / (sum by(gatewayID,networkID) (increase(mme_spgw_create_session_rsp{gatewayID=~"$gatewayID",networkID=~"$networkID"}[3h])))',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            yMax: 100,
          },
        ],
      },
    ],
  };
};

export const SubscriberDBData = (networkIDs: Array<string>): GrafanaDBData => {
  return {
    title: 'Subscribers',
    description:
      'Metrics relevant to subscribers. Do not edit: edits will be overwritten. Save this dashboard under another name to copy and edit.',
    templates: [getNetworkTemplate(networkIDs), msisdnTemplate],
    rows: [
      {
        title: '',
        panels: [
          {
            title: 'UE Data Usage In',
            targets: [
              {
                expr:
                  'sum(ue_reported_usage{msisdn=~"$msisdn", direction="down"}) by (IMSI, apn, msisdn)',
                legendFormat: '{{IMSI}}, MSISDN: {{msisdn}}, APN: {{apn}}',
              },
            ],
            unit: 'decbytes',
            description: 'Inbound data per subscriber measured in bytes.',
          },
          {
            title: 'UE Data Usage Out',
            targets: [
              {
                expr:
                  'sum(ue_reported_usage{msisdn=~"$msisdn", direction="up"}) by (IMSI, apn, msisdn)',
                legendFormat: '{{IMSI}}, MSISDN: {{msisdn}}, APN: {{apn}}',
              },
            ],
            unit: 'decbytes',
            description: 'Outbound data per subscriber measured in bytes.',
          },
          {
            title: 'Throughput In',
            targets: [
              {
                expr:
                  'avg(rate(ue_reported_usage{msisdn=~"$msisdn", direction="down"}[5m])) by (IMSI, apn, msisdn)',
                legendFormat: '{{IMSI}}, MSISDN: {{msisdn}}, APN: {{apn}}',
              },
            ],
            unit: 'Bps',
            description:
              'Inbound data rate per subscriber measured in bytes/second.',
          },
          {
            title: 'Throughput Out',
            targets: [
              {
                expr:
                  'avg(rate(ue_reported_usage{msisdn=~"$msisdn", direction="up"}[5m])) by (IMSI, apn, msisdn)',
                legendFormat: '{{IMSI}}, MSISDN: {{msisdn}}, APN: {{apn}}',
              },
            ],
            unit: 'Bps',
            description:
              'Outbound data rate per subscriber measured in bytes/second.',
          },
        ],
      },
    ],
  };
};

export const InternalDBData = (networkIDs: Array<string>): GrafanaDBData => {
  return {
    title: 'Internal',
    description:
      'Metrics relevant to the internals of gateways. Do not edit: edits will be overwritten. Save this dashboard under another name to copy and edit.',
    templates: [getNetworkTemplate(networkIDs), gatewayTemplate],
    rows: [
      {
        title: '',
        panels: [
          {
            title: 'Physical Memory Available Percent',
            targets: [
              {
                expr:
                  'mem_free{gatewayID=~"$gatewayID"}/mem_total{gatewayID=~"$gatewayID",networkID=~"$networkID"} * 100',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
          },
          {
            title: 'Temperature',
            targets: [
              {
                expr:
                  'temperature{gatewayID=~"$gatewayID",sensor="coretemp_0",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}} - {{sensor}}',
              },
            ],
            unit: 'percent',
          },
          {
            title: 'Virtual Memory Percent',
            targets: [
              {
                expr:
                  'virtual_memory_percent{gatewayID=~"$gatewayID",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            unit: 'percent',
          },
          {
            title: 'Backhaul Latency',
            targets: [
              {
                expr:
                  'magmad_ping_rtt_ms{gatewayID=~"$gatewayID",service="magmad",host="8.8.8.8",metric="rtt_ms",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}}',
              },
            ],
            unit: 's',
          },
          {
            title: 'System Uptime',
            targets: [
              {
                expr:
                  'process_uptime_seconds{gatewayID=~"$gatewayID",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}}-{{service}}',
              },
            ],
            unit: 's',
          },
          {
            title: 'Number of Service Restarts',
            targets: [
              {
                expr:
                  'unexpected_service_restarts{gatewayID=~"$gatewayID",networkID=~"$networkID"}',
                legendFormat: '{{networkID}} - {{gatewayID}}-{{service_name}}',
              },
            ],
          },
        ],
      },
    ],
  };
};

export function createDashboard(dbdata: GrafanaDBData) {
  const rows = dbdata.rows.map(conf => {
    const row = new Grafana.Row({title: conf.title});
    conf.panels.forEach(panel => row.addPanel(newPanel(panel)));
    return row;
  });
  const db = new Grafana.Dashboard({
    schemaVersion: 6,
    title: dbdata.title,
    templating: dbdata.templates,
    description: dbdata.description,
    rows,
  });
  db.state.editable = false;

  // Necessary to make custom templates display "all" option
  for (const template of db.state?.templating?.list) {
    if (template.type === 'custom' && template.includeAll) {
      template.options.unshift({selected: true, text: 'All', value: '$__all'});
    }
    template.current = template.options[0];
  }
  return db;
}

export type GrafanaDBData = {
  title: string,
  description: string,
  rows: Array<GrafanaDBRow>,
  templates: Array<TemplateConfig>,
};

type GrafanaDBRow = {
  title: string,
  panels: PanelParams[],
};

type PanelParams = {
  title: string,
  targets: Array<{expr: string, legendFormat?: string}>,
  unit?: string,
  yMin?: ?number,
  yMax?: number,
  aggregates?: {avg?: boolean, max?: boolean},
  description?: string,
};

function newPanel(params: PanelParams) {
  const pan = new Grafana.Panels.Graph({
    title: params.title,
    span: 6,
    datasource: 'default',
    description: params.description ?? '',
  });
  // Have to add this after to avoid grafana-dash-gen from forcing the target
  // into a Graphite format
  pan.state.targets = params.targets;

  // "short" is the default unit for grafana (no unit)
  pan.state.y_formats[0] = params.unit ?? 'short';

  // yMin should be 0 at minimum unless otherwise specified.
  // null is used to indicate 'auto' in grafana
  if (params.yMin === null) {
    pan.state.grid.leftMin = null;
  } else {
    pan.state.grid.leftMin = params.yMin ?? 0;
  }

  if (params.yMax !== null && params.yMax !== undefined) {
    pan.state.grid.leftMax = params.yMax;
  }

  pan.state.legend.avg = params.aggregates?.avg ?? false;
  pan.state.legend.max = params.aggregates?.max ?? false;
  return pan;
}

export type TemplateParams = {
  name: string,
  query?: string,
  options?: Array<string>,
  regex?: string,
  sort?: VariableSortOption,
  includeAll: boolean,
  type?: string,
};

type VariableSortOption =
  | 'none'
  | 'alpha-asc'
  | 'alpha-desc'
  | 'num-asc'
  | 'num-desc'
  | 'alpha-insensitive-asc'
  | 'alpha-insensitive-desc';

export function variableTemplate(params: TemplateParams): TemplateConfig {
  return {
    allValue: '.+',
    definition: params.query,
    hide: 0,
    includeAll: params.includeAll,
    allFormat: 'glob',
    multi: true,
    name: params.name,
    query: params.query ?? '',
    options: params.options ?? [],
    regex: params.regex,
    type: params.type ?? 'query',
    refresh: true,
    useTags: false,
    sort: params.sort ? variableSortNumbers[params.sort] : 0,
  };
}

export function customVariableTemplate(params: TemplateParams): TemplateConfig {
  return {
    options: params.options ?? [],
    includeAll: true,
    name: params.name,
    multi: true,
    allFormat: 'glob',
    allValue: '.+',
  };
}

export type TemplateConfig = {
  allValue: string,
  definition?: string,
  hide?: number,
  includeAll: boolean,
  allFormat: string,
  multi: boolean,
  name: string,
  query?: string,
  options: Array<string>,
  regex?: string,
  type?: string,
  refresh?: boolean,
  useTags?: boolean,
};
