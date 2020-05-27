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
  addTenantIdPrefix,
  anythingTo,
  assertAllowedSystemTask,
  createProxyOptionsBuffer,
  withInfixSeparator,
} from '../utils.js';

import type {
  AfterFun,
  BeforeFun,
  Task,
  TransformerRegistrationFun,
} from '../../types';

const logger = logging.getLogger(module);

// Gets all task definition
/*
curl  -H "x-auth-organization: fb-test" "localhost/proxy/api/metadata/taskdefs"
*/
const getAllTaskdefsAfter: AfterFun = (tenantId, groups, req, respObj) => {
  const tasks = anythingTo<Array<Task>>(respObj);
  // iterate over taskdefs, keep only those belonging to tenantId or global
  // remove tenantId prefix, keep GLOBAL_
  const tenantWithInfixSeparator = withInfixSeparator(tenantId);
  for (let idx = tasks.length - 1; idx >= 0; idx--) {
    const taskdef = tasks[idx];
    if (taskdef.name.indexOf(tenantWithInfixSeparator) == 0) {
      taskdef.name = taskdef.name.substr(tenantWithInfixSeparator.length);
    } else {
      // remove element
      tasks.splice(idx, 1);
    }
  }
};

// Used in POST and PUT
function sanitizeTaskdefBefore(tenantId: string, taskdef: Task): void {
  // only whitelisted system tasks are allowed
  assertAllowedSystemTask(taskdef);
  // prepend tenantId
  addTenantIdPrefix(tenantId, taskdef);
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
const postTaskdefsBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  // iterate over taskdefs, prefix with tenantId
  const reqObj = req.body;
  if (reqObj != null && Array.isArray(reqObj)) {
    for (let idx = 0; idx < reqObj.length; idx++) {
      const taskdef = anythingTo<Task>(reqObj[idx]);
      sanitizeTaskdefBefore(tenantId, taskdef);
    }
    proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
  } else {
    logger.error('Expected req.body to be array in postTaskdefsBefore');
    throw 'Expected req.body to be array in postTaskdefsBefore';
  }
};

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
const putTaskdefBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  const reqObj = req.body;
  if (reqObj != null && typeof reqObj === 'object') {
    const taskdef = anythingTo<Task>(reqObj);
    sanitizeTaskdefBefore(tenantId, taskdef);
    proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
  } else {
    logger.error('Expected req.body to be object in putTaskdefBefore');
    throw 'Expected req.body to be object in putTaskdefBefore';
  }
};

/*
curl -H "x-auth-organization: fb-test" \
 "localhost/proxy/api/metadata/taskdefs/frinx"
*/
// Gets the task definition
const getTaskdefByNameBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  req.params.name = withInfixSeparator(tenantId) + req.params.name;
  // modify url
  req.url = '/api/metadata/taskdefs/' + req.params.name;
  proxyCallback();
};

const getTaskdefByNameAfter: AfterFun = (
  tenantId,
  groups,
  req,
  respObj,
  res,
) => {
  const task = anythingTo<Task>(respObj);
  const tenantWithInfixSeparator = withInfixSeparator(tenantId);
  // remove prefix
  if (task.name.indexOf(tenantWithInfixSeparator) == 0) {
    task.name = task.name.substr(tenantWithInfixSeparator.length);
  } else {
    logger.error(
      `Tenant Id prefix '${tenantId}' not found, taskdef name: '${task.name}'`,
    );
    res.status(400);
    res.send('Prefix not found');
  }
};

// TODO: can this be disabled?
// Remove a task definition
/*
curl -H "x-auth-organization: fb-test" \
 "localhost/api/metadata/taskdefs/bar" -X DELETE -v
*/
const deleteTaskdefByNameBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  req.params.name = withInfixSeparator(tenantId) + req.params.name;
  // modify url
  req.url = '/api/metadata/taskdefs/' + req.params.name;
  proxyCallback();
};

const registration: TransformerRegistrationFun = () => [
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

export default registration;
