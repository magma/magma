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
import type {
  network_cellular_configs,
  network_dns_config,
  tier,
} from '@fbcnms/magma-api';

import asyncHandler from '@fbcnms/util/asyncHandler';
import express from 'express';

import MagmaV1API from '../magma';
import {AccessRoles} from '@fbcnms/auth/roles';
import {CELLULAR} from '@fbcnms/types/network';
import {access} from '@fbcnms/auth/access';

const logger = require('@fbcnms/logging').getLogger(module);

const router = express.Router();

const DEFAULT_CELLULAR_CONFIG: network_cellular_configs = {
  epc: {
    cloud_subscriberdb_enabled: false,
    default_rule_id: 'default_rule_1',
    lte_auth_amf: 'gAA=',
    lte_auth_op: 'EREREREREREREREREREREQ==',
    mcc: '001',
    mnc: '01',
    network_services: ['metering', 'dpi', 'policy_enforcement'],
    relay_enabled: false,
    sub_profiles: {},
    tac: 1,
  },
  features: {
    // A placeholder due to bug in serialization
    placeholder: 'true',
  },
  feg_network_id: '',
  ran: {
    bandwidth_mhz: 20,
    // TODO: Add option in UI for either fdd or tdd
    // plus config values
    // fdd_config: {
    //   earfcndl: 44590,
    //   earfcnul: 18000,
    // },
    tdd_config: {
      earfcndl: 44590,
      special_subframe_pattern: 7,
      subframe_assignment: 2,
    },
  },
};

const DEFAULT_DNS_CONFIG: network_dns_config = {
  enable_caching: false,
  local_ttl: 0,
  records: [],
};

const DEFAULT_UPGRADE_TIER: tier = {
  gateways: [],
  id: 'default',
  images: [],
  name: 'Default Tier',
  version: '0.0.0-0',
};

router.post(
  '/create',
  access(AccessRoles.SUPERUSER),
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const {networkID, data} = req.body;
    const {name, description} = data;

    let resp;
    try {
      if (data.features.networkType === CELLULAR) {
        resp = await MagmaV1API.postLte({
          lteNetwork: {
            cellular: DEFAULT_CELLULAR_CONFIG,
            dns: DEFAULT_DNS_CONFIG,
            id: networkID,
            name,
            description,
          },
        });
      } else {
        await MagmaV1API.postNetworks({
          network: {
            name,
            description,
            id: networkID,
            type: data.features.networkType,
            dns: DEFAULT_DNS_CONFIG,
          },
        });
      }

      MagmaV1API.postNetworksByNetworkIdTiers({
        networkId: networkID,
        tier: DEFAULT_UPGRADE_TIER,
      });

      // Add network to organization
      if (req.organization) {
        const organization = await req.organization();
        const networkIDs = [...organization.networkIDs, networkID];
        await organization.update({networkIDs});
      }
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
        apiResponse: resp,
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
      // Update network
      if (data.features.networkType === CELLULAR) {
        resp = await MagmaV1API.putLteByNetworkId({
          networkId: networkID,
          lteNetwork: data,
        });
      } else {
        resp = await MagmaV1API.putNetworksByNetworkId({
          networkId: networkID,
          network: data,
        });
      }
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
        apiResponse: resp,
      })
      .end();
  }),
);

export default router;
