/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

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
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import nullthrows from '@fbcnms/util/nullthrows';
import {AllNetworkTypes} from '@fbcnms/types/network';
import {MagmaAPIUrls} from '../../common/MagmaAPI';
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
  onSave: string => void,
};

export default function NetworkDialog(props: Props) {
  const classes = useStyles();
  const editingNetworkID = useRouter().match.params.networkID;
  const [name, setName] = useState('');
  const [networkID, setNetworkId] = useState<?string>(null);
  const [networkType, setNetworkType] = useState('');
  const [error, setError] = useState(null);

  useEffect(() => {
    if (editingNetworkID) {
      axios
        .get(MagmaAPIUrls.network(editingNetworkID))
        .then(({data}) => {
          setNetworkId(editingNetworkID);
          setName(data.name);
          setNetworkType(data?.features?.networkType || '');
        })
        .catch(error => setError(error));
    } else {
      setNetworkId('');
    }
  }, [editingNetworkID]);

  if (networkID === null) {
    return <LoadingFillerBackdrop />;
  }

  const successHandler = response => {
    if (response.data.success) {
      props.onSave(nullthrows(networkID));
    } else {
      setError(response.data.message);
    }
  };
  const onSave = () => {
    const payload = {
      networkID,
      data: {
        name,
        features: {
          networkType: networkType,
        },
      },
    };
    if (editingNetworkID) {
      axios.put('/nms/network/update/', payload).then(successHandler);
    } else {
      axios.post('/nms/network/create', payload).then(successHandler);
    }
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
          disabled={editingNetworkID}
        />
        <TextField
          name="name"
          label="Name"
          className={classes.input}
          value={name}
          onChange={({target}) => setName(target.value)}
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="types">Network Type</InputLabel>
          <Select
            value={networkType}
            onChange={({target}) => setNetworkType(target.value)}
            input={<Input id="tabs" />}>
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
