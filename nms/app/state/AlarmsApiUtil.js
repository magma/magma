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
 *
 * @flow strict-local
 * @format
 */

import MagmaV1API from '../../generated/WebClient';
import nullthrows from '../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPIFlow';

import type {ApiUtil} from '../views/alarms/components/AlarmsApi';

export const MagmaAlarmsApiUtil: ApiUtil = {
  useAlarmsApi: <TParams: {...}, TResponse>(
    func: TParams => Promise<TResponse>,
    params: TParams,
    cacheCounter?: string | number,
  ) => {
    return useMagmaAPI(func, params, undefined, cacheCounter);
  },
  // Alerts
  viewFiringAlerts: async ({networkId}) => {
    const alerts = await MagmaV1API.getNetworksByNetworkIdAlerts({
      networkId: nullthrows(networkId),
    });
    return alerts;
  },
  viewMatchingAlerts: async ({networkId: _, expression: __}) => {
    console.warn('not implemented');
    return [];
  },
  getTroubleshootingLink: async ({alertName}) => {
    return fetch('api/data/AlertLinks')
      .then(result => result.json())
      .then(result => {
        return {
          link: result[alertName] ?? '',
          title: 'View Troubleshooting Documentation',
        };
      });
  },

  // Alert Rules
  createAlertRule: async ({networkId, rule}) => {
    await MagmaV1API.postNetworksByNetworkIdPrometheusAlertConfig({
      networkId: nullthrows(networkId),
      alertConfig: rule,
    });
  },
  editAlertRule: async ({networkId, rule}) => {
    await MagmaV1API.putNetworksByNetworkIdPrometheusAlertConfigByAlertName({
      networkId: nullthrows(networkId),
      alertName: rule.alert,
      alertConfig: rule,
    });
  },
  getAlertRules: async ({networkId}) => {
    return await MagmaV1API.getNetworksByNetworkIdPrometheusAlertConfig({
      networkId: nullthrows(networkId),
    });
  },
  deleteAlertRule: async ({ruleName, networkId}) => {
    await MagmaV1API.deleteNetworksByNetworkIdPrometheusAlertConfig({
      networkId: nullthrows(networkId),
      alertName: ruleName,
    });
  },
  // Suppressions
  getSuppressions: () => {
    console.warn('not implemented');
    return Promise.resolve([]);
  },
  // Receivers
  createReceiver: async ({networkId, receiver}) => {
    await MagmaV1API.postNetworksByNetworkIdPrometheusAlertReceiver({
      networkId: nullthrows(networkId),
      // $FlowFixMe[prop-missing]: require_tls needs to be added
      receiverConfig: receiver,
    });
  },
  editReceiver: async ({networkId, receiver}) => {
    await MagmaV1API.putNetworksByNetworkIdPrometheusAlertReceiverByReceiver({
      networkId: nullthrows(networkId),
      receiver: receiver.name,
      // $FlowFixMe[prop-missing]: require_tls needs to be added
      receiverConfig: receiver,
    });
  },
  getReceivers: async ({networkId}) => {
    return await MagmaV1API.getNetworksByNetworkIdPrometheusAlertReceiver({
      networkId: nullthrows(networkId),
    });
  },
  deleteReceiver: async ({networkId, receiverName}) => {
    await MagmaV1API.deleteNetworksByNetworkIdPrometheusAlertReceiver({
      networkId: nullthrows(networkId),
      receiver: receiverName,
    });
  },
  // Routes
  getRouteTree: async ({networkId}) => {
    return await MagmaV1API.getNetworksByNetworkIdPrometheusAlertReceiverRoute({
      networkId: nullthrows(networkId),
    });
  },
  editRouteTree: async ({networkId, route}) => {
    await MagmaV1API.postNetworksByNetworkIdPrometheusAlertReceiverRoute({
      networkId: nullthrows(networkId),
      route: route,
    });
  },
  // Metric Series
  getMetricSeries: async ({networkId}) => {
    const series = await MagmaV1API.getNetworksByNetworkIdPrometheusSeries({
      networkId: nullthrows(networkId),
    });
    return series;
  },
  getMetricNames: async ({networkId}) => {
    const series = await MagmaV1API.getNetworksByNetworkIdPrometheusSeries({
      networkId: nullthrows(networkId),
    });
    const names = new Set([]);
    series.forEach(value => {
      names.add(value.__name__);
    });
    return Array.from(names);
  },

  //alertmanager global config
  getGlobalConfig: _ => Promise.reject('Disabled feature'),
  editGlobalConfig: _ => Promise.reject('Disabled feature'),

  // Tenants
  getTenants: _ => Promise.reject('Disabled feature'),
  getAlertmanagerTenancy: _ => Promise.reject('Disabled feature'),
  getPrometheusTenancy: _ => Promise.reject('Disabled feature'),
};
