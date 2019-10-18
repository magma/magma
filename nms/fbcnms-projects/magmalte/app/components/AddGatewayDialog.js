/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {lte_gateway} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React, {useState} from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
});

type Props = {|
  open: boolean,
  onClose: () => void,
  onSave: lte_gateway => void,
|};

export default function AddGatewayDialog(props: Props) {
  const classes = useStyles();

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [hardwareID, setHardwareID] = useState('');
  const [gatewayID, setGatewayID] = useState('');
  const [challengeKey, setChallengeKey] = useState('');

  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const networkID = nullthrows(match.params.networkId);

  const onSave = async () => {
    if (!name || !description || !hardwareID || !gatewayID || !challengeKey) {
      enqueueSnackbar('Please complete all fields', {variant: 'error'});
      return;
    }

    try {
      await MagmaV1API.postLteByNetworkIdGateways({
        networkId: networkID,
        gateway: {
          id: gatewayID,
          name,
          description,
          cellular: {
            epc: {nat_enabled: false, ip_block: '192.168.0.1/24'},
            ran: {pci: 260, transmit_enabled: false},
            non_eps_service: undefined,
          },
          magmad: {
            autoupgrade_enabled: true,
            autoupgrade_poll_interval: 300,
            checkin_interval: 60,
            checkin_timeout: 10,
          },
          device: {
            hardware_id: hardwareID,
            key: {
              key: challengeKey,
              key_type: 'SOFTWARE_ECDSA_SHA256', // default key/challenge type
            },
          },
          connected_enodeb_serials: [],
          tier: 'default',
        },
      });
      const gateway = await MagmaV1API.getLteByNetworkIdGatewaysByGatewayId({
        networkId: networkID,
        gatewayId: gatewayID,
      });
      props.onSave(gateway);
    } catch (e) {
      enqueueSnackbar(e?.response?.data?.message || e?.message || e, {
        variant: 'error',
      });
    }
  };

  return (
    <Dialog open={props.open} onClose={props.onClose}>
      <DialogTitle>Add Gateway</DialogTitle>
      <DialogContent>
        <TextField
          label="Gateway Name"
          className={classes.input}
          value={name}
          onChange={({target}) => setName(target.value)}
          placeholder="Gateway 1"
        />
        <TextField
          label="Gateway Description"
          className={classes.input}
          value={description}
          onChange={({target}) => setDescription(target.value)}
          placeholder="Sample Gateway description"
        />
        <TextField
          label="Hardware UUID"
          className={classes.input}
          value={hardwareID}
          onChange={({target}) => setHardwareID(target.value)}
          placeholder="Eg. 4dfe212f-df33-4cd2-910c-41892a042fee"
        />
        <TextField
          label="Gateway ID"
          className={classes.input}
          value={gatewayID}
          onChange={({target}) => setGatewayID(target.value)}
          placeholder="<country>_<org>_<location>_<sitenumber>"
        />
        <TextField
          label="Challenge Key"
          className={classes.input}
          value={challengeKey}
          onChange={({target}) => setChallengeKey(target.value)}
          placeholder="A base64 bytestring of the key in DER format"
        />
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
