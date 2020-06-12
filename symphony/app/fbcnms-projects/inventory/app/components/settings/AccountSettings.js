/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {User} from '../admin/userManagement/utils/UserManagementUtils';
import type {UserManagementContext_UserQuery} from '../admin/userManagement/__generated__/UserManagementContext_UserQuery.graphql';

import * as React from 'react';
import InventorySuspense from '../../common/InventorySuspense';
import UserAccountPane from '../admin/userManagement/users/UserAccountPane';
import ViewContainer from '@fbcnms/ui/components/design-system/View/ViewContainer';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {UserManagementContextProvider} from '../admin/userManagement/UserManagementContext';
import {graphql, useLazyLoadQuery} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';
import {useMainContext} from '../MainContext';

const userQuery = graphql`
  query AccountSettings_UserQuery($id: ID!) {
    node(id: $id) {
      ... on User {
        id
        authID
        firstName
        lastName
        email
        status
        role
        groups {
          id
          name
        }
        profilePhoto {
          id
          fileName
          storeKey
        }
      }
    }
  }
`;

const useStyles = makeStyles(() => ({
  settingsPage: {
    borderRadius: '4px',
    boxShadow: symphony.shadows.DP1,
    background: symphony.palette.white,
    flexGrow: 1,
    padding: '32px',
    maxWidth: '1024px',
  },
}));

function UserAccountWrapper() {
  const classes = useStyles();
  const mainContext = useMainContext();

  const loggedInUserID = mainContext.me?.user?.id;

  const userData = useLazyLoadQuery<UserManagementContext_UserQuery>(
    userQuery,
    {id: loggedInUserID},
  );

  const loggedInUser: User = userData?.node;

  if (loggedInUserID == null || loggedInUser == null) {
    return <fbt desc="">Failed to identify logged in user account</fbt>;
  }

  return (
    <ViewContainer
      header={{
        title: <fbt desc="">User Settings</fbt>,
        subtitle: <fbt desc="">Manage your own private settings.</fbt>,
      }}>
      <div className={classes.settingsPage}>
        <UserAccountPane user={loggedInUser} isForCurrentUserSettings={true} />
      </div>
    </ViewContainer>
  );
}

export default function AccountSettings() {
  return (
    <InventorySuspense>
      <UserManagementContextProvider>
        <UserAccountWrapper />
      </UserManagementContextProvider>
    </InventorySuspense>
  );
}
