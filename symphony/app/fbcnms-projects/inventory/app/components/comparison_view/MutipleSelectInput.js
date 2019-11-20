/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  textField: {
    backgroundColor: theme.palette.grey.A100,
  },
  input: {
    paddingTop: '6px',
    paddingBottom: '6px',
  },
}));

const ENTER_KEY_CODE = 13;

type Props = {
  onBlur: () => void,
  onSubmit: () => void,
  value: Array<string>,
  options: Array<{value: string | number, label: string}>,
  onChange: (newValue: Array<string>) => void,
};

const MutipleSelectInput = (props: Props) => {
  const classes = useStyles();
  const {value, onChange, onBlur, options} = props;
  return (
    <Select
      autoFocus={true}
      multiple
      onBlur={onBlur}
      onKeyDown={e => e.keyCode === ENTER_KEY_CODE && onBlur()}
      value={value}
      inputProps={{autoComplete: 'off', className: classes.input}}
      className={classes.textField}
      margin="none"
      variant="outlined"
      onChange={({target}) => {
        onChange((target.value: any));
      }}>
      {options.map(option => (
        <MenuItem key={option.value} value={option.value}>
          {option.label}
        </MenuItem>
      ))}
    </Select>
  );
};

export default MutipleSelectInput;
