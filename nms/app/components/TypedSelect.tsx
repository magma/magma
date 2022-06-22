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
 */

import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import {InputBaseComponentProps} from '@material-ui/core/InputBase/InputBase';

type Props<T extends string | number> = {
  value: T;
  onChange: (value: T) => void;
  items: Record<T, string>;
  className?: string;
  input?: React.ReactElement<any, any>;
  disabled?: boolean;
  fullWidth?: boolean;
  inputProps?: InputBaseComponentProps;
};

export default function TypedSelect<T extends string | number>(
  props: Props<T>,
) {
  const {onChange, ...otherProps} = props;
  return (
    <Select
      {...otherProps}
      onChange={({target}) => onChange(target.value as T)}>
      {(Object.keys(props.items) as Array<T>).map(key => (
        <MenuItem value={key} key={key}>
          {props.items[key]}
        </MenuItem>
      ))}
    </Select>
  );
}
