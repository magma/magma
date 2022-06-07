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

import {PromFiringAlert} from '../../../../generated-ts';
import type {
  AlertConfig,
  AlertManagerGlobalConfig,
  AlertReceiver,
  AlertRoutingTree,
  AlertSuppression,
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
    func: (params: TParams) => Promise<{data: TResponse}>,
    params: TParams,
    cacheCounter?: string | number,
  ) => {
    response: TResponse | null | undefined;
    error?: Error;
    isLoading: boolean;
  };

  //alerts
  viewFiringAlerts: (
    req: ApiRequest,
  ) => Promise<{data: Array<PromFiringAlert>}>;
  getTroubleshootingLink: (req: {
    alertName: string;
  }) => Promise<{data: TroubleshootingLinkType | undefined | null}>;
  //rules
  viewMatchingAlerts: (
    req: {expression: string} & ApiRequest,
  ) => Promise<{data: Array<PromFiringAlert>}>;
  createAlertRule: (
    req: {rule: AlertConfig} & ApiRequest,
  ) => Promise<{data: void}>;
  // TODO: support renaming alerts by passing old alert name separately
  editAlertRule: (
    req: {rule: AlertConfig} & ApiRequest,
  ) => Promise<{data: void}>;
  getAlertRules: (req: ApiRequest) => Promise<{data: Array<AlertConfig>}>;
  deleteAlertRule: (
    req: {ruleName: string} & ApiRequest,
  ) => Promise<{data: void}>;

  // suppressions
  getSuppressions: (
    req: ApiRequest,
  ) => Promise<{data: Array<AlertSuppression>}>;

  // receivers
  createReceiver: (
    req: {receiver: AlertReceiver} & ApiRequest,
  ) => Promise<{data: void}>;
  editReceiver: (
    req: {receiver: AlertReceiver} & ApiRequest,
  ) => Promise<{data: void}>;
  getReceivers: (req: ApiRequest) => Promise<{data: Array<AlertReceiver>}>;
  deleteReceiver: (
    req: {receiverName: string} & ApiRequest,
  ) => Promise<{data: void}>;

  // routes
  getRouteTree: (req: ApiRequest) => Promise<{data: AlertRoutingTree}>;
  editRouteTree: (
    req: {route: AlertRoutingTree} & ApiRequest,
  ) => Promise<{data: void}>;

  getMetricNames: (req: ApiRequest) => Promise<{data: Array<string>}>;
  getMetricSeries: (
    req: {name: string} & ApiRequest,
  ) => Promise<{data: Array<PrometheusLabelset>}>;

  //alertmanager global config
  getGlobalConfig: (
    req: ApiRequest,
  ) => Promise<{data: AlertManagerGlobalConfig}>;
  editGlobalConfig: (
    req: ApiRequest & {config: Partial<AlertManagerGlobalConfig>},
  ) => Promise<{data: void}>;

  // tenants
  getTenants: (req: ApiRequest) => Promise<{data: Array<string>}>;

  getAlertmanagerTenancy: (req: ApiRequest) => Promise<{data: TenancyConfig}>;
  getPrometheusTenancy: (req: ApiRequest) => Promise<{data: TenancyConfig}>;
};
