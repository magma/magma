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
  anythingTo,
  assertAllowedSystemTask,
  createProxyOptionsBuffer,
  isDecisionTask,
  isForkTask,
  isSubworkflowTask,
  objectToValues,
  withInfixSeparator,
} from '../utils.js';

import type {
  AfterFun,
  BeforeFun,
  Task,
  TransformerRegistrationFun,
  Workflow,
} from '../../types';

const logger = logging.getLogger(module);

// Utility used in PUT, POST before methods to check that submitted workflow
// and its tasks
// do not contain any prefix. Prefix is added to workflowdef if input is valid.
export function sanitizeWorkflowdefBefore(
  tenantId: string,
  workflowdef: Workflow,
) {
  // only whitelisted system tasks are allowed
  for (const task of workflowdef.tasks) {
    assertAllowedSystemTask(task);
  }
  // add prefix to tasks
  for (const task of workflowdef.tasks) {
    sanitizeWorkflowTaskdefBefore(tenantId, task);
  }
  // add prefix to workflow
  addTenantIdPrefix(tenantId, workflowdef);
}

function sanitizeWorkflowTaskdefBefore(tenantId: string, task: Task) {
  addTenantIdPrefix(tenantId, task);

  // add prefix to SUB_WORKFLOW tasks' referenced workflows
  if (isSubworkflowTask(task)) {
    addTenantIdPrefix(tenantId, anythingTo<Task>(task.subWorkflowParam));
  }

  // process decision tasks recursively
  if (isDecisionTask(task)) {
    const defaultCaseTasks = task.defaultCase ? task.defaultCase : [];
    for (const nestedTask of defaultCaseTasks) {
      sanitizeWorkflowTaskdefBefore(tenantId, nestedTask);
    }

    const decisionCaseIdToTasks = task.decisionCases ? task.decisionCases : {};
    for (const nestedTasks of objectToValues<string, Array<Task>>(
      decisionCaseIdToTasks,
    )) {
      for (const nestedTask of nestedTasks) {
        sanitizeWorkflowTaskdefBefore(tenantId, anythingTo<Task>(nestedTask));
      }
    }
  }

  // process fork tasks recursively
  if (isForkTask(task)) {
    const forkTasks: Array<Task> = task.forkTasks ? task.forkTasks : [];
    for (const nestedTask of forkTasks) {
      // handle sub lists in forked tasks
      if (Array.isArray(nestedTask)) {
        for (const l2nestedTask of nestedTask) {
          sanitizeWorkflowTaskdefBefore(tenantId, l2nestedTask);
        }
      } else {
        sanitizeWorkflowTaskdefBefore(tenantId, nestedTask);
      }
    }
  }
}

// Utility used after getting single or all workflowdefs to remove prefix from
// workflowdef names, taskdef names.
// Return true iif sanitization succeeded, false iif this
// workflowdef is invalid
export function sanitizeWorkflowdefAfter(
  tenantId: string,
  workflowdef: Workflow,
): boolean {
  const tenantWithInfixSeparator = withInfixSeparator(tenantId);
  if (workflowdef.name.indexOf(tenantWithInfixSeparator) == 0) {
    // keep only workflows with correct taskdefs,
    // allowed are GLOBAL and those with tenantId prefix which will be removed

    for (const task of workflowdef.tasks) {
      if (!sanitizeWorkflowTaskdefAfter(tenantWithInfixSeparator, task)) {
        return false;
      }
    }
    // remove prefix
    workflowdef.name = workflowdef.name.substr(tenantWithInfixSeparator.length);
    return true;
  } else {
    return false;
  }
}

