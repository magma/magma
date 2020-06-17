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

import type {EntsMap} from '../../../../common/EntUtils';
import type {GroupSearchContextQueryResponse} from './search/__generated__/GroupSearchContextQuery.graphql';
import type {
  UserManagementContextQueryResponse,
  UserRole,
  UserStatus,
  UsersGroupStatus,
} from '../__generated__/UserManagementContextQuery.graphql';

import fbt from 'fbt';
import {ent2EntsMap} from '../../../../common/EntUtils';
import {graphql} from 'relay-runtime';

graphql`
  fragment UserManagementUtils_user_base on User {
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
`;

graphql`
  fragment UserManagementUtils_user on User {
    ...UserManagementUtils_user_base @relay(mask: false)
    groups {
      ...UserManagementUtils_group_base @relay(mask: false)
    }
  }
`;

graphql`
  fragment UserManagementUtils_group_base on UsersGroup {
    id
    name
    description
    status
  }
`;

graphql`
  fragment UserManagementUtils_group on UsersGroup {
    ...UserManagementUtils_group_base @relay(mask: false)
    members {
      ...UserManagementUtils_user_base @relay(mask: false)
    }
    policies {
      ...UserManagementUtils_policies_base @relay(mask: false)
    }
  }
`;

graphql`
  fragment UserManagementUtils_inventoryPolicy on InventoryPolicy {
    read {
      isAllowed
    }
    location {
      create {
        isAllowed
      }
      update {
        isAllowed
      }
      delete {
        isAllowed
      }
    }
    equipment {
      create {
        isAllowed
      }
      update {
        isAllowed
      }
      delete {
        isAllowed
      }
    }
    equipmentType {
      create {
        isAllowed
      }
      update {
        isAllowed
      }
      delete {
        isAllowed
      }
    }
    locationType {
      create {
        isAllowed
      }
      update {
        isAllowed
      }
      delete {
        isAllowed
      }
    }
    portType {
      create {
        isAllowed
      }
      update {
        isAllowed
      }
      delete {
        isAllowed
      }
    }
    serviceType {
      create {
        isAllowed
      }
      update {
        isAllowed
      }
      delete {
        isAllowed
      }
    }
  }
`;

graphql`
  fragment UserManagementUtils_workforcePolicy on WorkforcePolicy {
    read {
      isAllowed
    }
    templates {
      create {
        isAllowed
      }
      update {
        isAllowed
      }
      delete {
        isAllowed
      }
    }
    data {
      create {
        isAllowed
      }
      update {
        isAllowed
      }
      delete {
        isAllowed
      }
      assign {
        isAllowed
      }
      transferOwnership {
        isAllowed
      }
    }
  }
`;

graphql`
  fragment UserManagementUtils_policies_base on PermissionsPolicy {
    id
    name
    description
    isGlobal
    policy {
      __typename
      ... on InventoryPolicy {
        ...UserManagementUtils_inventoryPolicy @relay(mask: false)
      }
      ... on WorkforcePolicy {
        ...UserManagementUtils_workforcePolicy @relay(mask: false)
      }
    }
  }
`;

graphql`
  fragment UserManagementUtils_policies on PermissionsPolicy {
    ...UserManagementUtils_policies_base @relay(mask: false)
    groups {
      ...UserManagementUtils_group @relay(mask: false)
    }
  }
`;

type KeyValueEnum<TValues> = {
  [key: TValues]: {
    key: TValues,
    value: string,
  },
};

export const USER_ROLES: KeyValueEnum<UserRole> = {
  USER: {
    key: 'USER',
    value: `${fbt('User', '')}`,
  },
  ADMIN: {
    key: 'ADMIN',
    value: `${fbt('Admin', '')}`,
  },
  OWNER: {
    key: 'OWNER',
    value: `${fbt('Owner', '')}`,
  },
};
export const USER_STATUSES: KeyValueEnum<UserStatus> = {
  ACTIVE: {
    key: 'ACTIVE',
    value: `${fbt('Active', '')}`,
  },
  DEACTIVATED: {
    key: 'DEACTIVATED',
    value: `${fbt('Deactivated', '')}`,
  },
};
export const EMPLOYMENT_TYPES: KeyValueEnum<string> = {
  FULL_TIME: {
    key: 'FULL_TIME',
    value: `${fbt('Full Time', '')}`,
  },
  CONTRUCTOR: {
    key: 'CONTRACTOR',
    value: `${fbt('Contractor', '')}`,
  },
};

export type EmploymentType = $Keys<typeof EMPLOYMENT_TYPES>;

export type UserGroups = {|
  +id: string,
  +name: string,
|};

export type User = {|
  id: string,
  authID: string,
  firstName: string,
  lastName: string,
  role: UserRole,
  status: UserStatus,
  photoId?: string,
  employmentType?: EmploymentType,
  employeeID?: string,
  jobTitle?: string,
  phoneNumber?: string,
  groups: $ReadOnlyArray<?UserGroups>,
|};

