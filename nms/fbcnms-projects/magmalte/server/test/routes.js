/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const express = require('express');

import type {NMSRequest} from '../../scripts/server';

const router = express.Router();

router.get('/', (req: NMSRequest, res) => {
  res.status(200).end('Success');
});

module.exports = router;
