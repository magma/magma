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
  WorkforceCUDPermissions,
  WorkforcePolicy,
} from '../data/PermissionsPolicies';

import * as React from 'react';
import HierarchicalCheckbox, {
  HIERARCHICAL_RELATION,
} from '../utils/HierarchicalCheckbox';
import PermissionsPolicyRulesSection, {
  DataRuleTitle,
} from './PermissionsPolicyRulesSection';
import PermissionsPolicyWorkforceDataRulesSpecification from './PermissionsPolicyWorkforceDataRulesSpecification';
import Switch from '@fbcnms/ui/components/design-system/switch/Switch';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {
  bool2PermissionRuleValue,
  permissionRuleValue2Bool,
} from '../data/PermissionsPolicies';
import {debounce} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';

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
    marginTop: '8px',
  },
  header: {
    marginBottom: '16px',
    marginLeft: '4px',
    display: 'flex',
    flexDirection: 'column',
  },
  policySpecificationContainer: {
    display: 'flex',
    flexDirection: 'column',
    padding: '16px',
    paddingBottom: '8px',
    backgroundColor: symphony.palette.background,
    borderStyle: 'solid',
    borderWidth: '1px',
    borderColor: symphony.palette.D100,
    borderLeftWidth: '2px',
    borderLeftColor: symphony.palette.primary,
    borderRadius: '2px',
    marginTop: '12px',
  },
  methodSelectionBox: {
    display: 'flex',
    flexDirection: 'column',
    width: 'fit-content',
    marginBottom: '16px',
    '& > *': {
      marginBottom: '4px',
    },
  },
  policyMethodSelect: {
    '&&': {
      paddingLeft: '8px',
      marginRight: '16px',
    },
  },
  templatesFieldContainer: {
    minHeight: '78px',
  },
  hidden: {
    display: 'none',
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
  rule: WorkforceCUDPermissions,
  disabled: boolean,
  onChange?: WorkforceCUDPermissions => void,
|}>;

type WorkforceCUDPermissionsKey = $Keys<WorkforceCUDPermissions>;

function WorkforceDataRulesSection(props: InventoryDataRulesSectionProps) {
  const {rule: ruleProp, disabled, onChange} = props;
  const classes = useStyles();

  const [rule, setRule] = useState<WorkforceCUDPermissions>(ruleProp);
  useEffect(() => setRule(ruleProp), [ruleProp]);

  const debouncedOnChange = useCallback(
    debounce(
      (newRuleValue: WorkforceCUDPermissions) =>
        onChange && onChange(newRuleValue),
      100,
    ),
    [],
  );

  const updateRuleChange = useCallback(
    (
      updates: Array<{|
        cudAction: string & WorkforceCUDPermissionsKey,
        actionValue: BasicPermissionRule,
      |}>,
    ) => {
      setRule(currentRule => {
        const newRuleValue: WorkforceCUDPermissions = updates.reduce(
          (ruleSoFar, update) => ({
            ...ruleSoFar,
            [update.cudAction]: update.actionValue,
          }),
          currentRule,
        );
        debouncedOnChange(newRuleValue);
        return newRuleValue;
      });
    },
    [debouncedOnChange],
  );

  return (
    <div className={classes.section}>
      <PermissionsPolicyRulesSection
        disabled={disabled}
        rule={{
          create: rule.create,
          delete: rule.delete,
          update: rule.update,
        }}
        className={classes.section}
        onChange={ruleCUD =>
          updateRuleChange([
            {
              cudAction: 'create',
              actionValue: ruleCUD.create,
            },
            {
              cudAction: 'update',
              actionValue: ruleCUD.update,
            },
            {
              cudAction: 'delete',
              actionValue: ruleCUD.delete,
            },
          ])
        }>
        <HierarchicalCheckbox
          id="assign"
          title={
            <DataRuleTitle>
              <fbt desc="">Assign</fbt>
            </DataRuleTitle>
          }
          disabled={disabled}
          value={permissionRuleValue2Bool(rule.assign.isAllowed)}
          onChange={
            onChange != null
              ? checked =>
                  onChange({
                    ...rule,
                    assign: {
                      isAllowed: bool2PermissionRuleValue(checked),
                    },
                  })
              : undefined
          }
          hierarchicalRelation={HIERARCHICAL_RELATION.PARENT_REQUIRED}
          className={classes.rule}
        />
        <HierarchicalCheckbox
          id="transferOwnership"
          title={
            <DataRuleTitle>
              <fbt desc="">Transfer ownership</fbt>
            </DataRuleTitle>
          }
          disabled={disabled}
          value={
            !disabled &&
            permissionRuleValue2Bool(rule.transferOwnership.isAllowed)
          }
          onChange={checked =>
            updateRuleChange([
              {
                cudAction: 'transferOwnership',
                actionValue: {
                  isAllowed: bool2PermissionRuleValue(checked),
                },
              },
            ])
          }
          hierarchicalRelation={HIERARCHICAL_RELATION.PARENT_REQUIRED}
          className={classes.rule}
        />
      </PermissionsPolicyRulesSection>
    </div>
  );
}

type Props = $ReadOnly<{|
  policy: ?WorkforcePolicy,
  onChange?: WorkforcePolicy => void,
  className?: ?string,
|}>;

export default function PermissionsPolicyWorkforceDataRulesTab(props: Props) {
  const {policy, onChange, className} = props;
  const classes = useStyles();

  const callOnChange = useCallback(
    (updatedPolicy: WorkforcePolicy) => {
      if (onChange == null) {
        return;
      }
      onChange(updatedPolicy);
    },
    [onChange],
  );

  if (policy == null) {
    return null;
  }

  const readAllowed = permissionRuleValue2Bool(policy.read.isAllowed);
  const isDisabled = onChange == null;

  return (
    <div className={classNames(classes.root, className)}>
      <div className={classes.header}>
        <Text variant="subtitle1">
          <fbt desc="">Workforce Data</fbt>
        </Text>
        <Text variant="body2" color="gray">
          <fbt desc="">
            Choose the permissions this policy should include for modifying the
            selected work orders and projects.
          </fbt>
        </Text>
      </div>
      <Switch
        className={classes.readRule}
        title={fbt('View work orders and projects', '')}
        checked={readAllowed}
        disabled={isDisabled}
        onChange={checked =>
          callOnChange({
            ...policy,
            read: {
              ...policy.read,
              isAllowed: bool2PermissionRuleValue(checked),
            },
          })
        }
      />
      <PermissionsPolicyWorkforceDataRulesSpecification
        policy={policy}
        onChange={callOnChange}
        disabled={isDisabled}
      />
      <WorkforceDataRulesSection
        disabled={isDisabled || !readAllowed}
        rule={policy.data}
        onChange={data =>
          callOnChange({
            ...policy,
            data,
          })
        }
      />
    </div>
  );
}
