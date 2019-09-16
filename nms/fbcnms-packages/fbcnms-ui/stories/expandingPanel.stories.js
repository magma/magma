/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import ExpandingPanel from '../components/ExpandingPanel.react';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import {storiesOf} from '@storybook/react';

storiesOf('ExpandingPanel', module).add('default', () => (
  <ExpandingPanel title="Expanding Panel">
    <Typography>This is the content</Typography>
  </ExpandingPanel>
));
