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
import ExpandingPanel from '../../components/ExpandingPanel.react';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/ExpandingPanel`, module)
  .add('default', () => (
    <ExpandingPanel title="Expanding Panel">
      <Typography>This is the content</Typography>
    </ExpandingPanel>
  ))
  .add('right button', () => (
    <ExpandingPanel title="Expanding Panel" rightContent={<AddIcon />}>
      <Typography>This is the content</Typography>
    </ExpandingPanel>
  ));
