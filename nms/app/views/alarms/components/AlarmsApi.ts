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
 */

import type {
  AlertConfig,
  AlertManagerGlobalConfig,
  AlertReceiver,
  AlertRoutingTree,
  AlertSuppression,
  FiringAlarm,
  PrometheusLabelset,
  TenancyConfig,
} from './AlarmAPIType';
import type {CancelToken} from 'axios';

export type ApiRequest = {
  networkId: string;
  cancelToken?: CancelToken;
};

type TroubleshootingLinkType = {
  link: string;
  title: string;
};

export type ApiUtil = {
  /**
   * React hook for loading data whenever a component mounts and refreshing in
   * response to prop changes. Do not use this in a class based component and
   * don't break the normal rules of react hooks.
   */
  useAlarmsApi: <TParams, TResponse>(
    func: (params: TParams) => Promise<TResponse>,
    params: TParams,
    cacheCounter?: string | number,
  ) => {
    response: TResponse | null | undefined;
    error: unknown;
    isLoading: boolean;
  };

  //alerts
  viewFiringAlerts: (req: ApiRequest) => Promise<Array<FiringAlarm>>;
  getTroubleshootingLink: (req: {
    alertName: string;
  }) => Promise<TroubleshootingLinkType | undefined | null>;
  //rules
  viewMatchingAlerts: (
    req: {expression: string} & ApiRequest,
  ) => Promise<Array<FiringAlarm>>;
  createAlertRule: (req: {rule: AlertConfig} & ApiRequest) => Promise<void>;
  // TODO: support renaming alerts by passing old alert name separately
  editAlertRule: (req: {rule: AlertConfig} & ApiRequest) => Promise<void>;
  getAlertRules: (req: ApiRequest) => Promise<Array<AlertConfig>>;
  deleteAlertRule: (req: {ruleName: string} & ApiRequest) => Promise<void>;

  // suppressions
  getSuppressions: (req: ApiRequest) => Promise<Array<AlertSuppression>>;

  // receivers
  createReceiver: (
    req: {receiver: AlertReceiver} & ApiRequest,
  ) => Promise<void>;
  editReceiver: (req: {receiver: AlertReceiver} & ApiRequest) => Promise<void>;
  getReceivers: (req: ApiRequest) => Promise<Array<AlertReceiver>>;
  deleteReceiver: (req: {receiverName: string} & ApiRequest) => Promise<void>;

  // routes
  getRouteTree: (req: ApiRequest) => Promise<AlertRoutingTree>;
  editRouteTree: (req: {route: AlertRoutingTree} & ApiRequest) => Promise<void>;

  getMetricNames: (req: ApiRequest) => Promise<Array<string>>;
  getMetricSeries: (
    req: {name: string} & ApiRequest,
  ) => Promise<Array<PrometheusLabelset>>;

  //alertmanager global config
  getGlobalConfig: (req: ApiRequest) => Promise<AlertManagerGlobalConfig>;
  editGlobalConfig: (
    req: ApiRequest & {config: Partial<AlertManagerGlobalConfig>},
  ) => Promise<void>;

  // tenants
  getTenants: (req: ApiRequest) => Promise<Array<string>>;

  getAlertmanagerTenancy: (req: ApiRequest) => Promise<TenancyConfig>;
  getPrometheusTenancy: (req: ApiRequest) => Promise<TenancyConfig>;
};
