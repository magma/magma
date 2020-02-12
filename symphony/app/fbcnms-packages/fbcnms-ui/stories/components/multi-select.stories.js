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
import Text from '../../components/design-system/Text';
import symphony from '../../theme/symphony';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  card: {
    marginBottom: '16px',
  },
  popover: {
    backgroundColor: symphony.palette.white,
    boxShadow: symphony.shadows.DP1,
    borderRadius: '4px',
  },
  select: {
    minWidth: '120px',
  },
}));

const Popover = () => {
  const classes = useStyles();
  return (
    <div className={classes.popover}>
      <Text variant="body2">
        Below the input, with the same width. Amazing.
      </Text>
    </div>
  );
};

const INITIAL_OPTIONS = [
  {
    label: 'Option 1',
    value: 'wow1',
  },
  {
    label: 'Option 2',
    value: 'wow2',
  },
];

const MultiSelectsRoot = () => {
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
      String(option.label)
        .toLowerCase()
        .includes(searchTerm.toLowerCase()),
    );

  return (
    <div className={classes.root}>
      <MultiSelect
        className={classes.select}
        popover={Popover}
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

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Multi Select', () => (
  <MultiSelectsRoot />
));
