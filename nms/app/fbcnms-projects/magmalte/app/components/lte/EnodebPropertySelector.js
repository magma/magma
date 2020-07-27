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

import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

import {useState} from 'react';

type Props = {
  titleLabel: string,
  value: number | string,
  valueOptionsByKey: {+[string]: number | string},
  onChange: (SyntheticInputEvent<>) => void,
  className: string,
};

export default function EnodebPropertySelector(props: Props) {
  const [open, setOpen] = useState(false);
  const {className, valueOptionsByKey} = props;
  const valueOptionsArr = [];
  for (const property in valueOptionsByKey) {
    if (valueOptionsByKey.hasOwnProperty(property)) {
      valueOptionsArr.push(valueOptionsByKey[property]);
    }
  }

  const menuItems = valueOptionsArr.map(valueOption => {
    return (
      <MenuItem key={valueOption} value={valueOption}>
        {valueOption}
      </MenuItem>
    );
  });

  return (
    <form autoComplete="off">
      <FormControl className={className}>
        <InputLabel htmlFor="demo-controlled-open-select">
          eNodeB DL/UL Bandwidth (MHz)
        </InputLabel>
        <Select
          open={open}
          onClose={() => setOpen(false)}
          onOpen={() => setOpen(true)}
          value={props.value}
          onChange={props.onChange}
          inputProps={{
            name: props.titleLabel,
            id: 'demo-controlled-open-select',
          }}>
          {menuItems}
        </Select>
      </FormControl>
    </form>
  );
}
