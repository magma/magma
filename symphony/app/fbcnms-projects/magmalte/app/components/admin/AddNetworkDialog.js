/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import nullthrows from '@fbcnms/util/nullthrows';
import {AllNetworkTypes} from '@fbcnms/types/network';
import {CWF, FEG} from '@fbcnms/types/network';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  onClose: () => void,
  onSave: string => void,
};

export default function NetworkDialog(props: Props) {
  const classes = useStyles();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [networkID, setNetworkId] = useState('');
  const [networkType, setNetworkType] = useState('');
  const [fegNetworkID, setFegNetworkID] = useState('');
  const [servedNetworkIDs, setServedNetworkIDs] = useState('');
  const [error, setError] = useState(null);

  const onSave = () => {
    const payload = {
      networkID,
      data: {
        name,
        description,
        networkType,
        fegNetworkID,
        servedNetworkIDs,
      },
    };
    axios
      .post('/nms/network/create', payload)
      .then(response => {
        if (response.data.success) {
          props.onSave(nullthrows(networkID));
        } else {
          setError(response.data.message);
        }
      })
      .catch(error => setError(error.response?.data?.error || error));
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>Add Network</DialogTitle>
      <DialogContent>
        {error && <FormLabel error>{error}</FormLabel>}
        <TextField
          name="name"
          label="Network ID"
          className={classes.input}
          value={networkID}
          onChange={({target}) => setNetworkId(target.value)}
        />
        <TextField
          name="name"
          label="Name"
          className={classes.input}
          value={name}
          onChange={({target}) => setName(target.value)}
        />
        <TextField
          name="description"
          label="Description"
          className={classes.input}
          value={description}
          onChange={({target}) => setDescription(target.value)}
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="types">Network Type</InputLabel>
          <Select
            value={networkType}
            onChange={({target}) => setNetworkType(target.value)}
            input={<Input id="types" />}>
            {AllNetworkTypes.map(type => (
              <MenuItem key={type} value={type}>
                <ListItemText primary={type} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        {networkType === CWF && (
          <TextField
            name="fegNetworkID"
            label="Federation Network ID"
            className={classes.input}
            value={fegNetworkID}
            onChange={({target}) => setFegNetworkID(target.value)}
          />
        )}
        {networkType === FEG && (
          <TextField
            placeholder="network1,network2"
            label="Served Network IDs"
            className={classes.input}
            value={servedNetworkIDs}
            onChange={({target}) => setServedNetworkIDs(target.value)}
          />
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}
