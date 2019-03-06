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
import CardFooter from '../components/CardFooter.react';
import React from 'react';

storiesOf('CardFooter', module)
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
