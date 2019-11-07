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
  AlertReceiver,
  AlertRoutingTree,
  AlertSuppression,
  FiringAlarm,
} from './AlarmAPIType';

export type ApiRequest = {
  networkId?: string,
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
  getReceivers: (req: ApiRequest) => Promise<Array<AlertReceiver>>,

  // routes
  getRoutes: (req: ApiRequest) => Promise<Array<AlertRoutingTree>>,
|};
