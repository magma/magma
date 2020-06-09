/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {TRefFor} from './types/TRefFor.flow';

import * as React from 'react';
import Text from './Text';
import classNames from 'classnames';
import symphony from '../../theme/symphony';
import {joinNullableStrings} from '@fbcnms/util/strings';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from './Form/FormElementContext';
import {useMemo} from 'react';

const useStyles = makeStyles(_theme => ({
  root: {
    border: 0,
    cursor: 'pointer',
    '&:focus': {
      outline: 'none',
    },
    flexShrink: 0,
    display: 'inline-flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    maxWidth: '100%',
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
    marginLeft: '8px',
  },
  leftIcon: {
    alignSelf: 'flex-start',
    marginRight: '8px',
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
  secondaryGraySkin: {},
  disabled: {},
  containedVariant: {
    height: '32px',
    minWidth: '88px',
    padding: '4px 12px',
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
          fill: symphony.palette.white,
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
          fill: symphony.palette.white,
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
          fill: symphony.palette.secondary,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
          fill: symphony.palette.primary,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.B700,
          fill: symphony.palette.B700,
        },
      },
    },
    '&$graySkin': {
      backgroundColor: symphony.palette.background,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.secondary,
          fill: symphony.palette.secondary,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
          fill: symphony.palette.primary,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.B700,
          fill: symphony.palette.B700,
        },
      },
    },
    '&$secondaryGraySkin': {
      backgroundColor: symphony.palette.background,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.D500,
          fill: symphony.palette.D500,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.D900,
          fill: symphony.palette.D900,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.B700,
          fill: symphony.palette.B700,
        },
      },
    },
    '&$disabled': {
      cursor: 'default',
      backgroundColor: symphony.palette.disabled,
      '& $buttonText, $icon': {
        color: symphony.palette.white,
        fill: symphony.palette.white,
      },
    },
  },
  buttonText: {
    maxHeight: '100%',
  },
  textVariant: {
    display: 'inline-flex',
    textAlign: 'left',
    background: 'none',
    padding: 0,
    height: '24px',
    maxWidth: '100%',
    '&$primarySkin': {
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
          fill: symphony.palette.primary,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.B700,
          fill: symphony.palette.B700,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.B700,
          fill: symphony.palette.B700,
        },
      },
    },
    '&$redSkin': {
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.R600,
          fill: symphony.palette.R600,
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
          fill: symphony.palette.secondary,
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
          color: symphony.palette.D500,
          fill: symphony.palette.D500,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
          fill: symphony.palette.primary,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
          fill: symphony.palette.primary,
        },
      },
    },
    '&$secondaryGraySkin': {
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.D500,
          fill: symphony.palette.D500,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.D900,
          fill: symphony.palette.D900,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: symphony.palette.primary,
          fill: symphony.palette.primary,
        },
      },
    },
    '&$disabled': {
      cursor: 'default',
      '& $buttonText, $icon': {
        color: symphony.palette.disabled,
        fill: symphony.palette.disabled,
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
  | 'secondaryGray'
  | 'orange'
  | 'green';
type SvgIcon = React$ComponentType<SvgIconExports>;

export type ButtonProps = $ReadOnly<{|
  skin?: ButtonSkin,
  variant?: ButtonVariant,
  useEllipsis?: ?boolean,
  disabled?: boolean,
  tooltip?: string,
|}>;

export type Props = $ReadOnly<{|
  className?: string,
  children: React.Node,
  onClick?:
    | void
    | (void | ((SyntheticMouseEvent<HTMLElement>) => void | Promise<void>)),

  leftIcon?: SvgIcon,
  leftIconClass?: string,
  rightIcon?: SvgIcon,
  rightIconClass?: string,
  ...ButtonProps,
|}>;

const Button = (props: Props, forwardedRef: TRefFor<HTMLButtonElement>) => {
  const {
    className,
    children,
    skin = 'primary',
    disabled: disabledProp = false,
    variant = 'contained',
    useEllipsis = true,
    onClick,
    leftIcon: LeftIcon = null,
    leftIconClass = null,
    rightIcon: RightIcon = null,
    rightIconClass = null,
    tooltip: tooltipProp,
  } = props;
  const classes = useStyles();

  const {
    disabled: contextDisabled,
    tooltip: contextTooltip,
  } = useFormElementContext();

  const disabled = useMemo(() => disabledProp || contextDisabled, [
    disabledProp,
    contextDisabled,
  ]);

  const tooltip = useMemo(
    () => joinNullableStrings([tooltipProp, contextTooltip]),
    [contextTooltip, tooltipProp],
  );

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
          color="inherit"
          className={classNames(classes.icon, classes.leftIcon, leftIconClass)}
          size="small"
        />
      ) : null}
      <Text
        variant="body2"
        weight="medium"
        useEllipsis={useEllipsis}
        className={classes.buttonText}>
        {children}
      </Text>
      {RightIcon ? (
        <RightIcon
          className={classNames(
            classes.icon,
            classes.rightIcon,
            rightIconClass,
          )}
          size="small"
          color="inherit"
        />
      ) : null}
    </button>
  );
};

export default React.forwardRef<Props, HTMLButtonElement>(Button);
