/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import CardFooter from '../../components/CardFooter';
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/CardFooter`, module)
  .add('left', () => (
    <CardFooter alignItems="left">
      <div>Option 1</div>
      <div>Option 2</div>
    </CardFooter>
  ))
  .add('right', () => (
    <CardFooter alignItems="right">
      <div>Option 1</div>
      <div>Option 2</div>
    </CardFooter>
  ));
