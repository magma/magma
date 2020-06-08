/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export type TextVariant =
  | 'h1'
  | 'h2'
  | 'h3'
  | 'h4'
  | 'h5'
  | 'subtitle1'
  | 'subtitle2'
  | 'body1'
  | 'body2'
  | 'body3'
  | 'code'
  | 'button'
  | 'caption'
  | 'overline';

// TODO: Load 'Inter' and 'Fira Code' fonts

export default {
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
      fontWeight: 400,
      fontSize: '56px',
      lineHeight: 1.33,
      letterSpacing: '-1px',
    },
    h2: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 300,
      fontSize: '48px',
      lineHeight: 1.33,
      letterSpacing: '0.5px',
    },
    h3: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '34px',
      lineHeight: 1.33,
      letterSpacing: '0.25px',
    },
    h4: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '24px',
      lineHeight: 1.33,
      letterSpacing: '0.15px',
    },
    h5: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 400,
      fontSize: '20px',
      lineHeight: 'normal',
      letterSpacing: '0.15px',
    },
    subtitle1: {
      fontFamily: '"Inter", sans-serif',
      fontWeight: 500,
      fontSize: '16px',
      lineHeight: 1.4,
      letterSpacing: '0.15px',
    },
    subtitle2: {
      fontFamily: '"Inter", sans-serif',
      fontWeight: 700,
      fontSize: '14px',
      lineHeight: 1.71,
    },
    body1: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 500,
      fontSize: '18px',
      lineHeight: 1.5,
      letterSpacing: '0.5px',
    },
    body2: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 500,
      fontSize: '14px',
      lineHeight: 1.43,
      letterSpacing: '0.25px',
    },
    body3: {
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: 500,
      fontSize: '12px',
      lineHeight: 1,
    },
    code: {
      fontFamily: '"Fira Code", sans-serif',
      fontWeight: 500,
      fontSize: '12px',
      lineHeight: 1,
    },
    button: {
      fontFamily: '"Inter", sans-serif',
      fontWeight: 600,
      fontSize: '12px',
      lineHeight: 1.33,
      letterSpacing: '0.5px',
    },
    caption: {
      fontFamily: '"Inter", sans-serif',
      fontWeight: 700,
      fontSize: '12px',
      lineHeight: 1.33,
      letterSpacing: '0.8px',
    },
    overline: {
      fontFamily: '"Inter", sans-serif',
      fontWeight: 500,
      fontSize: '12px',
      lineHeight: 0.66,
      letterSpacing: '0.4px',
    },
  },
};
