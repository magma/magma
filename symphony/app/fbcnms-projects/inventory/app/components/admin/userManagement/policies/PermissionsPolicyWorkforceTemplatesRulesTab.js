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
  WorkforcePolicy,
} from '../utils/UserManagementUtils';

import * as React from 'react';
import Checkbox from '@fbcnms/ui/components/design-system/Checkbox/Checkbox';
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
  rule: ?CUDPermissions,
  onChange: CUDPermissions => void,
|}>;

function WorkforceTemplatesRulesSection(props: InventoryDataRulesSectionProps) {
  const {title, subtitle, rule, onChange} = props;
  const classes = useStyles();

  if (rule == null) {
    return null;
  }

  return (
    <div className={classes.section}>
      <div className={classes.header}>
        <Text variant="subtitle1">{title}</Text>
        <Text variant="body2" color="gray">
          {subtitle}
        </Text>
      </div>
      <Checkbox
        className={classes.rule}
        title={fbt('Create', '')}
        checked={permissionRuleValue2Bool(rule.create.isAllowed)}
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
        checked={permissionRuleValue2Bool(rule.update.isAllowed)}
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
        title={fbt('Delete', '')}
        checked={permissionRuleValue2Bool(rule.delete.isAllowed)}
        onChange={selection =>
          onChange({
            ...rule,
            delete: {
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

export default function PermissionsPolicyWorkforceTemplatesRulesTab(
  props: Props,
) {
  const {policy, onChange} = props;
  const classes = useStyles();

  if (policy == null) {
    return null;
  }

  return (
    <div className={classes.root}>
      <WorkforceTemplatesRulesSection
        title={fbt('Workforce Templates', '')}
        subtitle={fbt(
          'Choose the permissions this policy should include for modifying the selected work orders and projects.',
          '',
        )}
        rule={policy.templates}
        onChange={templates =>
          onChange({
            ...policy,
            templates,
          })
        }
      />
    </div>
  );
}
