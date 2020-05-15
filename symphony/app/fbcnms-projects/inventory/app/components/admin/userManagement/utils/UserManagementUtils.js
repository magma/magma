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
  AddPermissionsPolicyInput,
  BasicPermissionRuleInput,
  InventoryPolicyInput,
  LocationCUDInput,
  WorkforceCUDInput,
  WorkforcePermissionRuleInput,
  WorkforcePolicyInput,
} from '../../../../mutations/__generated__/AddPermissionsPolicyMutation.graphql';
import type {GroupSearchContextQueryResponse} from './search/__generated__/GroupSearchContextQuery.graphql';
import type {
  PermissionValue,
  UserManagementContextQueryResponse,
  UserRole,
  UserStatus,
  UsersGroupStatus,
} from '../__generated__/UserManagementContextQuery.graphql';

import fbt from 'fbt';
import {graphql} from 'relay-runtime';

graphql`
  fragment UserManagementUtils_user on User {
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
`;

graphql`
  fragment UserManagementUtils_group on UsersGroup {
    id
    name
    description
    status
    members {
      id
      authID
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
  fragment UserManagementUtils_policies on PermissionsPolicy {
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
    groups {
      id
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
  ?UsersMap,
) => UserPermissionsGroup = (
  groupResponse: GroupReponseFieldsPart,
  usersMap?: ?UsersMap,
) => ({
  id: groupResponse.id,
  name: groupResponse.name,
  description: groupResponse.description || '',
  status: groupResponse.status,
  members: groupResponse.members,
  memberUsers:
    usersMap == null
      ? []
      : groupResponse.members
          .map(member => usersMap.get(member.id))
          .filter(Boolean),
});

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
    .map<UserPermissionsGroup>(gr => groupResponse2Group(gr, usersMap));
};

export const PERMISSION_RULE_VALUES = {
  YES: 'YES',
  NO: 'NO',
  BY_CONDITION: 'BY_CONDITION',
};

export type BasicPermissionRule = $ReadOnly<{|
  isAllowed: PermissionValue,
|}>;

export type CUDPermissions = $ReadOnly<{|
  create: BasicPermissionRule,
  update: BasicPermissionRule,
  delete: BasicPermissionRule,
|}>;

export type InventoryCatalogPolicy = $ReadOnly<{|
  equipmentType: CUDPermissions,
  locationType: CUDPermissions,
  portType: CUDPermissions,
  serviceType: CUDPermissions,
|}>;

export type InventoryPolicy = $ReadOnly<{|
  read: BasicPermissionRule,
  location: CUDPermissions,
  equipment: CUDPermissions,
  ...InventoryCatalogPolicy,
|}>;

export type WorkforceBasicPermissions = BasicPermissionRule &
  $ReadOnly<{|
    ...BasicPermissionRule,
    projectTypeIds?: ?$ReadOnlyArray<string>,
    workOrderTypeIds?: ?$ReadOnlyArray<string>,
  |}>;

export type WorkforceCUDPermissions = $ReadOnly<{|
  create: WorkforceBasicPermissions,
  update: WorkforceBasicPermissions,
  delete: WorkforceBasicPermissions,
  assign: WorkforceBasicPermissions,
  transferOwnership: WorkforceBasicPermissions,
|}>;

export type WorkforcePolicy = $ReadOnly<{|
  read: BasicPermissionRule,
  data: WorkforceCUDPermissions,
  templates: CUDPermissions,
|}>;

export type PermissionsPolicyRules = InventoryPolicy | WorkforcePolicy | {||};
export type PermissionsPolicy = $ReadOnly<{|
  id: string,
  name: string,
  description: string,
  type: PolicyTypes,
  isGlobal: boolean,
  groups: Array<UserPermissionsGroup>,
  inventoryRules?: ?InventoryPolicy,
  workforceRules?: ?WorkforcePolicy,
|}>;

function tryGettingInventoryPolicy(
  policyRules: ?PermissionsPolicyRules,
): ?InventoryPolicy {
  if (policyRules == null) {
    return null;
  }

  if (
    policyRules.read &&
    policyRules.location &&
    policyRules.equipment &&
    policyRules.equipmentType &&
    policyRules.locationType &&
    policyRules.portType &&
    policyRules.serviceType
  ) {
    return policyRules;
  }

  return null;
}

function tryGettingWorkforcePolicy(
  policyRules: ?PermissionsPolicyRules,
): ?WorkforcePolicy {
  if (policyRules == null) {
    return null;
  }

  if (policyRules.read && policyRules.data && policyRules.templates) {
    return policyRules;
  }

  return null;
}

// line was too long. So made it shorter...
type PPR2PP = PermissionsPoliciesReponseFieldsPart => PermissionsPolicy;
export const permissionsPolicyResponse2PermissionsPolicy: PPR2PP = (
  policyNode: PermissionsPoliciesReponseFieldsPart,
) => {
  const {__typename: type, ...policyRules} = policyNode.policy;
  return {
    id: policyNode.id,
    name: policyNode.name,
    description: policyNode.description || '',
    type,
    isGlobal: policyNode.isGlobal,
    groups: [], // policyNode.groups,
    inventoryRules: tryGettingInventoryPolicy(policyRules),
    workforceRules: tryGettingWorkforcePolicy(policyRules),
  };
};

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

export const permissionPolicyBasicRule2PermissionPolicyBasicRuleInput: (
  ?BasicPermissionRule,
) => BasicPermissionRuleInput = (rule: ?BasicPermissionRule) => {
  return {
    isAllowed: rule?.isAllowed ?? PERMISSION_RULE_VALUES.NO,
  };
};

export const permissionPolicyCUDRule2PermissionPolicyCUDRuleInput = (
  rule: ?CUDPermissions,
) => {
  return {
    create: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      rule?.create,
    ),
    update: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      rule?.update,
    ),
    delete: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      rule?.delete,
    ),
  };
};

export const permissionPolicyWFCUDRule2PermissionPolicyWFCUDRuleInput: (
  ?WorkforceCUDPermissions,
) => WorkforceCUDInput = (rule: ?WorkforceCUDPermissions) => {
  return {
    create: wfPermissionPolicyBasicRule2wfPermissionRuleInput(rule?.create),
    update: wfPermissionPolicyBasicRule2wfPermissionRuleInput(rule?.update),
    delete: wfPermissionPolicyBasicRule2wfPermissionRuleInput(rule?.delete),
    assign: wfPermissionPolicyBasicRule2wfPermissionRuleInput(rule?.assign),
    transferOwnership: wfPermissionPolicyBasicRule2wfPermissionRuleInput(
      rule?.transferOwnership,
    ),
  };
};

export const wfPermissionPolicyBasicRule2wfPermissionRuleInput: (
  ?WorkforceBasicPermissions,
) => WorkforcePermissionRuleInput = (rule: ?BasicPermissionRule) => {
  return {
    isAllowed: rule?.isAllowed ?? PERMISSION_RULE_VALUES.NO,
  };
};

export const initInventoryRulesInput: (
  ?InventoryPolicy,
) => InventoryPolicyInput = (policyRules?: ?InventoryPolicy) => {
  return {
    read: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      policyRules?.read,
    ),
    location: permissionPolicyCUDRule2LocationPermissionPolicyCUDRuleInput(
      policyRules?.location,
    ),
    equipment: permissionPolicyCUDRule2PermissionPolicyCUDRuleInput(
      policyRules?.equipment,
    ),
    equipmentType: permissionPolicyCUDRule2PermissionPolicyCUDRuleInput(
      policyRules?.equipmentType,
    ),
    locationType: permissionPolicyCUDRule2PermissionPolicyCUDRuleInput(
      policyRules?.locationType,
    ),
    portType: permissionPolicyCUDRule2PermissionPolicyCUDRuleInput(
      policyRules?.portType,
    ),
    serviceType: permissionPolicyCUDRule2PermissionPolicyCUDRuleInput(
      policyRules?.serviceType,
    ),
  };
};

export const permissionPolicyCUDRule2LocationPermissionPolicyCUDRuleInput: (
  ?CUDPermissions,
) => LocationCUDInput = (rule: ?CUDPermissions) => {
  const partialInput = permissionPolicyCUDRule2PermissionPolicyCUDRuleInput(
    rule,
  );
  return {
    create: {
      ...partialInput.create,
    },
    update: {
      ...partialInput.update,
    },
    delete: {
      ...partialInput.delete,
    },
  };
};

export const initWorkforceRulesInput: (
  ?WorkforcePolicy,
) => WorkforcePolicyInput = (policyRules?: ?WorkforcePolicy) => {
  return {
    read: wfPermissionPolicyBasicRule2wfPermissionRuleInput(policyRules?.read),
    data: permissionPolicyWFCUDRule2PermissionPolicyWFCUDRuleInput(
      policyRules?.data,
    ),
    templates: permissionPolicyCUDRule2PermissionPolicyCUDRuleInput(
      policyRules?.templates,
    ),
  };
};

// line was too long. So made it shorter...
type PP2APPI = PermissionsPolicy => AddPermissionsPolicyInput;
export const permissionsPolicy2PermissionsPolicyInput: PP2APPI = (
  policy: PermissionsPolicy,
) => {
  return {
    name: policy.name,
    description: policy.description,
    inventoryInput:
      policy.type === POLICY_TYPES.InventoryPolicy.key
        ? initInventoryRulesInput(policy.inventoryRules)
        : null,
    workforceInput:
      policy.type === POLICY_TYPES.WorkforcePolicy.key
        ? initWorkforceRulesInput(policy.workforceRules)
        : null,
  };
};

export function bool2PermissionRuleValue(value: ?boolean): PermissionValue {
  return value === true
    ? PERMISSION_RULE_VALUES.YES
    : PERMISSION_RULE_VALUES.NO;
}

export function permissionRuleValue2Bool(value: PermissionValue) {
  return value === PERMISSION_RULE_VALUES.YES;
}
