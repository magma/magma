/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import logging from '@fbcnms/logging';
import {checkWasmer} from './wasmer.js';
import {executePython, pythonHealthCheck} from './python.js';
import {executeQuickJs, quickJsHealthCheck} from './quickjs.js';
const logger = logging.getLogger(module);

const ConductorClient = require('conductor-client').default;
// properties
const conductorApiUrl =
  process.env.CONDUCTOR_API_URL || 'http://conductor-server:8080/api';
const maxRunner = process.env.MAX_RUNNER || 1;

const conductorClient = new ConductorClient({
  baseURL: conductorApiUrl,
});

function registerWasmWorker(workerSuffix, callback) {
  conductorClient.registerWatcher(
    'GLOBAL___' + workerSuffix,
    callback,
    {pollingIntervals: 1000, autoAck: true, maxRunner: maxRunner},
    true,
  );
}

async function createTaskResult(
  outputIsJson,
  outputData,
  stderr,
  updaterFun,
  reasonForIncompletion,
) {
  const logs = stderr.split('\n').filter(String);
  if (outputIsJson) {
    // convert back to object
    try {
      outputData.result = JSON.parse(outputData.result);
    } catch (e) {
      logs.push('Cannot convert stdout to json');
    }
  }
  logger.info('createTaskResult updating task', {outputData, updaterFun});
  await updaterFun({
    outputData,
    logs,
    reasonForIncompletion,
  });
}

async function checkAndRegister(wasmSuffix, healthCheckFn, executeFn) {
  if (!(await healthCheckFn())) {
    logger.warn(wasmSuffix + ' healthcheck failed');
  }
  registerWasmWorker(wasmSuffix, async (data, updater) => {
    logger.info(wasmSuffix + ' got new task', {inputData: data.inputData});
    const inputData = data.inputData;
    const args = inputData.args;
    const outputIsJson = inputData.outputIsJson === 'true';
    const scriptExpression = inputData.scriptExpression;
    try {
      const {stdout, stderr} = await executeFn(scriptExpression, args);
      await createTaskResult(
        outputIsJson,
        {result: stdout},
        stderr,
        updater.complete,
      );
    } catch (e) {
      logger.error('Task has failed', {
        killed: e.killed,
        code: e.code,
        signal: e.signal,
        cmd: e.cmd,
      });
      logger.debug('Task has failed', {error: e});
      let reasonForIncompletion = 'Unknown reason';
      if (e.killed) {
        reasonForIncompletion = 'Timeout';
      } else if (e.code != null && e.code != 0) {
        reasonForIncompletion = 'Exited with error ' + e.code;
      }
      await createTaskResult(
        outputIsJson,
        {result: e.stdout},
        e.stderr,
        updater.fail,
        reasonForIncompletion,
      );
    }
  });
}

async function registerTaskDefs() {
  const taskDefs = [
    {
      name: 'GLOBAL___js',
      type: 'SIMPLE',
      retryCount: 3,
      retryLogic: 'FIXED',
      retryDelaySeconds: 10,
      timeoutSeconds: 300,
      timeoutPolicy: 'TIME_OUT_WF',
      responseTimeoutSeconds: 180,
      ownerEmail: 'example@example.com',
    },
    {
      name: 'GLOBAL___py',
      type: 'SIMPLE',
      retryCount: 3,
      retryLogic: 'FIXED',
      retryDelaySeconds: 10,
      timeoutSeconds: 300,
      timeoutPolicy: 'TIME_OUT_WF',
      responseTimeoutSeconds: 180,
      ownerEmail: 'example@example.com',
    },
  ];
  await conductorClient.registerTaskDefs(taskDefs);
}

async function init() {
  await checkWasmer();

  await registerTaskDefs();

  const workers = new Map([
    ['js', {healthCheckFn: quickJsHealthCheck, executeFn: executeQuickJs}],
    ['py', {healthCheckFn: pythonHealthCheck, executeFn: executePython}],
  ]);

  for (const [wasmSuffix, {healthCheckFn, executeFn}] of workers) {
    try {
      await checkAndRegister(wasmSuffix, healthCheckFn, executeFn);
    } catch (error) {
      logger.warn('Error in checkAndRegister of ' + wasmSuffix, {error});
    }
  }
}

init();
