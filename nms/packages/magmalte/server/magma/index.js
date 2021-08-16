/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow
 * @format
 */

import MagmaAPIBindings from '@fbcnms/magma-api';
import axios from 'axios';
import https from 'https';
import nullthrows from '@fbcnms/util/nullthrows';
import {API_HOST, apiCredentials} from '@fbcnms/platform-server/config';

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
