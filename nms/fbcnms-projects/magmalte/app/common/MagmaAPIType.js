/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {gateway_device, gateway_wifi_configs} from '@fbcnms/magma-api';

// TODO: remove this when wifi is fully converted to V1 API
export type WifiConfig = gateway_wifi_configs;

export type Record = gateway_device;
