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
import PermissionsGroupsView from './PermissionsGroupsView';
import Strings from '../../../common/CommonStrings';
import UsersView from './UsersView';
import fbt from 'fbt';

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
    navigation: {
      label: USERS_HEADER,
      tooltip: `${USERS_HEADER}`,
    },
    header: {
      title: `${USERS_HEADER}`,
      subtitle:
        'Add and manage your organization users, and set their role to control their global settings',
      actionButtons: [
        {
          title: fbt('Add User', ''),
          action: () => {},
        },
      ],
    },
    children: <UsersView />,
  },
  {
    navigation: {
      label: PERMISSIONS_GROUPS_HEADER,
      tooltip: `${PERMISSIONS_GROUPS_HEADER}`,
    },
    header: {
      title: `${PERMISSIONS_GROUPS_HEADER}`,
      subtitle:
        'Create groups with different rules and add users to apply permissions',
      actionButtons: [
        {
          title: fbt('Create Group', ''),
          action: () => {},
        },
      ],
    },
    children: <PermissionsGroupsView />,
  },
];

export default function UserManaementView() {
  return (
    <NavigatableViews header={Strings.admin.users.viewHeader} views={VIEWS} />
  );
}
