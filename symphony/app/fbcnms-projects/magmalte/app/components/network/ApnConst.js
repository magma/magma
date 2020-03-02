/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

const ONE_MEGABYTE = 1000000;
export const DEFAULT_DATA_PLAN_ID = 'default';
export const BITRATE_MULTIPLIER = ONE_MEGABYTE;
export const DATA_PLAN_UNLIMITED_RATES = {
  max_bandwidth_dl: 200 * ONE_MEGABYTE,
  max_bandwidth_ul: 100 * ONE_MEGABYTE,
};
