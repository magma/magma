/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {enodeb_configuration} from '@fbcnms/magma-api';

export const EnodebDeviceClass: {
  [string]: $PropertyType<enodeb_configuration, 'device_class'>,
} = Object.freeze({
  BAICELLS_NOVA_233_2_OD_FDD: 'Baicells Nova-233 G2 OD FDD',
  BAICELLS_NOVA_243_OD_TDD: 'Baicells Nova-243 OD TDD',
  BAICELLS_ID: 'Baicells ID TDD/FDD',
  NURAN_CAVIUM_OC_LTE: 'NuRAN Cavium OC-LTE',
});

export const EnodebBandwidthOption: {
  [string]: $NonMaybeType<$PropertyType<enodeb_configuration, 'bandwidth_mhz'>>,
} = Object.freeze({
  '3': 3,
  '5': 5,
  '10': 10,
  '15': 15,
  '20': 20,
});
