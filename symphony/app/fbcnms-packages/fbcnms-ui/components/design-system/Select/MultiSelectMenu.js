/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {OptionProps} from './SelectMenu';

import * as React from 'react';
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

const useStyles = makeStyles({
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
});

type Props<TValue> = {
  className?: string,
  onChange: (option: OptionProps<TValue>) => void | (() => void),
  options: Array<OptionProps<TValue>>,
  searchable?: boolean,
  onOptionsFetchRequested?: (searchTerm: string) => void,
  selectedValues: Array<OptionProps<TValue>>,
  size?: 'normal' | 'full',
};

type MultiSelectMenuItemsProps<TValue> = {
  options: Array<OptionProps<TValue>>,
  isSelected: (option: OptionProps<TValue>) => boolean,
  onChange: (option: OptionProps<TValue>) => void | (() => void),
};

const MultiSelectMenuItems = <TValue>({
  options,
  isSelected,
  onChange,
}: MultiSelectMenuItemsProps<TValue>) =>
  options.map(option => (
    <SelectMenuItem
      key={`option_${String(option.value)}`}
      label={option.label}
      value={option.value}
      onClick={() => {
        onChange(option);
      }}
      isSelected={isSelected(option)}
    />
  ));

const MultiSelectMenu = <TValue>({
  className,
  options,
  onChange,
  selectedValues,
  size = 'full',
  searchable = false,
  onOptionsFetchRequested,
}: Props<TValue>) => {
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
