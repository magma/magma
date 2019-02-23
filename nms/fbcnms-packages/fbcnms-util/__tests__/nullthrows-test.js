/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import nullthrows from '../nullthrows';

test('valid values', () => {
  const str: ?string = 'string';
  expect(nullthrows(str)).toEqual('string');
  expect(nullthrows(0)).toEqual(0);
});

test('null values', () => {
  const str: ?string = null;
  expect(() => nullthrows(str)).toThrowError('[NullValueError]');
  expect(() => nullthrows(undefined)).toThrowError('[NullValueError]');
});

test('null with message', () => {
  expect(() => nullthrows(null, 'mymessage')).toThrowError(
    '[NullValueError] mymessage',
  );
});

test('nullable types', () => {});
