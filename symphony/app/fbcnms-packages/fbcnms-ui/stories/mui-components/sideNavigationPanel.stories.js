/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React, {useState} from 'react';
import SideNavigationPanel from '../../components/SideNavigationPanel';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

const TemplatesPanel = () => {
  const [selectedItem, setSelectedItem] = useState('0');
  return (
    <SideNavigationPanel
      title="Templates"
      items={[{key: '0', label: 'Work Orders'}, {key: '1', label: 'Projects'}]}
      selectedItemId={selectedItem}
      onItemClicked={({key}) => setSelectedItem(key)}
    />
  );
};

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/SideNavigationPanel`, module).add(
  'default',
  () => <TemplatesPanel />,
);
