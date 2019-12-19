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

const logger = require('@fbcnms/logging').getLogger(module);

type Payload = {
  tenant: string,
  alertname: string,
  networkID: string,
  labels: {[string]: string},
};

export async function triggerMagmaAlert(payload: Payload) {
  logger.info('sending payload: %s', payload);
}
