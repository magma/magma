/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {allowed_gre_peers} from '@fbcnms/magma-api';

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Button from '@material-ui/core/Button';
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
  allowedGREPeers: allowed_gre_peers,
  onChange: allowed_gre_peers => void,
};

export default function(props: Props) {
  const classes = useStyles();
  const {allowedGREPeers} = props;

  const onChange = (index, field: 'ip' | 'key', value) => {
    const newValue = [...allowedGREPeers];
    newValue[index] = {...allowedGREPeers[index]};
    if (field === 'key') {
      newValue[index].key = parseInt(value.replace(/[^0-9]*/g, '') || 0);
    } else {
      newValue[index][field] = value;
    }
    props.onChange(newValue);
  };

  const removePeer = index => {
    const newValue = [...allowedGREPeers];
    newValue.splice(index, 1);
    props.onChange(newValue);
  };

  const addPeer = () => {
    const newValue = [...allowedGREPeers];
    newValue.push({ip: '', key: 0});
    props.onChange(newValue);
  };

  if (allowedGREPeers.length === 0) {
    return (
      <Button color="primary" variant="contained" onClick={addPeer}>
        Add GRE Peer
      </Button>
    );
  }

  return (
    <div>
      {allowedGREPeers.map((peer, index) => (
        <div className={classes.container} key={index}>
          <TextField
            label="IP"
            margin="none"
            value={peer.ip}
            onChange={({target}) => onChange(index, 'ip', target.value)}
            className={classes.inputKey}
          />
          <TextField
            label="Key"
            margin="none"
            value={peer.key}
            onChange={({target}) => onChange(index, 'key', target.value)}
            className={classes.inputValue}
          />
          <IconButton
            onClick={() => removePeer(index)}
            className={classes.icon}>
            <RemoveCircleOutline />
          </IconButton>
          {index === allowedGREPeers.length - 1 && (
            <IconButton onClick={addPeer} className={classes.icon}>
              <AddCircleOutline />
            </IconButton>
          )}
        </div>
      ))}
    </div>
  );
}
