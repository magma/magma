/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItemFillingProps} from './CheckListItemFilling';

import * as React from 'react';
import MultiSelect from '@fbcnms/ui/components/design-system/Select/MultiSelect';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import fbt from 'fbt';
import {enumStringToArray} from '../ChecklistUtils';
import {makeStyles} from '@material-ui/styles';
import {useFormContext} from '../../../common/FormContext';
import {useMemo} from 'react';

const useStyles = makeStyles(() => ({
  select: {
    width: '100%',
  },
}));

const MultipleChoiceCheckListItemFilling = ({
  item,
  onChange,
}: CheckListItemFillingProps): React.Node => {
  const classes = useStyles();
  const form = useFormContext();
  const enumArrayToOptions = (enumString: ?string) =>
    enumStringToArray(enumString).map(v => ({
      key: v,
      label: v,
      value: v,
    }));
  const options = useMemo(() => enumArrayToOptions(item.enumValues), [
    item.enumValues,
  ]);
  const selectedValues = useMemo(
    () => enumArrayToOptions(item.selectedEnumValues),
    [item.selectedEnumValues],
  );
  const updateOnChange = (selectedEnumValues: string) => {
    if (!onChange) {
      return;
    }
    onChange({
      ...item,
      selectedEnumValues,
    });
  };

  return item.enumSelectionMode === 'single' ? (
    <Select
      className={classes.select}
      label={<fbt desc="">Select option</fbt>}
      options={options}
      selectedValue={
        selectedValues.length > 0 &&
        selectedValues[0].value != null &&
        selectedValues[0] !== '' &&
        selectedValues[0].value !== ''
          ? selectedValues[0].value
          : null
      }
      onChange={value => updateOnChange(value)}
      disabled={form.alerts.missingPermissions.detected}
    />
  ) : (
    <MultiSelect
      className={classes.select}
      label={<fbt desc="">Select options</fbt>}
      options={options}
      onChange={option => {
        updateOnChange(
          (selectedValues.map(v => v.value).includes(option.value)
            ? selectedValues.filter(v => v.value !== option.value)
            : [...selectedValues, option]
          )
            .map(v => v.value)
            .join(','),
        );
      }}
      selectedValues={selectedValues}
      disabled={form.alerts.missingPermissions.detected}
    />
  );
};

export default MultipleChoiceCheckListItemFilling;
