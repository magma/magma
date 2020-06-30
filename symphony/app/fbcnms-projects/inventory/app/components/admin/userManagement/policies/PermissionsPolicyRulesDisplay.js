/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionsPolicy} from '../data/PermissionsPolicies';

import * as React from 'react';
import InventorySuspense from '../../../../common/InventorySuspense';
import PermissionsPolicyInventoryCatalogRulesTab from './PermissionsPolicyInventoryCatalogRulesTab';
import PermissionsPolicyInventoryDataRulesTab from './PermissionsPolicyInventoryDataRulesTab';
import PermissionsPolicyWorkforceDataRulesTab from './PermissionsPolicyWorkforceDataRulesTab';
import PermissionsPolicyWorkforceTemplatesRulesTab from './PermissionsPolicyWorkforceTemplatesRulesTab';
import classNames from 'classnames';
import {POLICY_TYPES} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  container: {
    display: 'flex',
    flexDirection: 'column',
  },
  rule: {
    marginBottom: '32px',
  },
}));

type Props = $ReadOnly<{|
  policy: PermissionsPolicy,
  className?: ?string,
|}>;

export default function PermissionsPolicyRulesDisplay(props: Props) {
  const {policy, className} = props;
  const classes = useStyles();

  const rules = useMemo(() => {
    switch (policy.type) {
      case POLICY_TYPES.InventoryPolicy.key:
        return (
          <>
            <PermissionsPolicyInventoryDataRulesTab
              className={classes.rule}
              policy={policy.inventoryRules}
            />
            <PermissionsPolicyInventoryCatalogRulesTab
              className={classes.rule}
              policy={policy.inventoryRules}
            />
          </>
        );
      case POLICY_TYPES.WorkforcePolicy.key:
        return (
          <>
            <PermissionsPolicyWorkforceDataRulesTab
              className={classes.rule}
              policy={policy.workforceRules}
            />
            <PermissionsPolicyWorkforceTemplatesRulesTab
              className={classes.rule}
              policy={policy.workforceRules}
            />
          </>
        );
      default:
        return null;
    }
  }, [classes.rule, policy.inventoryRules, policy.type, policy.workforceRules]);

  return (
    <InventorySuspense>
      <div className={classNames(classes.container, className)}>{rules}</div>
    </InventorySuspense>
  );
}
