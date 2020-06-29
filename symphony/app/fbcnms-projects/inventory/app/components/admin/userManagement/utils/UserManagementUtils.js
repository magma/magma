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

import type {OptionalRefTypeWrapper} from '../../../../common/EntUtils';
import type {UserManagementUtils_user} from './__generated__/UserManagementUtils_user.graphql';
import type {UserManagementUtils_user_base} from './__generated__/UserManagementUtils_user_base.graphql';
import type {
  UserRole,
  UserStatus,
} from './__generated__/UserManagementUtils_user.graphql';
import type {UsersGroupStatus} from './__generated__/UserManagementUtils_group.graphql';

import fbt from 'fbt';
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
        locationTypeIds
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
      projectTypeIds
      workOrderTypeIds
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

const ADMIN_LABEL = `${fbt('Admin', '')}`;

export const USER_ROLES: KeyValueEnum<UserRole> = {
  USER: {
    key: 'USER',
    value: `${fbt('User', '')}`,
  },
  ADMIN: {
    key: 'ADMIN',
    value: ADMIN_LABEL,
  },
  OWNER: {
    key: 'OWNER',
    value: ADMIN_LABEL,
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

export type UserBase = OptionalRefTypeWrapper<UserManagementUtils_user_base>;
export type User = OptionalRefTypeWrapper<UserManagementUtils_user>;

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

export const userFullName = (user: $Shape<User>) =>
  `${user.firstName} ${user.lastName}`.trim() || '_';
