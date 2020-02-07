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
import {FutureStateValues} from '../../common/WorkOrder';

const PowerSearchLinkFutureStateFilter = (props: FilterProps) => {
  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;
  return (
    <PowerSearchFilter
      name="Future State"
      operator={value.operator}
      editMode={editMode}
      value={(value.idSet ?? [])
        .map(id => FutureStateValues.find(status => status.value === id)?.label)
        .join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <MutipleSelectInput
          options={FutureStateValues}
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

export default PowerSearchLinkFutureStateFilter;
