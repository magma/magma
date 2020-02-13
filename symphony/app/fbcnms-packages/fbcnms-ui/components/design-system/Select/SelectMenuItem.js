/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionHandlingProps} from '../Form/FormAction';

import * as React from 'react';
import CheckIcon from '@material-ui/icons/Check';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormElementContext from '@fbcnms/ui/components/design-system/Form/FormElementContext';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  option: {
    display: 'flex',
    alignItems: 'center',
    padding: '6px 16px',
    cursor: 'pointer',
    whiteSpace: 'nowrap',
    '&:not($disabled)&:hover': {
      backgroundColor: symphony.palette.B50,
    },
  },
  disabled: {
    opacity: 0.5,
  },
  label: {
    flexGrow: 1,
  },
  checkIcon: {
    marginLeft: '6px',
    color: symphony.palette.primary,
  },
});

type Props<TValue> = {|
  label: React.Node,
  value: TValue,
  onClick: (value: TValue) => void,
  isSelected?: boolean,
  className?: ?string,
  ...PermissionHandlingProps,
|};

const SelectMenuItem = <TValue>({
  label,
  value,
  onClick,
  isSelected = false,
  hideWhenDisabled = false,
  className,
  ...permissionHandlingProps
}: Props<TValue>) => {
  const classes = useStyles();
  return (
    <FormAction
      {...permissionHandlingProps}
      hideWhenDisabled={hideWhenDisabled}>
      <FormElementContext.Consumer>
        {context => {
          const disabled = context.disabled;
          return (
            <div
              className={classNames(classes.option, className, {
                [classes.disabled]: disabled,
              })}
              onClick={disabled ? null : () => onClick(value)}>
              {typeof label === 'string' ? (
                <Text className={classes.label} variant="body2">
                  {label}
                </Text>
              ) : (
                label
              )}
              {isSelected && (
                <CheckIcon className={classes.checkIcon} fontSize="small" />
              )}
            </div>
          );
        }}
      </FormElementContext.Consumer>
    </FormAction>
  );
};

export default SelectMenuItem;
