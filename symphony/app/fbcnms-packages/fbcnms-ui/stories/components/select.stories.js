/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React, {useState} from 'react';
import Select from '../../components/design-system/Select/Select';
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

const SelectsRoot = () => {
  const classes = useStyles();
  const [selectedValue, setSelectedValue] = useState(null);
  const renderSelect = (disabled: boolean = false, size = 'normal') => (
    <Select
      className={classes.select}
      label="Project"
      options={[
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
        {
          key: 'long_option',
          label: 'Option with a long label',
          value: 'wow3',
        },
      ]}
      disabled={disabled}
      selectedValue={selectedValue}
      size={size}
      onChange={value => setSelectedValue(value)}
    />
  );
  return (
    <div className={classes.root}>
      {renderSelect()}
      {renderSelect(false, 'full')}
      {renderSelect(true)}
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Select', () => (
  <SelectsRoot />
));
