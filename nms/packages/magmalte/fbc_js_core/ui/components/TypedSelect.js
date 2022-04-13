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

import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

type Props<T: string | number> = {
  value: T,
  onChange: T => void,
  items: {[T]: string},
};

export default function TypedSelect<T: string | number>(props: Props<T>) {
  const {onChange, ...otherProps} = props;
  return (
    <Select
      {...otherProps}
      // $FlowIgnore the selected values can only be the values in the MenuItems
      onChange={({target}) => onChange(((target.value: any): T))}>
      {Object.keys(props.items).map(key => (
        <MenuItem value={key} key={key}>
          {props.items[key]}
        </MenuItem>
      ))}
    </Select>
  );
}
