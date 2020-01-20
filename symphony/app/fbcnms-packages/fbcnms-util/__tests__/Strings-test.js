/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {base64ToHex, hexToBase64} from '../strings';

describe('base64ToHex', () => {
  test('converted to hex', () => {
    expect(base64ToHex('0w==')).toEqual('d3');
    expect(base64ToHex('0/w=')).toEqual('d3fc');
  });
  test('invalid input stripped', () => {
    expect(base64ToHex('0^^z==')).toEqual('d3');
  });
  test('empty input', () => {
    expect(base64ToHex('')).toEqual('');
  });
});

describe('hexToBase64', () => {
  test('converted to base64', () => {
    expect(hexToBase64('d3')).toEqual('0w==');
    expect(hexToBase64('d3fc')).toEqual('0/w=');
  });
  test('odd number', () => {
    expect(hexToBase64('0d')).toEqual(hexToBase64('d'));
    expect(hexToBase64('0d')).toEqual(hexToBase64('d'));
  });
  test('invalid input raises exception', () => {
    expect(() => hexToBase64('d$$3')).toThrow();
    expect(() => hexToBase64('d$3zzZ')).toThrow();
  });
  test('upper/lowercase doesnt matter', () => {
    expect(hexToBase64('D3aF')).toEqual(hexToBase64('d3af'));
  });
  test('empty input', () => {
    expect(hexToBase64('')).toEqual('');
  });
});
