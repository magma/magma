/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
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
