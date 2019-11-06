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
import axios from 'axios';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../../common/useMagmaAPI';
import {MagmaAPIUrls} from '../../../common/MagmaAPI';

import type {ApiUtil} from '@fbcnms/alarms/components/ApiUrls';
import type {Match} from 'react-router-dom';

//DEPRECATED
export const MagmaAlarmAPIUrls = {
  viewFiringAlerts: (nid: string | Match) =>
    `${MagmaAPIUrls.network(nid)}/alerts`,
  alertConfig: (nid: string | Match) =>
    `${MagmaAPIUrls.network(nid)}/prometheus/alert_config`,
  updateAlertConfig: (nid: string | Match, alertName: string) =>
    `${MagmaAlarmAPIUrls.alertConfig(nid)}/${alertName}`,
  bulkAlertConfig: (nid: string | Match) =>
    `${MagmaAlarmAPIUrls.alertConfig(nid)}/bulk`,
  receiverConfig: (nid: string | Match) =>
    `${MagmaAPIUrls.network(nid)}/prometheus/alert_receiver`,
  // get count of matching metrics
  viewMatchingAlerts: (nid: string | Match, alertName: string) =>
    `${MagmaAPIUrls.network(nid)}/matching_alerts/${alertName}`,
  receiverUpdate: (nid: string | Match, receiverName: string) =>
    `${MagmaAlarmAPIUrls.receiverConfig(nid)}/${receiverName}`,
  routeConfig: (nid: string | Match) =>
    `${MagmaAlarmAPIUrls.receiverConfig(nid)}/route`,
  viewSilences: (nid: string | Match) =>
    `${MagmaAPIUrls.network(nid)}/silences`, //TODO
  viewRoutes: (nid: string | Match) => `${MagmaAPIUrls.network(nid)}/routes`, //TODO
  viewReceivers: (nid: string | Match) =>
    `${MagmaAPIUrls.network(nid)}/receivers`, //TODO
};

export const MagmaAlarmsApiUtil: ApiUtil = {
  useAlarmsApi: <TParams: {...}, TResponse>(
    func: TParams => Promise<TResponse>,
    params: TParams,
    cacheCounter?: string | number,
  ) => {
    return useMagmaAPI(func, params, undefined, cacheCounter);
  },
  viewFiringAlerts: async ({networkId}) => {
    const alerts = await MagmaV1API.getNetworksByNetworkIdAlerts({
      networkId: nullthrows(networkId),
    });
    return alerts;
  },
  viewMatchingAlerts: async ({networkId, expression}) => {
    // TODO: switch to correct MagmaV1API
    const response = await axios.get(
      `${MagmaAPIUrls.network(
        nullthrows(networkId),
      )}/matching_alerts/${expression}`,
    );
    return response.data;
  },
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
};
