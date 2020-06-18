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

import Button from '../../components/Button';
import FormField from '../../components/FormField';
import {STORY_CATEGORIES} from '../storybookUtils';
import {action} from '@storybook/addon-actions';
import {storiesOf} from '@storybook/react';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Button`, module)
  .add('with text', () => (
    <Button onClick={action('clicked')}>Hello Button</Button>
  ))
  .add('with error', () => <Button error={true}>With Error</Button>)
  .add('with some emoji', () => (
    <Button onClick={action('clicked')}>
      <span role="img" aria-label="so cool">
        😀 😎 👍 💯
      </span>
    </Button>
  ));
storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/FormField`, module)
  .add('default', () => <FormField label="Hello Form" />)
  .add('with children', () => (
    <FormField label="Hello Form field with Button">
      <Button onClick={action('clicked')}>
        <span role="img" aria-label="so cool">
          😀 😎 👍 💯
        </span>
      </Button>
    </FormField>
  ));
