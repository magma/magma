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
 */

import OrchestratorAPI from '../api/OrchestratorAPI';
import asyncHandler from '../util/asyncHandler';
import logging from '../../shared/logging';
import {AccessRoles} from '../../shared/roles';
import {AxiosError} from 'axios';
import {CWF, FEG, FEG_LTE, LTE} from '../../shared/types/network';
import {Request, Router} from 'express';
import {access} from '../auth/access';
import {difference} from 'lodash';
import type {
  NetworkCellularConfigs,
  NetworkDnsConfig,
  Tier,
} from '../../generated';

const logger = logging.getLogger(module);

const router = Router();

const DEFAULT_CELLULAR_CONFIG: NetworkCellularConfigs = {
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

const DEFAULT_DNS_CONFIG: NetworkDnsConfig = {
  enable_caching: false,
  local_ttl: 0,
  records: [],
};

const DEFAULT_UPGRADE_TIER: Tier = {
  gateways: [],
  id: 'default',
  images: [],
  name: 'Default Tier',
  version: '0.0.0-0',
};

router.post(
  '/create',
  access(AccessRoles.SUPERUSER),
  asyncHandler(
    async (
      req: Request<
        never,
        any,
        {
          networkID: string;
          data: {
            name: string;
            description: string;
            fegNetworkID?: string;
            networkType: string;
            servedNetworkIDs?: string;
          };
        }
      >,
      res,
    ) => {
      const {networkID, data} = req.body;
      const {name, description} = data;
      const allowedNetworkTypes = ['LTE', 'FEG_LTE', 'CWF', 'FEG'];

      if (!allowedNetworkTypes.includes(data.networkType?.toUpperCase())) {
        res
          .status(400)
          .send(
            `please provide a valid network type like: LTE, FEG_LTE, CWF or FEG`,
          )
          .end();
        return;
      }
      const commonField = {
        name,
        description,
        id: networkID,
        ...NETWORK_FEATURES,
      };

      let resp;
      try {
        if (data.networkType === LTE) {
          resp = await OrchestratorAPI.lteNetworks.ltePost({
            lteNetwork: {
              ...commonField,
              cellular: DEFAULT_CELLULAR_CONFIG,
              dns: DEFAULT_DNS_CONFIG,
            },
          });
        } else if (data.networkType === FEG_LTE) {
          resp = await OrchestratorAPI.federatedLTENetworks.fegLtePost({
            lteNetwork: {
              ...commonField,
              cellular: {
                ...DEFAULT_CELLULAR_CONFIG,
                feg_network_id: data.fegNetworkID,
              },
              dns: DEFAULT_DNS_CONFIG,
              federation: {feg_network_id: data.fegNetworkID!},
            },
          });
        } else if (data.networkType === CWF) {
          resp = await OrchestratorAPI.carrierWifiNetworks.cwfPost({
            cwfNetwork: {
              ...commonField,
              dns: DEFAULT_DNS_CONFIG,
              federation: {feg_network_id: data.fegNetworkID!},
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
          resp = await OrchestratorAPI.federationNetworks.fegPost({
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
                served_network_ids: data.servedNetworkIDs!.split(','),
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

        await OrchestratorAPI.upgrades.networksNetworkIdTiersPost({
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
        const error = e as AxiosError<{message: string}>;
        logger.error(error.toString(), {
          response: error.response?.data,
        });
        res
          .status(200)
          .send({
            success: false,
            message: error.response?.data.message || error.toString(),
            apiResponse: error.response?.data,
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
  ),
);

router.post(
  '/delete',
  access(AccessRoles.SUPERUSER),
  asyncHandler(async (req: Request<never, any, {networkID: string}>, res) => {
    const {networkID} = req.body;

    try {
      await OrchestratorAPI.networks.networksNetworkIdDelete({
        networkId: networkID,
      });

      // Remove network from organization
      if (req.organization) {
        const organization = await req.organization();

        const networkIDs = difference(organization.networkIDs, [networkID]);
        await organization.update({networkIDs});
      }
    } catch (e) {
      const error = e as AxiosError<{message: string}>;
      logger.error(error.toString(), {
        response: error.response?.data,
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
