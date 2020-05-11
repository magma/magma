/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {CUDPermissions} from '../utils/UserManagementUtils';

import * as React from 'react';
import Checkbox from '@fbcnms/ui/components/design-system/Checkbox/Checkbox';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import {
  bool2PermissionRuleValue,
  permissionRuleValue2Bool,
} from '../utils/UserManagementUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
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
  dependantRules: {
    marginLeft: '34px',
    display: 'flex',
    flexDirection: 'column',
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

type Props = $ReadOnly<{|
  title?: React.Node,
  subtitle?: React.Node,
  mainCheckHeaderPrefix?: React.Node,
  rule: ?CUDPermissions,
  disabled?: ?boolean,
  className?: ?string,
  onChange: CUDPermissions => void,
|}>;

export default function PermissionsPolicyInventoryRulesSection(props: Props) {
  const {
    title,
    subtitle,
    mainCheckHeaderPrefix,
    rule,
    disabled,
    className,
    onChange,
  } = props;
  const classes = useStyles();

  if (rule == null) {
    return null;
  }

  const dependantDataRules: Array<{
    key: CUDPermissionsKey,
    title: React.Node,
  }> = [
    {
      key: 'create',
      title: fbt('Add', ''),
    },
    {
      key: 'delete',
      title: fbt('Delete', ''),
    },
  ];

  return (
    <div className={classNames(classes.section, className)}>
      <div className={classes.header}>
        <Text variant="subtitle1">{title}</Text>
        <Text variant="body2" color="gray">
          {subtitle}
        </Text>
      </div>
      <InventoryDataRule
        title={
          <>
            <Text variant="subtitle2">
              <fbt desc="">Edit</fbt>
            </Text>
            {mainCheckHeaderPrefix != null && (
              <span> {mainCheckHeaderPrefix}</span>
            )}
          </>
        }
        rule={rule}
        cudAction="update"
        disabled={disabled == true}
        onChange={onChange}
      />
      <div className={classes.dependantRules}>
        {dependantDataRules.map(dRule => (
          <InventoryDataRule
            title={<Text variant="subtitle2">{dRule.title}</Text>}
            rule={rule}
            cudAction={dRule.key}
            disabled={
              disabled == true ||
              !permissionRuleValue2Bool(rule.update.isAllowed)
            }
            onChange={onChange}
          />
        ))}
      </div>
    </div>
  );
}
