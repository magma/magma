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
  InventoryCatalogPolicy,
  InventoryPolicy,
} from '../utils/UserManagementUtils';

import * as React from 'react';
import PermissionsPolicyInventoryRulesSection from './PermissionsPolicyInventoryRulesSection';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    marginLeft: '4px',
    maxHeight: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    marginBottom: '4px',
    marginLeft: '4px',
    display: 'flex',
    flexDirection: 'column',
  },
}));

type RuleCatalogKey = string & $Keys<InventoryCatalogPolicy>;

const rules: Array<{key: RuleCatalogKey, title: React.Node}> = [
  {
    key: 'equipmentType',
    title: fbt('Equipment Types', ''),
  },
  {
    key: 'locationType',
    title: fbt('Location Types', ''),
  },
  {
    key: 'portType',
    title: fbt('Port Types', ''),
  },
  {
    key: 'serviceType',
    title: fbt('Service Types', ''),
  },
];

type Props = $ReadOnly<{|
  policy: ?InventoryPolicy,
  onChange: InventoryPolicy => void,
|}>;

export default function PermissionsPolicyInventoryCatalogRulesTab(
  props: Props,
) {
  const {policy, onChange} = props;
  const classes = useStyles();

  if (policy == null) {
    return null;
  }

  return (
    <div className={classes.root}>
      <div className={classes.header}>
        <Text variant="subtitle1">
          <fbt desc="">Inventory Catalog</fbt>
        </Text>
        <Text variant="body2" color="gray">
          <fbt desc="">
            Choose which sections of the catalog this group can modify.
          </fbt>
        </Text>
      </div>

      {rules.map(rule => (
        <PermissionsPolicyInventoryRulesSection
          mainCheckHeaderPrefix={rule.title}
          rule={policy[rule.key]}
          onChange={cud =>
            onChange({
              ...policy,
              [rule.key]: cud,
            })
          }
        />
      ))}
    </div>
  );
}
