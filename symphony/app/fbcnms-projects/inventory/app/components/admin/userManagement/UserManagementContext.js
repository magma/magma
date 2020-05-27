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

import type {AddPermissionsPolicyMutationResponse} from '../../../mutations/__generated__/AddPermissionsPolicyMutation.graphql';
import type {AddUsersGroupMutationResponse} from '../../../mutations/__generated__/AddUsersGroupMutation.graphql';
import type {DeletePermissionsPolicyMutationResponse} from '../../../mutations/__generated__/DeletePermissionsPolicyMutation.graphql';
import type {DeleteUsersGroupMutationResponse} from '../../../mutations/__generated__/DeleteUsersGroupMutation.graphql';
import type {EditPermissionsPolicyMutationResponse} from '../../../mutations/__generated__/EditPermissionsPolicyMutation.graphql';
import type {EditUserMutationResponse} from '../../../mutations/__generated__/EditUserMutation.graphql';
import type {EditUsersGroupMutationResponse} from '../../../mutations/__generated__/EditUsersGroupMutation.graphql';
import type {
  GroupsMap,
  PermissionsPolicy,
  User,
  UserPermissionsGroup,
  UsersMap,
} from './utils/UserManagementUtils';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';
import type {
  UserManagementContextQuery,
  UserRole,
} from './__generated__/UserManagementContextQuery.graphql';
import type {UserManagementContext_UserQuery} from './__generated__/UserManagementContext_UserQuery.graphql';

import * as React from 'react';
import AddPermissionsPolicyMutation from '../../../mutations/AddPermissionsPolicyMutation';
import AddUsersGroupMutation from '../../../mutations/AddUsersGroupMutation';
import DeletePermissionsPolicyMutation from '../../../mutations/DeletePermissionsPolicyMutation';
import DeleteUsersGroupMutation from '../../../mutations/DeleteUsersGroupMutation';
import EditPermissionsPolicyMutation from '../../../mutations/EditPermissionsPolicyMutation';
import EditUserMutation from '../../../mutations/EditUserMutation';
import EditUsersGroupMutation from '../../../mutations/EditUsersGroupMutation';
import InventorySuspense from '../../../common/InventorySuspense';
import RelayEnvironment from '../../../common/RelayEnvironment';
import axios from 'axios';
import nullthrows from 'nullthrows';
import {ConnectionHandler, fetchQuery, graphql} from 'relay-runtime';
import {LogEvents, ServerLogger} from '../../../common/LoggingUtils';
import {RelayEnvironmentProvider} from 'react-relay/hooks';
import {
  USER_ROLES,
  groupResponse2Group,
  groups2GroupsMap,
  groupsResponse2Groups,
  permissionsPoliciesResponse2PermissionsPolicies,
  permissionsPolicy2PermissionsPolicyInput,
  permissionsPolicyResponse2PermissionsPolicy,
  userResponse2User,
  users2UsersMap,
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
        resolve(groupResponse2Group(usersMap)(response.editUsersGroup));
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
          members: newGroupValue.members.map(m => m.id),
        },
      },
      callbacks,
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
        resolve(groupResponse2Group(usersMap)(response.addUsersGroup));
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
          members: newGroupValue.members.map(m => m.id),
        },
      },
      callbacks,
      addNewGroupToStore,
    );
  });
};

