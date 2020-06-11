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

const ConductorClient = require('conductor-client').default;

const conductorApiUrl =
  process.env.CONDUCTOR_API_URL || 'http://conductor-server:8080/api';

const conductorClient = new ConductorClient({
  baseURL: conductorApiUrl,
});

conductorClient.registerWatcher(
  'GLOBAL___js',
  (data, updater) => {
    console.log(data.taskType, data.inputData);
    updater.complete({});
  },
  {pollingIntervals: 1000, autoAck: true, maxRunner: 1},
  true,
);
