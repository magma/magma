/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  AlertConfig,
  AlertManagerGlobalConfig,
  AlertReceiver,
  AlertRoutingTree,
  AlertSuppression,
  FiringAlarm,
  PrometheusLabelset,
} from './AlarmAPIType';
import type {CancelToken} from 'axios';

export type ApiRequest = {
  networkId?: string,
  cancelToken?: CancelToken,
};

export type ApiUtil = {|
  /**
   * React hook for loading data whenever a component mounts and refreshing in
   * response to prop changes. Do not use this in a class based component and
   * don't break the normal rules of react hooks.
   */
  useAlarmsApi: <TParams: {...}, TResponse>(
    func: (TParams) => Promise<TResponse>,
    params: TParams,
    cacheCounter?: string | number,
    // eslint-disable-next-line flowtype/no-weak-types
  ) => {response: ?TResponse, error: any, isLoading: boolean},

  //alerts
  viewFiringAlerts: (req: ApiRequest) => Promise<Array<FiringAlarm>>,

  //rules
  viewMatchingAlerts: (
    req: {expression: string} & ApiRequest,
  ) => Promise<Array<FiringAlarm>>,
  createAlertRule: (req: {rule: AlertConfig} & ApiRequest) => Promise<void>,
  // TODO: support renaming alerts by passing old alert name separately
  editAlertRule: (req: {rule: AlertConfig} & ApiRequest) => Promise<void>,
  getAlertRules: (req: ApiRequest) => Promise<Array<AlertConfig>>,
  deleteAlertRule: (req: {ruleName: string} & ApiRequest) => Promise<void>,

  // suppressions
  getSuppressions: (req: ApiRequest) => Promise<Array<AlertSuppression>>,

  // receivers
  createReceiver: (
    req: {receiver: AlertReceiver} & ApiRequest,
  ) => Promise<void>,
  editReceiver: (req: {receiver: AlertReceiver} & ApiRequest) => Promise<void>,
  getReceivers: (req: ApiRequest) => Promise<Array<AlertReceiver>>,
  deleteReceiver: (req: {receiverName: string} & ApiRequest) => Promise<void>,

  // routes
  getRouteTree: (req: ApiRequest) => Promise<AlertRoutingTree>,
  editRouteTree: (req: {route: AlertRoutingTree} & ApiRequest) => Promise<void>,

  // metric series
  getMetricSeries: (req: ApiRequest) => Promise<Array<PrometheusLabelset>>,

  //alertmanager global config
  getGlobalConfig: (req: ApiRequest) => Promise<AlertManagerGlobalConfig>,
  editGlobalConfig: (
    req: ApiRequest & {config: AlertManagerGlobalConfig},
  ) => Promise<void>,

  // tenants
  getTenants: (req: ApiRequest) => Promise<Array<string>>,
|};
