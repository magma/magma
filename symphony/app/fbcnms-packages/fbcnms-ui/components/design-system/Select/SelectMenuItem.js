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
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  option: {
    display: 'flex',
    alignItems: 'center',
    padding: '6px 16px',
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

type Props<TValue> = {
  label: React.Node,
  value: TValue,
  onClick: (value: TValue) => void,
  isSelected?: boolean,
};

const SelectMenuItem = <TValue>({
  label,
  value,
  onClick,
  isSelected = false,
}: Props<TValue>) => {
  const classes = useStyles();
  return (
    <div className={classes.option} onClick={() => onClick(value)}>
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
};

export default SelectMenuItem;
