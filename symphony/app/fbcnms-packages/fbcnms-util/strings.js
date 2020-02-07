/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export function hexToBase64(hexString: string): string {
  let parsedValue;

  parsedValue = hexString.toLowerCase();
  if (parsedValue.length % 2 === 1) {
    parsedValue = '0' + parsedValue;
  }
  // Raise an exception if any bad value is entered
  if (!isValidHex(hexString)) {
    throw new Error('is not valid hex');
  }
  return Buffer.from(parsedValue, 'hex').toString('base64');
}

export function base64ToHex(base64String: string): string {
  return Buffer.from(base64String, 'base64').toString('hex');
}

export function isValidHex(hexString: string): boolean {
  return hexString.match(/^[a-fA-F0-9]*$/) !== null;
}

export function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}
