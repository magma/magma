/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
const {execFile} = require('child_process');
import logging from '@fbcnms/logging';
const logger = logging.getLogger(module);

const wasmerPath =
  process.env.WASMER_PATH ||
  '/app/fbcnms-projects/workflows-wasm-workers/.wasmer/bin/wasmer';
const enableTracing = process.env.TRACING == 'true' || false;
const maximumWasmerTimeoutMillis: number =
  parseInt(process.env.MAX_WASMER_TIMEOUT_MS) || 10000;

type WasmerError = {
  stdout: string,
  stderr: string,
  cause: child_process$Error,
  code: number | string | null,
  killed?: boolean,
  signal?: string | null,
};

export function executeWasmer(wasmerArgs: string[], stdin: ?string) {
  return new Promise<{stdout: string, stderr: string}>((resolve, reject) => {
    wasmerArgs.unshift('run', '--backend', 'cranelift');
    if (enableTracing) {
      console.debug('executeWasmer');
      wasmerArgs.forEach(it => console.debug(it));
    }

    const options = {
      timeout: maximumWasmerTimeoutMillis,
      killSignal: 'SIGKILL',
    };
    const child = execFile(
      wasmerPath,
      wasmerArgs,
      options,
      (error: ?child_process$Error, stdout, stderr) => {
        if (error) {
          logger.error('Rejecting execution', {error});
          const wasmerError: WasmerError = {
            stdout: stdout.toString(),
            stderr: stderr.toString(),
            code: error.code,
            killed: error.killed,
            signal: error.signal,
            cause: error,
          };
          reject(wasmerError);
        }
        resolve({stdout: stdout.toString(), stderr: stderr.toString()});
      },
    );
    if (stdin != null) {
      child.stdin.write(stdin);
      child.stdin.end();
    }
  });
}

export async function checkWasmer() {
  // will end with rejected promise if exit code != 0
  const {stdout} = await executeWasmer(['--version']);
  logger.info('Wasmer version: ' + stdout);
}
