/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {symphony_agent} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import FormLabel from '@material-ui/core/FormLabel';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';

import {DEFAULT_MAGMAD_CONFIGS} from './DevicesUtils';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  title: {
    margin: '15px 0 5px',
  },
}));

type Props = {
  onClose: () => void,
  onSave: symphony_agent => void,
};

export default function DevicesNewAgentDialog(props: Props) {
  const {match} = useRouter();
  const classes = useStyles();
  const [hardwareId, setHardwareId] = useState('');
  const [name, setName] = useState('');
  const [error, setError] = useState('');

  function onSave() {
    const sanitizedHardwareId = hardwareId.toLowerCase();
    const agentId = name.replace(/[^a-zA-z0-9]/g, '_').toLowerCase();

    if (
      sanitizedHardwareId.match(
        /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/,
      ) == null
    ) {
      throw 'Invalid UUID, should be in format of "01234567-0123-0123-0123-0123456789ab"';
    }

    const symphonyAgent: symphony_agent = {
      description: agentId,
      device: {
        hardware_id: sanitizedHardwareId,
        key: {key_type: 'ECHO'},
      },
      id: agentId,
      magmad: DEFAULT_MAGMAD_CONFIGS,
      managed_devices: [],
      name: name,
      tier: 'default',
    };

    MagmaV1API.postSymphonyByNetworkIdAgents({
      networkId: nullthrows(match.params.networkId),
      symphonyAgent,
    })
      .then(() => props.onSave(symphonyAgent))
      .catch(err => {
        setError(err.response?.data?.message || err.toString());
      });
  }

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>New Agent</DialogTitle>
      <DialogContent>
        {error ? <FormLabel error>{error}</FormLabel> : null}
        <FormGroup row>
          <TextField
            required
            className={classes.input}
            label="Hardware UUID"
            margin="normal"
            value={hardwareId}
            onChange={event => setHardwareId(event.target.value)}
          />
          <TextField
            required
            className={classes.input}
            label="ID"
            margin="normal"
            value={name}
            onChange={event => setName(event.target.value)}
          />
        </FormGroup>
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