const addPermissionsPolicy = (groupsMap: GroupsMap) => (
  newPolicyValue: PermissionsPolicy,
) => {
  return new Promise<PermissionsPolicy>((resolve, reject) => {
    const callbacks: MutationCallbacks<AddPermissionsPolicyMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(
          permissionsPolicyResponse2PermissionsPolicy(groupsMap)(
            response.addPermissionsPolicy,
          ),
        );
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };

    const addNewPolicyToStore = store => {
      const rootQuery = store.getRoot();
      const newNode = store.getRootField('addPermissionsPolicy');
      if (newNode == null) {
        return;
      }
      const policies = ConnectionHandler.getConnection(
        rootQuery,
        'UserManagementContext_permissionsPolicies',
      );
      if (policies == null) {
        return;
      }
      const edge = ConnectionHandler.createEdge(
        store,
        policies,
        newNode,
        'PermisionsPolicyEdge',
      );
      ConnectionHandler.insertEdgeAfter(policies, edge);
    };
    AddPermissionsPolicyMutation(
      {
        input: permissionsPolicy2PermissionsPolicyInput(newPolicyValue),
      },
      callbacks,
      addNewPolicyToStore,
    );
  });
};

const editPermissionsPolicy = (groupsMap: GroupsMap) => (
  newPolicyValue: PermissionsPolicy,
) => {
  return new Promise<PermissionsPolicy>((resolve, reject) => {
    type Callbacks = MutationCallbacks<EditPermissionsPolicyMutationResponse>;
    const callbacks: Callbacks = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(
          permissionsPolicyResponse2PermissionsPolicy(groupsMap)(
            response.editPermissionsPolicy,
          ),
        );
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };

    EditPermissionsPolicyMutation(
      {
        input: {
          id: newPolicyValue.id,
          ...permissionsPolicy2PermissionsPolicyInput(newPolicyValue),
        },
      },
      callbacks,
    );
  });
};

const deletePermissionsPolicy = (id: string) => {
  return new Promise((resolve, reject) => {
    const cbs: MutationCallbacks<DeletePermissionsPolicyMutationResponse> = {
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
    const removePolicyFromStore = store => {
      const rootQuery = store.getRoot();
      const policies = ConnectionHandler.getConnection(
        rootQuery,
        'UserManagementContext_permissionsPolicies',
      );
      if (policies == null) {
        return;
      }
      ConnectionHandler.deleteNode(policies, id);
      store.delete(id);
    };
    DeletePermissionsPolicyMutation({id}, cbs, removePolicyFromStore);
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
  addPermissionsPolicy: PermissionsPolicy => Promise<PermissionsPolicy>,
  editPermissionsPolicy: PermissionsPolicy => Promise<PermissionsPolicy>,
  deletePermissionsPolicy: (id: string) => Promise<void>,
};

const emptyUsersMap = new Map<string, User>();
const emptyGroupsMap = new Map<string, UserPermissionsGroup>();
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
  addPermissionsPolicy: addPermissionsPolicy(emptyGroupsMap),
  editPermissionsPolicy: editPermissionsPolicy(emptyGroupsMap),
  deletePermissionsPolicy,
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
    usersGroups(first: 500)
      @connection(key: "UserManagementContext_usersGroups") {
      edges {
        node {
          ...UserManagementUtils_group @relay(mask: false)
        }
      }
    }
    permissionsPolicies(first: 500)
      @connection(key: "UserManagementContext_permissionsPolicies") {
      edges {
        node {
          ...UserManagementUtils_policies @relay(mask: false)
        }
      }
    }
  }
`;

function ProviderWrap(props: Props) {
  const providerValue = (users, groups, policies, usersMap, groupsMap) => ({
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
    addPermissionsPolicy: addPermissionsPolicy(groupsMap),
    editPermissionsPolicy: editPermissionsPolicy(groupsMap),
    deletePermissionsPolicy,
  });

  const data = useLazyLoadQuery<UserManagementContextQuery>(dataQuery);

  const users = usersResponse2Users(data.users);
  const usersMap = users2UsersMap(users);
  const groups = groupsResponse2Groups(data.usersGroups, usersMap);
  const groupsMap = groups2GroupsMap(groups);
  const policies = permissionsPoliciesResponse2PermissionsPolicies(
    data.permissionsPolicies,
    groupsMap,
  );

  return (
    <UserManagementContext.Provider
      value={providerValue(users, groups, policies, usersMap, groupsMap)}>
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
