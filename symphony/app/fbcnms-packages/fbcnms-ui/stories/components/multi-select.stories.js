/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import MultiSelect from '../../components/design-system/Select/MultiSelect';
import React, {useMemo, useState} from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  select: {
    minWidth: '120px',
    marginBottom: '20px',
  },
}));

const INITIAL_OPTIONS = [
  {
    key: 'option_1',
    label: 'Option 1',
    value: 'wow1',
  },
  {
    key: 'option_2',
    label: 'Option 2',
    value: 'wow2',
  },
];

const BasicMultiSelect = () => {
  const classes = useStyles();
  const [selectedValues, setSelectedValues] = useState([]);
  return (
    <MultiSelect
      className={classes.select}
      label="Project"
      options={INITIAL_OPTIONS}
      onChange={option =>
        setSelectedValues(
          selectedValues.map(v => v.value).includes(option.value)
            ? selectedValues.filter(v => v.value !== option.value)
            : [...selectedValues, option],
        )
      }
      selectedValues={selectedValues}
    />
  );
};

const MultiSelectWithSearch = () => {
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
      <MultiSelect
        className={classes.select}
        label="Project"
        options={options}
        searchable={true}
        onOptionsFetchRequested={searchTerm =>
          setOptions(
            filterBySearchTerm(
              [...selectedOptions, ...unselectedOptions],
              searchTerm,
            ),
          )
        }
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

const MultiSelectsRoot = () => {
  return (
    <div>
      <BasicMultiSelect />
      <MultiSelectWithSearch />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Multi Select', () => (
  <MultiSelectsRoot />
));
