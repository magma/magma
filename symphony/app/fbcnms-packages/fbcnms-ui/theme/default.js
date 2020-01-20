/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import SymphonyTheme from './symphony';
import {
  blue05,
  blue30,
  blue60,
  blueGrayDark,
  fadedBlue,
  gray0,
  gray00,
  gray1,
  gray13,
  gray50,
  primaryText,
  red,
  redwood,
  white,
} from './colors';
import {createMuiTheme} from '@material-ui/core/styles';

export default createMuiTheme({
  symphony: SymphonyTheme,
  palette: {
    primary: {
      light: SymphonyTheme.palette.B300,
      main: SymphonyTheme.palette.B600,
      dark: SymphonyTheme.palette.B900,
    },
    secondary: {
      main: SymphonyTheme.palette.D900,
    },
    action: {
      hover: '#cfd8dc',
      selected: '#f2f3f5',
    },
    grey: {
      '50': '#e4f0f6',
      '100': '#eaeff2',
      '200': '#cfd8dc',
      '300': '#b9cad2',
      '400': '#8f9ea4',
      '500': '#89a1ac',
      '600': '#606770',
      '700': '#455a64',
      '800': '#4d4d4e',
      '900': '#263238',
      A100: '#ecf3ff',
      A200: '#8d949e',
      A700: '#444950',
    },
    typography: {
      ...SymphonyTheme.typography,
    },
    red: red,
    redwood: redwood,
    dark: '#1d2129',
    fadedBlue: fadedBlue,
    blueGrayDark: blueGrayDark,
    blue05: blue05,
    blue30: blue30,
    blue60: blue60,
    gray00: gray00,
    gray1: gray1,
    gray50: gray50,
    gray13,
    primaryText: primaryText,
  },
  overrides: {
    MuiAppBar: {
      colorPrimary: {
        backgroundColor: blue60,
        color: white,
      },
    },
    MuiButton: {
      root: {
        borderRadius: 4,
        cursor: 'pointer',
        fontSize: '14px',
        textTransform: 'capitalize',
        padding: '8px 12px',
        fontWeight: 500,
        lineHeight: '16px',
      },
      contained: {
        boxShadow: 'none',
      },
    },
    MuiFormControl: {
      marginDense: {
        marginTop: '0px',
        marginBottom: '0px',
      },
    },
    MuiToggleButtonGroup: {
      '&$selected': {
        boxShadow: 'none',
        borderRadius: 4,
        border: `1px solid ${blue60}`,
      },
    },
    MuiToggleButton: {
      root: {
        color: blue60,
        backgroundColor: white,
        textTransform: 'none',
        '&$selected': {
          color: white,
          backgroundColor: blue60,
        },
      },
    },
    MuiListItemText: {
      root: {
        marginBottom: '0px',
        marginTop: '0px',
      },
    },
    MuiSelect: {
      selectMenu: {
        height: '24px',
      },
    },
    MuiTableRow: {
      root: {
        backgroundColor: 'white',
      },
    },
    MuiIconButton: {
      root: {
        color: blue60,
      },
    },
    MuiAvatar: {
      colorDefault: {
        backgroundColor: '#e4f0f6',
        color: blue60,
      },
    },
    MuiInputLabel: {
      outlined: {
        transform: 'translate(14px, 16px) scale(1)',
      },
    },
    MuiOutlinedInput: {
      root: {
        '&$notchedOutline': {
          borderColor: '#CCD0D5',
        },
        '&$focused $notchedOutline': {
          borderColor: 'rgba(0, 0, 0, 0.87)',
          borderWidth: '1px',
        },
        '&$disabled': {
          background: gray0,
        },
      },
      input: {
        fontSize: '14px',
        lineHeight: '14px',
        paddingBottom: '12px',
        paddingTop: '12px',
        height: '24px',
        '&::-webkit-input-placeholder': {
          opacity: 1,
        },
        '&::-moz-placeholder': {
          opacity: 1,
        },
        '&::-ms-input-placeholder': {
          opacity: 1,
        },
      },
      inputMarginDense: {
        paddingTop: '9px',
        paddingBottom: '9px',
        fontSize: '14px',
        lineHeight: '14px',
        height: '14px',
        '&::placeholder': {
          color: '#8895ad',
        },
        '&::-webkit-input-placeholder': {
          opacity: 1,
        },
        '&::-moz-placeholder': {
          opacity: 1,
        },
        '&::-ms-input-placeholder': {
          opacity: 1,
        },
      },
    },
  },
});
