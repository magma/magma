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
import LoadingIndicator from '../../common/LoadingIndicator';
import UserAccountPane from '../admin/userManagement/users/UserAccountPane';
import fbt from 'fbt';
import {FormContextProvider} from '../../common/FormContext';
import {Suspense} from 'react';
import {UserManagementContextProvider} from '../admin/userManagement/UserManagementContext';
import {graphql, useLazyLoadQuery} from 'react-relay/hooks';
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

function UserAccountWrapper() {
  const mainContext = useMainContext();

  const loggedInUserID = mainContext.me?.user.id;

  const userData = useLazyLoadQuery<UserManagementContext_UserQuery>(
    userQuery,
    {id: loggedInUserID},
  );

  const loggedInUser: User = userData?.node;

  if (loggedInUserID == null || loggedInUser == null) {
    return <fbt desc="">Failed to identify logged in user account</fbt>;
  }

  return (
    <FormContextProvider ignorePermissions={true}>
      <UserAccountPane user={loggedInUser} isForCurrentUserSettings={true} />
    </FormContextProvider>
  );
}

export default function AccountSettings() {
  return (
    <Suspense fallback={<LoadingIndicator />}>
      <UserManagementContextProvider>
        <UserAccountWrapper />
      </UserManagementContextProvider>
    </Suspense>
  );
}
