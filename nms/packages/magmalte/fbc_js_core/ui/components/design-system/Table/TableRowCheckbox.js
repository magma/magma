/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
      title={null}
      onChange={(selection: SelectionType) => changeRowSelection(id, selection)}
      onClick={e => e.stopPropagation()}
    />
  );
};

export default TableRowCheckbox;
