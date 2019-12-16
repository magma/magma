/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import SelectMenu from '../../components/design-system/ContexualLayer/SelectMenu';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  select: {
    width: '150px',
  },
}));

const SelectMenuRoot = () => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <SelectMenu
        className={classes.select}
        label="Project"
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
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Select Menu', () => (
  <SelectMenuRoot />
));
