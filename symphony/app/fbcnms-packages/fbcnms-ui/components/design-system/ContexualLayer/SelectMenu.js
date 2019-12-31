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
import ClearIcon from '@material-ui/icons/Clear';
import InputAffix from '../Input/InputAffix';
import Text from '../Text';
import TextInput from '../Input/TextInput';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
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
  fullWidth: {
    width: '100%',
  },
  normalWidth: {
    width: '236px',
  },
  input: {
    padding: '16px',
  },
  clearIconContainer: {
    backgroundColor: symphony.palette.background,
    padding: '6px',
    borderRadius: '100%',
    width: '20px',
    height: '20px',
    boxSizing: 'border-box',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  },
  clearIcon: {
    color: symphony.palette.D800,
    fontSize: '13.66px',
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
  searchable?: boolean,
  onOptionsFetchRequested?: (searchTerm: string) => void,
  selectedValue?: ?TValue,
  size?: 'normal' | 'full',
};

const SelectMenu = <TValue>({
  className,
  options,
  onChange,
  selectedValue,
  size = 'full',
  searchable = false,
  onOptionsFetchRequested,
}: Props<TValue>) => {
  const classes = useStyles();
  const {onClose} = useMenuContext();
  const [searchTerm, setSearchTerm] = useState('');

  const updateSearchTerm = useCallback(
    searchTerm => {
      setSearchTerm(searchTerm.toLowerCase());
      onOptionsFetchRequested && onOptionsFetchRequested(searchTerm);
    },
    [onOptionsFetchRequested],
  );

  return (
    <div
      className={classNames(classes.root, className, {
        [classes.fullWidth]: size === 'full',
        [classes.normalWidth]: size === 'normal',
      })}>
      {searchable && (
        <TextInput
          className={classes.input}
          type="string"
          placeholder="Type to filter..."
          onChange={({target}) => updateSearchTerm(target.value)}
          value={searchTerm}
          suffix={
            searchTerm ? (
              <InputAffix
                onClick={() => updateSearchTerm('')}
                className={classes.clearIconContainer}>
                <ClearIcon className={classes.clearIcon} />
              </InputAffix>
            ) : null
          }
        />
      )}
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
