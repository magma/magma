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

const ConductorClient = require('conductor-client').default;
const util = require('util');
const execFile = util.promisify(require('child_process').execFile);
import logging from '@fbcnms/logging';
const logger = logging.getLogger(module);
// properties
const conductorApiUrl =
  process.env.CONDUCTOR_API_URL || 'http://conductor-server:8080/api';
const maxRunner = process.env.MAX_RUNNER || 1;
const wasmerPath = process.env.WASMER_PATH || '/root/.wasmer/bin/wasmer';
const quickJsPath = process.env.QUICKJS_PATH || 'wasm/quickjs/quickjs.wasm';
const maximumWasmerTimeoutMillis = process.env.MAX_WASMER_TIMEOUT_MS || 10000;
//

const conductorClient = new ConductorClient({
  baseURL: conductorApiUrl,
});

async function executeWasmer(wasmerArgs) {
  logger.info('executeWasmer', {wasmerArgs});
  try {
    return await execFile(wasmerPath, wasmerArgs, {
      timeout: maximumWasmerTimeoutMillis,
    });
  } catch (e) {
    logger.warn('executeWasmer failed', {wasmerArgs, e});
    throw e;
  }
}

async function checkWasmer() {
  // will end with rejected promise if exit code != 0
  const {stdout} = await executeWasmer(['--version']);
  logger.info('Wasmer version: ' + stdout);
}

function argsToJsonArray(args) {
  if (!Array.isArray(args)) {
    if (typeof args !== 'string') {
      // serialize it to a string
      args = JSON.stringify(args);
    }
    args = [args];
  }
  // 0-th argument is the program name
  args.unshift('script');
  return JSON.stringify(args);
}

async function executeQuickJs(script, args) {
  const preamble =
    `const process = {argv:${argsToJsonArray(args)}};\n` +
    `console.error = function(...args) { std.err.puts(args.join(' '));std.err.puts('\\n'); }\n`;
  script = preamble + script;
  const wasmerArgs = ['run', quickJsPath, '--', '--std', '-e', script];
  try {
    const {stdout, stderr} = await executeWasmer(wasmerArgs);
    logger.info('executeQuickJs succeeded', {stdout, stderr});
    return {stdout, stderr};
  } catch (e) {
    logger.warn('executeQuickJs failed', {script, args, e});
    throw e;
  }
}

async function quickJsHealthCheck() {
  const {stdout, stderr} = await executeQuickJs(
    `console.log('stdout');console.error('stderr');`,
    [],
  );
  if (stdout == 'stdout\n' && stderr == 'stderr\n') {
    return true;
  }
  logger.warn('Unexpected healthcheck result', {stdout, stderr});
  return false;
}

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

async function init() {
  await checkWasmer();
  if (!(await quickJsHealthCheck())) {
    logger.warn('QuickJs healthcheck failed');
  }

  // TODO conductorClient.registerTaskDefs(taskDefs)

  registerWasmWorker('js', async (data, updater) => {
    logger.info('Got new task', {inputData: data.inputData});
    const inputData = data.inputData;
    const args = inputData.args;
    const outputIsJson = inputData.outputIsJson === 'true';
    const script = inputData.script;
    try {
      const {stdout, stderr} = await executeQuickJs(script, args);
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
        reasonForIncompletion = 'Exited with error';
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

init();
