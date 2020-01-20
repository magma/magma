/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as beaver from 'beaver-logger';
import axios from 'axios';

export const Events = {
  DOCUMENTATION_LINK_CLICKED: 'documentation_link_clicked',
  SETTINGS_CLICKED: 'settings_clicked',
};

export const ServerLog = (topic: string) =>
  beaver.Logger({
    url: '/logger/' + topic,
    logLevel: beaver.LOG_LEVEL.INFO,
    flushInterval: 10 * 1000,
    // eslint-disable-next-line flowtype/no-weak-types
    transport: async ({url, method, json}): Promise<any> =>
      axios({
        method,
        url,
        data: json.events.map(e => {
          const {event, level, payload} = e;
          const {timestamp, data, user} = payload;
          return {
            event,
            level,
            ts: timestamp / 1000,
            data,
            ...user,
          };
        }),
      }),
  });

export const GeneralLogger = ServerLog('common');
