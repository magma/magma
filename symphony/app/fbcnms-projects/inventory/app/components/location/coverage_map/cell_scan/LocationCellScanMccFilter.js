/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

`use strict`;

import type {FilterProps} from '../../../comparison_view/ComparisonViewTypes';

import * as React from 'react';
import PowerSearchFilter from '../../../comparison_view/PowerSearchFilter';
import TextInput from '../../../comparison_view/TextInput';

const LocationCellScanMccFilter = (props: FilterProps) => {
  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;
  return (
    <PowerSearchFilter
      name="MCC"
      operator={value.operator}
      editMode={editMode}
      value={value.stringValue}
      onRemoveFilter={onRemoveFilter}
      input={
        <TextInput
          type="text"
          onSubmit={onInputBlurred}
          onBlur={onInputBlurred}
          value={value.stringValue ?? ''}
          onChange={newMcc =>
            onValueChanged({
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              stringValue: newMcc,
            })
          }
        />
      }
    />
  );
};

export default LocationCellScanMccFilter;
