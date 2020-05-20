/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FormActionProps} from '@fbcnms/ui/components/design-system/Form/FormAction';
import type {PermissionEnforcement} from '../components/admin/userManagement/utils/usePermissions';

import * as React from 'react';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import usePermissions from '../components/admin/userManagement/utils/usePermissions';

type Props = $ReadOnly<{|
  ...FormActionProps,
  permissions: PermissionEnforcement,
  includingParentFormPermissions?: boolean,
|}>;

export default function FormActionWithPermissions(props: Props) {
  const {includingParentFormPermissions, permissions, ...rest} = props;

  const permissionsRules = usePermissions();

  const missingPermissionsMessage = permissionsRules.check(
    permissions,
    includingParentFormPermissions === true
      ? `${JSON.stringify(permissions)}`
      : null,
  );

  if (missingPermissionsMessage) {
    return null;
  }

  return (
    <FormAction ignorePermissions={!includingParentFormPermissions} {...rest} />
  );
}
