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

import type {MetricGraphConfig} from '@fbcnms/ui/insights/Metrics';

import AppContext from '@fbcnms/ui/context/AppContext';
import AssessmentIcon from '@material-ui/icons/Assessment';
import ExploreIcon from '@material-ui/icons/Explore';
import Explorer from '../../views/metrics/Explorer';
import GatewayMetrics from '@fbcnms/ui/insights/GatewayMetrics';
import Grafana from '../Grafana';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import NetworkKPIs from './NetworkKPIs';
import React, {useContext} from 'react';
import TopBar from '../../components/TopBar';

import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

const CONFIGS: Array<MetricGraphConfig> = [
  {
    basicQueryConfigs: [
      {
        metric: 's1_connection',
        filters: [{name: 'service', value: 'mme'}],
      },
    ],
    label: 'E-Node B Status',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'ue_connected',
        filters: [{name: 'service', value: 'mme'}],
      },
    ],
    label: 'Connected Subscribers',
  },
  {
    customQueryConfigs: [
      {
        resolveQuery: gw =>
          `gtp_port_user_plane_dl_bytes{gatewayID="${gw}", service="pipelined"}/1000`,
      },
    ],
    basicQueryConfigs: [],
    label: 'Download Throughput',
    unit: ' Mbps',
  },
  {
    customQueryConfigs: [
      {
        resolveQuery: gw =>
          `gtp_port_user_plane_ul_bytes{gatewayID="${gw}", service="pipelined"}/1000`,
      },
    ],
    basicQueryConfigs: [],
    label: 'Upload Throughput',
    unit: ' Mbps',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'magmad_ping_rtt_ms',
        filters: [
          {name: 'service', value: 'magmad'},
          {name: 'metric', value: 'rtt_ms'},
        ],
      },
    ],
    label: 'Latency',
    unit: ' ms',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'cpu_percent',
        filters: [{name: 'service', value: 'magmad'}],
      },
    ],
    label: 'Gateway CPU (%)',
    unit: '%',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'temperature',
        filters: [{name: 'service', value: 'magmad'}],
      },
    ],
    label: 'Temperature (℃)',
    unit: '℃',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'disk_percent',
        filters: [{name: 'service', value: 'magmad'}],
      },
    ],
    label: 'Disk (%)',
    unit: '%',
  },
  {
    basicQueryConfigs: [
      {
        metric: 's6a_auth_success',
        filters: [{name: 'service', value: 'subscriberdb'}],
      },
    ],
    label: 's6a Auth Success',
    unit: '',
  },
  {
    basicQueryConfigs: [
      {
        metric: 's6a_auth_failure',
        filters: [{name: 'service', value: 'subscriberdb'}],
      },
    ],
    label: 's6a Auth Failure',
    unit: '',
  },
  {
    basicQueryConfigs: [
      {
        metric: 'enodeb_rf_tx_enabled',
        filters: [{name: 'service', value: 'enodebd'}],
      },
    ],
    label: 'E-NodeB Transmitting',
    unit: '',
  },
];

const INTERNAL_CONFIGS: Array<MetricGraphConfig> = [
  {
    label: 'Memory Utilization',
    basicQueryConfigs: [],
    filters: [],
    unit: '',
    customQueryConfigs: [
      {
        resolveQuery: gid =>
          `mem_free{gatewayID="${gid}"} / mem_total{gatewayID="${gid}"}`,
      },
    ],
  },
  {
    label: 'Temperature',
    basicQueryConfigs: [
      {
        metric: 'temperature',
        filters: [
          {name: 'service', value: 'magmad'},
          {name: 'sensor', value: 'coretemp_0'},
        ],
      },
    ],
    unit: 'C',
  },
  {
    label: 'Virtual Memory',
    basicQueryConfigs: [
      {
        metric: 'virtual_memory_percent',
        filters: [{name: 'service', value: 'magmad'}],
      },
    ],
    unit: '%',
  },
  {
    label: 'Backhaul Latency',
    basicQueryConfigs: [
      {
        metric: 'magmad_ping_rtt_ms',
        filters: [
          {name: 'service', value: 'magmad'},
          {name: 'host', value: '8.8.8.8'},
          {name: 'metric', value: 'rtt_ms'},
        ],
      },
    ],
    unit: 'ms',
  },
  {
    label: 'System Uptime',
    basicQueryConfigs: [
      {
        metric: 'process_uptime_seconds',
        filters: [{name: 'service', value: 'magmad'}],
      },
    ],
    unit: 's',
  },
  {
    label: 'Number of Service Restarts',
    basicQueryConfigs: [
      {
        metric: 'unexpected_service_restarts',
        filters: [{name: 'service', value: 'magmad'}],
      },
    ],
    unit: 'restarts',
  },
];

function GatewayMetricsGraphs() {
  return <GatewayMetrics configs={CONFIGS} />;
}

function InternalMetrics() {
  return <GatewayMetrics configs={INTERNAL_CONFIGS} />;
}

function GrafanaDashboard() {
  return <Grafana grafanaURL={'/grafana'} />;
}

export default function () {
  const lteNetworkMetrics = useContext(AppContext).isFeatureEnabled(
    'lte_network_metrics',
  );
  if (!lteNetworkMetrics) {
    return <GatewayMetricsGraphs />;
  }

  const {relativePath, relativeUrl} = useRouter();

  const grafanaEnabled =
    useContext(AppContext).isFeatureEnabled('grafana_metrics') &&
    useContext(AppContext).user.isSuperUser;

  const tabNames = ['gateways', 'network', 'internal'];
  if (grafanaEnabled) {
    tabNames.push('grafana');
  }

  let tabList = [];
  if (!grafanaEnabled) {
    tabList = [
      {
        component: {NestedRouteLink},
        label: 'Gateways',
        to: '/gateways',
      },
      {
        component: {NestedRouteLink},
        label: 'Internal',
        to: '/internal',
      },
    ];
  } else {
    tabList = [
      {
        icon: AssessmentIcon,
        component: {NestedRouteLink},
        label: 'Grafana',
        to: '/grafana',
      },
      {
        icon: ExploreIcon,
        component: {NestedRouteLink},
        label: 'Explorer',
        to: '/explorer',
      },
    ];
  }

  return (
    <>
      <TopBar header={'Metrics'} tabs={tabList} />
      <Switch>
        {!grafanaEnabled ? (
          <>
            <Route
              path={relativePath('/gateways')}
              component={GatewayMetricsGraphs}
            />
            <Route path={relativePath('/network')} component={NetworkKPIs} />
            <Route
              path={relativePath('/internal')}
              component={InternalMetrics}
            />
            <Redirect to={relativeUrl('/gateways')} />
          </>
        ) : (
          <>
            <Route
              path={relativePath('/grafana')}
              component={GrafanaDashboard}
            />
            <Route path={relativePath('/explorer')} component={Explorer} />
            <Redirect to={relativeUrl('/grafana')} />
          </>
        )}
      </Switch>
    </>
  );
}
