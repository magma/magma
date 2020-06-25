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
  AddPermissionsPolicyMutationResponse,
  InventoryPolicyInput,
  LocationCUDInput,
  WorkforceCUDInput,
  WorkforcePermissionRuleInput,
  WorkforcePolicyInput,
} from '../../../../mutations/__generated__/AddPermissionsPolicyMutation.graphql';
import type {DeletePermissionsPolicyMutationResponse} from '../../../../mutations/__generated__/DeletePermissionsPolicyMutation.graphql';
import type {EditPermissionsPolicyMutationResponse} from '../../../../mutations/__generated__/EditPermissionsPolicyMutation.graphql';
import type {KeyValueEnum} from '../../../../common/EntUtils';
import type {MutationCallbacks} from '../../../../mutations/MutationCallbacks.js';
import type {OptionalRefTypeWrapper} from '../../../../common/EntUtils';
import type {
  PermissionValue,
  PermissionsPoliciesQuery,
  PermissionsPoliciesQueryResponse,
} from './__generated__/PermissionsPoliciesQuery.graphql';
import type {PermissionsPoliciesSearchQuery} from './__generated__/PermissionsPoliciesSearchQuery.graphql';
import type {UserManagementUtils_policies as policies} from '../utils/__generated__/UserManagementUtils_policies.graphql';
import type {UserManagementUtils_policies_base as policies_base} from '../utils/__generated__/UserManagementUtils_policies_base.graphql';

import AddPermissionsPolicyMutation from '../../../../mutations/AddPermissionsPolicyMutation';
import DeletePermissionsPolicyMutation from '../../../../mutations/DeletePermissionsPolicyMutation';
import EditPermissionsPolicyMutation from '../../../../mutations/EditPermissionsPolicyMutation';
import fbt from 'fbt';
import {getGraphError} from '../../../../common/EntUtils';
import {graphql} from 'relay-runtime';
import {useLazyLoadQuery} from 'react-relay/hooks';

export type PermissionsPolicyBase = OptionalRefTypeWrapper<policies_base>;
export type PermissionsPolicyRaw = OptionalRefTypeWrapper<policies>;
export type PermissionsPolicy = $ReadOnly<{|
  ...PermissionsPolicyRaw,
  type: PolicyTypes,
  inventoryRules?: ?InventoryPolicy,
  workforceRules?: ?WorkforcePolicy,
  isSystemDefault?: ?true,
|}>;

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

export const PERMISSION_RULE_VALUES = {
  YES: 'YES',
  NO: 'NO',
  BY_CONDITION: 'BY_CONDITION',
};

export type BasicPermissionRule = $ReadOnly<{|
  isAllowed: PermissionValue,
|}>;

type BasicCUDPermissions = $ReadOnly<{|
  create: BasicPermissionRule,
  delete: BasicPermissionRule,
|}>;

export type CUDPermissions = $ReadOnly<{|
  ...BasicCUDPermissions,
  update: BasicPermissionRule,
|}>;

export type InventoryCatalogPolicy = $ReadOnly<{|
  equipmentType: CUDPermissions,
  locationType: CUDPermissions,
  portType: CUDPermissions,
  serviceType: CUDPermissions,
|}>;

type LocationUpdatePermissionRule = $ReadOnly<{|
  ...BasicPermissionRule,
  locationTypeIds: ?$ReadOnlyArray<string>,
|}>;

export type LocationCUDPermissions = $ReadOnly<{|
  ...BasicCUDPermissions,
  update: LocationUpdatePermissionRule,
|}>;

export type InventoryEntsPolicy = $ReadOnly<{|
  location: LocationCUDPermissions,
  equipment: CUDPermissions,
  ...InventoryCatalogPolicy,
|}>;

export type InventoryPolicy = $ReadOnly<{|
  read: BasicPermissionRule,
  ...InventoryEntsPolicy,
|}>;

export type WorkforceCUDPermissions = $ReadOnly<{|
  create: BasicPermissionRule,
  update: BasicPermissionRule,
  delete: BasicPermissionRule,
  assign: BasicPermissionRule,
  transferOwnership: BasicPermissionRule,
|}>;

export type WorkforceReadPermissionRule = $ReadOnly<{|
  isAllowed: PermissionValue,
  projectTypeIds: ?$ReadOnlyArray<string>,
  workOrderTypeIds: ?$ReadOnlyArray<string>,
|}>;

