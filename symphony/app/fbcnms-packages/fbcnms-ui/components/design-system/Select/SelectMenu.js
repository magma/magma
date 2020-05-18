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
import SelectMenuItem from './SelectMenuItem';
import SelectSearchInput from './SelectSearchInput';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useMenuContext} from './MenuContext';

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: symphony.palette.white,
    boxShadow: symphony.shadows.DP1,
    borderRadius: '4px',
    maxHeight: '322px',
    overflowY: 'auto',
  },
  fullWidth: {
    width: '100%',
  },
  normalWidth: {
    width: '236px',
  },
}));

export type OptionProps<TValue> = {|
  key: string,
  label: React.Node,
  searchTerm?: string,
  value: TValue,
  className?: ?string,
  ...PermissionHandlingProps,
|};

type Props<TValue> = {
  className?: string,
  onChange: (value: TValue) => void | (() => void),
  options: Array<OptionProps<TValue>>,
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
      {onOptionsFetchRequested && (
        <SelectSearchInput
          searchTerm={searchTerm}
          onChange={updateSearchTerm}
        />
      )}
      {options
        .map(option => {
          const {
            key,
            label,
            value,
            ignorePermissions,
            hideOnMissingPermissions,
            className,
          } = option;
          return (
            <SelectMenuItem
              key={key}
              label={label}
              value={value}
              className={className}
              ignorePermissions={ignorePermissions}
              hideOnMissingPermissions={hideOnMissingPermissions}
              onClick={value => {
                onChange(value);
                onClose();
              }}
              isSelected={selectedValue === option.value}
            />
          );
        })
        .filter(Boolean)}
    </div>
  );
};

export default SelectMenu;
