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

export async function injectOrganizationParams<T: {[string]: any}>(
  req: FBCNMSMiddleWareRequest,
  params: T,
): Promise<T & {organization?: string}> {
  if (req.organization) {
    const organization = await req.organization();
    return {
      ...params,
      organization: organization.name,
    };
  }
  return params;
}
