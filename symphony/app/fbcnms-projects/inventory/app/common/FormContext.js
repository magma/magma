/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  AlertRuleCheck,
  FormAlertsContextType,
} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import type {
  BasicPermissionRule,
  CUDPermissions,
  InventoryEntsPolicy,
  WorkforceCUDPermissions,
} from '../components/admin/userManagement/utils/UserManagementUtils';
import type {UserPermissions} from '../components/MainContext';

import * as React from 'react';
import FormAlertsContext, {
  DEFAULT_CONTEXT_VALUE as DEFAULT_ALERTS,
  FormAlertsContextProvider,
} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import fbt from 'fbt';
import useFeatureFlag from '@fbcnms/ui/context/useFeatureFlag';
import {createContext, useContext} from 'react';
import {permissionRuleValue2Bool} from '../components/admin/userManagement/utils/UserManagementUtils';
import {useMainContext} from '../components/MainContext';

type FromContextType = $ReadOnly<{|
  alerts: FormAlertsContextType,
|}>;

const DEFAULT_CONTEXT_VALUE = {
  alerts: DEFAULT_ALERTS,
};

const FormContext = createContext<FromContextType>(DEFAULT_CONTEXT_VALUE);

type BasePermissionEnforcement = $ReadOnly<{|
  ignore?: ?boolean,
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

type Props = $ReadOnly<{|
  children: React.Node,
  permissions: PermissionEnforcement,
|}>;

export function enforcePermissions(
  userPermissions: UserPermissions,
  permissionsEnforcement: PermissionEnforcement,
  ruleCheck: AlertRuleCheck,
) {
  if (permissionsEnforcement.adminRightsRequired === true) {
    const adminPermission = userPermissions.adminPolicy.access.isAllowed;
    ruleCheck({
      fieldId: 'System Rules',
      fieldDisplayName: 'Admin rights',
      value: permissionRuleValue2Bool(adminPermission),
      checkCallback: userIsAdmin =>
        userIsAdmin
          ? ''
          : `${fbt(
              'Admin rights are required. Contact your system administrator.',
              '',
            )}`,
    });

    return;
  }

  if (
    !permissionsEnforcement.entity ||
    !permissionsEnforcement.action ||
    permissionsEnforcement.ignore === true
  ) {
    return;
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
      return;
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
      return;
    }
    actionPermissionValue =
      userPermissions.inventoryPolicy[subjectEntity][enforcement.action];
  }

  if (actionPermissionValue == null) {
    return;
  }

  ruleCheck({
    fieldId: 'System Rules',
    fieldDisplayName: 'Permissions Check',
    value: permissionRuleValue2Bool(actionPermissionValue.isAllowed),
    checkCallback: passedPermissionsCheck =>
      passedPermissionsCheck
        ? ''
        : `${fbt(
            "User doesn't have sufficient permissions for this action",
            '',
          )}`,
  });
}

export function FormContextProvider(props: Props) {
  const {children, permissions} = props;
  const {me} = useMainContext();

  const permissionsEnforcementIsOn = useFeatureFlag(
    'permissions_ui_enforcement',
  );

  const userManagementDevModeIsOn = useFeatureFlag('user_management_dev');
  const shouldEnforcePermissions =
    permissionsEnforcementIsOn && permissions.ignore != true;

  return (
    <FormAlertsContextProvider>
      <FormAlertsContext.Consumer>
        {alerts => {
          if (shouldEnforcePermissions && me != null) {
            if (userManagementDevModeIsOn) {
              enforcePermissions(
                me.permissions,
                permissions,
                alerts.missingPermissions.check,
              );
            } else {
              alerts.missingPermissions.check({
                fieldId: 'System Rules',
                fieldDisplayName: 'Read Only User',
                value: me?.permissions.canWrite,
                checkCallback: canWrite =>
                  canWrite
                    ? ''
                    : `${fbt(
                        'Writing permissions are required. Contact your system administrator.',
                        '',
                      )}`,
              });
            }
          }
          return (
            <FormContext.Provider value={{alerts}}>
              {children}
            </FormContext.Provider>
          );
        }}
      </FormAlertsContext.Consumer>
    </FormAlertsContextProvider>
  );
}

export const useFormContext = () => useContext(FormContext);

export default FormContext;
