/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FBCNMSRequest} from '@fbcnms/auth/access';

import asyncHandler from '@fbcnms/util/asyncHandler';
import express from 'express';
import {Organization, jsonArrayContains} from '@fbcnms/sequelize-models';
import {masterOrgMiddleware} from '@fbcnms/platform-server/master/middleware';
import {triggerActionsAlert} from '../graphgrpc/magmaalert';

const router = express.Router();

router.post(
  '/magma',
  masterOrgMiddleware,
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const {status, alerts} = req.body;

    const error = (message: string) =>
      res.status(200).send({success: false, message}).end();

    if (status !== 'firing') {
      error('"firing" webhooks are only supported');
      return;
    }

    await Promise.all(
      alerts.map(async alert => {
        if (alert.status !== 'firing') {
          return;
        }
        const {alertname, networkID} = alert.labels;

        // Get all orgs that have this networkID
        const networkOrgs = await Organization.findAll({
          where: jsonArrayContains('networkIDs', networkID),
        });

        networkOrgs.map(networkOrg =>
          triggerActionsAlert({
            tenantID: networkOrg.name,
            alertname,
            networkID,
            labels: alert.labels,
          }),
        );
      }),
    );

    res.status(200).send({success: true, message: 'ok'}).end();
  }),
);

export default router;
