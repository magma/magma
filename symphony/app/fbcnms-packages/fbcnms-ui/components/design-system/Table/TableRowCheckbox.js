/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {SelectionType} from '../Checkbox/Checkbox';

import Checkbox from '../Checkbox/Checkbox';
import React, {useMemo} from 'react';
import {useSelection} from './TableSelectionContext';

type Props = {
  id: string | number,
};

const TableRowCheckbox = ({id}: Props) => {
  const {selectedIds, changeRowSelection} = useSelection();
  const checked = useMemo(() => selectedIds.includes(id), [selectedIds, id]);
  return (
    <Checkbox
      checked={checked}
      onChange={(selection: SelectionType) => changeRowSelection(id, selection)}
      onClick={e => e.stopPropagation()}
    />
  );
};

export default TableRowCheckbox;
