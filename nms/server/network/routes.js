/*
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

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '../auth/access';
import type {
  network_cellular_configs,
  network_dns_config,
  tier,
} from '../../generated/MagmaAPIBindings';

import asyncHandler from '../util/asyncHandler';
import express from 'express';

import MagmaV1API from '../magma';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AccessRoles} from '../../shared/roles';
// $FlowFixMe migrated to typescript
import {CWF, FEG, FEG_LTE, LTE, XWFM} from '../../shared/types/network';
import {access} from '../auth/access';
import {difference} from 'lodash';

const logger = require('../../shared/logging').getLogger(module);

const router: express.Router<FBCNMSRequest, ExpressResponse> = express.Router();

const DEFAULT_CELLULAR_CONFIG: network_cellular_configs = {
  epc: {
    cloud_subscriberdb_enabled: false,
    default_rule_id: 'default_rule_1',
    lte_auth_amf: 'gAA=',
    lte_auth_op: 'EREREREREREREREREREREQ==',
    mcc: '001',
    mnc: '01',
    network_services: ['policy_enforcement'],
    hss_relay_enabled: false,
    gx_gy_relay_enabled: false,
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
      } else if (data.networkType === FEG_LTE) {
        resp = await MagmaV1API.postFegLte({
          lteNetwork: {
            ...commonField,
            cellular: {
              ...DEFAULT_CELLULAR_CONFIG,
              feg_network_id: data.fegNetworkID,
            },
            dns: DEFAULT_DNS_CONFIG,
            federation: {feg_network_id: data.fegNetworkID},
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
      } else if (data.networkType === XWFM) {
        resp = await MagmaV1API.postCwf({
          cwfNetwork: {
            ...commonField,
            dns: DEFAULT_DNS_CONFIG,
            federation: {feg_network_id: data.fegNetworkID},
            carrier_wifi: {
              is_xwfm_variant: true,
              aaa_server: {
                accounting_enabled: true,
                create_session_on_auth: true,
                idle_session_timeout_ms: 500000,
              },
              default_rule_id: '',
              eap_aka: {},
              network_services: [],
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
      } else {
        res
          .status(400)
          .send(`Unsupported network type ${data.networkType}`)
          .end();
        return;
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

router.post(
  '/delete',
  access(AccessRoles.SUPERUSER),
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const {networkID} = req.body;

    try {
      await MagmaV1API.deleteNetworksByNetworkId({
        networkId: networkID,
      });

      // Remove network from organization
      if (req.organization) {
        const organization = await req.organization();

        const networkIDs = difference(organization.networkIDs, [networkID]);
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
        })
        .end();
      return;
    }

    res
      .status(200)
      .send({
        success: true,
      })
      .end();
  }),
);
export default router;
