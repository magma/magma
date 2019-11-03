/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

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
  primarySkin: {},
  redSkin: {},
  regularSkin: {},
  disabled: {},
  containedVariant: {
    height: '32px',
    minWidth: '88px',
    padding: '4px 18px',
    borderRadius: '4px',
    '&$primarySkin': {
      backgroundColor: symphony.palette.primary,
      '&:not($disabled) $buttonText': {
        color: symphony.palette.white,
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
      '&:not($disabled) $buttonText': {
        color: symphony.palette.white,
      },
      '&:hover:not($disabled)': {
        backgroundColor: symphony.palette.R700,
      },
      '&:active:not($disabled)': {
        backgroundColor: symphony.palette.R800,
      },
    },
    '&$regularSkin': {
      backgroundColor: symphony.palette.white,
      '&:not($disabled) $buttonText': {
        color: symphony.palette.secondary,
      },
      '&:hover:not($disabled) $buttonText': {
        color: symphony.palette.primary,
      },
      '&:active:not($disabled) $buttonText': {
        color: symphony.palette.B700,
      },
    },
    '&$disabled': {
      cursor: 'default',
      backgroundColor: symphony.palette.disabled,
      '& $buttonText': {
        color: symphony.palette.white,
      },
    },
  },
  buttonText: {},
  textVariant: {
    display: 'inline-block',
    background: 'none',
    padding: 0,
    '&$primarySkin': {
      '&:not($disabled) $buttonText': {
        color: symphony.palette.primary,
      },
      '&:hover:not($disabled) $buttonText': {
        opacity: 0.75,
      },
      '&:active:not($disabled) $buttonText': {
        opacity: 0.75,
      },
    },
    '&$redSkin': {
      '&:not($disabled) $buttonText': {
        color: symphony.palette.R600,
      },
      '&:hover:not($disabled) $buttonText': {
        opacity: 0.75,
      },
      '&:active:not($disabled) $buttonText': {
        opacity: 0.75,
      },
    },
    '&$regularSkin': {
      '&:not($disabled) $buttonText': {
        color: symphony.palette.secondary,
      },
      '&:hover:not($disabled) $buttonText': {
        opacity: 0.75,
      },
      '&:active:not($disabled) $buttonText': {
        opacity: 0.75,
      },
    },
    '&$disabled': {
      cursor: 'default',
      '& $buttonText': {
        color: symphony.palette.disabled,
      },
    },
  },
}));

type Props = {
  className?: string,
  children: React.Node,
  onClick?: void | (() => void | Promise<void>),
  skin: 'primary' | 'regular' | 'red',
  variant: 'contained' | 'text',
  disabled: boolean,
};

const Button = (props: Props) => {
  const {className, children, skin, disabled, variant, onClick} = props;
  const classes = useStyles();
  const textifiedChildren = Array.isArray(children) ? (
    children.map(c =>
      typeof c === 'string' ? (
        <Text variant="body2" weight="medium" className={classes.buttonText}>
          {c}
        </Text>
      ) : (
        c
      ),
    )
  ) : typeof children === 'string' ? (
    <Text variant="body2" weight="medium" className={classes.buttonText}>
      {children}
    </Text>
  ) : (
    children
  );

  return (
    <button
      className={classNames(
        classes.root,
        classes[`${skin}Skin`],
        classes[`${variant}Variant`],
        {[classes.disabled]: disabled},
        className,
      )}
      disabled={disabled}
      onClick={onClick}>
      {textifiedChildren}
    </button>
  );
};

Button.defaultProps = {
  skin: 'primary',
  disabled: false,
  variant: 'contained',
};

export default Button;
