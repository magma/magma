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
import Breadcrumbs from '../components/Breadcrumbs.react';
import React from 'react';

storiesOf('Breadcrumbs', module).add('collapsed', () => (
  <Breadcrumbs
    breadcrumbs={[
      {
        id: '1',
        name: 'Folder #1',
      },
      {
        id: '2',
        name: 'Folder #2',
      },
      {
        id: '3',
        name: 'Folder #3',
      },
      {
        id: '4',
        name: 'Folder #4',
      },
      {
        id: '5',
        name: 'Folder #5',
      },
      {
        id: '6',
        name: 'Folder #6',
      },
    ]}
    onBreadcrumbClicked={() => {}}
  />
));

storiesOf('Breadcrumbs', module).add('expanded', () => (
  <Breadcrumbs
    breadcrumbs={[
      {
        id: '1',
        name: 'Folder #1',
      },
      {
        id: '2',
        name: 'Folder #2',
      },
      {
        id: '3',
        name: 'Folder #3',
      },
    ]}
    onBreadcrumbClicked={() => {}}
  />
));

storiesOf('Breadcrumbs', module).add('subtext', () => (
  <Breadcrumbs
    breadcrumbs={[
      {
        id: '1',
        name: 'Folder #1',
        subtext: 'Photos',
      },
      {
        id: '2',
        name: 'Folder #2',
        subtext: 'Mexico',
      },
      {
        id: '3',
        name: 'Folder #3',
        subtext: 'Beach',
      },
    ]}
    onBreadcrumbClicked={() => {}}
  />
));

storiesOf('Breadcrumbs', module).add('small', () => (
  <Breadcrumbs
    breadcrumbs={[
      {
        id: '1',
        name: 'Folder #1',
      },
      {
        id: '2',
        name: 'Folder #2',
      },
      {
        id: '3',
        name: 'Folder #3',
      },
      {
        id: '4',
        name: 'Folder #4',
      },
      {
        id: '5',
        name: 'Folder #5',
      },
    ]}
    onBreadcrumbClicked={() => {}}
    size="small"
  />
));
