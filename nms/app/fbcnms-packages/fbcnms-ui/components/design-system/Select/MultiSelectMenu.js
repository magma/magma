/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MenuItemLeftAux} from './SelectMenuItem';
import type {OptionProps} from './SelectMenu';

import * as React from 'react';
import CheckBoxIcon from '@material-ui/icons/CheckBox';
import CheckBoxOutlineBlankIcon from '@material-ui/icons/CheckBoxOutlineBlank';
import SelectMenuItem from './SelectMenuItem';
import SelectSearchInput from './SelectSearchInput';
import Text from '../Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '../../../theme/symphony';
import useVerticalScrollingEffect from '../hooks/useVerticalScrollingEffect';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useMemo, useRef} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: symphony.palette.white,
    boxShadow: symphony.shadows.DP1,
    borderRadius: '4px',
  },
  itemsContainer: {
    overflowY: 'auto',
    maxHeight: '250px',
  },
  fullWidth: {
    width: '100%',
  },
  normalWidth: {
    width: '236px',
  },
  selectedTitle: {
    display: 'block',
    padding: '0px 16px',
    marginBottom: '8px',
  },
  separator: {
    borderTop: `1px solid ${symphony.palette.D50}`,
    margin: '8px 16px 12px 16px',
  },
  checkedIcon: {
    color: symphony.palette.primary,
  },
  uncheckedIcon: {
    color: symphony.palette.D200,
  },
}));

export type MultiSelectOptionProps<TValue> = $ReadOnly<
  $Rest<OptionProps<TValue>, {|leftAux?: MenuItemLeftAux|}>,
>;

export type MultiSelectMenuProps<TValue> = $ReadOnly<{|
  className?: string,
  onChange: (option: MultiSelectOptionProps<TValue>) => void | (() => void),
  options: Array<MultiSelectOptionProps<TValue>>,
  searchable?: boolean,
  onOptionsFetchRequested?: (searchTerm: string) => void,
  selectedValues: Array<MultiSelectOptionProps<TValue>>,
  size?: 'normal' | 'full',
|}>;

type MultiSelectMenuItemsProps<TValue> = $ReadOnly<{|
  options: Array<MultiSelectOptionProps<TValue>>,
  isSelected: (option: MultiSelectOptionProps<TValue>) => boolean,
  onChange: (option: MultiSelectOptionProps<TValue>) => void | (() => void),
|}>;

const MultiSelectMenuItems = <TValue>({
  options,
  isSelected,
  onChange,
}: MultiSelectMenuItemsProps<TValue>) => {
  const classes = useStyles();
  return options.map(option => {
    const isSelectedValue = isSelected(option);
    const Icon = isSelectedValue ? CheckBoxIcon : CheckBoxOutlineBlankIcon;
    return (
      <SelectMenuItem
        key={`option_${String(option.value)}`}
        label={option.label}
        value={option.value}
        onClick={() => {
          onChange(option);
        }}
        leftAux={{
          type: 'node',
          node: (
            <Icon
              className={
                isSelectedValue ? classes.checkedIcon : classes.uncheckedIcon
              }
              size="small"
            />
          ),
        }}
        isSelected={isSelectedValue}
      />
    );
  });
};

const MultiSelectMenu = <TValue>({
  className,
  options,
  onChange,
  selectedValues,
  size = 'full',
  searchable = false,
  onOptionsFetchRequested,
}: MultiSelectMenuProps<TValue>) => {
  const thisElement = useRef(null);
  const classes = useStyles();
  const [searchTerm, setSearchTerm] = useState('');
  useVerticalScrollingEffect(thisElement);

  const updateSearchTerm = useCallback(
    searchTerm => {
      setSearchTerm(searchTerm.toLowerCase());
      onOptionsFetchRequested && onOptionsFetchRequested(searchTerm);
    },
    [onOptionsFetchRequested],
  );

  const unselectedOptions = useMemo(
    () =>
      options.filter(
        option =>
          !selectedValues.map(value => value.value).includes(option.value),
      ),
    [options, selectedValues],
  );

  return (
    <div
      className={classNames(classes.root, className, {
        [classes.fullWidth]: size === 'full',
        [classes.normalWidth]: size === 'normal',
      })}>
      {searchable && (
        <SelectSearchInput
          searchTerm={searchTerm}
          onChange={updateSearchTerm}
        />
      )}
      <div className={classes.itemsContainer} ref={thisElement}>
        {searchable && selectedValues.length > 0 && (
          <div>
            <Text
              className={classes.selectedTitle}
              variant="caption"
              weight="bold">
              <fbt desc="Title shown above the selected items">Selected</fbt>
            </Text>
            <MultiSelectMenuItems
              options={selectedValues}
              isSelected={() => true}
              onChange={onChange}
            />
            <div
              className={classes.separator}
              hidden={unselectedOptions.length === 0}
            />
          </div>
        )}
        <MultiSelectMenuItems
          options={searchable ? unselectedOptions : options}
          isSelected={option =>
            selectedValues.map(v => v.value).includes(option.value)
          }
          onChange={onChange}
        />
      </div>
    </div>
  );
};

export default MultiSelectMenu;
