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
import {
  addTenantIdPrefix,
  anythingTo,
  createProxyOptionsBuffer,
  findValuesByJsonPath,
  removeTenantPrefix,
  withInfixSeparator,
} from '../utils.js';

import type {
  AfterFun,
  BeforeFun,
  ScheduleRequest,
  TransformerRegistrationFun,
} from '../../types';

const logger = logging.getLogger(module);
const versionSeparator = ':';
let schellarTarget: string;

// Used in POST and PUT
function sanitizeScheduleBefore(
  tenantId: string,
  schedule: ScheduleRequest,
): void {
  const expectedName =
    schedule.workflowName + versionSeparator + schedule.workflowVersion;
  if (schedule.name !== expectedName) {
    logger.error('Unexpected schedule name', {schedule, expectedName});
    throw 'Unexpected schedule name';
  }
  if (expectedName.indexOf('/') > -1) {
    // guard against https://github.com/flaviostutz/schellar/issues/4
    logger.error('Cannot create schedule with name containing "/" character');
    throw 'Cannot create schedule with name containing "/" character';
  }

  // add tenantId to name
  addTenantIdPrefix(tenantId, schedule);
  // add tenantId to workflowName
  schedule.workflowName = withInfixSeparator(tenantId) + schedule.workflowName;
}

/*
curl http://localhost/proxy/schedule \
  -H "X-Auth-Organization: fb-test" \
  -H 'Content-Type: application/json'
*/
const getAllBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  req.url = '/schedule';
  proxyCallback({target: schellarTarget});
};

function removeWrongPrefixesFromArray(
  tenantId,
  array: Array<mixed>,
  jsonPath: string,
) {
  const tenantWithInfixSeparator = withInfixSeparator(tenantId);
  for (let idx = array.length - 1; idx >= 0; idx--) {
    const item = array[idx];
    const leafValues = findValuesByJsonPath(item, jsonPath, 'value');
    let prefixFound = true;
    for (const leaf of leafValues) {
      if (!leaf || leaf.indexOf(tenantWithInfixSeparator) != 0) {
        prefixFound = false;
        break;
      }
    }

    if (!prefixFound) {
      array.splice(idx, 1);
    }
  }
}

const getAllAfter: AfterFun = (tenantId, req, respObj) => {
  if (respObj != null && Array.isArray(respObj)) {
    removeWrongPrefixesFromArray(tenantId, anythingTo(respObj), '$.name');
    removeTenantPrefix(tenantId, respObj, '$[*].workflowName', false);
    removeTenantPrefix(tenantId, respObj, '$[*].name', false);
  } else {
    logger.error('Unexpected response', {respObj});
    throw 'Unexpected response';
  }
};

/*
curl http://localhost/proxy/schedule/workflow1 \
  -H "X-Auth-Organization: fb-test" \
  -H 'Content-Type: application/json'
*/
const getBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const reqName = req.params.name;
  req.url = '/schedule/' + withInfixSeparator(tenantId) + reqName;
  proxyCallback({target: schellarTarget});
};
const getAfter: AfterFun = (tenantId, req, respObj) => {
  removeTenantPrefix(tenantId, respObj, '$.workflowName', false);
  removeTenantPrefix(tenantId, respObj, '$.name', false);
};

/*
curl -X POST http://localhost/proxy/schedule \
  -H "X-Auth-Organization: fb-test" \
  -H 'Content-Type: application/json' \
  -d '
  {
  "name": "workflow1:1",
  "enabled": true,
  "parallelRuns": false,
  "workflowName": "workflow1",
  "workflowVersion": "1",
  "cronString": "0 * * ? * *",
  "workflowContext": {
    "param1": "value1",
    "param2": "value2"
  },
  "fromDate": "2019-01-01T15:04:05Z",
  "toDate": "2029-07-01T15:04:05Z"
  }
'
*/
const postBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  req.url = '/schedule';
  const schedule = anythingTo<ScheduleRequest>(req.body);
  sanitizeScheduleBefore(tenantId, schedule);
  const buffer = createProxyOptionsBuffer(schedule, req);
  proxyCallback({target: schellarTarget, buffer});
};

/*
curl -X PUT http://localhost/proxy/schedule/workflow1 \
  -H "X-Auth-Organization: fb-test" \
  -H 'Content-Type: application/json' \
  -d '
  {
  "name": "workflow1:1",
  "enabled": true,
  "parallelRuns": false,
  "workflowName": "workflow1",
  "workflowVersion": "1",
  "cronString": "0 * * ? * *",
  "workflowContext": {
    "param1": "value1",
    "param2": "value2"
  },
  "fromDate": "2019-01-01T15:04:05Z",
  "toDate": "2029-07-01T15:04:05Z"
  }
'
*/
// Renaming is not supported by proxy - url name must be equal to workflowName
const putBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const schedule = anythingTo<ScheduleRequest>(req.body);
  let reqName = req.params.name;
  if (reqName !== schedule.name) {
    logger.error('Schedule name must be equal to name supplied in url', {
      schedule,
      reqName,
    });
    // TODO create Exception class
    throw 'Schedule name must be equal to name supplied in url';
  }
  sanitizeScheduleBefore(tenantId, schedule);
  reqName = schedule.name;
  req.url = '/schedule/' + reqName;
  const buffer = createProxyOptionsBuffer(schedule, req);
  proxyCallback({target: schellarTarget, buffer});
};

/*
curl -X DELETE \
  -H "X-Auth-Organization: fb-test" \
  -H 'Content-Type: application/json' \
  http://localhost/proxy/schedule/workflow1
*/
const deleteBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const reqName = req.params.name;
  req.url = '/schedule/' + withInfixSeparator(tenantId) + reqName;
  proxyCallback({target: schellarTarget});
};

const registration: TransformerRegistrationFun = function(ctx) {
  schellarTarget = ctx.schellarTarget;
  return [
    {
      method: 'get',
      url: '/schedule/?',
      before: getAllBefore,
      after: getAllAfter,
    },
    {
      method: 'get',
      url: '/schedule/:name',
      before: getBefore,
      after: getAfter,
    },
    {
      method: 'post',
      url: '/schedule/?',
      before: postBefore,
    },
    {
      method: 'put',
      url: '/schedule/:name',
      before: putBefore,
    },
    {
      method: 'delete',
      url: '/schedule/:name',
      before: deleteBefore,
    },
  ];
};

export default registration;
