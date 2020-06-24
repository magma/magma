/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {WorkforcePolicy} from '../data/PermissionsPolicies';

import * as React from 'react';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import ProjectTemplatesTokenizer from '../../../../common/ProjectTemplatesTokenizer';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import WorkOrderTemplatesTokenizer from '../../../../common/WorkOrderTemplatesTokenizer';
import classNames from 'classnames';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import useFeatureFlag from '@fbcnms/ui/context/useFeatureFlag';
import {PERMISSION_RULE_VALUES} from '../data/PermissionsPolicies';
import {makeStyles} from '@material-ui/styles';
import {permissionRuleValue2Bool} from '../data/PermissionsPolicies';
import {useCallback, useMemo, useState} from 'react';
import {useFormAlertsContext} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';

const useStyles = makeStyles(() => ({
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
}));

type Props = $ReadOnly<{|
  policy: ?WorkforcePolicy,
  disabled?: ?boolean,
  onChange: WorkforcePolicy => void,
  className?: ?string,
|}>;

const METHOD_ALL_TYPES_VALUE = 0;
const METHOD_SELECTED_TYPES_VALUE = 1;

export default function PermissionsPolicyWorkforceDataRulesSpecification(
  props: Props,
) {
  const {policy, onChange, disabled, className} = props;
  const classes = useStyles();

  const policyMethods = useMemo(() => {
    const methods = [];
    methods[METHOD_ALL_TYPES_VALUE] = {
      label: <fbt desc="">All Types</fbt>,
      value: METHOD_ALL_TYPES_VALUE,
      key: METHOD_ALL_TYPES_VALUE,
    };
    methods[METHOD_SELECTED_TYPES_VALUE] = {
      label: <fbt desc="">Based on selected templates</fbt>,
      value: METHOD_SELECTED_TYPES_VALUE,
      key: METHOD_SELECTED_TYPES_VALUE,
    };
    return methods;
  }, []);

  const selectedTypesCount =
    policy?.read.projectTypeIds?.length ||
    0 + policy?.read.workOrderTypeIds?.length ||
    0;
  const [policyMethod, setPolicyMethod] = useState(
    selectedTypesCount > 0
      ? METHOD_SELECTED_TYPES_VALUE
      : METHOD_ALL_TYPES_VALUE,
  );

  const callSetPolicyMethod = useCallback(
    newPolicyMethod => {
      if (policy == null) {
        return;
      }

      setPolicyMethod(newPolicyMethod);
      onChange({
        ...policy,
        read: {
          ...policy.read,
          isAllowed:
            newPolicyMethod === METHOD_SELECTED_TYPES_VALUE
              ? PERMISSION_RULE_VALUES.BY_CONDITION
              : PERMISSION_RULE_VALUES.YES,
        },
      });
    },
    [onChange, policy],
  );

  const alerts = useFormAlertsContext();
  const emptyRequiredTypesSelectionErrorMessage = alerts.error.check({
    fieldId: 'workforce_types_selection',
    fieldDisplayName: 'Policies applied workforce types selection',
    value:
      policy != null &&
      permissionRuleValue2Bool(policy.read.isAllowed) &&
      policyMethod === METHOD_SELECTED_TYPES_VALUE &&
      selectedTypesCount === 0,
    checkCallback: missingRequiredSelection =>
      missingRequiredSelection
        ? `${fbt(
            'At least one Work Order or Project type must be selected.',
            '',
          )}`
        : '',
  });

  const permissionPolicyPerTypeEnabled = useFeatureFlag(
    'permission_policy_per_type',
  );

  if (policy == null || !permissionPolicyPerTypeEnabled) {
    return null;
  }

  return (
    <div
      className={classNames(classes.policySpecificationContainer, className)}>
      <div className={classes.methodSelectionBox}>
        <Text>
          {disabled == true ? (
            <fbt desc="">Work Orders and Projects this policy applies to</fbt>
          ) : (
            <fbt desc="">
              Choose Work Orders and Projects this policy applies to
            </fbt>
          )}
        </Text>
        <FormField disabled={disabled}>
          <Select
            options={policyMethods}
            selectedValue={policyMethod}
            onChange={callSetPolicyMethod}
            className={classes.policyMethodSelect}
          />
        </FormField>
      </div>
      <div
        className={classNames({
          [classes.hidden]: policyMethod !== METHOD_SELECTED_TYPES_VALUE,
        })}>
        <FormField
          className={classes.templatesFieldContainer}
          disabled={disabled}
          hasError={emptyRequiredTypesSelectionErrorMessage != null}
          label={`${fbt('Work order templates', '')}`}>
          <WorkOrderTemplatesTokenizer
            selectedWorkOrderTemplateIds={policy.read.workOrderTypeIds}
            onSelectedWorkOrderTemplateIdsChange={newWorkOrderTypeIds =>
              onChange({
                ...policy,
                read: {
                  ...policy.read,
                  workOrderTypeIds: newWorkOrderTypeIds,
                },
              })
            }
          />
        </FormField>
        <FormField
          className={classes.templatesFieldContainer}
          disabled={disabled}
          hasError={emptyRequiredTypesSelectionErrorMessage != null}
          errorText={emptyRequiredTypesSelectionErrorMessage}
          label={`${fbt('Project templates', '')}`}>
          <ProjectTemplatesTokenizer
            selectedProjectTemplateIds={policy.read.projectTypeIds}
            onSelectedProjectTemplateIdsChange={newProjectTypeIds =>
              onChange({
                ...policy,
                read: {
                  ...policy.read,
                  projectTypeIds: newProjectTypeIds,
                },
              })
            }
          />
        </FormField>
      </div>
    </div>
  );
}