export type WorkforcePolicy = $ReadOnly<{|
  read: WorkforceReadPermissionRule,
  data: WorkforceCUDPermissions,
  templates: CUDPermissions,
|}>;

export type PermissionsPolicyRules = InventoryPolicy | WorkforcePolicy | {||};

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

function permissionsPolicy2PermissionsPolicyInput(
  policy: PermissionsPolicy,
): AddPermissionsPolicyInput {
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
    isGlobal: policy.isGlobal,
    groups: policy.groups.map(group => group.id),
  };
}

type PermissionPolicyBasicRuleInput = {|
  ...BasicPermissionRule,
|};

function permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
  rule: ?BasicPermissionRule,
): PermissionPolicyBasicRuleInput {
  return {
    isAllowed: parsePermissionValue(rule?.isAllowed),
  };
}

function parsePermissionValue(
  permissionValue?: ?PermissionValue,
): PermissionValue {
  return permissionValue ?? PERMISSION_RULE_VALUES.NO;
}

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

export function permissionPolicyWFCUDRule2PermissionPolicyWFCUDRuleInput(
  rule: ?WorkforceCUDPermissions,
): WorkforceCUDInput {
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
    assign: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      rule?.assign,
    ),
    transferOwnership: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      rule?.transferOwnership,
    ),
  };
}

export const initInventoryRulesInput: (
  ?InventoryPolicy,
) => InventoryPolicyInput = (policyRules?: ?InventoryPolicy) => {
  return {
    read: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      policyRules?.read,
    ),
    location: locationPolicyCUDRule2LocationPolicyCUDRuleInput(
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

export function locationPolicyCUDRule2LocationPolicyCUDRuleInput(
  rule: ?LocationCUDPermissions,
): LocationCUDInput {
  const updateIsAllowedValue = parsePermissionValue(rule?.update.isAllowed);

  return {
    create: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      rule?.create,
    ),
    update: {
      isAllowed: updateIsAllowedValue,
      locationTypeIds:
        updateIsAllowedValue === PERMISSION_RULE_VALUES.BY_CONDITION
          ? rule?.update.locationTypeIds
          : null,
    },
    delete: permissionPolicyBasicRule2PermissionPolicyBasicRuleInput(
      rule?.delete,
    ),
  };
}

function workforceReadPermissionRule2WorkforcePermissionRuleInput(
  policyRule?: ?WorkforceReadPermissionRule,
): WorkforcePermissionRuleInput {
  return {
    isAllowed: parsePermissionValue(policyRule?.isAllowed),
    projectTypeIds: policyRule?.projectTypeIds,
    workOrderTypeIds: policyRule?.workOrderTypeIds,
  };
}

export function initWorkforceRulesInput(
  policyRules?: ?WorkforcePolicy,
): WorkforcePolicyInput {
  return {
    read: workforceReadPermissionRule2WorkforcePermissionRuleInput(
      policyRules?.read,
    ),
    data: permissionPolicyWFCUDRule2PermissionPolicyWFCUDRuleInput(
      policyRules?.data,
    ),
    templates: permissionPolicyCUDRule2PermissionPolicyCUDRuleInput(
      policyRules?.templates,
    ),
  };
}

export function bool2PermissionRuleValue(value: ?boolean): PermissionValue {
  return value === true
    ? PERMISSION_RULE_VALUES.YES
    : PERMISSION_RULE_VALUES.NO;
}

export function permissionRuleValue2Bool(value: PermissionValue) {
  return value !== PERMISSION_RULE_VALUES.NO;
}

function response2PermissionsPolicy(
  policyResponse: PermissionsPolicyRaw | PermissionsPolicyBase,
): PermissionsPolicy {
  const {__typename: type, ...policyRules} = policyResponse.policy;
  return {
    groups: [],
    ...policyResponse,
    type,
    inventoryRules: tryGettingInventoryPolicy(policyRules),
    workforceRules: tryGettingWorkforcePolicy(policyRules),
  };
}

export function wrapRawPermissionsPolicies(
  rawPolicies: $ReadOnlyArray<PermissionsPolicyRaw | PermissionsPolicyBase>,
): Array<PermissionsPolicy> {
  return rawPolicies.map(policy => response2PermissionsPolicy(policy));
}

function permissionsPolicy2PermissionsPolicyBase(
  policy: PermissionsPolicy,
): PermissionsPolicyBase {
  const {
    type: _type,
    groups: _groups,
    inventoryRules: _inventoryRules,
    workforceRules: _workforceRules,
    isSystemDefault: _isSystemDefault,
    ...rawPolicy
  } = policy;
  return rawPolicy;
}

export function unwrapPermissionsPolicies(
  policies: $ReadOnlyArray<PermissionsPolicy>,
): $ReadOnlyArray<PermissionsPolicyBase> {
  return policies.map(policy =>
    permissionsPolicy2PermissionsPolicyBase(policy),
  );
}

export const EMPTY_POLICY = {
  __typename: '%other',
};

export const WORKORDER_SYSTEM_POLICY_ID = 'system_workorder';
export const WORKORDER_SYSTEM_POLICY = {
  id: WORKORDER_SYSTEM_POLICY_ID,
  name: `${fbt('Work orders editing', '')}`,
  description: `${fbt(
    'All active users can view and edit work orders and projects assigned to and owned by them (including changing assignment). An active user who owns the work order can transfer ownership to other user and even delete it. When a work order is part of a project, that project will be visible as well.',
    '',
  )}`,
  type: POLICY_TYPES.WorkforcePolicy.key,
  policy: EMPTY_POLICY,
  isGlobal: true,
  groups: [],
  isSystemDefault: true,
};

function response2PermissionsPolicies(
  policiesResponse: PermissionsPoliciesQueryResponse,
): Array<PermissionsPolicy> {
  const policiesData = policiesResponse.permissionsPolicies?.edges || [];
  const policies = wrapRawPermissionsPolicies(
    policiesData.map(p => p.node).filter(Boolean),
  );

  policies.unshift(WORKORDER_SYSTEM_POLICY);

  return policies;
}

export function addPermissionsPolicy(
  newPolicyValue: PermissionsPolicy,
): Promise<PermissionsPolicy> {
  return new Promise<PermissionsPolicy>((resolve, reject) => {
    const callbacks: MutationCallbacks<AddPermissionsPolicyMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(response2PermissionsPolicy(response.addPermissionsPolicy));
      },
      onError: e => {
        reject(getGraphError(e));
      },
    };
    AddPermissionsPolicyMutation(
      {
        input: permissionsPolicy2PermissionsPolicyInput(newPolicyValue),
      },
      callbacks,
    );
  });
}

