/*
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

import FormControl from '@material-ui/core/FormControl';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import type {SelectProps} from './ActionTable';

/**
 * SelectEdit is a dropdown menu used for editing a field.
 */
export function SelectEdit(props: SelectProps) {
  if (props.value === undefined || props.value === null) {
    if (props.defaultValue !== undefined) {
      props.onChange(props.defaultValue);
      return null;
    }
  }
  return (
    <FormControl>
      <Select
        data-testid={props.testId ?? ''}
        value={props.value}
        onChange={({target}) => props.onChange(target.value)}
        input={<OutlinedInput />}>
        {props.content.map((k: string, idx: number) => (
          <MenuItem key={idx} value={k}>
            {k}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
}
