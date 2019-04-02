/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {createMuiTheme} from '@material-ui/core/styles';
import {gray0, sea100, sea200, sea300, red, redwood, white} from './colors';

export default createMuiTheme({
  palette: {
    primary: {
      light: sea100,
      main: sea200,
      dark: sea300,
    },
    secondary: {
      main: '#606770',
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
    red: red,
    redwood: redwood,
  },
  overrides: {
    MuiAppBar: {
      colorPrimary: {
        backgroundColor: sea200,
        color: white,
      },
    },
    MuiButton: {
      root: {
        borderRadius: 4,
        cursor: 'pointer',
        fontSize: 14,
        padding: '6px 30px',
      },
    },
    MuiToggleButtonGroup: {
      '&$selected': {
        boxShadow: 'none',
        borderRadius: 4,
        border: `1px solid ${sea200}`,
      },
    },
    MuiToggleButton: {
      root: {
        color: sea200,
        backgroundColor: white,
        textTransform: 'none',
        '&$selected': {
          color: white,
          backgroundColor: sea200,
        },
      },
    },
    MuiIconButton: {
      root: {
        color: sea100,
      },
    },
    MuiAvatar: {
      colorDefault: {
        backgroundColor: '#e4f0f6',
        color: sea200,
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
      },
      input: {
        fontSize: '14px',
        lineHeight: '14px',
        paddingBottom: '15px',
        paddingTop: '15px',
      },
      inputMarginDense: {
        paddingTop: '9px',
        paddingBottom: '9px',
        fontSize: '14px',
        lineHeight: '14px',
        height: '14px',
        '&::placeholder': {
          color: 'rgba(0, 0, 0, 0.6)',
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
      '&$disabled': {
        background: gray0,
      },
    },
  },
  typography: {
    useNextVariants: true,
  },
});
