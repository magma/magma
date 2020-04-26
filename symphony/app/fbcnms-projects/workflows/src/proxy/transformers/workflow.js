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
  createProxyOptionsBuffer,
  removeTenantPrefix,
  removeTenantPrefixes,
  withInfixSeparator,
} from '../utils.js';

const logger = logging.getLogger(module);
// Search for workflows based on payload and other parameters
/*
 curl -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/workflow/search?query=status+IN+(FAILED)"
*/
function getSearchBefore(tenantId, req, res, proxyCallback) {
  // prefix query with workflowType STARTS_WITH tenantId_
  const originalQueryString = req._parsedUrl.query;
  const parsedQuery = qs.parse(originalQueryString);
  const limitToTenant = `workflowType STARTS_WITH \'${tenantId}_'`;
  let q = parsedQuery['query'];
  if (q) {
    // TODO: validate conductor query to prevent security issues
    q = limitToTenant + ' AND (' + q + ')';
  } else {
    q = limitToTenant;
  }
  parsedQuery['query'] = q;
  const newQueryString = qs.stringify(parsedQuery);
  logger.debug(
    `Transformed query string from ` +
      `'${originalQueryString}' to '${newQueryString}`,
  );
  req.url = req._parsedUrl.pathname + '?' + newQueryString;
  proxyCallback();
}

function getSearchAfter(tenantId, req, respObj) {
  removeTenantPrefix(tenantId, respObj, 'results[*].workflowType', false);
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
function postWorkflowBefore(tenantId, req, res, proxyCallback) {
  // name must start with prefix
  const tenantWithInfixSeparator = withInfixSeparator(tenantId);
  const reqObj = req.body;

  // workflowDef section is not allowed (no dynamic workflows)
  if (reqObj.workflowDef) {
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
}

// Gets the workflow by workflow id
/*
curl  -H "x-auth-organization: fb-test" \
    "localhost/proxy/api/workflow/c0a438d4-25b7-4c12-8a29-3473d98b1ad7"
*/
function getExecutionStatusAfter(tenantId, req, respObj) {
  const jsonPathToAllowGlobal = {
    workflowName: false,
    workflowType: false,
    'tasks[*].taskDefName': true,
    'tasks[*].workflowTask.name': true,
    'tasks[*].workflowTask.taskDefinition.name': true,
    'tasks[*].workflowType': false,
    'workflowDefinition.name': false,
    'workflowDefinition.tasks[*].name': true,
    'workflowDefinition.tasks[*].taskDefinition.name': true,
  };
  removeTenantPrefixes(tenantId, respObj, jsonPathToAllowGlobal);
}

// Removes the workflow from the system
/*
curl  -H "x-auth-organization: fb-test" \
    "localhost/proxy/api/workflow/2dbb6e3e-c45d-464b-a9c9-2bbb16b7ca71/remove" \
    -X DELETE
*/
function removeWorkflowBefore(tenantId, req, res, proxyCallback) {
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
}

let proxyTarget;

export default function(ctx) {
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
}
