/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

class NullValueError extends Error {
  constructor(message?: string) {
    super('[NullValueError]' + (message ? ' ' + message : ''));
  }
}

export default function nullthrows<TVal>(data: ?TVal, message?: string): TVal {
  if (data == null) {
    throw new NullValueError(message);
  }
  return data;
}
