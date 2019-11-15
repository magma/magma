/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Breadcrumbs from '../../components/Breadcrumbs';
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';

const onBreadcrumbClicked = b => window.alert(`Clicked ${b.name}`);

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Breadcrumbs`, module)
  .add('collapsed', () => (
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
      ].map(b => ({...b, onClick: () => onBreadcrumbClicked(b)}))}
    />
  ))
  .add('expanded', () => (
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
      ].map(b => ({...b, onClick: () => onBreadcrumbClicked(b)}))}
    />
  ))
  .add('subtext', () => (
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
      ].map(b => ({...b, onClick: () => onBreadcrumbClicked(b)}))}
    />
  ))
  .add('small', () => (
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
      ].map(b => ({...b, onClick: () => onBreadcrumbClicked(b)}))}
      size="small"
    />
  ));
