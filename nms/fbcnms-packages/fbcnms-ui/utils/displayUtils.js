/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export const sortLexicographically = (a: string, b: string) =>
  a.localeCompare(b, 'en', {numeric: true});
