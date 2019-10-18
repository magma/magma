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
import type {Record, WifiConfig} from './MagmaAPIType';
import type {magmad_gateway_configs} from '@fbcnms/magma-api';

import axios from 'axios';
import url from 'url';
import {get} from 'lodash';

function dedupeNetworkId(networkIdOrMatch: string | Match): string {
  return get(networkIdOrMatch, 'params.networkId', networkIdOrMatch);
}

export const MagmaAPIUrls = {
  networks: () => '/nms/apicontroller/magma/networks',
  network: (networkIdOrMatch: string | Match) =>
    `/nms/apicontroller/magma/networks/${dedupeNetworkId(networkIdOrMatch)}`,
  networkConfigsForType: (networkIdOrMatch: string | Match, type: 'wifi') =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/configs/${type}`,
  networkPolicyRules: (networkIdOrMatch: string | Match) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/policies/rules`,
  networkPolicyRule: (networkIdOrMatch: string | Match, ruleId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/policies/rules/${ruleId}`,
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
    type: 'wifi' | 'devmand',
  ) =>
    `${MagmaAPIUrls.network(
      networkIdOrMatch,
    )}/gateways/${gatewayId}/configs/${type}`,
  gatewayStatus: (networkIdOrMatch: string | Match, gatewayId: string) =>
    `${MagmaAPIUrls.network(networkIdOrMatch)}/gateways/${gatewayId}/status`,
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
  type: 'wifi' | 'devmand',
  configs: magmad_gateway_configs,
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
