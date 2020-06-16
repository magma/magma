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

const wasmerPath = process.env.WASMER_PATH || '/root/.wasmer/bin/wasmer';
const maximumWasmerTimeoutMillis = process.env.MAX_WASMER_TIMEOUT_MS || 10000;

export function executeWasmer(wasmerArgs, stdin) {
  return new Promise((resolve, reject) => {
    logger.info('executeWasmer', {wasmerArgs, stdin});

    const options = {
      timeout: maximumWasmerTimeoutMillis,
      killSignal: 'SIGKILL',
    };
    const child = execFile(
      wasmerPath,
      wasmerArgs,
      options,
      (error, stdout, stderr) => {
        if (error) {
          logger.error('Rejecting execution', {error});
          reject({stdout, stderr, ...error});
        }
        resolve({stdout, stderr});
      },
    );
    if (stdin != null) {
      child.stdin.setEncoding('utf-8');
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
