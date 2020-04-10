/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TimeRange} from './AsyncMetric';

import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

type Props = {
  value: TimeRange,
  onChange: TimeRange => void,
  className: string,
};

export default function TimeRangeSelector(props: Props) {
  return (
    <FormControl variant="filled" className={props.className}>
      <InputLabel htmlFor="time_range">Period</InputLabel>
      <Select
        inputProps={{id: 'time_range'}}
        value={props.value}
        onChange={event => props.onChange((event.target.value: any))}>
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
