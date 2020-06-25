/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {LocationCUDPermissions} from '../data/PermissionsPolicies';

import * as React from 'react';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import LocationTypesTokenizer from '../../../../common/LocationTypesTokenizer';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import fbt from 'fbt';
import useFeatureFlag from '@fbcnms/ui/context/useFeatureFlag';
import {
  PERMISSION_RULE_VALUES,
  permissionRuleValue2Bool,
} from '../data/PermissionsPolicies';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo, useState} from 'react';
import {useFormAlertsContext} from '@fbcnms/ui/components/design-system/Form/FormAlertsContext';

const ERROR_MESSAGE_HEIGHT = '6px';

const useStyles = makeStyles(() => ({
  policySpecificationContainer: {
    display: 'flex',
    flexDirection: 'column',
    padding: '2px 16px 0 7px',
  },
  policyMethodSelection: {
    display: 'flex',
    flexDirection: 'column',
    width: 'fit-content',
    '& > *': {
      marginBottom: '4px',
    },
  },
  permissionMethodSelect: {
    '&&': {
      paddingLeft: '8px',
      marginRight: '16px',
    },
  },
  locationTypesSelection: {
    marginTop: '16px',
    minHeight: '52px',
    marginBottom: `-${ERROR_MESSAGE_HEIGHT}`,
  },
  hidden: {
    display: 'none',
  },
}));

type Props = $ReadOnly<{|
  locationRule: LocationCUDPermissions,
  onChange: LocationCUDPermissions => void,
  disabled?: ?boolean,
  className?: ?string,
|}>;

const METHOD_ALL_LOCATIONS_VALUE = 0;
const METHOD_SELECTED_LOCATIONS_VALUE = 1;

export default function PermissionsPolicyLocationRulesSpecification(
  props: Props,
) {
  const {locationRule, onChange, disabled, className} = props;
  const classes = useStyles();

  const policyMethods = useMemo(() => {
    const methods = [];
    methods[METHOD_ALL_LOCATIONS_VALUE] = {
      label: <fbt desc="">All locations</fbt>,
      value: METHOD_ALL_LOCATIONS_VALUE,
      key: METHOD_ALL_LOCATIONS_VALUE,
    };
    methods[METHOD_SELECTED_LOCATIONS_VALUE] = {
      label: <fbt desc="">Selected</fbt>,
      value: METHOD_SELECTED_LOCATIONS_VALUE,
      key: METHOD_SELECTED_LOCATIONS_VALUE,
    };
    return methods;
  }, []);

  const selectedLocationTypesCount =
    locationRule.update.locationTypeIds?.length || 0;
  const [policyMethod, setPolicyMethod] = useState(
    selectedLocationTypesCount > 0
      ? METHOD_SELECTED_LOCATIONS_VALUE
      : METHOD_ALL_LOCATIONS_VALUE,
  );

  const callSetPermissionMethod = useCallback(
    newPermissionMethod => {
      setPolicyMethod(newPermissionMethod);
      onChange({
        ...locationRule,
        update: {
          ...locationRule.update,
          isAllowed:
            newPermissionMethod === METHOD_SELECTED_LOCATIONS_VALUE
              ? PERMISSION_RULE_VALUES.BY_CONDITION
              : PERMISSION_RULE_VALUES.YES,
        },
      });
    },
    [onChange, locationRule],
  );

  const alerts = useFormAlertsContext();
  const emptyRequiredTypesSelectionErrorMessage = alerts.error.check({
    fieldId: 'location_types_selection',
    fieldDisplayName: 'Policies applied location types selection',
    value:
      permissionRuleValue2Bool(locationRule.update.isAllowed) &&
      policyMethod === METHOD_SELECTED_LOCATIONS_VALUE &&
      selectedLocationTypesCount === 0,
    checkCallback: missingRequiredSelection =>
      missingRequiredSelection
        ? `${fbt('At least one location type must be selected.', '')}`
        : '',
  });

  const isPermissionPolicyPerTypeEnabled = useFeatureFlag(
    'permission_policy_per_type',
  );

  if (!isPermissionPolicyPerTypeEnabled) {
    return null;
  }

  return (
    <div
      className={classNames(classes.policySpecificationContainer, className)}>
      <div className={classes.policyMethodSelection}>
        <Text>
          {disabled == true ? (
            <fbt desc="">Location types this policy applies to</fbt>
          ) : (
            <fbt desc="">Choose location types this policy applies to</fbt>
          )}
        </Text>
        <FormField disabled={disabled}>
          <Select
            options={policyMethods}
            selectedValue={policyMethod}
            onChange={callSetPermissionMethod}
            className={classes.permissionMethodSelect}
          />
        </FormField>
      </div>
      <div
        className={classNames(classes.locationTypesSelection, {
          [classes.hidden]: policyMethod !== METHOD_SELECTED_LOCATIONS_VALUE,
        })}>
        <FormField
          disabled={disabled}
          errorText={emptyRequiredTypesSelectionErrorMessage}
          hasError={!!emptyRequiredTypesSelectionErrorMessage}>
          <LocationTypesTokenizer
            selectedLocationTypeIds={locationRule.update.locationTypeIds}
            onSelectedLocationTypesIdsChange={newLocationTypeIds =>
              onChange({
                ...locationRule,
                update: {
                  ...locationRule.update,
                  locationTypeIds: newLocationTypeIds,
                },
              })
            }
          />
        </FormField>
      </div>
    </div>
  );
}
