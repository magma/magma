/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {anythingTo, getUserEmail} from '../utils.js';
import {
  getExecutionStatusAfter as getExecutionStatusAfterDelegate,
  getSearchAfter as getSearchAfterDelegate,
  getSearchBefore as getSearchBeforeDelegate,
  postWorkflowBefore as postWorkflowBeforeDelegate,
  removeWorkflowBefore as removeWorkflowBeforeDelegate,
  updateQuery,
} from './workflow.js';
import type {
  AfterFun,
  BeforeFun,
  StartWorkflowRequest,
  TransformerRegistrationFun,
} from '../../types';

const postWorkflowBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  // FIXME verify workflow def has proper groups attached to it
  const reqObj = anythingTo<StartWorkflowRequest>(req.body);
  // Put userEmail into correlationId
  reqObj.correlationId = getUserEmail(req);
  postWorkflowBeforeDelegate(tenantId, groups, req, res, proxyCallback);
};

const getExecutionStatusAfter: AfterFun = (
  tenantId,
  groups,
  req,
  respObj,
  res,
) => {
  // FIXME verify workflow def has proper groups attached to it
  getExecutionStatusAfterDelegate(tenantId, groups, req, respObj, res);
};

export const removeWorkflowBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  // FIXME verify workflow def has proper groups attached to it
  removeWorkflowBeforeDelegate(tenantId, groups, req, res, proxyCallback);
};

export const getSearchBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  // Prefix query with correlationId == userEmail
  // This limits the search to workflows started by current user
  const userEmail = getUserEmail(req);
  const originalQueryString = req._parsedUrl.query;
  const limitToTenant = `correlationId = \'${userEmail}\'`;
  const newQueryString = updateQuery(originalQueryString, limitToTenant);
  req._parsedUrl.query = newQueryString;
  getSearchBeforeDelegate(tenantId, groups, req, res, proxyCallback);
};

const registration: TransformerRegistrationFun = function() {
  return [
    {
      method: 'get',
      url: '/api/workflow/search',
      before: getSearchBefore,
      after: getSearchAfterDelegate,
    },
    {
      method: 'post',
      url: '/api/workflow',
      before: postWorkflowBefore,
    },
    {
      method: 'get',
      url: '/api/workflow/:workflowId',
      after: getExecutionStatusAfter,
    },
    {
      method: 'delete',
      url: '/api/workflow/:workflowId/remove',
      before: removeWorkflowBefore,
    },
  ];
};

export default registration;
