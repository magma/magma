/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {Match} from 'react-router-dom';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import {DEFAULT_WIFI_GATEWAY_CONFIGS} from '../wifi/WifiUtils';
import {createDevice} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  title: {
    margin: '15px 0 5px',
  },
});

type Props = {
  onClose: () => void,
  onSave: any => void,
};

async function onSave(match: Match, hardwareId: string, name: string): any {
  const deviceId = name.replace(/[^a-zA-z0-9]/g, '_').toLowerCase();
  return await createDevice(
    deviceId,
    {
      hardware_id: hardwareId,
      key: {key_type: 'ECHO'},
    },
    'devmand', // type
    DEFAULT_WIFI_GATEWAY_CONFIGS,
    {managed_devices: []},
    match,
  );
}

export default function NewControllerDialog(props: Props) {
  const {match} = useRouter();
  const classes = useStyles();
  const [hardwareId, setHardwareId] = useState('');
  const [name, setName] = useState('');

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>New Controller</DialogTitle>
      <DialogContent>
        <FormGroup row>
          <TextField
            required
            className={classes.input}
            label="Hardware ID"
            margin="normal"
            value={hardwareId}
            onChange={event => setHardwareId(event.target.value)}
          />
          <TextField
            required
            className={classes.input}
            label="Name"
            margin="normal"
            value={name}
            onChange={event => setName(event.target.value)}
          />
        </FormGroup>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button
          onClick={() => onSave(match, hardwareId, name).then(props.onSave)}
          color="primary"
          variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
