/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import https from 'https';
import nullthrows from '@fbcnms/util/nullthrows';
import {API_HOST, apiCredentials} from '../config';

export const httpsAgent = new https.Agent({
  cert: apiCredentials().cert,
  key: apiCredentials().key,
  rejectUnauthorized: false,
});

export function apiUrl(path: string): string {
  return !/^https?\:\/\//.test(nullthrows(API_HOST))
    ? `https://${nullthrows(API_HOST)}${path}`
    : `${nullthrows(API_HOST)}${path}`;
}
