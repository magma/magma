/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TRefFor} from './types/TRefFor.flow';

import * as React from 'react';
import Text from './Text';
import classNames from 'classnames';
import symphony from '../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    border: 0,
    cursor: 'pointer',
    '&:focus': {
      outline: 'none',
    },
    display: 'inline-flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
  },
  icon: {},
  hasIcon: {
    justifyContent: 'start',
    '& $buttonText': {
      flexGrow: 1,
    },
  },
  rightIcon: {
    alignSelf: 'flex-end',
    marginLeft: '6px',
  },
  leftIcon: {
    alignSelf: 'flex-start',
    marginRight: '6px',
  },
  hasRightIcon: {
    '& $buttonText': {
      textAlign: 'left',
    },
  },
  hasLeftIcon: {
    '& $buttonText': {
      textAlign: 'right',
    },
  },
  primarySkin: {},
  redSkin: {},
  orangeSkin: {},
  greenSkin: {},
  regularSkin: {},
  graySkin: {},
  disabled: {},
  containedVariant: {
    height: '32px',
    minWidth: '88px',
    padding: '4px 18px',
    borderRadius: '4px',
    '&$hasRightIcon': {
      padding: '4px 6px 4px 12px',
    },
    '&$hasLeftIcon': {
      padding: '4px 12px 4px 6px',
    },
    '&$primarySkin': {
      backgroundColor: symphony.palette.primary,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.white,
        },
      },
      '&:hover:not($disabled)': {
        backgroundColor: symphony.palette.B700,
      },
      '&:active:not($disabled)': {
        backgroundColor: symphony.palette.B800,
      },
    },
    '&$redSkin': {
      backgroundColor: symphony.palette.R600,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.white,
        },
      },
      '&:hover:not($disabled)': {
        backgroundColor: symphony.palette.R700,
      },
      '&:active:not($disabled)': {
        backgroundColor: symphony.palette.R800,
      },
    },
    '&$orangeSkin': {
      backgroundColor: symphony.palette.Y600,
      '&:not($disabled) $buttonText': {
        color: symphony.palette.white,
      },
    },
    '&$greenSkin': {
      backgroundColor: symphony.palette.G600,
      '&:not($disabled) $buttonText': {
        color: symphony.palette.white,
      },
    },
    '&$regularSkin': {
      backgroundColor: symphony.palette.white,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.secondary,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.B700,
        },
      },
    },
    '&$graySkin': {
      backgroundColor: symphony.palette.background,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.secondary,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.B700,
        },
      },
    },
    '&$disabled': {
      cursor: 'default',
      backgroundColor: symphony.palette.disabled,
      '& $buttonText, $icon': {
        color: symphony.palette.white,
      },
    },
  },
  buttonText: {},
  textVariant: {
    display: 'inline-flex',
    textAlign: 'left',
    background: 'none',
    padding: 0,
    '&$primarySkin': {
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          opacity: 0.75,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          opacity: 0.75,
        },
      },
    },
    '&$redSkin': {
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.R600,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          opacity: 0.75,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          opacity: 0.75,
        },
      },
    },
    '&$regularSkin': {
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.secondary,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          opacity: 0.75,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          opacity: 0.75,
        },
      },
    },
    '&$graySkin': {
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.secondary,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          opacity: 0.75,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          opacity: 0.75,
        },
      },
    },
    '&$disabled': {
      cursor: 'default',
      '& $buttonText, $icon': {
        color: symphony.palette.disabled,
      },
    },
  },
}));

export type ButtonVariant = 'contained' | 'text';
export type ButtonSkin =
  | 'primary'
  | 'regular'
  | 'red'
  | 'gray'
  | 'orange'
  | 'green';
type SvgIcon = React$ComponentType<SvgIconExports>;

export type ButtonProps = {|
  skin?: ButtonSkin,
  variant?: ButtonVariant,
  disabled?: boolean,
|};

export type Props = {
  className?: string,
  children: React.Node,
  onClick?: void | (() => void | Promise<void>),
  leftIcon?: SvgIcon,
  rightIcon?: SvgIcon,
  tooltip?: string,
  ...ButtonProps,
};

const Button = (props: Props, forwardedRef: TRefFor<HTMLButtonElement>) => {
  const {
    className,
    children,
    skin = 'primary',
    disabled = false,
    variant = 'contained',
    onClick,
    leftIcon: LeftIcon = null,
    rightIcon: RightIcon = null,
    tooltip,
  } = props;
  const classes = useStyles();

  return (
    <button
      className={classNames(
        classes.root,
        classes[`${skin}Skin`],
        classes[`${variant}Variant`],
        {
          [classes.disabled]: disabled,
          [classes.hasIcon]: LeftIcon != null || RightIcon != null,
          [classes.hasLeftIcon]: LeftIcon != null,
          [classes.hasRightIcon]: RightIcon != null,
        },
        className,
      )}
      type="button"
      title={tooltip}
      disabled={disabled}
      onClick={onClick}
      ref={forwardedRef}>
      {LeftIcon ? (
        <LeftIcon
          className={classNames(classes.icon, classes.leftIcon)}
          size="small"
          color="inherit"
        />
      ) : null}
      <Text variant="body2" weight="medium" className={classes.buttonText}>
        {children}
      </Text>
      {RightIcon ? (
        <RightIcon
          className={classNames(classes.icon, classes.rightIcon)}
          size="small"
          color="inherit"
        />
      ) : null}
    </button>
  );
};

export default React.forwardRef<Props, HTMLButtonElement>(Button);
