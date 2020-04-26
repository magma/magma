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
    // TODO set timeouts
  });

  for (const entry of transformers) {
    logger.info(`Routing url:${entry.url}, method:${entry.method}`);
    router[entry.method](entry.url, async (req, res, next) => {
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
      // FIXME: this is a hack to be able to modify response.
      // It only works if response is not empty.
      res.write = function(data) {
        if (res.statusCode >= 200 && res.statusCode < 300 && entry.after) {
          let respObj = null;
          try {
            // conductor does not always send correct Content-Type, e.g. on 404
            respObj = JSON.parse(data);
          } catch (e) {
            logger.warn('Response is not JSON');
          }
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
        proxy.web(req, res, proxyOptions, function(e) {
          logger.error('Inline error handler', e);
          next(e);
        });
      };
      if (entry.before) {
        try {
          entry.before(tenantId, req, res, proxyCallback);
        } catch (err) {
          logger.error('Got error in beforeFun', err);
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