export const PermissionValues = {
  YES: 'YES',
  BY_CONDITION: 'BY_CONDITION',
  NO: 'NO',
};

export const NEW_DIALOG_PARAM = 'new';

export const GROUP_STATUSES: KeyValueEnum<UsersGroupStatus> = {
  ACTIVE: {
    key: 'ACTIVE',
    value: `${fbt('Active', '')}`,
  },
  DEACTIVATED: {
    key: 'DEACTIVATED',
    value: `${fbt('Inactive', '')}`,
  },
};

export type PolicyTypes = 'InventoryPolicy' | 'WorkforcePolicy' | '%other';
export const POLICY_TYPES: KeyValueEnum<PolicyTypes> = {
  InventoryPolicy: {
    key: 'InventoryPolicy',
    value: `${fbt('Inventory', '')}`,
  },
  WorkforcePolicy: {
    key: 'WorkforcePolicy',
    value: `${fbt('Workforce', '')}`,
  },
};

export type UserPermissionsGroupMember = $ReadOnly<{|
  +id: string,
  +authID: string,
|}>;
export type UserPermissionsGroup = $ReadOnly<{|
  id: string,
  name: string,
  description: string,
  status: UsersGroupStatus,
  members: $ReadOnlyArray<UserPermissionsGroupMember>,
  memberUsers: $ReadOnlyArray<User>,
|}>;
type UsersReponsePart = $ElementType<
  UserManagementContextQueryResponse,
  'users',
>;
type UsersEdgesResponsePart = $ElementType<
  $NonMaybeType<UsersReponsePart>,
  'edges',
>;
type UserNodeReponseFieldsPart = $ElementType<UsersEdgesResponsePart, number>;
type UsersReponseFieldsPart = $NonMaybeType<
  $ElementType<$NonMaybeType<UserNodeReponseFieldsPart>, 'node'>,
>;
type GroupsReponsePart = $ElementType<
  UserManagementContextQueryResponse,
  'usersGroups',
>;
type GroupsSearchReponsePart = $ElementType<
  GroupSearchContextQueryResponse,
  'usersGroupSearch',
>;
type GroupsEdgesResponsePart = $ElementType<
  $NonMaybeType<GroupsReponsePart>,
  'edges',
>;
type GroupNodeReponseFieldsPart = $ElementType<GroupsEdgesResponsePart, number>;
type GroupReponseFieldsPart = $NonMaybeType<
  $ElementType<$NonMaybeType<GroupNodeReponseFieldsPart>, 'node'>,
>;

export type UsersGroup = GroupReponseFieldsPart;

export const userResponse2User: UsersReponseFieldsPart => User = (
  userNode: UsersReponseFieldsPart,
) => ({
  id: userNode.id,
  authID: userNode.authID,
  firstName: userNode.firstName,
  lastName: userNode.lastName,
  role: userNode.role,
  status: userNode.status,
  groups: /* userNode.groups ?? */ [],
  photoId: userNode.profilePhoto?.id,
});

export const usersResponse2Users = (usersResponse: UsersReponsePart) =>
  usersResponse?.edges == null
    ? []
    : usersResponse?.edges
        .filter(Boolean)
        .map(ur => ur.node)
        .filter(Boolean)
        .map<User>(userResponse2User);

export type UsersMap = EntsMap<User>;
export const users2UsersMap = (users: Array<User>) => ent2EntsMap<User>(users);

export const userFullName = (user: User) =>
  `${user.firstName} ${user.lastName}`.trim() || '_';

export function groupResponse2Group(
  usersMap?: ?UsersMap,
): GroupReponseFieldsPart => UserPermissionsGroup {
  return (groupResponse: GroupReponseFieldsPart) => ({
    id: groupResponse.id,
    name: groupResponse.name,
    description: groupResponse.description || '',
    status: groupResponse.status,
    members: /* groupResponse.members*/ [],
    memberUsers:
      usersMap == null
        ? []
        : groupResponse.members
            .map(member => usersMap.get(member.id))
            .filter(Boolean),
  });
}

export const groupsResponse2Groups = (
  groupsResponse: GroupsReponsePart | GroupsSearchReponsePart,
  usersMap?: ?UsersMap,
) => {
  if (groupsResponse == null) {
    return [];
  }
  const resposeNodes =
    groupsResponse.edges != null
      ? groupsResponse?.edges.filter(Boolean).map(gr => gr.node)
      : groupsResponse.usersGroups != null
      ? groupsResponse.usersGroups
      : [];

  return resposeNodes
    .filter(Boolean)
    .map<UserPermissionsGroup>(groupResponse2Group(usersMap));
};

export type GroupsMap = EntsMap<UserPermissionsGroup>;
export const groups2GroupsMap = (groups: Array<UserPermissionsGroup>) =>
  ent2EntsMap<UserPermissionsGroup>(groups);
