/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import SectionedCard from '../../components/SectionedCard';
import Text from '@fbcnms/ui/components/design-system/Text';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/SectionedCard`, module).add(
  'string',
  () => (
    <div>
      <SectionedCard>
        <Text>Card 1</Text>
      </SectionedCard>
      <SectionedCard>
        <Text>Card 2</Text>
      </SectionedCard>
      <SectionedCard>
        <Text>Card 3</Text>
      </SectionedCard>
    </div>
  ),
);
