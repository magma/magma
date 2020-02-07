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

const SelectsRoot = () => {
  const classes = useStyles();
  const [selectedValue, setSelectedValue] = useState(null);

  return (
    <div className={classes.root}>
      <Select
        className={classes.select}
        popover={Popover}
        label="Project"
        options={[
          {
            label: 'Option 1',
            value: 'wow1',
          },
          {
            label: 'Option 2',
            value: 'wow2',
          },
        ]}
        selectedValue={selectedValue}
        onChange={value => setSelectedValue(value)}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Select', () => (
  <SelectsRoot />
));
