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
import Tokenizer from '../components/Tokenizer.react';
import {storiesOf} from '@storybook/react';

const entries = [
  {label: 'Chassis', id: '0'},
  {label: 'Rack', id: '1'},
  {label: 'Card', id: '2'},
  {label: 'AP', id: '3'},
];

storiesOf('Tokenizer', module).add('basic', () => (
  <div style={{width: '300px'}}>
    <Tokenizer searchEntries={entries} onEntriesRequested={() => {}} />
  </div>
));
