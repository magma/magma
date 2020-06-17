/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 *
 * @flow strict-local
 * @format
 */

// flowlint untyped-import:off

import type {EditUserMutationResponse} from '../../../mutations/__generated__/EditUserMutation.graphql';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';
import type {User} from './utils/UserManagementUtils';
import type {
  UserManagementContextQuery,
  UserRole,
} from './__generated__/UserManagementContextQuery.graphql';
import type {UserManagementContext_UserQuery} from './__generated__/UserManagementContext_UserQuery.graphql';

import * as React from 'react';
import EditUserMutation from '../../../mutations/EditUserMutation';
import InventorySuspense from '../../../common/InventorySuspense';
import RelayEnvironment from '../../../common/RelayEnvironment';
import axios from 'axios';
import nullthrows from 'nullthrows';
import {ConnectionHandler, fetchQuery, graphql} from 'relay-runtime';
import {LogEvents, ServerLogger} from '../../../common/LoggingUtils';
import {RelayEnvironmentProvider} from 'react-relay/hooks';
import {
  USER_ROLES,
  userResponse2User,
  usersResponse2Users,
} from './utils/UserManagementUtils';
import {getGraphError} from '../../../common/EntUtils';
import {useContext} from 'react';
import {useLazyLoadQuery} from 'react-relay/hooks';

const userQuery = graphql`
  query UserManagementContext_UserQuery($authID: String!) {
    user(authID: $authID) {
      id
    }
  }
`;

const roleToNodeRole = (role: UserRole) =>
  role === USER_ROLES.USER.key ? 0 : 3;

const createNewUserInNode = (newUserValue: User, password: string) => {
  const newUserPayload = {
    email: newUserValue.authID,
    password,
    role: roleToNodeRole(newUserValue.role),
    networkIDs: [],
  };
  return axios.post<empty, empty>('/user/async/', newUserPayload);
};

const getUserEntIdByAuthID = authID => {
  return fetchQuery<UserManagementContext_UserQuery>(
    RelayEnvironment,
    userQuery,
    {
      authID,
    },
  ).then(response => nullthrows(response.user?.id));
};

const setNewUserEntValues = (userEntId: string, userValues: User) => {
  const mutatedUser = {
    ...userValues,
    id: userEntId,
  };
  const addNewUserToStore = store => {
    const rootQuery = store.getRoot();
    const newNode = store.get(mutatedUser.id);
    if (newNode == null) {
      return;
    }
    const users = ConnectionHandler.getConnection(
      rootQuery,
      'UserManagementContext_users',
    );
    if (users == null) {
      return;
    }
    const edge = ConnectionHandler.createEdge(
      store,
      users,
      newNode,
      'UserEdge',
    );
    ConnectionHandler.insertEdgeAfter(users, edge);
  };
  return editUser(mutatedUser, addNewUserToStore);
};

const addUser = (newUserValue: User, password: string) => {
  return createNewUserInNode(newUserValue, password)
    .then(() => getUserEntIdByAuthID(newUserValue.authID))
    .then(userId => setNewUserEntValues(userId, newUserValue));
};

const updateUserInNode = (
  email: string,
  role?: UserRole,
  password?: string,
) => {
  const updateUserPayload = {};
  if (password != null) {
    updateUserPayload.password = password;
  }
  if (role != null) {
    updateUserPayload.role = roleToNodeRole(role);
  }
  if (Object.keys(updateUserPayload).length === 0) {
    return Promise.resolve();
  }
  return axios.put(`/user/set/${email}`, updateUserPayload);
};

const changeUserPassword = (user: User, password: string) => {
  return updateUserInNode(user.authID, undefined, password).then(() => user);
};

const changeCurrentUserPassword = (
  currentPassword: string,
  newPassword: string,
) => {
  const payload = {
    currentPassword,
    newPassword,
  };
  return axios.post(`/user/change_password`, payload).then(() => undefined);
};

const editUser = (newUserValue: User, updater?: SelectorStoreUpdater) => {
  return new Promise<User>((resolve, reject) => {
    const callbacks: MutationCallbacks<EditUserMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(userResponse2User(response.editUser));
        // TEMP: Need to update Node with the new role.
        // (Once Node is changed to take the role from graph,
        //  we can remove this)
        updateUserInNode(newUserValue.authID, newUserValue.role).catch(
          error => {
            ServerLogger.error(LogEvents.CLIENT_FATAL_ERROR, {
              message: error.message,
              stack: error.stack,
            });
          },
        );
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };
    EditUserMutation(
      {
        input: {
          id: newUserValue.id,
          firstName: newUserValue.firstName,
          lastName: newUserValue.lastName,
          role: newUserValue.role,
          status: newUserValue.status,
        },
      },
      callbacks,
      updater,
    );
  });
};

type UserManagementContextValue = {
  users: Array<User>,
  addUser: (user: User, password: string) => Promise<User>,
  editUser: (
    newUserValue: User,
    updater?: SelectorStoreUpdater,
  ) => Promise<User>,
  changeUserPassword: (user: User, password: string) => Promise<User>,
  changeCurrentUserPassword: (
    currentPassword: string,
    newPassword: string,
  ) => Promise<void>,
};

const UserManagementContext = React.createContext<UserManagementContextValue>({
  users: [],
  addUser,
  editUser,
  changeUserPassword,
  changeCurrentUserPassword,
});

type Props = {
  children: React.Node,
};

const dataQuery = graphql`
  query UserManagementContextQuery {
    users(first: 500) @connection(key: "UserManagementContext_users") {
      edges {
        node {
          ...UserManagementUtils_user @relay(mask: false)
        }
      }
    }
  }
`;

function ProviderWrap(props: Props) {
  const providerValue = users => ({
    users,
    addUser,
    editUser,
    changeUserPassword,
    changeCurrentUserPassword,
  });

  const data = useLazyLoadQuery<UserManagementContextQuery>(dataQuery);

  const users = usersResponse2Users(data.users);

  return (
    <UserManagementContext.Provider value={providerValue(users)}>
      {props.children}
    </UserManagementContext.Provider>
  );
}

export function UserManagementContextProvider(props: Props) {
  return (
    <RelayEnvironmentProvider environment={RelayEnvironment}>
      <InventorySuspense>
        <ProviderWrap {...props} />
      </InventorySuspense>
    </RelayEnvironmentProvider>
  );
}

export function useUserManagement() {
  return useContext(UserManagementContext);
}

export default UserManagementContext;
