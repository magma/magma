/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 *
 */

import React from 'react';
import type {ApiUtil} from './AlarmsApi';
import type {Labels} from './AlarmAPIType';
import type {RuleInterfaceMap} from './rules/RuleInterface';

export type AlarmContext = {|
  apiUtil: ApiUtil,
  thresholdEditorEnabled?: boolean,
  filterLabels?: (labels: Labels) => Labels,
  ruleMap: RuleInterfaceMap<*>,
|};

const emptyApiUtil = {
  useAlarmsApi: () => ({
    response: null,
    error: new Error('not implemented'),
    isLoading: false,
  }),
  viewFiringAlerts: (..._) => Promise.reject('not implemented'),
  viewMatchingAlerts: (..._) => Promise.reject('not implemented'),
  createAlertRule: (..._) => Promise.reject('not implemented'),
  editAlertRule: (..._) => Promise.reject('not implemented'),
  getAlertRules: (..._) => Promise.reject('not implemented'),
  deleteAlertRule: (..._) => Promise.reject('not implemented'),
  getSuppressions: (..._) => Promise.reject('not implemented'),
  createReceiver: (..._) => Promise.reject('not implemented'),
  editReceiver: (..._) => Promise.reject('not implemented'),
  getReceivers: (..._) => Promise.reject('not implemented'),
  deleteReceiver: (..._) => Promise.reject('not implemented'),
  getRouteTree: (..._) => Promise.reject('not implemented'),
  editRouteTree: (..._) => Promise.reject('not implemented'),
  getMetricSeries: (..._) => Promise.reject('not implemented'),
};

const context = React.createContext<AlarmContext>({
  apiUtil: emptyApiUtil,
  thresholdEditorEnabled: false,
  filterLabels: x => x,
  ruleMap: {},
});

export function useAlarmContext() {
  return React.useContext(context);
}

export default context;
