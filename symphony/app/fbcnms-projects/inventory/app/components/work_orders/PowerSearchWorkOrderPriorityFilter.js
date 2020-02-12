/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FilterProps} from '../comparison_view/ComparisonViewTypes';

import * as React from 'react';
import MutipleSelectInput from '../comparison_view/MutipleSelectInput';
import PowerSearchFilter from '../comparison_view/PowerSearchFilter';
import {priorityValues} from '../../common/WorkOrder';

const PowerSearchWorkOrderPriorityFilter = (props: FilterProps) => {
  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;
  return (
    <PowerSearchFilter
      name="Priority"
      operator={value.operator}
      editMode={editMode}
      onRemoveFilter={onRemoveFilter}
      value={(value.stringSet ?? [])
        .map(
          value =>
            priorityValues.find(priority => priority.value === value)?.label,
        )
        .join(', ')}
      input={
        <MutipleSelectInput
          options={priorityValues}
          onSubmit={onInputBlurred}
          onBlur={onInputBlurred}
          value={value.stringSet ?? []}
          onChange={newName =>
            onValueChanged({
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              stringSet: newName,
            })
          }
        />
      }
    />
  );
};

export default PowerSearchWorkOrderPriorityFilter;
