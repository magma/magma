/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FBCNMSRequest} from '@fbcnms/auth/access';
import type {FeatureID} from '@fbcnms/types/features';

import asyncHandler from '@fbcnms/util/asyncHandler';
import express from 'express';
import {FeatureFlag} from '@fbcnms/sequelize-models';
import {
  featureConfigs,
  getEnabledFeatures,
} from '@fbcnms/platform-server/features';

const router: express.Router<FBCNMSRequest, ExpressResponse> = express.Router();

router.get(
  '/async/',
  asyncHandler(async (req: FBCNMSRequest, res: ExpressResponse) => {
    const organization = await req.organization();
    const features = await getEnabledFeatures(req, organization.name, true);
    res.status(200).send({features});
  }),
);

router.post(
  '/async/:featureId',
  asyncHandler(async (req: FBCNMSRequest, res: ExpressResponse) => {
    // $FlowFixMe: Ensure it's a FeatureID
    const featureId: FeatureID = req.params.featureId;
    if (!featureConfigs[featureId]?.publicAccess) {
      res.status(401).send({error: 'Unauthorized to change this feature flag'});
    }
    const {enabled} = req.body;
    const organization = await req.organization();
    const flag = await FeatureFlag.findOne({
      where: {
        featureId,
        organization: organization.name,
      },
    });
    if (flag) {
      await flag.update({
        enabled,
      });
    } else {
      await FeatureFlag.create({
        featureId,
        organization: organization.name,
        enabled,
      });
    }
    res.status(200).send();
  }),
);

export default router;
