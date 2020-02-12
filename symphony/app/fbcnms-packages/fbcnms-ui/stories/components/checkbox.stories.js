/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Checkbox from '../../components/design-system/Checkbox/Checkbox';
import React, {useState} from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
}));

const TablesRoot = () => {
  const classes = useStyles();
  const [checked, setChecked] = useState(true);
  const [checkedIndeterminate, setCheckedIndeterminate] = useState(false);
  return (
    <div className={classes.root}>
      <Checkbox
        checked={checked}
        onChange={selection =>
          setChecked(selection === 'checked' ? true : false)
        }
      />
      <Checkbox
        checked={checkedIndeterminate}
        indeterminate={!checkedIndeterminate}
        onChange={selection =>
          setCheckedIndeterminate(selection === 'checked' ? true : false)
        }
      />
      <Checkbox checked={false} />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Checkbox', () => (
  <TablesRoot />
));
