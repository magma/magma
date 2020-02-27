/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
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
