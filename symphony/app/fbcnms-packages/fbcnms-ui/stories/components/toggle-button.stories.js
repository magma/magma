/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import BubbleChartIcon from '@material-ui/icons/BubbleChart';
import ListAltIcon from '@material-ui/icons/ListAlt';
import React, {useState} from 'react';
import ToggleButtonGroup from '../../components/design-system/ToggleButton/ToggleButtonGroup';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  toggleButtonContainer: {
    marginBottom: '24px',
  },
}));

const TablesRoot = () => {
  const classes = useStyles();
  const [selectedTextButtonId, setSelectedTextButtonId] = useState('0');
  const [selectedIconButtonId, setSelectedIconButtonId] = useState('0');

  return (
    <div className={classes.root}>
      <div className={classes.toggleButtonContainer}>
        <ToggleButtonGroup
          buttons={[
            {id: '0', item: 'Status'},
            {id: '1', item: 'Technician'},
            {id: '2', item: 'Project'},
          ]}
          selectedButtonId={selectedTextButtonId}
          onItemClicked={id => setSelectedTextButtonId(id)}
        />
      </div>
      <ToggleButtonGroup
        buttons={[
          {id: '0', item: <ListAltIcon />},
          {id: '1', item: <BubbleChartIcon />},
        ]}
        selectedButtonId={selectedIconButtonId}
        onItemClicked={id => setSelectedIconButtonId(id)}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add(
  'ToggleButtonGroup',
  () => <TablesRoot />,
);
