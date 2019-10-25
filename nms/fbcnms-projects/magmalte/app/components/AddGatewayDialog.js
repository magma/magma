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
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFillerBackdrop from '@fbcnms/ui/components/LoadingFillerBackdrop';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useState} from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../common/useMagmaAPI';
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

type GatewayData = {
  gatewayID: string,
  name: string,
  description: string,
  hardwareID: string,
  challengeKey: string,
  tier: string,
};

export const MAGMAD_DEFAULT_CONFIGS = {
  autoupgrade_enabled: true,
  autoupgrade_poll_interval: 300,
  checkin_interval: 60,
  checkin_timeout: 10,
};

type Props = {|
  onClose: () => void,
  onSave: GatewayData => Promise<void>,
|};

export default function AddGatewayDialog(props: Props) {
  const classes = useStyles();

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [hardwareID, setHardwareID] = useState('');
  const [gatewayID, setGatewayID] = useState('');
  const [challengeKey, setChallengeKey] = useState('');
  const [tier, setTier] = useState('');

  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const networkID = nullthrows(match.params.networkId);
  const {response: tiers, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdTiers,
    {networkId: networkID},
  );

  if (isLoading || !tiers) {
    return <LoadingFillerBackdrop />;
  }

  const onSave = async () => {
    if (!name || !description || !hardwareID || !gatewayID || !challengeKey) {
      enqueueSnackbar('Please complete all fields', {variant: 'error'});
      return;
    }

    try {
      await props.onSave({
        gatewayID,
        name,
        description,
        hardwareID,
        challengeKey,
        tier,
      });
    } catch (e) {
      enqueueSnackbar(e?.response?.data?.message || e?.message || e, {
        variant: 'error',
      });
    }
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
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
        <FormControl className={classes.input}>
          <InputLabel htmlFor="types">Upgrade Tier</InputLabel>
          <Select
            className={classes.input}
            value={tier}
            onChange={({target}) => setTier(target.value)}>
            {tiers.map(tier => (
              <MenuItem key={tier} value={tier}>
                {tier}
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
