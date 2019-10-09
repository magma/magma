/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

// TODO
import type {FBCNMSRequest} from '@fbcnms/auth/access';
import type {network_cellular_configs} from '@fbcnms/magmalte/app/common/__generated__/MagmaAPIBindings';

import asyncHandler from '@fbcnms/util/asyncHandler';
import axios from 'axios';
import express from 'express';

import {AccessRoles} from '@fbcnms/auth/roles';
import {CELLULAR} from '@fbcnms/types/network';
import {access} from '@fbcnms/auth/access';
import {apiUrl, httpsAgent} from '../magma';

const logger = require('@fbcnms/logging').getLogger(module);

const router = express.Router();

const DEFAULT_CELLULAR_CONFIG: network_cellular_configs = {
  epc: {
    cloud_subscriberdb_enabled: false,
    default_rule_id: '',
    mcc: '001',
    mnc: '01',
    tac: 1,
    lte_auth_amf: 'gAA=',
    lte_auth_op: 'EREREREREREREREREREREQ==',
    relay_enabled: false,
    sub_profiles: {},
  },
  ran: {
    bandwidth_mhz: 20,
    earfcndl: 44590,
    special_subframe_pattern: 7,
    subframe_assignment: 2,
    ul_dl_ratio: 1,
    tdd_config: {
      earfcndl: 44590,
      special_subframe_pattern: 7,
      subframe_assignment: 2,
    },
  },
  non_eps_service: null,
};

router.post(
  '/create',
  access(AccessRoles.SUPERUSER),
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const {networkID, data} = req.body;

    let resp;
    try {
      // Create network
      resp = await axios.post(apiUrl('/magma/networks'), data, {
        httpsAgent,
        params: {
          requested_id: networkID,
          new_workflow_flag: false,
        },
      });

      if (data.features.networkType === CELLULAR) {
        // Create default cellular config
        await axios.post(
          apiUrl(`/magma/networks/${networkID}/configs/cellular`),
          DEFAULT_CELLULAR_CONFIG,
          {httpsAgent},
        );
      }

      // Add network to organization
      const organization = await req.organization();
      const networkIDs = [...organization.networkIDs, networkID];
      await organization.update({networkIDs});
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
  }),
);

router.put(
  '/update',
  access(AccessRoles.SUPERUSER),
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const {networkID, data} = req.body;

    let resp;
    try {
      // Create network
      resp = await axios.put(apiUrl(`/magma/networks/${networkID}`), data, {
        httpsAgent,
        params: {
          requested_id: networkID,
          new_workflow_flag: false,
        },
      });
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
  }),
);

export default router;
