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

import type {
  UserManagementContextQueryResponse,
  UserRole,
  UserStatus,
  UsersGroupStatus,
} from '../__generated__/UserManagementContextQuery.graphql';

import fbt from 'fbt';

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
    value: `${fbt('Deactivated', '')}`,
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
type GroupsEdgesResponsePart = $ElementType<
  $NonMaybeType<GroupsReponsePart>,
  'edges',
>;
type GroupNodeReponseFieldsPart = $ElementType<GroupsEdgesResponsePart, number>;
type GroupReponseFieldsPart = $NonMaybeType<
  $ElementType<$NonMaybeType<GroupNodeReponseFieldsPart>, 'node'>,
>;
type PermissionsPoliciesReponsePart = $ElementType<
  UserManagementContextQueryResponse,
  'permissionsPolicies',
>;
type PoliciesEdgesResponsePart = $ElementType<
  $NonMaybeType<PermissionsPoliciesReponsePart>,
  'edges',
>;
type PolicyNodeReponseFieldsPart = $ElementType<
  PoliciesEdgesResponsePart,
  number,
>;
type PermissionsPoliciesReponseFieldsPart = $NonMaybeType<
  $ElementType<$NonMaybeType<PolicyNodeReponseFieldsPart>, 'node'>,
>;

export const userResponse2User: UsersReponseFieldsPart => User = (
  userNode: UsersReponseFieldsPart,
) => ({
  id: userNode.id,
  authID: userNode.authID,
  firstName: userNode.firstName,
  lastName: userNode.lastName,
  role: userNode.role,
  status: userNode.status,
  groups: userNode.groups ?? [],
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

export type UsersMap = Map<string, User>;
export const users2UsersMap: (Array<User>) => UsersMap = users =>
  new Map<string, User>(users.map(user => [user.id, user]));

export const userFullName = (user: User) =>
  `${user.firstName} ${user.lastName}`.trim() || '_';

export const groupResponse2Group: (
  GroupReponseFieldsPart,
  UsersMap,
) => UserPermissionsGroup = (groupResponse, usersMap) => ({
  id: groupResponse.id,
  name: groupResponse.name,
  description: groupResponse.description || '',
  status: groupResponse.status,
  members: groupResponse.members,
  memberUsers: groupResponse.members
    .map(member => usersMap.get(member.id))
    .filter(Boolean),
});

export const groupsResponse2Groups = (
  groupsResponse: GroupsReponsePart,
  usersMap: UsersMap,
) =>
  groupsResponse?.edges == null
    ? []
    : groupsResponse?.edges
        .filter(Boolean)
        .map(gr => gr.node)
        .filter(Boolean)
        .map<UserPermissionsGroup>(gr => groupResponse2Group(gr, usersMap));

export type PermissionsPolicy = $ReadOnly<{|
  id: string,
  name: string,
  description: string,
  type: PolicyTypes,
  isGlobal: boolean,
  groups: Array<UserPermissionsGroup>,
|}>;

// line was too long. So made it shorter...
type PPR2PP = PermissionsPoliciesReponseFieldsPart => PermissionsPolicy;
export const permissionsPolicyResponse2PermissionsPolicy: PPR2PP = (
  policyNode: PermissionsPoliciesReponseFieldsPart,
) => ({
  id: policyNode.id,
  name: policyNode.name,
  description: policyNode.description || '',
  type: policyNode.policy.__typename,
  isGlobal: policyNode.isGlobal,
  groups: [], // policyNode.groups,
});

export const permissionsPoliciesResponse2PermissionsPolicies = (
  policiesResponse: PermissionsPoliciesReponsePart,
) =>
  policiesResponse?.edges == null
    ? []
    : policiesResponse?.edges
        .filter(Boolean)
        .map(ur => ur.node)
        .filter(Boolean)
        .map<PermissionsPolicy>(permissionsPolicyResponse2PermissionsPolicy);
