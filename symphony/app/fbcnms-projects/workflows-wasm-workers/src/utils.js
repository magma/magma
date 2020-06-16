/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export function argsToJsonArray(args) {
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
