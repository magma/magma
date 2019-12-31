/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React, {useState} from 'react';
import SelectMenu from '../../components/design-system/ContexualLayer/SelectMenu';
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
    label: 'Option 1',
    value: '1',
  },
  {
    label: 'Option 2',
    value: '2',
  },
  {
    label: 'Option 3',
    value: '3',
  },
  {
    label: 'Option 4',
    value: '4',
  },
  {
    label: 'Option 5',
    value: '5',
  },
  {
    label: 'Option 6',
    value: '6',
  },
  {
    label: 'Option 7',
    value: '7',
  },
  {
    label: 'Option 8',
    value: '8',
  },
];

const SelectMenuRoot = () => {
  const classes = useStyles();
  const [options, setOptions] = useState(INITIAL_OPTIONS);

  return (
    <div className={classes.root}>
      <SelectMenu
        className={classes.select}
        label="Project"
        size="normal"
        options={[
          {
            label: 'Option 1',
            value: '1',
          },
          {
            label: 'Option 2',
            value: '2',
          },
        ]}
        onChange={value => window.alert(`Click option #${value}`)}
      />
      <SelectMenu
        className={classes.select}
        size="normal"
        label="Project"
        searchable={true}
        onOptionsFetchRequested={searchTerm =>
          setOptions(
            searchTerm === ''
              ? INITIAL_OPTIONS
              : INITIAL_OPTIONS.filter(option =>
                  String(option.label)
                    .toLowerCase()
                    .includes(searchTerm.toLowerCase()),
                ),
          )
        }
        options={options}
        onChange={value => window.alert(`Click option #${value}`)}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Select Menu', () => (
  <SelectMenuRoot />
));
