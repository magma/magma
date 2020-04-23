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
// Currently just filters result without passing prefix to conductor.
// TODO: implement querying by prefix in conductor

import qs from 'qs';
import {
  addTenantIdPrefix,
  assertAllowedSystemTask,
  createProxyOptionsBuffer,
  withUnderscore,
} from '../utils.js';

const logger = logging.getLogger(module);

// Utility used in PUT, POST before methods to check that submitted workflow
// and its tasks
// do not contain any prefix. Prefix is added to workflowdef if input is valid.
function sanitizeWorkflowdefBefore(tenantId, workflowdef) {
  // only whitelisted system tasks are allowed
  for (const task of workflowdef.tasks) {
    assertAllowedSystemTask(task);
  }
  // add prefix to tasks
  for (const task of workflowdef.tasks) {
    addTenantIdPrefix(tenantId, task);
  }
  // add prefix to workflow
  addTenantIdPrefix(tenantId, workflowdef);
}

// Utility used after getting single or all workflowdefs to remove prefix from
// workflowdef names, taskdef names.
// Return true iif sanitization succeeded, false iif this
// workflowdef is invalid
function sanitizeWorkflowdefAfter(tenantId, workflowdef) {
  const tenantWithUnderscore = withUnderscore(tenantId);
  if (workflowdef.name.indexOf(tenantWithUnderscore) == 0) {
    // keep only workflows with correct taskdefs,
    // allowed are GLOBAL and those with tenantId prefix which will be removed
    for (const task of workflowdef.tasks) {
      if (task.name.indexOf(tenantWithUnderscore) == 0) {
        // remove prefix
        task.name = task.name.substr(tenantWithUnderscore.length);
      } else {
        return false;
      }
    }
    // remove prefix
    workflowdef.name = workflowdef.name.substr(tenantWithUnderscore.length);
    return true;
  } else {
    return false;
  }
}

// Retrieves all workflow definition along with blueprint
/*
curl -H "x-auth-organization: fb-test" "localhost/proxy/api/metadata/workflow"
*/
function getAllWorkflowsAfter(tenantId, req, respObj) {
  // iterate over workflows, keep only those belonging to tenantId
  for (let workflowIdx = respObj.length - 1; workflowIdx >= 0; workflowIdx--) {
    const workflowdef = respObj[workflowIdx];
    const ok = sanitizeWorkflowdefAfter(tenantId, workflowdef, respObj);
    if (!ok) {
      logger.warn(
        `Removing workflow with invalid task or name: ${JSON.stringify(
          workflowdef,
        )}`,
      );
      // remove element
      respObj.splice(workflowIdx, 1);
    }
  }
}

// Removes workflow definition. It does not remove workflows associated
// with the definition.
// Version is passed as url parameter.
/*
curl -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/metadata/workflow/2/2" -X DELETE
*/
function deleteWorkflowBefore(tenantId, req, res, proxyCallback) {
  const tenantWithUnderscore = withUnderscore(tenantId);
  // change URL: add prefix to name
  const name = tenantWithUnderscore + req.params.name;
  const newUrl = `/api/metadata/workflow/${name}/${req.params.version}`;
  logger.debug(`Transformed url from '${req.url}' to '${newUrl}'`);
  req.url = newUrl;
  proxyCallback();
}

// Retrieves workflow definition along with blueprint
// Version is passed as query parameter.
/*
curl -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/metadata/workflow/fx3?version=1"
*/
function getWorkflowBefore(tenantId, req, res, proxyCallback) {
  const tenantWithUnderscore = withUnderscore(tenantId);
  const name = tenantWithUnderscore + req.params.name;
  let newUrl = `/api/metadata/workflow/${name}`;
  const originalQueryString = req._parsedUrl.query;
  const parsedQuery = qs.parse(originalQueryString);
  const version = parsedQuery['version'];
  if (version) {
    newUrl += '?version=' + version;
  }
  logger.debug(`Transformed url from '${req.url}' to '${newUrl}'`);
  req.url = newUrl;
  proxyCallback();
}

function getWorkflowAfter(tenantId, req, respObj) {
  const ok = sanitizeWorkflowdefAfter(tenantId, respObj);
  if (!ok) {
    logger.error(
      `Possible error in code: response contains invalid task or` +
        `workflowdef name, tenant Id: ${tenantId}`,
    );
    throw 'Possible error in code: response contains' +
      ' invalid task or workflowdef name'; // TODO create Exception class
  }
}

// Create or update workflow definition
// Underscore in name is not allowed.
/*
curl -X PUT -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/metadata/workflow" \
  -H 'Content-Type: application/json' -d '
[
    {
    "name": "fx3",
    "description": "foo1",
    "ownerEmail": "foo@bar.baz",
    "version": 1,
    "schemaVersion": 2,
    "tasks": [
        {
        "name": "bar",
        "taskReferenceName": "barref",
        "type": "SIMPLE",
        "inputParameters": {}
        }
    ]
    }
]'


curl -X PUT -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/metadata/workflow" \
  -H 'Content-Type: application/json' -d '
[
    {
    "name": "fx3",
    "description": "foo1",
    "ownerEmail": "foo@bar.baz",
    "version": 1,
    "schemaVersion": 2,
    "tasks": [
        {
        "name": "bar",
        "taskReferenceName": "barref",
        "type": "SIMPLE",
        "inputParameters": {}
        },
        {
        "name": "GLOBAL_GLOBAL1",
        "taskReferenceName": "globref",
        "type": "SIMPLE",
        "inputParameters": {}
        }
    ]
    }
]'
*/
function putWorkflowBefore(tenantId, req, res, proxyCallback) {
  const reqObj = req.body;
  for (const workflowdef of reqObj) {
    sanitizeWorkflowdefBefore(tenantId, workflowdef);
  }
  logger.debug(`Transformed request to ${JSON.stringify(reqObj)}`);
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
}

// Create a new workflow definition
// Underscore in name is not allowed.
/*
curl -X POST -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/metadata/workflow" \
  -H 'Content-Type: application/json' -d '

    {
    "name": "fx3",
    "description": "foo1",
    "ownerEmail": "foo@bar.baz",
    "version": 1,
    "schemaVersion": 2,
    "tasks": [
        {
        "name": "bar",
        "taskReferenceName": "barref",
        "type": "SIMPLE",
        "inputParameters": {}
        },
        {
        "name": "GLOBAL_GLOBAL1",
        "taskReferenceName": "globref",
        "type": "SIMPLE",
        "inputParameters": {}
        }
    ]
    }
'
*/
function postWorkflowBefore(tenantId, req, res, proxyCallback) {
  const reqObj = req.body;
  sanitizeWorkflowdefBefore(tenantId, reqObj);
  logger.debug(`Transformed request to ${JSON.stringify(reqObj)}`);
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
}

export default function() {
  return [
    {
      method: 'get',
      url: '/api/metadata/workflow',
      after: getAllWorkflowsAfter,
    },
    {
      method: 'delete',
      url: '/api/metadata/workflow/:name/:version',
      before: deleteWorkflowBefore,
    },
    {
      method: 'get',
      url: '/api/metadata/workflow/:name',
      before: getWorkflowBefore,
      after: getWorkflowAfter,
    },
    {
      method: 'put',
      url: '/api/metadata/workflow',
      before: putWorkflowBefore,
    },
    {
      method: 'post',
      url: '/api/metadata/workflow',
      before: postWorkflowBefore,
    },
  ];
}
