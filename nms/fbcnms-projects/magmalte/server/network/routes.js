/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CellularNetworkConfig} from '../../app/common/MagmaAPIType';

import axios from 'axios';
import express from 'express';
import https from 'https';
import {AccessRoles} from '@fbcnms/auth/roles';
import {access} from '@fbcnms/auth/access';

import {API_HOST, apiCredentials} from '../config';

const logger = require('@fbcnms/logging').getLogger(module);

import type {NMSRequest} from '../../scripts/server';

const httpsAgent = new https.Agent({
  cert: apiCredentials().cert,
  key: apiCredentials().key,
  rejectUnauthorized: false,
});

const router = express.Router();

const DEFAULT_CELLULAR_CONFIG: CellularNetworkConfig = {
  epc: {
    mcc: '001',
    mnc: '01',
    tac: 1,
    lte_auth_amf: 'gAA=',
    lte_auth_op: 'EREREREREREREREREREREQ==',
    sub_profiles: {},
  },
  ran: {
    bandwidth_mhz: 20,
    earfcndl: 44590,
    special_subframe_pattern: 7,
    subframe_assignment: 2,
    ul_dl_ratio: 1,
  },
  non_eps_service: null,
};

router.post(
  '/create',
  access(AccessRoles.SUPERUSER),
  async (req: NMSRequest, res) => {
    const {name} = req.body;

    let resp;
    try {
      // Create network
      resp = await axios.post(
        apiUrl('/magma/networks'),
        {name},
        {
          httpsAgent,
          params: {
            requested_id: name,
            new_workflow_flag: false,
          },
        },
      );

      // Create default cellular config
      await axios.post(
        apiUrl(`/magma/networks/${name}/configs/cellular`),
        DEFAULT_CELLULAR_CONFIG,
        {httpsAgent},
      );
    } catch (e) {
      logger.error(e, {
        response: e.response?.data,
      });
      res
        .status(200)
        .send({
          success: false,
          message: e.response?.data.message || e.toString(),
          apiResponse: e.response?.data,
        })
        .end();
      return;
    }

    res
      .status(200)
      .send({
        success: true,
        apiResponse: resp.data,
      })
      .end();
  },
);

const apiUrl = path =>
  !/^https?\:\/\//.test(API_HOST)
    ? `https://${API_HOST}${path}`
    : `${API_HOST}${path}`;

export default router;
