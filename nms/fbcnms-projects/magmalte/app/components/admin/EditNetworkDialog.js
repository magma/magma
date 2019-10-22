/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {network} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';

import {AllNetworkTypes} from '@fbcnms/types/network';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
});

type Props = {
  onClose: () => void,
  onSave: () => void,
};

export default function NetworkDialog(props: Props) {
  const classes = useStyles();
  const editingNetworkID = useRouter().match.params.networkID;
  const [networkConfig, setNetworkConfig] = useState<?network>(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    MagmaV1API.getNetworksByNetworkId({networkId: editingNetworkID})
      .then(setNetworkConfig)
      .catch(error => setError(error));
  }, [editingNetworkID]);

  if (!networkConfig) {
    return <LoadingFillerBackdrop />;
  }

  const updateNetwork = (
    field: 'name' | 'description' | 'type',
    value: string,
  ) =>
    setNetworkConfig({
      ...networkConfig,
      [field]: value,
    });

  const onSave = () => {
    MagmaV1API.putNetworksByNetworkId({
      networkId: networkConfig.id,
      network: networkConfig,
    })
      .then(props.onSave)
      .catch(error => setError(error.response?.data?.error || error));
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>Editing "{networkConfig.id}"</DialogTitle>
      <DialogContent>
        {error && <FormLabel error>{error}</FormLabel>}
        <TextField
          name="name"
          label="Name"
          className={classes.input}
          value={networkConfig.name}
          onChange={({target}) => updateNetwork('name', target.value)}
        />
        <TextField
          name="description"
          label="Description"
          className={classes.input}
          value={networkConfig.description}
          onChange={({target}) => updateNetwork('description', target.value)}
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="types">Network Type</InputLabel>
          <Select
            value={networkConfig.type}
            onChange={({target}) => updateNetwork('type', target.value)}
            input={<Input id="types" />}>
            {AllNetworkTypes.map(type => (
              <MenuItem key={type} value={type}>
                <ListItemText primary={type} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button onClick={onSave} color="primary" variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
