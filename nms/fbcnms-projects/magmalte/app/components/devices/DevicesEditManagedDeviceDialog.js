/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DevicesManagedDevice} from './DevicesUtils';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import FormLabel from '@material-ui/core/FormLabel';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import {MagmaAPIUrls} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useEffect, useState} from 'react';

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
  title: string,
  onCancel: () => void,
  onSave: string => void,
};

export default function DevicesEditManagedDeviceDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();

  const {deviceID: initialDeviceID} = match.params;

  const [deviceID, setDeviceID] = useState(initialDeviceID);
  const [error, setError] = useState('');

  // sectioned device configs
  const [typeTextbox, setTypeTextbox] = useState('cisco');
  const [hostTextbox, setHostTextbox] = useState('0.0.0.0');
  const [platformTextbox, setPlatformTextbox] = useState('Cisco');
  const [frinxTextbox, setFrinxTextbox] = useState(
    JSON.stringify(
      {
        authorization: 'Basic YWRtaW46YWRtaW4=', // admin:admin
        device_type: 'ios',
        device_version: '15.2',
        frinx_port: 8181,
        host: 'frinx',
        password: 'frinx',
        port: 23,
        transport_type: 'telnet',
        username: 'username',
      },
      null,
      2,
    ),
  );
  const [snmpTextbox, setSnmpTextbox] = useState(
    JSON.stringify(
      {
        community: 'public',
        version: 'v1',
      },
      null,
      2,
    ),
  );
  const [cambiumTextbox, setCambiumTextbox] = useState(
    JSON.stringify(
      {
        client_id: 'Ya4dNAAYSUoFMUSs',
        client_ip: '10.0.0.1',
        client_mac: '58:C1:7A:90:36:50',
        client_secret: '2kBgGnMr69NgNdeHB3s7x4GzUODLkc',
      },
      null,
      2,
    ),
  );
  const [otherChannelTextbox, setOtherChannelTextbox] = useState(
    JSON.stringify({}, null, 2),
  );
  const [configTextbox, setConfigTextbox] = useState('{}'); // device configs

  const genManagedDeviceConfig = (): DevicesManagedDevice => {
    const config: DevicesManagedDevice = {
      device_type: typeTextbox.length > 0 ? [typeTextbox] : [],
      host: hostTextbox,
      platform: platformTextbox,
      channels: {},
    };

    function setConfigIfExist(fieldName, textbox, setter) {
      if (!textbox || textbox.length === 0) {
        return;
      }
      try {
        setter(JSON.parse(textbox));
      } catch (err) {
        throw {message: `${fieldName}: ${err.message}`};
      }
    }

    setConfigIfExist('Frinx Channel', frinxTextbox, jsonValue => {
      if (config.channels) {
        config.channels.frinx_channel = jsonValue;
      }
    });

    setConfigIfExist('SNMP Channel', snmpTextbox, jsonValue => {
      if (config.channels) {
        config.channels.snmp_channel = jsonValue;
      }
    });

    setConfigIfExist('Cambium Channel', cambiumTextbox, jsonValue => {
      if (config.channels) {
        config.channels.cambium_channel = jsonValue;
      }
    });

    setConfigIfExist('Other Channel', otherChannelTextbox, jsonValue => {
      if (config.channels) {
        config.channels.other_channel = {
          channel_props: jsonValue,
        };
      }
    });

    setConfigIfExist(
      'Device Config',
      configTextbox,
      jsonValue => (config.device_config = JSON.stringify(jsonValue)),
    );

    return config;
  };

  const onEdit = async () => {
    try {
      const newDeviceConfig = genManagedDeviceConfig();
      await axios.put(MagmaAPIUrls.device(match, deviceID), newDeviceConfig);
      props.onSave(deviceID);
    } catch (error) {
      setError(`${error.message} ${error.response?.data?.message || ''}`);
      return;
    }
  };

  const onCreate = async () => {
    try {
      const newDeviceConfig = genManagedDeviceConfig();
      await axios.post(MagmaAPIUrls.devices(match, deviceID), newDeviceConfig);
      setDeviceID(deviceID);
      props.onSave(deviceID);
    } catch (error) {
      setError(`${error.message} ${error.response?.data?.message || ''}`);
      return;
    }
  };

  const {isLoading, error: responseError, response} = useAxios<
    null,
    DevicesManagedDevice,
  >({
    method: 'get',
    url: MagmaAPIUrls.device(match, initialDeviceID),
  });

  useEffect(() => {
    // TODO: separate out create from edit flow so we don't have garbage axios
    if (initialDeviceID) {
      const initialDeviceConfig = response?.data || {};

      // TODO: support more than 1 device_type in the list
      setTypeTextbox(initialDeviceConfig.device_type?.[0] || '');
      setHostTextbox(initialDeviceConfig.host || '');
      setPlatformTextbox(initialDeviceConfig.platform || '');
      setFrinxTextbox(
        JSON.stringify(initialDeviceConfig.channels?.frinx_channel, null, 2),
      );
      setSnmpTextbox(
        JSON.stringify(initialDeviceConfig.channels?.snmp_channel, null, 2),
      );
      setCambiumTextbox(
        JSON.stringify(initialDeviceConfig.channels?.cambium_channel, null, 2),
      );
      setOtherChannelTextbox(
        JSON.stringify(
          initialDeviceConfig.channels?.other_channel?.channel_props,
          null,
          2,
        ),
      );

      // special case config where config should be JSON.
      // This attempts to format as JSON, failing that, fill in text box as-is.
      try {
        setConfigTextbox(
          JSON.stringify(
            JSON.parse(initialDeviceConfig.device_config || ''),
            null,
            2,
          ),
        );
      } catch (err) {
        setConfigTextbox(initialDeviceConfig.device_config);
      }
    }
  }, [initialDeviceID, response]);

  useEffect(() => {
    // TODO: separate out create from edit flow so we don't have garbage axios
    if (initialDeviceID && responseError) {
      if (responseError.response.status === 404) {
        setError(
          'Warning! Missing config - please enter a new device config below.',
        );
      } else {
        setError(responseError.message);
      }
    }
  }, [initialDeviceID, responseError]);

  // if initialDeviceID is set, then don't allow modifying deviceID
  const deviceIDcontent = initialDeviceID ? (
    <>
      <FormLabel>Device ID: </FormLabel>
      <div>{initialDeviceID}</div>
    </>
  ) : (
    <TextField
      required
      className={classes.input}
      label="Device ID"
      margin="normal"
      onChange={({target}) => {
        setDeviceID(target.value);
      }}
      value={deviceID}
    />
  );

  const content = (
    <FormGroup>
      {deviceIDcontent}

      <TextField
        required
        label="Host IP"
        className={classes.input}
        onChange={({target}) => setHostTextbox(target.value)}
        value={hostTextbox}
      />

      <TextField
        required
        label="Platform"
        className={classes.input}
        onChange={({target}) => setPlatformTextbox(target.value)}
        value={platformTextbox}
      />

      <TextField
        label="Device Type"
        className={classes.input}
        onChange={({target}) => setTypeTextbox(target.value)}
        value={typeTextbox}
      />

      <TextField
        label="Frinx Channel Config"
        className={classes.input}
        multiline={true}
        lines={Infinity}
        onChange={({target}) => setFrinxTextbox(target.value)}
        value={frinxTextbox}
      />

      <TextField
        label="SNMP Channel Config"
        className={classes.input}
        multiline={true}
        lines={Infinity}
        onChange={({target}) => setSnmpTextbox(target.value)}
        value={snmpTextbox}
      />

      <TextField
        label="Cambium Channel Config"
        className={classes.input}
        multiline={true}
        lines={Infinity}
        onChange={({target}) => setCambiumTextbox(target.value)}
        value={cambiumTextbox}
      />

      <TextField
        label="Other Channel Config Props"
        className={classes.input}
        multiline={true}
        lines={Infinity}
        onChange={({target}) => setOtherChannelTextbox(target.value)}
        value={otherChannelTextbox}
      />

      <TextField
        label="Device Config"
        className={classes.input}
        multiline={true}
        lines={Infinity}
        onChange={({target}) => setConfigTextbox(target.value)}
        value={configTextbox}
      />
    </FormGroup>
  );

  return (
    <Dialog open={true} onClose={props.onCancel} fullWidth={true} scroll="body">
      <DialogTitle>{props.title}</DialogTitle>
      <DialogContent>
        {error ? <FormLabel error>{error}</FormLabel> : null}
        {isLoading ? <LoadingFiller /> : content}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onCancel} color="primary">
          Cancel
        </Button>
        <Button
          onClick={initialDeviceID ? onEdit : onCreate}
          color="primary"
          variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
