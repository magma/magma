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

export async function createGraphUser(tenant: string, email: string) {
  const user = caller(
    `${process.env.GRAPH_HOST || 'graph'}:443`,
    path.resolve(__dirname, 'graph.proto'),
    'UserService',
  );
  await user.Create({tenant, id: email}).catch(err => console.error(err));
}

export async function deleteGraphUser(tenant: string, email: string) {
  const user = caller(
    `${process.env.GRAPH_HOST || 'graph'}:443`,
    path.resolve(__dirname, 'graph.proto'),
    'UserService',
  );
  await user.Delete({tenant, id: email}).catch(err => console.error(err));
}
