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

import Checkbox from '../Checkbox/Checkbox';
import React from 'react';
import {useSelection} from './TableSelectionContext';

const TableHeaderCheckbox = () => {
  const {selectionMode, changeHeaderSelectionMode} = useSelection();
  return (
    <Checkbox
      indeterminate={selectionMode === 'some'}
      checked={selectionMode === 'all'}
      title={null}
      onChange={changeHeaderSelectionMode}
    />
  );
};

export default TableHeaderCheckbox;
