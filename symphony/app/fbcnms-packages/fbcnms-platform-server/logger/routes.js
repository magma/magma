/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
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
