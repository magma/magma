/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FBCNMSRequest} from '../access';

import {Issuer} from 'openid-client';

const _clientCache = {};

export async function clientFromRequest(req: FBCNMSRequest) {
  if (!req.organization) {
    throw new Error('Must be using organization');
  }

  const {
    name,
    ssoOidcClientID,
    ssoOidcClientSecret,
    ssoOidcConfigurationURL,
  } = await req.organization();

  if (_clientCache[name]) {
    return _clientCache[name];
  }

  const issuer = await Issuer.discover(ssoOidcConfigurationURL);
  _clientCache[name] = new issuer.Client({
    client_id: ssoOidcClientID,
    client_secret: ssoOidcClientSecret,
    response_types: ['code'],
  });
  return _clientCache[name];
}
