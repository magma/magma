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

import type {AddUsersGroupMutationResponse} from '../../../mutations/__generated__/AddUsersGroupMutation.graphql';
import type {DeleteUsersGroupMutationResponse} from '../../../mutations/__generated__/DeleteUsersGroupMutation.graphql';
import type {EditUserMutationResponse} from '../../../mutations/__generated__/EditUserMutation.graphql';
import type {EditUsersGroupMutationResponse} from '../../../mutations/__generated__/EditUsersGroupMutation.graphql';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';
import type {
  PermissionsPolicy,
  User,
  UserPermissionsGroup,
} from './utils/UserManagementUtils';
import type {SelectorStoreUpdater} from 'relay-runtime';
import type {UpdateUsersGroupMembersMutationResponse} from '../../../mutations/__generated__/UpdateUsersGroupMembersMutation.graphql';
import type {
  UserManagementContextQuery,
  UserRole,
} from './__generated__/UserManagementContextQuery.graphql';
import type {UserManagementContext_UserQuery} from './__generated__/UserManagementContext_UserQuery.graphql';
import type {UsersMap} from './utils/UserManagementUtils';

import * as React from 'react';
import AddUsersGroupMutation from '../../../mutations/AddUsersGroupMutation';
import DeleteUsersGroupMutation from '../../../mutations/DeleteUsersGroupMutation';
import EditUserMutation from '../../../mutations/EditUserMutation';
import EditUsersGroupMutation from '../../../mutations/EditUsersGroupMutation';
import LoadingIndicator from '../../../common/LoadingIndicator';
import RelayEnvironment from '../../../common/RelayEnvironment';
import UpdateUsersGroupMembersMutation from '../../../mutations/UpdateUsersGroupMembersMutation';
import axios from 'axios';
import nullthrows from 'nullthrows';
import {ConnectionHandler, fetchQuery, graphql} from 'relay-runtime';
import {LogEvents, ServerLogger} from '../../../common/LoggingUtils';
import {RelayEnvironmentProvider} from 'react-relay/hooks';
import {Suspense} from 'react';
import {USER_ROLES} from './utils/UserManagementUtils';
import {getGraphError} from '../../../common/EntUtils';
import {
  groupResponse2Group,
  groupsResponse2Groups,
  permissionsPoliciesResponse2PermissionsPolicies,
  userResponse2User,
  users2UsersMap,
  usersResponse2Users,
} from './utils/UserManagementUtils';
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

const deleteGroup = (id: string) => {
  return new Promise((resolve, reject) => {
    const callbacks: MutationCallbacks<DeleteUsersGroupMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve();
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };
    const removeGroupFromStore = store => {
      const rootQuery = store.getRoot();
      const groups = ConnectionHandler.getConnection(
        rootQuery,
        'UserManagementContext_usersGroups',
      );
      if (groups == null) {
        return;
      }
      ConnectionHandler.deleteNode(groups, id);
      store.delete(id);
    };
    DeleteUsersGroupMutation({id}, callbacks, removeGroupFromStore);
  });
};

const editGroup = (usersMap: UsersMap) => (
  newGroupValue: UserPermissionsGroup,
) => {
  return new Promise<UserPermissionsGroup>((resolve, reject) => {
    const callbacks: MutationCallbacks<EditUsersGroupMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(groupResponse2Group(response.editUsersGroup, usersMap));
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };
    EditUsersGroupMutation(
      {
        input: {
          id: newGroupValue.id,
          name: newGroupValue.name,
          description: newGroupValue.description,
          status: newGroupValue.status,
        },
      },
      callbacks,
    );
  });
};

const updateGroupMembers = (usersMap: UsersMap) => (
  group: UserPermissionsGroup,
  addUserIds: Array<string>,
  removeUserIds: Array<string>,
) => {
  return new Promise<UserPermissionsGroup>((resolve, reject) => {
    const cbs: MutationCallbacks<UpdateUsersGroupMembersMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(
          groupResponse2Group(response.updateUsersGroupMembers, usersMap),
        );
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };
    UpdateUsersGroupMembersMutation(
      {
        input: {
          id: group.id,
          addUserIds,
          removeUserIds,
        },
      },
      cbs,
    );
  });
};

