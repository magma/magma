/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {DevicesAgent} from './DevicesUtils';

import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormLabel from '@material-ui/core/FormLabel';
import ListFields from '@fbcnms/magmalte/app/components/ListFields';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
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
  onSave: (agentID: string) => void,
  agent: DevicesAgent,
  devmandManagedDevices: Array<string>,
};

export default function DevicesAgentFields(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();

  const [errorMessage, setErrorMessage] = useState('');
  const [
    devmandManagedDevices,
    setDevmandManagedDevices,
  ] = useState<?(string[])>(null);

  const onSave = () => {
    // distinguishes null (not modified) vs empty (devices were removed)
    if (devmandManagedDevices != null) {
      MagmaV1API.putSymphonyByNetworkIdAgentsByAgentIdManagedDevices({
        networkId: match.params.networkId,
        agentId: props.agent.id,
        managedDevices: devmandManagedDevices.filter(
          device => device.length > 0,
        ),
      })
        .then(() => props.onSave(props.agent.id))
        .catch(err => {
          setErrorMessage(
            err.toString() + ' ' + err.response?.data?.message || '',
          );
        });
    } else {
      () => props.onSave(props.agent.id);
    }
  };

  return (
    <>
      <DialogContent>
        <ListFields
          label="Devmand Managed Devices"
          className={classes.input}
          itemList={devmandManagedDevices || props.devmandManagedDevices}
          onChange={setDevmandManagedDevices}
        />
        <FormLabel error>{errorMessage}</FormLabel>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </>
  );
}
