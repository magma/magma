/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const querystring = require('querystring');
import {AccessRoles} from '@fbcnms/auth/roles';
import {oidcAccessToken} from '@fbcnms/auth/oidc/middleware';
import type {ClientRequest} from 'http';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

export function accessRoleToString(role: number): string {
  if (role === AccessRoles.SUPERUSER) {
    return 'OWNER';
  }
  return 'USER';
}

const onProxyReq = (proxyReq: ClientRequest, req: FBCNMSRequest): void => {
  if (req.user.organization) {
    proxyReq.setHeader('x-auth-organization', req.user.organization);
  }
  proxyReq.setHeader('x-auth-user-email', req.user.email);
  proxyReq.setHeader('x-auth-user-role', accessRoleToString(req.user.role));

  const accessToken = oidcAccessToken(req);
  if (accessToken != null) {
    proxyReq.setHeader('authorization', 'Bearer ' + accessToken);
  }

  if (!req.body || !Object.keys(req.body).length) {
    return;
  }

  const writeBody = (body: string) => {
    proxyReq.setHeader('Content-Length', Buffer.byteLength(body).toString());
    proxyReq.write(body);
    proxyReq.end();
  };

  const contentType = proxyReq.getHeader('Content-Type');
  if (contentType.includes('application/json')) {
    writeBody(JSON.stringify(req.body));
  } else if (contentType.includes('application/x-www-form-urlencoded')) {
    writeBody(querystring.stringify(req.body));
  }
};

export default onProxyReq;
