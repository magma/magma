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
 */

//  NOTE: Color Names generated from hex code at http://chir.ag/projects/name-that-color/

import {createTheme} from '@material-ui/core/styles';

export const colors = {
  primary: {
    white: '#FFFFFF',
    selago: '#F4F7FD',
    concrete: '#F2F2F2',
    mercury: '#E5E5E5',
    nobel: '#B3B3B3',
    gullGray: '#9DA7BB',
    comet: '#545F77',
    brightGray: '#323845',
    mirage: '#171B25',
  },
  secondary: {
    malibu: '#88B3F9',
    dodgerBlue: '#3984FF',
    mariner: '#1F5BC4',
  },
  button: {
    lightOutline: '#CCD0DB',
    fill: '#FAFAFB',
  },
  state: {
    positive: '#31BF56',
    positiveAlt: '#229A41',
    error: '#E52240',
    errorAlt: '#B21029',
    errorFill: '#FFF8F9',
    warning: '#F5DD5A',
    warningAlt: '#B69900',
    warningFill: '#FFFCED',
  },
  alerts: {
    severe: '#E52240',
    major: '#E36730',
    minor: '#F5DD5A',
    other: '#88B3F9',
  },
  data: {
    coral: '#FF824B',
    flamePea: '#E36730',
    portage: '#A07EEA',
    studio: '#6649A6',
  },
  code: {
    crusta: '#F76D47',
    pelorous: '#39B6C8',
    electricViolet: '#7D4DFF',
    orchid: '#DA70D6',
    chelseaCucumber: '#91B859',
    candlelight: '#FFD715',
    mischka: '#D4D8DE',
  },
};

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