export function editPermissionsPolicy(
  newPolicyValue: PermissionsPolicy,
): Promise<PermissionsPolicy> {
  return new Promise<PermissionsPolicy>((resolve, reject) => {
    type Callbacks = MutationCallbacks<EditPermissionsPolicyMutationResponse>;
    const callbacks: Callbacks = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          reject(getGraphError(errors[0]));
        }
        resolve(response2PermissionsPolicy(response.editPermissionsPolicy));
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
}

export function deletePermissionsPolicy(id: string) {
  return new Promise<void>((resolve, reject) => {
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
    DeletePermissionsPolicyMutation({id}, cbs);
  });
}

const policiesQuery = graphql`
  query PermissionsPoliciesQuery {
    permissionsPolicies(first: 500)
      @connection(key: "PermissionsPoliciesQuery_permissionsPolicies") {
      edges {
        node {
          ...UserManagementUtils_policies @relay(mask: false)
        }
      }
    }
  }
`;

export function usePermissionsPolicies(): $ReadOnlyArray<PermissionsPolicy> {
  const data = useLazyLoadQuery<PermissionsPoliciesQuery>(policiesQuery);
  return response2PermissionsPolicies(data);
}

const policyQuery = graphql`
  query PermissionsPoliciesSearchQuery($policyId: ID!) {
    policy: node(id: $policyId) {
      ... on PermissionsPolicy {
        ...UserManagementUtils_policies @relay(mask: false)
      }
    }
  }
`;

export function usePermissionsPolicy(policyId: string): ?PermissionsPolicy {
  const data = useLazyLoadQuery<PermissionsPoliciesSearchQuery>(policyQuery, {
    policyId,
  });
  return data.policy == null ? null : response2PermissionsPolicy(data.policy);
}
