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

const quickJsPath = process.env.QUICKJS_PATH || 'wasm/quickjs/quickjs.wasm';

export async function executeQuickJs(script: string, args: string[]) {
  const preamble = `
const argv = ${argsToJsonArray(args)};
console.error = function(...args) { std.err.puts(args.join(' '));std.err.puts('\\n'); }
`;
  script = preamble + script;
  const wasmerArgs = ['run', quickJsPath, '--', '--std', '-e', script];
  try {
    const {stdout, stderr} = await executeWasmer(wasmerArgs);
    logger.info('executeQuickJs succeeded', {stdout, stderr});
    return {stdout, stderr};
  } catch (error) {
    logger.warn('executeQuickJs failed', {script, args, error});
    throw error;
  }
}

export async function quickJsHealthCheck() {
  try {
    const {stdout, stderr} = await executeQuickJs(
      `console.log('stdout');console.error('stderr');`,
      [],
    );
    if (stdout == 'stdout\n' && stderr == 'stderr\n') {
      return true;
    }
    logger.warn('Unexpected healthcheck result', {stdout, stderr});
  } catch (error) {
    logger.warn('Unexpected healthcheck error', {error});
  }
  return false;
}