const addGroup = (usersMap: UsersMap) => (
  newGroupValue: UserPermissionsGroup,
) => {
  return new Promise<UserPermissionsGroup>((resolve, reject) => {
    const callbacks: MutationCallbacks<AddUsersGroupMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(groupResponse2Group(response.addUsersGroup, usersMap));
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };

    const addNewGroupToStore = store => {
      const rootQuery = store.getRoot();
      // eslint-disable-next-line no-warning-comments
      // $FlowFixMe (T62907961) Relay flow types
      const newNode = store.getRootField('addUsersGroup');
      if (newNode == null) {
        return;
      }
      const groups = ConnectionHandler.getConnection(
        rootQuery,
        'UserManagementContext_usersGroups',
      );
      if (groups == null) {
        return;
      }
      const edge = ConnectionHandler.createEdge(
        store,
        groups,
        newNode,
        'UsersGroupEdge',
      );
      ConnectionHandler.insertEdgeAfter(groups, edge);
    };
    AddUsersGroupMutation(
      {
        input: {
          name: newGroupValue.name,
          description: newGroupValue.description,
        },
      },
      callbacks,
      addNewGroupToStore,
    );
  });
};

type UserManagementContextValue = {
  policies: Array<PermissionsPolicy>,
  groups: Array<UserPermissionsGroup>,
  users: Array<User>,
  usersMap: UsersMap,
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
  addGroup: UserPermissionsGroup => Promise<UserPermissionsGroup>,
  editGroup: UserPermissionsGroup => Promise<UserPermissionsGroup>,
  deleteGroup: (id: string) => Promise<void>,
  updateGroupMembers: (
    group: UserPermissionsGroup,
    addUserIds: Array<string>,
    removeUserIds: Array<string>,
  ) => Promise<UserPermissionsGroup>,
};

const emptyUsersMap = new Map<string, User>();
const UserManagementContext = React.createContext<UserManagementContextValue>({
  policies: [],
  groups: [],
  users: [],
  usersMap: emptyUsersMap,
  addUser,
  editUser,
  changeUserPassword,
  changeCurrentUserPassword,
  addGroup: addGroup(emptyUsersMap),
  editGroup: editGroup(emptyUsersMap),
  deleteGroup,
  updateGroupMembers: updateGroupMembers(emptyUsersMap),
});

type Props = {
  children: React.Node,
};

const dataQuery = graphql`
  query UserManagementContextQuery {
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
    usersGroups(first: 500)
      @connection(key: "UserManagementContext_usersGroups") {
      edges {
        node {
          id
          name
          description
          status
          members {
            id
            authID
          }
        }
      }
    }
    permissionsPolicies(first: 50) {
      edges {
        node {
          id
          name
          description
          isGlobal
          policy {
            ... on InventoryPolicy {
              __typename
              read {
                isAllowed
              }
            }
            ... on WorkforcePolicy {
              __typename
              read {
                isAllowed
              }
            }
          }
          groups {
            id
          }
        }
      }
    }
  }
`;

function ProviderWrap(props: Props) {
  const providerValue = (users, groups, policies, usersMap) => ({
    policies,
    groups,
    users,
    usersMap,
    addUser,
    editUser,
    changeUserPassword,
    changeCurrentUserPassword,
    addGroup: addGroup(usersMap),
    editGroup: editGroup(usersMap),
    deleteGroup,
    updateGroupMembers: updateGroupMembers(usersMap),
  });

  const data = useLazyLoadQuery<UserManagementContextQuery>(dataQuery);

  const users = usersResponse2Users(data.users);
  const usersMap = users2UsersMap(users);
  const groups = groupsResponse2Groups(data.usersGroups, usersMap);
  const policies = permissionsPoliciesResponse2PermissionsPolicies(
    data.permissionsPolicies,
  );

  return (
    <UserManagementContext.Provider
      value={providerValue(users, groups, policies, usersMap)}>
      {props.children}
    </UserManagementContext.Provider>
  );
}

export function UserManagementContextProvider(props: Props) {
  return (
    <RelayEnvironmentProvider environment={RelayEnvironment}>
      <Suspense fallback={<LoadingIndicator />}>
        <ProviderWrap {...props} />
      </Suspense>
    </RelayEnvironmentProvider>
  );
}

export function useUserManagement() {
  return useContext(UserManagementContext);
}

export default UserManagementContext;
