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

import * as React from 'react';
import EditUserMutation from '../../../mutations/EditUserMutation';
import InventoryQueryRenderer from '../../InventoryQueryRenderer';
import RelayEnvironment from '../../../common/RelayEnvironment';
import axios from 'axios';
import nullthrows from 'nullthrows';
import {ConnectionHandler, fetchQuery, graphql} from 'relay-runtime';
import {USER_ROLES} from './TempTypes';
import {useContext} from 'react';

export type UserManagementContextValue = {
  users: Array<User>,
  addUser: (user: User, password: string) => Promise<User>,
  editUser: (newUserValue: User, updater?: StoreUpdater) => Promise<User>,
};

const editUser = (newUserValue: User, updater?: StoreUpdater) => {
  return new Promise<User>((resolve, reject) => {
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
};

const userQuery = graphql`
  query UserManagementContext_UserQuery($authID: String!) {
    user(authID: $authID) {
      id
    }
  }
`;

const createNewUserInNode = (newUserValue: User, password: string) => {
  const newUserPayload = {
    email: newUserValue.authID,
    password,
    role: newUserValue.role === USER_ROLES.USER.key ? 0 : 3,
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
const UserManagementContext = React.createContext<UserManagementContextValue>({
  users: [],
  addUser,
  editUser,
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
