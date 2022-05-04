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

import type {ExpressRequest, ExpressResponse} from 'express';

const express = require('express');
const proxy = require('express-http-proxy');

const {LOGGER_HOST} = require('../config');

const router: express.Router<
  ExpressRequest,
  ExpressResponse,
> = express.Router();

router.use('/', proxy(LOGGER_HOST));

module.exports = router;
