/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {NavigatableView} from '@fbcnms/ui/components/design-system/View/NavigatableViews';

import * as React from 'react';
import NavigatableViews from '@fbcnms/ui/components/design-system/View/NavigatableViews';
import NewUserDialog from './NewUserDialog';
import PermissionsGroupsView from './PermissionsGroupsView';
import Strings from '../../../common/CommonStrings';
import UsersView from './UsersView';
import emptyFunction from '@fbcnms/util/emptyFunction';
import fbt from 'fbt';
import {useState} from 'react';

const USERS_HEADER = fbt(
  'Users & Roles',
  'Header for view showing system users settings',
);
const PERMISSIONS_GROUPS_HEADER = fbt(
  'Permission Groups',
  'Header for view showing system permissions settings',
);

const VIEWS: Array<NavigatableView> = [
  {
    menuItem: {
      label: USERS_HEADER,
      tooltip: `${USERS_HEADER}`,
    },
    component: {
      header: {
        title: `${USERS_HEADER}`,
        subtitle:
          'Add and manage your organization users, and set their role to control their global settings',
        actionButtons: [
          {
            title: fbt('Add User', ''),
            action: emptyFunction,
          },
        ],
      },
      children: <UsersView />,
    },
  },
  {
    menuItem: {
      label: PERMISSIONS_GROUPS_HEADER,
      tooltip: `${PERMISSIONS_GROUPS_HEADER}`,
    },
    component: {
      header: {
        title: `${PERMISSIONS_GROUPS_HEADER}`,
        subtitle:
          'Create groups with different rules and add users to apply permissions',
        actionButtons: [
          {
            title: fbt('Create Group', ''),
            action: emptyFunction,
          },
        ],
      },
      children: <PermissionsGroupsView />,
    },
  },
];

export default function UserManaementView() {
  const [addingNewUser, setAddingNewUser] = useState(false);

  if (VIEWS[0].component != null) {
    const userActions = VIEWS[0].component.header?.actionButtons;
    if (userActions != null && userActions.length > 0) {
      userActions[0].action = () => setAddingNewUser(true);
    }
  }

  return (
    <>
      <NavigatableViews header={Strings.admin.users.viewHeader} views={VIEWS} />
      {addingNewUser && (
        <NewUserDialog
          isOpened={addingNewUser}
          onClose={() => setAddingNewUser(false)}
        />
      )}
    </>
  );
}
