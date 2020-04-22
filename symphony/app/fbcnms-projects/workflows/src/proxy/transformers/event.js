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
import {createProxyOptionsBuffer} from '../utils.js';

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
function postEventBefore(tenantId, req, res, proxyCallback) {
  const reqObj = req.body;
  logger.debug(`Transforming '${JSON.stringify(reqObj)}'`);
  // TODO: prefix name
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
}

function putEventBefore(tenantId, req, res, proxyCallback) {
  const reqObj = req.body;
  logger.debug(`Transforming '${JSON.stringify(reqObj)}'`);
  // TODO: prefix name
  proxyCallback({buffer: createProxyOptionsBuffer(reqObj, req)});
}

export default function() {
  return [
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
}
