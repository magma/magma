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
} from '../data/PermissionsPolicies';

import * as React from 'react';
import PermissionsPolicyRulesSection from './PermissionsPolicyRulesSection';
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
      <PermissionsPolicyRulesSection
        rule={{
          create: rule.create,
          delete: rule.delete,
          update: rule.update,
        }}
        className={classes.section}
        onChange={ruleCUD =>
          onChange({
            ...rule,
            ...ruleCUD,
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
          'Choose the permissions this policy should include for modifying work orders and projects templates.',
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
