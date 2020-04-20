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
 * @flow strict-local
 * @format
 */

const {LOG_FORMAT, LOG_LEVEL} = require('@fbcnms/platform-server/config');

// This must be done before any module imports to configure
// logging correctly
const logging = require('@fbcnms/logging');
logging.configure({
  LOG_FORMAT,
  LOG_LEVEL,
});

const {runMigrations} = require('./runMigrations');

runMigrations()
  .then(_ => console.log('Ran migrations successfully'))
  .catch(_ => console.error('Failed to run migrations'));
