/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import MenuItemPhoto from '../../components/design-system/Select/MenuItemPhoto';
import React, {useState} from 'react';
import SelectMenu from '../../components/design-system/Select/SelectMenu';
import {DeleteIcon, EditIcon} from '../../components/design-system/Icons';
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
];

const SelectMenuRoot = () => {
  const classes = useStyles();
  const [options, setOptions] = useState(INITIAL_OPTIONS);

  return (
    <div className={classes.root}>
      <SelectMenu
        className={classes.select}
        size="normal"
        options={[
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
            label: 'Option with a long long long label',
            value: '3',
          },
          {
            key: 'option_4',
            label: 'Option with icon',
            value: '4',
            leftAux: {
              type: 'icon',
              icon: EditIcon,
            },
          },
          {
            key: 'option_5',
            label: 'User Name',
            value: '5',
            leftAux: {
              type: 'node',
              node: (
                <MenuItemPhoto src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNc/R8AAlsBrPwu9ZMAAAAASUVORK5CYII=" />
              ),
            },
            secondaryText: 'User role',
          },
          {
            key: 'option_6',
            label: 'Disabled option',
            value: '6',
            disabled: true,
          },
          {
            key: 'option_7',
            label: 'Warning option',
            value: '7',
            skin: 'red',
            leftAux: {
              type: 'icon',
              icon: DeleteIcon,
            },
          },
        ]}
        onChange={value => window.alert(`Click option #${value}`)}
      />
      <SelectMenu
        className={classes.select}
        size="normal"
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
