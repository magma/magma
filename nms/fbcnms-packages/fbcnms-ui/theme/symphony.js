/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {hexToRgb} from '../utils/displayUtils';

export const DARK = {
  D10: '#F5F7FC',
  D50: '#EDF0F9',
  D100: '#D2DAE7',
  D200: '#B8C2D3',
  D300: '#9DA9BE',
  D400: '#8895AD',
  D500: '#73839E',
  D600: '#64748C',
  D700: '#536074',
  D800: '#434D5E',
  D900: '#303846',
};

export const BLUE = {
  B50: '#E4F2FF',
  B100: '#BDDEFF',
  B200: '#93C9FF',
  B300: '#66B4FF',
  B400: '#48A3FF',
  B500: '#3593FF',
  B600: '#3984FF',
  B700: '#3A71EA',
  B800: '#3A5FD7',
  B900: '#383DB7',
};

export default {
  palette: {
    primary: BLUE.B600,
    secondary: DARK.D900,
    ...DARK,
    ...BLUE,
    white: '#FFFFFF',
    background: DARK.D10,
    disabled: `rgba(${hexToRgb(DARK.D900)},0.38)`,
    R600: '#FA383E',
    G600: '#00AF5B',
    Y600: '#FFB63E',
    separator: 'rgba(0, 0, 0, 0.12)',
  },
  shadows: {
    DP1: '0px 1px 4px 0px rgba(0, 0, 0, 0.17)',
    DP2: '0px 2px 8px 1px rgba(0, 0, 0, 0.14)',
    DP3: '0px 3px 20px 0px rgba(0, 0, 0, 0.21)',
    DP4:
      '0px 5px 5px -3px rgba(0, 0, 0, 0.2),0px 3px 14px 2px rgba(0, 0, 0, 0.12),0px 8px 10px 1px rgba(0, 0, 0, 0.14)',
  },
};
