/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
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

export const RED = {
  R600: '#FA383E',
  R700: '#E03237',
  R800: '#C82C31',
};

export type TextVariant =
  | 'h1'
  | 'h2'
  | 'h3'
  | 'h4'
  | 'h5'
  | 'h6'
  | 'subtitle1'
  | 'subtitle2'
  | 'subtitle3'
  | 'body1'
  | 'body2'
  | 'caption'
  | 'overline';

export default {
  palette: {
    primary: BLUE.B600,
    secondary: DARK.D900,
    ...DARK,
    ...BLUE,
    white: '#FFFFFF',
    background: DARK.D10,
    disabled: `rgba(${hexToRgb(DARK.D900)},0.38)`,
    overlay: `rgba(${hexToRgb(DARK.D900)},0.5)`,
    ...RED,
    G600: '#00AF5B',
    Y600: '#FFB63E',
    separator: 'rgba(0, 0, 0, 0.12)',
    separatorLight: 'rgba(0, 0, 0, 0.06)',
  },
  shadows: {
    DP1: '0px 1px 4px 0px rgba(0, 0, 0, 0.17)',
    DP2: '0px 2px 8px 1px rgba(0, 0, 0, 0.14)',
    DP3: '0px 3px 20px 0px rgba(0, 0, 0, 0.21)',
    DP4:
      '0px 5px 5px -3px rgba(0, 0, 0, 0.2),0px 3px 14px 2px rgba(0, 0, 0, 0.12),0px 8px 10px 1px rgba(0, 0, 0, 0.14)',
  },
  typography: {
    h1: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 300,
      fontSize: '96px',
      lineHeight: 1.33,
      letterSpacing: '-1.5px',
    },
    h2: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 300,
      fontSize: '60px',
      lineHeight: 1.33,
      letterSpacing: '-0.5px',
    },
    h3: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '48px',
      lineHeight: 1.33,
      letterSpacing: '0px',
    },
    h4: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '32px',
      lineHeight: 1.25,
      letterSpacing: '0.25px',
    },
    h5: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '24px',
      lineHeight: 1.33,
      letterSpacing: '0px',
    },
    h6: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 500,
      fontSize: '20px',
      lineHeight: 'normal',
      letterSpacing: '0.25px',
    },
    body1: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '16px',
      lineHeight: 1.5,
      letterSpacing: '0.15px',
    },
    body2: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '14px',
      lineHeight: 1.43,
      letterSpacing: '0.25px',
    },
    subtitle1: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 500,
      fontSize: '16px',
      lineHeight: 1.5,
      letterSpacing: '0.15px',
    },
    subtitle2: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 500,
      fontSize: '14px',
      lineHeight: 1.71,
      letterSpacing: '0.1px',
    },
    subtitle3: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 500,
      fontSize: '14px',
      lineHeight: 1.14,
      letterSpacing: '1.25px',
      textTransform: 'uppercase',
    },
    caption: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '12px',
      lineHeight: 1.33,
      letterSpacing: 'normal',
    },
    overline: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 500,
      fontSize: '12px',
      lineHeight: 1.33,
      letterSpacing: '1px',
      textTransform: 'uppercase',
    },
  },
};
