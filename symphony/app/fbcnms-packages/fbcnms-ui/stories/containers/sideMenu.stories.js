/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import SideMenu from '../../components/design-system/Menu/SideMenu';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';
import {useState} from 'react';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    height: '100%',
    minHeight: '400px',
    maxHeight: '500px',
  },
  card: {
    marginBottom: '16px',
  },
}));

const SideMenuRoot = () => {
  const classes = useStyles();
  const [selectedItem, setSelectedItem] = useState(0);

  const menuItems = [
    {
      label: 'Home',
      tooltip: 'Go to home page',
    },
    {
      label: 'Products',
      tooltip: 'See what we offer',
    },
    {
      label: 'Downloads',
      tooltip: 'Drivers, Guids and stuff',
    },
    {
      label: 'About',
      tooltip: 'Who are we',
    },
  ];

  return (
    <div className={classes.root}>
      <SideMenu
        header="The Company"
        items={menuItems}
        activeItemIndex={selectedItem}
        onActiveItemChanged={(_item, index) => setSelectedItem(index)}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.CONTAINERS}`, module).add('SideMenu', () => (
  <SideMenuRoot />
));
