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

import {storiesOf} from '@storybook/react';

import AppContent from '../../components/layout/AppContent';
import AppDrawer from '../../components/layout/AppDrawer';
import CssBaseline from '@material-ui/core/CssBaseline';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import PublicIcon from '@material-ui/icons/Public';
import TopBar from '../../components/layout/TopBar';
import {STORY_CATEGORIES} from '../storybookUtils';
import {TopBarContextProvider} from '../../components/layout/TopBarContext';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/layout.TopBar`, module).add(
  'default',
  () => (
    <TopBarContextProvider>
      <CssBaseline />
      <div style={{display: 'flex'}}>
        <AppDrawer>
          <NavListItem label="Item 1" path="/item1" icon={<PublicIcon />} />
          <NavListItem label="Item 2" path="/item2" icon={<PublicIcon />} />
        </AppDrawer>
        <AppContent>
          <TopBar title="Title">Right hand content</TopBar>
          <div>Content</div>
        </AppContent>
      </div>
    </TopBarContextProvider>
  ),
);
