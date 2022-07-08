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

import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import type {TimeRange} from './AsyncMetric';

type Props = {
  value: TimeRange;
  onChange: (timeRange: TimeRange) => void;
  className: string;
};

export default function TimeRangeSelector(props: Props) {
  return (
    <FormControl variant="filled" className={props.className}>
      <InputLabel htmlFor="time_range">Period</InputLabel>
      <Select
        inputProps={{id: 'time_range'}}
        value={props.value}
        // TODO[ts-migration]: this was migrated from fbcnms-ui and types are broken
        // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
        onChange={event => props.onChange(event.target.value as any)}>
        <MenuItem value="3_hours">Last 3 hours</MenuItem>
        <MenuItem value="6_hours">Last 6 hours</MenuItem>
        <MenuItem value="12_hours">Last 12 hours</MenuItem>
        <MenuItem value="24_hours">Last 24 hours</MenuItem>
        <MenuItem value="7_days">Last 7 days</MenuItem>
        <MenuItem value="14_days">Last 14 days</MenuItem>
        <MenuItem value="30_days">Last 30 days</MenuItem>
      </Select>
    </FormControl>
  );
}
