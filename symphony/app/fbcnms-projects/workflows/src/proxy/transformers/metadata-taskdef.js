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
import {
  GLOBAL_PREFIX,
  INFIX_SEPARATOR,
  createProxyOptionsBuffer,
  isAllowedSystemTask,
  withUnderscore,
} from '../utils.js';

const logger = logging.getLogger(module);

// Gets all task definition
/*
curl  -H "x-auth-organization: fb-test" "localhost/proxy/api/metadata/taskdefs"
*/
function getAllTaskdefsAfter(tenantId, req, respObj) {
  // iterate over taskdefs, keep only those belonging to tenantId or global
  // remove tenantId prefix, keep GLOBAL_
  const tenantWithUnderscore = withUnderscore(tenantId);
  const globalWithUnderscore = withUnderscore(GLOBAL_PREFIX);
  for (let idx = respObj.length - 1; idx >= 0; idx--) {
    const taskdef = respObj[idx];
    if (taskdef.name.indexOf(tenantWithUnderscore) == 0) {
      taskdef.name = taskdef.name.substr(tenantWithUnderscore.length);
    } else if (taskdef.name.indexOf(globalWithUnderscore) == 0) {
      // noop
    } else {
      // remove element
      respObj.splice(idx, 1);
    }
  }
}

// Used in POST and PUT
function sanitizeTaskdefBefore(tenantId, taskdef) {
  const tenantWithUnderscore = withUnderscore(tenantId);
  if (taskdef.name.indexOf(INFIX_SEPARATOR) > -1) {
    logger.error(
      `Name of taskdef must not contain '${INFIX_SEPARATOR}': '${taskdef.name}'`,
    );
    throw 'Name must not contain underscore'; // TODO create Exception class
  }
  // only whitelisted system tasks are allowed
  if (!isAllowedSystemTask(taskdef)) {
    logger.error(
      `Task type is not allowed: '${tenantId}'` +
        ` in '${JSON.stringify(taskdef)}'`,
    );
    // TODO create Exception class
    throw 'Task type is not allowed';
  }
  // prepend tenantId
  taskdef.name = tenantWithUnderscore + taskdef.name;
}
// Create new task definition(s)
// Underscore in name is not allowed.
/*
curl -X POST -H "x-auth-organization: fb-test"  \
 "localhost/proxy/api/metadata/taskdefs" \
 -H 'Content-Type: application/json' -d '
[
    {
      "name": "bar",
      "retryCount": 3,
      "retryLogic": "FIXED",
      "retryDelaySeconds": 10,
      "timeoutSeconds": 300,
      "timeoutPolicy": "TIME_OUT_WF",
      "responseTimeoutSeconds": 180,
      "ownerEmail": "foo@bar.baz"
    }
]
'
*/
// TODO: should this be disabled?
function postTaskdefsBefore(tenantId, req, res, proxyCallback) {
  // iterate over taskdefs, prefix with tenantId
  const reqObj = req.body;
  for (let idx = 0; idx < reqObj.length; idx++) {
    const taskdef = reqObj[idx];
    sanitizeTaskdefBefore(tenantId, taskdef);
  }
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
}

// Update an existing task
// Underscore in name is not allowed.
/*
curl -X PUT -H "x-auth-organization: fb-test" \
 "localhost/proxy/api/metadata/taskdefs" \
 -H 'Content-Type: application/json' -d '
    {
      "name": "frinx",
      "retryCount": 3,
      "retryLogic": "FIXED",
      "retryDelaySeconds": 10,
      "timeoutSeconds": 400,
      "timeoutPolicy": "TIME_OUT_WF",
      "responseTimeoutSeconds": 180,
      "ownerEmail": "foo@bar.baz"
    }
'
*/
// TODO: should this be disabled?
function putTaskdefBefore(tenantId, req, res, proxyCallback) {
  const reqObj = req.body;
  const taskdef = reqObj;
  sanitizeTaskdefBefore(tenantId, taskdef);
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
}

/*
curl -H "x-auth-organization: fb-test" \
 "localhost/proxy/api/metadata/taskdefs/frinx"
*/
// Gets the task definition
function getTaskdefByNameBefore(tenantId, req, res, proxyCallback) {
  req.params.name = withUnderscore(tenantId) + req.params.name;
  // modify url
  req.url = '/api/metadata/taskdefs/' + req.params.name;
  proxyCallback();
}

function getTaskdefByNameAfter(tenantId, req, respObj, res) {
  if (res.status == 200) {
    const tenantWithUnderscore = withUnderscore(tenantId);
    // remove prefix
    if (respObj.name && respObj.name.indexOf(tenantWithUnderscore) == 0) {
      respObj.name = respObj.name.substr(tenantWithUnderscore.length);
    } else {
      logger.error(
        `Tenant Id prefix '${tenantId}' not found, taskdef name: '${respObj.name}'`,
      );
      res.status(400);
      res.send('Prefix not found'); // TODO: this exits the process
    }
  }
}

// TODO: can this be disabled?
// Remove a task definition
/*
curl -H "x-auth-organization: fb-test" \
 "localhost/api/metadata/taskdefs/bar" -X DELETE -v
*/
function deleteTaskdefByNameBefore(tenantId, req, res, proxyCallback) {
  req.params.name = withUnderscore(tenantId) + req.params.name;
  // modify url
  req.url = '/api/metadata/taskdefs/' + req.params.name;
  proxyCallback();
}

export default function() {
  return [
    {
      method: 'get',
      url: '/api/metadata/taskdefs',
      after: getAllTaskdefsAfter,
    },
    {
      method: 'post',
      url: '/api/metadata/taskdefs',
      before: postTaskdefsBefore,
    },
    {
      method: 'put',
      url: '/api/metadata/taskdefs',
      before: putTaskdefBefore,
    },
    {
      method: 'get',
      url: '/api/metadata/taskdefs/:name',
      before: getTaskdefByNameBefore,
      after: getTaskdefByNameAfter,
    },
    {
      method: 'delete',
      url: '/api/metadata/taskdefs/:name',
      before: deleteTaskdefByNameBefore,
    },
  ];
}
