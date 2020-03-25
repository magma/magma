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

import {action} from '@storybook/addon-actions';
import {storiesOf} from '@storybook/react';

import EditUserDialog from '../../components/auth/EditUserDialog';
import LoginForm from '../../components/auth/LoginForm';
import {STORY_CATEGORIES} from '../storybookUtils';

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/auth.EditUserDialog`, module)
  .add('default', () => (
    <EditUserDialog
      editingUser={null}
      open={true}
      onClose={action('close')}
      onEditUser={action('edit user')}
      onCreateUser={action('create user')}
      allNetworkIDs={['network1', 'network2']}
      ssoEnabled={false}
    />
  ))
  .add('no networks', () => (
    <EditUserDialog
      editingUser={null}
      open={true}
      onClose={action('close')}
      onEditUser={action('edit user')}
      onCreateUser={action('create user')}
      ssoEnabled={false}
    />
  ))
  .add('SSO user (i.e. no password)', () => (
    <EditUserDialog
      editingUser={null}
      open={true}
      onClose={action('close')}
      onEditUser={action('edit user')}
      onCreateUser={action('create user')}
      ssoEnabled={true}
    />
  ));

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/auth.LoginForm`, module).add(
  'default',
  () => (
    <LoginForm
      action="/test/user/login"
      title="My title"
      csrfToken="abcd1234"
    />
  ),
);
