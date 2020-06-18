/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export default function renderList(list: Array<string>): string {
  if (!Array.isArray(list)) {
    console.error(
      `renderList(): expected array, received ${list} (${typeof list})`,
    );
    return '';
  }

  if (list.length < 4) {
    return list.join(', ');
  }

  return `${list[0]}, ${list[1]} & ${list.length - 2} others`;
}
