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
  CUDPermissions,
  InventoryEntsPolicy,
  WorkforceCUDPermissions,
} from '../components/admin/userManagement/utils/UserManagementUtils';
import type {FormAlertsContextType} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import FormAlertsContext, {
  DEFAULT_CONTEXT_VALUE as DEFAULT_ALERTS,
  FormAlertsContextProvider,
} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import fbt from 'fbt';
import {createContext, useContext} from 'react';
import {useMainContext} from '../components/MainContext';

type FromContextType = $ReadOnly<{|
  alerts: FormAlertsContextType,
|}>;

const DEFAULT_CONTEXT_VALUE = {
  alerts: DEFAULT_ALERTS,
};

const FormContext = createContext<FromContextType>(DEFAULT_CONTEXT_VALUE);

type BasePermissionEnforcement = {
  ignore?: ?boolean,
};

type AdminPermissionEnforcement = {
  ...BasePermissionEnforcement,
  adminRightsRequired: true,
};

type InventoryEntName = $Keys<InventoryEntsPolicy> | 'service' | 'port';
type InventoryActionName = $Keys<CUDPermissions>;

type InventoryPermissionEnforcement = {
  ...BasePermissionEnforcement,
  entity: InventoryEntName,
  action?: ?InventoryActionName,
};

type WorkforceEntName =
  | 'workorder'
  | 'project'
  | 'workorderTemplate'
  | 'projectTemplate';

type WorkforceActionName = $Keys<WorkforceCUDPermissions>;

type WorkforcePermissionEnforcement = {
  ...BasePermissionEnforcement,
  entity: WorkforceEntName,
  action?: ?WorkforceActionName,
};

export type EntName = InventoryEntName | WorkforceEntName;
export type ActionName = InventoryActionName | WorkforceActionName;

export type PermissionEnforcement =
  | BasePermissionEnforcement
  | AdminPermissionEnforcement
  | InventoryPermissionEnforcement
  | WorkforcePermissionEnforcement;

type Props = {
  children: React.Node,
  permissions: PermissionEnforcement,
};

export function FormContextProvider(props: Props) {
  const {children, permissions} = props;
  const {me} = useMainContext();
  const {isFeatureEnabled} = useContext(AppContext);

  const permissionsEnforcementIsOn = isFeatureEnabled(
    'permissions_ui_enforcement',
  );
  const shouldEnforcePermissions =
    permissionsEnforcementIsOn && permissions.ignore != true;

  return (
    <FormAlertsContextProvider>
      <FormAlertsContext.Consumer>
        {alerts => {
          if (shouldEnforcePermissions) {
            alerts.editLock.check({
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
