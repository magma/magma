/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MagmadConfig, Record, WifiConfig} from './MagmaAPIType';
import type {Match} from 'react-router-dom';

import axios from 'axios';
import url from 'url';
import {get, map} from 'lodash';

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
    type: 'wifi' | 'cellular' | 'tarazed' | 'devmand',
  ) =>
    `${MagmaAPIUrls.network(
      networkIdOrMatch,
    )}/gateways/${gatewayId}/configs/${type}`,
  gatewayName: (networkIdOrMatch: string | Match, gatewayId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/gateways/${gatewayId}/name`,
  gatewayStatus: (networkIdOrMatch: string | Match, gatewayId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/gateways/${gatewayId}/status`,
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
  device: (networkIdOrMatch: string | Match, deviceId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/configs/devices/${deviceId}`,
  devices: (networkIdOrMatch: string | Match, deviceId?: string) => {
    return deviceId
      ? `${MagmaAPIUrls.network(
          networkIdOrMatch,
        )}/configs/devices?requested_id=${deviceId}`
      : `${MagmaAPIUrls.network(networkIdOrMatch)}/configs/devices`;
  },
  devicesDevmandConfigs: (
    networkIdOrMatch: string | Match,
    gatewayId: string,
  ) =>
    `${MagmaAPIUrls.network(
      networkIdOrMatch,
    )}/gateways/${gatewayId}/configs/devmand`,
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

export async function fetchDevice(
  networkIdOrMatch: string | Match,
  id: string,
): {[string]: any} {
  const response = await axios.get(
    MagmaAPIUrls.gatewaysSingle(networkIdOrMatch, id),
  );
  return response.data[0];
}

export async function createDevice(
  id: string,
  data: {
    ...Record,
    key: {key_type: string, key?: string},
  },
  type: 'wifi' | 'cellular' | 'tarazed' | 'devmand',
  configs: MagmadConfig,
  extraConfigs: WifiConfig | {[string]: mixed},
  networkIdOrMatch: string | Match,
): {[string]: any} {
  const uri = url.format({
    pathname: MagmaAPIUrls.gateways(networkIdOrMatch),
    query: {new_workflow_flag: true, requested_id: id},
  });

  // creating a device in Magma requires two steps:
  // 1st Step: creating the device object itself
  await axios.post(uri, data);

  // 2nd Step: creating the config objects
  await axios.all([
    axios.post(MagmaAPIUrls.gatewayConfigs(networkIdOrMatch, id), configs),
    axios.post(
      MagmaAPIUrls.gatewayConfigsForType(networkIdOrMatch, id, type),
      extraConfigs,
    ),
  ]);

  return await fetchDevice(networkIdOrMatch, id);
}

export async function updateGatewayName(
  gatewayId: string,
  name: string,
  networkIdOrMatch: string | Match,
): Promise<void> {
  await axios.put(
    MagmaAPIUrls.gatewayName(networkIdOrMatch, gatewayId),
    JSON.stringify(`"${name}"`),
    {
      headers: {'content-type': 'application/json'},
    },
  );
}
