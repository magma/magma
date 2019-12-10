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
import CheckIcon from '@material-ui/icons/Check';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useMenuContext} from './MenuContext';

const useStyles = makeStyles({
  root: {
    backgroundColor: symphony.palette.white,
    boxShadow: symphony.shadows.DP1,
    borderRadius: '4px',
  },
  option: {
    display: 'flex',
    alignItems: 'center',
    padding: '6px 8px',
    cursor: 'pointer',
    '&:hover': {
      backgroundColor: symphony.palette.B50,
    },
  },
  label: {
    flexGrow: 1,
  },
  checkIcon: {
    marginLeft: '6px',
    color: symphony.palette.primary,
  },
});

export type OptionProps<TValue> = {
  label: React.Node,
  value: TValue,
};

type Props<TValue> = {
  className?: string,
  onChange: (value: TValue) => void | (() => void),
  options: Array<OptionProps<TValue>>,
  selectedValue?: ?TValue,
};

const SelectMenu = <TValue>({
  className,
  options,
  onChange,
  selectedValue,
}: Props<TValue>) => {
  const classes = useStyles();
  const {onClose} = useMenuContext();
  return (
    <div className={classNames(classes.root, className)}>
      {options.map(option => (
        <div
          className={classes.option}
          onClick={() => {
            onChange(option.value);
            onClose();
          }}>
          {typeof option.label === 'string' ? (
            <Text className={classes.label} variant="body2">
              {option.label}
            </Text>
          ) : (
            option.label
          )}
          {option.value === selectedValue && (
            <CheckIcon className={classes.checkIcon} fontSize="small" />
          )}
        </div>
      ))}
    </div>
  );
};

export default SelectMenu;
