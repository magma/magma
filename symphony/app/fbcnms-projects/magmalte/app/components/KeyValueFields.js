/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import IconButton from '@material-ui/core/IconButton';
import React from 'react';
import RemoveCircleOutline from '@material-ui/icons/RemoveCircleOutline';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  container: {
    display: 'block',
    margin: '5px 0',
    whiteSpace: 'nowrap',
    width: '100%',
  },
  inputKey: {
    width: '245px',
    paddingRight: '10px',
  },
  inputValue: {
    width: '240px',
  },
  icon: {
    width: '30px',
    height: '30px',
    verticalAlign: 'bottom',
  },
});

type Props = {
  keyValuePairs: Array<[string, string]>,
  onChange: (Array<[string, string]>) => void,
};

export default function KeyValueFields(props: Props) {
  const classes = useStyles();
  const onChange = (index, subIndex, value) => {
    const keyValuePairs = [...props.keyValuePairs];
    keyValuePairs[index] = [keyValuePairs[index][0], keyValuePairs[index][1]];
    keyValuePairs[index][subIndex] = value;
    props.onChange(keyValuePairs);
  };

  const removeField = index => {
    const keyValuePairs = [...props.keyValuePairs];
    keyValuePairs.splice(index, 1);
    props.onChange(keyValuePairs);
  };

  const addField = () => {
    props.onChange([...props.keyValuePairs, ['', '']]);
  };

  return (
    <>
      {props.keyValuePairs.map((pair, index) => (
        <div className={classes.container} key={index}>
          <TextField
            label="Key"
            margin="none"
            value={pair[0]}
            onChange={({target}) => onChange(index, 0, target.value)}
            className={classes.inputKey}
          />
          <TextField
            label="Value"
            margin="none"
            value={pair[1]}
            onChange={({target}) => onChange(index, 1, target.value)}
            className={classes.inputValue}
          />
          {props.keyValuePairs.length !== 1 && (
            <IconButton
              onClick={() => removeField(index)}
              className={classes.icon}>
              <RemoveCircleOutline />
            </IconButton>
          )}
          {index === props.keyValuePairs.length - 1 && (
            <IconButton onClick={addField} className={classes.icon}>
              <AddCircleOutline />
            </IconButton>
          )}
        </div>
      ))}
    </>
  );
}