function sanitizeWorkflowTaskdefAfter(
  tenantWithInfixSeparator: string,
  task: Task,
): boolean {
  // TODO refactor the boolean return value
  if (isForkTask(task)) {
    const forkTasks: Array<Task> = task.forkTasks ? task.forkTasks : [];
    for (const nestedTask of forkTasks) {
      if (Array.isArray(nestedTask)) {
        for (const l2nestedTask of nestedTask) {
          if (
            !sanitizeWorkflowTaskdefAfter(
              tenantWithInfixSeparator,
              l2nestedTask,
            )
          ) {
            return false;
          }
        }
      } else {
        if (
          !sanitizeWorkflowTaskdefAfter(tenantWithInfixSeparator, nestedTask)
        ) {
          return false;
        }
      }
    }
  }

  if (isDecisionTask(task)) {
    const defaultCaseTasks = task.defaultCase ? task.defaultCase : [];
    for (const nestedTask of defaultCaseTasks) {
      if (!sanitizeWorkflowTaskdefAfter(tenantWithInfixSeparator, nestedTask)) {
        return false;
      }
    }

    const decisionCaseIdToTasks = task.decisionCases ? task.decisionCases : {};
    for (const nestedTasks of objectToValues<string, Array<Task>>(
      decisionCaseIdToTasks,
    )) {
      for (const nestedTask of nestedTasks) {
        if (
          !sanitizeWorkflowTaskdefAfter(tenantWithInfixSeparator, nestedTask)
        ) {
          return false;
        }
      }
    }
  }

  // remove prefix from SUB_WORKFLOW tasks' referenced workflows
  if (isSubworkflowTask(task)) {
    if (task.subWorkflowParam != null) {
      const namedObject: {name: string} = task.subWorkflowParam;
      if (namedObject.name.indexOf(tenantWithInfixSeparator) == 0) {
        namedObject.name = namedObject.name.substr(
          tenantWithInfixSeparator.length,
        );
      } else {
        return false;
      }
    }
  }

  if (task.name.indexOf(tenantWithInfixSeparator) == 0) {
    // remove prefix
    task.name = task.name.substr(tenantWithInfixSeparator.length);
  } else {
    return false;
  }

  return true;
}

// Retrieves all workflow definition along with blueprint
/*
curl -H "x-auth-organization: fb-test" "localhost/proxy/api/metadata/workflow"
*/
const getAllWorkflowsAfter: AfterFun = (tenantId, req, respObj) => {
  const workflows: Array<Workflow> = anythingTo<Array<Workflow>>(respObj);
  // iterate over workflows, keep only those belonging to tenantId
  for (
    let workflowIdx = workflows.length - 1;
    workflowIdx >= 0;
    workflowIdx--
  ) {
    const workflowdef = workflows[workflowIdx];
    const ok = sanitizeWorkflowdefAfter(tenantId, workflowdef);
    if (!ok) {
      logger.warn(
        `Removing workflow with invalid task or name: ${JSON.stringify(
          workflowdef,
        )}`,
      );
      // remove element
      workflows.splice(workflowIdx, 1);
    }
  }
};

// Removes workflow definition. It does not remove workflows associated
// with the definition.
// Version is passed as url parameter.
/*
curl -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/metadata/workflow/2/2" -X DELETE
*/
const deleteWorkflowBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const tenantWithInfixSeparator = withInfixSeparator(tenantId);
  // change URL: add prefix to name
  const name = tenantWithInfixSeparator + req.params.name;
  const newUrl = `/api/metadata/workflow/${name}/${req.params.version}`;
  logger.debug(`Transformed url from '${req.url}' to '${newUrl}'`);
  req.url = newUrl;
  proxyCallback();
};

// Retrieves workflow definition along with blueprint
// Version is passed as query parameter.
/*
curl -H "x-auth-organization: fb-test" \
  "localhost/proxy/api/metadata/workflow/fx3?version=1"
*/
const getWorkflowBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const tenantWithInfixSeparator = withInfixSeparator(tenantId);
  const name = tenantWithInfixSeparator + req.params.name;
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
};

const getWorkflowAfter: AfterFun = (tenantId, req, respObj) => {
  const workflow = anythingTo<Workflow>(respObj);
  const ok = sanitizeWorkflowdefAfter(tenantId, workflow);
  if (!ok) {
    logger.error(
      `Possible error in code: response contains invalid task or` +
        `workflowdef name, tenant Id: ${tenantId}`,
    );
    throw 'Possible error in code: response contains' +
      ' invalid task or workflowdef name'; // TODO create Exception class
  }
};

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
const putWorkflowBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const workflows: Array<Workflow> = anythingTo<Array<Workflow>>(req.body);
  for (const workflowdef of workflows) {
    sanitizeWorkflowdefBefore(tenantId, workflowdef);
  }
  logger.debug(`Transformed request to ${JSON.stringify(workflows)}`);
  proxyCallback({buffer: createProxyOptionsBuffer(workflows, req)});
};

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
const postWorkflowBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const workflow: Workflow = anythingTo<Workflow>(req.body);
  sanitizeWorkflowdefBefore(tenantId, workflow);
  logger.debug(`Transformed request to ${JSON.stringify(workflow)}`);
  proxyCallback({buffer: createProxyOptionsBuffer(workflow, req)});
};

const registration: TransformerRegistrationFun = () => [
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

export default registration;
