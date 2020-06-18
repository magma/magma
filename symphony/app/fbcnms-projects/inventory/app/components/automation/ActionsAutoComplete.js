/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ActionData} from './types';

import Autocomplete from '@material-ui/lab/Autocomplete';
import CancelIcon from '@material-ui/icons/Cancel';
import Chip from '@material-ui/core/Chip';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import {findIndex} from 'lodash';

type Props = {
  onChange: (SyntheticEvent<*>, string[]) => void,
  value: ActionData,
  options: string[],
};

export default function ActionsAutocomplete(props: Props) {
  const renderTags = (value, {className}) => {
    return value.map((option, index) => {
      const onDelete = evt => {
        const i = findIndex(value, item => item === option);
        if (i !== -1) {
          const newValue = [...value];
          newValue.splice(i, 1);
          props.onChange(evt, newValue);
        }
      };

      return (
        <Chip
          key={index}
          variant="outlined"
          tabIndex={-1}
          label={option}
          className={className}
          deleteIcon={<CancelIcon data-tag-index={index} />}
          onDelete={onDelete}
        />
      );
    });
  };

  return (
    <Autocomplete
      autoComplete
      filterSelectedOptions
      freeSolo
      multiple
      onChange={props.onChange}
      options={props.options}
      value={props.value}
      renderTags={renderTags}
      renderInput={params => {
        return (
          <TextField
            ref={params.ref}
            InputLabelProps={params.InputLabelProps}
            InputProps={params.InputProps}
            inputProps={params.inputProps}
            variant="standard"
            margin="normal"
            fullWidth
          />
        );
      }}
    />
  );
}
