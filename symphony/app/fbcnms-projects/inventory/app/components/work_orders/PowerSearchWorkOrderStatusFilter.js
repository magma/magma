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
import {statusValues} from '../../common/WorkOrder';

const PowerSearchWorkOrderStatusFilter = (props: FilterProps) => {
  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;
  return (
    <PowerSearchFilter
      name="Status"
      operator={value.operator}
      editMode={editMode}
      value={(value.idSet ?? [])
        .map(id => statusValues.find(status => status.value === id)?.label)
        .join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <MutipleSelectInput
          options={statusValues}
          onSubmit={onInputBlurred}
          onBlur={onInputBlurred}
          value={value.idSet ?? []}
          onChange={newEntries => {
            onValueChanged({
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              idSet: newEntries,
            });
          }}
        />
      }
    />
  );
};

export default PowerSearchWorkOrderStatusFilter;
