/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
'use strict';

// NOTE: Regex based routes for paths that don't require logged in user access
export default [
  /^\/favicon.ico$/,
  /^\/healthz$/,
  /^\/user\/login(\?.*)?/,
  /^\/user\/onboarding(\?.*)?/,
  /^\/([a-z_-]+\/)?static\/css/,
  /^\/([a-z_-]+\/)?static\/dist\/login.js/,
  /^\/([a-z_-]+\/)?static\/dist\/onboarding.js/,
  /^\/([a-z_-]+\/)?static\/dist\/vendor.js/,
  /^\/([a-z_-]+\/)?static\/fonts/,
  /^\/([a-z_-]+\/)?static\/images/,
  /^\/([a-z_-]+\/)?user\/login(\?.*)?$/,
  /^\/([a-z_-]+\/)?user\/login\/saml$/,
  /^\/([a-z_-]+\/)?user\/login\/saml\/callback/,
  /^\/([a-z_-]+\/)?user\/logout$/,
  /^\/([a-z_-]+\/)?__webpack_hmr.js/,
];
