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
  cellular: 'cellular',
  wifi: 'wifi',
  third_party: 'third_party',
  wac: 'wac',
  rhino: 'rhino',
  lte: 'lte',
};

export const CELLULAR = Networks.cellular;
export const WIFI = Networks.wifi;
export const THIRD_PARTY = Networks.third_party;
export const WAC = Networks.wac;
export const RHINO = Networks.rhino;
export const LTE = Networks.lte;

export const AllNetworkTypes: NetworkType[] = Object.keys(Networks);

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
    return 'wifi';
  }

  return null;
}
