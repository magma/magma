/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export type WifiConfig = {
  mesh_id?: string,
  info?: string,
  longitude?: number,
  latitude?: number,
  client_channel: string,
  is_production: boolean,
  additional_props: ?{[string]: string},
};

export type Record = {
  hardware_id: string,
};
