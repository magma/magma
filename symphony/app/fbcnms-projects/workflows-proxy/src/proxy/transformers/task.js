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
import request from 'request';
import {withInfixSeparator} from '../utils.js';
import type {BeforeFun, TransformerRegistrationFun} from '../../types';

const logger = logging.getLogger(module);
let proxyTarget;

export const getLogBefore: BeforeFun = (
  tenantId,
  groups,
  req,
  res,
  proxyCallback,
) => {
  const url = proxyTarget + '/api/tasks/' + req.params.taskId;
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
      const task = JSON.parse(body);
      // make sure name starts with prefix
      const tenantWithInfixSeparator = withInfixSeparator(tenantId);
      if (task.workflowType?.indexOf(tenantWithInfixSeparator) == 0) {
        proxyCallback();
      } else {
        logger.error(
          `Error trying to get task of different tenant: ${tenantId},`,
          {task},
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

const registration: TransformerRegistrationFun = function(ctx) {
  proxyTarget = ctx.proxyTarget;
  return [
    {
      method: 'get',
      url: '/api/tasks/:taskId/log',
      before: getLogBefore,
    },
  ];
};

export default registration;
