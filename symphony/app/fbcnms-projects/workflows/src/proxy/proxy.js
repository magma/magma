/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import Router from 'express';
import bodyParser from 'body-parser';
import httpProxy from 'http-proxy';
import logging from '@fbcnms/logging';
import transformerRegistry from './transformer-registry';
import {getTenantId} from './utils.js';

const logger = logging.getLogger(module);
const router = Router();
router.use(bodyParser.urlencoded({extended: false}));
router.use('/', bodyParser.json());

export default async function(proxyTarget) {
  const transformers = await transformerRegistry({proxyTarget});

  // Configure http-proxy
  const proxy = httpProxy.createProxyServer({
    target: proxyTarget,
  });

  for (const entry of transformers) {
    logger.info(`Routing url:${entry.url}, method:${entry.method}`);
    router[entry.method](entry.url, async (req, res) => {
      let tenantId;
      try {
        tenantId = getTenantId(req);
      } catch (err) {
        res.status(400);
        res.send('Cannot get tenantId:' + err);
        return;
      }
      // prepare 'after'
      const _write = res.write; // backup real write method
      // create wrapper that allows transforming output from target
      res.write = function(data) {
        if (entry.after) {
          // TODO: parse only if data is json
          const respObj = JSON.parse(data);
          entry.after(tenantId, req, respObj, res);
          data = JSON.stringify(respObj);
        }
        _write.call(res, data);
      };

      // start with 'before'
      logger.info(
        `REQ ${req.method} ${
          req.url
        } tenantId ${tenantId} body ${JSON.stringify(req.body)}`,
      );
      const proxyCallback = function(proxyOptions) {
        proxy.web(req, res, proxyOptions);
      };
      if (entry.before) {
        try {
          entry.before(tenantId, req, res, proxyCallback);
        } catch (err) {
          console.error('Got error in beforeFun', err);
          res.status(500);
          res.send('Cannot send request: ' + err);
          return;
        }
      } else {
        proxyCallback();
      }
    });
  }
  return router;
}
