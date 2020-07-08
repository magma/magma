/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

// Stringify object twice so that it can be dropped into source code.
// E.g. input is {"lambdaValue":"a"}
// result is "{\"lambdaValue\":\"a\"}"
// This is useful for embedding the value inside script to be evaluated:
// const scriptString = 'JSON.parse(' + result + ')';
// eval(scriptString)
export function escapeJson(inputData: mixed) {
  const firstJson = JSON.stringify(inputData) || '"{}"';
  const result = JSON.stringify(firstJson);
  return result;
}
