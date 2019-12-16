/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FilterProps} from '../comparison_view/ComparisonViewTypes';

import * as React from 'react';
import MenuItem from '@material-ui/core/MenuItem';
import PowerSearchFilter from '../comparison_view/PowerSearchFilter';
import Select from '@material-ui/core/Select';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  textField: {
    backgroundColor: symphony.palette.D10,
  },
  input: {
    paddingTop: '10px',
    paddingBottom: '4px',
    paddingRight: '30px',
  },
}));

const ENTER_KEY_CODE = 13;

const boolValues = [
  {
    value: 'false',
    label: 'False',
  },
  {
    value: 'true',
    label: 'True',
  },
];

type Props = FilterProps & {
  label: string,
};

const PowerSearchBoolFilter = (props: Props) => {
  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
    label,
  } = props;
  const classes = useStyles();
  return (
    <PowerSearchFilter
      name={label}
      operator={value.operator}
      editMode={editMode}
      value={(value.boolValue ?? '').toString()}
      onRemoveFilter={onRemoveFilter}
      input={
        <Select
          autoFocus={true}
          onBlur={onInputBlurred}
          onKeyDown={e => e.keyCode === ENTER_KEY_CODE && onInputBlurred()}
          value={(value.boolValue ?? '').toString()}
          inputProps={{autoComplete: 'off', className: classes.input}}
          className={classes.textField}
          margin="none"
          variant="outlined"
          onChange={event => {
            const newValue = {
              id: value.id,
              key: value.key,
              name: value.name,
              operator: value.operator,
              boolValue: event.target.value === 'true' ? true : false,
            };
            onValueChanged(newValue);
          }}>
          {boolValues.map(option => (
            <MenuItem key={option.value} value={option.value}>
              {option.label}
            </MenuItem>
          ))}
        </Select>
      }
    />
  );
};

export default PowerSearchBoolFilter;
