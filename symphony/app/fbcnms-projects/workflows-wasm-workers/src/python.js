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
const logger = logging.getLogger(module);
import {argsToJsonArray} from './utils.js';
import {executeWasmer} from './wasmer.js';

const pythonBinPath = process.env.PYTHON_PATH || 'wasm/python/bin/python.wasm';
const pythonLibPath = process.env.PYTHON_LIB_PATH || 'wasm/python/lib';

export async function executePython(script, args) {
  const preamble = `argv = ${argsToJsonArray(args)};\n`;
  script = preamble + script;
  const wasmerArgs = ['run', pythonBinPath, '--mapdir=lib:' + pythonLibPath];
  try {
    const {stdout, stderr} = await executeWasmer(wasmerArgs, script);
    logger.info('executePython succeeded', {stdout, stderr});
    return {stdout, stderr};
  } catch (error) {
    logger.warn('executePython failed', {script, args, error});
    throw error;
  }
}

export async function pythonHealthCheck() {
  try {
    const {stdout, stderr} = await executePython(`print('stdout');`, []);
    if (stdout == 'stdout\n') {
      return true;
    }
    logger.warn('Unexpected healthcheck result', {stdout, stderr});
  } catch (error) {
    logger.warn('Unexpected healthcheck error', {error});
  }
  return false;
}
