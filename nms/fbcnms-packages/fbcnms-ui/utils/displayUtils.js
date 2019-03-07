/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const KB = 1024;
const MB = 1024 * 1024;
const GB = 1024 * 1024 * 1024;

export const sortLexicographically = (a: string, b: string) =>
  a.localeCompare(b, 'en', {numeric: true});

export const formatFileSize = (sizeInBytes: number) => {
  if (sizeInBytes === 0) {
    return '0MB';
  }

  if (sizeInBytes >= GB) {
    return `${(sizeInBytes / GB).toFixed(2)}GB`;
  } else if (sizeInBytes >= MB) {
    return `${(sizeInBytes / MB).toFixed(2)}MB`;
  } else if (sizeInBytes >= KB) {
    return `${Math.round(sizeInBytes / KB)}KB`;
  } else {
    return `${sizeInBytes}B`;
  }
};
