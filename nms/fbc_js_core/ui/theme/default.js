/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import SymphonyTheme from './symphony';
import {
  blue05,
  blue30,
  blue60,
  blueGrayDark,
  brightGray,
  comet,
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
import {createTheme} from '@material-ui/core/styles';

export default createTheme({
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

    magmalte: {
      appbar: '#323845',
      background: '#171B25',
    },
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
        color: brightGray,
        backgroundColor: white,
        minHeight: '56px',
        borderColor: white,
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
        color: comet,
        fontFamily: '"Inter", sans-serif',
        fontWeight: 600,
        fontSize: '12px',
        lineHeight: 1.33,
        letterSpacing: '0.5px',
        lineHeight: '14px',
        paddingBottom: '12px',
        paddingTop: '12px',
        height: '24px',
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
    MuiDialogContent: {
      root: {
        padding: '0 32px',
      },
    },
  },
});
