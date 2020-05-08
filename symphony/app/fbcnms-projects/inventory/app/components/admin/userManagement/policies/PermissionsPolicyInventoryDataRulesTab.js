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
  InventoryPolicy,
} from '../utils/UserManagementUtils';

import * as React from 'react';
import Checkbox from '@fbcnms/ui/components/design-system/Checkbox/Checkbox';
import Switch from '@fbcnms/ui/components/design-system/switch/Switch';
import Text from '@fbcnms/ui/components/design-system/Text';
import fbt from 'fbt';
import {
  bool2PermissionRuleValue,
  permissionRuleValue2Bool,
} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    marginLeft: '4px',
    maxHeight: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  section: {
    display: 'flex',
    flexDirection: 'column',
    marginTop: '32px',
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
  readRule: {
    marginLeft: '4px',
  },
}));

type CUDPermissionsKey = $Keys<CUDPermissions>;

type InventoryDataRuleProps = $ReadOnly<{|
  title: React.Node,
  rule: CUDPermissions,
  cudAction: string & CUDPermissionsKey,
  disabled: boolean,
  onChange: CUDPermissions => void,
|}>;

function InventoryDataRule(props: InventoryDataRuleProps) {
  const {title, rule, cudAction, disabled, onChange} = props;
  const classes = useStyles();

  if (rule == null) {
    return null;
  }

  return (
    <Checkbox
      className={classes.rule}
      title={title}
      disabled={disabled}
      checked={!disabled && permissionRuleValue2Bool(rule[cudAction].isAllowed)}
      onChange={selection =>
        onChange({
          ...rule,
          [cudAction]: {
            isAllowed: bool2PermissionRuleValue(selection === 'checked'),
          },
        })
      }
    />
  );
}

type InventoryDataRulesSectionProps = $ReadOnly<{|
  title: React.Node,
  subtitle: React.Node,
  rule: ?CUDPermissions,
  disabled: boolean,
  onChange: CUDPermissions => void,
|}>;

function InventoryDataRulesSection(props: InventoryDataRulesSectionProps) {
  const {title, subtitle, rule, disabled, onChange} = props;
  const classes = useStyles();

  if (rule == null) {
    return null;
  }

  const dataRules: Array<{key: CUDPermissionsKey, title: React.Node}> = [
    {
      key: 'create',
      title: fbt('Add', ''),
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
    <div className={classes.section}>
      <div className={classes.header}>
        <Text variant="subtitle1">{title}</Text>
        <Text variant="body2" color="gray">
          {subtitle}
        </Text>
      </div>
      {dataRules.map(dRule => (
        <InventoryDataRule
          title={dRule.title}
          rule={rule}
          cudAction={dRule.key}
          disabled={disabled}
          onChange={onChange}
        />
      ))}
    </div>
  );
}

type Props = $ReadOnly<{|
  policy: ?InventoryPolicy,
  onChange: InventoryPolicy => void,
|}>;

export default function PermissionsPolicyInventoryDataRulesTab(props: Props) {
  const {policy, onChange} = props;
  const classes = useStyles();

  if (policy == null) {
    return null;
  }

  const readAllowed = permissionRuleValue2Bool(policy.read.isAllowed);

  return (
    <div className={classes.root}>
      <Switch
        className={classes.readRule}
        title={fbt('View inventory data', '')}
        checked={readAllowed}
        onChange={checked =>
          onChange({
            ...policy,
            read: {
              isAllowed: bool2PermissionRuleValue(checked),
            },
          })
        }
      />
      <InventoryDataRulesSection
        title={fbt('Locations', '')}
        subtitle={fbt(
          'Location data includes location details, properties, floor plans and coverage maps.',
          '',
        )}
        disabled={!readAllowed}
        rule={policy.location}
        onChange={location =>
          onChange({
            ...policy,
            location,
          })
        }
      />
      <InventoryDataRulesSection
        title={fbt('Equipment', '')}
        subtitle={fbt(
          'Equipment data includes equipment items, ports, links, services and network maps.',
          '',
        )}
        disabled={!readAllowed}
        rule={policy.equipment}
        onChange={equipment =>
          onChange({
            ...policy,
            equipment,
          })
        }
      />
    </div>
  );
}
