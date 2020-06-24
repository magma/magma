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
  LocationCUDPermissions,
} from '../data/PermissionsPolicies';
import type {PermissionsPolicyRulesSectionDisplayProps} from './PermissionsPolicyRulesSection';

import * as React from 'react';
import PermissionsPolicyLocationRulesSpecification from './PermissionsPolicyLocationRulesSpecification';
import PermissionsPolicyRulesSection from './PermissionsPolicyRulesSection';
import symphony from '@fbcnms/ui/theme/symphony';
import useFeatureFlag from '@fbcnms/ui/context/useFeatureFlag';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo} from 'react';

const useStyles = makeStyles(() => ({
  secondLevelBox: {
    backgroundColor: symphony.palette.background,
    borderStyle: 'solid',
    borderWidth: '1px',
    borderColor: symphony.palette.D100,
    borderLeftWidth: '2px',
    borderLeftColor: symphony.palette.primary,
    paddingTop: '16px',
    paddingBottom: '10px',
    borderRadius: '2px',
    marginTop: '8px',
  },
}));

type Props = $ReadOnly<{|
  ...PermissionsPolicyRulesSectionDisplayProps,
  locationRule: LocationCUDPermissions,
  onChange?: LocationCUDPermissions => void,
|}>;

export default function PermissionsPolicyLocationRulesSection(props: Props) {
  const {
    locationRule,
    onChange,
    disabled,
    ...permissionsPolicyRulesSectionDisplayProps
  } = props;
  const classes = useStyles();

  const rule: CUDPermissions = useMemo(
    () => ({
      create: locationRule.create,
      delete: locationRule.delete,
      update: {
        isAllowed: locationRule.update.isAllowed,
      },
    }),
    [locationRule.create, locationRule.delete, locationRule.update.isAllowed],
  );

  const callOnChange = useCallback(
    (updatedRule: LocationCUDPermissions) => {
      if (onChange == null) {
        return;
      }
      onChange(updatedRule);
    },
    [onChange],
  );

  const callOnChangeWithCUDPermissions = useCallback(
    (updatedRule: CUDPermissions) => {
      callOnChange({
        ...locationRule,
        ...updatedRule,
        update: {
          ...locationRule.update,
          ...updatedRule.update,
        },
      });
    },
    [locationRule, callOnChange],
  );

  const isPermissionPolicyPerTypeEnabled = useFeatureFlag(
    'permission_policy_per_type',
  );

  return (
    <PermissionsPolicyRulesSection
      rule={rule}
      onChange={callOnChangeWithCUDPermissions}
      secondLevelRulesClassName={
        isPermissionPolicyPerTypeEnabled ? classes.secondLevelBox : null
      }
      policySpecifications={
        <PermissionsPolicyLocationRulesSpecification
          locationRule={locationRule}
          onChange={callOnChange}
          disabled={disabled}
        />
      }
      disabled={disabled}
      {...permissionsPolicyRulesSectionDisplayProps}
    />
  );
}
