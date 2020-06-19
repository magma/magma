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
import {escapeJson} from './utils.js';
import {executeWasmer} from './wasmer.js';

const pythonBinPath = process.env.PYTHON_PATH || 'wasm/python/bin/python.wasm';
const pythonLibPath = process.env.PYTHON_LIB_PATH || 'wasm/python/lib';

function prefixLines(script: string, indent: string) {
  return script
    .split('\n')
    .map(it => indent + it)
    .join('\n');
}

export async function executePython(
  script: string,
  args: string[],
  inputData: mixed,
) {
  const escapedInputDataJson = escapeJson(inputData);
  script = `
import sys,json
inputData = json.loads(${escapedInputDataJson})
def log(*args, **kwargs):
  print(*args, file=sys.stderr, **kwargs)

def script_fun():
${prefixLines(script, '  ')}

result = script_fun()
if not result is None:
  if isinstance(result, str):
    sys.stdout.write(result)
  else:
    sys.stdout.write(json.dumps(result))
`;

  // options:
  // -q: quiet, do not print python version
  // -B: do not write .pyc files on import
  // -c script: execute passed script
  const wasmerArgs = [
    'run',
    pythonBinPath,
    '--mapdir=lib:' + pythonLibPath,
    '--',
    '-B',
    '-q',
    '-c',
    script,
  ];
  try {
    const {stdout, stderr} = await executeWasmer(wasmerArgs);
    logger.info('executePython succeeded', {stdout, stderr});
    return {stdout, stderr};
  } catch (error) {
    logger.warn('executePython failed', {script, args, error});
    throw error;
  }
}

export async function pythonHealthCheck() {
  try {
    const {stdout, stderr} = await executePython(
      `log('stderr');return 'stdout'`,
      [],
      {},
    );
    if (stdout == 'stdout' && stderr == 'stderr\n') {
      return true;
    }
    logger.warn('Unexpected healthcheck result', {stdout, stderr});
  } catch (error) {
    console.error('Unexpected healthcheck error', error);
  }
  return false;
}
