/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FilterProps} from './ComparisonViewTypes';

import * as React from 'react';
import PowerSearchFilter from './PowerSearchFilter';
import TextInput from './TextInput';

const PowerSearchEquipmentNameFilter = (props: FilterProps) => {
  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
  } = props;
  return (
    <PowerSearchFilter
      name="Equipment Name"
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
          onChange={newName =>
            onValueChanged({
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              stringValue: newName,
            })
          }
        />
      }
    />
  );
};

export default PowerSearchEquipmentNameFilter;
