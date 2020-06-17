/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

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
import {CWF, FEG, LTE, SYMPHONY} from '@fbcnms/types/network';
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
    network_services: ['policy_enforcement'],
    relay_enabled: false,
    sub_profiles: {},
    tac: 1,
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

// A placeholder due to bug in serialization
const NETWORK_FEATURES = {
  features: {
    features: {
      placeholder: 'true',
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
    const commonField = {
      name,
      description,
      id: networkID,
      ...NETWORK_FEATURES,
    };

    let resp;
    try {
      if (data.networkType === LTE) {
        resp = await MagmaV1API.postLte({
          lteNetwork: {
            ...commonField,
            cellular: DEFAULT_CELLULAR_CONFIG,
            dns: DEFAULT_DNS_CONFIG,
          },
        });
      } else if (data.networkType === CWF) {
        resp = await MagmaV1API.postCwf({
          cwfNetwork: {
            ...commonField,
            dns: DEFAULT_DNS_CONFIG,
            federation: {feg_network_id: data.fegNetworkID},
            carrier_wifi: {
              aaa_server: {
                accounting_enabled: true,
                create_session_on_auth: true,
                idle_session_timeout_ms: 500000,
              },
              default_rule_id: '',
              eap_aka: {},
              network_services: ['policy_enforcement', 'dpi'],
            },
          },
        });
      } else if (data.networkType === FEG) {
        resp = await MagmaV1API.postFeg({
          fegNetwork: {
            ...commonField,
            dns: DEFAULT_DNS_CONFIG,
            federation: {
              aaa_server: {},
              csfb: {},
              eap_aka: {},
              gx: {},
              gy: {},
              health: {},
              hss: {},
              s6a: {},
              served_network_ids: data.servedNetworkIDs.split(','),
              swx: {},
            },
          },
        });
      } else if (data.networkType === SYMPHONY) {
        resp = await MagmaV1API.postSymphony({
          symphonyNetwork: {
            ...commonField,
          },
        });
      } else {
        await MagmaV1API.postNetworks({
          network: {
            ...commonField,
            type: data.networkType,
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

export default router;
