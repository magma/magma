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
import type {TabProps} from '@fbcnms/ui/components/design-system/Tabs/TabsBar';

import * as React from 'react';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import PermissionsPolicyInventoryCatalogRulesTab from './PermissionsPolicyInventoryCatalogRulesTab';
import PermissionsPolicyInventoryDataRulesTab from './PermissionsPolicyInventoryDataRulesTab';
import PermissionsPolicyWorkforceDataRulesTab from './PermissionsPolicyWorkforceDataRulesTab';
import PermissionsPolicyWorkforceTemplatesRulesTab from './PermissionsPolicyWorkforceTemplatesRulesTab';
import TabsBar from '@fbcnms/ui/components/design-system/Tabs/TabsBar';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {POLICY_TYPES} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo, useState} from 'react';

const useStyles = makeStyles(() => ({
  tabsContainer: {
    paddingLeft: '24px',
  },
  viewContainer: {
    padding: '28px',
    overflowY: 'auto',
    background: symphony.palette.white,
  },
}));

type Props = $ReadOnly<{|
  policy: PermissionsPolicy,
  onChange: PermissionsPolicy => void,
  className?: ?string,
|}>;

type ViewTab = $ReadOnly<{|
  tab: TabProps,
  view: React.Node,
|}>;

export default function PermissionsPolicyRulesPane(props: Props) {
  const {policy, onChange, className} = props;
  const classes = useStyles();
  const [activeTab, setActiveTab] = useState(0);
  const callOnInventoryChange = useCallback(
    inventoryRules =>
      onChange({
        ...policy,
        inventoryRules,
      }),
    [onChange, policy],
  );
  const callOnWorkforceChange = useCallback(
    workforceRules =>
      onChange({
        ...policy,
        workforceRules,
      }),
    [onChange, policy],
  );

  const ruleTypes: Array<ViewTab> = useMemo(() => {
    switch (policy.type) {
      case POLICY_TYPES.InventoryPolicy.key:
        return [
          {
            tab: {
              label: `${fbt('Inventory Data', '')}`,
            },
            view: (
              <PermissionsPolicyInventoryDataRulesTab
                policy={policy.inventoryRules}
                onChange={callOnInventoryChange}
              />
            ),
          },
          {
            tab: {
              label: `${fbt('Inventory Catalog', '')}`,
            },
            view: (
              <PermissionsPolicyInventoryCatalogRulesTab
                policy={policy.inventoryRules}
                onChange={callOnInventoryChange}
              />
            ),
          },
        ];
      case POLICY_TYPES.WorkforcePolicy.key:
        return [
          {
            tab: {
              label: `${fbt('Workforce Data', '')}`,
            },
            view: (
              <PermissionsPolicyWorkforceDataRulesTab
                policy={policy.workforceRules}
                onChange={callOnWorkforceChange}
              />
            ),
          },
          {
            tab: {
              label: `${fbt('Workforce Templates', '')}`,
            },
            view: (
              <PermissionsPolicyWorkforceTemplatesRulesTab
                policy={policy.workforceRules}
                onChange={callOnWorkforceChange}
              />
            ),
          },
        ];
      default:
        return [];
    }
  }, [
    callOnInventoryChange,
    callOnWorkforceChange,
    policy.inventoryRules,
    policy.type,
    policy.workforceRules,
  ]);

  return (
    <Card className={className} margins="none">
      <TabsBar
        className={classes.tabsContainer}
        tabs={ruleTypes.map(type => type.tab)}
        activeTabIndex={activeTab}
        onChange={setActiveTab}
        spread={false}
      />
      <div className={classes.viewContainer}>{ruleTypes[activeTab].view}</div>
    </Card>
  );
}
