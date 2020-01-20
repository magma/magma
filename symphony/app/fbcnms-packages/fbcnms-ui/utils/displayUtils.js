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

export const isJSON = (text: ?string): boolean => {
  if (!text) {
    return false;
  }
  try {
    JSON.parse(text);
  } catch (e) {
    return false;
  }
  return true;
};

// formats server side timestamps (seonds from epoch)
// to text input required format dd-mm-yyyy
export const formatDateForTextInput = (dateValue: ?string) => {
  return !!dateValue ? dateValue.split('T')[0] : '';
};

export const formatMultiSelectValue = (
  options: Array<{value: string, label: string}>,
  value: string,
) => options.find(option => option.value === value)?.label;

export function hexToRgb(hexColor: string) {
  hexColor = hexColor.substr(1);

  const re = new RegExp(`.{1,${hexColor.length / 3}}`, 'g');
  let colors = hexColor.match(re);

  if (colors && colors[0].length === 1) {
    colors = colors.map(n => n + n);
  }

  return colors ? colors.map(n => parseInt(n, 16)).join(',') : '';
}
