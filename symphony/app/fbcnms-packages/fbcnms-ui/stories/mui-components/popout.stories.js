/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '@material-ui/core/Button';
import Popout from '../../components/Popout';
import React from 'react';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Popout`, module).add(
  'default',
  () => (
    <div style={{padding: 100}}>
      <Popout
        content={
          <div style={{padding: 20}}>
            <Text variant="body2">Content</Text>
          </div>
        }>
        <Button variant="contained" color="primary">
          Click me!
        </Button>
      </Popout>
    </div>
  ),
);
