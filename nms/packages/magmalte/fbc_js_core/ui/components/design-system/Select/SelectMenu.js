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
  key: string | number,
|}>;

export type SelectMenuProps<TValue> = $ReadOnly<{|
  onChange: (value: TValue) => void | (() => void),
  options: $ReadOnlyArray<OptionProps<TValue>>,
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
              disabled={option.disabled === true}
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
