/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import MagmaAPIBindings from '@fbcnms/magma-api';
import axios from 'axios';
import https from 'https';
import nullthrows from '@fbcnms/util/nullthrows';
import {API_HOST, apiCredentials} from '../config';

const httpsAgent = new https.Agent({
  cert: apiCredentials().cert,
  key: apiCredentials().key,
  rejectUnauthorized: false,
});

function apiUrl(): string {
  return !/^https?\:\/\//.test(nullthrows(API_HOST))
    ? `https://${nullthrows(API_HOST)}/magma/v1`
    : `${nullthrows(API_HOST)}/magma/v1`;
}

export default class NodeClient extends MagmaAPIBindings {
  static async request(
    path: string,
    method: 'POST' | 'GET' | 'PUT' | 'DELETE' | 'OPTIONS' | 'HEAD' | 'PATCH',
    query: {[string]: mixed},
    // eslint-disable-next-line flowtype/no-weak-types
    body?: any,
  ) {
    const response = await axios({
      baseURL: apiUrl(),
      url: path,
      method: (method: string),
      params: query,
      data: body,
      headers: {'content-type': 'application/json'},
      httpsAgent,
    });

    return response.data;
  }
}
