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

import bodyParser from 'body-parser';
import compression from 'compression';
import cookieParser from 'cookie-parser';
import express from 'express';
import helmet from 'helmet';
import logging from '../../shared/logging';

/**
 * General middleware that every application should use, and it should be the
 * first thing used.  These shouldn't have any side effects in the application
 * it should just introduce additional functionality
 */
// $FlowIgnore[value-as-type]
export default function appMiddleware(): Middleware {
  const router = express.Router();
  [
    helmet(),
    // parse json. Strict disabled because magma wants gateway name update
    // to be just a string (e.g. "name") which is not actually legit
    bodyParser.json({limit: '1mb', strict: false}),
    // parse application/x-www-form-urlencoded
    bodyParser.urlencoded({limit: '1mb', extended: false}),
    cookieParser(),
    compression(),
    logging.getHttpLogger(module),
  ].forEach(middleware => router.use(middleware));
  return router;
}
