/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import caller from 'grpc-caller';
import path from 'path';

const logger = require('@fbcnms/logging').getLogger(module);

type Payload = {
  tenantID: string,
  alertname: string,
  networkID: string,
  labels: {[string]: string},
};

export async function triggerActionsAlert(payload: Payload) {
  logger.info('sending payload: %s', payload);

  const actionsAlertService = caller(
    `${process.env.GRAPH_HOST || 'graph'}:443`,
    path.resolve(__dirname, 'graph.proto'),
    'ActionsAlertService',
  );
  await actionsAlertService.Trigger(payload).catch(err => console.error(err));
}
