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

import type {AddUsersGroupMutationResponse} from '../../../../mutations/__generated__/AddUsersGroupMutation.graphql';
import type {DeleteUsersGroupMutationResponse} from '../../../../mutations/__generated__/DeleteUsersGroupMutation.graphql';
import type {EditUsersGroupMutationResponse} from '../../../../mutations/__generated__/EditUsersGroupMutation.graphql';
import type {MutationCallbacks} from '../../../../mutations/MutationCallbacks.js';
import type {OptionalRefTypeWrapper} from '../../../../common/EntUtils';
import type {UsersGroupsQuery} from './__generated__/UsersGroupsQuery.graphql';
import type {UsersGroupsSearchQuery} from './__generated__/UsersGroupsSearchQuery.graphql';
import type {UserManagementUtils_group as group} from '../utils/__generated__/UserManagementUtils_group.graphql';
import type {UserManagementUtils_group_base as group_base} from '../utils/__generated__/UserManagementUtils_group_base.graphql';

import AddUsersGroupMutation from '../../../../mutations/AddUsersGroupMutation';
import DeleteUsersGroupMutation from '../../../../mutations/DeleteUsersGroupMutation';
import EditUsersGroupMutation from '../../../../mutations/EditUsersGroupMutation';
import {getGraphError} from '../../../../common/EntUtils';
import {graphql} from 'relay-runtime';
import {useLazyLoadQuery} from 'react-relay/hooks';

export type UsersGroup = OptionalRefTypeWrapper<group>;
export type UsersGroupBase = OptionalRefTypeWrapper<group_base>;

const groupsQuery = graphql`
  query UsersGroupsQuery {
    usersGroups(first: 500) @connection(key: "UsersGroups_usersGroups") {
      edges {
        node {
          ...UserManagementUtils_group @relay(mask: false)
        }
      }
    }
  }
`;

export function useUsersGroups(): $ReadOnlyArray<UsersGroup> {
  const data = useLazyLoadQuery<UsersGroupsQuery>(groupsQuery);
  const groupsData = data.usersGroups?.edges || [];
  return groupsData.map(p => p.node).filter(Boolean);
}

const groupQuery = graphql`
  query UsersGroupsSearchQuery($groupId: ID!) {
    group: node(id: $groupId) {
      ... on UsersGroup {
        ...UserManagementUtils_group @relay(mask: false)
      }
    }
  }
`;

export function useUsersGroup(groupId: string): UsersGroup {
  const data = useLazyLoadQuery<UsersGroupsSearchQuery>(groupQuery, {
    groupId,
  });
  return data.group;
}

export function deleteGroup(id: string): Promise<void> {
  return new Promise<void>((resolve, reject) => {
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
    DeleteUsersGroupMutation({id}, callbacks);
  });
}

export function editGroup(
  newGroupValue: UsersGroup,
): Promise<EditUsersGroupMutationResponse> {
  return new Promise<EditUsersGroupMutationResponse>((resolve, reject) => {
    const callbacks: MutationCallbacks<EditUsersGroupMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(response);
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
          policies: newGroupValue.policies.map(p => p.id),
        },
      },
      callbacks,
    );
  });
}

export function editGroupPolicies(
  id: string,
  policiesIDs: $ReadOnlyArray<string>,
): Promise<EditUsersGroupMutationResponse> {
  return new Promise<EditUsersGroupMutationResponse>((resolve, reject) => {
    const callbacks: MutationCallbacks<EditUsersGroupMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(response);
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };
    EditUsersGroupMutation(
      {
        input: {
          id,
          policies: policiesIDs,
        },
      },
      callbacks,
    );
  });
}

export function addGroup(
  newGroupValue: UsersGroup,
): Promise<AddUsersGroupMutationResponse> {
  return new Promise<AddUsersGroupMutationResponse>((resolve, reject) => {
    const callbacks: MutationCallbacks<AddUsersGroupMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(response);
      },
      onError: e => {
        reject(getGraphError(e));
      },
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
    );
  });
}
