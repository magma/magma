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

import PowerSearchWorkOrderGeneralUserFilter from './PowerSearchWorkOrderGeneralUserFilter';
import React from 'react';

const PowerSearchWorkOrderAssigneeFilter = (props: FilterProps) => {
  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
    config,
    onNewInputBlurred,
  } = props;

  return (
    <PowerSearchWorkOrderGeneralUserFilter
      config={config}
      value={value}
      onInputBlurred={onInputBlurred}
      onValueChanged={onValueChanged}
      onRemoveFilter={onRemoveFilter}
      onNewInputBlurred={value => onNewInputBlurred(value)}
      editMode={editMode}
      title={'Assignee'}
    />
  );
};

export default PowerSearchWorkOrderAssigneeFilter;
