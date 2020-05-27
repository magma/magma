/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {MenuBaseProps} from './MenuBase';
import type {SelectMenuItemBaseProps} from './SelectMenuItem';

import * as React from 'react';
import MenuBase from './MenuBase';
import SelectMenuItem from './SelectMenuItem';
import SelectSearchInput from './SelectSearchInput';
import {useCallback, useState} from 'react';
import {useMenuContext} from './MenuContext';

export type OptionProps<TValue> = $ReadOnly<{|
  ...SelectMenuItemBaseProps<TValue>,
  key: string,
|}>;

export type SelectMenuProps<TValue> = $ReadOnly<{|
  onChange: (value: TValue) => void | (() => void),
  options: Array<OptionProps<TValue>>,
  onOptionsFetchRequested?: (searchTerm: string) => void,
  selectedValue?: ?TValue,
  ...MenuBaseProps,
|}>;

const SelectMenu = <TValue>(props: SelectMenuProps<TValue>) => {
  const {
    options,
    onChange,
    selectedValue,
    onOptionsFetchRequested,
    ...menuBaseProps
  } = props;
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
    <MenuBase {...menuBaseProps}>
      {onOptionsFetchRequested && (
        <SelectSearchInput
          searchTerm={searchTerm}
          onChange={updateSearchTerm}
        />
      )}
      {options
        .map(option => {
          const {key, label, value, ...menuItemProps} = option;
          return (
            <SelectMenuItem
              key={key}
              label={label}
              value={value}
              onClick={value => {
                onChange(value);
                onClose();
              }}
              isSelected={selectedValue === option.value}
              {...menuItemProps}
            />
          );
        })
        .filter(Boolean)}
    </MenuBase>
  );
};

export default SelectMenu;
