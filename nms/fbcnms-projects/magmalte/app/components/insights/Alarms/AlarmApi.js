/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../../common/useMagmaAPI';

import type {ApiUtil} from '@fbcnms/alarms/components/AlarmsApi';

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
  // Alert Rules
  createAlertRule: async ({networkId, rule}) => {
    const rules = await MagmaV1API.postNetworksByNetworkIdPrometheusAlertConfig(
      {
        networkId: nullthrows(networkId),
        alertConfig: rule,
      },
    );
    return rules;
  },
  editAlertRule: async ({networkId, rule}) => {
    const rules = await MagmaV1API.putNetworksByNetworkIdPrometheusAlertConfigByAlertName(
      {
        networkId: nullthrows(networkId),
        alertName: rule.alert,
        alertConfig: rule,
      },
    );
    return rules;
  },
  getAlertRules: async ({networkId}) => {
    const rules = await MagmaV1API.getNetworksByNetworkIdPrometheusAlertConfig({
      networkId: nullthrows(networkId),
    });
    return rules;
  },
  deleteAlertRule: async ({ruleName, networkId}) => {
    await MagmaV1API.deleteNetworksByNetworkIdPrometheusAlertConfig({
      networkId: nullthrows(networkId),
      alertName: ruleName,
    });
  },
  getSuppressions: () => {
    console.warn('not implemented');
    return Promise.resolve([]);
  },
  getReceivers: () => {
    console.warn('not implemented');
    return Promise.resolve([]);
  },
  getRoutes: () => {
    console.warn('not implemented');
    return Promise.resolve([]);
  },
  getMetricSeries: async ({networkId}) => {
    const series = await MagmaV1API.getNetworksByNetworkIdPrometheusSeries({
      networkId: nullthrows(networkId),
    });
    return series;
  },
};
