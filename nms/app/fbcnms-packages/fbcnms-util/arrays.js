/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export function removeItem<T>(
  input: $ReadOnlyArray<T>,
  index: number,
): Array<T> {
  const newArray = [...input];
  newArray.splice(index, 1);
  return newArray;
}

export function setItem<TItem>(
  input: $ReadOnlyArray<TItem>,
  index: number,
  value: TItem,
): Array<TItem> {
  const newArray = [...input];
  newArray[index] = value;
  return newArray;
}

/**
 * Given an array of dicts, updates the property key at the
 * index to the value
 */
export function updateItem<TItem: {}, TProp: $Keys<TItem>>(
  input: $ReadOnlyArray<TItem>,
  index: number,
  prop: TProp,
  value: $ElementType<TItem, TProp>,
): Array<TItem> {
  const newArray = [...input];
  newArray[index] = {
    ...newArray[index],
    [prop]: value,
  };
  return newArray;
}
