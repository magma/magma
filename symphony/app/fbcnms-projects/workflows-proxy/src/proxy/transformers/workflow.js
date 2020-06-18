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

import qs from 'qs';
import request from 'request';
import {
  addTenantIdPrefix,
  anythingTo,
  createProxyOptionsBuffer,
  removeTenantPrefix,
  removeTenantPrefixes,
  withInfixSeparator,
} from '../utils.js';

import type {
  AfterFun,
  BeforeFun,
  StartWorkflowRequest,
  TransformerRegistrationFun,
} from '../../types';

const logger = logging.getLogger(module);
// Search for workflows based on payload and other parameters
/*
 curl -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/workflow/search?query=status+IN+(FAILED)"
*/
export const getSearchBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  // prefix query with workflowType STARTS_WITH tenantId_
  const originalQueryString = req._parsedUrl.query;
  const limitToTenant = `workflowType STARTS_WITH \'${tenantId}_'`;
  const newQueryString = updateQuery(originalQueryString, limitToTenant);
  req.url = req._parsedUrl.pathname + '?' + newQueryString;
  proxyCallback();
};

export const getSearchAfter: AfterFun = (tenantId, groups, req, respObj) => {
  removeTenantPrefix(tenantId, respObj, 'results[*].workflowType', false);
};

export function updateQuery(
  originalQueryString: string,
  queryExpanded: string,
): string {
  const parsedQuery = qs.parse(originalQueryString);
  let q = parsedQuery['query'];
  if (q) {
    // TODO: validate conductor query to prevent security issues
    q = `(${q} AND (${queryExpanded}))`;
  } else {
    q = `(${queryExpanded})`;
  }
  parsedQuery['query'] = q;
  const newQueryString = qs.stringify(parsedQuery);
  logger.debug(
    `Transformed query string from ` +
      `'${originalQueryString}' to '${newQueryString}`,
  );
  return newQueryString;
}

// Start a new workflow with StartWorkflowRequest, which allows task to be
// executed in a domain
/*
curl -X POST -H "x-auth-organization: fb-test" -H \
"Content-Type: application/json" "localhost/proxy/api/workflow" -d '
{
  "name": "fx3",
  "version": 1,
  "correlatonId": "corr1",
  "ownerApp": "my_owner_app",
  "input": {
  }
}
'
*/
export const postWorkflowBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  // name must start with prefix
  const tenantWithInfixSeparator = withInfixSeparator(tenantId);
  const reqObj = anythingTo<StartWorkflowRequest>(req.body);

  // workflowDef section is not allowed (no dynamic workflows)
  if (reqObj.workflowDef != null) {
    logger.error(
      `Section workflowDef is not allowed ${JSON.stringify(reqObj)}`,
    );
    throw 'Section workflowDef is not allowed';
  }
  // taskToDomain section is not allowed
  if (reqObj.taskToDomain) {
    logger.error(
      `Section taskToDomain is not allowed ${JSON.stringify(reqObj)}`,
    );
    throw 'Section taskToDomain is not allowed';
  }

  // add prefix
  addTenantIdPrefix(tenantId, reqObj);
  // add taskToDomain
  reqObj.taskToDomain = {};
  //TODO: is this OK?
  reqObj.taskToDomain[tenantWithInfixSeparator + '*'] = tenantId;
  logger.debug(`Transformed request to ${JSON.stringify(reqObj)}`);
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
};

// Gets the workflow by workflow id
/*
curl  -H "x-auth-organization: fb-test" \
    "localhost/proxy/api/workflow/c0a438d4-25b7-4c12-8a29-3473d98b1ad7"
*/
export const getExecutionStatusAfter: AfterFun = (
  tenantId,
  groups,
  req,
  respObj,
) => {
  const jsonPathToAllowGlobal = {
    workflowName: false,
    workflowType: false,
    'tasks[*].taskDefName': true,
    'tasks[*].workflowTask.name': true,
    'tasks[*].workflowTask.taskDefinition.name': true,
    'tasks[*].workflowType': false,
    'tasks[*].inputData.subWorkflowName': false,
    'tasks[*].workflowType': false,
    'tasks[*].outputData.workflowType': false,
    'tasks[*].workflowTask.subWorkflowParam.name': false,
    'output.workflowType': false,
    'workflowDefinition.name': false,
    'workflowDefinition.tasks[*].name': true,
    'workflowDefinition.tasks[*].taskDefinition.name': true,
    'workflowDefinition.tasks[*].subWorkflowParam.name': false,
  };
  removeTenantPrefixes(tenantId, respObj, jsonPathToAllowGlobal);
};

// Removes the workflow from the system
/*
curl  -H "x-auth-organization: fb-test" \
    "localhost/proxy/api/workflow/2dbb6e3e-c45d-464b-a9c9-2bbb16b7ca71/remove" \
    -X DELETE
*/
export const removeWorkflowBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  const url = proxyTarget + '/api/workflow/' + req.params.workflowId;
  // first make a HTTP request to validate that this workflow belongs to tenant
  const requestOptions = {
    url,
    method: 'GET',
    headers: {
      'Content-Type': 'application/javascript',
    },
  };
  logger.debug(`Requesting ${JSON.stringify(requestOptions)}`);
  request(requestOptions, function(error, response, body) {
    logger.debug(`Got status code: ${response.statusCode}, body: '${body}'`);
    if (response.statusCode == 200) {
      const workflow = JSON.parse(body);
      // make sure name starts with prefix
      const tenantWithInfixSeparator = withInfixSeparator(tenantId);
      if (workflow.workflowName.indexOf(tenantWithInfixSeparator) == 0) {
        proxyCallback();
      } else {
        logger.error(
          `Error trying to delete workflow of different tenant: ${tenantId},` +
            ` workflow: ${JSON.stringify(workflow)}`,
        );
        res.status(401);
        res.send('Unauthorized');
      }
    } else {
      res.status(response.statusCode);
      res.send(body);
    }
  });
};

let proxyTarget: string;

const registration: TransformerRegistrationFun = function(ctx) {
  proxyTarget = ctx.proxyTarget;
  return [
    {
      method: 'get',
      url: '/api/workflow/search',
      before: getSearchBefore,
      after: getSearchAfter,
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
