/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import PageFooter from '../../components/PageFooter.react';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/PageFooter`, module).add(
  'string',
  () => (
    <PageFooter>
      <Typography>Wow</Typography>
    </PageFooter>
  ),
);
