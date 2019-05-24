/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {storiesOf} from '@storybook/react';
import AppSideBar from '../components/layout/AppSideBar.react';
import AssignmentIcon from '@material-ui/icons/Assignment';
import ListIcon from '@material-ui/icons/List';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';

storiesOf('AppSideBar', module)
  .add('default', () => (
    <AppSideBar
      user={{email: 'user@fb.com', isSuperUser: true}}
      projects={[
        {
          id: 'inventory',
          name: 'Inventory',
          secondary: 'Inventory Management',
          url: '/',
        },
        {
          id: 'nms',
          name: 'NMS',
          secondary: 'Network Management',
          url: '/nms',
        },
      ]}
      mainItems={[
        <AssignmentIcon color="primary" style={{margin: '8px'}} />,
        <ListIcon color="primary" style={{margin: '8px'}} />,
      ]}
      secondaryItems={[
        <SettingsIcon color="primary" style={{margin: '8px'}} />,
      ]}
    />
  ))
  .add('Not super user', () => (
    <AppSideBar
      user={{email: 'user@fb.com', isSuperUser: false}}
      projects={[
        {
          id: 'inventory',
          name: 'Inventory',
          secondary: 'Inventory Management',
          url: '/',
        },
        {
          id: 'nms',
          name: 'NMS',
          secondary: 'Network Management',
          url: '/nms',
        },
      ]}
      mainItems={[
        <AssignmentIcon color="primary" style={{margin: '8px'}} />,
        <ListIcon color="primary" style={{margin: '8px'}} />,
      ]}
      secondaryItems={[]}
    />
  ))
  .add('No projects', () => (
    <AppSideBar
      user={{email: 'user@fb.com', isSuperUser: false}}
      projects={[]}
      mainItems={[
        <AssignmentIcon color="primary" style={{margin: '8px'}} />,
        <ListIcon color="primary" style={{margin: '8px'}} />,
      ]}
      secondaryItems={[]}
    />
  ));
