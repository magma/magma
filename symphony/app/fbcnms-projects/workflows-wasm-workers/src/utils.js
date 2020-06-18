/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export function argsToJsonArray(args: string[] | string | {}): string {
  let argsArray: string[];
  if (!Array.isArray(args)) {
    let argString: string;
    if (typeof args !== 'string') {
      // serialize it to a string
      argString = JSON.stringify(args) || '';
    } else {
      argString = args;
    }
    argsArray = [argString];
  } else {
    argsArray = args.slice();
  }
  // 0-th argument is the program name
  argsArray.unshift('script');
  return JSON.stringify(argsArray);
}
