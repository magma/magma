/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {BasicPermissionRule} from './UserManagementUtils';
import type {
  CUDPermissions,
  InventoryEntsPolicy,
  WorkforceCUDPermissions,
} from './UserManagementUtils';
import type {FormAlertsContextType} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import type {PermissionHandlingProps} from '@fbcnms/ui/components/design-system/Form/FormAction';
import type {UserPermissions} from '../../../MainContext';

import fbt from 'fbt';
import useFeatureFlag from '@fbcnms/ui/context/useFeatureFlag';
import {permissionRuleValue2Bool} from './UserManagementUtils';
import {useFormAlertsContext} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import {useMainContext} from '../../../MainContext';

type BasePermissionEnforcement = $ReadOnly<{|
  ...PermissionHandlingProps,
|}>;

type AdminPermissionEnforcement = $ReadOnly<{|
  ...BasePermissionEnforcement,
  adminRightsRequired: true,
|}>;

type InventoryEntWithPermission = $Keys<InventoryEntsPolicy>;
type InventoryEntName = InventoryEntWithPermission | 'service' | 'port';
type WorkforceTemplateEntName = 'workorderTemplate' | 'projectTemplate';
type InventoryActionName = $Keys<CUDPermissions>;

type InventoryPermissionEnforcement = $ReadOnly<{|
  ...BasePermissionEnforcement,
  entity: InventoryEntName | WorkforceTemplateEntName,
  action?: ?InventoryActionName,
|}>;

type WorkforceEntName = 'workorder' | 'project';

type WorkforceActionName = $Keys<WorkforceCUDPermissions>;

type WorkforcePermissionEnforcement = $ReadOnly<{|
  ...BasePermissionEnforcement,
  entity: WorkforceEntName,
  action?: ?WorkforceActionName,
|}>;

export type EntName = InventoryEntName | WorkforceEntName;
export type ActionName = InventoryActionName | WorkforceActionName;

export type PermissionEnforcement =
  | BasePermissionEnforcement
  | AdminPermissionEnforcement
  | InventoryPermissionEnforcement
  | WorkforcePermissionEnforcement;

const PASSED_VALUE = '';
const FAILED_ADMIN_VALUE = `${fbt(
  'Admin rights are required. Contact your system administrator.',
  '',
)}`;
const FAILED_REGULAR_VALUE = `${fbt(
  "User doesn't have sufficient permissions for this action.",
  '',
)}`;

const noEnforcement = () => (
  _permissionsEnforcement: PermissionEnforcement,
  _aggregationKey?: ?string,
) => '';

const performCheck = (
  userPermissions: UserPermissions,
  permissionsEnforcement: PermissionEnforcement,
) => {
  if (permissionsEnforcement.ignorePermissions === true) {
    return PASSED_VALUE;
  }

  if (permissionsEnforcement.adminRightsRequired === true) {
    const adminPermission = userPermissions.adminPolicy.access.isAllowed;
    return permissionRuleValue2Bool(adminPermission)
      ? PASSED_VALUE
      : FAILED_ADMIN_VALUE;
  }

  if (!permissionsEnforcement.entity || !permissionsEnforcement.action) {
    return PASSED_VALUE;
  }

  let actionPermissionValue: ?BasicPermissionRule;

  if (
    permissionsEnforcement.entity === 'workorder' ||
    permissionsEnforcement.entity === 'project'
  ) {
    const action: WorkforceActionName = permissionsEnforcement.action;
    actionPermissionValue = userPermissions.workforcePolicy.data[action];
  } else if (
    permissionsEnforcement.entity === 'workorderTemplate' ||
    permissionsEnforcement.entity === 'projectTemplate'
  ) {
    const enforcement: InventoryPermissionEnforcement = permissionsEnforcement;
    if (!enforcement.action) {
      return PASSED_VALUE;
    }
    actionPermissionValue =
      userPermissions.workforcePolicy.templates[enforcement.action];
  } else if (
    permissionsEnforcement.entity === 'location' ||
    permissionsEnforcement.entity === 'equipment' ||
    permissionsEnforcement.entity === 'port' ||
    permissionsEnforcement.entity === 'service' ||
    permissionsEnforcement.entity === 'locationType' ||
    permissionsEnforcement.entity === 'equipmentType' ||
    permissionsEnforcement.entity === 'portType' ||
    permissionsEnforcement.entity === 'serviceType'
  ) {
    const enforcement: InventoryPermissionEnforcement = permissionsEnforcement;
    const subjectEntity: InventoryEntWithPermission =
      permissionsEnforcement.entity === 'service' ||
      permissionsEnforcement.entity === 'port'
        ? 'equipment'
        : permissionsEnforcement.entity;
    if (!enforcement.action) {
      return PASSED_VALUE;
    }
    actionPermissionValue =
      userPermissions.inventoryPolicy[subjectEntity][enforcement.action];
  }

  if (actionPermissionValue == null) {
    return PASSED_VALUE;
  }

  return permissionRuleValue2Bool(actionPermissionValue.isAllowed)
    ? PASSED_VALUE
    : FAILED_REGULAR_VALUE;
};

const enforcePermissions = (
  userPermissions: UserPermissions,
  formAlertsContext: FormAlertsContextType,
) => (
  permissionsEnforcement: PermissionEnforcement,
  aggregationKey?: ?string,
) => {
  const checkResult = performCheck(userPermissions, permissionsEnforcement);

  if (formAlertsContext.isInitialized && aggregationKey != null) {
    formAlertsContext.missingPermissions.check({
      fieldId: aggregationKey,
      fieldDisplayName: 'Permissions Check',
      value: checkResult,
      checkCallback: () => checkResult,
    });
  }

  return checkResult;
};

export default function usePermissions() {
  const {me} = useMainContext();
  const formAlerts = useFormAlertsContext();
  const permissionsEnforcementIsOn = useFeatureFlag(
    'permissions_ui_enforcement',
  );

  const check =
    permissionsEnforcementIsOn && me != null
      ? enforcePermissions(me.permissions, formAlerts)
      : noEnforcement();

  return {
    check,
  };
}
