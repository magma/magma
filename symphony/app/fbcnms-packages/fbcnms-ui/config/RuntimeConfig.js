/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

const TEST_SUBDOMAIN = '-test';
const LOCALHOST = 'localhost';
const PHB_SUBDOMAIN = 'purpleheadband.cloud';

export function isTestEnv(): boolean {
  return (
    window.location.hostname.includes(TEST_SUBDOMAIN) ||
    window.location.hostname.includes(LOCALHOST)
  );
}

export function isPhbProdEnv(): boolean {
  return window.location.hostname.includes(PHB_SUBDOMAIN);
}
