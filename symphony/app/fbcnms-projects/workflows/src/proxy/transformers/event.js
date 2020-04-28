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
  assertValueIsWithoutInfixSeparator,
  createProxyOptionsBuffer,
  withInfixSeparator,
} from '../utils.js';

import type {BeforeFun, Event, TransformerRegistrationFun} from '../../types';

const logger = logging.getLogger(module);

/*
curl -H "x-auth-organization: fb-test" \
 "localhost/proxy/api/event" -X PUT -H 'Content-Type: application/json' -d '
{
    "actions": [
        {
            "action": "complete_task",
            "complete_task": {
                "output": {},
                "taskRefName": "${targetTaskRefName}",
                "workflowId": "${targetWorkflowId}"
            }
        }
    ],
    "active": true,
    "event": "conductor:event:eventTaskRefZUEX",
    "name": "event_eventTaskRefZUEX"
}
' -v
*/

function sanitizeEvent(tenantId: string, event: Event) {
  // prefix event name
  addTenantIdPrefix(tenantId, event);
  // 'event' attribute uses following format:
  // conductor:WORKFLOW_NAME:TASK_REFERENCE
  // workflow name must be prefixed.
  const split = event.event.split(':');
  if (split.length == 3 && split[0] === 'conductor') {
    let workflowName = split[1];
    assertValueIsWithoutInfixSeparator(workflowName);
    workflowName = withInfixSeparator(tenantId) + workflowName;
    event.event = split[0] + ':' + workflowName + ':' + split[2];
  } else {
    logger.error(
      `Tenant ${tenantId} sent invalid event ` + `${JSON.stringify(event)}`,
    );
  }
}

const postEventBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const reqObj = req.body;
  logger.debug(`Transforming '${JSON.stringify(reqObj)}'`);
  sanitizeEvent(tenantId, anythingTo<Event>(reqObj));
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
};

const putEventBefore: BeforeFun = (tenantId, req, res, proxyCallback) => {
  const reqObj = req.body;
  logger.debug(`Transforming '${JSON.stringify(reqObj)}'`);
  sanitizeEvent(tenantId, anythingTo<Event>(reqObj));
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
};

const registration: TransformerRegistrationFun = () => [
  {
    method: 'post',
    url: '/api/event',
    before: postEventBefore,
  },
  {
    method: 'put',
    url: '/api/event',
    before: putEventBefore,
  },
];

export default registration;
