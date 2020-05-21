/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FormAlertsContextType} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import type {PermissionEnforcement} from '../components/admin/userManagement/utils/usePermissions';

import * as React from 'react';
import fbt from 'fbt';
import useFeatureFlag from '@fbcnms/ui/context/useFeatureFlag';
import usePermissions from '../components/admin/userManagement/utils/usePermissions';
import {
  DEFAULT_CONTEXT_VALUE as DEFAULT_ALERTS,
  FormAlertsContextProvider,
  useFormAlertsContext,
} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';
import {createContext, useContext} from 'react';
import {useMainContext} from '../components/MainContext';

type FromContextType = $ReadOnly<{|
  alerts: FormAlertsContextType,
|}>;

const DEFAULT_CONTEXT_VALUE = {
  alerts: DEFAULT_ALERTS,
};

const FormContext = createContext<FromContextType>(DEFAULT_CONTEXT_VALUE);

type Props = $ReadOnly<{|
  children: React.Node,
  permissions: PermissionEnforcement,
|}>;

function FormWrapper(props: Props) {
  const {children, permissions} = props;
  const {me} = useMainContext();

  const permissionsEnforcementIsOn = useFeatureFlag(
    'permissions_ui_enforcement',
  );

  const permissionPoliciesMode = useFeatureFlag('permission_policies');
  const shouldEnforcePermissions =
    permissionsEnforcementIsOn && permissions.ignorePermissions != true;

  const permissionsRules = usePermissions();
  const alerts = useFormAlertsContext();

  if (permissionPoliciesMode) {
    permissionsRules.check(permissions, 'Form Permissions');
  } else if (shouldEnforcePermissions && me != null) {
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

  return (
    <FormContext.Provider value={{alerts}}>{children}</FormContext.Provider>
  );
}

export function FormContextProvider(props: Props) {
  return (
    <FormAlertsContextProvider>
      <FormWrapper {...props} />
    </FormAlertsContextProvider>
  );
}

export const useFormContext = () => useContext(FormContext);

export default FormContext;
