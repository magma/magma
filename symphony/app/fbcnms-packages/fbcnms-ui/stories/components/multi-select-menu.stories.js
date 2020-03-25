/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import MultiSelectMenu from '../../components/design-system/Select/MultiSelectMenu';
import React, {useMemo, useState} from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  select: {
    margin: '16px',
  },
}));

const INITIAL_OPTIONS = [
  {
    key: 'option_1',
    label: 'Option 1',
    value: '1',
  },
  {
    key: 'option_2',
    label: 'Option 2',
    value: '2',
  },
  {
    key: 'option_3',
    label: 'Option 3',
    value: '3',
  },
  {
    key: 'option_4',
    label: 'Option 4',
    value: '4',
  },
  {
    key: 'option_5',
    label: 'Option 5',
    value: '5',
  },
  {
    key: 'option_6',
    label: 'Option 6',
    value: '6',
  },
  {
    key: 'option_7',
    label: 'Option 7',
    value: '7',
  },
  {
    key: 'option_8',
    label: 'Option 8',
    value: '8',
  },
  {
    key: 'option_9',
    label: 'Option 9',
    value: '9',
  },
  {
    key: 'option_10',
    label: 'Option 10',
    value: '10',
  },
];

const SelectMenuRoot = () => {
  const classes = useStyles();
  const [options, setOptions] = useState(INITIAL_OPTIONS);
  const [selectedValues, setSelectedValues] = useState([]);
  const selectedOptions = useMemo(
    () =>
      INITIAL_OPTIONS.filter(option =>
        selectedValues.map(v => v.value).includes(option.value),
      ),
    [selectedValues],
  );
  const unselectedOptions = useMemo(
    () =>
      INITIAL_OPTIONS.filter(
        option => !selectedValues.map(v => v.value).includes(option.value),
      ),
    [selectedValues],
  );
  const filterBySearchTerm = (options, searchTerm) =>
    options.filter(option =>
      String(option.label).toLowerCase().includes(searchTerm.toLowerCase()),
    );

  return (
    <div className={classes.root}>
      <MultiSelectMenu
        className={classes.select}
        size="normal"
        label="Project"
        searchable={true}
        onOptionsFetchRequested={searchTerm =>
          setOptions(
            filterBySearchTerm(
              [...selectedOptions, ...unselectedOptions],
              searchTerm,
            ),
          )
        }
        options={options}
        onChange={option =>
          setSelectedValues(
            selectedValues.map(v => v.value).includes(option.value)
              ? selectedValues.filter(v => v.value !== option.value)
              : [...selectedValues, option],
          )
        }
        selectedValues={selectedValues}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add(
  'Multi Select Menu',
  () => <SelectMenuRoot />,
);
