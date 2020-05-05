/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionsPolicy} from '../utils/UserManagementUtils';

type Props = $ReadOnly<{|
  policy: PermissionsPolicy,
  onChange: PermissionsPolicy => void,
|}>;

export default function PermissionsPolicyWorkforceDataRulesTab(_props: Props) {
  return 'PermissionsPolicyWorkforceDataRulesTab';
}
