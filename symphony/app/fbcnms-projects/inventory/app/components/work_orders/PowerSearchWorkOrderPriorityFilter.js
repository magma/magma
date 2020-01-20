/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
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
      value={(value.idSet ?? [])
        .map(
          id => priorityValues.find(priority => priority.value === id)?.label,
        )
        .join(', ')}
      input={
        <MutipleSelectInput
          options={priorityValues}
          onSubmit={onInputBlurred}
          onBlur={onInputBlurred}
          value={value.idSet ?? []}
          onChange={newName =>
            onValueChanged({
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              idSet: newName,
            })
          }
        />
      }
    />
  );
};

export default PowerSearchWorkOrderPriorityFilter;
