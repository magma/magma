/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 *
 * @flow
 * @format
 */

import type {EditUserMutationResponse} from '../../../../mutations/__generated__/EditUserMutation.graphql';
import type {MutationCallbacks} from '../../../../mutations/MutationCallbacks';
import type {OptionalRefTypeWrapper} from '../../../../common/EntUtils';
import type {SelectorStoreUpdater} from 'relay-runtime';
import type {UserRole, UsersQuery} from './__generated__/UsersQuery.graphql';
import type {UsersByAuthIDQuery} from './__generated__/UsersByAuthIDQuery.graphql';
import type {UserManagementUtils_user as user} from '../utils/__generated__/UserManagementUtils_user.graphql';
import type {UserManagementUtils_user_base as user_base} from '../utils/__generated__/UserManagementUtils_user_base.graphql';

import EditUserMutation from '../../../../mutations/EditUserMutation';
import RelayEnvironment from '../../../../common/RelayEnvironment';
import axios from 'axios';
import nullthrows from 'nullthrows';
import {ConnectionHandler, fetchQuery, graphql} from 'relay-runtime';
import {LogEvents, ServerLogger} from '../../../../common/LoggingUtils';
import {USER_ROLES} from '../utils/UserManagementUtils';
import {UserRoles} from '@fbcnms/auth/types';
import {getGraphError} from '../../../../common/EntUtils';
import {useLazyLoadQuery} from 'react-relay/hooks';

export type User = OptionalRefTypeWrapper<user>;
export type UserBase = OptionalRefTypeWrapper<user_base>;

const usersQuery = graphql`
  query UsersQuery {
    users(first: 500) @connection(key: "Users_users") {
      edges {
        node {
          ...UserManagementUtils_user @relay(mask: false)
        }
      }
    }
  }
`;

export function useUsers(): $ReadOnlyArray<User> {
  const data = useLazyLoadQuery<UsersQuery>(usersQuery);
  const usersData = data.users?.edges || [];
  return usersData.map(p => p.node).filter(Boolean);
}

function roleToNodeRole(role: UserRole): number {
  return role === USER_ROLES.USER.key ? UserRoles.USER : UserRoles.SUPERUSER;
}

function createNewUserInPlatformServer(
  newUserValue: User,
  password: string,
): Promise<void> {
  const newUserPayload = {
    email: newUserValue.authID,
    password,
    role: roleToNodeRole(newUserValue.role),
    networkIDs: [],
  };
  return axios
    .post<empty, empty>('/user/async/', newUserPayload)
    .then(() => undefined);
}

const userQuery = graphql`
  query UsersByAuthIDQuery($authID: String!) {
    user(authID: $authID) {
      id
    }
  }
`;

function getUserEntIdByAuthID(authID: string): Promise<string> {
  return fetchQuery<UsersByAuthIDQuery>(RelayEnvironment, userQuery, {
    authID: authID.toLowerCase(),
  }).then(response => nullthrows(response.user?.id));
}

function setNewUserEntValues(
  userEntId: string,
  userValues: User,
): Promise<User> {
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
    const users = ConnectionHandler.getConnection(rootQuery, 'Users_users');
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
}

export function addUser(newUserValue: User, password: string): Promise<User> {
  return createNewUserInPlatformServer(newUserValue, password)
    .then(() => getUserEntIdByAuthID(newUserValue.authID))
    .then(userId => setNewUserEntValues(userId, newUserValue));
}

function updateUserInNode(
  email: string,
  role?: UserRole,
  password?: string,
): Promise<void> {
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
  return axios
    .put(`/user/set/${email}`, updateUserPayload)
    .then(() => undefined);
}

export function changeUserPassword(
  user: User,
  password: string,
): Promise<User> {
  return updateUserInNode(user.authID, undefined, password).then(() => user);
}

export function changeCurrentUserPassword(
  currentPassword: string,
  newPassword: string,
): Promise<void> {
  const payload = {
    currentPassword,
    newPassword,
  };
  return axios.post(`/user/change_password`, payload).then(() => undefined);
}

export function editUser(
  newUserValue: User,
  updater?: SelectorStoreUpdater,
): Promise<User> {
  return new Promise<User>((resolve, reject) => {
    const callbacks: MutationCallbacks<EditUserMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(response.editUser);
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
}
