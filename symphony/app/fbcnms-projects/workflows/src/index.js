/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import ExpressApplication from 'express';

const app = ExpressApplication();

app.get('/', (req, res) => {
  res.send('hello world');
});

app.get('/echo/:str', (req, res) => {
  res.send(req.params.str);
});

app.listen(80);
