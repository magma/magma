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
import type {TabProps} from '@fbcnms/ui/components/design-system/Tabs/TabsBar';

import * as React from 'react';
import PermissionsPolicyInventoryCatalogRulesTab from './PermissionsPolicyInventoryCatalogRulesTab';
import PermissionsPolicyInventoryDataRulesTab from './PermissionsPolicyInventoryDataRulesTab';
import PermissionsPolicyWorkforceDataRulesTab from './PermissionsPolicyWorkforceDataRulesTab';
import PermissionsPolicyWorkforceTemplatesRulesTab from './PermissionsPolicyWorkforceTemplatesRulesTab';
import TabsBar from '@fbcnms/ui/components/design-system/Tabs/TabsBar';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {POLICY_TYPES} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useMemo, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    maxHeight: '100%',
    overflow: 'hidden',
    display: 'flex',
    flexDirection: 'column',
  },
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
                policy={policy}
                onChange={onChange}
              />
            ),
          },
          {
            tab: {
              label: `${fbt('Inventory Catalog', '')}`,
            },
            view: (
              <PermissionsPolicyInventoryCatalogRulesTab
                policy={policy}
                onChange={onChange}
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
                policy={policy}
                onChange={onChange}
              />
            ),
          },
          {
            tab: {
              label: `${fbt('Workforce Templates', '')}`,
            },
            view: (
              <PermissionsPolicyWorkforceTemplatesRulesTab
                policy={policy}
                onChange={onChange}
              />
            ),
          },
        ];
      default:
        return [];
    }
  }, [onChange, policy]);

  return (
    <div className={classNames(classes.root, className)}>
      <TabsBar
        className={classes.tabsContainer}
        tabs={ruleTypes.map(type => type.tab)}
        activeTabIndex={activeTab}
        onChange={setActiveTab}
        spread={false}
      />
      <div className={classes.viewContainer}>{ruleTypes[activeTab].view}</div>
    </div>
  );
}
