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

const fs = require('fs');
const {getValidLogLevel} = require('../../fbc_js_core/logging');

const DEV_MODE = process.env.NODE_ENV !== 'production';
const LOG_FORMAT = DEV_MODE ? 'shell' : 'json';
const LOG_LEVEL = getValidLogLevel(process.env.LOG_LEVEL);

const MAPBOX_ACCESS_TOKEN = process.env.MAPBOX_ACCESS_TOKEN || '';

const LOGGER_HOST = process.env.LOGGER_HOST || 'fluentd:9880';

// NMS specific

const API_HOST = process.env.API_HOST || 'magma_test.local';

let _cachedApiCredentials = null;
function apiCredentials() {
  if (_cachedApiCredentials) {
    return _cachedApiCredentials;
  }

  let cert = process.env.API_CERT;
  if (process.env.API_CERT_FILENAME) {
    try {
      cert = fs.readFileSync(process.env.API_CERT_FILENAME);
    } catch (e) {
      console.warn('cannot read cert file', e);
    }
  }

  let key = process.env.API_PRIVATE_KEY;
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

module.exports = {
  apiCredentials,
  API_HOST,
  DEV_MODE,
  LOG_FORMAT,
  LOG_LEVEL,
  LOGGER_HOST,
  MAPBOX_ACCESS_TOKEN,
};
