/**
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

import fs from 'fs';
import {getValidLogLevel} from '../shared/logging';

export const DEV_MODE = process.env.NODE_ENV !== 'production';
export const LOG_FORMAT = DEV_MODE ? 'shell' : 'json';
export const LOG_LEVEL = getValidLogLevel(process.env.LOG_LEVEL);
export const LOGGER_HOST = process.env.LOGGER_HOST || 'fluentd:9880';
export const API_HOST = process.env.API_HOST || 'magma_test.local';

let _cachedApiCredentials: {
  cert: string | Buffer | undefined;
  key: string | Buffer | undefined;
} | null = null;
export function apiCredentials() {
  if (_cachedApiCredentials) {
    return _cachedApiCredentials;
  }

  let cert: string | Buffer | undefined = process.env.API_CERT;
  if (process.env.API_CERT_FILENAME) {
    try {
      cert = fs.readFileSync(process.env.API_CERT_FILENAME);
    } catch (e) {
      console.warn('cannot read cert file', e);
    }
  }

  let key: string | Buffer | undefined = process.env.API_PRIVATE_KEY;
  if (process.env.API_PRIVATE_KEY_FILENAME) {
    try {
      key = fs.readFileSync(process.env.API_PRIVATE_KEY_FILENAME);
    } catch (e) {
      console.warn('cannot read key file', e);
    }
  }

  return (_cachedApiCredentials = {
    cert,
    key,
  });
}
