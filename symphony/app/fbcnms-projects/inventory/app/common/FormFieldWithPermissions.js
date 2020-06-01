/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FormFieldProps} from '@fbcnms/ui/components/design-system/FormField/FormField';
import type {PermissionEnforcement} from '../components/admin/userManagement/utils/usePermissions';

import * as React from 'react';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import usePermissions from '../components/admin/userManagement/utils/usePermissions';

type Props = $ReadOnly<{|
  ...FormFieldProps,
  permissions: PermissionEnforcement,
  includingParentFormPermissions?: boolean,
|}>;

export default function FormFieldWithPermissions(props: Props) {
  const {
    includingParentFormPermissions,
    permissions,
    disabled,
    ...rest
  } = props;

  const permissionsRules = usePermissions();

  const missingPermissionsMessage = permissionsRules.check(permissions);

  return (
    <FormField
      ignorePermissions={!includingParentFormPermissions}
      disabled={disabled === true || !!missingPermissionsMessage}
      {...rest}
    />
  );
}
