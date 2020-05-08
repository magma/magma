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
  InventoryCatalogPolicy,
  InventoryPolicy,
} from '../utils/UserManagementUtils';

import * as React from 'react';
import HierarchicalCheckbox from '../utils/HierarchicalCheckbox';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {
  bool2PermissionRuleValue,
  permissionRuleValue2Bool,
} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';

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
  rule: {
    marginTop: '8px',
    marginLeft: '4px',
  },
}));

type RuleCUDKey = string & $Keys<CUDPermissions>;
type RuleCatalogKey = string & $Keys<InventoryCatalogPolicy>;

type CatalogRuleProps = $ReadOnly<{|
  title: React.Node,
  policy: InventoryPolicy,
  ruleCUD: RuleCUDKey,
  ruleCatalog: RuleCatalogKey,
  onChange: (
    ruleCUD: RuleCUDKey,
    ruleCatalog: RuleCatalogKey,
    value: ?boolean,
  ) => void,
|}>;

function CatalogRule(props: CatalogRuleProps) {
  const {title, policy, ruleCUD, ruleCatalog, onChange} = props;
  const classes = useStyles();

  return (
    <HierarchicalCheckbox
      id={`${ruleCatalog}_${ruleCUD}`}
      title={title}
      className={classes.rule}
      value={permissionRuleValue2Bool(policy[ruleCatalog][ruleCUD].isAllowed)}
      onChange={checked => onChange(ruleCUD, ruleCatalog, checked)}
    />
  );
}

type CatalogsTreeProps = $ReadOnly<{|
  title: React.Node,
  policy: InventoryPolicy,
  ruleCUD: string & $Keys<CUDPermissions>,
  onChange: InventoryPolicy => void,
|}>;

function CatalogsTree(props: CatalogsTreeProps) {
  const {title, policy: propPolicy, ruleCUD, onChange: propOnChange} = props;

  const [policy, setPolicy] = useState(propPolicy);
  useEffect(() => {
    setPolicy(propPolicy);
  }, [propPolicy]);

  const classes = useStyles();

  const onChange = useCallback(
    (ruleCUD: RuleCUDKey, ruleCatalog: RuleCatalogKey, value: ?boolean) => {
      setPolicy(currentPolicy => {
        const newPolicy = {
          ...currentPolicy,
          [ruleCatalog]: {
            ...currentPolicy[ruleCatalog],
            [ruleCUD]: {
              isAllowed: bool2PermissionRuleValue(value),
            },
          },
        };
        propOnChange(newPolicy);
        return newPolicy;
      });
    },
    [propOnChange],
  );

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

  return (
    <HierarchicalCheckbox id={ruleCUD} className={classes.rule} title={title}>
      {rules.map(rule => (
        <CatalogRule
          title={rule.title}
          policy={policy}
          ruleCatalog={rule.key}
          ruleCUD={ruleCUD}
          onChange={onChange}
        />
      ))}
    </HierarchicalCheckbox>
  );
}

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

  const trees: Array<{key: RuleCUDKey, title: React.Node}> = [
    {
      key: 'create',
      title: fbt('Create', ''),
    },
    {
      key: 'update',
      title: fbt('Edit', ''),
    },
    {
      key: 'delete',
      title: fbt('Delete', ''),
    },
  ];

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
      {trees.map(tree => (
        <CatalogsTree
          policy={policy}
          ruleCUD={tree.key}
          title={tree.title}
          onChange={onChange}
        />
      ))}
    </div>
  );
}
