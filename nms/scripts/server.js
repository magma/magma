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

if (!process.env.NODE_ENV) {
  process.env.BABEL_ENV = 'development';
  process.env.NODE_ENV = 'development';
} else {
  process.env.BABEL_ENV = process.env.NODE_ENV;
}

import app from '../server/app';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import logging from '../shared/logging';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {runMigrations} from './runMigrations';

const logger = logging.getLogger(module);
const port = parseInt(process.env.PORT || 80);

(async function main() {
  await runMigrations();
  app.listen(port, '', err => {
    if (err) {
      logger.error(err.toString());
    }
    if (process.env.NODE_ENV === 'development') {
      logger.info(`Development server started on port ${port}`);
    } else {
      logger.info(`Production server started on port ${port}`);
    }
  });
})().catch(error => {
  logger.error(error);
});
