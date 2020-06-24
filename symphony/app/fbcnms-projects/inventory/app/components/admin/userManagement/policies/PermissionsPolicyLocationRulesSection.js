/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {LocationCUDPermissions} from '../data/PermissionsPolicies';
import type {PermissionsPolicyRulesSectionDisplayProps} from './PermissionsPolicyRulesSection';

import * as React from 'react';
import PermissionsPolicyRulesSection from './PermissionsPolicyRulesSection';

type Props = $ReadOnly<{|
  ...PermissionsPolicyRulesSectionDisplayProps,
  rule: LocationCUDPermissions,
  onChange?: LocationCUDPermissions => void,
|}>;

export default function PermissionsPolicyLocationRulesSection(props: Props) {
  const {rule, onChange, ...permissionsPolicyRulesSectionDisplayProps} = props;

  return (
    <PermissionsPolicyRulesSection
      // eslint-disable-next-line no-warning-comments
      // $FlowFixMe: NEXT DIFF - Handle location types
      rule={rule}
      // eslint-disable-next-line no-warning-comments
      // $FlowFixMe: NEXT DIFF - Handle location types
      onChange={onChange}
      {...permissionsPolicyRulesSectionDisplayProps}
    />
  );
}
