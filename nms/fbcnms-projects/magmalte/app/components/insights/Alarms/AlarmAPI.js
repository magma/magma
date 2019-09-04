/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {MagmaAPIUrls} from '../../../common/MagmaAPI';
import type {Match} from 'react-router-dom';

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
  receiverUpdate: (nid: string | Match, receiverName: string) =>
    `${MagmaAlarmAPIUrls.receiverConfig(nid)}/${receiverName}`,
  routeConfig: (nid: string | Match) =>
    `${MagmaAlarmAPIUrls.receiverConfig(nid)}/route`,
};
