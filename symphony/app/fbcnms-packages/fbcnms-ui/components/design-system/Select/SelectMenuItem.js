/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  ErrorHandlingProps,
  PermissionHandlingProps,
} from '../Form/FormAction';

import * as React from 'react';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormElementContext from '@fbcnms/ui/components/design-system/Form/FormElementContext';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  option: {
    display: 'flex',
    alignItems: 'center',
    padding: '8px 16px',
    cursor: 'pointer',
    '&:not($disabled)&:hover': {
      backgroundColor: symphony.palette.background,
    },
    '&$optionWithLeftAux': {
      paddingLeft: '12px',
      paddingTop: '6px',
      paddingBottom: '6px',
    },
  },
  optionWithLeftAux: {},
  disabled: {
    opacity: 0.38,
    cursor: 'not-allowed',
  },
  label: {
    flexGrow: 1,
  },
  checkIcon: {
    marginLeft: '6px',
    color: symphony.palette.primary,
  },
  leftAux: {
    display: 'inline-flex',
    marginRight: '8px',
  },
  contentContainer: {
    display: 'flex',
    flexDirection: 'column',
  },
}));

export type MenuItemLeftAux = $ReadOnly<
  | {|
      type: 'icon',
      icon: React$ComponentType<SvgIconExports>,
    |}
  | {
      type: 'node',
      node: React.Node,
    },
>;

export type SelectMenuItemBaseProps<TValue> = $ReadOnly<{|
  label: React.Node,
  value: TValue,
  isSelected?: boolean,
  className?: ?string,
  leftAux?: MenuItemLeftAux,
  secondaryText?: React.Node,
  disabled?: boolean,
  skin?: 'regular' | 'red',
  ...PermissionHandlingProps,
  ...ErrorHandlingProps,
|}>;

type Props<TValue> = $ReadOnly<{|
  ...SelectMenuItemBaseProps<TValue>,
  onClick: (value: TValue) => void,
|}>;

const SelectMenuItem = <TValue>({
  label,
  value,
  onClick,
  isSelected = false,
  hideOnMissingPermissions = false,
  className,
  leftAux,
  secondaryText,
  skin = 'regular',
  disabled: disabledProp = false,
  ...actionProps
}: Props<TValue>) => {
  const classes = useStyles();
  const LeftIcon = leftAux?.type === 'icon' ? leftAux.icon : null;
  const coercedSkin = disabledProp
    ? 'regular'
    : skin === 'red'
    ? 'error'
    : skin;
  return (
    <FormAction
      {...actionProps}
      disabled={disabledProp}
      hideOnMissingPermissions={hideOnMissingPermissions}>
      <FormElementContext.Consumer>
        {({disabled}) => {
          return (
            <div
              className={classNames(classes.option, className, {
                [classes.disabled]: disabled,
                [classes.optionWithLeftAux]: leftAux != null,
              })}
              onClick={disabled ? null : () => onClick(value)}>
              {leftAux != null && (
                <div className={classes.leftAux}>
                  {leftAux.type === 'icon'
                    ? LeftIcon != null && (
                        <LeftIcon
                          color={isSelected ? 'primary' : coercedSkin}
                          size="small"
                        />
                      )
                    : leftAux.node}
                </div>
              )}
              <div className={classes.contentContainer}>
                <Text
                  className={classes.label}
                  variant="body2"
                  color={isSelected ? 'primary' : coercedSkin}>
                  {label}
                </Text>
                {secondaryText != null && (
                  <Text color="gray" variant="caption">
                    {secondaryText}
                  </Text>
                )}
              </div>
            </div>
          );
        }}
      </FormElementContext.Consumer>
    </FormAction>
  );
};

export default SelectMenuItem;
