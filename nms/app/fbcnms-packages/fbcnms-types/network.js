/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const Networks = {
  carrier_wifi_network: 'carrier_wifi_network',
  xwfm: 'xwfm',
  feg: 'feg',
  lte: 'lte',
  rhino: 'rhino',
  symphony: 'symphony',
  third_party: 'third_party', // TODO: deprecate third_party in lieu of symphony
  wifi_network: 'wifi_network',
};

export const CWF = Networks.carrier_wifi_network;
export const XWFM = Networks.xwfm;
export const FEG = Networks.feg;
export const LTE = Networks.lte;
export const RHINO = Networks.rhino;
export const SYMPHONY = Networks.symphony;
export const THIRD_PARTY = Networks.third_party;
export const WIFI = Networks.wifi_network;

export const AllNetworkTypes: NetworkType[] = Object.keys(Networks).sort();
export const V1NetworkTypes: NetworkType[] = [
  CWF,
  FEG,
  LTE,
  SYMPHONY,
  WIFI,
  XWFM,
];

export type NetworkType = $Keys<typeof Networks>;

export function coalesceNetworkType(
  networkID: string,
  networkType: ?string,
): ?NetworkType {
  if (networkType && Networks[networkType]) {
    return (networkType: any);
  }

  // backwards compatibility
  if (networkID.startsWith('mesh_')) {
    return 'wifi_network';
  }

  return null;
}
