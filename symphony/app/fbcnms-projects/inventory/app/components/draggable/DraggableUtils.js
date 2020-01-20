/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export function reorder<T>(
  items: $ReadOnlyArray<T>,
  startIndex: number,
  endIndex: number,
): Array<T> {
  const movedItem = items[startIndex];
  const newArray = [...items];
  newArray.splice(startIndex, 1);
  newArray.splice(endIndex, 0, movedItem);
  return newArray;
}

export const sortByIndex = (
  a: $ReadOnly<{index?: ?number}>,
  b: $ReadOnly<{index?: ?number}>,
) => (a.index ?? 0) - (b.index ?? 0);
