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
import {serviceStatusToVisibleNames} from '../../common/Service';

const PowerSearchServiceStatusFilter = (props: FilterProps) => {
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
      value={(value.stringSet ?? [])
        .map(status => serviceStatusToVisibleNames[status])
        .join(', ')}
      onRemoveFilter={onRemoveFilter}
      input={
        <MutipleSelectInput
          options={Object.entries(serviceStatusToVisibleNames).map(entry => {
            // $FlowFixMe - Flow doesn't value type well from object
            return {value: entry[0], label: entry[1]};
          })}
          onSubmit={onInputBlurred}
          onBlur={onInputBlurred}
          value={value.stringSet ?? []}
          onChange={newEntries => {
            onValueChanged({
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              stringSet: newEntries,
            });
          }}
        />
      }
    />
  );
};

export default PowerSearchServiceStatusFilter;
