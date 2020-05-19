/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import Router from 'express';
import {getUserGroups, getUserRole} from './utils.js';
import type {AuthorizationCheck, GroupLoadingStrategy} from '../types';

const rbacRouter = Router();

export default async function(
  authorizationCheck: AuthorizationCheck,
  groupLoadingStrategy: GroupLoadingStrategy,
) {
  rbacRouter.get('/editableworkflows', async (req, res, _) => {
    const role = getUserRole(req);
    const groups = await getUserGroups(req, groupLoadingStrategy);

    if (authorizationCheck(role, groups)) {
      res.status(200).send(true);
    } else {
      res.status(200).send(false);
    }
  });

  return rbacRouter;
}
