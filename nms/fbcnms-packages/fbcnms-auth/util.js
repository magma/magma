/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FBCNMSMiddleWareRequest} from '@fbcnms/express-middleware';
import type {UserType} from '@fbcnms/sequelize-models/models/user.js';

import querystring from 'querystring';
import {format, parse} from 'url';
import {injectOrganizationParams} from './organization';
import {User} from '@fbcnms/sequelize-models';

export function addQueryParamsToUrl(
  url: string,
  params: {[string]: any},
): string {
  const parsedUrl = parse(url, true /* parseQueryString */);
  if (params) {
    parsedUrl.search = querystring.stringify({
      ...parsedUrl.query,
      ...params,
    });
  }
  return format(parsedUrl);
}

export async function getUserFromRequest(
  req: FBCNMSMiddleWareRequest,
  email: string,
): Promise<?UserType> {
  const where = await injectOrganizationParams(req, {email});
  return await User.findOne({where});
}
