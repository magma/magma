/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Match} from 'react-router-dom';

import type {AlertConfig, FiringAlarm} from './AlarmAPIType';

// DEPRECATED
export type ApiUrls = {
  //alerts
  viewFiringAlerts: (nid: string | Match) => string,
  //rules
  viewMatchingAlerts: (nid: string | Match, alertName: string) => string,
  alertConfig: (nid: string | Match) => string,
  updateAlertConfig: (nid: string | Match, alertName: string) => string,
  bulkAlertConfig: (nid: string | Match) => string,
  //routes
  routeConfig: (nid: string | Match) => string,
  viewRoutes: (nid: string | Match) => string,
  //silences
  viewSilences: (nid: string | Match) => string,
  // receivers
  viewReceivers: (nid: string | Match) => string,
  receiverConfig: (nid: string | Match) => string,
  receiverUpdate: (nid: string | Match, receiverName: string) => string,
};

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
  editAlertRule: (req: {rule: AlertConfig} & ApiRequest) => Promise<void>,
  getAlertRules: (req: ApiRequest) => Promise<Array<AlertConfig>>,
  deleteAlertRule: (req: {ruleName: string} & ApiRequest) => Promise<void>,
|};
