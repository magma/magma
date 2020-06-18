/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import CheckIcon from '@material-ui/icons/Check';
import EditIcon from '@material-ui/icons/Edit';
import React, {useState} from 'react';
import Text from './design-system/Text';
import TextField from '@material-ui/core/TextField';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    alignItems: 'center',
    display: 'flex',
    flexDirection: 'row',
    '&:hover': {
      '& $button': {
        display: 'block',
      },
    },
  },
  button: {
    cursor: 'pointer',
    fontSize: '20px',
    color: theme.palette.primary.main,
    marginLeft: '8px',
    display: 'none',
  },
  input: {
    display: 'inline-block',
    marginBottom: '0',
    width: '100%',
  },
  inputContainer: {
    display: 'inline-block',
    width: '200px',
  },
  emptyString: {
    width: '200px',
    height: '100%',
    display: 'inline-block',
    marginBottom: '0',
  },
}));

type Props = {
  value: ?string,
  onSave: string => boolean,
  type: 'string' | 'date',
  editDisabled: boolean,
};

export default function EditableField(props: Props) {
  const classes = useStyles();
  const [isEditing, setIsEditing] = useState(false);
  const [editedText, setEditedText] = useState(props.value);

  function onClickSave() {
    if (editedText != null) {
      if (!props.onSave(editedText)) {
        return;
      }
    }
    setIsEditing(false);
    setEditedText(null);
  }

  if (props.editDisabled || !isEditing) {
    return (
      <div className={classes.root}>
        {props.value || props.editDisabled ? (
          <Text variant="body2">{props.value}</Text>
        ) : (
          <Text variant="body2" color="regular">
            Set...
          </Text>
        )}
        <EditIcon
          className={classes.button}
          onClick={() => setIsEditing(true)}
        />
      </div>
    );
  }

  return (
    <div className={classes.root}>
      <div className={classes.inputContainer}>
        {props.type == 'date' ? (
          <TextField
            type="date"
            className={classes.input}
            InputLabelProps={{
              shrink: true,
            }}
            onChange={event => setEditedText(event.target.value)}
          />
        ) : (
          <TextField
            defaultValue={props.value || ''}
            className={classes.input}
            variant="outlined"
            margin="dense"
            onChange={event => setEditedText(event.target.value)}
          />
        )}
      </div>
      <CheckIcon className={classes.button} onClick={onClickSave} />
    </div>
  );
}
