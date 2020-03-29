/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {EditUserMutationResponse} from '../../../mutations/__generated__/EditUserMutation.graphql';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';
import type {StoreUpdater} from '../../../common/RelayEnvironment';
import type {User} from './TempTypes';
import type {UserManagementContext_UserQuery} from './__generated__/UserManagementContext_UserQuery.graphql';
import type {UserManagementContext_UsersQueryResponse} from './__generated__/UserManagementContext_UsersQuery.graphql';
import type {UserRole} from './__generated__/UserManagementContext_UsersQuery.graphql';

import * as React from 'react';
import EditUserMutation from '../../../mutations/EditUserMutation';
import InventoryQueryRenderer from '../../InventoryQueryRenderer';
import RelayEnvironment from '../../../common/RelayEnvironment';
import axios from 'axios';
import nullthrows from 'nullthrows';
import {ConnectionHandler, fetchQuery, graphql} from 'relay-runtime';
import {LogEvents, ServerLogger} from '../../../common/LoggingUtils';
import {USER_ROLES} from './TempTypes';
import {useContext} from 'react';

export type UserManagementContextValue = {
  users: Array<User>,
  addUser: (user: User, password: string) => Promise<User>,
  editUser: (newUserValue: User, updater?: StoreUpdater) => Promise<User>,
  changeUserPassword: (user: User, password: string) => Promise<User>,
};

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
  return axios.post('/user/async/', newUserPayload);
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
  userValues.id = userEntId;
  const addNewUserToStore = store => {
    const rootQuery = store.getRoot();
    const newNode = store.get(userValues.id);
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
  return editUser(userValues, addNewUserToStore);
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

const editUser = (newUserValue: User, updater?: StoreUpdater) => {
  const graphCall = new Promise<User>((resolve, reject) => {
    const callbacks: MutationCallbacks<EditUserMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(errors[0].message);
        }
        resolve({
          id: response.editUser.id,
          authID: response.editUser.authID,
          firstName: response.editUser.firstName,
          lastName: response.editUser.lastName,
          role: response.editUser.role,
          status: response.editUser.status,
        });
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
        reject(e.message);
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

  return graphCall.then();
};
const UserManagementContext = React.createContext<UserManagementContextValue>({
  users: [],
  addUser,
  editUser,
  changeUserPassword,
});

export function useUserManagement() {
  return useContext(UserManagementContext);
}

type Props = {
  children: React.Node,
};

const usersQuery = graphql`
  query UserManagementContext_UsersQuery {
    users(first: 500) @connection(key: "UserManagementContext_users") {
      edges {
        node {
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
  }
`;

export function UserManagementContextProvider(props: Props) {
  const providerValue = users => ({
    users,
    addUser,
    editUser,
    changeUserPassword,
  });
  return (
    <InventoryQueryRenderer
      query={usersQuery}
      variables={{}}
      render={(response: UserManagementContext_UsersQueryResponse) => {
        const users: Array<User> = [];
        if (response.users != null) {
          // using 'for' and not simple 'map' beacuse of flow.
          for (let i = 0; i < response.users.edges.length; i++) {
            const userNode = response.users.edges[i].node;
            if (userNode == null) {
              continue;
            }
            users.push({
              id: userNode.id,
              authID: userNode.authID,
              firstName: userNode.firstName,
              lastName: userNode.lastName,
              role: userNode.role,
              status: userNode.status,
            });
          }
        }
        return (
          <UserManagementContext.Provider value={providerValue(users)}>
            {props.children}
          </UserManagementContext.Provider>
        );
      }}
    />
  );
}

export default UserManagementContext;
