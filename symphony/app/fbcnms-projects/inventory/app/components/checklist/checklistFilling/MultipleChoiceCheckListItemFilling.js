/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItem} from '../checkListCategory/ChecklistItemsDialogMutateState';
import type {OptionProps} from '@fbcnms/ui/components/design-system/Select/SelectMenu';

import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import MultiSelect from '@fbcnms/ui/components/design-system/Select/MultiSelect';
import React, {useContext, useMemo} from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import fbt from 'fbt';
import {enumStringToArray} from '../ChecklistUtils';
import {makeStyles} from '@material-ui/styles';

type Props = {
  item: CheckListItem,
  onChange?: (updatedChecklistItem: CheckListItem) => void,
};

const useStyles = makeStyles(() => ({
  select: {
    width: '100%',
  },
}));

const MultipleChoiceCheckListItemFilling = ({item, onChange}: Props) => {
  const classes = useStyles();
  const validationContext = useContext(FormValidationContext);
  const enumArrayToOptions = (enumString: ?string) =>
    enumStringToArray(enumString).map(v => ({
      key: v,
      label: v,
      value: v,
    }));
  const options: Array<OptionProps<string>> = useMemo(
    () => enumArrayToOptions(item.enumValues),
    [item.enumValues],
  );
  const selectedValues: Array<OptionProps<string>> = useMemo(
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
      disabled={validationContext.editLock.detected}
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
      disabled={validationContext.editLock.detected}
    />
  );
};

export default MultipleChoiceCheckListItemFilling;
