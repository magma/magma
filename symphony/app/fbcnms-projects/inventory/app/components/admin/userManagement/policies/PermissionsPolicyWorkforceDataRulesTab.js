/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {WorkforceCUD, WorkforcePolicy} from '../utils/UserManagementUtils';

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

type InventoryDataRulesSectionProps = $ReadOnly<{|
  title: React.Node,
  subtitle: React.Node,
  rule: ?WorkforceCUD,
  disabled: boolean,
  onChange: WorkforceCUD => void,
|}>;

function WorkforceDataRulesSection(props: InventoryDataRulesSectionProps) {
  const {title, subtitle, rule, disabled, onChange} = props;
  const classes = useStyles();

  if (rule == null) {
    return null;
  }

  return (
    <div className={classes.section}>
      <div className={classes.header}>
        <Text variant="subtitle1">{title}</Text>
        <Text variant="body2">{subtitle}</Text>
      </div>
      <Checkbox
        className={classes.rule}
        title={fbt('Create', '')}
        disabled={disabled}
        checked={!disabled && permissionRuleValue2Bool(rule.create.isAllowed)}
        onChange={selection =>
          onChange({
            ...rule,
            create: {
              isAllowed: bool2PermissionRuleValue(selection === 'checked'),
            },
          })
        }
      />
      <Checkbox
        className={classes.rule}
        title={fbt('Edit', '')}
        disabled={disabled}
        checked={!disabled && permissionRuleValue2Bool(rule.update.isAllowed)}
        onChange={selection =>
          onChange({
            ...rule,
            update: {
              isAllowed: bool2PermissionRuleValue(selection === 'checked'),
            },
          })
        }
      />
      <Checkbox
        className={classes.rule}
        title={fbt('Assign', '')}
        disabled={disabled}
        checked={!disabled && permissionRuleValue2Bool(rule.assign.isAllowed)}
        onChange={selection =>
          onChange({
            ...rule,
            assign: {
              isAllowed: bool2PermissionRuleValue(selection === 'checked'),
            },
          })
        }
      />
      <Checkbox
        className={classes.rule}
        title={fbt('Delete', '')}
        disabled={disabled}
        checked={!disabled && permissionRuleValue2Bool(rule.delete.isAllowed)}
        onChange={selection =>
          onChange({
            ...rule,
            delete: {
              isAllowed: bool2PermissionRuleValue(selection === 'checked'),
            },
          })
        }
      />
      <Checkbox
        className={classes.rule}
        title={fbt('Transfer Ownership', '')}
        disabled={disabled}
        checked={
          !disabled &&
          permissionRuleValue2Bool(rule.transferOwnership.isAllowed)
        }
        onChange={selection =>
          onChange({
            ...rule,
            transferOwnership: {
              isAllowed: bool2PermissionRuleValue(selection === 'checked'),
            },
          })
        }
      />
    </div>
  );
}

type Props = $ReadOnly<{|
  policy: ?WorkforcePolicy,
  onChange: WorkforcePolicy => void,
|}>;

export default function PermissionsPolicyWorkforceDataRulesTab(props: Props) {
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
        title={fbt('View work orders and projects', '')}
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
      <WorkforceDataRulesSection
        title={fbt('Modifying', '')}
        subtitle={fbt(
          'Choose the permissions this policy should include for modifying the selected work orders and projects.',
          '',
        )}
        disabled={!readAllowed}
        rule={policy.data}
        onChange={data =>
          onChange({
            ...policy,
            data,
          })
        }
      />
    </div>
  );
}
