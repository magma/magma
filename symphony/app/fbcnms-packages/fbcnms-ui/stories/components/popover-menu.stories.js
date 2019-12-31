/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AddIcon from '@material-ui/icons/Add';
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';
import PopoverMenu from '../../components/design-system/ContexualLayer/PopoverMenu';
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles({
  root: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
  },
  popoverMenu: {
    marginRight: '16px',
  },
  moreIcon: {
    padding: '6px',
    backgroundColor: 'white',
    borderRadius: '100%',
    cursor: 'pointer',
  },
});

const PopoverMenuRoot = () => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <PopoverMenu
        className={classes.popoverMenu}
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
        onChange={value => window.alert(`Clicked on item #${value}`)}>
        <MoreHorizIcon className={classes.moreIcon} />
      </PopoverMenu>
      <PopoverMenu
        variant="contained"
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
        onChange={value => window.alert(`Clicked on item #${value}`)}
        rightIcon={AddIcon}>
        Add Filter
      </PopoverMenu>
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Popover Menu', () => (
  <PopoverMenuRoot />
));