export const typography = {
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
  h6: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    fontWeight: 400,
    fontSize: '16px',
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
  subtitle3: {
    fontFamily: '"Inter", sans-serif',
    fontWeight: 700,
    fontSize: '12px',
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
};

export const shadows = {
  DP1:
    '0px 0px 2px 0px rgba(0, 0, 0, 0.14), 0px 2px 2px 0px rgba(0, 0, 0, 0.12), 0px 1px 3px 0px rgba(0, 0, 0, 0.20)',
  DP2:
    '0px 2px 4px 0px rgba(0, 0, 0, 0.14), 0px 3px 4px 0px rgba(0, 0, 0, 0.12), 0px 1px 5px 0px rgba(0, 0, 0, 0.20)',
  DP3:
    '0px 2px 4px 0px rgba(0, 0, 0, 0.14), 0px 4px 5px 0px rgba(0, 0, 0, 0.12), 0px 1px 10px 0px rgba(0, 0, 0, 0.20)',
  DP4:
    '0px 6px 10px 0px rgba(0, 0, 0, 0.14), 0px 1px 18px 0px rgba(0, 0, 0, 0.12), 0px 3px 5px 0px rgba(0, 0, 0, 0.20)',
  DP5:
    '0px 24px 38px 0px rgba(0, 0, 0, 0.14), 0px 9px 46px 0px rgba(0, 0, 0, 0.12), 0px 11px 15px 0px rgba(0, 0, 0, 0.20)',
};

export default createTheme({
  palette: {
    primary: {
      light: colors.secondary.malibu,
      main: colors.secondary.dodgerBlue,
      dark: colors.secondary.mariner,
    },
  },
  typography: {
    ...typography,
  },
  props: {
    MuiDialogTitle: {
      disableTypography: true, // No more wrapping children in typography component for dialog title
    },
  },
  overrides: {
    MuiAppBar: {
      colorPrimary: {
        color: colors.primary.white,
      },
    },
    MuiAccordionSummary: {
      content: {
        margin: '0px',
      },
    },
    // @ts-ignore
    MuiAlert: {
      root: {
        marginBottom: '16px',
      },
    },
    MuiBackdrop: {
      root: {
        backgroundColor: `rgba(50, 56, 69, 0.8)`, // colors.primary.brightGray RGB value
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
      containedPrimary: {
        backgroundColor: colors.primary.comet,
        '&:hover, &:focus': {
          backgroundColor: colors.primary.brightGray,
          boxShadow: 'none',
        },
      },
      outlinedPrimary: {
        border: `1px solid ${colors.primary.comet}`,

        backgroundColor: colors.primary.white,
        '&:hover, &:focus': {
          backgroundColor: colors.primary.white,
          boxShadow: 'none',
          border: `1px solid ${colors.primary.comet}`,
        },
        color: colors.primary.comet,
      },
    },
    MuiCheckbox: {
      colorSecondary: {
        '&$checked': {
          color: colors.secondary.dodgerBlue,
        },
      },
    },
    MuiDivider: {
      root: {
        backgroundColor: colors.primary.concrete,
        height: '2px',
      },
    },
    MuiDialog: {
      paper: {
        backgroundColor: colors.primary.concrete,
      },
    },
    MuiDialogActions: {
      root: {
        backgroundColor: colors.primary.white,
        boxShadow: shadows.DP3,
        padding: '20px',
        zIndex: 1,
      },
    },
    MuiDialogContent: {
      root: {
        padding: '32px',
        minHeight: '480px',
      },
    },
    MuiDialogTitle: {
      root: {
        backgroundColor: colors.primary.mirage,
        padding: '16px 24px',
        color: colors.primary.white,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        width: '100%',
      },
    },
    MuiDropzoneArea: {
      root: {
        color: colors.primary.gullGray,
        backgroundColor: colors.primary.concrete,
        marginBottom: '16px',
      },
      icon: {
        color: colors.primary.gullGray,
      },
    },
    MuiFormControl: {
      marginDense: {
        marginTop: '0px',
        marginBottom: '0px',
      },
    },
    MuiFormLabel: {
      root: {
        '&$error': {...typography.overline, color: colors.state.error},
      },
    },
    MuiToggleButton: {
      root: {
        color: colors.secondary.dodgerBlue,
        backgroundColor: colors.primary.white,
        textTransform: 'none',
        '&$selected': {
          color: colors.primary.white,
          backgroundColor: colors.secondary.dodgerBlue,
        },
      },
    },
    MuiListItem: {
      root: {
        paddingTop: '16px',
        paddingBottom: '16px',
      },
    },
    MuiListItemText: {
      root: {
        marginBottom: '0px',
        marginTop: '0px',
        ...typography.body3,
      },
    },
    MuiSelect: {
      root: {
        display: 'flex',
        alignItems: 'center',
      },
      outlined: {
        '&:focus': {
          backgroundColor: colors.primary.white,
        },
        '&:read-only': {
          opacity: 1,
        },
      },
    },
    MuiSwitch: {
      colorSecondary: {
        '&$checked': {
          color: colors.secondary.dodgerBlue,
        },
        '&$checked + .MuiSwitch-track': {
          backgroundColor: colors.secondary.dodgerBlue,
          opacity: 0.5,
        },
      },
    },
    MuiTabs: {
      indicator: {
        height: '4px',
      },
    },
    MuiTab: {
      root: {
        textTransform: 'none',
        ...typography.body1,
        padding: '12px 24px',
      },
    },
    MuiTableRow: {
      root: {
        backgroundColor: 'white',
      },
    },
    MuiIconButton: {
      root: {
        color: colors.secondary.dodgerBlue,
      },
    },
    MuiAvatar: {
      colorDefault: {
        backgroundColor: '#e4f0f6',
        color: colors.secondary.dodgerBlue,
      },
    },
    MuiInputLabel: {
      outlined: {
        transform: 'translate(14px, 16px) scale(1)',
      },
    },
    MuiOutlinedInput: {
      root: {
        minHeight: '36px',
        borderRadius: '2px',
        boxSizing: 'border-box',
        padding: '8px 16px',
        color: colors.primary.brightGray,
        backgroundColor: colors.primary.white,
      },
      multiline: {
        padding: '8px 16px',
      },
      input: {
        padding: 0,
        ...typography.button,
        '&::-webkit-input-placeholder': {
          opacity: 0.5,
        },
        '&::-moz-placeholder': {
          opacity: 0.5,
        },
        '&::-ms-input-placeholder': {
          opacity: 0.5,
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
