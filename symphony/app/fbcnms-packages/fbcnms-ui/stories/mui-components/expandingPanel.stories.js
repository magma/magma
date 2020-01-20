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
import ExpandingPanel from '../../components/ExpandingPanel';
import React from 'react';
import Text from '../../components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/ExpandingPanel`, module)
  .add('default', () => (
    <ExpandingPanel title="Expanding Panel">
      <Text>This is the content</Text>
    </ExpandingPanel>
  ))
  .add('right button', () => (
    <ExpandingPanel title="Expanding Panel" rightContent={<AddIcon />}>
      <Text>This is the content</Text>
    </ExpandingPanel>
  ));
