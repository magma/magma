/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export const CELLULAR = 'cellular';
export const WIFI = 'wifi';
export const THIRD_PARTY = 'third_party';
export const TARAZED = 'tarazed';
export const WAC = 'wac';

export const AllNetworkTypes: NetworkType[] = [
  CELLULAR,
  WIFI,
  THIRD_PARTY,
  TARAZED,
  WAC,
];

export type NetworkType =
  | 'cellular'
  | 'wifi'
  | 'third_party'
  | 'tarazed'
  | 'wac';
