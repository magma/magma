/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import ConfigureTitle from '../../components/ConfigureTitle';
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/ConfigureTitle`, module).add(
  'default',
  () => (
    <ConfigureTitle
      title={'This is a title'}
      subtitle={'This is a descriptive subtitle'}
    />
  ),
);
