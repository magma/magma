/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import React from 'react';
import TreeView from '../../components/TreeView';
import {STORY_CATEGORIES} from '../storybookUtils';
import {action} from '@storybook/addon-actions';
import {storiesOf} from '@storybook/react';

const levels = [1, 2, 3, 4];

const treeLeaves = levels.map(lvl => {
  return {
    name: 'Leaf ' + lvl,
    subtitle: 'Go ' + lvl,
    children: [],
  };
});

const treeLvls = levels.map(lvl => {
  return {
    name: 'Level ' + lvl,
    subtitle: 'Go ' + lvl,
    children: treeLeaves,
  };
});

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/TreeView`, module)
  .add('default', () => (
    <TreeView onClick={() => {}} tree={treeLvls} selectedId={null} />
  ))
  .add('with actions', () => (
    <TreeView onClick={action('clicked')} tree={treeLvls} selectedId={null} />
  ))
  .add('with subtitle', () => (
    <TreeView
      onClick={() => {}}
      subtitlePropertyGetter={(node: Object) => node.subtitle}
      tree={treeLvls}
      selectedId={null}
    />
  ));
