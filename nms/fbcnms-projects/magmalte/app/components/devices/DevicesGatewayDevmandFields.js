/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DevicesGateway} from './DevicesUtils';

import Button from '@material-ui/core/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormLabel from '@material-ui/core/FormLabel';
import ListFields from '@fbcnms/magmalte/app/components/ListFields';
import React from 'react';
import axios from 'axios';
import {MagmaAPIUrls} from '@fbcnms/magmalte/app/common/MagmaAPI';

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
});

type Props = {
  onClose: () => void,
  onSave: (gatewayID: string) => void,
  gateway: DevicesGateway,
  devmandManagedDevices: Array<string>,
};

export default function DevicesGatewayDevmandFields(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();

  const [errorMessage, setErrorMessage] = useState('');
  const [devmandManagedDevices, setDevmandManagedDevices] = useState<string[]>(
    props.devmandManagedDevices,
  );

  const onSave = () => {
    axios
      .put(MagmaAPIUrls.devicesDevmandConfigs(match, props.gateway.id), {
        managed_devices: devmandManagedDevices.filter(
          device => device.length > 0,
        ),
      })
      .then(() => props.onSave(props.gateway.id))
      .catch(err => {
        setErrorMessage(
          err.toString() + ' ' + err.response?.data?.message || '',
        );
      });
  };

  return (
    <>
      <DialogContent>
        <ListFields
          label="Devmand Managed Devices"
          className={classes.input}
          itemList={devmandManagedDevices}
          onChange={setDevmandManagedDevices}
        />
        <FormLabel error>{errorMessage}</FormLabel>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button onClick={onSave} color="primary" variant="contained">
          Save
        </Button>
      </DialogActions>
    </>
  );
}
