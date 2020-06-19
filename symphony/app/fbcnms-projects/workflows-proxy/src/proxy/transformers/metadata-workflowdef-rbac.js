/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import logging from '@fbcnms/logging';

import {anythingTo, isLabeledWithGroup} from '../utils.js';
import {
  getAllWorkflowsAfter as getAllWorkflowsAfterDelegate,
  getWorkflowAfter as getWorkflowAfterDelegate,
  getWorkflowBefore as getWorkflowBeforeDelegate,
} from './metadata-workflowdef.js';
import type {AfterFun, TransformerRegistrationFun, Workflow} from '../../types';

const logger = logging.getLogger(module);

const getAllWorkflowsAfter: AfterFun = (
  tenantId,
  groups,
  req,
  respObj,
  res,
) => {
  const workflows: Array<Workflow> = anythingTo<Array<Workflow>>(respObj);
  for (
    let workflowIdx = workflows.length - 1;
    workflowIdx >= 0;
    workflowIdx--
  ) {
    const workflowdef = workflows[workflowIdx];
    if (!isLabeledWithGroup(workflowdef, groups)) {
      workflows.splice(workflowIdx, 1);
    }
  }
  getAllWorkflowsAfterDelegate(tenantId, groups, req, respObj, res);
};

export const getWorkflowAfter: AfterFun = (
  tenantId,
  groups,
  req,
  respObj,
  res,
) => {
  const workflowdef: Workflow = anythingTo<Workflow>(respObj);
  if (!isLabeledWithGroup(workflowdef, groups)) {
    // fail if workflow is outside of user's groups
    logger.error(
      `User accessing unauthorized workflow: ${workflowdef.name} for tenant: ${tenantId}`,
    );
    res.status(401).send('User unauthorized to access this endpoint');
  }

  getWorkflowAfterDelegate(tenantId, groups, req, respObj, res);
};

const registration: TransformerRegistrationFun = () => [
  {
    method: 'get',
    url: '/api/metadata/workflow',
    after: getAllWorkflowsAfter,
  },
  {
    method: 'get',
    url: '/api/metadata/workflow/:name',
    before: getWorkflowBeforeDelegate,
    after: getWorkflowAfter,
  },
];

export default registration;
