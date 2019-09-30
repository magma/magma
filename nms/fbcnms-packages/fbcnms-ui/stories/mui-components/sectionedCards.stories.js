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
import SectionedCard from '../../components/SectionedCard.react';
import Typography from '@material-ui/core/Typography';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/SectionedCard`, module).add(
  'string',
  () => (
    <div>
      <SectionedCard>
        <Typography>Card 1</Typography>
      </SectionedCard>
      <SectionedCard>
        <Typography>Card 2</Typography>
      </SectionedCard>
      <SectionedCard>
        <Typography>Card 3</Typography>
      </SectionedCard>
    </div>
  ),
);
