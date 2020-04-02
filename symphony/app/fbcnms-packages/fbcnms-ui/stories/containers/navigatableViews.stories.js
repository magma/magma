/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import NavigatableViews from '../../components/design-system/View/NavigatableViews';
import React from 'react';
import symphony from '../../theme/symphony';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    width: '100%',
    minHeight: '400px',
    maxHeight: 'calc(100vh - 16px)',
  },
  childView: {
    display: 'flex',
    height: '100%',
    backgroundColor: symphony.palette.white,

    height: '2000px',
  },
}));

const ViewHeaderRoot = () => {
  const classes = useStyles();

  const views = [
    {
      menuItem: {
        label: 'Home',
        tooltip: 'Go to home page',
      },
      component: {
        header: {
          title: 'Home',
          subtitle:
            'The Company is a secret group of multinational corporate alliances known only by those who work for them or oppose them. Its influence and power over individuals stretches to the White House, controlling every decision the country makes.',
          actionButtons: [
            {
              title: 'Open Door',
              action: () =>
                console.log(
                  'If this is home, we probably need to be able to open the door...',
                ),
            },
          ],
        },
        children: <div className={classes.childView} />,
      },
    },
    {
      menuItem: {
        label: 'Products',
        tooltip: 'See what we offer',
      },
      component: {
        header: {
          title: 'Products',
          subtitle: 'Ever heard of Sona..?',
          actionButtons: [
            {
              title: 'Purchase',
              action: () => console.log('I want to buy things!!'),
            },
          ],
        },
        children: <div className={classes.childView} />,
      },
    },
    {
      menuItem: {
        label: 'Downloads',
        tooltip: 'Drivers, Guids and stuff',
      },
      component: {
        header: {
          title: 'Downloads',
          actionButtons: [
            {
              title: 'Go Torrent!',
              action: () => console.log('sh...'),
            },
          ],
        },
        children: <div className={classes.childView} />,
      },
    },
    {
      menuItem: {
        label: 'About',
        tooltip: 'Who are we',
      },
      component: {
        header: {
          title: 'About',
          actionButtons: [
            {
              title: 'Contact',
              action: () => console.log('Send some mail..'),
            },
          ],
        },
        children: <div className={classes.childView} />,
      },
    },
  ];

  return (
    <div className={classes.root}>
      <NavigatableViews header="The Company" views={views} />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.CONTAINERS}`, module).add(
  'NavigatableViews',
  () => <ViewHeaderRoot />,
);
