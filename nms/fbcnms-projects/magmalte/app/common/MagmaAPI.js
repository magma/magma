/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Match} from 'react-router-dom';
import type {WifiConfig, MagmadConfig, Record} from './MagmaAPIType';

import axios from 'axios';
import {get, map} from 'lodash';
import url from 'url';

import type {CheckindGateway, NetworkUpgradeTier} from './MagmaAPIType';

function dedupeNetworkId(networkIdOrMatch: string | Match): string {
  return get(networkIdOrMatch, 'params.networkId', networkIdOrMatch);
}

export const MagmaAPIUrls = {
  networks: () => '/nms/apicontroller/magma/networks',
  network: (networkIdOrMatch: string | Match) =>
    `/nms/apicontroller/magma/networks/${dedupeNetworkId(networkIdOrMatch)}`,
  networkConfigsForType: (
    networkIdOrMatch: string | Match,
    type: 'wifi' | 'cellular',
  ) => `${MagmaAPIUrls.network(networkIdOrMatch)}/configs/${type}`,
  networkPolicyRules: (networkIdOrMatch: string | Match) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/policies/rules`,
  networkPolicyRule: (networkIdOrMatch: string | Match, ruleId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/policies/rules/${ruleId}`,
  enodeb: (networkIdOrMatch: string | Match, enodebId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/configs/enodeb/${enodebId}`,
  enodebs: (networkIdOrMatch: string | Match) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/configs/enodeb`,
  gateway: (networkIdOrMatch: string | Match, gatewayId: string) => {
    return `${MagmaAPIUrls.network(networkIdOrMatch)}/gateways/${gatewayId}`;
  },
  gateways: (networkIdOrMatch: string | Match, viewFull: boolean = false) => {
    const params = viewFull ? '?view=full' : '';
    return `${MagmaAPIUrls.network(networkIdOrMatch)}/gateways${params}`;
  },
  gatewaysSingle: (networkIdOrMatch: string | Match, gatewayId: string) =>
    `${MagmaAPIUrls.network(
      networkIdOrMatch,
    )}/gateways?view=full&gateway_ids[0]=${gatewayId}`,
  gatewayConfigs: (networkIdOrMatch: string | Match, gatewayId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/gateways/${gatewayId}/configs`,
  gatewayConfigsForType: (
    networkIdOrMatch: string | Match,
    gatewayId: string,
    type: 'wifi' | 'cellular',
  ) =>
    `${MagmaAPIUrls.network(
      networkIdOrMatch,
    )}/gateways/${gatewayId}/configs/${type}`,
  prometheusQueryRange: (networkIdOrMatch: string | Match) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/prometheus/query_range`,
  graphiteQuery: (networkIdOrMatch: string | Match) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/graphite/query`,
  networkTiers: (networkIdOrMatch: string | Match) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/tiers`,
  networkTier: (networkIdOrMatch: string | Match, tierId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/tiers/${tierId}`,
  subscribers: (networkIdOrMatch: string | Match) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/subscribers`,
  subscriber: (networkIdOrMatch: string | Match, subscriberId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/subscribers/${subscriberId}`,
  upgradeChannel: (channel: string) =>
    `/nms/apicontroller/magma/channels/${channel}`,
  command: (
    networkIdOrMatch: string | Match,
    gatewayId: string,
    commandName: string,
  ) =>
    `${MagmaAPIUrls.network(
      networkIdOrMatch,
    )}/gateways/${gatewayId}/command/${commandName}`,
};

export async function fetchAllNetworkUpgradeTiers(
  networkId: string,
): Promise<Array<NetworkUpgradeTier>> {
  const tierIdsResponse = await axios.get(MagmaAPIUrls.networkTiers(networkId));
  const tierIds = tierIdsResponse.data;
  const tierResponses = await axios.all(
    map(tierIds, tierId =>
      axios.get(MagmaAPIUrls.networkTier(networkId, tierId)),
    ),
  );
  return map(tierResponses, tierResponse => tierResponse.data);
}

export async function fetchNetworkUpgradeTier(
  networkId: string,
  tierId: string,
): Promise<NetworkUpgradeTier> {
  const tierResponse = await axios.get(
    MagmaAPIUrls.networkTier(networkId, tierId),
  );
  return tierResponse.data;
}

export async function fetchAllGateways(
  networkId: string,
): Promise<Array<CheckindGateway>> {
  const response = await axios.get(MagmaAPIUrls.gateways(networkId, true));
  const gateways: Array<CheckindGateway> = response.data;
  return gateways.filter(gateway => gateway.record !== null);
}

export async function fetchDevice(match: Match, id: string): {[string]: any} {
  const response = await axios.get(MagmaAPIUrls.gatewaysSingle(match, id));
  return response.data[0];
}

export async function createDevice(
  id: string,
  data: {
    ...Record,
    key: {key_type: string, key?: string},
  },
  type: 'wifi' | 'cellular',
  configs: MagmadConfig,
  extraConfigs: WifiConfig | {[string]: mixed},
  match: Match,
): {[string]: any} {
  const uri = url.format({
    pathname: MagmaAPIUrls.gateways(match),
    query: {new_workflow_flag: true, requested_id: id},
  });

  // creating a device in Magma requires two steps:
  // 1st Step: creating the device object itself
  await axios.post(uri, data);

  // 2nd Step: creating the config objects
  await axios.all([
    axios.post(MagmaAPIUrls.gatewayConfigs(match, id), configs),
    axios.post(
      MagmaAPIUrls.gatewayConfigsForType(match, id, type),
      extraConfigs,
    ),
  ]);

  return await fetchDevice(match, id);
}
