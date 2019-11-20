/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import caller from 'grpc-caller';
import path from 'path';

export async function createGraphTenant(name: string) {
  console.log(`Creating graph tenant: name=${name}`);
  const tenant = caller(
    `${process.env.GRAPH_HOST || 'graph'}:443`,
    path.resolve(__dirname, 'graph.proto'),
    'TenantService',
  );
  await tenant.Create({value: name}).catch(err => console.error(err));
}

export async function deleteGraphTenant(name: string) {
  console.log(`Getting graph tenant: name=${name}`);
  const tenant = caller(
    `${process.env.GRAPH_HOST || 'graph'}:443`,
    path.resolve(__dirname, 'graph.proto'),
    'TenantService',
  );
  const gt = await tenant.Get({value: name}).catch(err => console.error(err));
  console.log(`Deleting graph tenant: id=${gt.id}, name=${name}`);
  await tenant.Delete({value: gt.id}).catch(err => console.error(err));
}
