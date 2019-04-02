/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {access} from '@fbcnms/auth/access';
import {AccessRoles} from '@fbcnms/auth/roles';
import axios from 'axios';
import express from 'express';
import https from 'https';

import {apiCredentials, API_HOST} from '../config';

const httpsAgent = new https.Agent({
  cert: apiCredentials().cert,
  key: apiCredentials().key,
  rejectUnauthorized: false,
});

const router = express.Router();

router.post('/create', access(AccessRoles.SUPERUSER), (req, res) => {
  const {name} = req.body;

  axios
    .post(
      apiUrl('/magma/networks'),
      {name},
      {
        httpsAgent,
        params: {
          requested_id: name,
          new_workflow_flag: false,
        },
      },
    )
    .then(resp =>
      res
        .status(200)
        .send({
          success: true,
          apiResponse: resp.data,
        })
        .end(),
    )
    .catch(e =>
      res
        .status(200)
        .send({
          success: false,
          message: e.response?.data.message || e.toString(),
          apiResponse: e.response?.data,
        })
        .end(),
    );
});

const apiUrl = path =>
  !/^https?\:\/\//.test(API_HOST)
    ? `https://${API_HOST}${path}`
    : `${API_HOST}${path}`;

export default router;
