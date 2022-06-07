/*
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
 */

import MagmaAPI from '../../api/MagmaAPI';
import nullthrows from '../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPI';
import {AlertRoutingTree} from '../views/alarms/components/AlarmAPIType';
import {AxiosResponse} from 'axios';

import type {ApiUtil} from '../views/alarms/components/AlarmsApi';

export const MagmaAlarmsApiUtil: ApiUtil = {
  useAlarmsApi: <TParams, TResponse>(
    func: (params: TParams) => Promise<{data: TResponse}>,
    params: TParams,
    cacheCounter?: string | number,
  ) => {
    return useMagmaAPI(func, params, undefined, cacheCounter);
  },
  // Alerts
  viewFiringAlerts: async ({networkId}) => {
    const alerts = await MagmaAPI.alerts.networksNetworkIdAlertsGet({
      networkId: nullthrows(networkId),
    });
    return alerts;
  },
  viewMatchingAlerts: () => {
    console.warn('not implemented');
    return Promise.resolve({data: []});
  },
  getTroubleshootingLink: async ({alertName}) => {
    return fetch('api/data/AlertLinks')
      .then(result => result.json() as Promise<Record<string, string>>)
      .then(result => {
        return {
          data: {
            link: result[alertName] ?? '',
            title: 'View Troubleshooting Documentation',
          },
        };
      });
  },
  // Alert Rules
  createAlertRule: async ({networkId, rule}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertConfigPost({
      networkId: nullthrows(networkId),
      alertConfig: rule,
    });
  },
  editAlertRule: async ({networkId, rule}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertConfigAlertNamePut(
      {
        networkId: nullthrows(networkId),
        alertName: rule.alert,
        alertConfig: rule,
      },
    );
  },
  getAlertRules: async ({networkId}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertConfigGet({
      networkId: nullthrows(networkId),
    });
  },
  deleteAlertRule: async ({ruleName, networkId}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertConfigDelete({
      networkId: nullthrows(networkId),
      alertName: ruleName,
    });
  },
  // Suppressions
  getSuppressions: () => {
    console.warn('not implemented');
    return Promise.resolve({data: []});
  },
  // Receivers
  createReceiver: async ({networkId, receiver}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertReceiverPost({
      networkId: nullthrows(networkId),
      // $FlowFixMe[prop-missing]: require_tls needs to be added
      receiverConfig: receiver,
    });
  },
  editReceiver: async ({networkId, receiver}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertReceiverReceiverPut(
      {
        networkId: nullthrows(networkId),
        receiver: receiver.name,
        // $FlowFixMe[prop-missing]: require_tls needs to be added
        receiverConfig: receiver,
      },
    );
  },
  getReceivers: async ({networkId}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertReceiverGet({
      networkId: nullthrows(networkId),
    });
  },
  deleteReceiver: async ({networkId, receiverName}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertReceiverDelete(
      {
        networkId: nullthrows(networkId),
        receiver: receiverName,
      },
    );
  },
  // Routes
  getRouteTree: async ({networkId}) => {
    return (await MagmaAPI.alerts.networksNetworkIdPrometheusAlertReceiverRouteGet(
      {
        networkId: nullthrows(networkId),
      },
      // TODO[TS-migration] it looks like the type for match in the swagger spec might be wrong
      //  see https://github.com/facebookincubator/prometheus-configmanager/blob/main/alertmanager/config/route.go
    )) as AxiosResponse<AlertRoutingTree>;
  },
  editRouteTree: async ({networkId, route}) => {
    return await MagmaAPI.alerts.networksNetworkIdPrometheusAlertReceiverRoutePost(
      {
        networkId: nullthrows(networkId),
        route: route,
      },
    );
  },
  // Metric Series
  getMetricSeries: async ({networkId}) => {
    const series = await MagmaAPI.metrics.networksNetworkIdPrometheusSeriesGet({
      networkId: nullthrows(networkId),
    });
    return series;
  },
  getMetricNames: async ({networkId}) => {
    const series = (
      await MagmaAPI.metrics.networksNetworkIdPrometheusSeriesGet({
        networkId: nullthrows(networkId),
      })
    ).data;
    const names = new Set<string>([]);
    series.forEach(value => {
      names.add(value.__name__);
    });
    return {data: Array.from(names)};
  },

  //alertmanager global config
  getGlobalConfig: () => Promise.reject('Disabled feature'),
  editGlobalConfig: () => Promise.reject('Disabled feature'),

  // Tenants
  getTenants: () => Promise.reject('Disabled feature'),
  getAlertmanagerTenancy: () => Promise.reject('Disabled feature'),
  getPrometheusTenancy: () => Promise.reject('Disabled feature'),
};
