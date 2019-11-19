/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Checkbox from '../Checkbox/Checkbox';
import React from 'react';
import {useSelection} from './TableSelectionContext';

const TableHeaderCheckbox = () => {
  const {selectionMode, changeHeaderSelectionMode} = useSelection();
  return (
    <Checkbox
      indeterminate={selectionMode === 'some'}
      checked={selectionMode === 'all'}
      onChange={changeHeaderSelectionMode}
    />
  );
};

export default TableHeaderCheckbox;
