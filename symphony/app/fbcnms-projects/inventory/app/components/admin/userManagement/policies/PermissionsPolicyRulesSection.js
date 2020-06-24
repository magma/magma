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
  BasicPermissionRule,
  CUDPermissions,
} from '../data/PermissionsPolicies';

import * as React from 'react';
import HierarchicalCheckbox, {
  HIERARCHICAL_RELATION,
} from '../utils/HierarchicalCheckbox';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import {
  bool2PermissionRuleValue,
  permissionRuleValue2Bool,
} from '../data/PermissionsPolicies';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';

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
  onChange: BasicPermissionRule => void,
  children?: React.Node,
|}>;

function InventoryDataRule(props: InventoryDataRuleProps) {
  const {title, rule, cudAction, disabled, onChange, children} = props;
  const classes = useStyles();

  if (rule == null) {
    return null;
  }

  return (
    <HierarchicalCheckbox
      id={cudAction}
      title={title}
      disabled={disabled}
      value={permissionRuleValue2Bool(rule[cudAction].isAllowed)}
      onChange={checked =>
        onChange({
          isAllowed: bool2PermissionRuleValue(checked),
        })
      }
      hierarchicalRelation={HIERARCHICAL_RELATION.PARENT_REQUIRED}
      className={classes.rule}>
      {children}
    </HierarchicalCheckbox>
  );
}

type DataRuleTitleProps = $ReadOnly<{|
  children: React.Node,
|}>;

export function DataRuleTitle(props: DataRuleTitleProps) {
  const {children} = props;
  return (
    <Text variant="subtitle2" color="inherit">
      {children}
    </Text>
  );
}

type Props = $ReadOnly<{|
  title?: React.Node,
  subtitle?: React.Node,
  mainCheckHeaderPrefix?: React.Node,
  rule: CUDPermissions,
  disabled?: ?boolean,
  className?: ?string,
  onChange?: CUDPermissions => void,
  children?: React.Node,
|}>;

export default function PermissionsPolicyRulesSection(props: Props) {
  const {
    title,
    subtitle,
    mainCheckHeaderPrefix,
    rule: ruleProp,
    disabled,
    className,
    onChange,
    children,
  } = props;
  const classes = useStyles();
  const [rule, setRule] = useState<CUDPermissions>(ruleProp);
  useEffect(() => setRule(ruleProp), [ruleProp]);

  const updateRuleChange = useCallback(
    (cudAction: string & CUDPermissionsKey) => (
      actionValue: BasicPermissionRule,
    ) => {
      setRule(currentRule => {
        const newRuleValue: CUDPermissions = {
          ...currentRule,
          [cudAction]: actionValue,
        };
        if (onChange != null) {
          onChange(newRuleValue);
        }
        return newRuleValue;
      });
    },
    [onChange],
  );

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
            <DataRuleTitle>
              <fbt desc="">Edit</fbt>
            </DataRuleTitle>
            {mainCheckHeaderPrefix != null && (
              <span> {mainCheckHeaderPrefix}</span>
            )}
          </>
        }
        rule={rule}
        cudAction="update"
        disabled={disabled == true}
        onChange={updateRuleChange('update')}>
        {dependantDataRules.map(dRule => (
          <InventoryDataRule
            title={<DataRuleTitle>{dRule.title}</DataRuleTitle>}
            rule={rule}
            cudAction={dRule.key}
            disabled={disabled == true}
            onChange={updateRuleChange(dRule.key)}
          />
        ))}
        {children}
      </InventoryDataRule>
    </div>
  );
}
