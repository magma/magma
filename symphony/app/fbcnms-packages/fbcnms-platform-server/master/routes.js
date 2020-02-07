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
import type {FeatureID} from '@fbcnms/types/features';

import MagmaV1API from '../magma';
import asyncHandler from '@fbcnms/util/asyncHandler';
import express from 'express';
import featureConfigs from '../features';

import {FeatureFlag, Organization} from '@fbcnms/sequelize-models';
import {User} from '@fbcnms/sequelize-models';
import {getPropsToUpdate} from '@fbcnms/auth/util';

const logger = require('@fbcnms/logging').getLogger(module);

const router = express.Router();

router.get(
  '/organization/async',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organizations = await Organization.findAll();
    res.status(200).send({organizations});
  }),
);

router.get(
  '/organization/async/:name',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await Organization.findOne({
      where: {name: req.params.name},
    });
    res.status(200).send({organization});
  }),
);

const configFromFeatureFlag = flag => ({
  id: flag.id,
  enabled: flag.enabled,
});
router.get(
  '/feature/async',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    // $FlowFixMe: results needs to be typed correctly
    const results: {[string]: any} = {...featureConfigs};
    Object.keys(results).forEach(id => (results[id].config = {}));
    const featureFlags = await FeatureFlag.findAll();
    featureFlags.forEach(flag => {
      if (!results[flag.featureId]) {
        logger.error(
          'feature config is missing for featureId: ' + flag.featureId,
        );
      } else {
        results[flag.featureId].config[
          flag.organization
        ] = configFromFeatureFlag(flag);
      }
    });
    res.status(200).send(Object.values(results));
  }),
);

router.post(
  '/feature/async/:featureId',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    // $FlowFixMe: Ensure it's a FeatureID
    const featureId: FeatureID = req.params.featureId;
    // $FlowFixMe: results needs to be typed correctly
    const results: {[string]: any} = {...featureConfigs};
    results.config = {};
    const {toUpdate, toDelete, toCreate} = req.body;
    const featureFlags = await FeatureFlag.findAll({where: {featureId}});
    await Promise.all(
      featureFlags.map(async flag => {
        if (toUpdate[flag.id]) {
          const newFlag = await flag.update({
            enabled: toUpdate[flag.id].enabled,
          });
          results.config[flag.organization] = configFromFeatureFlag(newFlag);
        } else if (toDelete[flag.id] !== undefined) {
          await FeatureFlag.destroy({where: {id: flag.id}});
        }
      }),
    );

    await Promise.all(
      toCreate.map(async data => {
        const flag = await FeatureFlag.create({
          featureId,
          organization: data.organization,
          enabled: data.enabled,
        });

        results.config[flag.organization] = configFromFeatureFlag(flag);
      }),
    );

    res.status(200).send(results);
  }),
);

router.post(
  '/organization/async',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await Organization.create({
      name: req.body.name,
      networkIDs: req.body.networkIDs,
      customDomains: req.body.customDomains,
      tabs: req.body.tabs,
      csvCharset: '',
      ssoCert: '',
      ssoEntrypoint: '',
      ssoIssuer: '',
    });
    res.status(200).send({organization});
  }),
);

router.put(
  '/organization/async/:name',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await Organization.findOne({
      where: {name: req.params.name},
    });
    if (!organization) {
      return res.status(404).send({error: 'Organization does not exist'});
    }
    const updated = await organization.update(req.body);
    res.status(200).send({organization: updated});
  }),
);

const USER_PROPS = ['email', 'networkIDs', 'password', 'role', 'organization'];

router.post(
  '/organization/async/:name/add_user',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await Organization.findOne({
      where: {name: req.params.name},
    });
    if (!organization) {
      return res.status(404).send({error: 'Organization does not exist'});
    }

    try {
      let props = {organization: req.params.name, ...req.body};
      props = await getPropsToUpdate(USER_PROPS, props, async params => ({
        ...params,
        organization: req.params.name,
      }));

      // this happens when the user is being added to an organization that
      // uses SSO for login, give it a random password
      if (props.password === undefined) {
        const organization = await Organization.findOne({
          where: {name: req.params.name},
        });
        if (organization && organization.ssoEntrypoint) {
          props.password = Math.random().toString(36);
        }
      }

      const user = await User.create(props);
      res.status(200).send({user});
    } catch (error) {
      res.status(400).send({error: error.toString()});
    }
  }),
);

router.delete(
  '/organization/async/:id',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    await Organization.destroy({
      where: {id: req.params.id},
      individualHooks: true,
    });
    res.status(200).send({success: true});
  }),
);

router.get(
  '/networks/async',
  asyncHandler(async (_: FBCNMSRequest, res) => {
    const networks = await MagmaV1API.getNetworks();
    res.status(200).send(networks);
  }),
);

export default router;
