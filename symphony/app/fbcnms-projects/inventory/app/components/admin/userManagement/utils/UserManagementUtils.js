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

export const NEW_GROUP_DIALOG_PARAM = 'new';

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

export type UserPermissionsGroupMember = {|
  +id: string,
  +authID: string,
|};
export type UserPermissionsGroup = {|
  id: string,
  name: string,
  description: string,
  status: UsersGroupStatus,
  members: $ReadOnlyArray<UserPermissionsGroupMember>,
  memberUsers: $ReadOnlyArray<User>,
|};
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
