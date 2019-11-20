/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ActionData} from './types';

import Autocomplete from '@material-ui/lab/Autocomplete';
import CancelIcon from '@material-ui/icons/Cancel';
import Chip from '@material-ui/core/Chip';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  popup: {
    zIndex: theme.zIndex.modal + 100,
  },
}));

type Props = {
  onChange: (SyntheticInputEvent<*>, string[]) => void,
  value: ActionData,
  options: string[],
};

export default function ActionsAutocomplete(props: Props) {
  const classes = useStyles();

  const renderTags = (value, {className, onDelete}) => {
    return value.map((option, index) => {
      return (
        <Chip
          key={index}
          variant="outlined"
          tabIndex={-1}
          label={option}
          className={className}
          deleteIcon={<CancelIcon data-tag-index={index} />}
          onDelete={evt => onDelete(evt)}
        />
      );
    });
  };

  return (
    <Autocomplete
      classes={{
        popup: classes.popup,
      }}
      autoComplete
      filterSelectedOptions
      freeSolo
      multiple
      onChange={props.onChange}
      options={props.options}
      value={props.value}
      renderTags={renderTags}
      renderInput={params => (
        <TextField {...params} variant="standard" margin="normal" fullWidth />
      )}
    />
  );
}
