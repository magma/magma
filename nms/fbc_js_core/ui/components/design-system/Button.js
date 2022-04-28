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

import typeof SvgIcon from '@material-ui/core/@@SvgIcon';
import type {TRefFor} from './types/TRefFor.flow';

import * as React from 'react';
import Text from '../../../../app/theme/design-system/Text';
import classNames from 'classnames';
import {colors} from '../../../../app/theme/default';
import {joinNullableStrings} from '../../../../fbc_js_core/util/strings';
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
  cometSkin: {},
  disabled: {},
  containedVariant: {
    color: 'white',
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
    '&$cometSkin': {
      backgroundColor: colors.primary.comet,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.primary.white,
          fill: colors.primary.white,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.primary.white,
          fill: colors.primary.white,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.secondary.mariner,
          fill: colors.secondary.mariner,
        },
      },
    },
    '&$primarySkin': {
      backgroundColor: colors.secondary.dodgerBlue,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.primary.white,
          fill: colors.primary.white,
        },
      },
      '&:hover:not($disabled)': {
        backgroundColor: colors.secondary.mariner,
      },
      '&:active:not($disabled)': {
        backgroundColor: colors.secondary.mariner,
      },
    },
    '&$redSkin': {
      backgroundColor: colors.state.error,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.primary.white,
          fill: colors.primary.white,
        },
      },
      '&:hover:not($disabled)': {
        backgroundColor: colors.state.errorAlt,
      },
      '&:active:not($disabled)': {
        backgroundColor: colors.state.errorAlt,
      },
    },
    '&$regularSkin': {
      backgroundColor: colors.primary.white,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.primary.brightGray,
          fill: colors.primary.brightGray,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.primary.comet,
          fill: colors.secondary.dodgerBlue,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.secondary.mariner,
          fill: colors.secondary.mariner,
        },
      },
    },
    '&$graySkin': {
      backgroundColor: colors.primary.selago,
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.primary.brightGray,
          fill: colors.primary.brightGray,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.secondary.dodgerBlue,
          fill: colors.secondary.dodgerBlue,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.secondary.mariner,
          fill: colors.secondary.mariner,
        },
      },
    },
    '&$disabled': {
      cursor: 'default',
      backgroundColor: colors.primary.gullGray,
      '& $buttonText, $icon': {
        color: colors.primary.white,
        fill: colors.primary.white,
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
          color: colors.secondary.dodgerBlue,
          fill: colors.secondary.dodgerBlue,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.secondary.mariner,
          fill: colors.secondary.mariner,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.secondary.mariner,
          fill: colors.secondary.mariner,
        },
      },
    },
    '&$redSkin': {
      '&:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.state.error,
          fill: colors.state.error,
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
          color: colors.primary.brightGray,
          fill: colors.primary.brightGray,
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
          color: colors.primary.brightGray,
          fill: colors.primary.brightGray,
        },
      },
      '&:hover:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.secondary.dodgerBlue,
          fill: colors.secondary.dodgerBlue,
        },
      },
      '&:active:not($disabled)': {
        '& $buttonText, $icon': {
          color: colors.secondary.dodgerBlue,
          fill: colors.secondary.dodgerBlue,
        },
      },
    },
    '&$disabled': {
      cursor: 'default',
      '& $buttonText, $icon': {
        color: colors.primary.gullGray,
        fill: colors.primary.gullGray,
      },
    },
  },
}));

export type ButtonVariant = 'contained' | 'text';
export type ButtonSkin = 'primary' | 'regular' | 'red' | 'gray' | 'comet';

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
